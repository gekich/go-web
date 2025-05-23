-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: ExistsUserByEmail :one
SELECT EXISTS (
    SELECT 1 FROM users WHERE email = $1
) AS exists;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY name;

-- name: CreateUser :one
INSERT INTO users (
    name, email, password, bio
) VALUES (
          $1, $2, $3, $4
         )
RETURNING *;

-- name: UpdateUser :exec
UPDATE users
set name = $2,
    bio = $3
WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;