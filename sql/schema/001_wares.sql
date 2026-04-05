-- +goose Up
PRAGMA foreign_keys = ON;
CREATE TABLE wares(
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    name TEXT NOT NULL UNIQUE,
    volume INTEGER NOT NULL CHECK(volume >= 0),
    price INTEGER NOT NULL CHECK(price >=0)
);

-- +goose Down
DROP TABLE wares;
PRAGMA foreign_keys = OFF;
