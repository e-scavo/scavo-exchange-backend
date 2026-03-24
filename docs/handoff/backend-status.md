# Backend Status

## Project

SCAVO Exchange - Backend

---

## Current Stage

Stage 0 - Foundation

## Current Phase

Phase 0.2 - Core Infrastructure

## Current Subphase

Phase 0.2.1 - Core Infrastructure Layout and Foundation

---

## Real Current Code Baseline

The current backend already includes:

- application bootstrap
- config loading
- structured logger
- HTTP router
- middleware
- WebSocket handler, hub, client, dispatcher, protocol
- JWT token service
- development HTTP login
- system WebSocket handler
- auth WebSocket handler registration

Current real modules:

- `system`
- `auth`

Current real core packages:

- `config`
- `logger`
- `httpx`
- `auth`
- `ws`

---

## Decisions Locked

The following decisions are officially locked at this point:

- modular monolith architecture
- DEX-first product strategy
- AMM v1 as the initial DEX model
- SCAVIUM as the primary chain
- self-custody first for initial DEX scope
- PostgreSQL as primary database target
- Redis as secondary infrastructure store
- REST and WebSocket as first-class transports
- matching engine out of initial scope
- migrations as schema source of truth
- platform adapters for chain and contract integrations
- infrastructure baseline before major feature expansion

---

## What This Subphase Added

This subphase formally defined:

- target technical repository layout
- growth direction for `internal/core`
- growth direction for `internal/modules`
- future `internal/platform` adapter role
- migration-based persistence direction
- observability direction
- jobs/background processing direction
- explicit development rules for safe expansion

---

## What Is Still Not Implemented

Not implemented yet:

- database integration
- migrations
- PostgreSQL wiring
- Redis wiring
- repository scaffolding
- local Docker environment
- chain client
- asset registry
- portfolio aggregation
- indexer
- DEX contracts
- quote engine
- routing engine
- tx tracking
- audit persistence
- metrics and readiness support
- test harness

---

## Recommended Next Step

Phase 0.2.2 - Persistence and Environment Baseline

Recommended scope:

- introduce DB and cache scaffolding direction into the project
- define migration workflow
- define initial local infrastructure workflow
- prepare repository-ready technical base
- keep current bootstrap stable while introducing infrastructure support

---

## Notes for Next Chat

The next step should remain infrastructure-oriented.

The project should still avoid jumping into:

- chain-heavy implementation
- DEX contracts
- quote logic
- wallet linking
- indexer logic

Before those, the backend should gain a stable baseline for:

- database
- migrations
- cache
- repository structure
- local environment