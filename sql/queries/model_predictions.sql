-- name: GetModelPredictions :many
SELECT glucose, time_predicted FROM model_predictions
WHERE patient_id=$1 
AND time_predicted>$2
ORDER BY time_predicted ASC;