-- +goose Up
CREATE TABLE odu_predictions (
    id UUID PRIMARY KEY,
    patient_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    glucose NUMERIC(5,2) NOT NULL,
    time_predicted TIMESTAMPTZ NOT NULL,
    generated_at TIMESTAMPTZ NOT NULL,
    predicted_value NUMERIC(5,2) NOT NULL,
    accuracy FLOAT
);

-- +goose Down
DROP TABLE odu_predictions;