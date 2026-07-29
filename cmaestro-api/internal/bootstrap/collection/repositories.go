package collection

import (
	"cmaestro-api/internal/config"
	"context"
)

type Repositories struct {
}

func NewRepositories(ctx context.Context, runtimeConfig *config.RuntimeConfig, conns *Connections) (*Repositories, error) {
	return &Repositories{}, nil
}
