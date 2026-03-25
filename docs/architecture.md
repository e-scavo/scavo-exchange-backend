# Architecture Overview

## Objective

The SCAVO Exchange backend is designed as a modular monolith that can evolve from an initial DEX-first backend into a broader hybrid exchange platform.

The system must support:

- DEX-first non-custodial trading
- Direct SCAVIUM blockchain integration
- Wallet-native user flows
- Future custodial and hybrid trading models
- REST and WebSocket APIs
- Background processing
- Strong observability and security boundaries

---

## Architecture Style

The official architecture style is:

**Modular Monolith**

This means:

- a single deployable backend binary
- clear domain module boundaries
- shared core infrastructure
- internal package-level composition
- no early microservice fragmentation

This approach is intentionally chosen to maximize delivery speed, simplify local development, reduce operational overhead, and preserve the ability to split services later if the product requires it.

---

## High-Level Structure

The backend is logically divided into the following layers:

1. Entry Layer
2. Application Composition Layer
3. Core Infrastructure Layer
4. Domain Modules
5. Persistence and External Integrations
6. Blockchain and Contract Integration
7. Background Jobs and Internal Processing
8. Observability and Validation Layer

---

## Entry Layer

This layer exposes the backend to clients and external systems.

Initial responsibilities:

- HTTP routing
- REST endpoints
- WebSocket upgrade and session handling
- request middleware
- response serialization
- request correlation and recovery

Current implementation already includes:

- `/health`
- `/readiness`
- `/version`
- `/auth/login`
- `/auth/me`
- `/ws`

---

## Application Composition Layer

This layer wires the application together.

Responsibilities:

- config loading
- logger initialization
- token service initialization
- WebSocket hub creation
- dispatcher registration
- module registration
- server boot and shutdown lifecycle
- repository/service wiring across modules

Current implementation:

- `internal/app`

This layer must remain orchestration-only and must not contain domain business logic.

---

## Core Infrastructure Layer

This layer contains reusable technical building blocks used by domain modules.

Current and planned responsibilities:

- configuration
- structured logging
- JWT/token services
- auth transport helpers
- auth claims context helpers
- HTTP helpers and middleware
- WebSocket protocol and routing
- database helpers
- cache helpers
- future queue/worker helpers
- future observability helpers

Current implementation:

- `internal/core/config`
- `internal/core/logger`
- `internal/core/httpx`
- `internal/core/auth`
- `internal/core/ws`
- `internal/core/db`
- `internal/core/cache`
- `internal/core/status`

This layer must remain domain-agnostic.

---

## Domain Modules

Domain modules contain business logic and represent functional areas of the exchange.

Current modules:

- `system`
- `auth`
- `user`

Planned modules:

- `wallet`
- `chain`
- `asset`
- `portfolio`
- `dex`
- `liquidity`
- `quote`
- `routing`
- `txtracking`
- `indexer`
- `ledger`
- `p2p`
- `admin`
- `audit`
- `compliance`

Each module should expose:

- handlers or transport adapters
- service layer
- repository contracts
- DTOs or transport payloads
- internal module-specific helpers

---

## Current Auth/User Boundary

At the current stage:

- the `auth` module owns token-oriented login orchestration
- the `user` module owns user persistence and user identity reads
- HTTP handlers remain thin
- HTTP auth enforcement is handled through reusable middleware
- authenticated claims travel through request context
- token extraction is shared instead of duplicated across transports
- WebSocket auth attachment follows the same extraction strategy
- persisted development login remains enabled
- authenticated identity read is available through `GET /auth/me`

This is intentionally still a bootstrap auth model.

The project is not yet implementing:

- refresh token persistence
- wallet signature login
- role/permission enforcement
- revocation flows
- full session persistence

Those will come later.

---

## Blockchain Integration

SCAVIUM is the primary chain for the exchange backend.

The backend will interact directly with:

- SCAVIUM RPC endpoints
- deployed DEX contracts
- token contracts
- indexed chain data
- transaction receipts and logs

The blockchain integration layer must support:

- direct read operations
- quote assistance
- allowance checks
- gas estimation
- transaction tracking
- event ingestion
- future failover across multiple RPC endpoints

The backend does not initially custody keys for DEX users.

---

## DEX Model

The initial product focus is DEX-first and non-custodial.

The first protocol model will be:

**AMM v1**

This includes:

- pools
- liquidity provisioning
- swaps
- quote generation
- routing support
- allowance inspection
- settlement through on-chain contracts

The backend will prepare and expose trading intelligence, but the user wallet remains the signing authority.

---

## Hybrid Growth Model

Although the initial implementation is DEX-first, the system is intentionally designed to support future hybrid expansion.

Future hybrid scope includes:

- custodial accounts
- internal balances
- internal ledger
- deposit and withdrawal orchestration
- P2P support