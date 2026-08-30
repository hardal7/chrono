CREATE TABLE feature_requests (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255),
    email VARCHAR(255),
    title VARCHAR(255) NOT NULL,
    problem TEXT NOT NULL,
    feature TEXT NOT NULL,
    priority VARCHAR(20) NOT NULL DEFAULT 'medium',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT feature_requests_priority_check
        CHECK (priority IN ('low', 'medium', 'high'))
);

CREATE TABLE bug_reports (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255),
    email VARCHAR(255),
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    steps TEXT,
    environment VARCHAR(255),
    additional TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
