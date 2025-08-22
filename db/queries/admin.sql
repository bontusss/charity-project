-- name: CreateAdmin :one
INSERT INTO admins (email, password)
VALUES ($1, $2)
RETURNING id, email, created_at;

-- name: GetAdminByEmail :one
SELECT * FROM admins
WHERE email = $1 LIMIT 1;

-- name: DeleteAlladmins :exec
DELETE FROM admins;