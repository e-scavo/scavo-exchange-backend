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

#### 0.9.3 Architectural Consequence

With 0.9.3, the authenticated surface exposed through that router is no longer treated as a provisional or merely mirrored route set. The backend now explicitly freezes the currently stabilized authenticated behavior as canonical `v1` semantics.

That freeze covers the authenticated bootstrap, profile, settings, session and wallet inventory/read surface already stabilized in earlier phases. Architecturally, this means:

- canonical `/api/v1/auth/...` routes are the official `v1` transport face of that authenticated surface
- legacy `/auth/...` routes remain supported as compatibility aliases over the same contract
- route versioning does not create independent handler/application semantics
- later authorization work in 0.10 can build on a contract boundary that is both exposed and explicitly frozen

---
#### 0.9.4 Architectural Consequence

With 0.9.4, the versioning model is no longer protected only by documentation and route registration. The transport layer now has representative automated checks proving that selected legacy and canonical entry points continue to expose the same observable contract shape for both standardized error and success scenarios.

Architecturally, this matters because it reduces the risk that future transport work silently separates the two supported route spaces while both still coexist during Stage 0.

#### 0.9.5 Consolidation Consequence

With 0.9.5, Phase 0.9 becomes a closed architectural foundation rather than a partially scattered implementation. The important consequence is documentary but still architectural: the repository now records one coherent answer to all of the following questions at once:

- what the canonical route space is
- why legacy routes still exist
- which authenticated surface is frozen as `v1`
- how representative parity is protected in automated coverage
- why the frontend still remains aligned only up to backend Phase 0.6 until Stage 0 closes

That consolidation matters because Phase 0.10 should build authorization on top of a stable and unambiguous transport boundary rather than on partially divergent status narratives.


### Impact on Future Phases

This foundation enables:

- Phase 0.8 → standardized error model
- Phase 0.9 → versioning strategy
- Phase 0.10 → authorization layer

Without the risk of cross-layer inconsistencies.


## Phase 0.11 — Domain Module Pattern (Deep Architectural Definition)

Phase 0.11 formalizes a repository-level expectation that the current Stage 0 domain-facing modules must stop evolving as mostly handler-centric or historically shaped folders and instead expose a consistent internal boundary model.

### Target Internal Shape

```text
internal/modules/<module>/
    http/
        handlers.go
        dto.go

    app/
        service.go
        usecases.go

    domain/
        model.go
        contracts.go

    repository/
        repository.go   (when it adds real clarity)
```

This is not a claim that every module must contain the same amount of code in each folder or that every module must introduce a repository abstraction regardless of need. The architectural requirement is consistency of responsibility, not cosmetic symmetry.

### Responsibility Split

#### HTTP

The `http` boundary should own:

- request parsing
- transport-level validation
- response mapping
- transport-facing error projection

It should not own authentication flow orchestration, user/business semantics or persistence decisions.

#### APP

The `app` boundary should own:

- use-case orchestration
- ordering of module-internal steps
- coordination with other modules through explicit dependencies

This is where the backend should express things such as authenticated bootstrap composition, settings fetch/update orchestration or authentication entry flows as use cases rather than as heavy handlers.

#### DOMAIN

The `domain` boundary should own:

- models
- invariants
- module-local vocabulary
- explicit contracts that define what the module needs from supporting implementations or collaborating modules

The domain layer must not become a mirror of HTTP DTOs or a storage-schema dump. Its value is to preserve semantic ownership.

#### REPOSITORY

The `repository` boundary should appear only when it improves clarity by isolating persistence-facing implementation details. It is intentionally optional because Phase 0.11 is not a dogmatic rewrite into interface-per-file architecture.

### Dependency Direction

The intended dependency direction is:

`HTTP → APP → DOMAIN`

Repository implementations support those layers but should not pull transport concerns back inward. This matters because the codebase already contains completed transport, error and authorization foundations; Phase 0.11 should build on them rather than collapse them by moving orchestration back into handlers.

### Module-Specific Consequences

#### `internal/modules/user`

The user module becomes the first concrete reference implementation of the pattern. Its handlers should become thinner, while user-related orchestration moves into `app` and user-specific semantic ownership becomes explicit in `domain`.

