package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/dunamismax/go-web-server/internal/testutil"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func TestUserHandler_Users(t *testing.T) {
	store := testutil.TestStore(t)
	handler := NewUserHandler(store)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.Users(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "User Management")
}

func TestUserHandler_UserList(t *testing.T) {
	store := testutil.TestStore(t)
	handler := NewUserHandler(store)

	// Create test user
	testutil.CreateTestUser(t, store)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/users/list", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.UserList(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "Test User")
	assert.Contains(t, rec.Body.String(), "test@example.com")
}

func TestUserHandler_UserCount(t *testing.T) {
	store := testutil.TestStore(t)
	handler := NewUserHandler(store)

	// Create test user
	testutil.CreateTestUser(t, store)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/users/count", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.UserCount(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "1")
}

func TestUserHandler_UserForm(t *testing.T) {
	store := testutil.TestStore(t)
	handler := NewUserHandler(store)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/users/form", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.UserForm(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "form")
}

func TestUserHandler_EditUserForm(t *testing.T) {
	store := testutil.TestStore(t)
	handler := NewUserHandler(store)

	// Create test user
	user := testutil.CreateTestUser(t, store)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/users/1/edit", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := handler.EditUserForm(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), user.Name)
	assert.Contains(t, rec.Body.String(), user.Email)
}

func TestUserHandler_EditUserForm_InvalidID(t *testing.T) {
	store := testutil.TestStore(t)
	handler := NewUserHandler(store)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/users/invalid/edit", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("invalid")

	err := handler.EditUserForm(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "Invalid user ID")
}

