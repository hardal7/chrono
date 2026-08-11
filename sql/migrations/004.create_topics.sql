CREATE TABLE topics (
    id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    name VARCHAR(64) NOT NULL,
    streak INT NOT NULL DEFAULT 0,
    today_time_tracked_seconds INT NOT NULL DEFAULT 0,
    total_time_tracked_seconds INT NOT NULL DEFAULT 0,
    created_by_userid uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
