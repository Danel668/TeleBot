CREATE TYPE banned_section AS ENUM ('all', 'recommendation');

CREATE TABLE IF NOT EXISTS banned_users (
    user_id BIGINT PRIMARY KEY,
    reason TEXT,
    banned_at TIMESTAMPTZ NOT NULL,
    banned_section banned_section NOT NULL
);
