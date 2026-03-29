-- +goose Up
CREATE TABLE wares(
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    name TEXT UNIQUE NOT NULL,
    category TEXT,
    min_price INTEGER NOT NULL,
    max_price INTEGER NOT NULL,
    volume INTEGER NOT NULL CHECK(volume >= 0),
    CONSTRAINT price_check CHECK(min_price <= max_price)
) STRICT;

-- +goose Down
DROP TABLE wares;
