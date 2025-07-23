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

	// Delete user
	err = store.DeleteUser(ctx, user.ID)
	require.NoError(t, err)

	// Verify user is deleted
	_, err = store.GetUser(ctx, user.ID)
	assert.Error(t, err)
}

func TestNewStore_InvalidDatabaseURL(t *testing.T) {
	_, err := NewStore("/invalid/path/that/does/not/exist/test.db")
	assert.Error(t, err)
	// SQLite driver may return different error messages, so check for any error
	assert.NotEmpty(t, err.Error())
}

func TestNewStoreWithDB(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer db.Close()

	store := NewStoreWithDB(db)
	assert.NotNil(t, store)
	assert.Equal(t, db, store.DB())
}

func TestStore_DB(t *testing.T) {
	store, err := NewStore(":memory:")
	require.NoError(t, err)
	defer store.Close()

	db := store.DB()
	assert.NotNil(t, db)

	// Test that we can use the DB connection
	err = db.Ping()
	assert.NoError(t, err)
}

func TestStore_Close(t *testing.T) {
	store, err := NewStore(":memory:")
	require.NoError(t, err)

	err = store.Close()
	assert.NoError(t, err)

	// Verify connection is closed by attempting to ping
	err = store.DB().Ping()
	assert.Error(t, err)
}

func TestStore_InitSchema_ErrorHandling(t *testing.T) {
	// Create a store with a closed database to trigger an error
	store, err := NewStore(":memory:")
	require.NoError(t, err)

	// Close the database before calling InitSchema
	store.Close()

	err = store.InitSchema()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to initialize schema")
}

func TestQueries_WithTx(t *testing.T) {
	store, err := NewStore(":memory:")
	require.NoError(t, err)
	defer store.Close()

	err = store.InitSchema()
	require.NoError(t, err)

	// Test WithTx function
	tx, err := store.DB().Begin()
	require.NoError(t, err)
	defer tx.Rollback()

	queries := store.WithTx(tx)
	assert.NotNil(t, queries)

	// Test that we can use queries with transaction
	ctx := context.Background()
	user, err := queries.CreateUser(ctx, CreateUserParams{
		Email:     "tx@example.com",
		Name:      "TX User",
		Bio:       sql.NullString{String: "TX bio", Valid: true},
		AvatarUrl: sql.NullString{},
	})
	require.NoError(t, err)
	assert.NotZero(t, user.ID)
}

func TestQueries_GetUserByEmail(t *testing.T) {
	store, err := NewStore(":memory:")
	require.NoError(t, err)
	defer store.Close()

	err = store.InitSchema()
	require.NoError(t, err)

	ctx := context.Background()

	// Create a user first
	user, err := store.CreateUser(ctx, CreateUserParams{
		Email:     "email@example.com",
		Name:      "Email User",
		Bio:       sql.NullString{String: "Email bio", Valid: true},
		AvatarUrl: sql.NullString{},
	})
	require.NoError(t, err)

	// Test GetUserByEmail
	fetchedUser, err := store.GetUserByEmail(ctx, "email@example.com")
	require.NoError(t, err)
	assert.Equal(t, user.ID, fetchedUser.ID)
	assert.Equal(t, user.Email, fetchedUser.Email)

	// Test with non-existent email
	_, err = store.GetUserByEmail(ctx, "nonexistent@example.com")
	assert.Error(t, err)
}

func TestQueries_ListAllUsers(t *testing.T) {
	store, err := NewStore(":memory:")
	require.NoError(t, err)
	defer store.Close()

	err = store.InitSchema()
	require.NoError(t, err)

	ctx := context.Background()

	// Create active and inactive users
	activeUser, err := store.CreateUser(ctx, CreateUserParams{
		Email:     "active@example.com",
		Name:      "Active User",
		Bio:       sql.NullString{String: "Active bio", Valid: true},
		AvatarUrl: sql.NullString{},
	})
	require.NoError(t, err)

	inactiveUser, err := store.CreateUser(ctx, CreateUserParams{
		Email:     "inactive@example.com",
		Name:      "Inactive User",
		Bio:       sql.NullString{String: "Inactive bio", Valid: true},
		AvatarUrl: sql.NullString{},
	})
	require.NoError(t, err)

	// Deactivate one user
	err = store.DeactivateUser(ctx, inactiveUser.ID)
	require.NoError(t, err)

	// Test ListAllUsers (should include both active and inactive)
	allUsers, err := store.ListAllUsers(ctx)
	require.NoError(t, err)
	assert.Len(t, allUsers, 2)

	// Test ListUsers (should only include active)
	activeUsers, err := store.ListUsers(ctx)
	require.NoError(t, err)
	assert.Len(t, activeUsers, 1)
	assert.Equal(t, activeUser.ID, activeUsers[0].ID)
}
