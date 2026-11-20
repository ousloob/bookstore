-- +goose Up
CREATE TABLE users (
    user_id             UUID        PRIMARY KEY,
    name                TEXT        NOT NULL CHECK (name <> ''),
    email               TEXT        NOT NULL UNIQUE CHECK (email = lower(email)),
    roles               TEXT[]      NOT NULL CHECK (cardinality(roles) > 0),
    password_hash       TEXT        NOT NULL,
    status              TEXT        NOT NULL,
    date_created        TIMESTAMPTZ NOT NULL,
    date_updated        TIMESTAMPTZ NOT NULL,
    date_status_changed TIMESTAMPTZ NOT NULL,
    date_last_active    TIMESTAMPTZ NULL
);

-- +goose Down
DROP TABLE users;
