# Git Commit Conventions

## Format

```
<type>(<scope>): <summary>

(optional) Body — explain WHY, not what
```

## Types

| Type       | When to use                                |
| ---------- | ------------------------------------------ |
| `feat`     | New feature                                |
| `fix`      | Bug fix                                    |
| `refactor` | Code change with no behavior change        |
| `perf`     | Performance improvement                    |
| `style`    | Formatting, whitespace (no logic change)   |
| `docs`     | Documentation only                         |
| `test`     | Adding or fixing tests                     |
| `chore`    | Dependencies, configs, tooling             |
| `ci`       | CI/CD pipeline changes                     |
| `build`    | Build system changes                       |
| `db`       | Migrations, schema changes                 |
| `revert`   | Revert a previous commit                   |
| `init`     | Initial project scaffold                   |
| `wip`      | Work in progress — avoid merging into main |

## Scope (optional)

Module or layer affected, in parentheses after the type:

```
feat(auth): add jwt middleware
fix(db): handle connection timeout
```

Common scopes: `auth`, `user`, `db`, `config`, `router`, `middleware`, `ci`, `logger`

## Summary rules

- English, imperative mood: `add`, `fix`, `remove` — not `added`, `fixed`
- Max 72 characters
- No period at the end
- Answer: "What does this commit do?"

## Examples

```
init: project scaffold with health endpoint
feat(auth): add register and login endpoints
feat(user): add get profile endpoint
fix(db): handle nil connection on startup
refactor(router): extract middleware registration
ci: add deploy workflow with ghcr
docs: add commit and branch naming conventions
chore: update golangci-lint to v2.5.0
db: add users table migration
test(auth): add unit tests for jwt token creation
```

## Anti-patterns

| Bad                 | Why                        |
| ------------------- | -------------------------- |
| `fix: bug`          | Not descriptive            |
| `update` / `commit` | Meaningless                |
| `add new file`      | Which file? Why?           |
| `wip: test`         | Do not merge WIP into main |
| `debug`             | Remove before merging      |
