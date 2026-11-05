-- +goose Up
CREATE TABLE users (
    user_id       UUID        PRIMARY KEY,
    name          TEXT        NOT NULL,
    email         TEXT        NOT NULL UNIQUE,
    roles         TEXT[]      NOT NULL,
    password_hash BYTEA       NOT NULL,
    date_created  TIMESTAMPTZ NOT NULL,
    date_updated  TIMESTAMPTZ NOT NULL
);

-- +goose Down
DROP TABLE users;
