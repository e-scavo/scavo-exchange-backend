# Backend Status

## Project

SCAVO Exchange - Backend

---

## Current Stage

Stage 0 - Foundation

## Current Phase

Phase 0.1 - Baseline and Documentation

## Current Subphase

Phase 0.1.2 - Architecture Definition

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

The following decisions are now considered officially locked for the current project direction:

- modular monolith architecture
- DEX-first product strategy
- AMM v1 as the initial DEX model
- SCAVIUM as the primary chain
- self-custody first for initial DEX scope
- PostgreSQL as primary database target
- Redis as secondary infrastructure store
- REST and WebSocket as first-class transports
- matching engine out of initial scope

---

## What This Subphase Added

This subphase formally defined:

- architecture style
- internal layer boundaries
- present and future domain modules
- backend role relative to smart contracts
- hybrid expansion compatibility
- request flow direction
- DEX execution boundary
- persistence direction

---

## What Is Still Not Implemented

Not implemented yet:

- database integration
- migrations
- repositories
- users domain
- wallet linking
- chain client
- asset registry
- portfolio aggregation
- indexer
- DEX contracts
- quote engine
- routing engine
- liquidity flows
- tx tracking
- audit persistence
- Redis integration
- observability stack
- test harness

---

## Recommended Next Step

Phase 0.2.1 - Core Infrastructure Layout and Foundation

Recommended scope:

- formalize target project structure
- introduce database and persistence direction scaffolding
- prepare shared application wiring for future modules
- keep the current bootstrap stable while expanding safely

---

## Notes for Next Chat

The next implementation-oriented step should remain conservative.

The project should not jump directly into DEX contracts or chain integration before the core infrastructure baseline exists.

The safest next move is to establish:

- structure
- DB foundation
- migration path
- shared technical interfaces
- environment-ready local stack