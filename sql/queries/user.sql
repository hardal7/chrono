-- name: CreateUser :exec
INSERT INTO users(email, username, password, country)
VALUES($1, $2, $3, $4);
-- name: UpdateUser :exec
UPDATE users
SET username = COALESCE($2, username), password = $3, updated_at = now()
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
-- name: UpdateUserActivity :exec
UPDATE users
SET last_seen_at = now()
WHERE id = $1;
-- name: GetTopUsers :many
SELECT * FROM users
WHERE 
    total_time_tracked_seconds < $1
    AND username ILIKE sqlc.arg(match_name) || '%'
    AND users.hide_user = FALSE
ORDER BY total_time_tracked_seconds DESC
LIMIT $2;
-- name: GetTopUsersLocal :many
SELECT users.* FROM users
JOIN users AS target_user ON target_user.id = $1
WHERE 
    users.country = target_user.country
    AND users.total_time_tracked_seconds < $2
    AND username ILIKE sqlc.arg(match_name) || '%'
    AND users.hide_country = FALSE
    AND users.hide_user = FALSE
ORDER BY users.total_time_tracked_seconds DESC
LIMIT $3;
-- name: TrackUserTime :exec
UPDATE users
SET
    total_time_tracked_seconds = total_time_tracked_seconds + sqlc.arg(time_tracked),
    today_time_tracked_seconds = today_time_tracked_seconds + sqlc.arg(time_tracked)
WHERE id = $1;
-- name: ResetUserTimeTrackedToday :exec
UPDATE users
SET today_time_tracked_seconds = 0;
-- name: GetUsersAll :many
SELECT * FROM users;
-- name: IncreaseStreak :exec
UPDATE users
SET streak = streak + 1
WHERE id = $1;
-- name: LoseStreak :exec
UPDATE users
SET streak = 0
WHERE id = $1;
