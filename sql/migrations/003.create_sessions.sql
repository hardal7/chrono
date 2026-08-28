CREATE TABLE sessions (
    id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    owner_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(64) NOT NULL,
    max_participants INT,
    password TEXT,
    expires_at TIMESTAMPTZ,
    topic VARCHAR(64),
    total_time_tracked_seconds INT NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE (name, owner_id)
);

CREATE TABLE session_participants (
    id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    session_id uuid NOT NULL REFERENCES sessions(id) ON DELETE CASCADE,
    total_time_tracked_seconds INT NOT NULL DEFAULT 0,
    today_time_tracked_seconds INT NOT NULL DEFAULT 0,

    UNIQUE (session_id, user_id)
);

