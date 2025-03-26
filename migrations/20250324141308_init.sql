-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS users (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    name varchar(255),
    email citext UNIQUE NOT NULL,
    encrypted_password bytea NOT NULL,
    activated bool NOT NULL,
    version integer NOT NULL DEFAULT 1
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
