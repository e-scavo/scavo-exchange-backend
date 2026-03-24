# Backend Status

## Project

SCAVO Exchange - Backend

---

## Current Stage

Stage 0 - Foundation

## Current Phase

Phase 0.3 - Infrastructure Bootstrap

## Current Subphase

Phase 0.3.3 - Repository and First Persistence Module

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
- first migration-backed domain table
- first repository-backed domain service

Current real modules:

- `system`
- `auth`
- `user`

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

- first real domain migration: `users`
- first repository contract and PostgreSQL repository
- first persistence-backed domain service
- safe integration between auth login and persisted users
- fallback compatibility when PostgreSQL is not configured

---

## What Is Still Not Implemented

Not implemented yet:

- refresh tokens persistence
- repository test suite
- migrations runner integration in app lifecycle
- docker-compose validation
- Redis-backed features
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

Phase 0.3.4 - Repository Validation and Migration Workflow Hardening

Recommended scope:

- validate migration flow end to end
- add repository tests
- add local execution notes for DB-backed login
- prepare the codebase for the next real persisted module