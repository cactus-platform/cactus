package config

import "cmaestro-api/internal/config/repositories"

type Config struct {
	Repositories *repositories.Config
}

func Load() *Config {
	return &Config{
		Repositories: repositories.Load(),
	}
}
