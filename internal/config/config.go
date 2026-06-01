package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port           string
	DataSourceName string
}

func Load() (*Config, error) {
	_ = godotenv.Load() // not fatal; env vars may be set directly

	cfg := &Config{
		Port:           getEnv("PORT", "8080"),
		DataSourceName: getEnv("DATA_SOURCE_NAME", ""),
	}

	if cfg.DataSourceName == "" {
		return nil, fmt.Errorf("DATA_SOURCE_NAME is required")
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}
