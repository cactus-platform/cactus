package config

type AppContext struct {
	Config *Config
}

func NewAppContext() *AppContext {
	return &AppContext{
		Config: Load(),
	}
}
