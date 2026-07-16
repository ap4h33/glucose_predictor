-- name: CreateHospital :one
INSERT INTO hospitals(id, name)
VALUES ($1, $2)
RETURNING *;

-- name: GetHospital :one
SELECT * 
FROM hospitals
WHERE name=$1;

-- name: GetAllHospitals :many
SELECT *
FROM hospitals;
