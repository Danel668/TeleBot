CREATE TABLE IF NOT EXISTS users (
    user_id BIGINT PRIMARY KEY,
    timezone TEXT NOT NULL,
    registration_at TIMESTAMPTZ NOT NULL
);
