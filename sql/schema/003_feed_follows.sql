-- +goose Up
CREATE TABLE feed_follows
(
    id         uuid PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    user_id    uuid      NOT NULL REFERENCES gator.public.users (id) ON DELETE CASCADE,
    feed_id    uuid      NOT NULL REFERENCES gator.public.feeds (id) ON DELETE CASCADE,
    CONSTRAINT user_feed_unique UNIQUE (user_id, feed_id)
);

-- +goose Down
DROP TABLE feed_follows;