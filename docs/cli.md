# Webhix CLI design

This document describes the intended command-line interface for Webhix.

Webhix is a self-hosted webhook inspector. It runs as a single binary, stores data in SQLite, embeds the web UI, and does not require external services by default. The primary product flow is simple: create an endpoint, point a webhook sender at `/r/<token>`, inspect captured requests in the UI, optionally replay them, configure custom responses, or forward them to a local service.

The CLI should preserve that simplicity. Common local and production workflows must stay short, while advanced operational tasks should be explicit, scriptable, and safe.

## Project context

Webhix exists as a private, self-hosted alternative to hosted webhook inspection tools. Hosted tools are convenient, but they receive sensitive payloads such as Stripe events, OAuth tokens, internal IDs, PII, and integration secrets. Webhix keeps those payloads inside the user's infrastructure.

The project currently targets:

- A single binary distribution.
- SQLite storage by default.
- No required external dependencies.
- Embedded web UI through Go `embed`.
- Local development with no public domain.
- Production usage behind Caddy, Nginx, Traefik, or another reverse proxy.
- Live request inspection.
- Custom endpoint responses.
- Request replay.
- CLI forwarding to local ports.
- Exporting captured requests as runnable commands.

The README currently documents these configuration values:

| Environment variable | Default                 | Description                                        |
| -------------------- | ----------------------- | -------------------------------------------------- |
| `WEBHIX_BASE_URL`    | `http://localhost:8080` | Public base URL used for generated endpoint links. |
| `WEBHIX_ADDR`        | `:8080`                 | Address the HTTP server listens on.                |
| `WEBHIX_DB_PATH`     | `./data`                | Path to the SQLite database directory.             |

The CLI should map these environment variables to clear flags and should allow configuration through the following priority order:

1. Explicit CLI flags.
2. Environment variables.
3. Config file values.
4. Built-in defaults.

## CLI principles

### Keep the happy path short

The most important commands should be memorable:

```sh
webhix serve
webhix endpoint create
webhix forward <token> --to localhost:3000
webhix request tail <token>
```

Users should not need a long config file or multiple setup steps to start inspecting webhooks locally.

### Make production configuration explicit

Production users need predictable behavior around bind addresses, public URLs, storage paths, reverse proxies, retention, and authentication. These settings should be available as both environment variables and flags.

### Be scriptable

Commands that return data should support machine-readable output through `--output json` and `--output yaml`. Commands that mutate state should support non-interactive usage through `--yes` where confirmation would otherwise be required.

### Be careful with captured secrets

Captured webhook requests often contain sensitive values. Commands that print, export, replay, or forward requests should provide redaction and filtering options where practical.

### Separate stable core from roadmap features

The core CLI should cover the current product:

- Server startup.
- Endpoint management.
- Request inspection.
- Forwarding.
- Replay.
- Export.
- Database maintenance.
- Diagnostics.

Roadmap commands can be designed now but implemented later:

- Multi-user and RBAC.
- Signature verification.
- Schema validation.
- Notifications.
- Optional Postgres.
- Tunnel mode.

## Global command shape

```sh
webhix [command] [flags]
```

The root command should provide product-level help, version information, and command groups.

### Global flags

These flags should be available to all commands where they make sense.

#### `--config <path>`

Path to a config file.

Recommended default:

```sh
config/.env
```

This flag is useful when running multiple Webhix instances on one machine or when a service manager passes an explicit config path.

Example:

```sh
webhix serve --config /etc/webhix/webhix.env
```

#### `--db-path <path>`

Path to the SQLite database directory.

Environment equivalent:

```sh
WEBHIX_DB_PATH
```

Default:

```sh
./data
```

Example:

```sh
webhix serve --db-path /var/lib/webhix
```

#### `--base-url <url>`

Public base URL used when Webhix generates endpoint links.

Environment equivalent:

```sh
WEBHIX_BASE_URL
```

Default:

```sh
http://localhost:8080
```

Example:

```sh
webhix serve --base-url https://hooks.example.com
```

If the endpoint token is `abc123`, the public capture URL becomes:

```text
https://hooks.example.com/r/abc123
```

#### `--addr <addr>`

HTTP listen address.

Environment equivalent:

```sh
WEBHIX_ADDR
```

Default:

```sh
:8080
```

Examples:

```sh
webhix serve --addr :8080
webhix serve --addr 127.0.0.1:8080
webhix serve --addr 0.0.0.0:9000
```

#### `--log-level <debug|info|warn|error>`

Controls log verbosity.

Recommended default:

```text
info
```

Examples:

```sh
webhix serve --log-level debug
webhix serve --log-level warn
```

#### `--log-format <text|json>`

Controls log output format.

Recommended default:

```text
text
```

For Docker, Kubernetes, and structured log collectors, `json` is usually better:

```sh
webhix serve --log-format json
```

#### `--no-color`

Disables colored terminal output.

Useful in CI or log files:

```sh
webhix doctor --no-color
```

#### `--yes`, `-y`

Skips confirmation prompts for commands that mutate or delete data.

Example:

```sh
webhix request clear abc123 --yes
```

This flag should only be available on commands that normally ask for confirmation.

