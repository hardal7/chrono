-- name: CreateUser :exec
INSERT INTO users(email, username, password)
VALUES($1, $2, $3);
-- name: UpdateUser :exec
UPDATE users
SET email = $2, username = $3, password = $4, updated_at = now()
WHERE id = $1;
-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;
-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = $1;
-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;
-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1;
-- name: GetTopUsers :many
SELECT * FROM users
WHERE total_time_tracked_seconds < $1 
ORDER BY total_time_tracked_seconds DESC
LIMIT $2;
-- name: TrackUserTime :exec
UPDATE users
SET
  total_time_tracked_seconds = total_time_tracked_seconds + $2,
  today_time_tracked_seconds = today_time_tracked_seconds + $2
WHERE id = $1;
-- name: ResetTimeTrackedToday :exec
UPDATE users
SET today_time_tracked_seconds = 0;
-- name: GetAvatarFromUserID :one
SELECT * FROM avatars
WHERE user_id = $1;
-- name: CreateAvatar :exec
INSERT INTO avatars (user_id)
VALUES ($1);
