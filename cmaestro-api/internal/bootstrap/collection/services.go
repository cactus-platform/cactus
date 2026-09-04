package collection

import (
	"cmaestro-api/internal/config"
	"context"

	coreServices "github.com/cactus-platform/cmaestro-core/services"
)

type Services struct {
	Repository coreServices.RepositoryService
	Ingest     coreServices.IngestService
}

func NewServices(ctx context.Context, runtimeConfig *config.RuntimeConfig, conns *Connections, repos *Repositories) (*Services, error) {
	return &Services{
		Repository: coreServices.NewRepositoryService(repos.Repository),
		Ingest:     coreServices.NewIngestService(conns.KeyVal),
	}, nil
}
