# Backend Status

## Project

SCAVO Exchange - Backend

---

## Current Stage

Stage 0 - Foundation

## Current Phase

Phase 0.2 - Core Infrastructure

## Current Subphase

Phase 0.2.3 - Observability and Test Bootstrap

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
- PostgreSQL as durable source of truth
- Redis for ephemeral and coordination state
- environment-driven local infrastructure baseline
- health and readiness must be separated
- observability is a first-class infrastructure concern
- testing must grow with the architecture

---

## What This Subphase Added

This subphase formally defined:

- observability direction
- logging baseline expectations
- health versus readiness distinction
- metrics direction
- testing layer model
- validation baseline before heavy implementation
- diagnostic expectations for future infrastructure phases

---

## What Is Still Not Implemented

Not implemented yet:

- database integration
- migrations
- PostgreSQL wiring
- Redis wiring
- repository scaffolding
- local Docker environment
- migration runner
- health/readiness expanded endpoints
- metrics endpoint
- tracing
- test harness
- chain client
- asset registry
- portfolio aggregation
- indexer
- DEX contracts
- quote engine
- routing engine
- tx tracking
- audit persistence

---

## Recommended Next Step

Phase 0.3.1 - Implementation Bootstrap for Persistence and Health Infrastructure

Recommended scope:

- introduce core DB scaffolding
- introduce core cache scaffolding
- prepare migration structure
- prepare health and readiness baseline implementation
- keep current bootstrap stable while making infrastructure integration possible

---

## Notes for Next Chat

Stage 0 documentation foundation is now materially complete enough to begin conservative implementation-oriented infrastructure work.

The next step should still avoid:

- DEX feature implementation
- chain-heavy module implementation
- wallet linking
- quote logic
- indexer logic

The safest next move is infrastructure bootstrap with health-aware design.