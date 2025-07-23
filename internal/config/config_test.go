package config

import (
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	// Clear environment to test defaults
	clearEnv(t)

	cfg := New()

	// Test server defaults
	assert.Equal(t, "8080", cfg.Port)
	assert.Equal(t, "", cfg.Host)
	assert.Equal(t, 10*time.Second, cfg.ReadTimeout)
	assert.Equal(t, 10*time.Second, cfg.WriteTimeout)
	assert.Equal(t, 30*time.Second, cfg.ShutdownTimeout)

	// Test database defaults
	assert.Equal(t, "data.db", cfg.DatabaseURL)
	assert.Equal(t, 25, cfg.DatabaseMaxConnections)
	assert.Equal(t, 30*time.Second, cfg.DatabaseTimeout)
	assert.True(t, cfg.RunMigrations)

	// Test application defaults
	assert.Equal(t, "development", cfg.Environment)
	assert.False(t, cfg.Debug)
	assert.Equal(t, slog.LevelInfo, cfg.LogLevel)
	assert.Equal(t, "text", cfg.LogFormat)

	// Test security defaults
	assert.Equal(t, []string{"127.0.0.1"}, cfg.TrustedProxies)
	assert.True(t, cfg.EnableCORS)
	assert.Equal(t, []string{"*"}, cfg.AllowedOrigins)

	// Test feature flags
	assert.False(t, cfg.EnableMetrics)
	assert.False(t, cfg.EnablePprof)
}

func TestNew_WithEnvironmentVariables(t *testing.T) {
	tests := []struct {
		name     string
		envVars  map[string]string
		validate func(t *testing.T, cfg *Config)
	}{
		{
			name: "custom server config",
			envVars: map[string]string{
				"PORT":             "9000",
				"HOST":             "localhost",
				"READ_TIMEOUT":     "30s",
				"WRITE_TIMEOUT":    "45s",
				"SHUTDOWN_TIMEOUT": "60s",
			},
			validate: func(t *testing.T, cfg *Config) {
				assert.Equal(t, "9000", cfg.Port)
				assert.Equal(t, "localhost", cfg.Host)
				assert.Equal(t, 30*time.Second, cfg.ReadTimeout)
				assert.Equal(t, 45*time.Second, cfg.WriteTimeout)
				assert.Equal(t, 60*time.Second, cfg.ShutdownTimeout)
			},
		},
		{
			name: "custom database config",
			envVars: map[string]string{
				"DATABASE_URL":             "test.db",
				"DATABASE_MAX_CONNECTIONS": "50",
				"DATABASE_TIMEOUT":         "60s",
				"RUN_MIGRATIONS":           "false",
			},
			validate: func(t *testing.T, cfg *Config) {
				assert.Equal(t, "test.db", cfg.DatabaseURL)
				assert.Equal(t, 50, cfg.DatabaseMaxConnections)
				assert.Equal(t, 60*time.Second, cfg.DatabaseTimeout)
				assert.False(t, cfg.RunMigrations)
			},
		},
		{
			name: "debug log level",
			envVars: map[string]string{
				"LOG_LEVEL": "debug",
			},
			validate: func(t *testing.T, cfg *Config) {
				assert.Equal(t, slog.LevelDebug, cfg.LogLevel)
				assert.True(t, cfg.Debug)
			},
		},
		{
			name: "warn log level",
			envVars: map[string]string{
				"LOG_LEVEL": "warn",
			},
			validate: func(t *testing.T, cfg *Config) {
				assert.Equal(t, slog.LevelWarn, cfg.LogLevel)
				assert.False(t, cfg.Debug)
			},
		},
		{
			name: "error log level",
			envVars: map[string]string{
				"LOG_LEVEL": "error",
			},
			validate: func(t *testing.T, cfg *Config) {
				assert.Equal(t, slog.LevelError, cfg.LogLevel)
				assert.False(t, cfg.Debug)
			},
		},
		{
			name: "feature flags enabled",
			envVars: map[string]string{
				"ENABLE_METRICS": "true",
				"ENABLE_PPROF":   "true",
				"ENABLE_CORS":    "false",
			},
			validate: func(t *testing.T, cfg *Config) {
				assert.True(t, cfg.EnableMetrics)
				assert.True(t, cfg.EnablePprof)
				assert.False(t, cfg.EnableCORS)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearEnv(t)

			// Set environment variables
			for key, value := range tt.envVars {
				t.Setenv(key, value)
			}

			cfg := New()
			tt.validate(t, cfg)
		})
	}
}

