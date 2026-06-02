---
generated: true
source_command: "$project-skill-builder create"
source_paths:
  - "AGENTS.md"
  - "README.md"
  - "go.mod"
  - "main.go"
  - "router/main.go"
  - "model/main.go"
  - "common/json.go"
---

# Backend Module

Last verified: 2026-06-02
Verified against commit: 7aaa5332

## Main Files

- `main.go`
- `router/`
- `controller/`
- `service/`
- `model/`
- `middleware/`
- `setting/`
- `common/`
- `dto/`
- `types/`

## Responsibility

The backend is a Go API gateway and management service. It exposes dashboard APIs, OpenAI-compatible and provider-specific relay APIs, user and token management, billing, rate limiting, setup, OAuth, passkeys, task polling, and embedded frontend assets.

The source of truth for the Go version is `go.mod`, which currently declares `go 1.25.1`.

## Startup Path

- `main.go` embeds `web/default/dist` and `web/classic/dist`.
- `InitResources()` loads `.env`, initializes environment settings, logging, ratio settings, HTTP client, token encoders, SQL databases, option map, pricing, Redis, performance metrics, system monitoring, i18n, and custom OAuth providers.
- `main()` starts recurring tasks for cache sync, options sync, quota data, channel tests, channel upstream model updates, Codex credential refresh, subscription quota resets, task polling, optional batch updates, optional pprof, and optional Pyroscope.
- Gin is configured with panic recovery, request IDs, Powered-By, i18n middleware, logging, sessions, analytics injection, and the router tree.

## Layer Boundaries

- `router/` wires URL groups and middleware.
- `controller/` handles HTTP request/response details and orchestrates services.
- `service/` contains business logic, billing sessions, quota, channel selection, task billing, OAuth helpers, token counting, and provider-independent conversion logic.
- `model/` owns GORM models, migrations, database access, cache-backed lookups, and persistence compatibility.
- `relay/` owns upstream request conversion, provider adaptors, streaming, task relay, and upstream response handling.
- `setting/` owns runtime configuration groups and ratio/billing settings.

## Invariants

- Business JSON marshal/unmarshal calls should go through `common/json.go` wrappers.
- Persistence must support SQLite, MySQL, and PostgreSQL simultaneously.
- Embedded frontend assets mean production backend builds expect both theme dist directories to exist or be intentionally handled by build tooling.
- `common.IsMasterNode` gates migrations and several recurring tasks.
- `RedisEnabled` forces memory cache compatibility at startup.

## Tests

Focused Go tests exist across `common/`, `controller/`, `dto/`, `middleware/`, `model/`, `pkg/billingexpr/`, `relay/`, `service/`, and `setting/`.

Use focused packages first, for example:

```bash
go test ./service ./relay/common ./pkg/billingexpr
```
