---
generated: true
source_command: "$project-skill-builder create"
source_paths:
  - "main.go"
  - ".env.example"
  - "docker-compose.dev.yml"
  - "web/default/package.json"
  - "web/default/rsbuild.config.ts"
  - "model/main.go"
  - "router/web-router.go"
---

# Development Troubleshooting

Last verified: 2026-06-02
Verified against commit: 7aaa5332

## Backend Cannot Start

Checks:

```bash
go test ./...
```

If the error mentions missing embedded assets, check whether `web/default/dist` and `web/classic/dist` exist. `main.go` embeds both theme dist directories.

If the error is database-related, inspect `SQL_DSN`, `SQLITE_PATH`, `LOG_SQL_DSN`, and the database selection rules in `model/main.go`.

## Frontend Dev Server Cannot Reach Backend

Checks from `web/default/`:

```bash
bun run dev
```

`web/default/rsbuild.config.ts` proxies `/api`, `/mj`, and `/pg` to `http://localhost:3000` unless `VITE_REACT_APP_SERVER_URL` overrides it.

## TypeScript Or Build Failure

Checks from `web/default/`:

```bash
bun run typecheck
bun run build:check
```

For user-visible text failures or missing locale keys, use `$i18n-translate` and run `bun run i18n:sync`.

## Relay Request Fails Before Provider Call

Likely areas:

- `middleware.TokenAuth()`
- `middleware.ModelRequestRateLimit()`
- `middleware.Distribute()`
- token model restrictions
- group/channel eligibility
- request body storage and validation
- pre-consume billing

Start from `router/relay-router.go`, `middleware/distributor.go`, and `controller/relay.go`.

## Tiered Billing Looks Wrong

Read `pkg/billingexpr/expr.md` first. Then check:

- whether the model is in `tiered_expr` mode
- expression variables used by the compiled expression
- `BillingSnapshot` captured during pre-consume
- `BuildTieredTokenParams()` normalization
- whether tier conditions use `len` instead of `p`

Focused tests:

```bash
go test ./pkg/billingexpr ./service ./relay/helper
```

## Database Compatibility Concern

Before changing migrations or raw SQL, verify:

- SQLite support
- MySQL support
- PostgreSQL support
- quoting for reserved columns such as `group` and `key`
- boolean literal differences

Prefer focused model or controller tests over manual DB assumptions.
