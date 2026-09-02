CREATE TABLE leaderboard_snapshots (
    id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
INSERT INTO leaderboard_snapshots DEFAULT VALUES;

CREATE TABLE leaderboard_users (
    id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    snapshot_id uuid NOT NULL REFERENCES leaderboard_snapshots(id) ON DELETE CASCADE,
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    rank INT NOT NULL,
    rank_change INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE (snapshot_id, rank),
    UNIQUE (snapshot_id, user_id)
);
