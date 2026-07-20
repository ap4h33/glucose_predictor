-- +goose Up
CREATE TABLE model_accuracy(
    id UUID PRIMARY KEY,
    model_id UUID NOT NULL REFERENCES models(id),
    rmse NUMERIC(10,6),
    mae NUMERIC(10,6),
    assesment_date TIMESTAMPTZ NOT NULL
);

-- +goose Down
DROP TABLE model_accuracy;