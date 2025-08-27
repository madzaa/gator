-- +goose Up
CREATE TABLE posts
(
    id           UUID PRIMARY KEY,
    created_at   TIMESTAMP NOT NULL,
    updated_at   TIMESTAMP NOT NULL,
    title        VARCHAR   NOT NULL,
    url          VARCHAR   NOT NULL UNIQUE,
    description  VARCHAR,
    published_at TIMESTAMP,
    feed_id      UUID      NOT NULL REFERENCES gator.public.feeds (id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE posts;