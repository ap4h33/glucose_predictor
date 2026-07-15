-- name: GetModelPredictions :many
SELECT glucose_predicted, time_predicted FROM predictions
WHERE model_id=$1
AND patient_id=$2 
AND time_predicted>$3
ORDER BY time_predicted ASC;