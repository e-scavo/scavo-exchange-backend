# Development Guide

## Objective

This document defines the development rules and structural conventions for the SCAVO Exchange backend.

It exists to keep implementation, documentation, and architectural direction aligned from the beginning of the project.

---

## Current Development Stage

Current project status:

- foundation and documentation phases
- minimal backend bootstrap already exists
- architecture and target structure are formally defined
- persistence and environment direction are formally defined
- observability and testing direction are now formally defined
- implementation of DB, cache, migrations, metrics, and test scaffolding is still pending

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

Repository organization rule:

- repositories should remain module-owned
- shared DB primitives may live under core infrastructure
- repository interfaces should be explicit
- storage mechanics should not leak into transport layer

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

The project should gradually standardize environment configuration around:

- application settings
- auth settings
- PostgreSQL settings
- Redis settings
- RPC settings
- logging settings

---

## Migration Rule

Database schema changes must be versioned.

Rules:

- no undocumented schema drift
- no silent structural DB changes
- migrations should live under `migrations/`
- repository changes and migrations should evolve together
- local development should be able to reproduce schema state deterministically

This rule becomes mandatory once persistence implementation starts.

---

## Cache Rule

Redis is a supporting infrastructure store, not the system of record.

Rules:

- do not place critical durable state only in Redis
- cache usage must remain explicit
- cache invalidation responsibility must be clear
- Redis-backed coordination should remain optional where possible

---

## Logging Rule

Structured logging is the default.

Rules:

- log meaningful operational events
- avoid noisy logs without value
- do not leak secrets
- preserve correlation context where possible
- keep logs consistent across modules
- infrastructure and dependency failures should be clearly distinguishable

Logging must support both local debugging and future operational observability.

---

## Health and Readiness Rule

The project must distinguish between health and readiness.

### Health
Represents whether the backend process is alive and able to answer basic status requests.

### Readiness
Represents whether the backend is actually ready to serve its intended workload, including infrastructure dependencies when required.

This distinction becomes important once the project includes:

- PostgreSQL
- Redis
- chain integrations
- migrations
- background jobs

The backend should not treat "process is up" as equivalent to "system is ready."

---

## Testing Rule

Testing must grow together with the system.

The project should evolve through multiple testing layers:

- unit tests
- service tests
- repository tests
- integration tests
- chain integration tests
- contract/backend end-to-end tests
- internal environment smoke tests

Rules:

- tests should be phase-appropriate
- new infrastructure should become testable as it is introduced
- tests should validate behavior, not only implementation details
- critical flows should gain regression coverage as soon as they stabilize

---

## Error Handling Rule

Errors must remain explicit.

Rules:

- no swallowed errors
- no panic-based business flow
- wrap or classify errors where helpful
- transport layer should map domain and infrastructure errors responsibly
- future error catalog should remain aligned across REST and WebSocket

Infrastructure errors should remain distinguishable from domain errors where possible.

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
- `docs/development.md`
- `docs/development-environment.md`
- `docs/observability.md`
- `docs/testing.md`
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

- `docs(phase-0.2.3): define observability and testing baseline`
- `refactor(core): prepare health and readiness structure`
- `test(auth): add service-level login coverage`

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

## Validation Baseline

Before major product features are added, the backend should gain a minimum validation baseline for:

- startup behavior
- configuration loading
- dependency failure visibility
- endpoint smoke validation
- future repository readiness
- future chain integration readiness

This reduces drift and makes infrastructure work safer.

---

## Safe Expansion Principle

Before adding large product features, the project should stabilize:

- layout
- config
- persistence direction
- migration direction
- environment baseline
- observability baseline
- testing direction
- module boundaries
- local environment rules

This principle is especially important for the current stage.

---

## Current Recommended Next Step

After this phase, the recommended next move is:

Phase 0.3.1 - Implementation Bootstrap for Persistence and Health Infrastructure

This should introduce:

- DB core scaffolding
- cache core scaffolding
- migration workflow bootstrap
- health/readiness baseline implementation
- testable infrastructure wiring direction