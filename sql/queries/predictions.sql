-- name: GetModelPredictions :many
SELECT glucose_predicted, time_predicted FROM predictions
WHERE model_id=$1
AND patient_id=$2 
AND time_predicted>$3
ORDER BY time_predicted ASC;


-- name: AddPrediction :one
INSERT INTO predictions (
    id,
    model_id, 
    patient_id,
    glucose_predicted,
    time_predicted,
    generated_at
)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;