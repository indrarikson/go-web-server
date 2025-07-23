package handler

import (
	"testing"

	"github.com/dunamismax/go-web-server/internal/testutil"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewHandlers(t *testing.T) {
	store := testutil.TestStore(t)

	handlers := NewHandlers(store)

	assert.NotNil(t, handlers)
	assert.NotNil(t, handlers.Home)
	assert.NotNil(t, handlers.User)
}

func TestRegisterRoutes(t *testing.T) {
	store := testutil.TestStore(t)
	handlers := NewHandlers(store)

	e := echo.New()
	err := RegisterRoutes(e, handlers)

	require.NoError(t, err)

	// Verify routes are registered by checking the router
	routes := e.Routes()
	assert.NotEmpty(t, routes)

	// Check for key routes
	routePaths := make(map[string]bool)
	for _, route := range routes {
		routePaths[route.Method+":"+route.Path] = true
	}

	expectedRoutes := []string{
		"GET:/",
		"GET:/health",
		"GET:/users",
		"GET:/users/list",
		"GET:/users/form",
		"GET:/users/:id/edit",
		"POST:/users",
		"PUT:/users/:id",
		"PATCH:/users/:id/deactivate",
		"DELETE:/users/:id",
		"GET:/api/users/count",
		"GET:/static/*",
	}

	for _, expectedRoute := range expectedRoutes {
		assert.True(t, routePaths[expectedRoute], "Route %s should be registered", expectedRoute)
	}
}

func TestRegisterRoutes_StaticFileError(t *testing.T) {
	// This test would require mocking the fs.Sub function which is complex
	// The error path in RegisterRoutes (lines 31-34) is hard to trigger
	// since ui.StaticFiles is embedded at compile time
	// For now, we'll document that this error case is difficult to test
	// without significant refactoring of the RegisterRoutes function

	store := testutil.TestStore(t)
	handlers := NewHandlers(store)

	e := echo.New()
	err := RegisterRoutes(e, handlers)

	// Normal case should always succeed
	require.NoError(t, err)
}
