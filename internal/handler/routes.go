package handler

import (
	"io/fs"
	"net/http"

	"github.com/dunamismax/go-web-server/internal/store"
	"github.com/dunamismax/go-web-server/ui"
	"github.com/labstack/echo/v4"
	"log/slog"
)

// Handlers holds all the application handlers.
type Handlers struct {
	Home *HomeHandler
	User *UserHandler
}

// NewHandlers creates a new handlers instance with the given store.
func NewHandlers(s *store.Store) *Handlers {
	return &Handlers{
		Home: NewHomeHandler(),
		User: NewUserHandler(s),
	}
}

// RegisterRoutes sets up all application routes.
func RegisterRoutes(e *echo.Echo, handlers *Handlers) error {
	// Serve static files
	staticFS, err := fs.Sub(ui.StaticFiles, "static")
	if err != nil {
		slog.Error("failed to create static file system", "error", err)
		return err
	}
	e.GET("/static/*", echo.WrapHandler(http.StripPrefix("/static/", http.FileServer(http.FS(staticFS)))))

	// Home routes
	e.GET("/", handlers.Home.Home)
	e.GET("/health", handlers.Home.Health)

	// User management routes
	e.GET("/users", handlers.User.Users)
	e.GET("/users/list", handlers.User.UserList)
	e.GET("/users/form", handlers.User.UserForm)
	e.GET("/users/:id/edit", handlers.User.EditUserForm)
	e.POST("/users", handlers.User.CreateUser)
	e.PUT("/users/:id", handlers.User.UpdateUser)
	e.PATCH("/users/:id/deactivate", handlers.User.DeactivateUser)
	e.DELETE("/users/:id", handlers.User.DeleteUser)

	// API routes
	api := e.Group("/api")
	api.GET("/users/count", handlers.User.UserCount)

	return nil
}
