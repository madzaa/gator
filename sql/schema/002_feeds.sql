-- +goose Up
CREATE TABLE feeds
(
    id         uuid PRIMARY KEY,
    created_at TIMESTAMP    NOT NULL,
    updated_at TIMESTAMP    NOT NULL,
    name       VARCHAR(255) NOT NULL,
    url        VARCHAR(255) NOT NULL UNIQUE,
    user_id    uuid         NOT NULL REFERENCES gator.public.users (id) ON DELETE CASCADE
);
-- Use an ON DELETE CASCADE constraint on the user_id foreign key so that if a user is deleted, all of their feeds are automatically deleted as well.
-- This will ensure we have no orphaned records and that deleting the users in the reset command also deletes all of their feeds.

-- +goose Down
DROP TABLE feeds;