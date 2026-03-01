-- name: CreateUser :one
INSERT INTO users(username, email, password)
VALUES ($1, $2, $3)
    RETURNING id, username, email, created_at, updated_at;

-- name: GetUser :one
SELECT id, username, email, created_at, updated_at
FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT id, username, email, created_at, updated_at
FROM users
WHERE email = $1;


-- name: GetUserByEmailIncludingPassword :one
SELECT id, username, email, password, created_at, updated_at
FROM users
WHERE email = $1;

-- name: ListUsers :many
SELECT id, username, email, created_at, updated_at 
FROM users
ORDER BY id
LIMIT $1 OFFSET $2;

-- name: GetTotalUserCound :one
SELECT COUNT(*) AS total FROM users;

-- name: CreateUserProfile :one
INSERT INTO user_profiles (user_id, profile_image)
VALUES ($1, $2)
RETURNING id, user_id, profile_image;

-- name: GetUserProfileByUserId :one
SELECT id, user_id, profile_image, created_at, updated_at
FROM user_profiles
WHERE user_id = $1;