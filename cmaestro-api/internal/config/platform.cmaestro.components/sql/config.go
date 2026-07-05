package sql

type Config struct {
	DatabaseName                string
	EndpointEnvironmentVariable string
	UsernameEnvironmentVariable string
	PasswordEnvironmentVariable string
}

func Load() *Config {
	return &Config{
		DatabaseName:                "platform.cmaestro",
		EndpointEnvironmentVariable: "CMAESTRO_SQL_ENDPOINT",
		UsernameEnvironmentVariable: "CMAESTRO_SQL_USER",
		PasswordEnvironmentVariable: "CMAESTRO_SQL_PASSWD",
	}
}
