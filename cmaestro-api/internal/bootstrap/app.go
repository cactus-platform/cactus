package bootstrap

import (
	"cmaestro-api/internal/bootstrap/collection"
	"cmaestro-api/internal/config"
	"context"
	"errors"
	"fmt"
	"log"
)

type App struct {
	StaticConfig  *config.StaticConfig
	RuntimeConfig *config.RuntimeConfig

	Connections  *collection.Connections
	Repositories *collection.Repositories
	Services     *collection.Services

	closeFns []func() error
}

func NewFromEnv(ctx context.Context) (*App, error) {
	staticConfig := config.Load()

	runtimeConfig, err := staticConfig.ResolveRuntime()
	if err != nil {
		return nil, fmt.Errorf("resolve configuration: %w", err)
	}

	return New(ctx, staticConfig, runtimeConfig)
}

// New creates all application dependencies.
// It uses a local app variable rather than a named return value so deferred
// cleanup never attempts to call Close on nil.
func New(ctx context.Context, staticCfg *config.StaticConfig, runtimeCfg *config.RuntimeConfig) (*App, error) {
	log.Println("initializing application")
	if ctx == nil {
		return nil, errors.New("context cannot be nil")
	}

	app := &App{
		StaticConfig:  staticCfg,
		RuntimeConfig: runtimeCfg,
	}

	success := false
	defer func() {
		if !success {
			_ = app.Close()
		}
	}()

	conn, err := collection.NewCollections(ctx, runtimeCfg)
	if err != nil {
		return nil, fmt.Errorf("create collections: %w", err)
	}
	log.Printf("Collections: ok")
	app.Connections = conn

	repos, err := collection.NewRepositories(ctx, runtimeCfg, app.Connections)
	if err != nil {
		return nil, fmt.Errorf("create repositories: %w", err)
	}
	log.Printf("Repositories: ok")
	app.Repositories = repos

	svcs, err := collection.NewServices(ctx, runtimeCfg, app.Connections, app.Repositories)
	if err != nil {
		return nil, fmt.Errorf("create services: %w", err)
	}
	log.Printf("Services: ok")
	app.Services = svcs

	success = true
	log.Printf("Application created!")
	return app, nil
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