#### `internal/modules/usersettings`

The user settings module should align to the same pattern without being treated as a mere extension of the user entity. Its domain boundary should preserve settings/configuration semantics such as effective values, consistency and update behavior.

#### `internal/modules/auth`

The auth module requires the most conservative alignment. It already sits on top of stabilized challenge/verify semantics, authenticated identity resolution and authorization-adjacent boundaries. In 0.11 it should be reorganized structurally, not functionally: authentication flows become clearer use cases, but the already delivered public behavior remains unchanged.

0.11.4 completes that conservative alignment in repository state. The final structure is intentionally staged:

- `auth/app` is the runtime/application owner for application orchestration, service behavior, wallet services, response models and reusable composition helpers
- `auth/domain` is the canonical owner of wallet contracts and base wallet types
- root `auth` remains a compatibility and transport surface where existing callers, sentinels and HTTP handlers still depend on package-level access
- `auth/repository` is a transitional façade over root store implementations rather than the active implementation owner
- root wallet store files remain the accepted runtime implementations until a dedicated repository migration phase moves them safely

This preserves the public wallet bootstrap auth flow, authenticated wallet-management behavior, session behavior and error semantics while reducing the pre-0.11 handler/root ownership model.

### Cross-Module Contract Boundaries

One of the deep architectural goals of 0.11 is to reduce concrete cross-module knowledge between `auth`, `user` and `usersettings`.

The important distinction is:

- coordination is allowed
- ownership transfer is not

That means:

- `auth` may coordinate bootstrap-facing reads without owning the `user` or `usersettings` domains
- `user` must not absorb authentication semantics
- `usersettings` must not become informal profile storage

Where module interaction is real, it should be expressed through minimal explicit contracts rather than through direct knowledge of another module's internal models or implementation details.

0.11.5 turns this principle into completed implementation state. The subphase was intentionally split into:

- 0.11.5.0 — Subphase Definition & Documentation Alignment ✔ Completed
- 0.11.5.1 — Dependency Mapping ✔ Completed
- 0.11.5.2 — Contract Extraction ✔ Completed
- 0.11.5.3 — Interface Alignment ✔ Completed
- 0.11.5.4 — Runtime Compatibility Validation ✔ Completed

The delivered result is that `auth/app` coordinates with `user` and `usersettings` through explicit minimal contracts, while root compatibility layers preserve existing service wiring and externally observable behavior.

The expected result is not a broad abstraction rewrite. It is a controlled consolidation of the real contracts needed for `auth`, `user` and `usersettings` to collaborate without reintroducing concrete ownership confusion.

### Semantic Ownership

The pattern is reinforced by explicit ownership boundaries:

- `auth` owns authentication flows, challenge/verify semantics, and authenticated identity entry
- `user` owns the platform user model and user-facing metadata/profile surface
- `usersettings` owns authenticated user configuration, preferences, and settings consistency

Modules may coordinate with each other, but they must not silently take ownership of another module's domain.

### Architectural Non-Goals

Phase 0.11 does not claim that the backend is adopting a full clean architecture, CQRS/event-sourcing model or generic cross-cutting module framework. The deep architectural purpose is narrower and more pragmatic:

- normalize the current module set
- preserve Stage 0 external stability
- reduce accidental coupling
- create clearer foundations for later work

### Relationship to Earlier Stage 0 Work

Phase 0.11 depends on the earlier Stage 0 foundations already being in place:

- the authenticated application surface consolidated in 0.6
- the application-layer intent introduced in 0.7
- the standardized error boundary introduced in 0.8
- the canonical route/versioning discipline introduced in 0.9
- the explicit authorization boundary completed in 0.10

Because those foundations already exist, 0.11 can stay strictly structural. It does not need to reopen public-contract decisions in order to improve module clarity.

### Current Repository State

0.11.1 is completed and has already established the deep architectural definition of the Domain Module Pattern.

0.11.2 is also completed in runtime repository state. `internal/modules/user` is the first concrete module aligned to this pattern, which means the repository has already moved from pure definition into controlled execution.

0.11.3 is also completed in runtime repository state. `internal/modules/usersettings` now follows the same structural pattern with explicit `app`, `domain` and `repository` boundaries while preserving settings-specific orchestration and external contract behavior.