#### `--output <text|table|json|yaml>`

Controls command output.

Recommended behavior:

- `text` for simple human-readable output.
- `table` for list commands.
- `json` for scripts and integrations.
- `yaml` for readable structured output.

Examples:

```sh
webhix endpoint list --output table
webhix endpoint list --output json
webhix request get req_123 --output yaml
```

#### `--help`, `-h`

Shows help for the current command.

Examples:

```sh
webhix --help
webhix serve --help
webhix endpoint create --help
```

#### `--version`, `-v`

Shows the Webhix binary version.

Example:

```sh
webhix --version
```

## Core commands

## `webhix serve`

Starts the Webhix HTTP server.

This command runs:

- The web UI.
- The webhook capture endpoint at `/r/<token>`.
- API routes used by the UI and CLI.
- Live request updates.
- SQLite storage.

Basic local usage:

```sh
webhix serve
```

Expected local output:

```text
Listening on http://localhost:8080
```

Production usage:

```sh
webhix serve --addr 0.0.0.0:8080 --base-url https://hooks.example.com --db-path /var/lib/webhix
```

### Flags

#### `--addr <addr>`, `-a <addr>`

Address to listen on.

Examples:

```sh
webhix serve --addr :8080
webhix serve -a 127.0.0.1:8080
webhix serve -a 0.0.0.0:9000
```

#### `--base-url <url>`

Public URL used for endpoint links.

Example:

```sh
webhix serve --base-url https://hooks.yourdomain.com
```

#### `--db-path <path>`

SQLite database directory.

Example:

```sh
webhix serve --db-path /data
```

#### `--password <value>`

Sets the single-user password.

Environment equivalent:

```sh
WEBHIX_PASSWORD
```

Example:

```sh
webhix serve --password "change-me"
```

This flag is convenient locally, but production deployments should prefer environment variables or `--password-file` to avoid leaking secrets through shell history or process lists.

#### `--password-file <path>`

Reads the single-user password from a file.

Example:

```sh
webhix serve --password-file /run/secrets/webhix_password
```

This is preferable for Docker, Docker Compose, Kubernetes, and systemd credentials.

#### `--trusted-proxies <cidr>[,<cidr>...]`

Trusts forwarded headers only from matching proxy networks.

Can be passed as a comma-separated list.

Example:

```sh
webhix serve --trusted-proxies 127.0.0.1/32,10.0.0.0/8 --base-url https://hooks.example.com
```

#### `--max-body-size <size>`

Maximum accepted request body size.

Examples:

```sh
webhix serve --max-body-size 1mb
webhix serve --max-body-size 25mb
webhix serve --max-body-size 100mb
```

This protects the instance from accidentally large webhook payloads.

#### `--retention <duration>`

How long captured requests are kept.

Examples:

```sh
webhix serve --retention 24h
webhix serve --retention 7d
webhix serve --retention 30d
webhix serve --retention 0
```

Recommended behavior:

- `0` disables automatic retention cleanup.
- Non-zero values periodically delete old requests.

#### `--max-requests <number>`

Maximum number of captured requests to retain.

Example:

```sh
webhix serve --max-requests 10000
```

This may be global or per endpoint. The implementation should document which one it uses.

#### `--readonly`

Starts Webhix in read-only mode.

In this mode the UI and API should allow inspection but block mutations such as:

- Creating endpoints.
- Deleting requests.
- Updating custom responses.
- Replaying requests.

Example:

```sh
webhix serve --readonly
```

Useful for debugging backups or sharing temporary read-only access.

#### `--open`

Opens the web UI in the default browser after startup.

Example:

```sh
webhix serve --open
```

This is primarily a local development convenience.

## `webhix version`

Prints version information.

Example:

```sh
webhix version
```

Recommended output:

```text
webhix 0.1.0
commit: abc123
built: 2026-05-17T10:00:00Z
go: go1.22.0
```

### Flags

#### `--output <text|json|yaml>`

Allows scripts to parse version details.

Example:

```sh
webhix version --output json
```

## Endpoint commands

Endpoint commands manage webhook receiver URLs.

An endpoint token maps to a public URL:

```text
<base-url>/r/<token>
```

Example:

```text
https://hooks.example.com/r/abc123
```

## `webhix endpoint create`

Creates a new endpoint.

Basic usage:

```sh
webhix endpoint create
```

Example output:

```text
Endpoint created
Token: abc123
URL:   http://localhost:8080/r/abc123
```

JSON usage:

```sh
webhix endpoint create --name stripe-dev --output json
```

### Flags

#### `--name <name>`

Sets a human-readable endpoint name.

Example:

```sh
webhix endpoint create --name stripe-dev
```

#### `--token <token>`

Uses a specific token instead of generating one.

Example:

```sh
webhix endpoint create --token stripe-local
```

This is useful for stable development URLs. The implementation should reject tokens that are already used or contain unsafe URL characters.

#### `--response-status <code>`

Sets the default HTTP status returned by the endpoint.

Example:

```sh
webhix endpoint create --response-status 202
```

#### `--response-header <key=value>`

Adds a header to the endpoint's custom response.

Can be repeated.

