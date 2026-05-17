# Webhix

[![RU](https://img.shields.io/badge/lang-ru-blue)](docs/README.ru.md) [![Contributing](https://img.shields.io/badge/contributing-guide-brightgreen)](CONTRIBUTING.md)

Self-hosted webhook inspector. Single binary, SQLite, no external dependencies.

webhook.site is the go-to tool for debugging webhooks, but it sends your data to someone else's server. Stripe payloads, OAuth tokens, PII — all of it leaves your network. A lot of companies block it outright for that reason. Webhix runs on your own infrastructure, stores everything locally, and stays out of your way.

## What it does

You create an endpoint, point your webhook source at it, and watch requests come in. Every request is captured in full — headers, body, query params, IP, timestamp, content type, size. The UI updates live without a page refresh.

Beyond just inspecting requests, you can:

- **Replay** any captured request with one click
- **Custom responses** — configure the status code, headers, and body your endpoint returns to senders (useful as a lightweight mock server)
- **CLI forwarding** — pipe incoming requests to a local port: `webhix forward <token> --to localhost:3000`
- **Export as curl** — copy any request as a runnable curl command
- **Search and filter** — filter requests by text or HTTP method

## Quick start

### Binary

```sh
curl -fsSL https://webhix.dev/install.sh | sh
webhix serve --base-url https://hooks.yourdomain.com
```

### Docker

```sh
docker run -p 8080:8080 -v webhix-data:/data \
  -e WEBHIX_BASE_URL=https://hooks.yourdomain.com \
  ghcr.io/gaisbax/webhix
```

### Docker Compose

```yaml
services:
  webhix:
    image: ghcr.io/gaisbax/webhix
    ports: ["8080:8080"]
    volumes: ["./data:/data"]
    environment:
      WEBHIX_BASE_URL: https://hooks.yourdomain.com
```

### Local dev (no domain needed)

```sh
webhix serve
# Listening on http://localhost:8080
```

Endpoint URLs follow the pattern `https://<base-url>/r/<token>`.

## Auth

Auth is required. Set at least one of:

```sh
# Basic auth password (browser login)
WEBHIX_PASSWORD=yourpassword webhix serve

# Secret key for API / CLI access (Authorization: Bearer or X-Webhix-Key header)
WEBHIX_SECRET_KEY=yourkey webhix serve

# Both at once
webhix serve --password yourpassword --secret-key yourkey
```

Webhook capture URLs (`/r/<token>`) are always public — no auth required there.

## Reverse proxy

Works behind Caddy, Nginx, Traefik. Reads `X-Forwarded-*` headers automatically. Set `--base-url` or `WEBHIX_BASE_URL` to match your public domain.

## Configuration

| Env var | Default | Description |
| ------- | ------- | ----------- |
| `WEBHIX_BASE_URL` | `http://localhost:8080` | Public base URL for generated endpoint links |
| `WEBHIX_ADDR` | `:8080` | Address to listen on (e.g. `0.0.0.0:9000`) |
| `WEBHIX_DB_PATH` | `./data` | Path to SQLite database directory |

## Technical notes

- Written in Go, compiles to a single binary
- SQLite by default, no external database required
- UI is embedded in the binary via `go:embed`
- Runs on Linux, macOS, Windows (amd64 + arm64)
- Memory usage under 50 MB at idle

## Roadmap

### v0.2

- Multi-user support with basic RBAC
- Webhook signature verification (Stripe, GitHub style)
- Schema validation
- Notifications on new requests (Slack, Telegram, Discord)
- Optional Postgres support
- Auto-HTTPS via Let's Encrypt (no reverse proxy needed)

### v0.3+

- Tunnel mode — connect Webhix to a managed relay and get a public URL without a server

## License

[AGPL-3.0](LICENSE). Self-hosted use is always free and open source.

If you want to run Webhix as a network service and keep your changes private, contact us for a commercial license.
