# Webhix Tunnel Protocol v1

## Overview

The Webhix Tunnel Protocol (WTP) enables a local `webhix tunnel` client to receive
HTTP requests forwarded from the public relay server (`relay.webhix.online`) and proxy
them to a local port, returning the real response back through the relay to the
original HTTP caller.

```
[Webhook sender]
      │
      ▼
abc123.webhix.online  (relay server, managed by Webhix)
      │  WebSocket (wss://)
      ▼
webhix tunnel 3000   (client on user's machine)
      │
      ▼
localhost:3000        (user's local app)
```

## Transport

- **WebSocket over TLS** (`wss://`)
- Default relay endpoint: `wss://relay.webhix.online/tunnel`
- Text frames, JSON-encoded messages
- One persistent WebSocket connection per tunnel
- All messages are UTF-8 JSON

## Message Format

All messages are JSON objects with a required `type` field:

```json
{
  "type": "<message_type>",
  "...": "...fields"
}
```

---

## Handshake

### Step 1 — Client → Relay: `register`

Sent immediately after WebSocket connection is established.

```json
{
  "type": "register",
  "version": "1",
  "auth_token": "tok_xxxx",
  "subdomain": "myapp"
}
```

| Field        | Required | Description                                 |
| ------------ | -------- | ------------------------------------------- |
| `version`    | yes      | Protocol version, must be `"1"`             |
| `auth_token` | no       | Pro auth token from webhix.online dashboard |
| `subdomain`  | no       | Requested reserved subdomain (Pro only)     |

### Step 2a — Relay → Client: `registered` (success)

```json
{
  "type": "registered",
  "subdomain": "abc123",
  "url": "https://abc123.webhix.online",
  "tier": "free",
  "expires_at": "2024-01-01T20:00:00Z",
  "expires_in_seconds": 7200
}
```

| Field                | Description                                                 |
| -------------------- | ----------------------------------------------------------- |
| `subdomain`          | Assigned subdomain (may differ from requested)              |
| `url`                | Full public URL to share                                    |
| `tier`               | `"free"` or `"pro"`                                         |
| `expires_at`         | ISO 8601 UTC, free tier session expiry; `null` for Pro      |
| `expires_in_seconds` | Seconds until expiry, for countdown display; `null` for Pro |

### Step 2b — Relay → Client: `error` (failure)

```json
{
  "type": "error",
  "code": "AUTH_FAILED",
  "message": "Invalid or expired auth token"
}
```

Connection is closed by the relay after sending an error.

---

## Request Forwarding

### Relay → Client: `request`

Sent when an HTTP request arrives at `<subdomain>.webhix.online`.

```json
{
  "type": "request",
  "id": "01HV2K9XMABCDEF123456789",
  "method": "POST",
  "path": "/webhook/stripe",
  "query": "sig=abc123",
  "headers": {
    "Content-Type": ["application/json"],
    "Stripe-Signature": ["t=1700000000,v1=abc..."]
  },
  "body": "eyJldmVudCI6InBheW1lbnQuc3VjY2VlZGVkIn0="
}
```

| Field     | Description                                                  |
| --------- | ------------------------------------------------------------ |
| `id`      | ULID or UUID v4, unique per request. Used to match responses |
| `method`  | HTTP method                                                  |
| `path`    | Request path, always starts with `/`                         |
| `query`   | Raw query string without `?`, may be empty string            |
| `headers` | Map of header name → array of values                         |
| `body`    | Base64-encoded request body, may be empty string             |

### Client → Relay: `response`

Sent after the client proxies the request to the local port and receives a response.

```json
{
  "type": "response",
  "id": "01HV2K9XMABCDEF123456789",
  "status": 200,
  "headers": {
    "Content-Type": ["application/json"]
  },
  "body": "eyJvayI6dHJ1ZX0="
}
```

