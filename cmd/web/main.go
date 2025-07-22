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

	// Middleware
	if cfg.Debug {
		e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
			Format: "method=${method}, uri=${uri}, status=${status}, latency=${latency_human}\n",
		}))
	}

	e.Use(middleware.Recover())
	e.Use(middleware.Secure())

	if cfg.EnableCORS {
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: cfg.AllowedOrigins,
			AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete},
			AllowHeaders: []string{"*"},
		}))
	}

	// Request ID middleware for tracing
	e.Use(middleware.RequestID())

	// Rate limiting (basic)
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))

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

// runMigrations runs database migrations if available
func runMigrations(databaseURL string) error {
	// Placeholder for migration logic
	// In a real implementation, you would use golang-migrate here
	slog.Info("Migration system ready (placeholder - implement with golang-migrate)")
	return nil
}
