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

## Phase 0.2.2 Scope

This phase defines the official persistence and local environment direction for the backend.

The purpose of this phase is to lock:

- persistence roles
- storage boundaries
- migration workflow
- local infrastructure expectations
- environment configuration baseline
- repository preparation rules

This phase does not yet require the full implementation of DB access or cache integration, but it establishes the official structure those implementations must follow.

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
- no environment-specific branching outside startup concerns

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

---

## Persistence Direction

The persistence model is intentionally split into two roles.

### PostgreSQL

PostgreSQL is the primary system of record.

It is responsible for durable, queryable, and auditable data.

Primary persistence targets:

- users
- linked wallets
- sessions or future refresh token state if persisted
- assets metadata
- chain cursors
- indexed blockchain events
- tracked transactions
- quote-related persistent metadata if needed later
- audit records
- future ledger entities
- future P2P entities
- operational state that must survive restarts

PostgreSQL is the source of truth for critical backend state.

### Redis

Redis is the secondary infrastructure store.

It is responsible for short-lived and coordination-oriented state.

Primary Redis targets:

- short-lived cache
- rate-limit counters
- temporary coordination keys
- optional locks
- optional quote cache
- optional WebSocket-related ephemeral support
- optional job coordination
- optional chain polling coordination

Redis is not the system of record.

Persistence-critical or audit-critical data must not rely exclusively on Redis.

---

## Persistence Boundary Rules

The following rules are official:

- durable data belongs in PostgreSQL
- ephemeral data belongs in Redis only when necessary
- no critical state should exist only in memory if it must survive process restart
- no critical business truth should exist only in Redis
- repositories must make storage ownership explicit
- cache use must remain optional and replaceable

These boundaries are important to avoid architectural confusion later.

---

## Repository Direction

The project should prepare for repository-based persistence.

Repository responsibilities include:

- persistence reads and writes
- query isolation
- transaction handling
- storage-specific mapping

Repository boundaries should remain module-oriented.

Example direction:

- user repositories belong to the user module
- wallet repositories belong to the wallet module
- chain cursor repositories belong to chain or indexer
- audit repositories belong to audit

Common DB wiring should remain inside reusable core infrastructure.

The project should avoid one giant shared repository package that centralizes all domain persistence in a single place.

---

## Migration Direction

Schema evolution will be controlled through versioned migrations.

A dedicated `migrations/` directory is part of the official repository structure.

Migration principles:

- all schema changes are versioned
- environments must be reproducible
- manual schema drift is discouraged
- migrations must be reviewable
- migrations must be aligned with repository evolution
- migration history becomes part of the project record

The migration system should support:

- up migrations
- down migrations when safe and practical
- local development usage
- CI usage later
- internal testing environment setup later

---

## Configuration Direction

The backend configuration must remain explicit and environment-driven.

The configuration baseline should include at least:

- application environment
- HTTP bind configuration
- JWT settings
- PostgreSQL settings
- Redis settings
- SCAVIUM RPC settings
- CORS-related settings if needed
- development flags where appropriate

Configuration should distinguish between:

- local development
- internal testing
- production-oriented environments later

No environment-specific behavior should depend on undocumented constants in code.

---

## Expected Base Environment Variables

The exact implementation may evolve, but the baseline configuration model should reserve space for variables such as:

- `APP_ENV`
- `APP_NAME`
- `HTTP_ADDR`
- `JWT_SECRET`
- `JWT_ISSUER`
- `JWT_TTL_MINUTES`
- `POSTGRES_HOST`
- `POSTGRES_PORT`
- `POSTGRES_DB`
- `POSTGRES_USER`
- `POSTGRES_PASSWORD`
- `POSTGRES_SSLMODE`
- `REDIS_ADDR`
- `REDIS_PASSWORD`
- `REDIS_DB`
- `SCAVIUM_RPC_URL`
- `SCAVIUM_CHAIN_ID`
- `LOG_LEVEL`

These names define the configuration direction and may later be expanded with more granular settings.

---

## Local Environment Direction

The project should support a consistent local development baseline.

The minimum local environment target is:

- backend app
- PostgreSQL
- Redis

Optional later additions may include:

- migration runner
- local contract deployment tooling
- seeded development data
- local observability helpers

The first local environment goal is not maximum completeness.

The first goal is deterministic development setup.

---

## Docker and Local Infrastructure Direction

The project should support Docker-based local infrastructure.

The target local infrastructure direction includes:

- a backend service
- a PostgreSQL service
- a Redis service
- optional bind-mounted config or environment file support
- optional one-shot migration service later

Docker support is intended for:

- reproducible onboarding
- consistent local setup
- future CI alignment
- easier internal testing preparation

A Docker-based development environment is recommended even if developers can also run services manually.

---

## Observability Direction

Observability remains part of the infrastructure baseline.

Planned observability components include:

- structured logs
- request correlation
- health and readiness visibility
- metrics later
- tracing later

Persistence-related observability should eventually make it possible to identify:

- DB connection failures
- migration failures
- Redis connectivity problems
- chain integration degradation
- queue or worker issues later

---

## Job and Background Processing Direction

The backend is initially single-process, but it must remain ready for background work.

Persistence and cache definitions should preserve room for future jobs such as:

- chain polling
- event indexing
- receipt tracking
- quote cache refresh
- cleanup tasks
- notification dispatch later

Redis may support coordination for some of these jobs, but durable execution state should not depend on Redis alone when long-lived traceability is required.

---

## Request Flow Direction

