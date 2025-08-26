package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type EnvConfig struct {
	Postgres struct {
		HOST     string
		Database string
		Username string
		Password string
		Port     string
	}
	JWT struct {
		SecretKey string
		Algorithm string
		Expire    int
	}
	CORS struct {
		AllowDomains string
		GlobalDomain string
	}
	Redis struct {
		Address  string
		Password string
		Database int
	}

	PrivateKey string

	Grafana struct {
		OTLPEndpoint string
		ServiceName  string
	}

	Environment struct {
		Mode  string
		Group string
	}
}

func LoadEnvConfig() (*EnvConfig, error) {
	var config EnvConfig
	var missing []string

	// Postgres
	config.Postgres.HOST = os.Getenv("PGPOOL_HOST")
	config.Postgres.Database = os.Getenv("PGPOOL_DB")
	config.Postgres.Username = os.Getenv("PGPOOL_USER")
	config.Postgres.Password = os.Getenv("PGPOOL_PASSWORD")
	config.Postgres.Port = os.Getenv("PGPOOL_PORT")

	// Check required Postgres envs
	if config.Postgres.HOST == "" {
		missing = append(missing, "PGPOOL_HOST")
	}
	if config.Postgres.Database == "" {
		missing = append(missing, "PGPOOL_DB")
	}
	if config.Postgres.Username == "" {
		missing = append(missing, "PGPOOL_USER")
	}
	if config.Postgres.Password == "" {
		missing = append(missing, "PGPOOL_PASSWORD")
	}
	if config.Postgres.Port == "" {
		missing = append(missing, "PGPOOL_PORT")
	}

	// JWT
	config.JWT.SecretKey = os.Getenv("JWT_SECRET_KEY")
	config.JWT.Algorithm = os.Getenv("JWT_ALGORITHM")
	if config.JWT.SecretKey == "" {
		missing = append(missing, "JWT_SECRET_KEY")
	}
	if config.JWT.Algorithm == "" {
		config.JWT.Algorithm = "HS256" // default fallback
	}

	if val := os.Getenv("JWT_EXPIRE"); val != "" {
		if _, err := fmt.Sscanf(val, "%d", &config.JWT.Expire); err != nil {
			return nil, fmt.Errorf("invalid JWT_EXPIRE: %w", err)
		}
	} else {
		config.JWT.Expire = 3600 * 24 * 7 // Default to 7 days
	}

	// CORS
	config.CORS.AllowDomains = os.Getenv("ALLOWED_DOMAINS")
	config.CORS.GlobalDomain = os.Getenv("GLOBAL_DOMAIN")

	config.PrivateKey = os.Getenv("PRIVATE_KEY")

	// Redis
	config.Redis.Address = os.Getenv("REDIS_ADDRESS")
	config.Redis.Password = os.Getenv("REDIS_PASSWORD")
	if dbStr := os.Getenv("REDIS_DB"); dbStr != "" {
		db, err := strconv.Atoi(dbStr)
		if err != nil {
			return nil, fmt.Errorf("invalid REDIS_DB: %w", err)
		}
		config.Redis.Database = db
	} else {
		config.Redis.Database = 0
	}
	if config.Redis.Address == "" {
		missing = append(missing, "REDIS_ADDRESS")
	}

	// Grafana/OpenTelemetry
	grafanaEndpoint := os.Getenv("GRAFANA_OTLP_ENDPOINT")
	if grafanaEndpoint == "" {
		grafanaEndpoint = "https://grafana.gauas.online"
	}
	// Remove protocol for OpenTelemetry client to avoid duplicate protocols
	if strings.HasPrefix(grafanaEndpoint, "https://") {
		config.Grafana.OTLPEndpoint = strings.TrimPrefix(grafanaEndpoint, "https://")
	} else if strings.HasPrefix(grafanaEndpoint, "http://") {
		config.Grafana.OTLPEndpoint = strings.TrimPrefix(grafanaEndpoint, "http://")
	} else {
		config.Grafana.OTLPEndpoint = grafanaEndpoint
	}
	config.Grafana.ServiceName = os.Getenv("SERVICE_NAME")
	if config.Grafana.ServiceName == "" {
		config.Grafana.ServiceName = "gau-account-service"
	}

	config.Environment.Mode = os.Getenv("DEPLOY_ENV")
	if config.Environment.Mode == "" {
		config.Environment.Mode = "development"
	}

	config.Environment.Group = os.Getenv("GROUP_NAME")
	if config.Environment.Group == "" {
		config.Environment.Group = "local"
	}
	if len(missing) > 0 {
		return nil, fmt.Errorf("missing required env vars: %v", strings.Join(missing, ", "))
	}

	return &config, nil
}
