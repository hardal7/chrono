CREATE TABLE users (
    id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    email VARCHAR(64) NOT NULL UNIQUE,
    username VARCHAR(64) NOT NULL UNIQUE,
    password TEXT NOT NULL,
    total_time_tracked_seconds INT NOT NULL DEFAULT 0,
    today_time_tracked_seconds INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE friends (
    id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    sender_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    recipient_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    is_accepted BOOLEAN NOT NULL DEFAULT FALSE,

    UNIQUE (sender_id, recipient_id)
);

CREATE TABLE sessions (
    id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    owner_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(64) NOT NULL,
    max_participants INT,
    password TEXT,
    expires_at TIMESTAMPTZ,
    topic VARCHAR(64),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE session_participants (
    id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    session_id uuid NOT NULL REFERENCES sessions(id) ON DELETE CASCADE,
    last_seen_at TIMESTAMPTZ NOT NULL,
    total_time_tracked_seconds INT NOT NULL DEFAULT 0,
    today_time_tracked_seconds INT NOT NULL DEFAULT 0,

    UNIQUE (session_id, user_id)
);

CREATE TABLE topics (
    id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    name VARCHAR(64) NOT NULL,
    streak INT NOT NULL DEFAULT 0,
    today_time_tracked_seconds INT NOT NULL DEFAULT 0,
    total_time_tracked_seconds INT NOT NULL DEFAULT 0,
    owner_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE topic_events (
    id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    topic_id uuid NOT NULL REFERENCES topics(id) ON DELETE CASCADE,
    time_tracked_seconds INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
