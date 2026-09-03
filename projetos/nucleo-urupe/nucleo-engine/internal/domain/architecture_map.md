# Maze Architecture Map

This document maps the core domain logic located in `internal/domain`. Maze follows the **Pulse & Gateway** model, where all cognitive turns are sequential and linear.

## Domain Overview

| Domain | Responsibility | Role in Pulse & Gateway |
| :--- | :--- | :--- |
| **`gateway`** | **The Core**. Sequential worker and turn orchestration. | Consumes events and executes structured turns. |
| **`persona`** | **Identity Resolution**. Identities and Style Overlays. | Decides "who" the agent is for the current turn. |
| **`memory`** | **Fractal Memory**. Episodic and Semantic compression. | Provides Associative and Episodic context. |
| **`intelligence`** | **Cognitive Logic**. Classification and diagnostic logic. | Handles the "internal monologue" and reasoning. |
| **`identity`** | **Social Intelligence**. Participant stance tracking and projections. | Injects participant metadata and stance history into prompts. |

## Interaction Flow (Sequential)

1. **Ingestion**: Discord Event -> Persisted in `messages` -> Enqueued to `GatewayWorker`.
2. **Pulse Check**: `GatewayWorker` checks if turn is warranted (Reactive vs. Ambient Gate).
3. **Resolution**: `persona.Resolver` selects active Identity and Style Overlay.
4. **Assembly**: `PayloadAssembler` builds **Socialized Triple-Context** (Participants + Memory).
5. **Execution**: `gateway.ExecuteStructured` calls the LLM via modular pipeline.
6. **Persistence**: Turn details, cost, and monologue saved to `messages` and `budget_events`.

## Design Principles

- **Sequentiality**: Only one `GatewayWorker` per channel. No concurrent "thinking".
- **Dumb Ingestion**: Discord events are saved immediately; cognitive processing happens later.
- **Social First**: Every turn is calibrated against participant metadata (stance, interests, etc.).
- **Modular Intelligence**: LLM providers are pooled and rotated automatically on failure.
