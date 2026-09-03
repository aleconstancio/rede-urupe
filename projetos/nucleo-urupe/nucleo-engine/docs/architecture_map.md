# Maze Architecture Map

> Oficina de Dédalo central hub. Evolving from Discord cognitive bot to unified suite interface.

## System Overview

Maze is the central hub of the Oficina de Dédalo suite — a collection of agentic AI products for the Brazilian market. It provides a unified interface for product switching, user management, and intelligence engine oversight.

The Discord bot remains as one feature, powered by the Talos Agent pipeline (7-stage). The hub is built on Go 1.26 + Svelte 5 + Bindrunes.

```
┌──────────────────────────────────────────────────────────────────┐
│                        Discord Gateway                           │
└──────────┬───────────────────────────────────────┬────────────────┘
           │ events                               │ API calls
           ▼                                       ▼
┌──────────────────────┐          ┌──────────────────────────────┐
│  internal/presentation│          │  internal/presentation       │
│  /discord/handler.go  │          │  /api/server.go              │
│  OnMessageCreate      │          │  REST + SSE on :8080         │
│  OnInteractionCreate  │          │  /api/persona, /api/forum    │
│  OnReactionAdd        │          │  /api/memory, /events        │
└──────────┬────────────┘          └──────────────┬───────────────┘
           │ triggerChan                           │ serves
           ▼                                       ▼
┌──────────────────────────────────────────────────────────────────┐
│                    internal/domain                                │
│                                                                  │
│  ┌──────────┐  ┌──────────┐  ┌────────┐  ┌───────────────────┐  │
│  │ gateway  │  │minotaur  │  │memory  │  │    intelligence    │  │
│  │ Pulse &  │  │Identity  │  │Capsule │  │   classification   │  │
│  │ Gateway  │  │Overlays  │  │+Com-   │  │   tone analysis    │  │
│  │ Worker   │  │Policy    │  │pactor  │  │                    │  │
│  └──────────┘  └──────────┘  └────────┘  └───────────────────┘  │
│                                                                  │
│  ┌──────────┐  ┌──────────┐  ┌──────────────────────────────┐   │
│  │ identity  │  │ commands │  │        forums                │   │
│  │ Stance    │  │ /maze    │  │  Template rendering +        │   │
│  │ Ledger    │  │ /maze-   │  │  scheduled publishing        │   │
│  │           │  │ config   │  │                              │   │
│  └──────────┘  └──────────┘  └──────────────────────────────┘   │
│                                                                  │
│  ┌──────────────────────────────────────────────────────────┐    │
│  │                    moderation                             │    │
│  │              spam / toxicity / attacks                    │    │
│  └──────────────────────────────────────────────────────────┘    │
└──────────────────────────────┬───────────────────────────────────┘
                               │
                               ▼
┌──────────────────────────────────────────────────────────────────┐
│                    internal/data                                  │
│                                                                  │
│  ┌──────────────────────────────────────────────────────────┐    │
│  │  sqlite/                                                  │    │
│  │  Repository — 25 tables (see below)                      │    │
│  │  Schema V5, FTS5 full-text search                        │    │
│  └──────────────────────────────────────────────────────────┘    │
│                                                                  │
│  ┌──────────────────────────────────────────────────────────┐    │
│  │  llm/                                                     │    │
│  │  Client — Multi-provider (Gemini, Groq, OpenRouter)      │    │
│  │  Provider pooling, model rotation, concurrency semaphore │    │
│  └──────────────────────────────────────────────────────────┘    │
└──────────────────────────────────────────────────────────────────┘
                               │
                               ▼
┌──────────────────────────────────────────────────────────────────┐
│                    frontend/ (Svelte 5 + Vite)                    │
│  App.svelte → views (Dashboard, Memory, Persona)                │
│  → built to internal/presentation/web/static_dist               │
│  ├── / (landing)              # Marketing page                  │
│  ├── /projects                # DB-driven showcase              │
│  ├── /knowledge               # Quartz knowledge map            │
│  ├── /dashboard/*             # Auth-guarded dashboard          │
│  └── /login                   # Supabase auth                   │
│  └── knowledge_dist/          # Built Quartz output              │
└──────────────────────────────────────────────────────────────────┘
```

