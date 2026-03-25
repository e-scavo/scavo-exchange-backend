# Backend Status

## Project

SCAVO Exchange - Backend

---

## Current Stage

Stage 0 - Foundation

## Current Phase

Phase 0.4 - Auth and User Stabilization

## Current Subphase

Phase 0.4.5 - Wallet Signature Verification and Token Issuance

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
- wallet challenge verification REST path through `POST /auth/wallet/verify`
- reusable HTTP auth middleware
- shared token extraction utilities for HTTP and WebSocket transports
- auth claims context utilities in the core auth layer
- enriched WebSocket session metadata derived from JWT claims
- authenticated WebSocket actions `auth.whoami` and `auth.session`
- wallet challenge service with stable message generation
- cryptographically secure nonce generation
- in-memory wallet challenge store with TTL cleanup behavior
- EVM wallet signature recovery and address verification
- replay-safe challenge consumption for successful wallet sign-in
- wallet-authenticated JWT issuance with wallet metadata claims
- wallet-aware REST and WebSocket session metadata
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

- EVM-style wallet signature verification through `POST /auth/wallet/verify`
- wallet challenge consumption with replay prevention in the bootstrap memory store
- wallet-auth JWT issuance after successful challenge verification
- wallet metadata propagation in JWT claims, REST session responses, and WebSocket session metadata
- wallet-auth fallback user resolution without prematurely requiring database persistence
- regression coverage for signature recovery, successful verification, and replay rejection

---

## What Is Still Not Implemented

Not implemented yet:

- challenge persistence in PostgreSQL or Redis
- durable wallet identity persistence and user-wallet linking
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
- unit validation of wallet challenge verification and replay protection
- unit validation of wallet challenge HTTP contracts
- unit validation of user service behavior
- unit validation of status/readiness logic
- integration validation baseline for PostgreSQL-backed user repository

The backend now exposes a complete bootstrap wallet-auth flow from challenge issuance to signature verification and JWT minting, while still keeping challenge and wallet identity storage in the current non-durable bootstrap layer.

---

## Recommended Next Step

Phase 0.4.6 - Wallet Identity Persistence and Durable Challenge Storage

Recommended scope:

- persist challenges beyond process memory using PostgreSQL or Redis
- introduce durable wallet identity storage and wallet-to-user linking direction
- preserve replay protection semantics across restarts
- prepare the backend for multi-wallet account evolution
