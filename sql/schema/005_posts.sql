-- +goose Up
CREATE TABLE posts(
    id uuid PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    title TEXT NOT NULL UNIQUE,
    url TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL UNIQUE,
    published_at TIMESTAMP NOT NULL,
    feed_id uuid NOT NULL,
    CONSTRAINT fk_posts_feed FOREIGN KEY (feed_id)
    REFERENCES feeds(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE posts;