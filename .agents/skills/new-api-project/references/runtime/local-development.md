---
generated: true
source_command: "$project-skill-builder create"
source_paths:
  - "README.md"
  - ".env.example"
  - "docker-compose.yml"
  - "docker-compose.dev.yml"
  - "Dockerfile"
  - "Dockerfile.dev"
  - "web/package.json"
  - "web/default/package.json"
  - "web/classic/package.json"
---

# Runtime And Local Development

Last verified: 2026-06-02
Verified against commit: 7aaa5332

## Main Files

- `.env.example`
- `docker-compose.yml`
- `docker-compose.dev.yml`
- `Dockerfile`
- `Dockerfile.dev`
- `main.go`
- `web/default/package.json`
- `web/classic/package.json`

## Default Runtime

- Backend port defaults to `PORT` or `3000`.
- SQLite is used when `SQL_DSN` is empty.
- `FRONTEND_BASE_URL` redirects web fallback on non-master nodes.
- `ENABLE_PPROF=true` starts pprof on `0.0.0.0:8005`.
- `REDIS_CONN_STRING` enables Redis.
- `SESSION_SECRET` should be configured for production and multi-node deployments.

## Docker

Production-style compose:

```bash
docker compose up -d
```

Development compose:

```bash
docker compose -f docker-compose.dev.yml up -d
```

The development compose builds the backend from local source, exposes backend port `3000`, and starts Redis and PostgreSQL. The frontend dev server is run separately.

## Frontend Development

`web/` is a Bun workspace containing `default` and `classic`. Script definitions live in the theme package directories. For the default theme:

```bash
cd web/default
bun install
bun run dev
```

The default Rsbuild config proxies `/api`, `/mj`, and `/pg` to `http://localhost:3000` unless `VITE_REACT_APP_SERVER_URL` overrides it.

## Backend Development

When running backend source directly, ensure embedded frontend dist directories exist or use the project's Docker/build flow. `main.go` embeds both:

- `web/default/dist`
- `web/classic/dist`

Focused backend tests are usually faster and clearer than running the full suite:

```bash
go test ./controller ./service ./relay/common
```

## Environment Groups

Common `.env.example` categories:

- port and frontend URL
- debug, pprof, and Pyroscope
- SQL and log SQL
- Redis and sync frequency
- task/update settings
- relay and streaming timeouts
- session secret
- redirect domain trust list

Do not record secrets in project-skill references or generated docs.
