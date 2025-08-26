-- name: CreateBlogPost :one
INSERT INTO blog_posts (title, body, image_link)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetBlogPost :one
SELECT * FROM blog_posts
WHERE id = $1
LIMIT 1;

-- name: ListBlogPosts :many
SELECT * FROM blog_posts
ORDER BY created_at DESC;

-- name: UpdateBlogPost :one
UPDATE blog_posts
SET title = $2,
    body = $3,
    image_link = $4,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteBlogPost :exec
DELETE FROM blog_posts
WHERE id = $1;