Example:

```sh
webhix endpoint create \
  --response-header "Content-Type=application/json" \
  --response-header "X-Mock=webhix"
```

#### `--response-body <text>`

Sets the response body directly.

Example:

```sh
webhix endpoint create --response-body '{"ok":true}'
```

#### `--response-body-file <path>`

Reads the response body from a file.

Example:

```sh
webhix endpoint create --response-body-file ./fixtures/ok.json
```

This is better than `--response-body` for multiline JSON or larger mock payloads.

#### `--output <text|json|yaml>`

Controls output format.

Example:

```sh
webhix endpoint create --output json
```

## `webhix endpoint list`

Lists endpoints.

Basic usage:

```sh
webhix endpoint list
```

Recommended table output:

```text
TOKEN       NAME         REQUESTS   LAST REQUEST
abc123      stripe-dev   42         2026-05-17 10:45:00
github-ci   github-ci    7          2026-05-17 09:12:30
```

### Flags

#### `--with-stats`

Includes endpoint statistics.

Example:

```sh
webhix endpoint list --with-stats
```

Stats may include:

- Request count.
- Last request timestamp.
- Total stored body size.
- Response status.

#### `--output <table|json|yaml>`

Controls output format.

Example:

```sh
webhix endpoint list --output json
```

## `webhix endpoint get <token>`

Shows details for one endpoint.

Example:

```sh
webhix endpoint get abc123
```

Recommended output:

```text
Token:           abc123
Name:            stripe-dev
URL:             http://localhost:8080/r/abc123
Response status: 200
Requests:        42
Created:         2026-05-17 10:00:00
Last request:    2026-05-17 10:45:00
```

### Flags

#### `--output <text|json|yaml>`

Controls output format.

Example:

```sh
webhix endpoint get abc123 --output json
```

## `webhix endpoint update <token>`

Updates endpoint settings.

Example:

```sh
webhix endpoint update abc123 --name stripe-staging --response-status 202
```

### Flags

#### `--name <name>`

Updates the display name.

#### `--response-status <code>`

Updates the custom response status.

#### `--response-header <key=value>`

Adds or replaces a response header.

Can be repeated.

Example:

```sh
webhix endpoint update abc123 \
  --response-header "Content-Type=application/json" \
  --response-header "X-Webhook-Inspector=webhix"
```

#### `--clear-response-headers`

Removes all configured custom response headers.

Example:

```sh
webhix endpoint update abc123 --clear-response-headers
```

#### `--response-body <text>`

Updates the response body.

#### `--response-body-file <path>`

Updates the response body from a file.

#### `--output <text|json|yaml>`

Controls output format.

## `webhix endpoint delete <token>`

Deletes an endpoint.

Example:

```sh
webhix endpoint delete abc123
```

Recommended default behavior:

- Ask for confirmation.
- Delete endpoint settings.
- Delete captured requests for that endpoint unless `--keep-requests` is provided.

### Flags

#### `--keep-requests`

Deletes the endpoint but preserves captured request history if the storage model supports that.

Example:

```sh
webhix endpoint delete abc123 --keep-requests
```

#### `--yes`, `-y`

Skips confirmation.

Example:

```sh
webhix endpoint delete abc123 --yes
```

## Request commands

Request commands inspect, stream, delete, clear, and export captured webhook requests.

Request identifiers should be stable and safe to copy. Suggested format:

```text
req_<id>
```

## `webhix request list <token>`

Lists captured requests for an endpoint.

Example:

```sh
webhix request list abc123
```

Recommended table output:

```text
ID        METHOD   PATH        CONTENT TYPE       SIZE    RECEIVED
req_001   POST     /r/abc123    application/json   2.1KB   2026-05-17 10:45:00
req_002   POST     /r/abc123    application/json   1.8KB   2026-05-17 10:46:10
```

### Flags

#### `--limit <number>`

Limits the number of returned requests.

Example:

```sh
webhix request list abc123 --limit 20
```

#### `--since <duration|timestamp>`

Returns requests received after a relative duration or timestamp.

Examples:

```sh
webhix request list abc123 --since 1h
webhix request list abc123 --since 2026-05-17T10:00:00Z
```

#### `--method <method>`

Filters by HTTP method.

Example:

```sh
webhix request list abc123 --method POST
```

#### `--status <code>`

Filters by response status if Webhix stores response status.

Example:

```sh
webhix request list abc123 --status 200
```

#### `--content-type <value>`

Filters by content type.

Example:

```sh
webhix request list abc123 --content-type application/json
```

#### `--ip <address>`

Filters by sender IP.

Example:

```sh
webhix request list abc123 --ip 203.0.113.10
```

#### `--search <text>`

Searches captured request data.

Suggested search scope:

- Headers.
- Query parameters.
- Body text.
- Path.

Example:

```sh
webhix request list abc123 --search checkout.session.completed
```

#### `--output <table|json|yaml>`

Controls output format.

Example:

```sh
webhix request list abc123 --output json
```

## `webhix request get <request-id>`

Shows one captured request in detail.

Example:

```sh
webhix request get req_001
```

Recommended output should include:

