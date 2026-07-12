-- +goose Up
CREATE TABLE model_predictions (
    id UUID PRIMARY KEY,
    patient_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    glucose NUMERIC(5,2) NOT NULL, 
    model_version VARCHAR NOT NULL,
    time_predicted TIMESTAMPTZ NOT NULL,
    generated_at TIMESTAMPTZ NOT NULL,
    predicted_value NUMERIC(5,2) NOT NULL,
    accuracy FLOAT,
    odu_accuracy FLOAT
);

-- +goose Down
DROP TABLE model_predictions;