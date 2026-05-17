-- name: CreateHook :one
INSERT INTO hooks (token, name)
VALUES (?, ?)
RETURNING id, token, name, created_at, updated_at;

-- name: GetHookByToken :one
SELECT id, token, name, created_at, updated_at
FROM hooks
WHERE token = ?;

-- name: CreateWebhookRequest :one
INSERT INTO webhook_requests (
    hook_id,
    method,
    path,
    query,
    headers,
    body,
    remote_addr,
    content_type,
    body_size
)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING id, hook_id, method, path, query, headers, body, remote_addr, content_type, body_size, received_at;

-- name: ListWebhookRequestsByHookID :many
SELECT id, hook_id, method, path, query, headers, body, remote_addr, content_type, body_size, received_at
FROM webhook_requests
WHERE hook_id = ?
ORDER BY received_at DESC, id DESC;

-- name: ListWebhookRequestsByTime :many
SELECT id, token, name, created_at, updated_at
FROM hooks
WHERE created_at <= datetime('now', ?);

-- name: DeleteWebhookRequestsOlderThan :execresult
DELETE FROM hooks
WHERE created_at < datetime('now', ?);

-- name: UpsertHookResponse :one
INSERT INTO hook_responses (hook_id, status_code, headers, body)
VALUES (?, ?, ?, ?)
ON CONFLICT (hook_id) DO UPDATE SET
    status_code = excluded.status_code,
    headers = excluded.headers,
    body = excluded.body,
    updated_at = CURRENT_TIMESTAMP
RETURNING id, hook_id, status_code, headers, body, created_at, updated_at;

-- name: GetHookResponseByHookID :one
SELECT id, hook_id, status_code, headers, body, created_at, updated_at
FROM hook_responses
WHERE hook_id = ?;
