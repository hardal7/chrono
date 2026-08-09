CREATE TABLE sessions (
    id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    name VARCHAR(64) NOT NULL,
    password TEXT,
    owner_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    participant_ids uuid[] NOT NULL DEFAULT '{}',
    expiry TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