- Request ID.
- Endpoint token.
- Timestamp.
- Sender IP.
- Method.
- Path.
- Query parameters.
- Headers.
- Body.
- Content type.
- Size.

### Flags

#### `--body-only`

Prints only the request body.

Example:

```sh
webhix request get req_001 --body-only
```

#### `--headers-only`

Prints only request headers.

Example:

```sh
webhix request get req_001 --headers-only
```

#### `--raw`

Prints a raw-style representation with minimal formatting.

Example:

```sh
webhix request get req_001 --raw
```

#### `--redact <target>`

Redacts sensitive values from output.

Can be repeated.

Examples:

```sh
webhix request get req_001 --redact header:Authorization
webhix request get req_001 --redact header:X-Signature --redact query:token
```

Suggested target formats:

- `header:<name>`
- `query:<name>`
- `body:<json.path>`

#### `--output <text|json|yaml>`

Controls output format.

## `webhix request tail <token>`

Streams incoming requests for an endpoint in real time.

Example:

```sh
webhix request tail abc123
```

This is the terminal equivalent of watching the live UI.

### Flags

#### `--body`

Includes request body in streamed output.

Example:

```sh
webhix request tail abc123 --body
```

#### `--headers`

Includes request headers.

Example:

```sh
webhix request tail abc123 --headers
```

#### `--pretty`

Pretty-prints JSON bodies when possible.

Example:

```sh
webhix request tail abc123 --body --pretty
```

#### `--filter <expr>`

Filters streamed requests.

Examples:

```sh
webhix request tail abc123 --filter method=POST
webhix request tail abc123 --filter header:X-GitHub-Event=push
webhix request tail abc123 --filter content-type=application/json
```

The first implementation can support simple equality expressions. More complex filtering can be added later.

## `webhix request delete <request-id>`

Deletes one captured request.

Example:

```sh
webhix request delete req_001
```

### Flags

#### `--yes`, `-y`

Skips confirmation.

Example:

```sh
webhix request delete req_001 --yes
```

## `webhix request clear <token>`

Deletes captured requests for an endpoint.

Example:

```sh
webhix request clear abc123
```

### Flags

#### `--before <duration|timestamp>`

Deletes only requests older than the specified duration or timestamp.

Examples:

```sh
webhix request clear abc123 --before 7d
webhix request clear abc123 --before 2026-05-01T00:00:00Z
```

#### `--yes`, `-y`

Skips confirmation.

Example:

```sh
webhix request clear abc123 --yes
```

## Forwarding commands

Forwarding connects captured incoming webhooks to a local development service.

The README already defines this important command:

```sh
webhix forward <token> --to localhost:3000
```

## `webhix forward <token> --to <url>`

Forwards incoming requests for an endpoint to another HTTP server.

Example:

```sh
webhix forward abc123 --to localhost:3000
```

If the local service expects a full URL:

```sh
webhix forward abc123 --to http://localhost:3000/webhooks/stripe
```

### Flags

#### `--to <url>`

Target URL for forwarded requests.

Required.

Examples:

```sh
webhix forward abc123 --to localhost:3000
webhix forward abc123 --to http://localhost:3000/webhook
webhix forward abc123 --to https://localhost:8443/webhook
```

If the scheme is omitted, Webhix can default to `http://`.

#### `--server <url>`

Webhix server URL to connect to.

Example:

```sh
webhix forward abc123 --server https://hooks.example.com --to localhost:3000
```

This allows the CLI to run on a developer machine while the Webhix instance runs elsewhere.

#### `--auth-token <token>`

Authentication token for connecting to a remote Webhix server.

Example:

```sh
webhix forward abc123 --server https://hooks.example.com --auth-token "$WEBHIX_TOKEN" --to localhost:3000
```

For production, environment variables or token files are safer than literal shell arguments.

#### `--auth-token-file <path>`

Reads the auth token from a file.

Example:

```sh
webhix forward abc123 --auth-token-file ~/.config/webhix/token --to localhost:3000
```

#### `--rewrite-host`

Rewrites the `Host` header to match the forwarding target.

Example:

```sh
webhix forward abc123 --to localhost:3000 --rewrite-host
```

Useful when local frameworks route or validate by host.

#### `--preserve-host`

Preserves the original `Host` header.

Example:

```sh
webhix forward abc123 --to localhost:3000 --preserve-host
```

This is useful when the target service wants to see the original public host.

`--rewrite-host` and `--preserve-host` should be mutually exclusive.

#### `--path <path>`

Overrides the forwarded request path.

Example:

```sh
webhix forward abc123 --to localhost:3000 --path /webhooks/stripe
```

This lets a public Webhix endpoint forward to a specific local application route.

#### `--header <key=value>`

Adds or replaces a header in forwarded requests.

Can be repeated.

Example:

```sh
webhix forward abc123 --to localhost:3000 --header "X-Forwarded-By=webhix"
```

#### `--drop-header <key>`

Removes a header before forwarding.

Can be repeated.

Example:

```sh
webhix forward abc123 --to localhost:3000 --drop-header Authorization
```

#### `--timeout <duration>`

Timeout for each forwarded request.

Example:

```sh
webhix forward abc123 --to localhost:3000 --timeout 10s
```

