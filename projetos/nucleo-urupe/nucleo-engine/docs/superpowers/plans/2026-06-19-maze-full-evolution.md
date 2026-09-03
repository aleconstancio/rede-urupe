# Maze Full Evolution — Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Transform Maze from a static dashboard with placeholder shells into a real-time command center with live data, working moderation enforcement, analytics visualization, secure admin endpoints, and multi-channel support.

**Architecture:** Five phases executed sequentially. Phase 1 builds the foundation (API client, SSE, auth middleware). Phase 2 wires all dashboard pages to real data. Phase 3 adds analytics charts. Phase 4 enables actual moderation enforcement via Discord API. Phase 5 refactors the single-channel architecture to support multiple channels. All UI uses bindrunes components exclusively.

**Tech Stack:** Go 1.25+ (net/http middleware), Svelte 5 + SvelteKit, bindrunes (160+ components, DataChart, DataTable, MetricCard, Suspense, createQuery), Supabase JWT, SQLite WAL, discordgo.

---

## File Map

### New Files (21)
| File | Purpose |
|---|---|
| `frontend/src/lib/api/client.ts` | REST API client for Go backend |
| `frontend/src/lib/api/sse.ts` | SSE client with auto-reconnect |
| `frontend/src/lib/api/hooks.ts` | Svelte reactive wrappers around API client |
| `internal/presentation/api/middleware.go` | Auth middleware + rate limiter |
| `internal/domain/moderation/enforcer_session.go` | Discord session wiring for enforcer |
| `frontend/src/routes/dashboard/analytics/+page.svelte` | Analytics dashboard with charts |
| `frontend/src/routes/dashboard/moderation/+page.svelte` | Moderation panel |
| `frontend/src/routes/dashboard/memory/+page.svelte` | Memory capsule browser |
| `frontend/src/routes/dashboard/persona/+page.svelte` | Persona editor |

### Modified Files (18)
| File | Change |
|---|---|
| `internal/presentation/api/server.go` | Add auth middleware, rate limiter, new routes |
| `internal/domain/moderation/enforcer.go` | Wire Discord session, implement delete/timeout |
| `internal/domain/moderation/moderation.go` | Add escalation rules (delete at 3 warnings, timeout at 5) |
| `internal/config/config.go` | Add `ADMIN_API_KEY` or JWT validation config |
| `internal/data/sqlite/repository.go` | Add multi-channel query helpers |
| `internal/domain/gateway/worker.go` | Support multi-channel event processing |
| `internal/presentation/discord/handler.go` | Multi-channel message filtering |
| `frontend/src/routes/dashboard/+layout.svelte` | Add analytics/moderation/memory/persona nav items |
| `frontend/src/routes/dashboard/+page.svelte` | Wire to real API data via client |
| `frontend/src/routes/dashboard/maze/+page.svelte` | Wire to real API data |
| `frontend/src/routes/dashboard/vico/+page.svelte` | Wire to real API data |
| `frontend/src/routes/dashboard/orb/+page.svelte` | Wire to real API data |
| `frontend/src/routes/dashboard/settings/+page.svelte` | Wire to real API data |
| `frontend/src/lib/stores/dashboard.ts` | Refactor to use new API client |
| `docs/architecture_map.md` | Update with new routes and middleware |
| `docs/frontend.md` | Update with new dashboard pages |
| `Makefile` | No change needed |
| `justfile` | No change needed |

---

## Phase 1: Foundation (API Client + SSE + Auth + Rate Limiting)

### Task 1: REST API Client
- Create: `frontend/src/lib/api/client.ts`

### Task 2: SSE Client
- Create: `frontend/src/lib/api/sse.ts`

### Task 3: Reactive Hooks
- Create: `frontend/src/lib/api/hooks.ts`

### Task 4: Auth Middleware + Rate Limiter
- Create: `internal/presentation/api/middleware.go`
- Modify: `internal/presentation/api/server.go`
- Modify: `internal/config/config.go`

---

## Phase 2: Dashboard Wiring

### Task 5: Dashboard Overview Page
- Modify: `frontend/src/routes/dashboard/+page.svelte`

### Task 6: Maze Bot Sub-Page
- Modify: `frontend/src/routes/dashboard/maze/+page.svelte`

### Task 7: Vico Sub-Page
- Modify: `frontend/src/routes/dashboard/vico/+page.svelte`

### Task 8: Orb Sub-Page
- Modify: `frontend/src/routes/dashboard/orb/+page.svelte`

### Task 9: Settings Page
- Modify: `frontend/src/routes/dashboard/settings/+page.svelte`

---

## Phase 3: Analytics Visualization

### Task 10: Analytics Layout + Navigation
- Create: `frontend/src/routes/dashboard/analytics/+layout.svelte`
- Modify: `frontend/src/routes/dashboard/+layout.svelte`

### Task 11: Sentiment Chart
- Create: `frontend/src/routes/dashboard/analytics/+page.svelte`

### Task 12: Token Usage Chart
- Create: `frontend/src/routes/dashboard/analytics/tokens/+page.svelte`

### Task 13: Growth Chart
- Create: `frontend/src/routes/dashboard/analytics/growth/+page.svelte`

### Task 14: Channel Health
- Create: `frontend/src/routes/dashboard/analytics/channels/+page.svelte`

---

## Phase 4: Moderation Enforcement

### Task 15: Wire Enforcer to Discord Session
- Create: `internal/domain/moderation/enforcer_session.go`
- Modify: `internal/domain/moderation/enforcer.go`
- Modify: `internal/presentation/discord/handler.go`

### Task 16: Implement Delete/Timeout Enforcement
- Modify: `internal/domain/moderation/enforcer.go`
- Modify: `internal/domain/moderation/moderation.go`

### Task 17: Moderation Dashboard Page
- Create: `frontend/src/routes/dashboard/moderation/+page.svelte`
- Modify: `frontend/src/routes/dashboard/+layout.svelte`

---

## Phase 5: Multi-Channel Support

### Task 18: Refactor Config for Multi-Channel
- Modify: `internal/config/config.go`
- Modify: `internal/data/sqlite/guild_repo.go`
- Create: `internal/data/sqlite/migrations/000003_active_channels.up.sql`
- Create: `internal/data/sqlite/migrations/000003_active_channels.down.sql`

### Task 19: Multi-Channel Gateway Worker
- Modify: `internal/domain/gateway/worker.go`
- Modify: `internal/presentation/discord/handler.go`

### Task 20: Channel Selector in Dashboard
- Modify: `frontend/src/routes/dashboard/+layout.svelte`
- Modify: `frontend/src/lib/api/client.ts`

---

## Phase 6: Documentation + Final Verification

### Task 21: Documentation Updates
- Modify: `docs/architecture_map.md`
- Modify: `docs/frontend.md`
- Modify: `docs/backend.md`

### Task 22: Final Verification
- Run all checks and build
