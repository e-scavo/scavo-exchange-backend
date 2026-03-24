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

## Phase 0.2.1 Scope

This phase defines the official technical growth layout for the backend.

It does not require immediate implementation of all directories or modules, but it locks the intended structure so that future code can be added without architectural drift.

The objective is to preserve the current bootstrap while introducing a stable structural direction for:

- persistence
- repositories
- chain integrations
- shared services
- module growth
- local infrastructure
- future workers
- observability

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
- no persistence logic
- no blockchain logic
- no transport payload decisions beyond application wiring

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
- future `internal/core/clock`
- future `internal/core/ids`

Rules:

- must not depend on business modules
- must remain reusable
- must not encode exchange domain workflows

---

### 3. Transport Layer

Primary responsibility:

- translate external requests into internal module calls
- validate transport input
- serialize output
- preserve request and security context

Transport forms:

- REST handlers
- WebSocket actions
- future internal admin routes
- future internal operational endpoints

Rules:

- thin by default
- transport-only validation at the boundary
- no repository orchestration directly from handlers unless explicitly justified

---

### 4. Module Service Layer

Primary responsibility:

- business rules
- orchestration inside each domain module
- interaction with repositories and integrations
- consistency validation

Each module should converge toward this structure:

- transport
- service
- repository contracts
- DTOs or transport payloads
- internal helpers
- module models

Rules:

- service layer owns business decisions
- transport layer does not own workflows
- services may combine persistence and blockchain integrations when necessary

---

### 5. Repository Layer

Primary responsibility:

- persistence access
- query composition
- transaction boundaries
- storage isolation

Repository layer targets:

- users
- linked wallets
- assets
- chain cursors
- indexed events
- tracked transactions
- audit entries
- future ledger entries
- operational state

Rules:

- repositories must not format API responses
- repositories must not contain transport logic
- repositories should remain storage-specific and testable

---

### 6. Integration Layer

Primary responsibility:

- communication with external systems

Primary external systems:

- SCAVIUM RPC
- DEX smart contracts
- token contracts
- PostgreSQL
- Redis
- future compliance providers
- future notification providers

Rules:

- external integrations should be isolated behind clients or adapters
- raw low-level protocol code should not spread across services

---

## Official Target Layout

The backend should gradually evolve toward the following structure:

- `cmd/scavo-server`
- `internal/app`
- `internal/core/config`
- `internal/core/logger`
- `internal/core/httpx`
- `internal/core/auth`
- `internal/core/ws`
- `internal/core/db`
- `internal/core/cache`
- `internal/core/observability`
- `internal/core/jobs`
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
- `internal/platform/scavium`
- `internal/platform/contracts`
- future `internal/platform/notifications`
- `migrations`
- `deployments`
- `scripts`
- `docs`

This structure is the official architectural direction.

It does not mean every directory must be created immediately, but it defines where future work belongs.

---

## Module Boundaries

### auth
Scope:

- current development login
- future production login
- JWT issuance
- refresh and session evolution
- wallet signature authentication in later phases

### system
Scope:

- health
- version
- diagnostics
- operational sanity checks

### user
Scope:

- user identity
- profile metadata
- linked entities
- future permissions and roles

### wallet
Scope:

- self-custody wallet linking
- ownership verification
- wallet preferences
- multi-wallet support

### chain
Scope:

- RPC coordination
- network awareness
- gas estimation
- native balance reads
- token balance reads
- allowance reads

### asset
Scope:

- asset registry
- token metadata
- decimals
- symbols
- chain-native asset definitions

### portfolio
Scope:

- aggregated wallet balances
- allowance summaries
- frontend-ready portfolio view
- token holdings snapshots

### dex
Scope:

- contract-facing DEX backend logic
- pool discovery
- pair state reads
- swap preparation support

### liquidity
Scope:

- add liquidity support
- remove liquidity support
- LP position reads
- pool share calculations

### quote
Scope:

- quote generation
- min-out modeling
- fee estimation
- slippage-aware responses

### routing
Scope:

