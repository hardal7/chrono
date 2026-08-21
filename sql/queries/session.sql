-- name: CreateSession :exec
INSERT INTO sessions(name, owner_id, max_participants, password, expires_at, topic)
VALUES($1, $2, $3, $4, $5, $6);
-- name: UpdateSession :exec
UPDATE sessions
SET name = $2, max_participants = $3, password = $4, expires_at = $5, updated_at = now()
WHERE id = $1;
-- name: DeleteSession :exec
DELETE FROM sessions
WHERE id = $1;
-- name: GetSessionByNameAndOwnerName :one
SELECT sessions.* FROM sessions
JOIN users ON users.id = sessions.owner_id
WHERE sessions.name = $1 AND users.username = $2;
-- name: GetSessionByNameAndOwnerID :one
SELECT * FROM sessions
WHERE name = $1 AND owner_id = $2;
-- name: JoinSession :exec
INSERT INTO session_participants(user_id, session_id, last_seen_at)
VALUES($1, $2, $3);
-- name: GetSessionParticipants :many
SELECT * FROM session_participants
WHERE session_id = $1;
