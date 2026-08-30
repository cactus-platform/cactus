package collection

import (
	"cmaestro-api/internal/config"
	"cmaestro-api/internal/pkg/services"
	"context"
)

type Services struct {
	Artifact services.ArtifactService
	Ingest   services.IngestService
}

func NewServices(ctx context.Context, runtimeConfig *config.RuntimeConfig, conns *Connections, repos *Repositories) (*Services, error) {
	return &Services{
		Artifact: services.NewArtifactService(repos.Artifact),
		Ingest:   services.NewIngestService(conns.KeyVal),
	}, nil
}