| Field     | Description                                       |
| --------- | ------------------------------------------------- |
| `id`      | Same ID as the corresponding `request` message    |
| `status`  | HTTP status code                                  |
| `headers` | Map of header name → array of values              |
| `body`    | Base64-encoded response body, may be empty string |

**If the local port is unreachable** (connection refused, timeout), the client MUST
send a 502 response rather than dropping the message:

```json
{
  "type": "response",
  "id": "01HV2K9XMABCDEF123456789",
  "status": 502,
  "headers": {},
  "body": ""
}
```

### Concurrent requests

Multiple `request` messages may be in-flight simultaneously. Responses can arrive
in any order. The relay matches requests to responses by `id`. Clients SHOULD
process each request in its own goroutine.

---

## Keepalive

- The relay sends WebSocket **Ping** frames every **30 seconds**
- The client MUST respond with **Pong** frames (standard WebSocket behavior)
- If no Pong is received within **10 seconds**, the relay closes the connection
- The client SHOULD reconnect with exponential backoff: 1s → 2s → 4s → 8s → max 30s

---

## Auth & Tiers

### Free tier (no `auth_token`)

- Random subdomain, changes on every reconnect
- 1 concurrent tunnel per IP
- Rate limit: **60 requests/minute** (matching ngrok free tier)
- Bandwidth: **1 GB/day**
- TTL: **2 hours** per session (main motivator to upgrade)
- No `auth_token` field in `register`

### Pro tier ($5–8/month)

- Provide `auth_token` obtained from [webhix.online](https://webhix.online) dashboard
- Optional `subdomain` for a **reserved subdomain** (main paid feature)
- Up to **3 concurrent tunnels**
- No rate limit, no bandwidth cap
- No TTL — `expires_at` and `expires_in_seconds` are `null` in `registered` response

---

## Limits

These limits apply at the relay level regardless of tier.

| Limit                                        | Value     |
| -------------------------------------------- | --------- |
| Max request body                             | 10 MB     |
| Max response body                            | 10 MB     |
| Max headers size                             | 64 KB     |
| Max WebSocket frame                          | 16 MB     |
| Idle timeout (no traffic)                    | 5 minutes |
| Max concurrent in-flight requests per tunnel | 100       |

Requests exceeding body or headers limits receive a `413 Request Entity Too Large`
response from the relay without being forwarded to the client.

---

## Error Codes

| Code                  | Description                                         |
| --------------------- | --------------------------------------------------- |
| `AUTH_FAILED`         | Invalid or expired auth token                       |
| `SUBDOMAIN_TAKEN`     | Requested subdomain is currently in use             |
| `SUBDOMAIN_RESERVED`  | Requested subdomain requires Pro tier               |
| `TUNNEL_LIMIT`        | Maximum concurrent tunnels reached for this account |
| `RATE_LIMITED`        | Rate limit exceeded (free tier)                     |
| `RELAY_ERROR`         | Internal relay server error                         |
| `UNSUPPORTED_VERSION` | Protocol version not supported by this relay        |

---

## Versioning

The protocol version is declared in the `register` message (`"version": "1"`).

The relay MUST:

- Accept `"1"` (this document)
- Reject unknown versions with `{"type": "error", "code": "UNSUPPORTED_VERSION"}`

Future versions will be documented in this file with a changelog section.

---

## Security

- All communication is over TLS (`wss://`)
- Auth tokens are opaque Bearer tokens; clients MUST NOT log them
- The relay validates tokens against the billing service on every `register`
- Tunnel traffic is isolated: one client cannot read another client's requests
- The relay does not store request/response bodies; they are forwarded in memory only

---

## Self-Hosted Relay

The tunnel client (`webhix tunnel`) supports a `--relay` flag to point to a
custom relay server. This enables enterprise deployments with a private relay.

```bash
webhix tunnel 3000 --relay wss://tunnel.mycompany.com/tunnel --auth-token tok_xxx
```

The relay server implementation is not part of this open-source repository.
The protocol spec is open so that third-party relay implementations are possible.
