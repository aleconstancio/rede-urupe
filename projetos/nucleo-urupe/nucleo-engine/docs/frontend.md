# Maze Frontend

> Svelte 5 hub: marketing, knowledge map, and real-time dashboard.

## Stack
- Svelte 5 with runes (SvelteKit, adapter-static)
- Tailwind CSS v4
- Bindrunes component library (160+ components)
- Supabase auth

## Routes
| Route | Auth | Purpose |
|-------|------|---------|
| `/` | No | Marketing landing page |
| `/projects` | No | Project showcase (API-driven) |
| `/knowledge` | No | Link to Quartz knowledge map |
| `/login` | No | Supabase email/password auth |
| `/dashboard` | Yes | Overview with live metrics (SSE) |
| `/dashboard/maze` | Yes | Maze bot status, persona, feed |
| `/dashboard/vico` | Yes | Vico service status |
| `/dashboard/orb` | Yes | Orb service status |
| `/dashboard/analytics` | Yes | Sentiment, tokens, growth, channels |
| `/dashboard/moderation` | Yes | Mod log and audit trail |
| `/dashboard/settings` | Yes | Account, theme, logout |

## API Integration
- `lib/api/client.ts` — REST client for Go backend
- `lib/api/sse.ts` — SSE client with auto-reconnect (events: feed, metrics, persona, projects)
- `lib/api/hooks.ts` — `createLiveQuery()` reactive hook combining REST + SSE
- `lib/stores/dashboard.ts` — Legacy store (still used for product status cards)

## State
- `createLiveQuery(fetcher, refreshOn)` — primary data fetching pattern
- `localStorage` — channel selector persistence
- Supabase auth via `bindrunes/createAuth`

## Rules
- Svelte 5 runes only — no legacy stores for new code
- All data fetching via `createLiveQuery` with SSE refresh
- Public routes (/, /projects, /knowledge) have no auth requirement
- Dashboard is real-time, not static
