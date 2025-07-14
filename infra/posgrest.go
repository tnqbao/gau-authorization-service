package infra

import (
	"fmt"
	"github.com/tnqbao/gau-authorization-service/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"time"
)

type PostgresClient struct {
	DB *gorm.DB
}

func InitPostgresClient(cfg *config.EnvConfig) (*PostgresClient, error) {
	pgUser := cfg.Postgres.Username
	pgPassword := cfg.Postgres.Password
	pgHost := cfg.Postgres.HOST
	pgDB := cfg.Postgres.Database
	pgPort := cfg.Postgres.Port

	if pgUser == "" || pgPassword == "" || pgHost == "" || pgDB == "" || pgPort == "" {
		return nil, fmt.Errorf("one or more required Postgres configs are missing")
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Ho_Chi_Minh",
		pgHost, pgUser, pgPassword, pgDB, pgPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Postgres: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get raw DB from GORM: %w", err)
	}

	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("PostgreSQL connected at", pgHost)

	return &PostgresClient{DB: db}, nil
}
