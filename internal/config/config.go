package config

import (
	"log/slog"
	"os"
)

type Config struct {
	Port        string
	DatabaseURL string
	Environment string
	LogLevel    slog.Level
}

func New() *Config {
	cfg := &Config{
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", "data.db"),
		Environment: getEnv("ENVIRONMENT", "development"),
	}

	switch getEnv("LOG_LEVEL", "info") {
	case "debug":
		cfg.LogLevel = slog.LevelDebug
	case "warn":
		cfg.LogLevel = slog.LevelWarn
	case "error":
		cfg.LogLevel = slog.LevelError
	default:
		cfg.LogLevel = slog.LevelInfo
	}

	return cfg
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