0.11.4 is also completed in runtime repository state. `internal/modules/auth` now follows the same pattern conservatively, with `auth/app` as runtime/application owner, `auth/domain` as canonical contract owner, root `auth` as compatibility/transport surface and `auth/repository` as a documented transitional façade.

Phase 0.11 is therefore **completed**. Cross-module contracts have been consolidated and the trunk documentation now records the completed phase state coherently.

---

## Phase 0.12 — Read / Write Model Separation (Deep Architectural Definition)

Phase 0.12 introduces a narrower model-responsibility discipline on top of the already completed module-boundary discipline from 0.11.

The key distinction is not storage architecture. The backend is not moving to full CQRS. Instead, the distinction is about preventing the same structure from silently representing several different concepts at once.

### Model Categories

0.12.1 must classify real structures using the following categories:

- **READ**: output/view/response-oriented shape
- **WRITE**: input/command/mutation-oriented shape
- **DOMAIN**: module-owned semantic model or invariant carrier
- **CONTRACT**: explicit interface or minimal cross-module contract
- **INFRASTRUCTURE**: persistence, transport, token, config or framework support shape
- **HYBRID / TRANSITIONAL**: structure currently carrying more than one responsibility

The classification is a documentation and design tool. It must be based on actual repository usage.

### Mapping Discipline

The target direction for later 0.12 subphases is explicit transformation:

```text
Write Model -> Domain / Application Use Case
Domain / Application Result -> Read Model
```

Mappings may be local to the module while the module set is still small. The phase does not require a global shared mapper package unless the real code proves that such extraction is useful.

### Non-Goals

Phase 0.12 does not introduce:

- independent read and write databases
- event streams
- projection rebuilds
- public API payload migration
- new versioned route spaces
- frontend migration requirements

### Deep Compatibility Constraint

The separation must be invisible to current clients. Any extracted read or write model must preserve the public shape currently emitted or accepted by the relevant handler. Where names change internally, mapping must absorb the change.

### 0.12.0 Result

0.12.0 completes the documentary definition for this discipline. Code changes begin only after the 0.12.1 audit identifies real models and real hybrid points.


### 0.12.1.0 — Audit-Level Documentation Lock

The Model Classification Audit is itself split into controlled internal steps because an incomplete audit would make later model extraction unsafe. 0.12.1.0 therefore locks the audit process before any inventory output is generated.

The required audit sequence is:

1. extract the complete model inventory from real code
2. classify each structure using the Phase 0.12 vocabulary
3. analyze cross-layer usage
4. identify hybrid/transitional risk
5. define the target separation direction
6. consolidate the audit before implementation begins

The deep architectural constraint is that classification is descriptive first. A structure marked HYBRID / TRANSITIONAL is not automatically changed during the audit. It becomes a documented candidate for later read/write/domain separation with explicit mapping.

This preserves the compatibility rule that public payloads, route behavior, authentication behavior, authorization behavior, error envelopes and API versioning remain unchanged while the audit is performed.

### 0.12.1 Deep Audit Result

0.12.1 completed the deep model audit required before Read / Write Model Separation can touch code.

Audit outputs now present in the repository:

- full model inventory
- model classification report
- cross-layer usage analysis
- problem detection and risk mapping
- target separation definition
- audit consolidation and closure

The audit baseline records:

- Total model-like structs: 125
- Hybrid / transitional structures: 11
- Classification groups: READ, WRITE, HYBRID and INTERNAL
- Risk levels used: CRITICAL, HIGH, MEDIUM and LOW

This result does not itself change architecture at runtime. It constrains the next implementation phases so extraction work is based on documented evidence rather than assumption.

### 0.12.2.0 — Read Model Extraction Definition Lock

Read Model Extraction is the first Phase 0.12 subphase that prepares for code-level model separation. The `.0` step is intentionally documentation-only so the extraction path is constrained before implementation begins.

The deep architectural rule for 0.12.2 is that read models represent output intent. They may serve response assembly, view-style application results or transport-facing response shapes, but they must not become mutation/input models and must not absorb domain invariants that belong to module-owned domain structures.

The implementation sequence is locked as:

