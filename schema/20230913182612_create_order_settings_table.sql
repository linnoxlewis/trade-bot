-- +goose Up
CREATE TABLE IF NOT EXISTS order_settings (
    id               BIGSERIAL PRIMARY KEY,
    order_id         INTEGER NOT NULL,
    tp_percent       VARCHAR(255) DEFAULT NULL,
    sl_percent       VARCHAR(255) DEFAULT NULL,
    ts               VARCHAR(255) DEFAULT NULL,
    tp_price         VARCHAR(255) DEFAULT NULL,
    sl_price         VARCHAR(255) DEFAULT NULL,
    tp_type          VARCHAR(255) DEFAULT NULL,
    sl_type          VARCHAR(255) DEFAULT NULL,
    created_at       TIMESTAMP WITH TIME ZONE,
    deleted_at       TIMESTAMP WITH TIME ZONE,
    updated_at       TIMESTAMP WITH TIME ZONE,
    FOREIGN KEY (order_id) REFERENCES orders (id) ON DELETE CASCADE);

CREATE UNIQUE INDEX uniq_order_id_value ON order_settings (order_id) WHERE (deleted_at IS NULL);

-- +goose Down
DROP INDEX IF EXISTS uniq_order_id_value;
DROP TABLE IF EXISTS order_settings;
