CREATE TABLE IF NOT EXISTS recommendations (
    user_id BIGINT NOT NULL,
    recommendation TEXT NOT NULL,
    send_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX idx_send_at_recommendations ON recommendations (send_at);
