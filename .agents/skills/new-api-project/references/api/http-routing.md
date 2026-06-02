---
generated: true
source_command: "$project-skill-builder create"
source_paths:
  - "router/main.go"
  - "router/api-router.go"
  - "router/relay-router.go"
  - "router/web-router.go"
  - "controller/"
  - "middleware/"
---

# HTTP Routing

Last verified: 2026-06-02
Verified against commit: 7aaa5332

## Main Files

- `router/main.go`
- `router/api-router.go`
- `router/relay-router.go`
- `router/web-router.go`
- `router/dashboard-router.go`
- `router/video-router.go`

## Router Composition

`router.SetRouter()` registers API, dashboard, relay, video, and web routes. If `FRONTEND_BASE_URL` is set on a non-master node, unmatched routes redirect to that external frontend. Otherwise the backend serves embedded theme assets.

## Route Families

- `/api`: dashboard/admin/user/status/setup/payment/config APIs.
- `/v1`, `/v1beta`, `/pg`, `/mj`, `/suno`: relay APIs.
- Web fallback: embedded `web/default/dist` or `web/classic/dist`, selected by theme.

## Auth Boundaries

- Public status/setup/legal/content routes live under `/api`.
- User routes use `middleware.UserAuth()`.
- Admin routes use `middleware.AdminAuth()`.
- Root-only routes use `middleware.RootAuth()`.
- Relay routes generally use `middleware.TokenAuth()`, except playground uses `UserAuth()`.
- Critical user/payment/security routes combine auth with rate-limit, Turnstile, secure verification, or disable-cache middleware as needed.

## Middleware Patterns

- `/api` uses route tagging, gzip, body storage cleanup, and global API rate limiting.
- Relay routing applies CORS, decompression, body storage cleanup, and stats globally, then applies token/user auth, performance checks, model rate limits, and `Distribute()` by route group.
- Web routing applies gzip, web rate limits, cache middleware, and static serving.

## When Adding Routes

- Put route wiring in the right router file rather than controller files.
- Apply the minimum required auth middleware explicitly.
- Keep wildcard routes after specific routes.
- For relay routes, choose the correct `types.RelayFormat` and ensure `middleware.Distribute()` can extract or default the model.
- Avoid duplicating a long route table in documentation; use `router/*.go` as source of truth.
