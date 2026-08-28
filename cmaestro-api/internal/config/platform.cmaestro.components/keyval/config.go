package keyval

type Config struct {
	EndpointEnvironemtnVariable string
	UserNameEnvironmentVariable string
	PasswordEnvironmentVariable string
}

func Load() *Config {
	return &Config{
		EndpointEnvironemtnVariable: "CMAESTRO_KEYVAL_ENDPOINT",
		UserNameEnvironmentVariable: "CMAESTRO_KEYVAL_USERNAME",
		PasswordEnvironmentVariable: "CMAESTRO_KEYVAL_PASSWORD",
	}
}
