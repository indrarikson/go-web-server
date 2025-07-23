// Package store provides data access layer interfaces and implementations.
// This file contains mock implementations for testing.
package store

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// UserStore defines the interface for user-related database operations.
// This interface allows for easy mocking in tests.
type UserStore interface {
	CreateUser(ctx context.Context, params CreateUserParams) (User, error)
	GetUser(ctx context.Context, id int64) (User, error)
	ListUsers(ctx context.Context) ([]User, error)
	UpdateUser(ctx context.Context, params UpdateUserParams) (User, error)
	DeactivateUser(ctx context.Context, id int64) error
	DeleteUser(ctx context.Context, id int64) error
	CountUsers(ctx context.Context) (int64, error)
}

// MockUserStore is a mock implementation of UserStore for testing.
type MockUserStore struct {
	mock.Mock
}

func (m *MockUserStore) CreateUser(ctx context.Context, params CreateUserParams) (User, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(User), args.Error(1)
}

func (m *MockUserStore) GetUser(ctx context.Context, id int64) (User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(User), args.Error(1)
}

func (m *MockUserStore) ListUsers(ctx context.Context) ([]User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]User), args.Error(1)
}

func (m *MockUserStore) UpdateUser(ctx context.Context, params UpdateUserParams) (User, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(User), args.Error(1)
}

func (m *MockUserStore) DeactivateUser(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserStore) DeleteUser(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserStore) CountUsers(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}
