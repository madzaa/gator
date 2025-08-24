-- +goose Up
CREATE TABLE users(
    id uuid PRIMARY KEY,
    created_at  TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL ,
    name VARCHAR(255) NOT NULL UNIQUE
);

-- +goose Down
DROP TABLE users;