func TestNew_ProductionEnvironment(t *testing.T) {
	clearEnv(t)
	t.Setenv("ENVIRONMENT", "production")

	cfg := New()

	assert.Equal(t, "production", cfg.Environment)
	assert.False(t, cfg.Debug)
	assert.Equal(t, "json", cfg.LogFormat)
	assert.Equal(t, []string{}, cfg.AllowedOrigins)
	assert.False(t, cfg.RunMigrations)
}

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue string
		envValue     string
		expected     string
	}{
		{
			name:         "returns environment value when set",
			key:          "TEST_KEY",
			defaultValue: "default",
			envValue:     "env_value",
			expected:     "env_value",
		},
		{
			name:         "returns default when env not set",
			key:          "MISSING_KEY",
			defaultValue: "default",
			envValue:     "",
			expected:     "default",
		},
		{
			name:         "returns empty string when env is empty",
			key:          "EMPTY_KEY",
			defaultValue: "default",
			envValue:     "",
			expected:     "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				t.Setenv(tt.key, tt.envValue)
			}

			result := getEnv(tt.key, tt.defaultValue)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetBoolEnv(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue bool
		envValue     string
		expected     bool
	}{
		{
			name:         "returns true when env is true",
			key:          "TEST_BOOL",
			defaultValue: false,
			envValue:     "true",
			expected:     true,
		},
		{
			name:         "returns false when env is false",
			key:          "TEST_BOOL",
			defaultValue: true,
			envValue:     "false",
			expected:     false,
		},
		{
			name:         "returns default when env is invalid",
			key:          "TEST_BOOL",
			defaultValue: true,
			envValue:     "invalid",
			expected:     true,
		},
		{
			name:         "returns default when env not set",
			key:          "MISSING_BOOL",
			defaultValue: false,
			envValue:     "",
			expected:     false,
		},
		{
			name:         "handles 1 as true",
			key:          "TEST_BOOL",
			defaultValue: false,
			envValue:     "1",
			expected:     true,
		},
		{
			name:         "handles 0 as false",
			key:          "TEST_BOOL",
			defaultValue: true,
			envValue:     "0",
			expected:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				t.Setenv(tt.key, tt.envValue)
			}

			result := getBoolEnv(tt.key, tt.defaultValue)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetIntEnv(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue int
		envValue     string
		expected     int
	}{
		{
			name:         "returns parsed int when valid",
			key:          "TEST_INT",
			defaultValue: 10,
			envValue:     "42",
			expected:     42,
		},
		{
			name:         "returns default when env is invalid",
			key:          "TEST_INT",
			defaultValue: 10,
			envValue:     "invalid",
			expected:     10,
		},
		{
			name:         "returns default when env not set",
			key:          "MISSING_INT",
			defaultValue: 99,
			envValue:     "",
			expected:     99,
		},
		{
			name:         "handles negative numbers",
			key:          "TEST_INT",
			defaultValue: 10,
			envValue:     "-5",
			expected:     -5,
		},
		{
			name:         "handles zero",
			key:          "TEST_INT",
			defaultValue: 10,
			envValue:     "0",
			expected:     0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				t.Setenv(tt.key, tt.envValue)
			}

			result := getIntEnv(tt.key, tt.defaultValue)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetDurationEnv(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue time.Duration
		envValue     string
		expected     time.Duration
	}{
		{
			name:         "returns parsed duration when valid",
			key:          "TEST_DURATION",
			defaultValue: 10 * time.Second,
			envValue:     "30s",
			expected:     30 * time.Second,
		},
		{
			name:         "returns default when env is invalid",
			key:          "TEST_DURATION",
			defaultValue: 10 * time.Second,
			envValue:     "invalid",
			expected:     10 * time.Second,
		},
		{
			name:         "returns default when env not set",
			key:          "MISSING_DURATION",
			defaultValue: 5 * time.Minute,
			envValue:     "",
			expected:     5 * time.Minute,
		},
		{
			name:         "handles minutes",
			key:          "TEST_DURATION",
			defaultValue: 10 * time.Second,
			envValue:     "2m",
			expected:     2 * time.Minute,
		},
		{
			name:         "handles hours",
			key:          "TEST_DURATION",
			defaultValue: 10 * time.Second,
			envValue:     "1h",
			expected:     1 * time.Hour,
		},
		{
			name:         "handles milliseconds",
			key:          "TEST_DURATION",
			defaultValue: 10 * time.Second,
			envValue:     "500ms",
			expected:     500 * time.Millisecond,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				t.Setenv(tt.key, tt.envValue)
			}

			result := getDurationEnv(tt.key, tt.defaultValue)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// clearEnv clears relevant environment variables for testing
func clearEnv(t *testing.T) {
	t.Helper()

	envVars := []string{
		"PORT", "HOST", "READ_TIMEOUT", "WRITE_TIMEOUT", "SHUTDOWN_TIMEOUT",
		"DATABASE_URL", "DATABASE_MAX_CONNECTIONS", "DATABASE_TIMEOUT", "RUN_MIGRATIONS",
		"ENVIRONMENT", "DEBUG", "LOG_LEVEL", "LOG_FORMAT",
		"ENABLE_CORS", "ENABLE_METRICS", "ENABLE_PPROF",
	}

	for _, env := range envVars {
		os.Unsetenv(env)
	}
}
