-- +goose Up
CREATE TABLE sessions (
    id UUID PRIMARY KEY,
    score INT NOT NULL,
    accuracy DECIMAL NOT NULL,
    created_at TIMESTAMP NOT NULL,
    player_id UUID NOT NULL REFERENCES users ON DELETE CASCADE,
    FOREIGN KEY (player_id) REFERENCES users(id)
);

-- +goose Down
DROP TABLE sessions;