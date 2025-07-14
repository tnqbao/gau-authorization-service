package config

type Config struct {
	EnvConfig *EnvConfig `json:"env_config"`
}

func InitConfig() *Config {
	EnvConfig := LoadEnvConfig()
	return &Config{
		EnvConfig: EnvConfig,
	}
}
