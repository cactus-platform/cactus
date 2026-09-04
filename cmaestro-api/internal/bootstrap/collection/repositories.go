package collection

import (
	"cmaestro-api/internal/config"
	"context"

	coreRepositories "github.com/cactus-platform/cmaestro-core/repositories"
)

type Repositories struct {
	Repository coreRepositories.RepositoryRepository
}

func NewRepositories(ctx context.Context, runtimeConfig *config.RuntimeConfig, conns *Connections) (*Repositories, error) {
	return &Repositories{
		Repository: coreRepositories.NewRepositoryRepository(conns.SQL),
	}, nil
}