1. define read model design targets from the audit evidence
2. introduce explicit read model structures where required
3. introduce Domain/Application → Read mapping where the flow is currently implicit
4. align existing responses without changing public payload semantics
5. validate runtime compatibility and `go test ./...`
6. close documentation cumulatively

The 0.12.1 audit remains the evidence source. Any read model extraction that cannot be connected to the inventory/classification/cross-layer/risk/target documents must be deferred instead of guessed.


### 0.12.2.6 — Read Model Extraction Closure

0.12.2 closes with the read-side boundary implemented and validated. The extracted read models are projection-oriented structures used to assemble responses and application-facing output without reusing mutation/input structures.

The implemented direction is Domain/Application → Read. Mapping is explicit and additive, which keeps public response compatibility stable while making ownership clearer inside `auth`, `user` and `usersettings`. The next structural risk is write-side reuse, which is reserved for 0.12.3 — Write Model Isolation.

### 0.12.3.0 — Write Model Isolation Definition Lock

Write Model Isolation is constrained by the same compatibility discipline used for read model extraction. The goal is to prevent mutation/input flows from depending on read-oriented response structures or domain-owned objects while keeping accepted request behavior stable.

The implementation sequence is locked as:

1. define write model design targets from audit evidence and current handlers
2. introduce explicit write model structures where required
3. introduce Write → Domain/Application mapping where input flow is currently implicit
4. align handlers without changing public request semantics
5. validate runtime compatibility and `go test ./...`
6. close documentation cumulatively

The 0.12.2 read model packages remain read-only projection structures. They must not become input models during 0.12.3.

### 0.12.3.6 — Write Model Isolation Closure

Phase 0.12.3 completes the write-side counterpart to the read model extraction. The repository now includes `writemodels` packages for `auth`, `user` and `usersettings`, domain write input structures and explicit Write → Domain mappers.

Handler alignment preserves public request payload semantics while introducing an internal input boundary. Read models remain output-only and are not reused as command/input structures. Validation for the completed sequence is recorded in `docs/phase0_12_3_validation_compatibility.md` and confirmed by `go test ./...` in the development environment.

### 0.12.4.0 — Mapping Layer Introduction Definition Lock

After 0.12.2 and 0.12.3, the repository has explicit model separation but mapper ownership is still distributed: read model packages own output projection helpers, write model packages own input-to-domain helpers and some application support code still participates in response assembly.

0.12.4 defines the consolidation path. The target is a dedicated module-local mapping package:

```text
internal/modules/<module>/mappers/
```

This package becomes the owner of explicit boundary transformations while model packages remain structural model definitions. The intended directions are:

- Write → Domain for command/input translation
- Domain/Application → Read for response projection

The `.0` step is documentation-only. It prevents implementation from starting before the ownership rule, compatibility constraints and sub-subphase sequence are locked.


### 0.12.4.6 — Mapping Layer Closure

The mapping layer introduction sequence is completed. The backend now has centralized module-local packages for explicit transformation directions: Write → Domain and Domain → Read. Application-layer helper code delegates to these packages where 0.12.4 targeted consolidation, while compatibility wrappers remain only where they preserve stable public behavior or keep staged refactor risk low.

Validation status: `go test ./...` passed in the developer environment after 0.12.4.5 and remained valid for the closure baseline.

## Phase 0.12.5 Deep Contract Alignment Model

Phase 0.12.5 is the bridge between the model separation work and the contracts introduced during the Domain Module Pattern phase. The purpose is not to redesign business behavior. The purpose is to verify and align internal contracts so that providers, services and application boundaries describe their read/write intent explicitly.

The alignment work must inspect:

- contract interfaces introduced during Phase 0.11
- provider implementations used by auth, user and usersettings
- read-side return values that should remain response-compatible
- write-side inputs that should not leak transport or read-model concerns
- mapping call sites that should remain owned by module-local mapper packages

The phase must preserve compatibility while making contract intent explicit enough for later hardening.

## Phase 0.12.5 Contract Alignment Closure

The contract alignment sequence confirms that handler-facing request contracts are centralized, read/write transformation ownership remains in module mapper packages and provider-facing boundaries no longer obscure model direction. This completes the Phase 0.12 contract alignment path without requiring additional runtime contract drift.
