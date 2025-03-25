-- name: CreateUser :one
INSERT INTO users (
    name,
    email,
    password
) VALUES (
    sqlc.arg(name),
    sqlc.arg(email),
    sqlc.arg(password)
) RETURNING *;
-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = sqlc.arg(email);
-- name: GetUserByID :one
SELECT * FROM users WHERE id = sqlc.arg(id);
-- name: UpdateUser :one
UPDATE users SET
    name = COALESCE(name, sqlc.arg(name)),
    email = COALESCE(email, sqlc.arg(email)),
    password = COALESCE(password, sqlc.arg(password)),
    updated_at = CURRENT_TIMESTAMP
WHERE id = sqlc.arg(id)
RETURNING *;
-- name: DeleteUser :exec
DELETE FROM users WHERE id = sqlc.arg(id);


