CREATE TABLE sessions (
    id SERIAL PRIMARY KEY,
    name VARCHAR(64) NOT NULL,
    password TEXT,
    owner_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    participant_ids INT[] NOT NULL DEFAULT '{}',
    expiry TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
