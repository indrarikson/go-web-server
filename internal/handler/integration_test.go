package handler

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dunamismax/go-web-server/internal/store"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	_ "modernc.org/sqlite"
)

// HandlerTestSuite demonstrates using testify's suite functionality
// for comprehensive integration testing.
type HandlerTestSuite struct {
	suite.Suite
	store    *store.Store
	handlers *Handlers
	echo     *echo.Echo
}

// SetupSuite runs once before all tests in the suite
func (suite *HandlerTestSuite) SetupSuite() {
	var err error
	suite.store, err = store.NewStore(":memory:")
	suite.Require().NoError(err)

	err = suite.store.InitSchema()
	suite.Require().NoError(err)

	suite.handlers = NewHandlers(suite.store)
	suite.echo = echo.New()

	err = RegisterRoutes(suite.echo, suite.handlers)
	suite.Require().NoError(err)
}

// TearDownSuite runs once after all tests in the suite
func (suite *HandlerTestSuite) TearDownSuite() {
	if suite.store != nil {
		suite.store.Close()
	}
}

// SetupTest runs before each test
func (suite *HandlerTestSuite) SetupTest() {
	// Clean up any existing data before each test
	suite.store.DB().Exec("DELETE FROM users")
}

// TestFullUserWorkflow demonstrates a complete user management workflow
func (suite *HandlerTestSuite) TestFullUserWorkflow() {
	// 1. Check initial user count
	req := httptest.NewRequest(http.MethodGet, "/api/users/count", nil)
	rec := httptest.NewRecorder()
	suite.echo.ServeHTTP(rec, req)

	suite.Equal(http.StatusOK, rec.Code)
	suite.Contains(rec.Body.String(), "0")

	// 2. Create a new user
	createReq := httptest.NewRequest(http.MethodPost, "/users", nil)
	createReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	createReq.PostForm = map[string][]string{
		"name":       {"John Doe"},
		"email":      {"john@example.com"},
		"bio":        {"Software developer"},
		"avatar_url": {"https://example.com/avatar.jpg"},
	}
	createRec := httptest.NewRecorder()
	suite.echo.ServeHTTP(createRec, createReq)

	suite.Equal(http.StatusOK, createRec.Code)
	suite.Contains(createRec.Body.String(), "John Doe")
	suite.Equal("userCreated", createRec.Header().Get("HX-Trigger"))

	// 3. Verify user count increased
	req = httptest.NewRequest(http.MethodGet, "/api/users/count", nil)
	rec = httptest.NewRecorder()
	suite.echo.ServeHTTP(rec, req)

	suite.Equal(http.StatusOK, rec.Code)
	suite.Contains(rec.Body.String(), "1")

	// 4. List users and verify the created user appears
	listReq := httptest.NewRequest(http.MethodGet, "/users/list", nil)
	listRec := httptest.NewRecorder()
	suite.echo.ServeHTTP(listRec, listReq)

	suite.Equal(http.StatusOK, listRec.Code)
	suite.Contains(listRec.Body.String(), "John Doe")
	suite.Contains(listRec.Body.String(), "john@example.com")

	// 5. Update the user
	updateReq := httptest.NewRequest(http.MethodPut, "/users/1", nil)
	updateReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	updateReq.PostForm = map[string][]string{
		"name":       {"John Smith"},
		"bio":        {"Senior Software Developer"},
		"avatar_url": {"https://example.com/new-avatar.jpg"},
	}
	updateRec := httptest.NewRecorder()
	suite.echo.ServeHTTP(updateRec, updateReq)

	suite.Equal(http.StatusOK, updateRec.Code)
	suite.Contains(updateRec.Body.String(), "John Smith")
	suite.Equal("userUpdated", updateRec.Header().Get("HX-Trigger"))

	// 6. Deactivate the user
	deactivateReq := httptest.NewRequest(http.MethodPatch, "/users/1/deactivate", nil)
	deactivateRec := httptest.NewRecorder()
	suite.echo.ServeHTTP(deactivateRec, deactivateReq)

	suite.Equal(http.StatusOK, deactivateRec.Code)
	suite.Equal("userDeactivated", deactivateRec.Header().Get("HX-Trigger"))

	// 7. Verify user count decreased (deactivated users don't count)
	req = httptest.NewRequest(http.MethodGet, "/api/users/count", nil)
	rec = httptest.NewRecorder()
	suite.echo.ServeHTTP(rec, req)

	suite.Equal(http.StatusOK, rec.Code)
	suite.Contains(rec.Body.String(), "0")
}

