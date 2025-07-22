-- name: GetUser :one
SELECT * FROM users WHERE id = ? LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = ? LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users 
WHERE is_active = 1 
ORDER BY created_at DESC;

-- name: ListAllUsers :many
SELECT * FROM users ORDER BY created_at DESC;

-- name: CreateUser :one
INSERT INTO users (email, name, bio, avatar_url) 
VALUES (?, ?, ?, ?)
RETURNING *;

-- name: UpdateUser :one
UPDATE users 
SET name = ?, bio = ?, avatar_url = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING *;

-- name: DeactivateUser :exec
UPDATE users 
SET is_active = 0, updated_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = ?;

-- name: CountUsers :one
SELECT COUNT(*) FROM users WHERE is_active = 1;