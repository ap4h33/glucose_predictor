-- name: AddReading :one
INSERT INTO readings(id, patient_id, time_of_reading, glucose, basal_rate, bolus, carbs, exercise_duration, exercise_intensity)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: SendReadingToModel :exec
UPDATE readings
SET in_the_model=true
WHERE id=$1;

-- name: GetReadings :many
SELECT * FROM readings
WHERE patient_id=$1;

-- name: GetUnseenReadings :many
SELECT * FROM readings
WHERE patient_id=$1
AND in_the_model=false;

-- name: GetLastReadings :many
SELECT * FROM READINGS 
WHERE patient_id=$1
AND in_the_model=false
ORDER BY time_of_reading DESC
LIMIT 12;