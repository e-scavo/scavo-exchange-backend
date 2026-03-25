# Backend Status

## Project

SCAVO Exchange - Backend

---

## Current Stage

Stage 0 - Foundation

## Current Phase

Phase 0.4 - Auth and User Stabilization

## Current Subphase

Phase 0.4.3 - Session Evolution and Wallet Auth Preparation

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
- centralized auth service for development login orchestration
- development HTTP login
- authenticated current-user REST path through `GET /auth/me`
- authenticated session REST path through `GET /auth/session`
- reusable HTTP auth middleware
- shared token extraction utilities for HTTP and WebSocket transports
- auth claims context utilities in the core auth layer
- enriched WebSocket session metadata derived from JWT claims
- authenticated WebSocket actions `auth.whoami` and `auth.session`
- system WebSocket handler
- PostgreSQL core scaffolding
- Redis core scaffolding
- status service for health and readiness
- readiness-aware router wiring
- first migration-backed domain table
- first repository-backed domain service
- user identity retrieval through repository and service expansion
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

- a formal shared session representation through `SessionView`
- an authenticated REST session introspection endpoint at `GET /auth/session`
- richer WebSocket session state including issuer, subject, and expiration metadata
- a new authenticated WebSocket action `auth.session`
- preservation of `auth.whoami` while expanding it with more session-aware fields
- WebSocket auth registration with access to the real auth service
- preparation for future wallet challenge and signature flows without coupling the system to them prematurely

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

The project should now support validation for:

- successful `go build ./...`
- successful `go test ./...`
- unit validation of auth login orchestration
- unit validation of current-user resolution
- unit validation of session view resolution
- unit validation of user service behavior
- unit validation of status/readiness logic
- integration validation baseline for PostgreSQL-backed user repository
- authenticated session introspection through both REST and WebSocket transport surfaces

The backend now has a coherent session-shaped auth surface that is still bootstrap-safe and ready for the first wallet-auth contract design step.

---

## Recommended Next Step

Phase 0.4.4 - Wallet Challenge Contract and Nonce Bootstrap

Recommended scope:

- define request and response contracts for wallet challenge creation
- introduce nonce/challenge generation rules and lifecycle expectations
- keep signature verification for the next step after nonce contracts are stable