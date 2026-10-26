-- +goose Up
DELETE FROM snippets;

ALTER TABLE snippets ADD
user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE;

-- +goose Down
ALTER TABLE snippets DROP
COLUMN user_id;
