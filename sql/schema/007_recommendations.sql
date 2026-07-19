-- +goose Up
CREATE TABLE recommendations(
    id UUID PRIMARY KEY,
    patient_id INT NOT NULL REFERENCES users(id),
    is_safe BOOLEAN NOT NULL,
    action TEXT,
    message TEXT,
    created_at TIMESTAMPTZ
);

-- +goose Down
DROP TABLE recommendations;