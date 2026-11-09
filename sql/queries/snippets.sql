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

-- name: DeleteSnippetByUser :execrows
DELETE FROM snippets
WHERE id = $1 AND user_id = $2;

-- name: DeleteSnippet :execrows
DELETE FROM snippets
WHERE id = $1;

-- name: GetSnippetWithUser :one
SELECT * FROM snippets_with_user
WHERE id = $1;

-- name: GetSnippetsWithUser :many
SELECT * FROM snippets_with_user
ORDER BY created_at DESC;

-- name: SearchSnippetsWithUser :many
SELECT * FROM snippets_with_user
WHERE title ILIKE ALL (
    SELECT '%' || word || '%'
    FROM unnest(regexp_split_to_array(trim($1::text), '\s+')) AS word
)
ORDER BY created_at DESC;
