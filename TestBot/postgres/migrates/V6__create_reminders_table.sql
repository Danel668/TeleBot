CREATE TABLE IF NOT EXISTS reminders (
    user_id BIGINT,
    send_at TIMESTAMPTZ NOT NULL,
    expire_at TIMESTAMPTZ NOT NULL,
    reminder TEXT NOT NULL,

    PRIMARY KEY(user_id, send_at)
);
