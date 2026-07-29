package collection

import (
	"cmaestro-api/internal/config"
	"context"
)

type Services struct {
}

func NewServices(ctx context.Context, runtimeConfig *config.RuntimeConfig, conns *Connections, repos *Repositories) (*Services, error) {
	return &Services{}, nil
}
