CREATE TABLE IF NOT EXISTS schema_migrations (
    version TEXT PRIMARY KEY,
    applied_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS users (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    discord_id    TEXT    UNIQUE NOT NULL,
    email         TEXT    UNIQUE NOT NULL,
    display_name  TEXT    NOT NULL,
    password_hash TEXT    NOT NULL,
    status        TEXT    NOT NULL DEFAULT 'waiting'
                  CHECK(status IN ('waiting','approved','rejected','suspended')),
    role          TEXT    NOT NULL DEFAULT 'user'
                  CHECK(role IN ('user','admin')),
    created_at    TEXT    NOT NULL DEFAULT (datetime('now')),
    updated_at    TEXT    NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_users_discord_id ON users(discord_id);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_status ON users(status);

CREATE TABLE IF NOT EXISTS bot_credentials (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id         INTEGER UNIQUE NOT NULL
                    REFERENCES users(id) ON DELETE CASCADE,
    application_id  TEXT    NOT NULL,
    bot_token       TEXT    NOT NULL,
    guild_id        TEXT    NOT NULL,
    public_key      TEXT    NOT NULL,
    approved_by     INTEGER REFERENCES users(id),
    approved_at     TEXT,
    created_at      TEXT    NOT NULL DEFAULT (datetime('now')),
    updated_at      TEXT    NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_bot_credentials_user_id ON bot_credentials(user_id);
