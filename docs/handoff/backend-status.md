# Backend Status

## Project

SCAVO Exchange - Backend

---

## Current Stage

Stage 0 - Foundation

## Current Phase

Phase 0.4 - Auth and User Stabilization

## Current Subphase

Phase 0.4.1 - Auth and User Module Stabilization

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
- auth service boundary for development login orchestration
- current authenticated identity resolution through bearer token
- system WebSocket handler
- auth WebSocket handler registration
- PostgreSQL core scaffolding
- Redis core scaffolding
- status service for health and readiness
- readiness-aware router wiring
- first migration-backed domain table
- first repository-backed domain service
- minimal authenticated REST identity read path through `GET /auth/me`
- unit and integration validation baseline for first persistence path
- auth and user regression coverage expanded

---

## Current Modules

- system
- auth
- user

---

## Current Core Packages

- config
- logger
- httpx
- auth
- ws
- db
- cache
- status

---

## What This Subphase Implemented

This subphase implemented:

- extraction of development login orchestration into an auth service
- formalized auth-to-user interaction through service boundaries
- `user` identity retrieval path through repository and service expansion
- `GET /auth/me` as the first minimal authenticated REST identity endpoint
- auth service tests for login and current-user resolution
- HTTP handler tests for login and current-user responses
- additional user service tests for identity retrieval behavior

---

## What Is Still Not Implemented

Not implemented yet:

- refresh token persistence
- token revocation
- session persistence
- wallet challenge generation
- wallet signature verification
- role/permission model
- Redis-backed auth/session coordination
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

## Validation Status

The project now supports:

- unit validation of core domain logic
- integration validation of PostgreSQL repository layer
- readiness validation with dependency awareness
- migration workflow execution and inspection
- smoke-level validation of login flow
- unit validation of auth login orchestration
- unit validation of current-user resolution from JWT
- HTTP handler validation for `/auth/login` and `/auth/me`

The backend now has a clearer auth-user boundary and a minimal authenticated read path, while remaining in controlled development-bootstrap mode.

---

## Recommended Next Step

Phase 0.4.2 - Token Lifecycle and Auth Transport Hardening

Recommended scope:

- formalize token extraction/parsing helpers for HTTP
- prepare auth middleware strategy without overcommitting too early
- define token lifecycle expectations before refresh-token persistence
- keep wallet auth deferred until the bootstrap auth boundary is fully stable