// Package seaweedfs provides a small client for SeaweedFS's S3-compatible API.
// It can idempotently initialise a bucket and an optional virtual root prefix,
// then upload ZIP files with the AWS SDK v2 multipart uploader.
package bucket

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awscfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"
)

const (
	defaultRegion      = "us-east-1"
	defaultPartSize    = 16 * 1024 * 1024 // 16 MiB
	defaultConcurrency = 4
)

// Config contains the connection and upload settings.
//
// Endpoint example: http://seaweedfs-s3.cmaestro-db.svc.cluster.local:8333
// RootPrefix example: root or application-data. It is not a real directory;
// it becomes a prefix in object keys.
type Config struct {
	Endpoint    string
	Region      string
	AccessKey   string
	SecretKey   string
	Bucket      string
	RootPrefix  string
	PartSize    int64
	Concurrency int
}

// Client owns an S3 client and a multipart uploader configured for SeaweedFS.
type Client struct {
	bucket string
	root   string
	s3     *s3.Client
	upload *manager.Uploader
}

// Object is the useful result of a successful upload.
type Object struct {
	Bucket    string
	Key       string
	ETag      string
	VersionID string
	Location  string
	Size      int64
}

// New creates a SeaweedFS S3 client. It deliberately uses path-style S3
// requests, which work reliably for in-cluster SeaweedFS service endpoints.
func New(ctx context.Context, cfg Config) (*Client, error) {
	if strings.TrimSpace(cfg.Endpoint) == "" {
		return nil, errors.New("SeaweedFS endpoint is required")
	}
	endpoint, err := url.ParseRequestURI(cfg.Endpoint)
	if err != nil || endpoint.Scheme == "" || endpoint.Host == "" {
		return nil, fmt.Errorf("invalid SeaweedFS endpoint %q", cfg.Endpoint)
	}
	if strings.TrimSpace(cfg.Bucket) == "" {
		return nil, errors.New("SeaweedFS bucket is required")
	}
	if cfg.AccessKey == "" || cfg.SecretKey == "" {
		return nil, errors.New("SeaweedFS access key and secret key are required")
	}

	root, err := normalizePrefix(cfg.RootPrefix)
	if err != nil {
		return nil, err
	}

	region := cfg.Region
	if region == "" {
		region = defaultRegion
	}
	if cfg.PartSize == 0 {
		cfg.PartSize = defaultPartSize
	}
	if cfg.PartSize < 5*1024*1024 {
		return nil, fmt.Errorf("multipart part size must be at least 5 MiB, got %d bytes", cfg.PartSize)
	}
	if cfg.Concurrency == 0 {
		cfg.Concurrency = defaultConcurrency
	}
	if cfg.Concurrency < 1 {
		return nil, fmt.Errorf("upload concurrency must be positive, got %d", cfg.Concurrency)
	}

	awsConfig, err := awscfg.LoadDefaultConfig(
		ctx,
		awscfg.WithRegion(region),
		awscfg.WithRetryMaxAttempts(5),
		awscfg.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(cfg.AccessKey, cfg.SecretKey, ""),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("load AWS SDK config: %w", err)
	}

	s3Client := s3.NewFromConfig(awsConfig, func(options *s3.Options) {
		options.BaseEndpoint = aws.String(strings.TrimRight(cfg.Endpoint, "/"))
		options.UsePathStyle = true
	})

	return &Client{
		bucket: cfg.Bucket,
		root:   root,
		s3:     s3Client,
		upload: manager.NewUploader(s3Client, func(u *manager.Uploader) {
			u.PartSize = cfg.PartSize
			u.Concurrency = cfg.Concurrency
		}),
	}, nil
}

// EnsureInitialRoot prepares first-run storage. It is safe to call from
// multiple application replicas: it creates the bucket when absent and uses
// an empty .keep object only to make an otherwise-empty root prefix visible.
func (c *Client) EnsureInitialRoot(ctx context.Context) error {
	if err := c.ensureBucket(ctx); err != nil {
		return err
	}
	if c.root == "" {
		return nil
	}

	markerKey := c.root + ".keep"
	_, err := c.s3.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(markerKey),
	})
	if err == nil {
		return nil
	}
	if !isNotFound(err) {
		return fmt.Errorf("check root marker s3://%s/%s: %w", c.bucket, markerKey, err)
	}

	_, err = c.s3.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(c.bucket),
		Key:           aws.String(markerKey),
		Body:          bytes.NewReader(nil),
		ContentLength: aws.Int64(0),
		ContentType:   aws.String("application/octet-stream"),
		Metadata: map[string]string{
			"managed-by": "seaweedfs-go-util",
		},
	})
	if err != nil {
		// Another replica may have created the marker after the HeadObject call.
		if _, checkErr := c.s3.HeadObject(ctx, &s3.HeadObjectInput{
			Bucket: aws.String(c.bucket),
			Key:    aws.String(markerKey),
		}); checkErr == nil {
			return nil
		}
		return fmt.Errorf("create root marker s3://%s/%s: %w", c.bucket, markerKey, err)
	}
	return nil
}

