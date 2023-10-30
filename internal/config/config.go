package config

type Config struct {
	ExperimentalAirSupport bool
}

var (
	// store config internally
	config = &Config{}
)

func GetConfig() *Config {
	return config
}
