-- +goose Up
CREATE VIEW snippets_with_user AS
SELECT s.*, u.name AS user_name
FROM snippets s
JOIN users u ON u.id = s.user_id;

-- +goose Down
DROP VIEW snippets_with_user;
