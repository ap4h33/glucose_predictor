-- name: GetOduPredictions :many
SELECT glucose, time_predicted FROM odu_predictions
WHERE patient_id=$1 
AND time_predicted>$2
ORDER BY time_predicted ASC;