- path selection
- single-hop first
- multi-hop expansion later
- deterministic route selection

### txtracking
Scope:

- transaction lifecycle registration
- pending/confirmed/failed state tracking
- receipt status updates

### indexer
Scope:

- contract event ingestion
- chain cursor tracking
- local read-model updates
- reorg-aware evolution later

### audit
Scope:

- auditable operational events
- sensitive action traceability
- structured internal event records

### admin
Scope:

- operational internal endpoints
- maintenance controls
- protected internal actions

### ledger
Future hybrid scope:

- custodial balances
- reservations
- internal settlement entries
- future exchange-controlled accounting

### p2p
Future hybrid scope:

- offers
- escrow abstractions
- disputes
- fiat-adjacent coordination

### compliance
Future hybrid scope:

- KYC hooks
- AML hooks
- reporting/export boundaries
- address screening integration points

---

## Platform and Adapter Direction

Some integrations are too specific to belong directly inside domain modules.

To avoid coupling every module to low-level protocol code, the project should introduce platform-specific packages such as:

- `internal/platform/scavium`
- `internal/platform/contracts`

These packages are intended for:

- RPC clients
- contract call helpers
- ABI-based integrations
- chain-specific retry and failover logic
- event decoding helpers

Modules should consume these through explicit interfaces whenever possible.

---

## Persistence Direction

PostgreSQL is the target primary datastore.

Redis is the target secondary infrastructure store.

### PostgreSQL responsibilities

- user records
- linked wallets
- assets metadata
- chain sync cursors
- indexed on-chain events
- tracked transactions
- audit records
- future internal ledger entities
- future P2P entities

### Redis responsibilities

- short-lived cache
- coordination keys
- rate-limit counters
- temporary quote cache if needed
- optional worker coordination
- optional lock support

Redis is not the system of record.

---

## Migrations Direction

The project should formally adopt a migration-based persistence workflow.

A dedicated `migrations/` directory should exist and become the source of truth for schema evolution.

Migration rules:

- no undocumented schema drift
- no manual production-only schema edits
- all structural DB changes must be versioned
- migrations must align with repository behavior and documentation

---

## Observability Direction

Observability will become part of the infrastructure baseline.

Planned observability components:

- structured logs
- request correlation
- internal domain event logging
- metrics
- health and readiness signals
- future tracing

Observability support should live in reusable technical infrastructure, not inside each module ad hoc.

---

## Job and Background Processing Direction

The system is initially single-process, but it must remain ready for background work.

Future job categories may include:

- chain polling
- event indexing
- receipt tracking
- cache refresh
- operational cleanup
- notifications

The initial direction is to keep jobs inside the same codebase, with the option to split process roles later if needed.

---

## Request Flow Direction

The standard backend flow should be:

transport -> service -> repository and client adapters -> result mapping -> response

This must remain valid for REST and WebSocket flows.

The standard write flow should be:

transport -> service -> transactional persistence and external coordination -> audit/event record -> response

---

## DEX Backend Boundary

The backend does not replace smart contracts.

Its role is to provide:

- contract-aware reads
- pool discovery
- quote generation
- routing assistance
- allowance inspection
- portfolio aggregation
- transaction preparation support
- transaction lifecycle tracking
- frontend-ready state views

The backend should not initially:

- hold user private keys for DEX usage
- internally settle DEX swaps
- execute matching-engine style trades
- replace wallet signing authority

---

## Evolution Principle

The architecture must evolve safely in this order:

1. documentation and alignment
2. infrastructure layout and shared foundation
3. persistence and environment baseline
4. identity and wallet support
5. chain reads and indexing baseline
6. DEX contracts
7. DEX backend logic
8. frontend-ready contract stabilization
9. hybrid expansion preparation

This order is intentional and should only change for a strong technical reason.

---

## Non-Goals for This Stage

This stage does not require:

- microservices
- Kubernetes-first design
- matching engine
- order-book implementation
- custodial ledger implementation
- fiat operations
- compliance implementation
- production-grade indexing yet

This phase only defines the structure required to support those later if the project grows into them.