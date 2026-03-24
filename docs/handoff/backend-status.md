# Backend Status

## Project

SCAVO Exchange - Backend

---

## Current Stage

Stage 0 - Foundation

## Current Phase

Phase 0.3 - Infrastructure Bootstrap

## Current Subphase

Phase 0.3.1 - Implementation Bootstrap for Persistence and Health Infrastructure

---

## Real Current Code Baseline

The backend now includes:

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
- PostgreSQL core scaffolding
- Redis core scaffolding
- status service for health and readiness
- readiness-aware router wiring

Current real modules:

- `system`
- `auth`

Current real core packages:

- `config`
- `logger`
- `httpx`
- `auth`
- `ws`
- `db`
- `cache`
- `status`

---

## What This Subphase Implemented

This subphase implemented:

- infrastructure-oriented config expansion
- PostgreSQL bootstrap client scaffolding
- Redis bootstrap client scaffolding
- explicit health versus readiness separation
- reusable dependency checker model
- readiness endpoint wiring
- startup visibility for enabled infrastructure

---

## What Is Still Not Implemented

Not implemented yet:

- migrations
- repository scaffolding
- docker-compose local stack
- real DB-backed modules
- real Redis-backed coordination
- metrics endpoint
- tracing
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

Phase 0.3.2 - Migration Bootstrap and Local Infrastructure Layout

Recommended scope:

- create migrations directory baseline
- define migration execution workflow
- prepare local Docker stack direction
- keep current bootstrap stable while making persistence reproducible