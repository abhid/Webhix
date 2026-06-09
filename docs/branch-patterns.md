# Git Branch Conventions

## Format

```
<type>/<short-description>
```

With task tracker:

```
<type>/<ticket-id>-<short-description>
```

## Types

| Type       | When to use                          | Examples                                          |
| ---------- | ------------------------------------ | ------------------------------------------------- |
| `feature`  | New functionality                    | `feature/user-auth`, `feature/add-exercises`      |
| `fix`      | Bug fix                              | `fix/login-error`, `fix/db-connection-timeout`    |
| `refactor` | Code improvement, no behavior change | `refactor/simplify-auth-handler`                  |
| `chore`    | Dependencies, configs, tooling       | `chore/update-dependencies`, `chore/setup-linter` |
| `docs`     | Documentation updates                | `docs/update-readme`, `docs/add-api-guide`        |
| `test`     | Adding or fixing tests               | `test/auth-unit-tests`, `test/fix-db-mocks`       |
| `ci`       | CI/CD configuration                  | `ci/add-deploy-workflow`, `ci/fix-build`          |
| `db`       | Migrations, schema changes           | `db/add-users-table`, `db/add-progress-index`     |
| `hotfix`   | Urgent production fix                | `hotfix/fix-payment-crash`                        |
| `release`  | Release preparation                  | `release/1.0.0`                                   |
| `wip`      | Temporary / experimental branch      | `wip/try-new-cache-strategy`                      |

## Rules

- Use hyphens, not underscores: `feature/add-user-login`, not `feature/add_user_login`
- Keep names short and descriptive
- Lowercase only

## Examples

```
feature/user-registration
feature/get-exercises
fix/jwt-expiration-parsing
refactor/split-usecase-layer
db/add-users-migration
ci/add-registry-token
docs/add-commit-conventions
chore/update-golangci-lint
hotfix/nil-pointer-on-startup
```