The standard backend flow remains:

transport -> service -> repository and client adapters -> result mapping -> response

For persistence-aware write operations, the direction becomes:

transport -> service -> repository transaction and external coordination -> audit event when applicable -> response

This flow is important for future database-backed modules.

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
4. observability and test bootstrap
5. identity and wallet support
6. chain reads and indexing baseline
7. DEX contracts
8. DEX backend logic
9. frontend-ready contract stabilization
10. hybrid expansion preparation

This order is intentional.

---

## Non-Goals for This Stage

This stage does not require:

- full SQL schema design
- full repository implementation
- Redis-based features in production form
- Docker orchestration implementation
- microservices
- matching engine
- custodial ledger implementation
- fiat operations
- compliance implementation
- production-grade indexing yet

This phase only locks the persistence and environment model that future implementation must follow.

## Phase 0.7 — Application Layer Foundation (Deep Analysis)

### Context

Prior to Phase 0.7, the backend exhibited a mixed orchestration model where HTTP handlers partially controlled execution flow, directly invoking services and assembling responses.

This led to:

- duplicated orchestration logic
- inconsistent handler responsibilities
- weak separation between transport and domain execution

---

### Architectural Shift

Phase 0.7 introduces a strict separation:

Transport Layer (HTTP Handlers)  
→ Application Layer (Use-case orchestration)  
→ Services / Stores (Execution)

---

### Application Layer Responsibilities

The application layer now:

- defines explicit use cases
- orchestrates multiple services
- centralizes flow control
- ensures deterministic execution paths

It does NOT:

- handle HTTP concerns
- perform low-level DB logic
- define transport formats

---

### Transport Layer Responsibilities

Handlers are now strictly limited to:

- request decoding
- claims extraction
- invoking application methods
- mapping errors to HTTP responses
- writing JSON responses

---

### Phase Breakdown

#### 0.7.1 — Boundary Introduction

- introduced `Application` struct
- bootstrap flow used as initial entry point

#### 0.7.2 — Authenticated Surface Extraction

- login, session, profile flows moved into application
- handlers reduced to delegators

#### 0.7.3 — Wallet Consolidation

- all wallet-related flows moved into application
- removed remaining handler-driven orchestration

#### 0.7.4 — Handler Simplification

- unified transport helpers
- removed duplication
- enforced consistent handler pattern

---

### Resulting Execution Model

All flows now follow:

HTTP → Handler → Application → Services → Response

With:

- single orchestration layer
- consistent transport behavior
- no cross-layer leakage

---

### Architectural Guarantees

- strict separation of concerns
- deterministic execution paths
- contract preservation across refactors
- extensibility for future phases

---


### Phase 0.9 — API Versioning Strategy (Deep Analysis)

#### Context

By the end of Phase 0.8, the backend already exposes a stable authenticated HTTP surface with a frozen structured error envelope, but it still does so through unversioned public routes. At that point the system has enough transport maturity that contract evolution becomes a first-class architectural concern rather than a future detail.

A versioning phase becomes necessary because the backend now has to balance three realities simultaneously:

- stable current consumers on the existing route set
- future transport and authorization growth after 0.8
- a frontend that intentionally remains aligned only up to backend Phase 0.6 until Stage 0 is fully closed

#### Architectural Problem

Without explicit versioning semantics, any later transport evolution would overload the meaning of the current route set:

- new canonical behaviors would have no formal home
- legacy compatibility would be ambiguous
- breaking changes could be introduced informally
- future authorization or domain-layer additions could create accidental contract drift

The problem after 0.8 is therefore not route proliferation, but uncontrolled contract evolution.

#### Architectural Direction

Phase 0.9 defines a two-surface model:

- canonical versioned API surface under `/api/v1/...`
- legacy compatibility surface under the current unversioned routes

This does not create two business implementations. It creates one business/application surface with two transport entry paths, where the legacy entry remains compatibility-oriented and the canonical entry becomes the authoritative contract target.

#### Design Consequences

This decision establishes the following architectural rules:

- API versioning belongs to the transport contract, not to application/service orchestration
- `v1` freezes the current success-payload semantics and the 0.8 standardized error envelope
- legacy routes are preserved as non-canonical compatibility routes while Stage 0 remains in flight
- future breaking transport changes must target a new API version rather than silently mutating `v1`
- authorization in 0.10 must layer on top of a defined versioning policy instead of inventing one during permission work

#### Why This Matters Before 0.10

Authorization introduces contract-sensitive behavior: forbidden/unauthorized handling, future role-sensitive route behavior, and additional policy signals. Without a defined versioning model, those additions could blur the boundary between “current stable behavior” and “future evolved behavior”.

Phase 0.9 avoids that by fixing the transport evolution model first.



#### 0.9.2 Implementation Consequence

With 0.9.2, the two-surface model is no longer only architectural policy. The HTTP router now materializes both entry paths:

- legacy compatibility exposure under `/auth/...`
- canonical exposure under `/api/v1/auth/...`

Crucially, this is implemented as shared transport registration over the same handlers rather than duplicated route-specific code paths. That preserves the 0.7 application boundary and keeps versioning as a projection concern in the transport layer.

This matters because later phases can now reason about canonical route behavior against a real router surface instead of a merely documented target.

---
### Impact on Future Phases

This foundation enables:

- Phase 0.8 → standardized error model
- Phase 0.9 → versioning strategy
- Phase 0.10 → authorization layer

Without the risk of cross-layer inconsistencies.
