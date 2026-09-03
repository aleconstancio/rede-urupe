CREATE TABLE IF NOT EXISTS active_channels (
    guild_id TEXT NOT NULL,
    channel_id TEXT NOT NULL,
    is_active INTEGER DEFAULT 1,
    registered_at TEXT DEFAULT (datetime('now')),
    PRIMARY KEY (guild_id, channel_id)
);
