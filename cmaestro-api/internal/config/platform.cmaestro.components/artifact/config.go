package artifact

type Config struct {
	BucketName                   string
	RootPrefix                   string
	EndpointEnvironmentVariable  string
	AccessKeyEnvironmentVariable string
	SecretKeyEnvironmentVariable string
}

func Load() *Config {
	return &Config{
		BucketName:                   "platform.cmaestro",
		RootPrefix:                   "artifacts",
		EndpointEnvironmentVariable:  "CMAESTRO_ARTIFACT_ENDPOINT",
		AccessKeyEnvironmentVariable: "CMAESTRO_ARTIFACT_ACCESS_KEY",
		SecretKeyEnvironmentVariable: "CMAESTRO_ARTIFACT_SECRET_KEY",
	}
}
