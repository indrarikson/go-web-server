package main

import (
	"io/fs"
	"log/slog"
	"net/http"
	"os"

	"github.com/dunamismax/go-web-server/internal/config"
	"github.com/dunamismax/go-web-server/internal/handler"
	"github.com/dunamismax/go-web-server/internal/store"
	"github.com/dunamismax/go-web-server/ui"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

//go:generate go install github.com/a-h/templ/cmd/templ@latest
//go:generate templ generate
//go:generate sqlc generate

func main() {
	cfg := config.New()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: cfg.LogLevel,
	}))
	slog.SetDefault(logger)

	db, err := store.NewDatabase(cfg.DatabaseURL)
	if err != nil {
		slog.Error("failed to connect to database", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer db.Close()

	if err := db.InitSchema(); err != nil {
		slog.Error("failed to initialize schema", slog.String("error", err.Error()))
		os.Exit(1)
	}

	e := echo.New()

	e.HideBanner = true
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	staticFS, err := fs.Sub(ui.StaticFiles, "static")
	if err != nil {
		slog.Error("failed to create static file system", slog.String("error", err.Error()))
		os.Exit(1)
	}
	e.GET("/static/*", echo.WrapHandler(http.StripPrefix("/static/", http.FileServer(http.FS(staticFS)))))

	homeHandler := handler.NewHomeHandler()
	userHandler := handler.NewUserHandler(db)

	e.GET("/", homeHandler.Home)
	e.GET("/health", homeHandler.Health)
	e.GET("/users", userHandler.Users)
	e.GET("/users/list", userHandler.UserList)
	e.POST("/users", userHandler.CreateUser)
	e.DELETE("/users/:id", userHandler.DeleteUser)

	slog.Info("starting server", slog.String("port", cfg.Port), slog.String("environment", cfg.Environment))

	if err := e.Start(":" + cfg.Port); err != nil {
		slog.Error("failed to start server", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
