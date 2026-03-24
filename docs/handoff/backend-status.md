# Backend Status

## Project

SCAVO Exchange - Backend

---

## Current Stage

Stage 0 - Foundation

## Current Phase

Phase 0.2 - Core Infrastructure

## Current Subphase

Phase 0.2.2 - Persistence and Environment Baseline

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

---

## What This Subphase Added

This subphase formally defined:

- persistence role separation between PostgreSQL and Redis
- migration workflow direction
- local environment baseline
- base environment variable direction
- repository preparation rules
- Docker-oriented local infrastructure recommendation
- persistence boundary rules for future implementation

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

Phase 0.2.3 - Observability and Test Bootstrap

Recommended scope:

- define observability baseline in more detail
- define health, readiness, and operational visibility direction
- define testing structure and initial harness direction
- prepare the project for safe infrastructure implementation after the documentation baseline is complete

---

## Notes for Next Chat

The project should still remain in foundation mode.

The next step should continue preparing the backend before entering implementation-heavy stages.

The current order remains:

- define foundation
- define persistence and environment
- define observability and testing baseline
- then move into implementation-oriented infrastructure work