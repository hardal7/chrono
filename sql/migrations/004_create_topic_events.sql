CREATE TABLE topic_events (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    topic_id INT NOT NULL REFERENCES topics(id) ON DELETE CASCADE,
    time_tracked_seconds INT NOT NULL,
    date TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
