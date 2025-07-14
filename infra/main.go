package infra

import (
	"fmt"
	"github.com/tnqbao/gau-authorization-service/config"
	"sync"
)

type Infra struct {
	Redis    *RedisClient
	Postgres *PostgresClient
}

var infraInstance *Infra

var infraOnce sync.Once
var infraErr error

func InitInfra(cfg *config.Config) (*Infra, error) {
	infraOnce.Do(func() {
		redis, err := InitRedisClient(cfg.EnvConfig)
		if err != nil {
			infraErr = fmt.Errorf("failed to initialize Redis: %w", err)
			return
		}

		postgres, err := InitPostgresClient(cfg.EnvConfig)
		if err != nil {
			infraErr = fmt.Errorf("failed to initialize Postgres: %w", err)
			return
		}

		infraInstance = &Infra{
			Redis:    redis,
			Postgres: postgres,
		}
	})

	return infraInstance, infraErr
}

func GetClient() *Infra {
	if infraInstance == nil {
		panic("Infra not initialized. Call InitInfra() first.")
	}
	return infraInstance
}
