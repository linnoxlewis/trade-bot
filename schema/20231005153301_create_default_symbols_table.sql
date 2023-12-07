-- +goose Up
CREATE TABLE IF NOT EXISTS default_symbols (
    symbol           VARCHAR(255) PRIMARY KEY,
    created_at       TIMESTAMP WITH TIME ZONE,
    deleted_at       TIMESTAMP WITH TIME ZONE,
    updated_at       TIMESTAMP WITH TIME ZONE);

-- +goose Down
DROP INDEX IF EXISTS default_symbols;

