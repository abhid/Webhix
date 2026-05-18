# Webhix

[![Release](https://img.shields.io/github/v/release/gaisbax/webhix)](https://github.com/gaisbax/webhix/releases)
[![License](https://img.shields.io/badge/license-AGPL--3.0-blue)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/gaisbax/webhix)](https://goreportcard.com/report/github.com/gaisbax/webhix)
[![Docker](https://img.shields.io/badge/docker-ghcr.io-blue)](https://github.com/gaisbax/webhix/pkgs/container/webhix)
[![RU](https://img.shields.io/badge/lang-ru-blue)](docs/README.ru.md)
[![Contributing](https://img.shields.io/badge/contributing-guide-brightgreen)](CONTRIBUTING.md)

Self-hosted webhook inspector. Single binary, SQLite, no external dependencies.

webhook.site is the go-to tool for debugging webhooks, but it sends your data to someone else's server. Stripe payloads, OAuth tokens, PII — all of it leaves your network. A lot of companies block it outright for that reason. Webhix runs on your own infrastructure and stores everything locally.

![Webhix UI](docs/screenshot.png)

## Features

- 📡 Capture any HTTP method — headers, body, query params, IP, timestamp, content type, size
- 🔴 Live UI updates via SSE — no page refresh needed
- 🪞 Replay any request with one click
- 🎭 Custom responses — configure status, headers, and body (lightweight mock server)
- 🔁 CLI forwarding to localhost: `webhix forward <token> --to localhost:3000`
- 📋 Export as curl — copy any request as a runnable command
- 🔍 Full-text search and filter by HTTP method
- 🔒 Basic auth out of the box
- 🐳 Docker, Compose, or standalone binary
- 💾 SQLite by default — no Redis or Postgres required

## Why not webhook.site / smee.io / webhook-tester?

|                 | Webhix        | webhook.site (self-hosted)    | smee.io        | tarampampam/webhook-tester |
| --------------- | ------------- | ----------------------------- | -------------- | -------------------------- |
| Self-hosted     | ✅            | ✅                            | ❌             | ✅                         |
| Single binary   | ✅            | ❌ PHP + Composer + MySQL     | ❌             | ❌ Redis or fs driver      |
| Request history | ✅            | ✅                            | ❌             | ✅                         |
| Live UI         | ✅            | ✅                            | ❌             | ✅                         |
| Replay          | ✅            | ❌                            | ❌             | ❌                         |
| CLI forwarding  | ✅ built-in   | ❌                            | ✅ only this   | ❌ needs ngrok             |
| Custom responses| ✅            | ❌                            | ❌             | ❌                         |

Webhix is the only tool combining single-binary deployment, request replay, and custom responses — no Redis, no PHP, no external tunnel services.

## Quick start

### Binary

```sh
curl -fsSL https://webhix.online/install.sh | sh
webhix serve --base-url https://hooks.yourdomain.com
```

Or download manually from [releases](https://github.com/gaisbax/webhix/releases/latest).

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
