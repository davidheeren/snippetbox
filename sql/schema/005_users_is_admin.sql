-- +goose Up
ALTER TABLE users ADD
is_admin BOOLEAN NOT NULL DEFAULT FALSE;

-- +goose Down
ALTER TABLE users DROP
COLUMN is_admin;
