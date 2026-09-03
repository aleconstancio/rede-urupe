# Maze: The Cognitive Gateway

> Oficina de Dédalo central hub. Evolving from Discord cognitive bot to unified suite interface.

## 1. Project Philosophy & Scope

Maze is the central hub of the Oficina de Dédalo suite — a collection of agentic AI products for the Brazilian market. Originally a Discord management platform, it is evolving into the unified interface for product switching, user management, and intelligence engine oversight.

The Discord bot remains as one feature among many, powered by the Talos Agent pipeline (7-stage). It provides a high-fidelity LLM-powered agent with dynamic personas. The hub prioritizes maintainability and architectural clarity over hyper-optimized complexity.

- **Core Motto**: "Readability and Order > Speed and Complex Interconnections".
- **Execution Model**: Linear and sequential. No asynchronous state reconstruction.
- **Persistence**: SQLite (WAL enabled) serves as the append-only system of record.
- **Dumb Ingestion**: Message intake is decoupled from cognitive processing. We save first, think later.

---

## 2. Architecture & Systems

### 2.1 The Pulse & Gateway Model

The cognitive engine utilizes a direct "Pulse & Gateway" model for all orchestration:

1.  **Reactive Pulse**: Direct @mentions, replies, or a keyword trigger.
2.  **Volume Pulse**: Every batch of N messages (default 10) triggers an ambient pulse.
3.  **Time Pulse**: A silence period (default 1h) triggers a proactivity check.

### 2.2 The Talos Agent Pipeline

The cognitive engine uses Talos's 7-stage Agent framework for turn orchestration:

1. **Ingest** — Converts Discord events to PipelineInput
2. **Classify** — Content classification via Talos intelligence
3. **Retrieve** — 3-layer memory context assembly
4. **PersonaSelect** — Resolves identity via minotaur.Resolver
5. **Gate** — Social triage (should the bot speak?)
6. **Generate** — LLM reply with structured output
7. **Validate** — Post-generation grounding check

A single, sequential `GatewayWorker` consumes persisted message events. This worker is the **Single Point of Truth** for turn orchestration, preventing concurrent hot-path execution and ensuring linear history.

### 2.3 Social Intelligence

The agent tracks **Member Profiles** to calibrate social interactions:

- **Metadata**: Roles, Age, Interests, Religion, and persistent Notes.
- **Injection**: Profiles of active participants are injected into the `<PARTICIPANTS>` tag of the prompt context.
- **Impact**: Enables context-aware responses that respect user preferences and social standing.

### 2.4 Dynamic Persona Framework (minotaur)

Identity is resolved at runtime via the minotaur module, not static:

- **Core Identity**: The fundamental "Who am I?" (e.g., AIrelius, Talos, Eris).
- **Style Overlay**: Situational masks that modify rhetorical style and tone.
- **Policy**: Channel-specific rules that map intents to specific persona configurations.

### 2.5 Fractal Memory

Memory is treated as an **episode-first compression task**:

- **Associative Memory**: Historical recall via SQLite FTS5 index.
- **Episodic Memory**: Today's active episode records, selected by recency and topical relevance.
- **Working Memory**: Raw message window since the last pulse.

---

## 3. Backend Implementation (Go & SQLite)

### 3.1 Module Map

- `cmd/bot`: Main entry point. Initializes repository, LLM client, and background workers.
- `internal/data/sqlite`: Unified repository layer owning schema, migrations, and social metadata.
- `internal/data/llm`: Modular client with provider pooling, model rotation, and dynamic concurrency.
- `internal/domain/gateway`: Sequential worker and prompt assembly logic.
- `internal/domain/minotaur`: Identity management and adaptive policy resolution (renamed from `persona`).
- `internal/domain/memory`: Fractal Memory workers (CapsuleWorker and Compactor).
- `internal/domain/identity`: Deterministic stance tracking and participant projection.
- `internal/domain/intelligence`: Keyword-based classification, tone analysis, stopword dictionary.
- `frontend/`: Svelte 5 hub — marketing, knowledge map, and dashboard.

### 3.2 Frontend Routes

- `frontend/src/routes/` — SvelteKit route definitions:

| Path | Handler | Auth | Description |
|------|---------|------|-------------|
| `/` | Svelte SPA | No | Marketing landing page |
| `/projects` | Svelte SPA | No | Project showcase (DB-driven) |
| `/knowledge/*` | Quartz static | No | Obsidian knowledge map |
| `/login` | Svelte SPA | No | Supabase auth |
| `/dashboard/*` | Svelte SPA | Yes | Auth-guarded dashboard |
| `/api/*` | Go REST | No* | API endpoints |
| `/events` | Go SSE | No* | Real-time event stream |

### 3.2 Persistence Model

- **`messages`**: Canonical append-only log of all conversational facts.
- **`member_profiles`**: Persistent social metadata for room participants.
- **`memory_capsules`**: Episode records storing source row spans and structured takeaways.
- **`budget_events`**: Immutable ledger of token usage and costs in Brasilia Time (UTC-3).

---

## 4. Cognitive Architecture (LLM & Memory)

### 4.1 The Socialized Triple-Context

Context is partitioned using the **Tagged Anchor Strategy** (XML tags):

- `<PARTICIPANTS>`: Social intelligence metadata for active speakers.
- `<ASSOCIATIVE_MEMORY>`: FTS5 past capsules for historical analogy.
- `<EPISODIC_MEMORY>`: Today's selected episode records.
- `<WORKING_MEMORY>`: Raw dialogue since last pulse.

### 4.2 Modular LLM Pipeline

- **Provider Pooling**: Parallel management of Gemini, Groq, and OpenRouter backends.
- **Model Rotation**: Automatic fallback sequence upon spending caps or 503 errors.
- **Dynamic Concurrency**: Semaphore-based throttling adjustable via dashboard config.

---

## 5. Frontend Dashboard (Svelte 5)

- **Tech**: Svelte 5 (Runes), Vite, Bun.
- **Design**: "Midnight Onyx" high-contrast glassmorphism.
- **Styling**: Vanilla CSS. **TailwindCSS is Forbidden.**

---

## 6. Design Constraints & Invariants

- **Temporal Integrity**: All timestamps and logic are anchored to **Brasilia Time (UTC-3)**.
- **Bounded Monologue**: Every turn requires a social diagnostic worksheet before generating text.
- **Memory Anchors**: Claims must be justified via a structured `grounding_ledger`.
- **FORBIDDEN**: Manual mutexes for concurrency (use the `GatewayWorker`).
- **FORBIDDEN**: Vector Databases (use SQLite FTS5).
- **FORBIDDEN**: TailwindCSS.
