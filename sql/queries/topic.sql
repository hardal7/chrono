-- name: CreateTopic :exec
INSERT INTO topics(name, owner_id)
VALUES($1, $2);
-- name: UpdateTopic :exec
UPDATE topics
SET name = $2, updated_at = now()
WHERE id = $1;
-- name: DeleteTopic :exec
DELETE FROM topics
WHERE id = $1;
-- name: GetTopicByID :one
SELECT * FROM topics
WHERE id = $1;
-- name: GetTopicsAll :many
SELECT * FROM topics;
-- name: GetTopicsByOwner :many
SELECT * FROM topics
WHERE owner_id = $1;
-- name: GetTopicByOwnerAndName :one
SELECT * FROM topics
WHERE name = $1 AND owner_id = $2;
-- name: TrackTopicTime :exec
UPDATE topics
SET 
    total_time_tracked_seconds = total_time_tracked_seconds + sqlc.arg(time_tracked),
    today_time_tracked_seconds = today_time_tracked_seconds + sqlc.arg(time_tracked)
WHERE id = $1 AND owner_id = $2;
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
