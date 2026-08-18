-- name: CreateFriendRequest :exec
INSERT INTO friends(sender_id, recipient_id)
SELECT $1, id
FROM users
WHERE username = $2;
-- name: DeleteFriend :exec
DELETE FROM friends
USING users
WHERE 
    users.username = $2 AND
    ((sender_id = $1 AND recipient_id = $2)
    OR (sender_id = $2 AND recipient_id = $1));
-- name: AcceptFriendRequest :exec
UPDATE friends
SET is_accepted = true
FROM users WHERE 
    users.id = friends.sender_id
    AND users.username = $1 AND recipient_id = $2;
-- name: GetFriendRequests :many
SELECT * FROM friends
WHERE recipient_id = $1;
-- name: GetSentFriendRequests :many
SELECT * FROM friends
WHERE sender_id = $1;
-- name: GetTopFriends :many
SELECT users.* FROM friends
JOIN users ON 
    users.id = friends.recipient_id
    OR users.id = friends.sender_id
WHERE 
    users.id = $1
    AND friends.is_accepted = TRUE
    AND users.total_time_tracked_seconds < $2
ORDER BY users.total_time_tracked_seconds DESC
LIMIT $3;