#### `--retry <number>`

Number of retry attempts when forwarding fails.

Example:

```sh
webhix forward abc123 --to localhost:3000 --retry 2
```

Retries should be used carefully because some webhook handlers are not idempotent.

#### `--print`

Prints one line for each forwarded request.

Example:

```sh
webhix forward abc123 --to localhost:3000 --print
```

#### `--print-body`

Prints forwarded request bodies.

Example:

```sh
webhix forward abc123 --to localhost:3000 --print-body
```

This may expose secrets in terminal logs.

#### `--insecure`

Allows forwarding to HTTPS targets with invalid or self-signed certificates.

Example:

```sh
webhix forward abc123 --to https://localhost:8443 --insecure
```

This should be clearly marked as a local-development option.

## Replay commands

Replay sends a captured request again. It can target the original destination, a custom URL, or a local service.

## `webhix replay <request-id>`

Replays one captured request.

Example:

```sh
webhix replay req_001 --to http://localhost:3000/webhook
```

### Flags

#### `--to <url>`

Target URL for replay.

Example:

```sh
webhix replay req_001 --to http://localhost:3000/webhooks/stripe
```

If omitted, Webhix may use the original request URL if it was stored and is safe to replay.

#### `--method <method>`

Overrides the HTTP method.

Example:

```sh
webhix replay req_001 --method POST --to http://localhost:3000/webhook
```

#### `--header <key=value>`

Adds or replaces a header.

Can be repeated.

Example:

```sh
webhix replay req_001 --header "X-Replayed-By=webhix"
```

#### `--drop-header <key>`

Removes a header from the replayed request.

Can be repeated.

Example:

```sh
webhix replay req_001 --drop-header Stripe-Signature
```

#### `--body <text>`

Overrides the request body.

Example:

```sh
webhix replay req_001 --body '{"test":true}' --to http://localhost:3000/webhook
```

#### `--body-file <path>`

Overrides the request body from a file.

Example:

```sh
webhix replay req_001 --body-file ./fixtures/event.json --to http://localhost:3000/webhook
```

#### `--timeout <duration>`

Timeout for the replayed request.

Example:

```sh
webhix replay req_001 --timeout 5s
```

#### `--repeat <number>`

Sends the replay multiple times.

Example:

```sh
webhix replay req_001 --repeat 10 --to http://localhost:3000/webhook
```

#### `--interval <duration>`

Wait time between repeated replays.

Example:

```sh
webhix replay req_001 --repeat 5 --interval 2s --to http://localhost:3000/webhook
```

#### `--dry-run`

Prints the replay request without sending it.

Example:

```sh
webhix replay req_001 --to http://localhost:3000/webhook --dry-run
```

This is useful before replaying requests that may trigger side effects.

## Export commands

Export converts captured requests into reusable representations.

Supported formats should include:

- `curl`
- `httpie`
- `raw`
- `json`

## `webhix export <request-id>`

Exports a captured request.

Example:

```sh
webhix export req_001 --format curl
```

Example output:

```sh
curl -X POST 'http://localhost:8080/r/abc123' \
  -H 'Content-Type: application/json' \
  --data '{"ok":true}'
```

### Flags

#### `--format <curl|httpie|raw|json>`

Export format.

Examples:

```sh
webhix export req_001 --format curl
webhix export req_001 --format httpie
webhix export req_001 --format json
```

#### `--output-file <path>`, `-o <path>`

Writes export output to a file.

Example:

```sh
webhix export req_001 --format curl -o replay.sh
```

#### `--include-headers`

Includes request headers.

Example:

```sh
webhix export req_001 --format json --include-headers
```

#### `--include-body`

Includes request body.

Example:

```sh
webhix export req_001 --format json --include-body
```

#### `--redact <target>`

Redacts sensitive values from exported output.

Can be repeated.

Examples:

```sh
webhix export req_001 --format curl --redact header:Authorization
webhix export req_001 --format json --redact query:token --redact body:data.object.customer_email
```

#### `--pretty`

Pretty-prints JSON where possible.

Example:

```sh
webhix export req_001 --format json --pretty
```

## Auth commands

Current README describes single-user authentication. Webhix can generate a password on first run or accept one from the environment.

The CLI should support explicit password management.

## `webhix auth set-password`

Sets or changes the single-user password.

Examples:

```sh
webhix auth set-password
webhix auth set-password --password "change-me"
webhix auth set-password --password-file /run/secrets/webhix_password
webhix auth set-password --generate
```

### Flags

#### `--password <value>`

Sets the password directly.

This is convenient but less secure than `--password-file` or an interactive prompt.

#### `--password-file <path>`

Reads the password from a file.

Example:

```sh
webhix auth set-password --password-file ./secret.txt
```

#### `--generate`

Generates a secure random password and prints it once.

Example:

```sh
webhix auth set-password --generate
```

Recommended output:

```text
Password updated.
Generated password: <value>
Store this password now. It will not be shown again.
```

## User commands

These commands belong to the v0.2 multi-user and RBAC roadmap.

## `webhix user create <username>`

Creates a user.

Example:

```sh
webhix user create alice --role admin
```

