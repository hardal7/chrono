-- name: CreateSessionToken :exec
INSERT INTO session_tokens(user_id, hash, expiry)
VALUES($1, $2, $3);
-- name: DeleteSessionToken :exec
DELETE FROM session_tokens
WHERE id = $1;
-- name: GetSessionToken :one
SELECT * FROM session_tokens
WHERE hash = $1 AND expiry > now();
