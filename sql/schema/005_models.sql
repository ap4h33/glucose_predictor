-- +goose Up
CREATE TABLE models (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    version TEXT NOT NULL
);


-- +goose Down
DROP TABLE models;