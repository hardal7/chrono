-- name: CreateLeaderboardSnapshot :one
INSERT INTO leaderboard_snapshots
DEFAULT VALUES
RETURNING id;
-- name: GetLastLeaderboardUsers :many
SELECT * FROM leaderboard_users
WHERE snapshot_id = (
    SELECT id
    FROM leaderboard_snapshots
    ORDER BY created_at DESC
    LIMIT 1
);
-- name: CreateLeaderboardUser :exec
INSERT INTO leaderboard_users(snapshot_id, user_id, rank, rank_change)
VALUES($1, $2, $3, $4);
