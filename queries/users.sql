-- name: CreateUser :one
INSERT INTO users (
    user_id,
    first_name,
    last_name,
    email,
    phone,
    age,
    status
)
VALUES ( $1, $2, $3, $4, $5, $6, $7)
    RETURNING user_id, email, status;

-- name: GetUser :one
SELECT * FROM users WHERE user_id = $1 LIMIT 1;

-- name: GetAllUsers :many
SELECT * FROM users ORDER BY first_name;

-- name: UpdateUser :one
UPDATE users
SET
    first_name = COALESCE($2, first_name),
    last_name  = COALESCE($3, last_name),
    email      = COALESCE($4, email),
    phone      = COALESCE($5, phone),
    age        = COALESCE($6, age),
    status     = COALESCE($7, status)
WHERE user_id = $1
    RETURNING user_id, email, status;

-- name: DeleteUser :one
DELETE FROM users
WHERE user_id = $1
    RETURNING user_id;