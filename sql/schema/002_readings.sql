-- +goose Up
CREATE TABLE readings (
    id UUID PRIMARY KEY,
    patient_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    time_of_reading TIMESTAMPTZ UNIQUE NOT NULL,
    glucose NUMERIC(5,2) NOT NULL,
    basal_rate NUMERIC(5,2) NOT NULL,
    bolus NUMERIC(5,2) NOT NULL,
    carbs NUMERIC(5,2) NOT NULL,
    exercise_duration INT,
    exercise_intensity INT,
    in_the_model BOOLEAN DEFAULT false
);

-- +goose Down
DROP TABLE readings;