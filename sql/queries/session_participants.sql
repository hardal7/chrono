-- name: GetSessionParticipantsAsUsers :many
SELECT users.* FROM session_participants
JOIN users ON session_participants.user_id = users.id
WHERE 
    session_id = $1
    AND users.hide_user = FALSE;
-- name: JoinSession :exec
INSERT INTO session_participants(user_id, session_id)
VALUES($1, $2);
-- name: KickFromSession :exec
DELETE FROM session_participants
USING users, sessions
WHERE
    session_participants.user_id = users.id
    AND session_participants.session_id = sessions.id
    AND sessions.owner_id = $1 
    AND sessions.name = $2
    AND users.username = sqlc.arg(participant_username);
-- name: TrackSessionParticipantTime :exec
UPDATE session_participants
SET
    total_time_tracked_seconds = total_time_tracked_seconds + sqlc.arg(time_tracked),
    today_time_tracked_seconds = today_time_tracked_seconds + sqlc.arg(time_tracked)
WHERE user_id = $1 AND session_id = $2;
-- name: ResetSessionParticipantTimeTrackedToday :exec
UPDATE session_participants
SET today_time_tracked_seconds = 0;
