-- name: FindAllUsers :many
SELECT * FROM users;

-- name: FindUserByEmail :one
SELECT id, created_at, name, email, encrypted_password, activated, provider 
FROM users 
WHERE email = $1;

-- name: CreateUser :one
INSERT INTO users (name, email, encrypted_password, activated, provider) 
VALUES ($1, $2, $3, $4, $5) 
RETURNING id, created_at, version;