# Maze Architecture

> Oficina de Dédalo central hub. Evolving from Discord cognitive bot to unified suite interface.
> System wiring, package layout, data flow.

## Top-Level

```
maze/
├── cmd/
│   ├── bot/                     # Main entry point
│   ├── analyze/
│   ├── deep_analyze/
│   ├── export/
│   └── list-models/
├── internal/
│   ├── config/
│   ├── data/
│   │   ├── llm/                 # Multi-provider client (Gemini/OpenRouter/Groq)
│   │   ├── postgres/
│   │   └── sqlite/              # 25 tables, FTS5, schema V5
│   ├── domain/
│   │   ├── commands/            # Discord slash commands
│   │   ├── forums/              # Template rendering, scheduled publishing
│   │   ├── gateway/             # Pulse & Gateway worker
│   │   ├── memory/              # Capsule worker, compactor
│   │   └── moderation/          # Spam, toxicity, personal attacks
│   ├── pkg/timeutil/            # Brasilia timezone helpers
│   └── presentation/
│       ├── api/                 # REST + SSE server (:8080)
│       ├── discord/             # Gateway event handler
│       └── web/                 # Built Svelte 5 assets
├── frontend/                    # Svelte 5 dashboard
├── apps/
│   └── siren/                   # Marketing suite (Next.js, NestJS, Temporal, Python)
├── docs/
│   ├── maze.md                  # Architectural bible
│   ├── cognitive.md             # Behavioral contract
│   ├── architecture_map.md      # System wiring diagram
│   └── CONTRIBUTING.md          # Contributor guidelines
├── backend.md                   # Backend rules
├── frontend.md                  # Frontend rules
└── AGENTS.md                    # Project laws
```

## Three-Layer Architecture

```
Presentation (Discord/API/Web)
    ↓
Domain (gateway, persona, memory, intelligence, identity, commands, forums, moderation)
    ↓
Data (SQLite, LLM)
```

## Key Components

| Component | Purpose |
|-----------|---------|
| GatewayWorker | Sequential cognitive pipeline, Single Point of Truth |
| Persona | Runtime identity: Core Identity + Style Overlay + Policy |
| Memory | Fractal: Associative (FTS5) + Episodic (today) + Working (raw window) |
| Intelligence | Keyword classification, tone analysis |
| Identity | Deterministic stance tracking |

## Marketing Suite (Siren)

Siren is the marketing engine of the Oficina de Dédalo suite, integrated into the Maze hub.

### Stack
- **Frontend**: Next.js 16, React 19, Tailwind CSS
- **Backend**: NestJS, Prisma ORM, PostgreSQL
- **Video Engine**: Python 3.11, MoviePy, FFmpeg, edge-tts
- **Workflow**: Temporal (TypeScript + Python SDK)
- **Go Services**: Email, webhooks, token refresh

### Features
- **14+ Platform Integrations** — Instagram, YouTube, TikTok, LinkedIn, X, Facebook, and more
- **AI Video Generation** — topic to complete short-form video with script, voice, footage, subtitles
- **Content Scheduling** — calendar view, queue management, drag-and-drop
- **Analytics** — post performance metrics, engagement tracking
- **Self-hosted** — Docker Compose deployment, full data ownership

### Integration with Hub
- Accessible via Maze dashboard
- Uses shared Supabase Auth
- Publishes marketing content for all suite products
- Built for Vico's public launch
