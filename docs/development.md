# Development Guide

## Objective

This document defines the development rules and structural conventions for the SCAVO Exchange backend.

It exists to keep implementation, documentation, and architectural direction aligned from the beginning of the project.

---

## Current Development Stage

Current project status:

- foundation and documentation phases
- minimal backend bootstrap already exists
- architecture and target structure are now formally defined
- persistence and infrastructure implementation are still pending

The current goal is to expand the backend safely without destabilizing the existing bootstrap.

---

## Source of Truth

The source of truth for the project is always:

- the latest project archive provided for review
- the current repository state after approved changes
- the aligned documentation under `README.md` and `docs/`

No previous assumption should override the real attached project state.

---

## Development Model

The project follows a phase-driven model.

Work must be organized as:

- stages
- phases
- subphases

Each subphase should aim to be:

- safe
- reviewable
- testable
- documented
- commit-friendly

Large unstructured jumps should be avoided.

---

## Architecture Rule

The backend architecture is a modular monolith.

This means:

- one main backend application
- explicit module boundaries
- shared technical infrastructure
- no premature microservice split
- clear room for future extraction if needed

All implementation decisions must respect this.

---

## Code Layout Direction

The official growth direction of the repository is:

- `cmd/` for entrypoints
- `internal/app` for composition and application wiring
- `internal/core/` for reusable technical infrastructure
- `internal/modules/` for domain modules
- `internal/platform/` for chain and external integration adapters
- `migrations/` for DB schema evolution
- `deployments/` for environment and deployment assets
- `scripts/` for local development support
- `docs/` for aligned project documentation

Not every directory must exist immediately, but new code should follow this direction.

---

## Handler Rule

Handlers must stay thin.

Handlers may:

- decode input
- validate request format
- invoke service layer
- encode output
- propagate request context

Handlers must not:

- implement business workflows
- coordinate complex persistence
- embed blockchain protocol logic
- accumulate unrelated module responsibilities

---

## Service Rule

Service layer owns business orchestration.

Services may:

- validate business conditions
- coordinate repositories
- coordinate external integrations
- assemble domain outputs
- emit audit-oriented events when appropriate

Services should remain readable and explicit.

---

## Repository Rule

Repositories are responsible for persistence access only.

Repositories may:

- create, read, update, delete domain data
- handle transaction boundaries where appropriate
- isolate SQL and storage-specific details

Repositories must not:

- know about HTTP or WebSocket payloads
- format API responses
- implement business policy that belongs in services

---

## Integration Rule

External systems must be integrated through dedicated clients or adapters.

Examples:

- SCAVIUM RPC
- contract calls
- Redis
- PostgreSQL driver wiring
- future compliance providers

Low-level external protocol code should not be duplicated across modules.

---

## Configuration Rule

Configuration must remain explicit and environment-driven.

Rules:

- no hidden magic constants for environment-specific behavior
- secrets must come from environment or secret managers later
- defaults may exist for local development only when clearly documented
- config changes must remain discoverable and auditable

---

## Logging Rule

Structured logging is the default.

Rules:

- log meaningful operational events
- avoid noisy logs without value
- do not leak secrets
- preserve correlation context where possible
- keep logs consistent across modules

Logging must support both local debugging and future operational observability.

---

## Error Handling Rule

Errors must remain explicit.

Rules:

- no swallowed errors
- no panic-based business flow
- wrap or classify errors where helpful
- transport layer should map domain and infrastructure errors responsibly
- future error catalog should remain aligned across REST and WebSocket

---

## Documentation Rule

Documentation must always evolve with the code.

Required rule:

- no code change that silently invalidates docs
- no documentation that contradicts implementation
- all new structural decisions must be reflected in docs when relevant

At minimum, the following should stay aligned:

- `README.md`
- `docs/index.md`
- `docs/roadmap.md`
- `docs/architecture.md`
- `docs/architecture-deep.md`
- `docs/flows.md`
- `docs/decisions.md`
- `docs/handoff/backend-status.md`
- `docs/phase-status.md`

---

## Commit Rule

Each subphase should include a suggested commit message.

Commit messages should be:

- scoped
- phase-aware when useful
- readable
- consistent with the type of work performed

Examples:

- `docs(phase-0.2.1): define infrastructure layout and repository direction`
- `refactor(core): prepare application structure for persistence bootstrap`
- `feat(auth): add wallet challenge generation endpoint`

---

## Local Development Direction

The backend should eventually support a consistent local environment including:

- backend app
- PostgreSQL
- Redis
- migration workflow
- seeded development baseline when needed

This will be introduced incrementally and documented as the project advances.

---

## Testing Direction

Testing will also evolve incrementally.

Planned test layers include:

- unit tests
- service tests
- repository tests
- chain integration tests
- contract/backend integration tests
- internal environment smoke tests

Testing support should grow together with the infrastructure baseline.

---

## Safe Expansion Principle

Before adding large product features, the project should stabilize:

- layout
- config
- persistence direction
- module boundaries
- local environment rules

This principle is especially important for the current stage.

---

## Current Recommended Next Step

After this structural phase, the recommended next move is:

Phase 0.2.2 - Persistence and Environment Baseline

This should introduce:

- database wiring direction
- migration strategy
- initial local infrastructure setup
- repository-ready technical scaffolding