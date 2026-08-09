-- name: CreateFriend :exec
INSERT INTO friends(owner_id, recipient_id)
VALUES($1, $2);
-- name: DeleteFriend :exec
DELETE FROM friends
WHERE owner_id = $1 AND recipient_id = $2;
-- name: AcceptFriendRequest :exec
UPDATE friends
SET is_accepted = true
WHERE owner_id = $1 AND recipient_id = $2;
-- name: GetFriendRequests :many
SELECT * FROM friends
WHERE recipient_id = $1;
-- name: GetOwnedFriendRequests :many
SELECT * FROM friends
WHERE owner_id = $1;
