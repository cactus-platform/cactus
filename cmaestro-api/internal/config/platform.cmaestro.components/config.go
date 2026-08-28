package platform_cmaestro_components

import (
	"cmaestro-api/internal/config/platform.cmaestro.components/artifact"
	"cmaestro-api/internal/config/platform.cmaestro.components/keyval"
	"cmaestro-api/internal/config/platform.cmaestro.components/sql"
)

type Config struct {
	Artifact *artifact.Config
	SQL      *sql.Config
	KeyVal   *keyval.Config
}

func Load() *Config {
	return &Config{
		Artifact: artifact.Load(),
		SQL:      sql.Load(),
		KeyVal:   keyval.Load(),
	}
}
