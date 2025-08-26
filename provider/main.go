package provider

import (
	"github.com/tnqbao/gau-authorization-service/config"
)

type Provider struct {
	LoggerProvider *LoggerProvider
}

var provider *Provider

func InitProvider(cfg *config.EnvConfig) *Provider {
	loggerProvider := NewLoggerProvider()
	provider = &Provider{
		LoggerProvider: loggerProvider,
	}

	return provider
}

func GetProvider() *Provider {
	if provider == nil {
		panic("Provider not initialized")
	}
	return provider
}
