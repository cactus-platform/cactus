package platform_cmaestro_components

import (
	"cmaestro-api/internal/config/platform.cmaestro.components/artifact"
)

type Config struct {
	Artifact *artifact.Config
}

func Load() *Config {
	return &Config{
		Artifact: artifact.Load(),
	}
}
