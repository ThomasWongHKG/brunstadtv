-- name: getProgressForProfile :many
SELECT episode_id, progress
FROM "users"."progress"
WHERE profile_id = $1::uuid
  AND episode_id = ANY ($2::int[])
ORDER BY updated_at DESC;

-- name: saveProgress :exec
INSERT INTO "users"."progress" (profile_id, episode_id, progress, updated_at)
VALUES ($1::uuid, $2::int, $3::time, NOW())
ON CONFLICT (profile_id, episode_id) DO UPDATE SET progress = EXCLUDED.progress;

-- name: deleteProgress :exec
DELETE
FROM "users"."progress"
WHERE profile_id = $1::uuid
  AND episode_id = $2::int;
