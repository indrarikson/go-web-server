-- Initial schema for users table
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

-- Insert sample data for development
INSERT INTO users (email, name, bio) VALUES 
    ('alice@example.com', 'Alice Johnson', 'Full-stack developer with a passion for Go'),
    ('bob@example.com', 'Bob Smith', 'DevOps engineer who loves automation'),
    ('charlie@example.com', 'Charlie Brown', 'UI/UX designer focused on user experience')
ON CONFLICT(email) DO NOTHING;