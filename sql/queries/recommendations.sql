-- name: GetRecommendation :one
SELECT * 
FROM recommendations
WHERE patient_id=$1
ORDER BY created_at DESC
LIMIT 1;

-- name: AddRecommendation :one
INSERT INTO recommendations(id, patient_id, is_safe, action, message, created_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;