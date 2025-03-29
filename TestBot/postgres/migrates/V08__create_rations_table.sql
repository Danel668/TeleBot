CREATE TABLE IF NOT EXISTS rations(
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    ration TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX user_id_created_at_rations ON rations (user_id, created_at);
