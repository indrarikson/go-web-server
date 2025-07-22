package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/dunamismax/go-web-server/internal/config"
	"github.com/dunamismax/go-web-server/internal/handler"
	"github.com/dunamismax/go-web-server/internal/store"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

//go:generate go install github.com/a-h/templ/cmd/templ@latest
//go:generate go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
//go:generate templ generate
//go:generate sh -c "cd ../../ && sqlc generate"

func main() {
	// Load configuration
	cfg := config.New()

	// Setup structured logging
	var logger *slog.Logger
	if cfg.LogFormat == "json" {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: cfg.LogLevel,
		}))
	} else {
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: cfg.LogLevel,
		}))
	}
	slog.SetDefault(logger)

	slog.Info("Starting Go Web Server",
		"version", "1.0.0",
		"environment", cfg.Environment,
		"go_version", "1.24+",
		"port", cfg.Port,
		"debug", cfg.Debug)

	// Initialize database store
	store, err := store.NewStore(cfg.DatabaseURL)
	if err != nil {
		slog.Error("failed to connect to database", "error", err, "database_url", cfg.DatabaseURL)
		os.Exit(1)
	}
	defer func() {
		if err := store.Close(); err != nil {
			slog.Error("failed to close database connection", "error", err)
		}
	}()

	// Run migrations if enabled
	if cfg.RunMigrations {
		slog.Info("Running database migrations")
		if err := runMigrations(cfg.DatabaseURL); err != nil {
			slog.Error("failed to run migrations", "error", err)
			os.Exit(1)
		}
	}

	// Initialize schema (fallback if migrations not used)
	if err := store.InitSchema(); err != nil {
		slog.Error("failed to initialize schema", "error", err)
		os.Exit(1)
	}

	// Create Echo instance
	e := echo.New()
	e.HideBanner = true
	e.Debug = cfg.Debug

	// Configure timeouts
	e.Server.ReadTimeout = cfg.ReadTimeout
	e.Server.WriteTimeout = cfg.WriteTimeout

	// Middleware stack (order matters)

	// Recovery middleware (should be first)
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		LogErrorFunc: func(c echo.Context, err error, stack []byte) error {
			slog.Error("panic recovered",
				"error", err,
				"path", c.Request().URL.Path,
				"method", c.Request().Method,
				"request_id", c.Response().Header().Get(echo.HeaderXRequestID))
			return nil
		},
	}))

	// Request ID middleware for tracing
	e.Use(middleware.RequestID())

	// Structured logging middleware
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:    true,
		LogURI:       true,
		LogError:     true,
		LogMethod:    true,
		LogLatency:   true,
		LogRemoteIP:  true,
		LogUserAgent: cfg.Debug,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				slog.Info("request",
					"method", v.Method,
					"uri", v.URI,
					"status", v.Status,
					"latency", v.Latency.String(),
					"remote_ip", v.RemoteIP,
					"request_id", v.RequestID)
			} else {
				slog.Error("request error",
					"method", v.Method,
					"uri", v.URI,
					"status", v.Status,
					"latency", v.Latency.String(),
					"remote_ip", v.RemoteIP,
					"request_id", v.RequestID,
					"error", v.Error)
			}
			return nil
		},
	}))

	// Security middleware
	e.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection:         "1; mode=block",
		ContentTypeNosniff:    "nosniff",
		XFrameOptions:         "DENY",
		HSTSMaxAge:            31536000,
		ContentSecurityPolicy: "default-src 'self'; style-src 'self' 'unsafe-inline'; script-src 'self' 'unsafe-inline'",
	}))

	// CORS middleware
	if cfg.EnableCORS {
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: cfg.AllowedOrigins,
			AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete},
			AllowHeaders: []string{"*"},
			MaxAge:       86400,
		}))
	}

	// Rate limiting
	e.Use(middleware.RateLimiterWithConfig(middleware.RateLimiterConfig{
		Store: middleware.NewRateLimiterMemoryStore(20),
		IdentifierExtractor: func(c echo.Context) (string, error) {
			return c.RealIP(), nil
		},
		ErrorHandler: func(context echo.Context, err error) error {
			return context.JSON(http.StatusTooManyRequests, map[string]string{
				"error": "rate limit exceeded",
			})
		},
	}))

	// Timeout middleware
	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout: cfg.ReadTimeout,
	}))

	// Initialize handlers and register routes
	handlers := handler.NewHandlers(store)
	if err := handler.RegisterRoutes(e, handlers); err != nil {
		slog.Error("failed to register routes", "error", err)
		os.Exit(1)
	}

	// Graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Start server in goroutine
	go func() {
		address := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
		slog.Info("Server starting", "address", address)

		if err := e.Start(address); err != nil && err != http.ErrServerClosed {
			slog.Error("failed to start server", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal
	<-ctx.Done()

	slog.Info("Shutting down server...")

	// Graceful shutdown with timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	if err := e.Shutdown(shutdownCtx); err != nil {
		slog.Error("failed to shutdown server gracefully", "error", err)
		os.Exit(1)
	}

	slog.Info("Server shutdown complete")
}

// runMigrations runs database migrations using golang-migrate
func runMigrations(databaseURL string) error {
	m, err := migrate.New(
		"file://internal/store/migrations",
		fmt.Sprintf("sqlite://%s", databaseURL),
	)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}
	defer m.Close()

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	if err == migrate.ErrNoChange {
		slog.Info("No new migrations to apply")
	} else {
		slog.Info("Database migrations completed successfully")
	}

	return nil
}
