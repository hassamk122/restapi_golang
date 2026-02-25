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

-- name: ListUsers :many
SELECT id, username, email, created_at, updated_at 
FROM users
ORDER BY id;

-- name: CreateBlog :one
INSERT INTO blogs(title, content, user_id)
VALUES ($1, $2, $3)
    RETURNING id, title, content, user_id, created_at, updated_at;

-- name: ListBlogs :many
SELECT id, title, content, user_id, created_at, updated_at
FROM blogs
ORDER BY id;