package infra

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/tnqbao/gau-authorization-service/config"
	"log"
)

type RedisClient struct {
	Client *redis.Client
}

func InitRedisClient(cfg *config.EnvConfig) (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Address,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.Database,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	log.Println("Connected to Redis:", cfg.Redis.Address)

	return &RedisClient{Client: client}, nil
}

func (r *RedisClient) Set(key string, value string) error {
	return r.Client.Set(context.Background(), key, value, 0).Err()
}

func (r *RedisClient) Get(key string) (string, error) {
	return r.Client.Get(context.Background(), key).Result()
}

func (r *RedisClient) Delete(key string) error {
	return r.Client.Del(context.Background(), key).Err()
}
