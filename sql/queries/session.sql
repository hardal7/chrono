-- name: CreateSession :exec
INSERT INTO sessions(name, owner_id, max_participants, password, expires_at, topic)
VALUES($1, $2, $3, $4, $5, $6);
-- name: UpdateSession :exec
UPDATE sessions
SET name = COALESCE(sqlc.arg(new_name), name), max_participants = $3, password = $4, expires_at = $5, updated_at = now()
WHERE owner_id = $1 AND name = $2;
-- name: DeleteSession :exec
DELETE FROM sessions
WHERE owner_id = $1 AND name = $2;
-- name: TrackSessionTime :exec
UPDATE sessions
SET total_time_tracked_seconds = total_time_tracked_seconds + sqlc.arg(time_tracked)
WHERE id = $1;
-- name: GetSessionByNameAndOwnerName :one
SELECT sessions.* FROM sessions
JOIN users ON users.id = sessions.owner_id
WHERE 
    sessions.name = $1
    AND users.username = $2
    AND users.hide_user = FALSE;
-- name: GetJoinedSessions :many
SELECT sessions.* FROM session_participants
JOIN sessions ON sessions.id = session_participants.session_id
WHERE user_id = $1;
-- name: GetSessionByNameAndOwnerID :one
SELECT * FROM sessions
WHERE 
    name = $1 
    AND owner_id = $2
    AND users.hide_user = FALSE;
-- name: GetSessionsAllByFriends :many
WITH friend_users AS (
    SELECT
        CASE
            WHEN sender_id = $1 THEN recipient_id
            ELSE sender_id
        END
    AS friend_id
    FROM friends
    WHERE
        (sender_id = $1 OR recipient_id = $1)
        AND is_accepted = TRUE
)
SELECT 
    sessions.*,
    users.username AS owner_username,
    users.id AS owner_id
FROM sessions
JOIN friend_users ON sessions.owner_id = friend_users.friend_id
JOIN users ON sessions.owner_id = users.id
WHERE 
    sessions.is_active = TRUE
    AND users.hide_user = FALSE;
