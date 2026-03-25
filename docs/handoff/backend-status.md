# Backend Status

## Project

SCAVO Exchange - Backend

---

## Current Stage

Stage 0 - Foundation

## Current Phase

Phase 0.4 - Auth and User Stabilization

## Current Subphase

Phase 0.4.4 - Wallet Challenge Contract and Nonce Bootstrap

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
- wallet challenge bootstrap REST path through `POST /auth/wallet/challenge`
- reusable HTTP auth middleware
- shared token extraction utilities for HTTP and WebSocket transports
- auth claims context utilities in the core auth layer
- enriched WebSocket session metadata derived from JWT claims
- authenticated WebSocket actions `auth.whoami` and `auth.session`
- wallet challenge service with stable message generation
- cryptographically secure nonce generation
- in-memory wallet challenge store with TTL cleanup behavior
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

- the first wallet-auth bootstrap contract through `POST /auth/wallet/challenge`
- wallet challenge request and response DTOs
- secure nonce generation
- stable signable challenge message construction
- bootstrap challenge TTL configuration
- bootstrap in-memory challenge storage
- address-format validation for EVM-style wallet addresses
- preparation for the next step where wallet signatures will actually be verified

---

## What Is Still Not Implemented

Not implemented yet:

- wallet signature verification
- wallet-based JWT issuance
- challenge consumption and replay prevention persistence
- challenge persistence in PostgreSQL or Redis
- refresh token persistence
- token revocation
- session persistence
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
- unit validation of wallet challenge generation
- unit validation of wallet challenge HTTP contract
- unit validation of user service behavior
- unit validation of status/readiness logic
- integration validation baseline for PostgreSQL-backed user repository

The backend now has a stable bootstrap contract for wallet-auth challenge issuance without prematurely enabling real signature-based authentication.

---

## Recommended Next Step

Phase 0.4.5 - Wallet Signature Verification and Token Issuance

Recommended scope:

- verify EVM signatures against issued wallet challenges
- consume valid challenges safely
- resolve or create wallet-authenticated identities
- mint JWT from wallet-authenticated flows