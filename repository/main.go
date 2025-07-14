package repository

import (
	"errors"

	"github.com/redis/go-redis/v9"
	"github.com/tnqbao/gau-authorization-service/infra"
	"gorm.io/gorm"
)

type Repository struct {
	DB      *gorm.DB
	CacheDB *redis.Client
}

// InitRepository initializes a new Repository instance using infrastructure dependencies.
// Returns an error if Postgres or Redis are not properly initialized.
func InitRepository(infra *infra.Infra) (*Repository, error) {
	if infra.Postgres == nil || infra.Postgres.DB == nil {
		return nil, errors.New("Postgres client is not initialized")
	}
	if infra.Redis == nil || infra.Redis.Client == nil {
		return nil, errors.New("Redis client is not initialized")
	}

	return &Repository{
		DB:      infra.Postgres.DB,
		CacheDB: infra.Redis.Client,
	}, nil
}
