package store

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

// Store provides all functions to execute db queries.
type Store struct {
	*Queries // Embed sqlc-generated queries
	db       *sql.DB
}

// NewStore creates a new store instance with database connection.
func NewStore(databaseURL string) (*Store, error) {
	db, err := sql.Open("sqlite", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Store{
		db:      db,
		Queries: New(db),
	}, nil
}

// NewStoreWithDB creates a new store instance with an existing database connection.
func NewStoreWithDB(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// Close closes the database connection.
func (s *Store) Close() error {
	return s.db.Close()
}

// DB returns the underlying database connection for advanced operations.
func (s *Store) DB() *sql.DB {
	return s.db
}

// InitSchema initializes the database schema using the schema.sql file.
// This is kept here for compatibility, but migrations are preferred.
func (s *Store) InitSchema() error {
	schema := `
		-- Enhanced users table with additional fields
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			email TEXT UNIQUE NOT NULL,
			name TEXT NOT NULL,
			avatar_url TEXT,
			bio TEXT,
			is_active BOOLEAN DEFAULT 1,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);

		-- Index for faster email lookups
		CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

		-- Index for active users
		CREATE INDEX IF NOT EXISTS idx_users_active ON users(is_active);
	`

	_, err := s.db.Exec(schema)
	if err != nil {
		return fmt.Errorf("failed to initialize schema: %w", err)
	}

	return nil
}
