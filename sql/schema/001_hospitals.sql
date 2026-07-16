-- +goose Up
CREATE TABLE hospitals(
    id UUID PRIMARY KEY,
    name TEXT UNIQUE NOT NULL
);

-- +goose Down
DROP TABLE hospitals;