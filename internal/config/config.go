package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port           string
	DataSourceName string
	Copium         CopiumConfig
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

	copium, err := loadCopiumConfig()
	if err != nil {
		return nil, err
	}
	cfg.Copium = copium

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}
