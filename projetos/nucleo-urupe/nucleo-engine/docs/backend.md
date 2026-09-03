# Maze Backend

> Pulse & Gateway model, cognitive pipeline, SQLite persistence.

## Core Philosophy
"Readability and Order > Speed and Complex Interconnections."
Linear and sequential. No asynchronous state reconstruction.

## Pulse & Gateway Model
1. **Reactive Pulse**: Direct @mentions, replies, keyword triggers
2. **Volume Pulse**: Every N messages (default 10) triggers ambient pulse
3. **Time Pulse**: Silence period (default 1h) triggers proactivity check

Single `GatewayWorker` per channel — Single Point of Truth for turn orchestration.

## Talos Agent Integration
The cognitive pipeline now uses Talos's 7-stage Agent framework:
1. **Ingest** — Converts Discord events to PipelineInput
2. **Classify** — Content classification via Talos intelligence
3. **Retrieve** — 3-layer memory context assembly
4. **PersonaSelect** — Resolves identity via minotaur.Resolver
5. **Gate** — Social triage (should the bot speak?)
6. **Generate** — LLM reply with structured output
7. **Validate** — Post-generation grounding check

Personas are defined in `talos/persona/` (renamed from `persona`).

## Social Intelligence
Member Profiles: roles, age, interests, religion, notes.
Injected into `<PARTICIPANTS>` tag of prompt context.

## Dynamic Persona
- Core Identity: fundamental "who am I?"
- Style Overlay: situational masks
- Policy: channel-specific rules mapping intents to personas

## Fractal Memory
- **Associative**: Historical recall via SQLite FTS5
- **Episodic**: Today's active records, recency + topical relevance
- **Working**: Raw message window since last pulse

## Moderation Enforcement
- Detection: spam (duplicates + flood), toxicity (tone analysis), personal attacks (keyword match)
- Enforcement: warn → delete (3+ warnings) → timeout 15min (5+ warnings) → timeout 1h (8+ warnings)
- Manual actions: kick, ban via Discord API
- Session wired via `DiscordSession` adapter

## Multi-Channel Mode
- Set `DISCORD_CHANNEL_ID=""` to enable multi-channel mode
- Channels auto-register on first message
- Dashboard channel selector persists choice in localStorage
- Per-channel persona, analytics, and moderation

## Security
- Admin endpoints protected by `ADMIN_API_KEY` (Bearer token)
- SSE endpoint rate limited (60 req/min per IP)
- Auth bypass when `ADMIN_API_KEY` is empty (open mode)

## Persistence
SQLite (WAL enabled), append-only system of record.
Key tables: `messages`, `member_profiles`, `memory_capsules`, `budget_events`.

## Canonical Docs
Always update when behavior changes:
- `docs/maze.md` — architectural bible
- `docs/cognitive.md` — behavioral contract
- `docs/architecture_map.md` — system wiring
