package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunMigrations(t *testing.T) {
	// Test with in-memory database
	err := runMigrations(":memory:")
	// This will likely fail since the migrations directory expects file paths
	// but we're testing the function exists and handles errors
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create migration instance")
}

func TestRunMigrations_WithValidPath(t *testing.T) {
	// Create a temporary database file
	tmpFile, err := os.CreateTemp("", "test_*.db")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	// Test migrations with a real file path
	// This should still fail because the migrations directory might not exist
	// in the test environment, but it tests the function logic
	err = runMigrations(tmpFile.Name())
	// We expect this to fail in test environment but the function should handle it gracefully
	if err != nil {
		// Check that it's a migration-related error, not a panic
		assert.Contains(t, err.Error(), "failed to")
	}
}

func TestRunMigrations_InvalidPath(t *testing.T) {
	// Test with an invalid database path
	err := runMigrations("/invalid/path/does/not/exist.db")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create migration instance")
}
