-- name: CreateHospital :one
INSERT INTO hospitals(id, name)
VALUES ($1, $2)
RETURNING *;

-- name: GetHospitalsByName :many
SELECT * 
FROM hospitals
WHERE name=$1;

-- name: GetAllHospitals :many
SELECT *
FROM hospitals;