### Flags

#### `--password <value>`

Sets the user's password.

#### `--password-file <path>`

Reads the password from a file.

#### `--role <admin|viewer|operator>`

Sets the user's role.

Suggested roles:

- `admin`: full access.
- `operator`: can inspect, replay, forward, and manage endpoints.
- `viewer`: read-only access to UI and request data.

#### `--disabled`

Creates the user in a disabled state.

## `webhix user list`

Lists users.

Example:

```sh
webhix user list
```

### Flags

#### `--output <table|json|yaml>`

Controls output format.

## `webhix user update <username>`

Updates a user.

Examples:

```sh
webhix user update alice --role viewer
webhix user update alice --disabled
webhix user update alice --password-file ./new-password.txt
```

### Flags

#### `--password <value>`

Changes the password.

#### `--password-file <path>`

Changes the password from a file.

#### `--role <admin|viewer|operator>`

Changes the role.

#### `--disabled`

Disables the user.

#### `--enabled`

Enables the user.

## `webhix user delete <username>`

Deletes a user.

Example:

```sh
webhix user delete alice
```

### Flags

#### `--yes`, `-y`

Skips confirmation.

## Database commands

Database commands operate on the local Webhix storage.

The default storage is SQLite in `WEBHIX_DB_PATH`.

## `webhix db migrate`

Applies database migrations.

Example:

```sh
webhix db migrate
```

This command is useful for:

- Manual upgrades.
- Deployment scripts.
- CI checks.
- Debugging migration failures.

### Flags

#### `--dry-run`

Shows pending migrations without applying them.

Example:

```sh
webhix db migrate --dry-run
```

#### `--target <version>`

Migrates to a specific schema version.

Example:

```sh
webhix db migrate --target 20260515141636
```

## `webhix db status`

Shows database status.

Example:

```sh
webhix db status
```

Recommended output:

```text
Path:       ./data
Driver:     sqlite
Schema:     20260515141636
Size:       12.4MB
Endpoints:  5
Requests:   1204
```

### Flags

#### `--output <text|json|yaml>`

Controls output format.

## `webhix db backup <path>`

Creates a database backup.

Example:

```sh
webhix db backup ./backups/webhix-2026-05-17.sqlite
```

### Flags

#### `--compress`

Compresses the backup.

Example:

```sh
webhix db backup ./backups/webhix.sqlite.gz --compress
```

#### `--overwrite`

Allows overwriting an existing backup file.

Example:

```sh
webhix db backup ./backups/latest.sqlite --overwrite
```

## `webhix db restore <path>`

Restores a database backup.

Example:

```sh
webhix db restore ./backups/webhix-2026-05-17.sqlite
```

### Flags

#### `--yes`, `-y`

Skips confirmation.

Example:

```sh
webhix db restore ./backups/webhix.sqlite --yes
```

## `webhix db vacuum`

Runs SQLite vacuum to reclaim space and optimize the database after large deletes.

Example:

```sh
webhix db vacuum
```

## `webhix db prune`

Deletes old captured requests.

Example:

```sh
webhix db prune --older-than 30d
```

### Flags

#### `--older-than <duration>`

Deletes requests older than the specified duration.

Examples:

```sh
webhix db prune --older-than 7d
webhix db prune --older-than 90d
```

#### `--endpoint <token>`

Limits pruning to one endpoint.

Example:

```sh
webhix db prune --endpoint abc123 --older-than 7d
```

#### `--yes`, `-y`

Skips confirmation.

## Config and diagnostic commands

These commands help users understand how Webhix is configured and whether the runtime environment is healthy.

## `webhix config print`

Prints the effective configuration.

Example:

```sh
webhix config print
```

Recommended output:

```text
addr:     :8080
base_url: http://localhost:8080
db_path:  ./data
```

### Flags

#### `--redact`

Redacts secrets.

Example:

```sh
webhix config print --redact
```

This should be enabled by default if secrets are included in config output.

#### `--output <text|json|yaml>`

Controls output format.

Example:

```sh
webhix config print --output json
```

## `webhix config validate`

Validates configuration.

Example:

```sh
webhix config validate
```

Checks should include:

- `WEBHIX_ADDR` parses as a listen address.
- `WEBHIX_BASE_URL` is a valid URL.
- `WEBHIX_DB_PATH` exists or can be created.
- Database directory is writable.
- Password configuration is valid if auth is enabled.
- Reverse proxy settings are consistent.

## `webhix doctor`

Runs a diagnostic check.

Example:

```sh
webhix doctor
```

Suggested checks:

- Binary version.
- Operating system and architecture.
- Go build metadata.
- Config file loaded or missing.
- Effective listen address.
- Effective base URL.
- Database path.
- Database access.
- Migration status.
- Disk free space.
- Whether configured port can be bound.
- Reverse proxy hints.

### Flags

#### `--verbose`

Prints detailed diagnostic output.

Example:

```sh
webhix doctor --verbose
```

#### `--output <text|json>`

Controls output format.

Example:

```sh
webhix doctor --output json
```

## Signature verification commands

These commands belong to the v0.2 roadmap.

They support testing webhook signature verification for providers such as Stripe and GitHub.

