---
generated: true
source_command: "$project-skill-builder create"
source_paths:
  - "web/default/AGENTS.md"
  - "web/package.json"
  - "web/default/package.json"
  - "web/default/rsbuild.config.ts"
  - "web/default/src/main.tsx"
  - "web/default/src/routes/__root.tsx"
  - "web/default/src/i18n/"
---

# Default Frontend Module

Last verified: 2026-06-02
Verified against commit: 7aaa5332

## Main Files

- `web/default/package.json`
- `web/default/rsbuild.config.ts`
- `web/default/src/main.tsx`
- `web/default/src/routes/`
- `web/default/src/features/`
- `web/default/src/components/`
- `web/default/src/lib/api.ts`
- `web/default/src/i18n/`
- `web/default/AGENTS.md`

## Responsibility

`web/default` is the primary React frontend. It handles the dashboard, setup, auth, channels, keys, models, pricing, usage logs, subscriptions, system settings, performance metrics, rankings, playground, and other management UI.

## Stack

- React 19 and TypeScript
- Rsbuild 2
- TanStack Router and TanStack Query
- Zustand
- Base UI and Tailwind CSS
- i18next and react-i18next
- VChart/Recharts for charts
- Bun for package management and scripts

## Entrypoints

- `src/main.tsx` initializes frontend cache, build metadata, React Query, TanStack Router, theme/font/direction providers, i18n, branding from `/api/status`, and global error handling.
- `src/routes/__root.tsx` is the root route. It loads system config, stores affiliate codes, guards setup redirect, mounts common UI, and wires route-level errors.
- `rsbuild.config.ts` proxies `/api`, `/mj`, and `/pg` to the backend dev server, and defines route-based code splitting behavior.

## Scripts

Run from `web/default/` unless a workspace-level command is explicitly intended:

```bash
bun run dev
bun run build
bun run build:check
bun run typecheck
bun run lint
bun run i18n:sync
```

## Invariants

- User-visible text must use i18n via `useTranslation()` in React components.
- New i18n keys must be present for `en`, `zh`, `fr`, `ja`, `ru`, and `vi`; use `$i18n-translate` for translation work.
- TypeScript/TSX edits require a typecheck before completion when feasible.
- Avoid `any`; prefer explicit types or `unknown`.
- Feature code belongs under `src/features/<feature>/`; shared UI and utilities belong under `src/components/` and `src/lib/`.
- Use the project `api` instance and `handleServerError` patterns rather than ad hoc request handling.

## Classic Theme

`web/classic` is a separate theme. Do not port or compare classic changes manually when the user asks for parity; use `$classic-to-default-sync`.

## Tests

Existing default frontend tests include component tests such as `web/default/src/components/ui/dropdown-menu.test.tsx`. The package scripts expose typecheck, lint, and build checks.