func (c *Client) ensureBucket(ctx context.Context) error {
	_, err := c.s3.HeadBucket(ctx, &s3.HeadBucketInput{Bucket: aws.String(c.bucket)})
	if err == nil {
		return nil
	}
	if !isNotFound(err) {
		return fmt.Errorf("check bucket %q: %w", c.bucket, err)
	}

	_, createErr := c.s3.CreateBucket(ctx, &s3.CreateBucketInput{Bucket: aws.String(c.bucket)})
	if createErr == nil {
		return nil
	}

	// A concurrent replica may have created the bucket between HeadBucket and
	// CreateBucket. Re-checking avoids treating other creation errors as success.
	if _, checkErr := c.s3.HeadBucket(ctx, &s3.HeadBucketInput{Bucket: aws.String(c.bucket)}); checkErr == nil {
		return nil
	}
	return fmt.Errorf("create bucket %q: %w", c.bucket, createErr)
}

// UploadZipFile uploads a local ZIP file below RootPrefix. objectName must be
// relative to RootPrefix, e.g. "archives/database-2026-06-26.zip".
func (c *Client) UploadZipFile(ctx context.Context, filePath, objectName string) (Object, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return Object{}, fmt.Errorf("open ZIP %q: %w", filePath, err)
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return Object{}, fmt.Errorf("stat ZIP %q: %w", filePath, err)
	}
	if !info.Mode().IsRegular() {
		return Object{}, fmt.Errorf("ZIP path %q is not a regular file", filePath)
	}
	return c.UploadZip(ctx, file, info.Size(), objectName)
}

// UploadZip uploads ZIP data. For files, prefer UploadZipFile: an *os.File is
// seekable, allowing the SDK to retry upload work more efficiently.
func (c *Client) UploadZip(ctx context.Context, body io.Reader, size int64, objectName string) (Object, error) {
	if body == nil {
		return Object{}, errors.New("ZIP body is required")
	}
	if size < 0 {
		return Object{}, fmt.Errorf("ZIP size cannot be negative: %d", size)
	}

	key, err := c.objectKey(objectName)
	if err != nil {
		return Object{}, err
	}

	result, err := c.upload.Upload(ctx, &s3.PutObjectInput{
		Bucket:             aws.String(c.bucket),
		Key:                aws.String(key),
		Body:               body,
		ContentLength:      aws.Int64(size),
		ContentType:        aws.String(zipContentType(objectName)),
		ContentDisposition: aws.String(fmt.Sprintf("attachment; filename=%q", path.Base(key))),
		Metadata: map[string]string{
			"managed-by":  "seaweedfs-go-util",
			"uploaded-at": time.Now().UTC().Format(time.RFC3339),
		},
	})
	if err != nil {
		return Object{}, fmt.Errorf("upload s3://%s/%s: %w", c.bucket, key, err)
	}

	return Object{
		Bucket:    c.bucket,
		Key:       key,
		ETag:      aws.ToString(result.ETag),
		VersionID: aws.ToString(result.VersionID),
		Location:  result.Location,
		Size:      size,
	}, nil
}

func (c *Client) objectKey(name string) (string, error) {
	name = strings.TrimSpace(strings.ReplaceAll(name, "\\", "/"))
	name = strings.Trim(name, "/")
	if name == "" {
		return "", errors.New("object name is required")
	}
	if strings.HasSuffix(name, "/") {
		return "", fmt.Errorf("object name %q points to a prefix, not a ZIP file", name)
	}
	for _, segment := range strings.Split(name, "/") {
		if segment == "." || segment == ".." || segment == "" {
			return "", fmt.Errorf("object name %q has an invalid path segment", name)
		}
	}
	if !strings.HasSuffix(strings.ToLower(name), ".zip") {
		return "", fmt.Errorf("object name %q must end in .zip", name)
	}
	return c.root + name, nil
}

func normalizePrefix(prefix string) (string, error) {
	prefix = strings.TrimSpace(strings.ReplaceAll(prefix, "\\", "/"))
	prefix = strings.Trim(prefix, "/")
	if prefix == "" {
		return "", nil
	}
	for _, segment := range strings.Split(prefix, "/") {
		if segment == "." || segment == ".." || segment == "" {
			return "", fmt.Errorf("root prefix %q has an invalid path segment", prefix)
		}
	}
	return prefix + "/", nil
}

func zipContentType(name string) string {
	if contentType := mime.TypeByExtension(path.Ext(name)); contentType != "" {
		return contentType
	}
	return "application/zip"
}

func isNotFound(err error) bool {
	var statusCoder interface{ HTTPStatusCode() int }
	if errors.As(err, &statusCoder) && statusCoder.HTTPStatusCode() == http.StatusNotFound {
		return true
	}

	var apiErr smithy.APIError
	if errors.As(err, &apiErr) {
		switch apiErr.ErrorCode() {
		case "NotFound", "NoSuchBucket", "NoSuchKey", "404":
			return true
		}
	}
	return false
}
