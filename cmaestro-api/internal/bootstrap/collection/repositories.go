package collection

import (
	"cmaestro-api/internal/config"
	"cmaestro-db/repositories"
	"context"
)

type Repositories struct {
	Artifact repositories.ArtifactRepository
}

func NewRepositories(ctx context.Context, runtimeConfig *config.RuntimeConfig, conns *Connections) (*Repositories, error) {
	return &Repositories{
		repositories.NewArtifactRepository(conns.SQL),
	}, nil
}
