package handler

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/dunamismax/go-web-server/internal/store"
	"github.com/dunamismax/go-web-server/internal/view"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	store *store.Store
}

func NewUserHandler(s *store.Store) *UserHandler {
	return &UserHandler{
		store: s,
	}
}

// Users renders the main user management page
func (h *UserHandler) Users(c echo.Context) error {
	component := view.Users()
	return component.Render(c.Request().Context(), c.Response().Writer)
}

// UserList returns the list of users as HTML fragment
func (h *UserHandler) UserList(c echo.Context) error {
	ctx := context.Background()
	users, err := h.store.ListUsers(ctx)
	if err != nil {
		slog.Error("Failed to fetch users", "error", err)
		return c.String(http.StatusInternalServerError, "Failed to fetch users")
	}

	component := view.UserList(users)
	return component.Render(c.Request().Context(), c.Response().Writer)
}

// UserCount returns the count of active users
func (h *UserHandler) UserCount(c echo.Context) error {
	ctx := context.Background()
	count, err := h.store.CountUsers(ctx)
	if err != nil {
		slog.Error("Failed to count users", "error", err)
		return c.String(http.StatusInternalServerError, "Failed to count users")
	}

	component := view.UserCount(count)
	return component.Render(c.Request().Context(), c.Response().Writer)
}

// UserForm renders the user creation/edit form
func (h *UserHandler) UserForm(c echo.Context) error {
	component := view.UserForm(nil)
	return component.Render(c.Request().Context(), c.Response().Writer)
}

// EditUserForm renders the user edit form with existing data
func (h *UserHandler) EditUserForm(c echo.Context) error {
	ctx := context.Background()
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid user ID")
	}

	user, err := h.store.GetUser(ctx, id)
	if err != nil {
		slog.Error("Failed to fetch user", "id", id, "error", err)
		return c.String(http.StatusNotFound, "User not found")
	}

	component := view.UserForm(&user)
	return component.Render(c.Request().Context(), c.Response().Writer)
}

// CreateUser creates a new user
func (h *UserHandler) CreateUser(c echo.Context) error {
	ctx := context.Background()

	name := c.FormValue("name")
	email := c.FormValue("email")
	bio := c.FormValue("bio")
	avatarUrl := c.FormValue("avatar_url")

	if name == "" || email == "" {
		return c.String(http.StatusBadRequest, "Name and email are required")
	}

	var bioSql sql.NullString
	if bio != "" {
		bioSql = sql.NullString{String: bio, Valid: true}
	}

	var avatarUrlSql sql.NullString
	if avatarUrl != "" {
		avatarUrlSql = sql.NullString{String: avatarUrl, Valid: true}
	}

	params := store.CreateUserParams{
		Email:     email,
		Name:      name,
		Bio:       bioSql,
		AvatarUrl: avatarUrlSql,
	}

	_, err := h.store.CreateUser(ctx, params)
	if err != nil {
		slog.Error("Failed to create user", "name", name, "email", email, "error", err)
		return c.String(http.StatusInternalServerError, "Failed to create user")
	}

	slog.Info("User created successfully", "name", name, "email", email)

	// Trigger custom event for HTMX
	c.Response().Header().Set("HX-Trigger", "userCreated")

	users, err := h.store.ListUsers(ctx)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to fetch users")
	}

	component := view.UserList(users)
	return component.Render(c.Request().Context(), c.Response().Writer)
}

// UpdateUser updates an existing user
func (h *UserHandler) UpdateUser(c echo.Context) error {
	ctx := context.Background()
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid user ID")
	}

	name := c.FormValue("name")
	bio := c.FormValue("bio")
	avatarUrl := c.FormValue("avatar_url")

	if name == "" {
		return c.String(http.StatusBadRequest, "Name is required")
	}

	var bioSql sql.NullString
	if bio != "" {
		bioSql = sql.NullString{String: bio, Valid: true}
	}

	var avatarUrlSql sql.NullString
	if avatarUrl != "" {
		avatarUrlSql = sql.NullString{String: avatarUrl, Valid: true}
	}

	params := store.UpdateUserParams{
		Name:      name,
		Bio:       bioSql,
		AvatarUrl: avatarUrlSql,
		ID:        id,
	}

	_, err = h.store.UpdateUser(ctx, params)
	if err != nil {
		slog.Error("Failed to update user", "id", id, "error", err)
		return c.String(http.StatusInternalServerError, "Failed to update user")
	}

	slog.Info("User updated successfully", "id", id, "name", name)

	// Trigger custom event for HTMX
	c.Response().Header().Set("HX-Trigger", "userUpdated")

	users, err := h.store.ListUsers(ctx)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to fetch users")
	}

	component := view.UserList(users)
	return component.Render(c.Request().Context(), c.Response().Writer)
}

// DeactivateUser deactivates a user instead of deleting
func (h *UserHandler) DeactivateUser(c echo.Context) error {
	ctx := context.Background()
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid user ID")
	}

	err = h.store.DeactivateUser(ctx, id)
	if err != nil {
		slog.Error("Failed to deactivate user", "id", id, "error", err)
		return c.String(http.StatusInternalServerError, "Failed to deactivate user")
	}

	slog.Info("User deactivated successfully", "id", id)

	// Get the updated user and return the row
	user, err := h.store.GetUser(ctx, id)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to fetch updated user")
	}

	// Trigger custom event for HTMX
	c.Response().Header().Set("HX-Trigger", "userDeactivated")

	component := view.UserRow(user)
	return component.Render(c.Request().Context(), c.Response().Writer)
}

// DeleteUser permanently deletes a user
func (h *UserHandler) DeleteUser(c echo.Context) error {
	ctx := context.Background()
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid user ID")
	}

	err = h.store.DeleteUser(ctx, id)
	if err != nil {
		slog.Error("Failed to delete user", "id", id, "error", err)
		return c.String(http.StatusInternalServerError, "Failed to delete user")
	}

	slog.Info("User deleted successfully", "id", id)

	// Trigger custom event for HTMX
	c.Response().Header().Set("HX-Trigger", "userDeleted")

	// Return empty response since the row should be removed
	return c.NoContent(http.StatusOK)
}
