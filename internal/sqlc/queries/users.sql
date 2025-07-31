-- EngLog User Management Queries
-- User authentication and profile management

-- name: CreateUser :one
INSERT INTO users (
    email, password_hash, first_name, last_name, timezone, preferences
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1;

-- name: UpdateUserLastLogin :exec
UPDATE users
SET last_login_at = NOW(), updated_at = NOW()
WHERE id = $1;

-- name: UpdateUserProfile :one
UPDATE users
SET first_name = $2, last_name = $3, timezone = $4,
    preferences = $5, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdateUserPassword :exec
UPDATE users
SET password_hash = $2, updated_at = NOW()
WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: GetUserCount :one
SELECT COUNT(*) FROM users;

-- name: GetRecentUsers :many
SELECT id, email, first_name, last_name, created_at
FROM users
ORDER BY created_at DESC
LIMIT $1;
