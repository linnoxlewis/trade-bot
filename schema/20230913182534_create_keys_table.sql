-- +goose Up
CREATE TABLE IF NOT EXISTS user_api_keys
(
    user_id integer NOT NULL,
    exchange varchar(255) NOT NULL,
    pub_key text NOT NULL,
    priv_key text NOT NULL,
    passphrase text NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,
    PRIMARY KEY(user_id, exchange)
);

CREATE UNIQUE INDEX uniq_keys_value ON user_api_keys (pub_key, priv_key) WHERE (deleted_at IS NULL);

-- +goose Down

