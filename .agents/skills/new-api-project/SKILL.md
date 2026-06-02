---
name: new-api-project
description: Use this skill when working inside /root/project/new-api, including backend gateway architecture, relay adapters, billing, database compatibility, web/default frontend, runtime, tests, and project-specific conventions.
---

# new-api Project

Use this skill when working inside `/root/project/new-api` or when the user asks about new-api architecture, runtime, relay providers, billing, persistence, routing, or the default frontend.

## Location

This project skill is repository-scoped at `/root/project/new-api/.agents/skills/new-api-project/`.

## First Reads

- Always read `AGENTS.md` before changing project behavior.
- For `web/default` changes, also read `web/default/AGENTS.md`.
- For tiered or expression-based billing, read `pkg/billingexpr/expr.md` before editing.
- For frontend i18n work, use `$i18n-translate`.
- For shadcn/ui or registry work, use `$shadcn-ui`.

## Work Protocol

- Treat source code and tests as the source of truth. Use README and generated docs as supporting context only.
- Keep edits scoped to the requested module and preserve existing layer boundaries: router -> controller -> service -> model, with relay provider adapters under `relay/`.
- Do not rename, remove, or replace protected project or organization identifiers.
- Use `common.Marshal`, `common.Unmarshal`, `common.UnmarshalJsonStr`, and `common.DecodeJson` for JSON operations in business code.
- Keep database code compatible with SQLite, MySQL, and PostgreSQL; prefer GORM and use DB-specific helpers from `model/main.go` when raw SQL is unavoidable.
- Preserve explicit zero values in upstream relay request DTOs by using pointer optional scalars with `omitempty`.
- Prefer Bun for frontend work and run type checks after TypeScript or TSX edits.
- Do not run Git commit, push, branch operations, or hard resets unless explicitly requested.

## Navigation

- Backend module: `references/modules/backend.md`
- Relay module: `references/modules/relay.md`
- Billing module: `references/modules/billing.md`
- Default frontend: `references/modules/frontend-default.md`
- HTTP routing: `references/api/http-routing.md`
- Persistence and migrations: `references/data/persistence.md`
- Relay request flow: `references/flows/relay-request.md`
- Runtime and local development: `references/runtime/local-development.md`
- Project invariants ADR: `references/adr/project-invariants.md`
- Troubleshooting: `references/troubleshooting/development.md`
- Glossary: `references/glossary.md`

## Verification Hints

- Backend: prefer focused `go test ./path/to/package` before broad `go test ./...`.
- Frontend default: run commands from `web/default/` unless the user explicitly works in `web/classic/`.
- i18n: run `bun run i18n:sync` from `web/default/` after adding user-visible default frontend text.
