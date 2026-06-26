package config

import (
	platform_cmaestro_components "cmaestro-api/internal/config/platform.cmaestro.components"
	"cmaestro-api/internal/config/platform.cmaestro/repositories"
)

type Config struct {
	Repositories *repositories.Config
	_Components  *platform_cmaestro_components.Config
}

func Load() *Config {
	return &Config{
		Repositories: repositories.Load(),
		_Components:  platform_cmaestro_components.Load(),
	}
}
