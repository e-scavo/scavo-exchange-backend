# Backend Status

## Project

SCAVO Exchange - Backend

---

## Current Stage

Stage 0 - Foundation

## Current Phase

Phase 0.4 - Auth and User Stabilization

## Current Subphase

Phase 0.4.2 - Token Lifecycle and Auth Transport Hardening

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
- reusable HTTP auth middleware
- shared token extraction utilities for HTTP and WebSocket transports
- auth claims context utilities in the core auth layer
- system WebSocket handler
- auth WebSocket handler registration
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

- shared token extraction helpers for transport consistency
- auth claims storage and retrieval through request context
- reusable HTTP auth middleware with claims injection
- `GET /auth/me` protection through middleware instead of local token parsing
- WebSocket auth alignment using the same token extraction strategy
- removal of transport coupling that caused cyclic dependency risk
- rescue and stabilization of the 0.4.2 implementation until build and tests were both green

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

- successful `go build ./...`
- successful `go test ./...`
- unit validation of auth login orchestration
- unit validation of current-user resolution
- unit validation of user service behavior
- unit validation of status/readiness logic
- integration validation baseline for PostgreSQL-backed user repository
- middleware-protected authenticated identity flow through the current transport model

The backend now has a coherent minimal auth transport layer without prematurely introducing session persistence or wallet-signature authentication.

---

## Recommended Next Step

Phase 0.4.3 - Session Evolution and Wallet Auth Preparation

Recommended scope:

- define a stable authenticated session model across HTTP and WebSocket surfaces
- prepare the backend for future wallet-based authentication flows
- avoid refresh/revocation persistence until session shape and auth lifecycle expectations are clearer