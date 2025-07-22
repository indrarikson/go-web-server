package store

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func TestNewStore(t *testing.T) {
	store, err := NewStore(":memory:")
	require.NoError(t, err)
	assert.NotNil(t, store)

	err = store.Close()
	assert.NoError(t, err)
}

func TestStore_InitSchema(t *testing.T) {
	store, err := NewStore(":memory:")
	require.NoError(t, err)
	defer store.Close()

	err = store.InitSchema()
	assert.NoError(t, err)

	// Verify table was created
	var count int
	err = store.DB().QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='users'").Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 1, count)
}

func TestStore_UserOperations(t *testing.T) {
	store, err := NewStore(":memory:")
	require.NoError(t, err)
	defer store.Close()

	err = store.InitSchema()
	require.NoError(t, err)

	ctx := context.Background()

	// Create a user
	user, err := store.CreateUser(ctx, CreateUserParams{
		Email:     "test@example.com",
		Name:      "Test User",
		Bio:       sql.NullString{String: "Test bio", Valid: true},
		AvatarUrl: sql.NullString{},
	})
	require.NoError(t, err)
	assert.NotZero(t, user.ID)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, "Test User", user.Name)

	// Get user
	fetchedUser, err := store.GetUser(ctx, user.ID)
	require.NoError(t, err)
	assert.Equal(t, user.ID, fetchedUser.ID)
	assert.Equal(t, user.Email, fetchedUser.Email)

	// List users
	users, err := store.ListUsers(ctx)
	require.NoError(t, err)
	assert.Len(t, users, 1)

	// Count users
	count, err := store.CountUsers(ctx)
	require.NoError(t, err)
	assert.Equal(t, int64(1), count)

	// Update user
	updatedUser, err := store.UpdateUser(ctx, UpdateUserParams{
		ID:        user.ID,
		Name:      "Updated Name",
		Bio:       sql.NullString{String: "Updated bio", Valid: true},
		AvatarUrl: sql.NullString{},
	})
	require.NoError(t, err)
	assert.Equal(t, "Updated Name", updatedUser.Name)

	// Deactivate user
	err = store.DeactivateUser(ctx, user.ID)
	require.NoError(t, err)

	// Verify user count after deactivation
	count, err = store.CountUsers(ctx)
	require.NoError(t, err)
	assert.Equal(t, int64(0), count)
}
