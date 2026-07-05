package config

import (
	"context"
	"errors"
	"fmt"
	"log"

	"cmaestro-db/bucket"
	"cmaestro-db/sql"
)

type App struct {
	StaticConfig  StaticConfig
	RuntimeConfig RuntimeConfig
	ArtifactStore *bucket.Client
	SQL           *sql.Client

	closeFns []func() error
}

func NewFromEnv(ctx context.Context) (*App, error) {
	rawConfig := Load()

	runtimeConfig, err := rawConfig.ResolveRuntime()
	if err != nil {
		return nil, fmt.Errorf("resolve configuration: %w", err)
	}

	return New(ctx, runtimeConfig)
}

// New creates all application dependencies.
// It uses a local app variable rather than a named return value so deferred
// cleanup never attempts to call Close on nil.
func New(ctx context.Context, cfg RuntimeConfig) (*App, error) {
	log.Println("initializing application")
	if ctx == nil {
		return nil, errors.New("context cannot be nil")
	}

	app := &App{
		RuntimeConfig: cfg,
	}

	success := false
	defer func() {
		if !success {
			_ = app.Close()
		}
	}()

	log.Println("creating [ArtifactStore] connection...")
	artifactStore, err := initArtifactStore(ctx, cfg.Artifact)
	if err != nil {
		return nil, err
	}

	app.ArtifactStore = artifactStore
	app.addCloser("ArtifactStore", artifactStore)

	log.Println("creating [SQL] connection...")
	sqlClient, err := initSQL(cfg.SQL)
	if err != nil {
		return nil, err
	}

	app.SQL = sqlClient
	app.addCloser("SQL Database", sqlClient)

	success = true
	log.Println("Application created!")
	return app, nil
}

func initArtifactStore(
	ctx context.Context,
	cfg ArtifactRuntimeConfig,
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

func initSQL(cfg SQLRuntimeConfig) (*sql.Client, error) {
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

	return client, nil
}

func (a *App) addCloser(eventName string, value any) {
	log.Printf("adding [%s] to pools.event.close", eventName)
	closer, ok := value.(interface {
		Close() error
	})
	if ok {
		a.closeFns = append(a.closeFns, closer.Close)
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

func (a *App) Close() error {
	log.Println("closing application...")
	if a == nil {
		return nil
	}

	var errs []error

	for i := len(a.closeFns) - 1; i >= 0; i-- {
		if err := a.closeFns[i](); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}
