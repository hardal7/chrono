-- name: CreateTopic :exec
INSERT INTO topics(name, time_tracked_seconds, created_by_userid)
VALUES($1, $2, $3);
-- name: UpdateTopic :exec
UPDATE topics
SET name = $2, time_tracked_seconds = $3, updated_at = now()
WHERE id = $1;
-- name: DeleteTopic :exec
DELETE FROM topics
WHERE id = $1;
-- name: GetTopicByID :one
SELECT * FROM topics
WHERE id = $1;
-- name: GetTopicOfUserByName :one
SELECT * FROM topics
WHERE name = $1 AND created_by_userid = $2;
-- name: TrackTopicTime :exec
UPDATE topics
SET time_tracked_seconds = time_tracked_seconds + $3
WHERE id = $1 AND created_by_userid = $2;
