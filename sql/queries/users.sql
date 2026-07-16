-- name: CreateUser :one
INSERT INTO users (id, hospital_id, name)
SELECT 
    CASE WHEN $1::INT = 0 THEN nextval(pg_get_serial_sequence('users', 'id')) ELSE $1::INT END, 
    $2, $3
RETURNING *;

-- name: GetUsers :many
SELECT *
FROM users
WHERE users.hospital_id=(
    SELECT admins.hospital_id
    FROM admins
    WHERE admins.id=$1
);

-- name: GetUser :one
SELECT *
FROM users
WHERE id=$1;