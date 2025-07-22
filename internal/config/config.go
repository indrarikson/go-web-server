package config

import (
	"log/slog"
	"os"
	"strconv"
	"time"
)

type Config struct {
	// Server configuration
	Port            string
	Host            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration

	// Database configuration
	DatabaseURL            string
	DatabaseMaxConnections int
	DatabaseTimeout        time.Duration
	RunMigrations          bool

	// Application configuration
	Environment string
	Debug       bool
	LogLevel    slog.Level
	LogFormat   string // "json" or "text"

	// Security configuration
	TrustedProxies []string
	EnableCORS     bool
	AllowedOrigins []string

	// Feature flags
	EnableMetrics bool
	EnablePprof   bool
}

func New() *Config {
	cfg := &Config{
		// Server defaults
		Port:            getEnv("PORT", "8080"),
		Host:            getEnv("HOST", ""),
		ReadTimeout:     getDurationEnv("READ_TIMEOUT", 10*time.Second),
		WriteTimeout:    getDurationEnv("WRITE_TIMEOUT", 10*time.Second),
		ShutdownTimeout: getDurationEnv("SHUTDOWN_TIMEOUT", 30*time.Second),

		// Database defaults
		DatabaseURL:            getEnv("DATABASE_URL", "data.db"),
		DatabaseMaxConnections: getIntEnv("DATABASE_MAX_CONNECTIONS", 25),
		DatabaseTimeout:        getDurationEnv("DATABASE_TIMEOUT", 30*time.Second),
		RunMigrations:          getBoolEnv("RUN_MIGRATIONS", true),

		// Application defaults
		Environment: getEnv("ENVIRONMENT", "development"),
		Debug:       getBoolEnv("DEBUG", false),
		LogFormat:   getEnv("LOG_FORMAT", "text"),

		// Security defaults
		TrustedProxies: []string{"127.0.0.1"},
		EnableCORS:     getBoolEnv("ENABLE_CORS", true),
		AllowedOrigins: []string{"*"},

		// Feature flags
		EnableMetrics: getBoolEnv("ENABLE_METRICS", false),
		EnablePprof:   getBoolEnv("ENABLE_PPROF", false),
	}

	// Set log level
	switch getEnv("LOG_LEVEL", "info") {
	case "debug":
		cfg.LogLevel = slog.LevelDebug
		cfg.Debug = true
	case "warn":
		cfg.LogLevel = slog.LevelWarn
	case "error":
		cfg.LogLevel = slog.LevelError
	default:
		cfg.LogLevel = slog.LevelInfo
	}

	// Production overrides
	if cfg.Environment == "production" {
		cfg.Debug = false
		cfg.LogFormat = "json"
		cfg.AllowedOrigins = []string{}
		cfg.RunMigrations = false
	}

	return cfg
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getBoolEnv(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if b, err := strconv.ParseBool(value); err == nil {
			return b
		}
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if d, err := time.ParseDuration(value); err == nil {
			return d
		}
	}
	return defaultValue
}
