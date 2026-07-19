-- name: GetRecommendation :one
SELECT * 
FROM recommendations
WHERE patient_id=$1
ORDER BY created_at DESC
LIMIT 1;