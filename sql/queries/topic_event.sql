-- name: CreateTopicEvent :exec
INSERT INTO topic_events(user_id, topic_id, time_tracked_seconds, created_at)
VALUES($1, $2, $3, $4);
-- name: Upcreated_atTopicEvent :exec
UPDATE topic_events
SET time_tracked_seconds = $2, created_at = $3
WHERE id = $1;
-- name: DeleteTopicEvent :exec
DELETE FROM topic_events
WHERE id = $1;
-- name: GetAllTopicEvents :many
SELECT * FROM topic_events
WHERE user_id = $1;
-- name: GetTopicEventsToday :many
SELECT * FROM topic_events
WHERE user_id = $1 AND created_at(date) = CURRENT_DATE;
