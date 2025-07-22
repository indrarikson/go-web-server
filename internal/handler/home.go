package handler

import (
	"net/http"

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

func (h *HomeHandler) Health(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status":  "ok",
		"message": "Go Web Server is running",
	})
}