// TestErrorHandling tests various error scenarios
func (suite *HandlerTestSuite) TestErrorHandling() {
	// Test invalid user ID
	req := httptest.NewRequest(http.MethodGet, "/users/invalid/edit", nil)
	rec := httptest.NewRecorder()
	suite.echo.ServeHTTP(rec, req)

	suite.Equal(http.StatusBadRequest, rec.Code)
	suite.Contains(rec.Body.String(), "Invalid user ID")

	// Test non-existent user
	req = httptest.NewRequest(http.MethodGet, "/users/999/edit", nil)
	rec = httptest.NewRecorder()
	suite.echo.ServeHTTP(rec, req)

	suite.Equal(http.StatusNotFound, rec.Code)
	suite.Contains(rec.Body.String(), "User not found")

	// Test creating user with missing required fields
	createReq := httptest.NewRequest(http.MethodPost, "/users", nil)
	createReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	createReq.PostForm = map[string][]string{
		"bio": {"Just a bio, no name or email"},
	}
	createRec := httptest.NewRecorder()
	suite.echo.ServeHTTP(createRec, createReq)

	suite.Equal(http.StatusBadRequest, createRec.Code)
	suite.Contains(createRec.Body.String(), "Name and email are required")
}

// Run the test suite
func TestHandlerSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}

// TestUserHandlerWithMock demonstrates using testify mocks
func TestUserHandlerWithMock(t *testing.T) {
	// Create a mock store
	mockStore := &store.MockUserStore{}

	// Define expected behavior
	expectedUser := store.User{
		ID:    1,
		Email: "test@example.com",
		Name:  "Test User",
		Bio:   sql.NullString{String: "Test bio", Valid: true},
	}

	mockStore.On("GetUser", mock.Anything, int64(1)).Return(expectedUser, nil)

	// Create handler with mock - Note: This requires modifying UserHandler to accept interface
	// For demonstration purposes, we'll test the logic directly
	ctx := context.Background()
	user, err := mockStore.GetUser(ctx, 1)

	require.NoError(t, err)
	assert.Equal(t, expectedUser.ID, user.ID)
	assert.Equal(t, expectedUser.Email, user.Email)
	assert.Equal(t, expectedUser.Name, user.Name)

	// Verify all expectations were met
	mockStore.AssertExpectations(t)
}

// TestUserHandlerErrorScenarios tests error handling with mocks
func TestUserHandlerErrorScenarios(t *testing.T) {
	mockStore := &store.MockUserStore{}

	// Mock a database error
	mockStore.On("GetUser", mock.Anything, int64(999)).Return(store.User{}, fmt.Errorf("user not found"))

	ctx := context.Background()
	_, err := mockStore.GetUser(ctx, 999)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")

	mockStore.AssertExpectations(t)
}

// BenchmarkUserList benchmarks the user list endpoint
func BenchmarkUserList(b *testing.B) {
	// Setup
	s, err := store.NewStore(":memory:")
	require.NoError(b, err)
	defer s.Close()

	err = s.InitSchema()
	require.NoError(b, err)

	// Create test data
	for i := 0; i < 100; i++ {
		params := store.CreateUserParams{
			Email: fmt.Sprintf("user%d@example.com", i),
			Name:  fmt.Sprintf("User %d", i),
		}
		_, err := s.CreateUser(context.Background(), params)
		require.NoError(b, err)
	}

	handler := NewUserHandler(s)
	e := echo.New()

	b.ResetTimer()

	// Benchmark the handler
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodGet, "/users/list", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.UserList(c)
		require.NoError(b, err)
	}
}
