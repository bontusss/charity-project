-- name: CreateProject :one
INSERT INTO projects (name)
VALUES ($1)
RETURNING *;

-- name: GetProject :one
SELECT * FROM projects
WHERE id = $1
LIMIT 1;

-- name: ListProjects :many
SELECT * FROM projects
ORDER BY created_at DESC;

-- name: UpdateProject :one
UPDATE projects
SET name = $2,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteProject :exec
DELETE FROM projects
WHERE id = $1;

-- name: CreateProjectBefore :one
INSERT INTO project_before (project_id, body, estimated_target, video_link)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetProjectBefore :one
SELECT * FROM project_before
WHERE project_id = $1
LIMIT 1;

-- name: UpdateProjectBefore :one
UPDATE project_before
SET body = $2,
    estimated_target = $3,
    video_link = $4,
    updated_at = NOW()
WHERE project_id = $1
RETURNING *;

-- name: DeleteProjectBefore :exec
DELETE FROM project_before
WHERE project_id = $1;

-- name: CreateProjectAfter :one
INSERT INTO project_after (project_id, body, project_cost, video_link)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetProjectAfter :one
SELECT * FROM project_after
WHERE project_id = $1
LIMIT 1;

-- name: UpdateProjectAfter :one
UPDATE project_after
SET body = $2,
    project_cost = $3,
    video_link = $4,
    updated_at = NOW()
WHERE project_id = $1
RETURNING *;

-- name: DeleteProjectAfter :exec
DELETE FROM project_after
WHERE project_id = $1;

-- name: AddProjectImage :one
INSERT INTO project_images (project_id, phase, image_link)
VALUES ($1, $2, $3)
RETURNING *;

-- name: ListProjectImages :many
SELECT * FROM project_images
WHERE project_id = $1
ORDER BY created_at ASC;

-- name: ListProjectImagesByPhase :many
SELECT * FROM project_images
WHERE project_id = $1
  AND phase = $2
ORDER BY created_at ASC;

-- name: DeleteProjectImage :exec
DELETE FROM project_images
WHERE id = $1;

-- name: DeleteAllProjectImages :exec
DELETE FROM project_images
WHERE project_id = $1;
