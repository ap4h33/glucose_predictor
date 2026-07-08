-- +goose Up
CREATE TABLE model_predictions (
    id UUID PRIMARY KEY,
    patient_id UUID REFERENCES users(id) ON DELETE CASCADE,
    model_version VARCHAR NOT NULL,
    time_predicted TIMESTAMPTZ NOT NULL,
    generated_at TIMESTAMPTZ UNIQUE NOT NULL,
    predicted_value NUMERIC(5,2) NOT NULL,
    accuracy FLOAT,
    odu_accuracy FLOAT
);



-- +goose Down
DROP TABLE model_predictions;