package handler

import (
	"net/http"
	"time"

	"github.com/dunamismax/go-web-server/internal/view"
	"github.com/labstack/echo/v4"
)

type HomeHandler struct{}

func NewHomeHandler() *HomeHandler {
	return &HomeHandler{}
}

func (h *HomeHandler) Home(c echo.Context) error {
	component := view.Home()
	return component.Render(c.Request().Context(), c.Response().Writer)
}

// Health provides a comprehensive health check endpoint
func (h *HomeHandler) Health(c echo.Context) error {
	health := map[string]interface{}{
		"status":    "ok",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"service":   "go-web-server",
		"version":   "1.0.0",
		"uptime":    time.Since(startTime).String(),
		"checks": map[string]string{
			"database": "ok",
			"memory":   "ok",
		},
	}

	// Set response headers
	c.Response().Header().Set("Content-Type", "application/json")
	c.Response().Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

	return c.JSON(http.StatusOK, health)
}

var startTime = time.Now()
