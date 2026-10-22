-- name: CreateSnippet :one
INSERT INTO snippets (id, created_at, updated_at, expires_at, title, content)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2,
    $3
) RETURNING *;

-- name: GetSnippet :one
SELECT * FROM snippets
WHERE id = $1;

-- name: DeleteSnippet :exec
DELETE FROM snippets
WHERE id = $1;

-- name: GetSnippets :many
SELECT * FROM snippets;
