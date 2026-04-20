package config

import (
	"log"
	"os"
)

type Config struct {
	HTTPAddr      string
	PostgresURL   string
	AppEnv        string
	MigrationsDir string
	SeedsDir      string
}

func MustLoad() Config {
	cfg := Config{
		HTTPAddr:      getenv("HTTP_ADDR", ":8080"),
		PostgresURL:   getenv("POSTGRES_URL", "postgres://cee4ai:cee4ai@localhost:5432/cee4ai?sslmode=disable"),
		AppEnv:        getenv("APP_ENV", "development"),
		MigrationsDir: getenv("MIGRATIONS_DIR", "./migrations"),
		SeedsDir:      getenv("SEEDS_DIR", "./seeds"),
	}

	if cfg.PostgresURL == "" {
		log.Fatal("POSTGRES_URL is required")
	}

	return cfg
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}
