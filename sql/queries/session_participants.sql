-- name: GetSessionParticipants :many
SELECT users.username, users.last_seen_at, session_participants.*
FROM session_participants
JOIN users ON session_participants.user_id = users.id
WHERE 
    session_id = $1
    AND users.hide_user = FALSE
ORDER BY session_participants.total_time_tracked_seconds DESC;
-- name: JoinSession :one
INSERT INTO session_participants (user_id, session_id)
SELECT
    sqlc.arg(user_id),
    sessions.id
FROM sessions
JOIN users ON users.id = sessions.owner_id
WHERE sessions.name = sqlc.arg(session_name)
    AND users.username_normalized = LOWER(sqlc.arg(owner_username))
    AND users.hide_user = FALSE
    AND (
        sessions.expires_at IS NULL
        OR sessions.expires_at > NOW()
    )
    AND (
        sessions.max_participants IS NULL
        OR (
            SELECT COUNT(*)
            FROM session_participants
            WHERE session_id = sessions.id
        ) < sessions.max_participants
    )
RETURNING (
    SELECT topic
    FROM sessions
    WHERE sessions.id = session_participants.session_id
) AS topic;
-- name: LeaveSession :exec
DELETE FROM session_participants
WHERE (
    user_id = $1
    AND session_id = (
        SELECT sessions.id
        FROM sessions
        JOIN users ON users.id = sessions.owner_id
        WHERE 
          sessions.name = $2
          AND users.username_normalized = LOWER(sqlc.arg(owner_username))
    )
);
-- name: KickFromSession :exec
DELETE FROM session_participants
USING users, sessions
WHERE
    session_participants.user_id = users.id
    AND session_participants.session_id = sessions.id
    AND sessions.owner_id = $1 
    AND sessions.name = $2
    AND users.username_normalized = LOWER(sqlc.arg(participant_username));
-- name: TrackSessionParticipantTime :exec
UPDATE session_participants
SET
    total_time_tracked_seconds = total_time_tracked_seconds + sqlc.arg(time_tracked),
    today_time_tracked_seconds = today_time_tracked_seconds + sqlc.arg(time_tracked)
WHERE user_id = $1 AND session_id = $2;
-- name: ResetSessionParticipantTimeTrackedToday :exec
UPDATE session_participants
SET today_time_tracked_seconds = 0;
-- name: DeleteExpiredSessionParticipants :exec
DELETE FROM session_participants
USING sessions
WHERE 
  session_participants.session_id = sessions.id
  AND (sessions.expires_at IS NOT NULL AND sessions.expires_at < NOW());
