-- +goose Up
CREATE TABLE hook_notifications (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    hook_id INTEGER NOT NULL UNIQUE,
    telegram_bot_token TEXT NOT NULL DEFAULT '',
    telegram_chat_id TEXT NOT NULL DEFAULT '',
    proxy_url TEXT NOT NULL DEFAULT '',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (hook_id) REFERENCES hooks(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS hook_notifications;
