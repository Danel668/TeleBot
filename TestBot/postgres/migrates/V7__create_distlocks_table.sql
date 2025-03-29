CREATE TABLE IF NOT EXISTS distlocks(
    lock_name TEXT PRIMARY KEY,
    owner_id TEXT NOT NULL,
    locked_at TIMESTAMPTZ DEFAULT NOW()
);
