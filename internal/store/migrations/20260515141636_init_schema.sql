-- +goose Up
PRAGMA foreign_keys = ON;

CREATE TABLE hooks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    token TEXT NOT NULL UNIQUE,
    name TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE hook_responses (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    hook_id INTEGER NOT NULL UNIQUE,
    status_code INTEGER NOT NULL DEFAULT 200 CHECK (status_code BETWEEN 100 AND 599),
    headers TEXT NOT NULL DEFAULT '{}',
    body BLOB,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (hook_id) REFERENCES hooks(id) ON DELETE CASCADE
);

CREATE TABLE webhook_requests (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    hook_id INTEGER NOT NULL,
    method TEXT NOT NULL,
    path TEXT NOT NULL,
    query TEXT NOT NULL DEFAULT '',
    headers TEXT NOT NULL DEFAULT '{}',
    body BLOB,
    remote_addr TEXT,
    content_type TEXT,
    body_size INTEGER NOT NULL DEFAULT 0 CHECK (body_size >= 0),
    received_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (hook_id) REFERENCES hooks(id) ON DELETE CASCADE
);

CREATE INDEX idx_hooks_token ON hooks(token);
CREATE INDEX idx_webhook_requests_hook_id_received_at ON webhook_requests(hook_id, received_at DESC);

-- +goose Down
DROP INDEX IF EXISTS idx_webhook_requests_hook_id_received_at;
DROP INDEX IF EXISTS idx_hooks_token;
DROP TABLE IF EXISTS webhook_requests;
DROP TABLE IF EXISTS hook_responses;
DROP TABLE IF EXISTS hooks;
