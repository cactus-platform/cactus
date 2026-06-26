package config

import (
	"context"
	"fmt"
	"os"

	"cmaestro-db/bucket"
)

type AppContext struct {
	Context    context.Context
	Config     *Config
	ArtifactDB *bucket.Client
	//InstantDB  *keyval.Client
}

func NewAppContext(ctx context.Context) (*AppContext, error) {
	cfg := Load()
	endpoint := os.Getenv(cfg._Components.Artifact.EndpointEnvironmentVariable)
	accessKey := os.Getenv(cfg._Components.Artifact.AccessKeyEnvironmentVariable)
	secretKey := os.Getenv(cfg._Components.Artifact.SecretKeyEnvironmentVariable)

	if endpoint == "" || accessKey == "" || secretKey == "" {
		return nil, fmt.Errorf("CMAESTRO_ARTIFACT env variables not defined or empty"+
			", "+
			"%s.size=[%dB]"+
			", "+
			"%s.size=[%dB]"+
			", "+
			"%s.size=[%dB]",
			cfg._Components.Artifact.EndpointEnvironmentVariable, len(endpoint),
			cfg._Components.Artifact.AccessKeyEnvironmentVariable, len(accessKey),
			cfg._Components.Artifact.SecretKeyEnvironmentVariable, len(secretKey))
	}

	artifactDB, err := bucket.New(ctx, bucket.Config{
		Endpoint:   endpoint,
		AccessKey:  accessKey,
		SecretKey:  secretKey,
		Bucket:     cfg._Components.Artifact.BucketName,
		RootPrefix: cfg._Components.Artifact.RootPrefix,
	})
	if err != nil {
		return nil, fmt.Errorf("create SeaweedFS artifact client: %w", err)
	}

	// Creates bucket "cmaestro" if needed and creates root/.keep.
	if err := artifactDB.EnsureInitialRoot(ctx); err != nil {
		return nil, fmt.Errorf("initialize SeaweedFS storage: %w", err)
	}

	return &AppContext{
		Context:    ctx,
		Config:     cfg,
		ArtifactDB: artifactDB,
		// InstantDB
	}, nil
}
