# Deep Architecture

## Objective

This document defines the detailed internal architecture of the SCAVO Exchange backend and explains how the project should evolve from the current bootstrap into a production-ready DEX-first platform.

---

## Current Real Baseline

The current codebase already establishes a minimal but coherent technical baseline.

Present structure:

- `cmd/scavo-server`
- `internal/app`
- `internal/core/config`
- `internal/core/logger`
- `internal/core/httpx`
- `internal/core/auth`
- `internal/core/ws`
- `internal/modules/auth`
- `internal/modules/system`

This means the project already has:

- application bootstrap
- config loading
- structured logging
- HTTP routing
- middleware
- WebSocket hub and dispatcher
- JWT token service
- basic auth endpoint
- module registration pattern

This is the correct foundation for the chosen modular monolith architecture.

---

## Layer Model

The backend should evolve through the following internal layers.

### 1. Bootstrap Layer

Primary responsibility:

- process startup
- config loading
- dependency creation
- module registration
- server lifecycle

Primary location:

- `cmd/scavo-server`
- `internal/app`

Rules:

- no business logic
- no transport-specific branching beyond startup concerns
- no persistence logic
- no blockchain logic

---

### 2. Core Technical Layer

Primary responsibility:

- reusable technical primitives
- shared infrastructure
- domain-independent helpers

Current and planned packages:

- `internal/core/config`
- `internal/core/logger`
- `internal/core/httpx`
- `internal/core/auth`
- `internal/core/ws`
- future `internal/core/db`
- future `internal/core/cache`
- future `internal/core/observability`
- future `internal/core/jobs`

Rules:

- must not depend on business modules
- must remain reusable
- must not embed exchange domain decisions

---

### 3. Transport Layer

Primary responsibility:

- translate external requests into internal module calls
- serialize responses
- validate input at transport boundary
- preserve correlation and security context

Transport forms:

- REST handlers
- WebSocket actions
- future admin endpoints
- future internal job triggers

Rules:

- should be thin
- should not contain domain workflows
- should call services, not repositories directly unless explicitly justified

---

### 4. Module Service Layer

Primary responsibility:

- business rules
- orchestration inside each domain module
- consistency validation
- interaction with repositories and external clients

Each module should converge toward this pattern:

- transport
- service
- repository contracts
- models / DTOs
- adapters if needed

Rules:

- service layer owns domain rules
- transport does not decide business outcomes
- service can coordinate multiple repositories and external integrations

---

### 5. Repository Layer

Primary responsibility:

- persistence access
- query composition
- transactional boundaries
- isolation of storage mechanics

Planned repository targets:

- users
- wallets
- assets
- chain cursors
- indexed events
- tracked transactions
- quotes cache
- liquidity metadata
- audit events
- future ledger entries

Rules:

- repositories should not contain transport logic
- repositories should not format API responses
- repositories should remain storage-specific

---

### 6. Integration Layer

Primary responsibility:

- communication with external systems

Primary external systems:

- SCAVIUM RPC
- DEX smart contracts
- token contracts
- Redis
- PostgreSQL
- future third-party compliance or notification providers

Rules:

- external integration code should be isolated behind clients or adapters
- module services must not embed raw low-level protocol details everywhere

---

## Target Package Direction

The project should gradually evolve toward a structure conceptually similar to:

- `cmd/scavo-server`
- `internal/app`
- `internal/core/...`
- `internal/modules/auth`
- `internal/modules/system`
- `internal/modules/user`
- `internal/modules/wallet`
- `internal/modules/chain`
- `internal/modules/asset`
- `internal/modules/portfolio`
- `internal/modules/dex`
- `internal/modules/liquidity`
- `internal/modules/quote`
- `internal/modules/routing`
- `internal/modules/txtracking`
- `internal/modules/indexer`
- `internal/modules/audit`
- `internal/modules/admin`
- future `internal/modules/ledger`
- future `internal/modules/p2p`
- future `internal/modules/compliance`

This does not require immediate creation of all directories. It defines the official growth direction.

---

## Module Responsibilities

### system
Scope:

- health-like system actions
- connectivity sanity checks
- simple operational diagnostics

