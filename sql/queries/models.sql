-- name: AddModel :one
INSERT INTO models (id, name, version)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetAllModels :many
SELECT * FROM models;