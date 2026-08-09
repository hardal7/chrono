CREATE TABLE friends (
    id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    owner_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    recipient_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    is_accepted BOOLEAN NOT NULL DEFAULT FALSE
);