## `webhix verify test <provider>`

Tests a webhook signature.

Examples:

```sh
webhix verify test stripe --request req_001 --secret-file ./stripe_secret.txt
webhix verify test github --payload-file payload.json --secret "$GITHUB_WEBHOOK_SECRET" --header "X-Hub-Signature-256=sha256=..."
```

### Arguments

#### `<provider>`

Verification provider.

Suggested values:

- `stripe`
- `github`
- `generic-hmac`

### Flags

#### `--request <request-id>`

Uses a captured request as input.

Example:

```sh
webhix verify test stripe --request req_001 --secret-file ./secret.txt
```

#### `--payload-file <path>`

Uses a local payload file.

Example:

```sh
webhix verify test github --payload-file ./payload.json --secret-file ./secret.txt
```

#### `--secret <value>`

Secret used for verification.

#### `--secret-file <path>`

Reads the secret from a file.

Recommended for production-like use.

#### `--header <key=value>`

Supplies signature-related headers.

Can be repeated.

Example:

```sh
webhix verify test github \
  --payload-file ./payload.json \
  --secret-file ./secret.txt \
  --header "X-Hub-Signature-256=sha256=abc"
```

#### `--algorithm <sha256|sha1>`

Algorithm for `generic-hmac`.

Example:

```sh
webhix verify test generic-hmac --algorithm sha256 --payload-file payload.json --secret-file secret.txt
```

## Schema validation commands

These commands belong to the v0.2 roadmap.

They allow endpoints to validate incoming JSON payloads against schemas.

## `webhix schema attach <token> <schema-file>`

Attaches a schema to an endpoint.

Example:

```sh
webhix schema attach abc123 ./schemas/stripe-event.schema.json
```

### Flags

#### `--name <name>`

Human-readable schema name.

Example:

```sh
webhix schema attach abc123 ./schema.json --name stripe-event
```

#### `--mode <warn|reject>`

Validation mode.

`warn` stores validation failures but still accepts the request:

```sh
webhix schema attach abc123 ./schema.json --mode warn
```

`reject` returns an error to the webhook sender when validation fails:

```sh
webhix schema attach abc123 ./schema.json --mode reject
```

Recommended default:

```text
warn
```

## `webhix schema list <token>`

Lists schemas attached to an endpoint.

Example:

```sh
webhix schema list abc123
```

### Flags

#### `--output <table|json|yaml>`

Controls output format.

## `webhix schema remove <token> <schema-id>`

Removes a schema from an endpoint.

Example:

```sh
webhix schema remove abc123 schema_001
```

### Flags

#### `--yes`, `-y`

Skips confirmation.

## `webhix schema test <schema-file>`

Tests a payload against a schema.

Example:

```sh
webhix schema test ./schema.json --payload-file ./payload.json
```

### Flags

#### `--payload-file <path>`

Path to a payload file.

Required unless payload is read from stdin.

Example:

```sh
webhix schema test ./schema.json --payload-file ./payload.json
```

#### `--request <request-id>`

Uses a captured request body as the payload.

Example:

```sh
webhix schema test ./schema.json --request req_001
```

## Notification commands

These commands belong to the v0.2 roadmap.

They configure notifications for new requests or validation outcomes.

Supported providers may include:

- Slack.
- Telegram.
- Discord.

## `webhix notify add <token> <provider>`

Adds a notification target to an endpoint.

Examples:

```sh
webhix notify add abc123 slack --webhook-url https://hooks.slack.com/services/...
webhix notify add abc123 discord --webhook-url https://discord.com/api/webhooks/...
```

### Arguments

#### `<token>`

Endpoint token.

#### `<provider>`

Notification provider.

Suggested values:

- `slack`
- `telegram`
- `discord`

### Flags

#### `--webhook-url <url>`

Provider webhook URL.

Example:

```sh
webhix notify add abc123 slack --webhook-url https://hooks.slack.com/services/...
```

#### `--secret-file <path>`

Reads provider credentials from a file.

Example:

```sh
webhix notify add abc123 telegram --secret-file ./telegram-token.txt
```

#### `--on <any|valid|invalid|error>`

Controls when notification fires.

Examples:

```sh
webhix notify add abc123 slack --webhook-url "$SLACK_URL" --on any
webhix notify add abc123 slack --webhook-url "$SLACK_URL" --on invalid
```

Suggested meanings:

- `any`: notify for every request.
- `valid`: notify only for requests that pass validation.
- `invalid`: notify only for requests that fail validation.
- `error`: notify only when Webhix fails to process something.

#### `--template <path>`

Path to a message template.

Example:

```sh
webhix notify add abc123 slack --webhook-url "$SLACK_URL" --template ./slack-template.md
```

## `webhix notify list <token>`

Lists notification targets for an endpoint.

Example:

```sh
webhix notify list abc123
```

### Flags

#### `--output <table|json|yaml>`

Controls output format.

## `webhix notify test <notification-id>`

Sends a test notification.

Example:

```sh
webhix notify test notif_001
```

### Flags

#### `--message <text>`

Custom test message.

Example:

```sh
webhix notify test notif_001 --message "Webhix test notification"
```

