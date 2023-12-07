-- +goose Up
CREATE TABLE IF NOT EXISTS orders
(
    id              BIGSERIAL PRIMARY KEY,
    user_id         INTEGER NOT NULL,
    exchange        VARCHAR(255) NOT NULL,
    symbol          VARCHAR(255) NOT NULL,
    status          VARCHAR(255) NOT NULL,
    side            VARCHAR(255) NOT NULL,
    order_type      VARCHAR(255) NOT NULL,
    quantity        VARCHAR(255) NOT NULL,
    price           VARCHAR(255) DEFAULT NULL,
    exec_order_id    VARCHAR(255) NOT NULL,
    time_in_force   VARCHAR(255),
    tp_sl           VARChAR(255) NOT NULL,
    created_at      TIMESTAMP WITH TIME ZONE,
    deleted_at      TIMESTAMP WITH TIME ZONE,
    updated_at      TIMESTAMP WITH TIME ZONE
    );

CREATE INDEX IF NOT EXISTS "symbol_index"  ON "orders"("symbol", "exchange");
CREATE INDEX IF NOT EXISTS "id_symbol_index"  ON "orders"("id","symbol", "exchange");

-- +goose Down
DROP INDEX IF EXISTS "symbol_index";
DROP INDEX IF EXISTS "id_symbol_index";
DROP TABLE IF EXISTS virtual_orders;