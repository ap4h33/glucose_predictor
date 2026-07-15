-- +goose Up
CREATE TABLE predictions(
    id UUID PRIMARY KEY,
    model_id UUID NOT NULL REFERENCES models(id),
    patient_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    glucose_predicted NUMERIC(5,2) not null,
    time_predicted TIMESTAMPTZ NOT NULL,
    generated_at TIMESTAMPTZ NOT NULL
);

-- +goose Down
DROP TABLE predictions;