package store

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestMockUserStore_CreateUser(t *testing.T) {
	mockStore := &MockUserStore{}
	ctx := context.Background()

	params := CreateUserParams{
		Email:     "test@example.com",
		Name:      "Test User",
		Bio:       sql.NullString{String: "Test bio", Valid: true},
		AvatarUrl: sql.NullString{},
	}

	expectedUser := User{
		ID:    1,
		Email: "test@example.com",
		Name:  "Test User",
		Bio:   sql.NullString{String: "Test bio", Valid: true},
	}

	mockStore.On("CreateUser", ctx, params).Return(expectedUser, nil)

	user, err := mockStore.CreateUser(ctx, params)

	require.NoError(t, err)
	assert.Equal(t, expectedUser, user)
	mockStore.AssertExpectations(t)
}

func TestMockUserStore_CreateUser_Error(t *testing.T) {
	mockStore := &MockUserStore{}
	ctx := context.Background()

	params := CreateUserParams{
		Email: "test@example.com",
		Name:  "Test User",
	}

	expectedError := errors.New("create user failed")
	mockStore.On("CreateUser", ctx, params).Return(User{}, expectedError)

	_, err := mockStore.CreateUser(ctx, params)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	mockStore.AssertExpectations(t)
}

func TestMockUserStore_GetUser(t *testing.T) {
	mockStore := &MockUserStore{}
	ctx := context.Background()
	userID := int64(1)

	expectedUser := User{
		ID:    userID,
		Email: "test@example.com",
		Name:  "Test User",
	}

	mockStore.On("GetUser", ctx, userID).Return(expectedUser, nil)

	user, err := mockStore.GetUser(ctx, userID)

	require.NoError(t, err)
	assert.Equal(t, expectedUser, user)
	mockStore.AssertExpectations(t)
}

func TestMockUserStore_ListUsers(t *testing.T) {
	mockStore := &MockUserStore{}
	ctx := context.Background()

	expectedUsers := []User{
		{ID: 1, Email: "user1@example.com", Name: "User 1"},
		{ID: 2, Email: "user2@example.com", Name: "User 2"},
	}

	mockStore.On("ListUsers", ctx).Return(expectedUsers, nil)

	users, err := mockStore.ListUsers(ctx)

	require.NoError(t, err)
	assert.Equal(t, expectedUsers, users)
	mockStore.AssertExpectations(t)
}

func TestMockUserStore_UpdateUser(t *testing.T) {
	mockStore := &MockUserStore{}
	ctx := context.Background()

	params := UpdateUserParams{
		ID:        1,
		Name:      "Updated Name",
		Bio:       sql.NullString{String: "Updated bio", Valid: true},
		AvatarUrl: sql.NullString{},
	}

	expectedUser := User{
		ID:    1,
		Email: "test@example.com",
		Name:  "Updated Name",
		Bio:   sql.NullString{String: "Updated bio", Valid: true},
	}

	mockStore.On("UpdateUser", ctx, params).Return(expectedUser, nil)

	user, err := mockStore.UpdateUser(ctx, params)

	require.NoError(t, err)
	assert.Equal(t, expectedUser, user)
	mockStore.AssertExpectations(t)
}

func TestMockUserStore_DeactivateUser(t *testing.T) {
	mockStore := &MockUserStore{}
	ctx := context.Background()
	userID := int64(1)

	mockStore.On("DeactivateUser", ctx, userID).Return(nil)

	err := mockStore.DeactivateUser(ctx, userID)

	assert.NoError(t, err)
	mockStore.AssertExpectations(t)
}

func TestMockUserStore_DeleteUser(t *testing.T) {
	mockStore := &MockUserStore{}
	ctx := context.Background()
	userID := int64(1)

	mockStore.On("DeleteUser", ctx, userID).Return(nil)

	err := mockStore.DeleteUser(ctx, userID)

	assert.NoError(t, err)
	mockStore.AssertExpectations(t)
}

func TestMockUserStore_CountUsers(t *testing.T) {
	mockStore := &MockUserStore{}
	ctx := context.Background()
	expectedCount := int64(5)

	mockStore.On("CountUsers", ctx).Return(expectedCount, nil)

	count, err := mockStore.CountUsers(ctx)

	require.NoError(t, err)
	assert.Equal(t, expectedCount, count)
	mockStore.AssertExpectations(t)
}

func TestMockUserStore_AllMethods_WithErrors(t *testing.T) {
	mockStore := &MockUserStore{}
	ctx := context.Background()
	expectedError := errors.New("database error")

	// Test all methods returning errors
	mockStore.On("CreateUser", ctx, mock.Anything).Return(User{}, expectedError)
	mockStore.On("GetUser", ctx, mock.Anything).Return(User{}, expectedError)
	mockStore.On("ListUsers", ctx).Return([]User{}, expectedError)
	mockStore.On("UpdateUser", ctx, mock.Anything).Return(User{}, expectedError)
	mockStore.On("DeactivateUser", ctx, mock.Anything).Return(expectedError)
	mockStore.On("DeleteUser", ctx, mock.Anything).Return(expectedError)
	mockStore.On("CountUsers", ctx).Return(int64(0), expectedError)

	// Test CreateUser error
	_, err := mockStore.CreateUser(ctx, CreateUserParams{})
	assert.Equal(t, expectedError, err)

	// Test GetUser error
	_, err = mockStore.GetUser(ctx, 1)
	assert.Equal(t, expectedError, err)

	// Test ListUsers error
	_, err = mockStore.ListUsers(ctx)
	assert.Equal(t, expectedError, err)

	// Test UpdateUser error
	_, err = mockStore.UpdateUser(ctx, UpdateUserParams{})
	assert.Equal(t, expectedError, err)

	// Test DeactivateUser error
	err = mockStore.DeactivateUser(ctx, 1)
	assert.Equal(t, expectedError, err)

	// Test DeleteUser error
	err = mockStore.DeleteUser(ctx, 1)
	assert.Equal(t, expectedError, err)

	// Test CountUsers error
	_, err = mockStore.CountUsers(ctx)
	assert.Equal(t, expectedError, err)

	mockStore.AssertExpectations(t)
}
