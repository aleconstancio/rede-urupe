# Contributing to Maze

Thank you for contributing! Maze is built on the **"Pulse & Gateway"** architecture. We prioritize maintainability and simplicity over complex asynchronous state management.

## Core Principles

1. **Readability and Order > Speed**.
2. **Canonical Doc Policy**: Always update `maze.md`, `cognitive.md`, and `architecture_map.md` when changing system behavior.
3. **Linear Cognition**: Avoid introducing new concurrent loops unless they fit into the Gateway or Capsule Worker patterns.
4. **Explicit over Implicit**: Use descriptive names and document the *why*.

## Development Workflow

1. **Docs First**: Define the architectural change in the canonical docs.
2. **Implement**: Keep components small and focused.
3. **Verify**:
   - Run `go fmt ./...`
   - Run `go build ./cmd/bot`
   - Ensure `.\build.ps1` passes if changing the frontend.

## Testing

- **Unit Tests**: Required for logic in `internal/domain`.
- **Integration**: Verify new Pulse triggers locally before committing.

## Commit Style

We follow [Conventional Commits](https://www.conventionalcommits.org/).

- `feat(gateway): add new stance extraction rule`
- `fix(memory): resolve capsule worker overlap`
- `docs(maze): update cognitive contract`

## Maintenance

- **Purge legacy**: If you see code referencing "Orchestrator", "Episodes", or "Snapshots" as standalone loops, it is a legacy bug. Remove it.
