-- +goose Up
CREATE TABLE addresses (
    address_id   UUID        PRIMARY KEY,
    user_id      UUID        NOT NULL REFERENCES users (user_id) ON DELETE CASCADE,
    label        TEXT        NOT NULL,
    recipient    TEXT        NOT NULL,
    phone        TEXT        NULL,
    line1        TEXT        NOT NULL,
    line2        TEXT        NULL,
    postal_code  TEXT        NOT NULL,
    city         TEXT        NOT NULL,
    country      TEXT        NOT NULL CHECK (country ~ '^[A-Z]{2}$'),
    is_default   BOOLEAN     NOT NULL,
    date_created TIMESTAMPTZ NOT NULL,
    date_updated TIMESTAMPTZ NOT NULL
);

CREATE INDEX addresses_user_id_idx ON addresses (user_id);

-- A user has at most one default address, and no row at all until they need
-- one shipped: collecting in the shop asks for no address.
CREATE UNIQUE INDEX addresses_one_default_idx
    ON addresses (user_id)
    WHERE is_default;

-- +goose Down
DROP TABLE addresses;