func TestUserHandler_EditUserForm_UserNotFound(t *testing.T) {
	store := testutil.TestStore(t)
	handler := NewUserHandler(store)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/users/999/edit", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("999")

	err := handler.EditUserForm(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Contains(t, rec.Body.String(), "User not found")
}

func TestUserHandler_CreateUser(t *testing.T) {
	store := testutil.TestStore(t)
	handler := NewUserHandler(store)

	form := url.Values{}
	form.Add("name", "John Doe")
	form.Add("email", "john@example.com")
	form.Add("bio", "Software developer")
	form.Add("avatar_url", "https://example.com/avatar.jpg")

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.CreateUser(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "John Doe")
	assert.Contains(t, rec.Body.String(), "john@example.com")
	assert.Equal(t, "userCreated", rec.Header().Get("HX-Trigger"))
}

func TestUserHandler_CreateUser_MissingRequiredFields(t *testing.T) {
	store := testutil.TestStore(t)
	handler := NewUserHandler(store)

	tests := []struct {
		name     string
		formData url.Values
	}{
		{
			name:     "missing name",
			formData: url.Values{"email": []string{"test@example.com"}},
		},
		{
			name:     "missing email",
			formData: url.Values{"name": []string{"Test User"}},
		},
		{
			name:     "empty name",
			formData: url.Values{"name": []string{""}, "email": []string{"test@example.com"}},
		},
		{
			name:     "empty email",
			formData: url.Values{"name": []string{"Test User"}, "email": []string{""}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(tt.formData.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := handler.CreateUser(c)

			assert.NoError(t, err)
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Contains(t, rec.Body.String(), "Name and email are required")
		})
	}
}

func TestUserHandler_UpdateUser(t *testing.T) {
	store := testutil.TestStore(t)
	handler := NewUserHandler(store)

	// Create test user
	user := testutil.CreateTestUser(t, store)

	form := url.Values{}
	form.Add("name", "Updated Name")
	form.Add("bio", "Updated bio")
	form.Add("avatar_url", "https://example.com/new-avatar.jpg")

	e := echo.New()
	req := httptest.NewRequest(http.MethodPut, "/users/1", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := handler.UpdateUser(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "Updated Name")
	assert.Equal(t, "userUpdated", rec.Header().Get("HX-Trigger"))

	// Verify user was actually updated
	updatedUser, err := store.GetUser(context.Background(), user.ID)
	require.NoError(t, err)
	assert.Equal(t, "Updated Name", updatedUser.Name)
}

func TestUserHandler_UpdateUser_InvalidID(t *testing.T) {
	store := testutil.TestStore(t)
	handler := NewUserHandler(store)

	form := url.Values{}
	form.Add("name", "Updated Name")

	e := echo.New()
	req := httptest.NewRequest(http.MethodPut, "/users/invalid", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("invalid")

	err := handler.UpdateUser(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "Invalid user ID")
}

func TestUserHandler_UpdateUser_MissingName(t *testing.T) {
	store := testutil.TestStore(t)
	handler := NewUserHandler(store)

	// Create test user
	testutil.CreateTestUser(t, store)

	form := url.Values{}
	form.Add("bio", "Updated bio")

	e := echo.New()
	req := httptest.NewRequest(http.MethodPut, "/users/1", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := handler.UpdateUser(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "Name is required")
}

func TestUserHandler_DeactivateUser(t *testing.T) {
	store := testutil.TestStore(t)
	handler := NewUserHandler(store)

	// Create test user
	user := testutil.CreateTestUser(t, store)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPatch, "/users/1/deactivate", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := handler.DeactivateUser(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "userDeactivated", rec.Header().Get("HX-Trigger"))

	// Verify user count decreased
	count, err := store.CountUsers(context.Background())
	require.NoError(t, err)
	assert.Equal(t, int64(0), count)

	// Verify user still exists but is deactivated
	deactivatedUser, err := store.GetUser(context.Background(), user.ID)
	require.NoError(t, err)
	assert.False(t, deactivatedUser.IsActive.Bool)
}

func TestUserHandler_DeactivateUser_InvalidID(t *testing.T) {
	store := testutil.TestStore(t)
	handler := NewUserHandler(store)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPatch, "/users/invalid/deactivate", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("invalid")

	err := handler.DeactivateUser(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "Invalid user ID")
}

func TestUserHandler_DeleteUser(t *testing.T) {
	store := testutil.TestStore(t)
	handler := NewUserHandler(store)

	// Create test user
	user := testutil.CreateTestUser(t, store)

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/users/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := handler.DeleteUser(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "userDeleted", rec.Header().Get("HX-Trigger"))
	assert.Empty(t, rec.Body.String())

	// Verify user was actually deleted
	_, err = store.GetUser(context.Background(), user.ID)
	assert.Error(t, err) // Should return error because user doesn't exist
}

func TestUserHandler_DeleteUser_InvalidID(t *testing.T) {
	store := testutil.TestStore(t)
	handler := NewUserHandler(store)

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/users/invalid", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("invalid")

	err := handler.DeleteUser(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "Invalid user ID")
}

func TestUserHandler_UserList_DatabaseError(t *testing.T) {
	store := testutil.TestStore(t)
	handler := NewUserHandler(store)

	// Close the database to trigger an error
	store.Close()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/users/list", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.UserList(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "Failed to fetch users")
}

func TestUserHandler_UserCount_DatabaseError(t *testing.T) {
	store := testutil.TestStore(t)
	handler := NewUserHandler(store)

	// Close the database to trigger an error
	store.Close()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/users/count", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.UserCount(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "Failed to count users")
}

func TestUserHandler_CreateUser_DatabaseError(t *testing.T) {
	store := testutil.TestStore(t)
	handler := NewUserHandler(store)

	form := url.Values{}
	form.Add("name", "John Doe")
	form.Add("email", "john@example.com")

	// Close the database to trigger an error
	store.Close()

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.CreateUser(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "Failed to create user")
}

func TestUserHandler_UpdateUser_DatabaseError(t *testing.T) {
	store := testutil.TestStore(t)
	handler := NewUserHandler(store)

	// Create test user first
	_ = testutil.CreateTestUser(t, store)

	form := url.Values{}
	form.Add("name", "Updated Name")

	// Close the database to trigger an error
	store.Close()

	e := echo.New()
	req := httptest.NewRequest(http.MethodPut, "/users/1", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := handler.UpdateUser(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "Failed to update user")
}

func TestUserHandler_DeactivateUser_DatabaseError(t *testing.T) {
	store := testutil.TestStore(t)
	handler := NewUserHandler(store)

	// Create test user first
	_ = testutil.CreateTestUser(t, store)

	// Close the database to trigger an error
	store.Close()

	e := echo.New()
	req := httptest.NewRequest(http.MethodPatch, "/users/1/deactivate", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := handler.DeactivateUser(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "Failed to deactivate user")
}

func TestUserHandler_DeleteUser_DatabaseError(t *testing.T) {
	store := testutil.TestStore(t)
	handler := NewUserHandler(store)

	// Create test user first
	_ = testutil.CreateTestUser(t, store)

	// Close the database to trigger an error
	store.Close()

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/users/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := handler.DeleteUser(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "Failed to delete user")
}
