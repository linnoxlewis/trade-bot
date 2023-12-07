-- +goose Up
CREATE TABLE IF NOT EXISTS users
(
    id INTEGER NOT NULL PRIMARY KEY,
    username VARCHAR(255),
    is_admin BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- +goose Down
DROP TABLE IF EXISTS users;
