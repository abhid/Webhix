# Webhix

Self-hosted webhook inspector. Single binary, SQLite, no external dependencies.

webhook.site is the go-to tool for debugging webhooks, but it sends your data to someone else's server. Stripe payloads, OAuth tokens, PII — all of it leaves your network. A lot of companies block it outright for that reason. Webhix runs on your own infrastructure, stores everything locally, and stays out of your way.

## What it does

You create an endpoint, point your webhook source at it, and watch requests come in. Every request is captured in full — headers, body, query params, IP, timestamp, content type, size. The UI updates live without a page refresh.

Beyond just inspecting requests, you can:

- **Replay** any request with one click, optionally with edits
- **Custom responses** — configure the status code, headers, and body your endpoint returns (useful as a lightweight mock server)
- **CLI forwarding** — pipe incoming requests to a local port: `webhix forward <token> --to localhost:3000`
- **Export as curl or HTTPie** -- copy any request as a runnable command

## Quick start

**Binary**

```sh
curl -fsSL https://webhix.dev/install.sh | sh
webhix serve --base-url https://hooks.yourdomain.com
```

**Docker**

```sh
docker run -p 8080:8080 -v webhix-data:/data \
  -e WEBHIX_BASE_URL=https://hooks.yourdomain.com \
  ghcr.io/gaisbax/webhix
```

**Docker Compose**

```yaml
services:
  webhix:
    image: ghcr.io/gaisbax/webhix
    ports: ["8080:8080"]
    volumes: ["./data:/data"]
    environment:
      WEBHIX_BASE_URL: https://hooks.yourdomain.com
```

**Local dev (no domain needed)**

```sh
webhix serve
# Listening on http://localhost:8080
```

Endpoint URLs follow the pattern `https://<base-url>/r/<token>`.

## Auth

Single-user by default. Set a password via env or let Webhix generate one on first run:

```sh
WEBHIX_PASSWORD=yourpassword webhix serve
```

## Reverse proxy

Works behind Caddy, Nginx, Traefik. Reads `X-Forwarded-*` headers automatically. Set `--base-url` or `WEBHIX_BASE_URL` to match your public domain.

## Configuration

| Flag | Env var | Default | Description |
|------|---------|---------|-------------|
| `--base-url` | `WEBHIX_BASE_URL` | `http://localhost:8080` | Public base URL for generated endpoint links |
| `--port` | `WEBHIX_PORT` | `8080` | Port to listen on |
| `--data` | `WEBHIX_DATA` | `./data` | Directory for SQLite database |
| `--password` | `WEBHIX_PASSWORD` | auto-generated | UI password |

## Technical notes

- Written in Go, compiles to a single binary
- SQLite by default, no external database required
- UI is embedded in the binary via `go:embed`
- Runs on Linux, macOS, Windows (amd64 + arm64)
- Memory usage under 50 MB at idle

## Roadmap

**v0.2**
- Multi-user support with basic RBAC
- Webhook signature verification (Stripe, GitHub style)
- Schema validation
- Notifications on new requests (Slack, Telegram, Discord)
- Optional Postgres support
- Auto-HTTPS via Let's Encrypt (no reverse proxy needed)

**v0.3+**
- Tunnel mode — connect Webhix to a managed relay and get a public URL without a server

## License

[AGPL-3.0](LICENSE). Self-hosted use is always free and open source.

If you want to run Webhix as a network service and keep your changes private, contact us for a commercial license.