## Layer Map

### Presentation Layer

| Package | Responsibility |
|---------|---------------|
| `internal/presentation/discord` | Discord gateway event handler — message create/update/delete, reactions, interactions |
| `internal/presentation/api` | REST API server (status, metrics, persona, memory, forum, SSE events) |
| `internal/presentation/web/static_dist` | Built Svelte 5 frontend assets |

### Domain Layer

| Package | Responsibility |
|---------|---------------|
| `internal/domain/gateway` | Pulse & Gateway cognitive engine — turn orchestration, structured LLM prompts, context assembly, Discord messenger |
| `internal/domain/minotaur` | Persona system — core identities (AIrelius, Talos, Eris), style overlays, adaptive memory, policy resolution (renamed from `persona`) |
| `internal/domain/intelligence` | Content classification (25 Portuguese categories, keyword-based), tone analysis, stopword filtering |
| `internal/domain/identity` | Event-sourced participant stance ledger — claims, arguments, concessions per topic per user |
| `internal/domain/memory` | Fractal memory — capsule segmenter, capsule worker, compactor, FTS5 search |
| `internal/domain/commands` | Discord slash command router — `/maze status`, `/maze config`, etc. |
| `internal/domain/forums` | Forum automation — template rendering with variable substitution, scheduled publishing |
| `internal/domain/moderation` | Rule-based moderation — spam detection, toxicity scoring, personal attack detection |

### Data Layer

| Package | Responsibility |
|---------|---------------|
| `internal/data/sqlite` | SQLite repository — all schema definitions and CRUD operations across 25 tables |
| `internal/data/llm` | Multi-provider LLM client — Gemini (native), OpenRouter, Groq; structured outputs, embeddings, provider rotation |

### Infrastructure

| Package | Responsibility |
|---------|---------------|
| `internal/config` | Environment variable loading into typed Config struct |
| `internal/pkg/timeutil` | Brasilia Time (UTC-3) normalization for all timestamps |

## CLI Tools

| Command | Description |
|---------|-------------|
| `cmd/bot` | Main Maze bot binary — wires all subsystems, opens Discord session |
| `cmd/export` | Message exporter — outputs JSONL fine-tuning pairs grouped by channel/day |
| `cmd/analyze` | Basic DB analysis — volume stats, active persona, top users/categories |
| `cmd/deep_analyze` | Deep intelligence report — hourly heatmap, participant rankings, category distribution, bigrams, token costs, interaction network |
| `cmd/list-models` | OpenRouter model lister — prints available model IDs |

## Database Tables (Schema V5)

### Core Conversation

| Table | Purpose |
|-------|---------|
| `messages` | Append-only message log (discord_id unique, channel, author, content, category, reactions, internal_monologue, grounding_ledger) |
| `message_artifacts` | Per-message extracted artifacts (type, content, metadata) |
| `message_annotations` | Enriched metacognitive metadata (topic/stance/style tags, durability/retrieval/humor/sarcasm scores) |
| `pending_turns` | Queued turn processing requests (retry count, status, error tracking) |
| `gateway_state` | Pulse/gateway watermark per channel (last_pulse_row_id, last_pulse_msg_id) |

### Memory

| Table | Purpose |
|-------|---------|
| `memory_capsules` | Episodic memory snapshots (day_date, time_span, kind, participants, main_topic, mood, key_facts, open_loops) |
| `memory_search` | FTS5 virtual table over memory_capsules (participants, main_topic, key_facts, category) with auto-sync triggers |
| `capsule_hwm` | Capsule processing high-water mark per channel |

### Persona System (minotaur)

| Table | Purpose |
|-------|---------|
| `core_identity_profiles` | Identity family definitions (name, identity_prompt, core_values_json, is_default) |
| `persona_overlays` | Style/posture modules (identity_id FK, style_prompt, traits_json, allowed_intents_json, sort_order) |
| `adaptive_persona_memory` | Per-channel adaptive style adjustments (adaptive_style_json, confidence, source, expires_at) |
| `persona_delta_proposals` | Metacognitive change proposals (target, proposed_changes_json, reason, evidence_message_ids, confidence, status) |
| `persona_policy` | Channel-level resolution rules (default_identity, selection_mode, allowed lists, intent→persona map, manual override) |

