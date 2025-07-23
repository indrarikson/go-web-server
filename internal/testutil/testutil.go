// Package testutil provides testing utilities and helpers for the application.
package testutil

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/dunamismax/go-web-server/internal/store"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

// TestStore creates an in-memory SQLite database for testing.
// It automatically initializes the schema and sets up cleanup.
func TestStore(t *testing.T) *store.Store {
	t.Helper()

	s, err := store.NewStore(":memory:")
	require.NoError(t, err)

	err = s.InitSchema()
	require.NoError(t, err)

	t.Cleanup(func() {
		s.Close()
	})

	return s
}

// CreateTestUser creates a test user with default values for testing.
// Returns the created user.
func CreateTestUser(t *testing.T, s *store.Store) store.User {
	t.Helper()

	user, err := s.CreateUser(context.Background(), store.CreateUserParams{
		Email:     "test@example.com",
		Name:      "Test User",
		Bio:       sql.NullString{String: "Test bio", Valid: true},
		AvatarUrl: sql.NullString{},
	})
	require.NoError(t, err)

	return user
}

// CreateTestUserWithParams creates a test user with custom parameters.
func CreateTestUserWithParams(t *testing.T, s *store.Store, params store.CreateUserParams) store.User {
	t.Helper()

	user, err := s.CreateUser(context.Background(), params)
	require.NoError(t, err)

	return user
}

// CreateMultipleTestUsers creates multiple test users for testing.
// Returns a slice of created users.
func CreateMultipleTestUsers(t *testing.T, s *store.Store, count int) []store.User {
	t.Helper()

	users := make([]store.User, count)
	for i := 0; i < count; i++ {
		user, err := s.CreateUser(context.Background(), store.CreateUserParams{
			Email:     fmt.Sprintf("user%d@example.com", i+1),
			Name:      fmt.Sprintf("Test User %d", i+1),
			Bio:       sql.NullString{String: fmt.Sprintf("Bio for user %d", i+1), Valid: true},
			AvatarUrl: sql.NullString{},
		})
		require.NoError(t, err)
		users[i] = user
	}

	return users
}

// AssertUserCount verifies the user count in the store.
func AssertUserCount(t *testing.T, s *store.Store, expected int64) {
	t.Helper()

	count, err := s.CountUsers(context.Background())
	require.NoError(t, err)
	require.Equal(t, expected, count)
}

// AssertUserExists verifies that a user exists in the store.
func AssertUserExists(t *testing.T, s *store.Store, userID int64) store.User {
	t.Helper()

	user, err := s.GetUser(context.Background(), userID)
	require.NoError(t, err)

	return user
}

// AssertUserNotExists verifies that a user does not exist in the store.
func AssertUserNotExists(t *testing.T, s *store.Store, userID int64) {
	t.Helper()

	_, err := s.GetUser(context.Background(), userID)
	require.Error(t, err)
}
