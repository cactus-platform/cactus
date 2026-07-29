package collection

import (
	"cmaestro-api/internal/config"
	"cmaestro-api/internal/pkg/services"
	"context"
)

type Services struct {
	Artifact services.ArtifactService
}

func NewServices(ctx context.Context, runtimeConfig *config.RuntimeConfig, conns *Connections, repos *Repositories) (*Services, error) {
	return &Services{
		Artifact: services.NewArtifactService(repos.Artifact),
	}, nil
}