### auth
Scope:

- current development login
- future real login
- JWT issuance
- refresh/session evolution
- wallet-signature-based auth in later phase

### user
Scope:

- user identity
- linked entities
- account profile metadata
- future roles and permissions

### wallet
Scope:

- linked self-custody wallets
- wallet ownership verification
- primary wallet selection
- multi-wallet support

### chain
Scope:

- SCAVIUM RPC client coordination
- chain health
- latest block awareness
- gas estimation helpers
- allowance and balance reads

### asset
Scope:

- asset registry
- token metadata
- decimals
- symbols
- icons and external metadata references if needed later

### portfolio
Scope:

- wallet balances
- token balances
- allowance views
- portfolio snapshots
- aggregated read model for frontend use

### dex
Scope:

- DEX protocol-facing backend logic
- pool reads
- pair discovery
- swap preparation
- on-chain contract coordination

### liquidity
Scope:

- add/remove liquidity support
- LP position read model
- liquidity calculations and projections

### quote
Scope:

- swap quote generation
- slippage projections
- fee projection
- min-out or max-in modeling

### routing
Scope:

- single-hop path support first
- multi-hop expansion later
- path scoring
- deterministic route selection

### txtracking
Scope:

- track submitted on-chain operations
- monitor receipts
- expose pending/confirmed/failed states

### indexer
Scope:

- ingest contract events
- maintain cursors
- handle reorg-safe synchronization in later phase
- build queryable local read models

### audit
Scope:

- operational event logging
- security-sensitive actions
- traceable state changes

### admin
Scope:

- internal operational endpoints
- future maintenance controls
- controlled access actions

### ledger
Future hybrid scope:

- custodial balances
- movements
- reservations
- settlement entries

### p2p
Future hybrid scope:

- offers
- escrow abstractions
- dispute-oriented flows
- fiat interaction placeholders

### compliance
Future hybrid scope:

- KYC hooks
- AML event points
- risk controls
- reporting/export preparation

---

## Request Flow Direction

The standard backend request direction should be:

transport -> service -> repository/client -> result mapping -> transport response

This rule should remain valid for both REST and WebSocket flows.

---

## WebSocket Direction

The current WebSocket architecture already follows a good pattern:

- handler accepts connection
- hub manages clients
- dispatcher routes actions
- session is attached if token exists
- module-specific actions are registered externally

This pattern should remain.

Future additions may include:

- authenticated subscriptions
- user channel routing
- tx status push
- quote stream push
- pool event push

---

## DEX Backend Role

The backend does not replace smart contracts.

Its role is to provide:

- contract-aware reads
- quote generation
- route calculation
- allowance inspection
- portfolio views
- transaction preparation inputs
- transaction lifecycle tracking
- frontend-ready aggregation

The backend should not initially hold private keys for DEX users.

---

## Smart Contract Boundary

DEX settlement will happen on-chain.

Therefore the backend boundary is:

- prepare
- inform
- validate
- observe
- index
- track

But not:

- sign user transactions
- custody user DEX funds by default
- internally settle DEX swaps off-chain

---

## Persistence Direction

PostgreSQL is the target primary datastore.

Redis is the target secondary infrastructure store.

### PostgreSQL should hold:
- users
- wallets
- assets
- chain state
- tracked transactions
- indexed events
- audit records
- future internal ledger entities

### Redis should hold:
- cache
- rate limit counters
- short-lived coordination
- optional ephemeral quote cache
- optional job coordination

---

## Evolution Principle

The architecture must evolve safely in this order:

1. documentation and alignment
2. core infrastructure foundation
3. identity and wallet support
4. chain reads and indexing
5. DEX contracts
6. DEX backend logic
7. frontend contract stabilization
8. hybrid expansion preparation

This ordering is intentional and must be preserved unless there is a strong technical reason to change it.

---

## Non-Goals for This Stage

At this stage, the architecture does not require:

- microservices
- Kubernetes-first design
- full event-driven architecture
- matching engine
- custodial ledger implementation
- fiat processing
- compliance subsystem implementation

Those may come later, but they are not part of the current architectural baseline.