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
- `/version`
- `/auth/login`
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
- HTTP helpers and middleware
- WebSocket protocol and routing
- future database helpers
- future cache helpers
- future queue/worker helpers
- future observability helpers

Current implementation:

- `internal/core/config`
- `internal/core/logger`
- `internal/core/httpx`
- `internal/core/auth`
- `internal/core/ws`

This layer must remain domain-agnostic.

---

## Domain Modules

Domain modules contain business logic and represent functional areas of the exchange.

Current modules:

- `system`
- `auth`

Planned modules:

- `user`
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
- fiat ramps
- optional order-book trading

These future capabilities must be added without breaking the DEX-first architecture.

---

## Persistence

Persistence is not fully implemented yet, but the target model is:

- PostgreSQL as primary relational database
- Redis for cache, locks, short-lived coordination, and future queue support

The persistence layer will support:

- users
- linked wallets
- asset metadata
- indexed blockchain events
- transaction tracking
- liquidity and pool metadata
- future custodial ledger entries
- audit events
- operational state

---

## Communication Model

The backend exposes two primary communication models:

### REST
Used for:

- request/response operations
- reads
- session creation
- metadata retrieval
- quote requests
- portfolio data
- admin operations

### WebSocket
Used for:

- real-time status
- user session awareness
- live transaction updates
- quote stream possibilities
- pool updates
- internal product interactivity

---

## Observability and Validation

Observability is a first-class architectural concern.

The backend must evolve with explicit support for:

- structured logs
- request correlation
- recoverable operational diagnostics
- health and readiness boundaries
- metrics
- future tracing
- testable infrastructure seams
- validation across unit, integration, and end-to-end layers

This is required because the project will later depend on:

- blockchain integrations
- persistent state
- event ingestion
- background jobs
- hybrid growth paths

Without early observability and test structure, those phases would become fragile and difficult to operate.

---

## Security Model

The architecture is designed around explicit security boundaries.

Current and planned security controls include:

- JWT authentication
- wallet-signature-based authentication
- CORS control
- request ID propagation
- panic recovery
- structured logging
- future rate limiting
- future abuse prevention
- future role and permission controls
- future audit trail enforcement

---

## Deployment Model

The deployment model is intentionally simple at the start:

- single backend process
- container-friendly
- environment-driven configuration
- future support for worker modes if needed

This keeps development and internal testing fast while preserving room for later operational separation.

---

## Architectural Direction

The architecture is deliberately optimized for:

- fast initial delivery
- safe iterative growth
- low operational complexity
- strong documentation alignment
- backend/frontend contract clarity
- future hybrid expansion without premature overengineering
- early operational visibility
- testable growth across infrastructure and domain phases