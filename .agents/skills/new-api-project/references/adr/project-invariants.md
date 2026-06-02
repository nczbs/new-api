---
generated: true
source_command: "$project-skill-builder create"
source_paths:
  - "AGENTS.md"
  - "common/json.go"
  - "model/main.go"
  - "dto/openai_request_zero_value_test.go"
  - "pkg/billingexpr/expr.md"
  - "web/default/AGENTS.md"
---

# ADR: Project Invariants For Safe Changes

Last verified: 2026-06-02
Verified against commit: 7aaa5332

## Context

new-api spans a Go backend, provider relay adapters, multiple database engines, billing, and two frontend themes. Small local changes can break upstream compatibility, billing correctness, database compatibility, or frontend i18n.

## Decision

Project work should preserve these invariants:

- Protected project and organization identifiers must not be removed, renamed, or replaced.
- JSON business operations use wrappers in `common/json.go`.
- Database code must support SQLite, MySQL, and PostgreSQL.
- Optional upstream relay DTO scalar fields use pointer types with `omitempty` so explicit zero and false values survive re-marshaling.
- New relay channels must verify `StreamOptions` support before adding to `streamSupportedChannels`.
- Tiered billing work must start from `pkg/billingexpr/expr.md`.
- Default frontend work must follow `web/default/AGENTS.md`, use Bun, and keep i18n complete across supported locales.

## Consequences

- Implementations may need slightly more explicit types and compatibility branches.
- Raw SQL and direct JSON operations should be rare and justified.
- Tests should focus on compatibility edges: zero values, DB behavior, token normalization, relay conversion, i18n, and auth boundaries.

## Verification Sources

- `AGENTS.md`
- `common/json.go`
- `model/main.go`
- `dto/openai_request_zero_value_test.go`
- `pkg/billingexpr/expr.md`
- `web/default/AGENTS.md`
