-- name: CreateSnippet :one
INSERT INTO snippets (id, created_at, updated_at, expires_at, title, content, user_id)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2,
    $3,
    $4
) RETURNING *;

-- name: GetSnippet :one
SELECT * FROM snippets
WHERE id = $1;

-- name: GetSnippetWithUser :one
SELECT s.*, u.name AS user_name
FROM snippets s
JOIN users u ON u.id = s.user_id
WHERE s.id = $1;

-- name: DeleteSnippetByUser :execrows
DELETE FROM snippets
WHERE id = $1 AND user_id = $2;

-- name: DeleteSnippet :execrows
DELETE FROM snippets
WHERE id = $1;

-- name: GetSnippets :many
SELECT * FROM snippets;

-- name: GetSnippetsByUser :many
SELECT s.*, u.name AS user_name
FROM snippets s
JOIN users u ON u.id = s.user_id;
