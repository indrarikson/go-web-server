package testutil

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/dunamismax/go-web-server/internal/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func TestTestStore(t *testing.T) {
	s := TestStore(t)
	assert.NotNil(t, s)

	// Verify schema is initialized
	var count int
	err := s.DB().QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='users'").Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 1, count)

	// Verify we can perform basic operations
	userCount, err := s.CountUsers(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, int64(0), userCount)
}

func TestCreateTestUser(t *testing.T) {
	s := TestStore(t)

	user := CreateTestUser(t, s)

	assert.NotZero(t, user.ID)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, "Test User", user.Name)
	assert.True(t, user.Bio.Valid)
	assert.Equal(t, "Test bio", user.Bio.String)
	assert.False(t, user.AvatarUrl.Valid)
}

func TestCreateTestUserWithParams(t *testing.T) {
	s := TestStore(t)

	params := store.CreateUserParams{
		Email:     "custom@example.com",
		Name:      "Custom User",
		Bio:       sql.NullString{String: "Custom bio", Valid: true},
		AvatarUrl: sql.NullString{String: "https://example.com/avatar.jpg", Valid: true},
	}

	user := CreateTestUserWithParams(t, s, params)

	assert.NotZero(t, user.ID)
	assert.Equal(t, "custom@example.com", user.Email)
	assert.Equal(t, "Custom User", user.Name)
	assert.True(t, user.Bio.Valid)
	assert.Equal(t, "Custom bio", user.Bio.String)
	assert.True(t, user.AvatarUrl.Valid)
	assert.Equal(t, "https://example.com/avatar.jpg", user.AvatarUrl.String)
}

func TestCreateMultipleTestUsers(t *testing.T) {
	s := TestStore(t)

	users := CreateMultipleTestUsers(t, s, 3)

	assert.Len(t, users, 3)

	for i, user := range users {
		assert.NotZero(t, user.ID)
		assert.Equal(t, int64(i+1), user.ID) // Auto-increment IDs
		assert.Equal(t, fmt.Sprintf("user%d@example.com", i+1), user.Email)
		assert.Equal(t, fmt.Sprintf("Test User %d", i+1), user.Name)
	}

	// Verify count
	count, err := s.CountUsers(context.Background())
	require.NoError(t, err)
	assert.Equal(t, int64(3), count)
}

func TestCreateMultipleTestUsers_Zero(t *testing.T) {
	s := TestStore(t)

	users := CreateMultipleTestUsers(t, s, 0)

	assert.Len(t, users, 0)

	// Verify count
	count, err := s.CountUsers(context.Background())
	require.NoError(t, err)
	assert.Equal(t, int64(0), count)
}

func TestAssertUserCount(t *testing.T) {
	s := TestStore(t)

	// Initially should be 0
	AssertUserCount(t, s, 0)

	// Create a user
	CreateTestUser(t, s)
	AssertUserCount(t, s, 1)

	// Create more users
	CreateMultipleTestUsers(t, s, 2)
	AssertUserCount(t, s, 3)
}

func TestAssertUserExists(t *testing.T) {
	s := TestStore(t)

	user := CreateTestUser(t, s)

	existingUser := AssertUserExists(t, s, user.ID)
	assert.Equal(t, user.ID, existingUser.ID)
	assert.Equal(t, user.Email, existingUser.Email)
	assert.Equal(t, user.Name, existingUser.Name)
}

func TestAssertUserNotExists(t *testing.T) {
	s := TestStore(t)

	// Create and then delete a user
	user := CreateTestUser(t, s)
	err := s.DeleteUser(context.Background(), user.ID)
	require.NoError(t, err)

	// Assert user no longer exists
	AssertUserNotExists(t, s, user.ID)

	// Also test with a user that never existed
	AssertUserNotExists(t, s, 999)
}

func TestHelperFunctions_Integration(t *testing.T) {
	s := TestStore(t)

	// Start with empty store
	AssertUserCount(t, s, 0)

	// Create multiple users
	users := CreateMultipleTestUsers(t, s, 5)
	AssertUserCount(t, s, 5)

	// Verify each user exists
	for _, user := range users {
		AssertUserExists(t, s, user.ID)
	}

	// Delete a user
	err := s.DeleteUser(context.Background(), users[0].ID)
	require.NoError(t, err)
	AssertUserCount(t, s, 4)
	AssertUserNotExists(t, s, users[0].ID)

	// Verify remaining users still exist
	for i := 1; i < len(users); i++ {
		AssertUserExists(t, s, users[i].ID)
	}
}
