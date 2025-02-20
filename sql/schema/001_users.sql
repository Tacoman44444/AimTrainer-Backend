-- +goose Up
CREATE TABLE users (
    id UUID PRIMARY KEY,
    username TEXT NOT NULL,
    hashed_password TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    UNIQUE(username)
);

-- +goose Down
DROP TABLE users;