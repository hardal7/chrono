-- name: CreateFriendRequest :exec
INSERT INTO friends(sender_id, recipient_id)
SELECT $1, id
FROM users
WHERE username_normalized = LOWER(sqlc.arg(username));
-- name: DeleteFriend :exec
DELETE FROM friends
USING users 
WHERE users.username_normalized = LOWER(sqlc.arg(username))
  AND (
    (friends.sender_id = $1 AND friends.recipient_id = users.id)
    OR
    (friends.sender_id = users.id AND friends.recipient_id = $1)
  );
-- name: AcceptFriendRequest :exec
UPDATE friends
SET 
    is_accepted = true,
    updated_at = NOW()
FROM users WHERE
    users.id = friends.sender_id
    AND users.username_normalized = LOWER(sqlc.arg(username)) AND recipient_id = $1;
-- name: GetFriendRequests :many
SELECT friends.created_at, senders.username FROM friends
JOIN users AS senders ON senders.id = friends.sender_id
JOIN users AS recipients ON recipients.id = friends.recipient_id
WHERE 
    recipients.id = $1
    AND friends.is_accepted = FALSE;
-- name: GetFriendStatus :one
SELECT friends.is_accepted FROM friends
JOIN users AS senders ON senders.id = friends.sender_id
JOIN users AS recipients ON recipients.id = friends.recipient_id
WHERE 
    recipients.username_normalized = LOWER(sqlc.arg(username))
    OR senders.username_normalized = LOWER(sqlc.arg(username));
-- name: GetTopFriends :many
SELECT users.* FROM friends
JOIN users ON 
    users.id = friends.recipient_id
    OR users.id = friends.sender_id
WHERE 
    users.id = $1
    AND friends.is_accepted = TRUE
    AND users.week_time_tracked_seconds < sqlc.arg(cursor)
    AND users.username_normalized ILIKE sqlc.arg(match_name) || '%'
    AND users.hide_user = FALSE
ORDER BY users.week_time_tracked_seconds DESC
LIMIT $2;
