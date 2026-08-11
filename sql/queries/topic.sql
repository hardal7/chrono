-- name: CreateTopic :exec
INSERT INTO topics(name, created_by_userid)
VALUES($1, $2);
-- name: UpdateTopic :exec
UPDATE topics
SET name = $2, total_time_tracked_seconds = $3, updated_at = now()
WHERE id = $1;
-- name: DeleteTopic :exec
DELETE FROM topics
WHERE id = $1;
-- name: GetTopicByID :one
SELECT * FROM topics
WHERE id = $1;
-- name: GetAllTopics :many
SELECT * FROM topics;
-- name: GetTopicOfUserByName :one
SELECT * FROM topics
WHERE name = $1 AND created_by_userid = $2;
-- name: TrackTopicTime :exec
UPDATE topics
SET 
    total_time_tracked_seconds = total_time_tracked_seconds + sqlc.arg(time_tracked),
    today_time_tracked_seconds = today_time_tracked_seconds + sqlc.arg(time_tracked)
WHERE id = $1 AND created_by_userid = $2;
-- name: ResetTopicTimeTrackedToday :exec
UPDATE topics
SET today_time_tracked_seconds = 0;
-- name: IncreaseStreak :exec
UPDATE topics
SET streak = streak + 1
WHERE id = $1;
-- name: LoseStreak :exec
UPDATE topics
SET streak = 0
WHERE id = $1;
