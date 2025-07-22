package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/dunamismax/go-web-server/internal/store"
	"github.com/dunamismax/go-web-server/internal/view"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	db *store.Database
}

func NewUserHandler(db *store.Database) *UserHandler {
	return &UserHandler{
		db: db,
	}
}

func (h *UserHandler) Users(c echo.Context) error {
	ctx := context.Background()
	users, err := h.db.ListUsers(ctx)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to fetch users")
	}

	component := view.Users(users)
	return component.Render(c.Request().Context(), c.Response().Writer)
}

func (h *UserHandler) UserList(c echo.Context) error {
	ctx := context.Background()
	users, err := h.db.ListUsers(ctx)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to fetch users")
	}

	component := view.UserList(users)
	return component.Render(c.Request().Context(), c.Response().Writer)
}

func (h *UserHandler) CreateUser(c echo.Context) error {
	ctx := context.Background()

	name := c.FormValue("name")
	email := c.FormValue("email")

	if name == "" || email == "" {
		return c.String(http.StatusBadRequest, "Name and email are required")
	}

	params := store.CreateUserParams{
		Name:  name,
		Email: email,
	}

	_, err := h.db.CreateUser(ctx, params)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to create user")
	}

	users, err := h.db.ListUsers(ctx)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to fetch users")
	}

	component := view.UserList(users)
	return component.Render(c.Request().Context(), c.Response().Writer)
}

func (h *UserHandler) DeleteUser(c echo.Context) error {
	ctx := context.Background()

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid user ID")
	}

	err = h.db.DeleteUser(ctx, id)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to delete user")
	}

	users, err := h.db.ListUsers(ctx)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to fetch users")
	}

	component := view.UserList(users)
	return component.Render(c.Request().Context(), c.Response().Writer)
}
