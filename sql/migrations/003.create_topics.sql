CREATE TABLE topics (
    id SERIAL PRIMARY KEY,
    name VARCHAR(64) NOT NULL,
    time_tracked_seconds INT NOT NULL,
    created_by_userid INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