### Identity / Social

| Table | Purpose |
|-------|---------|
| `stance_events` | Append-only participant stance ledger (topic, position, action, evidence_type, confidence, source_msg_id) |
| `behavior_events` | Participant behavior archetype records (archetype, rhetorical_style, status) |
| `member_profiles` | Enriched participant profiles (discord_id, roles, age, interests, religion, notes) |
| `agent_profiles` | Legacy agent profile storage |

### Forum Automation

| Table | Purpose |
|-------|---------|
| `forum_templates` | Forum post templates with scheduling (guild_id, channel_id, title, body, tags, variables, schedule, is_enabled) |
| `forum_posts` | Published forum posts (template_id FK, guild_id, discord_message_id, discord_thread_id, status, error) |

### Guild Management

| Table | Purpose |
|-------|---------|
| `guild_config` | Per-guild bot configuration (prefix, default_channel, mod_log_channel, gate/reply/memory models, mod_enabled, max_warnings) |
| `guild_warnings` | Moderation warning records (user_id, channel_id, reason, severity) |

### Operational

| Table | Purpose |
|-------|---------|
| `budget_events` | Token cost ledger (input_tokens, output_tokens, model, reason, trigger_type) |
| `system_metrics` | Operational metrics log (metric_type, sub_type, value, latency_ms) |
| `system_config` | Key-value configuration store |
| `sync_state` | Discord history sync watermark (latest_synced_msg_id, oldest_synced_msg_id) |

### Showcase

| Table | Purpose |
|-------|---------|
| `showcase_projects` | Project listings for the marketing showcase (name, slug, description, category, GitHub URL, featured flag) |

### Multi-Channel
| Table | Purpose |
|-------|---------|
| `active_channels` | Registered channels per guild for multi-channel mode |

## Key Design Constraints

- **Temporal Integrity**: All timestamps in Brasilia Time (UTC-3) via `internal/pkg/timeutil`
- **Single Point of Truth**: `GatewayWorker` is the sole orchestrator — no concurrent hot-path execution
- **Append-Only**: Messages, stance events, and budget events are never mutated
- **No Vector DB**: SQLite FTS5 for full-text search
- **No Manual Mutexes**: Concurrency handled via the sequential worker model
- **No TailwindCSS**: Vanilla CSS with "Midnight Onyx" glassmorphism theme

## Route Map (Hub Architecture)

| Path | Handler | Auth | Description |
|------|---------|------|-------------|
| `/` | Svelte SPA | No | Marketing landing page |
| `/projects` | Svelte SPA | No | Project showcase (DB-driven) |
| `/knowledge/*` | Quartz static | No | Obsidian knowledge map |
| `/login` | Svelte SPA | No | Supabase auth |
| `/dashboard/*` | Svelte SPA | Yes | Auth-guarded dashboard |
| `/api/*` | Go REST | No* | API endpoints |
| `/api/channels` | Go REST | No | List active channels |
| `/api/projects` | Go REST | No | Showcase projects CRUD |
| `/api/projects/` | Go REST | No | Project by slug |
| `/events` | Go SSE | No* | Real-time event stream |

## Middleware

| Middleware | Applied To | Purpose |
|------------|-----------|---------|
| `requireAuth` | `/api/admin/*` | Bearer token validation (ADMIN_API_KEY env var) |
| `RateLimiter` | `/events` | 60 req/min per IP on SSE endpoint |

## Environment Variables

See [.env.example](../.env.example) for the full list of 30+ configuration variables. Key groups:
- **Required**: `DISCORD_TOKEN`, `OPENROUTER_API_KEY` or `GEMINI_API_KEY`, `DISCORD_CHANNEL_ID`
- **Model Assignment**: `GATEWAY_GATE_MODEL`, `GATEWAY_REPLY_MODEL`, `MEMORY_MODEL`
- **Dashboard**: `DASHBOARD_ENABLED`, `DASHBOARD_PORT`
- **Database**: `SQLITE_PATH` (default: `data/memory.db`)
