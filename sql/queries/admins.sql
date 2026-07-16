-- name: CreateAdmin :one
INSERT INTO admins(id, hospital_id)
VALUES ($1, $2)
RETURNING *;

-- name: GetAdmin :one
SELECT *
FROM admins
WHERE id=$1;