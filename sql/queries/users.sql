-- name: CreateUser :one
INSERT INTO users (id, name)
SELECT 
    CASE WHEN $1::INT = 0 THEN nextval(pg_get_serial_sequence('users', 'id')) ELSE $1::INT END, 
    $2
RETURNING *;