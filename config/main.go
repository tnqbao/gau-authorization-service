package config

import "fmt"

type Config struct {
	EnvConfig *EnvConfig `json:"env_config"`
}

func InitConfig() (*Config, error) {
	envConfig, err := LoadEnvConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load environment configuration: %w", err)
	}

	return &Config{
		EnvConfig: envConfig,
	}, nil
}
