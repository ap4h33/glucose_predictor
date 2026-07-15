-- +goose Up 
CREATE TABLE admins(
    id UUID PRIMARY KEY,
    hospital_id UUID REFERENCES hospitals(id)
);

-- +goose Down
DROP TABLE admins;