## `webhix notify remove <notification-id>`

Removes a notification target.

Example:

```sh
webhix notify remove notif_001
```

### Flags

#### `--yes`, `-y`

Skips confirmation.

## Tunnel command

Tunnel mode belongs to the v0.3+ roadmap.

It exposes a local port via a public URL on the Webhix managed relay (`webhix.online`),
without requiring a VPS or reverse proxy. Protocol spec: `docs/tunnel-protocol.md`.

## `webhix tunnel <port>`

Exposes a local port to the internet.

Examples:

```sh
webhix tunnel 3000
webhix tunnel 3000 --subdomain myapp
webhix tunnel 3000 --auth-token tok_xxx
webhix tunnel 3000 --relay wss://my.relay/tunnel
```

### Flags

#### `--relay <url>`

Relay server WebSocket URL. Defaults to `wss://relay.webhix.online/tunnel`.

Example:

```sh
webhix tunnel 3000 --relay wss://tunnel.mycompany.com/tunnel
```

#### `--auth-token <token>`

Pro auth token from webhix.online dashboard (env: `WEBHIX_TUNNEL_TOKEN`).

Example:

```sh
webhix tunnel 3000 --auth-token tok_xxx
```

#### `--subdomain <name>`

Reserved subdomain to request (Pro only).

Example:

```sh
webhix tunnel 3000 --subdomain myapp
```

## Suggested command groups

The Cobra command tree should group commands by user task.

### Core

```sh
webhix serve
webhix version
webhix help
```

### Endpoints

```sh
webhix endpoint create
webhix endpoint list
webhix endpoint get <token>
webhix endpoint update <token>
webhix endpoint delete <token>
```

### Requests

```sh
webhix request list <token>
webhix request get <request-id>
webhix request tail <token>
webhix request delete <request-id>
webhix request clear <token>
```

### Forwarding, replay, and export

```sh
webhix forward <token> --to <url>
webhix replay <request-id>
webhix export <request-id>
```

### Auth and users

```sh
webhix auth set-password
webhix user create <username>
webhix user list
webhix user update <username>
webhix user delete <username>
```

### Database

```sh
webhix db migrate
webhix db status
webhix db backup <path>
webhix db restore <path>
webhix db vacuum
webhix db prune
```

### Config and diagnostics

```sh
webhix config print
webhix config validate
webhix doctor
```

### Roadmap v0.2

```sh
webhix verify test <provider>
webhix schema attach <token> <schema-file>
webhix schema list <token>
webhix schema remove <token> <schema-id>
webhix schema test <schema-file>
webhix notify add <token> <provider>
webhix notify list <token>
webhix notify test <notification-id>
webhix notify remove <notification-id>
```

### Roadmap v0.3+

```sh
webhix tunnel
```

## Minimal first implementation

The full CLI can be implemented gradually. A practical first milestone should include:

```sh
webhix serve
webhix version
webhix config print
webhix config validate
webhix doctor
```

Then add endpoint and request operations:

```sh
webhix endpoint create
webhix endpoint list
webhix endpoint get <token>
webhix request list <token>
webhix request get <request-id>
webhix request tail <token>
```

Then add developer workflow commands:

```sh
webhix forward <token> --to <url>
webhix replay <request-id>
webhix export <request-id>
```

Finally add maintenance and roadmap commands:

```sh
webhix db migrate
webhix db status
webhix db backup <path>
webhix db restore <path>
webhix db vacuum
webhix db prune
webhix auth set-password
webhix user create <username>
webhix user list
webhix user update <username>
webhix user delete <username>
webhix verify test <provider>
webhix schema attach <token> <schema-file>
webhix notify add <token> <provider>
webhix tunnel
```

## Naming notes

Recommended names:

- Use singular command groups: `endpoint`, `request`, `schema`, `notify`.
- Use verbs for actions: `create`, `list`, `get`, `update`, `delete`, `clear`, `tail`.
- Keep `forward`, `replay`, and `export` as top-level commands because they are core workflows, not only subcommands of `request`.
- Prefer `--to` for target URLs because it reads naturally in `forward` and `replay`.
- Prefer `--output` over format-specific flags for structured command output.
- Prefer `--format` when the command is converting data, such as `webhix export --format curl`.

## Safety notes

Commands that may expose secrets:

```sh
webhix request get
webhix request tail --body
webhix export
webhix replay --dry-run
webhix forward --print-body
webhix config print
```

These commands should support redaction where relevant and should avoid printing secrets by default when a safer default is possible.

Commands that should ask for confirmation by default:

```sh
webhix endpoint delete <token>
webhix request delete <request-id>
webhix request clear <token>
webhix db restore <path>
webhix db prune
webhix user delete <username>
webhix schema remove <token> <schema-id>
webhix notify remove <notification-id>
```

Commands that should support `--yes`:

```sh
webhix endpoint delete <token> --yes
webhix request delete <request-id> --yes
webhix request clear <token> --yes
webhix db restore <path> --yes
webhix db prune --yes
webhix user delete <username> --yes
webhix schema remove <token> <schema-id> --yes
webhix notify remove <notification-id> --yes
```
