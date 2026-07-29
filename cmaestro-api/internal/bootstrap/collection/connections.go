package collection

import (
	"cmaestro-api/internal/config"
	"cmaestro-db/bucket"
	"cmaestro-db/models"
	"cmaestro-db/sql"
	"context"
	"errors"
	"fmt"
	"log"
)

type Connections struct {
	SQL           *sql.Client
	ArtifactStore *bucket.Client

	closeFns []func() error
}

func NewCollections(ctx context.Context, runtimeConfig *config.RuntimeConfig) (*Connections, error) {
	conn := &Connections{}

	success := false
	defer func() {
		if !success {
			_ = conn.Close()
		}
	}()

	log.Println("creating [ArtifactStore] connection...")
	artifactStore, err := initArtifactStore(ctx, runtimeConfig.Artifact)
	if err != nil {
		return nil, err
	}
	conn.ArtifactStore = artifactStore
	conn.addCloser("ArtifactStore", artifactStore)

	log.Println("creating [SQL] connection...")
	sqlClient, err := initSQL(runtimeConfig.SQL)
	if err != nil {
		return nil, err
	}

	conn.SQL = sqlClient
	conn.addCloser("SQL Database", sqlClient)

	success = true
	log.Println("Connections created!")
	return conn, nil
}

func initArtifactStore(
	ctx context.Context,
	cfg config.ArtifactRuntimeConfig,
) (*bucket.Client, error) {
	client, err := bucket.New(ctx, bucket.Config{
		Endpoint:   cfg.Endpoint,
		AccessKey:  cfg.AccessKey,
		SecretKey:  cfg.SecretKey,
		Bucket:     cfg.BucketName,
		RootPrefix: cfg.RootPrefix,
	})
	if err != nil {
		return nil, fmt.Errorf("create artifact storage client: %w", err)
	}

	if err := client.EnsureInitialRoot(ctx); err != nil {
		closeIfPossible(client)
		return nil, fmt.Errorf("initialize artifact storage: %w", err)
	}

	return client, nil
}

func initSQL(cfg config.SQLRuntimeConfig) (*sql.Client, error) {
	client, err := sql.New(sql.Config{
		Endpoint: cfg.Endpoint,
		Username: cfg.Username,
		Password: cfg.Password,
		Database: cfg.Database,
		Port:     cfg.Port,
	})
	if err != nil {
		return nil, fmt.Errorf("create SQL client: %w", err)
	}

	err = client.AutoMigrate(
		&models.Artifact{},
	)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (conn *Connections) addCloser(eventName string, value any) {
	log.Printf("adding [%s] to pools.event.close", eventName)
	closer, ok := value.(interface {
		Close() error
	})
	if ok {
		conn.closeFns = append(conn.closeFns, closer.Close)
	}
}

func closeIfPossible(value any) {
	closer, ok := value.(interface {
		Close() error
	})
	if ok {
		_ = closer.Close()
	}
}

func (conn *Connections) Close() error {
	log.Println("closing connections...")
	if conn == nil {
		return nil
	}

	var errs []error

	for i := len(conn.closeFns) - 1; i >= 0; i-- {
		if err := conn.closeFns[i](); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}
