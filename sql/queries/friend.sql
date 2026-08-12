-- name: CreateFriend :exec
INSERT INTO friends(sender_id, recipient_id)
VALUES($1, $2);
-- name: DeleteFriend :exec
DELETE FROM friends
WHERE sender_id = $1 AND recipient_id = $2;
-- name: AcceptFriendRequest :exec
UPDATE friends
SET is_accepted = true
WHERE sender_id = $1 AND recipient_id = $2;
-- name: GetFriendRequests :many
SELECT * FROM friends
WHERE recipient_id = $1;
-- name: GetSentriendRequests :many
SELECT * FROM friends
WHERE sender_id = $1;
-- name: GetTopFriends :many
SELECT users.* FROM friends
JOIN users ON users.id = friends.recipient_id
WHERE 
    friends.is_accepted = TRUE
    AND users.id = $1
    AND users.total_time_tracked_seconds < $2
ORDER BY users.total_time_tracked_seconds DESC
LIMIT $3;
