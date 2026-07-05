package config

import (
	platform_cmaestro_components "cmaestro-api/internal/config/platform.cmaestro.components"
	"cmaestro-api/internal/config/platform.cmaestro/repositories"
)

type StaticConfig struct {
	Repositories *repositories.Config
	Components   *platform_cmaestro_components.Config
}

func Load() *StaticConfig {
	return &StaticConfig{
		Repositories: repositories.Load(),
		Components:   platform_cmaestro_components.Load(),
	}
}
