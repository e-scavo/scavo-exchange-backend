# Architectural Decisions

## ADR-0001 - Modular Monolith

### Decision
The backend architecture is a modular monolith.

### Reason
The current project is still at an early stage and needs fast iteration, simple deployment, low operational complexity, and strong alignment between implementation and documentation.

### Impact
- single deployable backend
- easier local development
- simpler testing
- easier refactoring during early product discovery
- future service extraction remains possible if needed

---

## ADR-0002 - DEX-First Product Strategy

### Decision
The backend will prioritize DEX functionality before custodial or CEX features.

### Reason
The current exchange focus is explicitly DEX-oriented, aligned with SCAVIUM and wallet-native interaction patterns.

### Impact
- self-custody comes first
- on-chain settlement is primary
- wallet integration is central
- custodial ledger is deferred
- order-book trading is deferred

---

## ADR-0003 - AMM v1 as Initial DEX Model

### Decision
The first DEX implementation will use an AMM v1 model rather than an order book.

### Reason
AMM is simpler, faster to deliver, more consistent with the initial non-custodial scope, and more realistic for a first end-to-end implementation.

### Impact
- pools become primary liquidity source
- quotes are derived from pool state
- routing starts simple and expands later
- smart contracts are more straightforward
- backend responsibility is easier to define

---

## ADR-0004 - SCAVIUM as Primary Chain

### Decision
SCAVIUM is the primary blockchain network for the backend.

### Reason
The exchange is being built around the SCAVIUM ecosystem and must integrate directly with its chain, contracts, RPC nodes, and future exchange-specific infrastructure.

### Impact
- chain abstraction must still exist
- initial implementations can optimize for SCAVIUM
- dedicated RPC capacity may be introduced when needed

---

## ADR-0005 - Self-Custody First

### Decision
DEX users remain self-custodial in the initial product model.

### Reason
This reduces operational and regulatory complexity while matching the DEX-first architecture.

### Impact
- backend does not initially hold user private keys
- user wallet signs transactions
- backend provides reads, quotes, routing help, and tracking
- future custodial features must remain explicitly separated

---

## ADR-0006 - PostgreSQL as Primary Database Target

### Decision
PostgreSQL is the target primary relational database.

### Reason
The project will require reliable relational modeling, transaction support, operational consistency, and flexible querying for blockchain-indexed data and future ledger features.

### Impact
- repository design should target PostgreSQL compatibility
- migrations will become a formal part of the project in Phase 0.2
- JSONB can be used where appropriate without abandoning relational discipline

---

## ADR-0007 - Redis as Secondary Infrastructure Store

### Decision
Redis will be used for cache and operational coordination.

### Reason
The project will need short-lived state, cache, rate-limit counters, and lightweight coordination beyond the primary relational database.

### Impact
- cache boundaries must remain explicit
- Redis is not the system of record
- persistence-critical data stays in PostgreSQL

---

## ADR-0008 - REST and WebSocket as First-Class Interfaces

### Decision
The backend will expose both REST and WebSocket interfaces as official first-class transports.

### Reason
The product needs both request/response behavior and real-time interaction.

### Impact
- transport contracts must be documented separately
- module services should remain transport-agnostic
- real-time state updates can grow without redesigning the application

---

## ADR-0009 - No Matching Engine in Initial Scope

### Decision
Matching engine and order-book trading are out of scope for the first product phases.

### Reason
They add major domain, operational, reconciliation, and custodial complexity that is not required for the initial DEX-first objective.

### Impact
- no internal execution engine in early phases
- no internal balance matching in early phases
- hybrid expansion remains possible later

---

## ADR-0010 - Migrations as Schema Source of Truth

### Decision
Database schema evolution will be managed through versioned migrations.

### Reason
The project needs reproducible environments, auditable schema changes, and a safe way to evolve persistence without undocumented drift.

### Impact
- schema changes must be versioned
- manual environment-only schema drift is discouraged
- repository evolution must remain aligned with migration history

---

## ADR-0011 - Platform Adapters for Chain and Contract Integration

### Decision
Chain-specific and contract-specific integrations should be isolated into dedicated platform packages or adapters.

### Reason
This reduces low-level protocol leakage into domain services and keeps blockchain integrations testable, replaceable, and easier to evolve.

### Impact
- domain services should not directly embed raw RPC logic everywhere
- chain integrations can evolve independently
- smart contract wiring remains more maintainable

---

## ADR-0012 - Infrastructure Before Feature Expansion

### Decision
Persistence and infrastructure baseline work must be completed before major product feature implementation begins.

### Reason
Jumping directly into DEX features, contracts, or chain-heavy logic before core infrastructure is stable would create architectural drift and rework.

### Impact
- Phase 0.2 remains mandatory
- major DEX features are deferred until infrastructure is ready
- implementation order remains disciplined and safer

---

## ADR-0013 - PostgreSQL as Durable Source of Truth

### Decision
Durable backend state must live in PostgreSQL unless there is a strong reason otherwise.

### Reason
The project will need durable operational state, reproducible environments, auditable records, and survivable data across restarts and deployments.

### Impact
- durable module state must not be Redis-only
- long-lived operational truth belongs in relational persistence
- Redis usage must remain clearly secondary

---

## ADR-0014 - Redis for Ephemeral and Coordination State

### Decision
Redis usage is limited to ephemeral, cache, or coordination-oriented concerns.

### Reason
This keeps the persistence model clean and avoids accidental business dependence on non-authoritative state.

### Impact
- Redis is optional for some flows
- cached values must remain reconstructable
- critical business truth should not depend on Redis persistence semantics

---

## ADR-0015 - Environment-Driven Local Infrastructure Baseline

### Decision
The project must support a reproducible local environment driven by explicit configuration.

### Reason
The backend will require DB, cache, and chain-oriented setup. Reproducible local environments reduce onboarding friction and prevent hidden configuration drift.

### Impact
- environment variables become part of the formal project contract
- Docker-based local infrastructure is recommended
- local development setup must be documented and repeatable

---

## ADR-0016 - Health and Readiness Must Be Explicitly Separated

### Decision
The backend must distinguish liveness from readiness as infrastructure and integrations are introduced.

### Reason
A running process is not necessarily capable of serving intended workloads once the application depends on DB, cache, chain RPC, migrations, and background coordination.

### Impact
- health endpoints and readiness checks should evolve separately
- operational diagnostics become clearer
- rollout and internal testing behavior become safer

---

## ADR-0017 - Observability Is a First-Class Infrastructure Concern

### Decision
Observability must be treated as part of the core backend foundation rather than an optional later add-on.

### Reason
The system will depend on blockchain integrations, persistence, background processing, and hybrid growth paths that are difficult to diagnose without structured visibility.

### Impact
- logging, correlation, and metrics are part of the architecture
- new infrastructure work should remain diagnosable
- failures should become visible earlier in development

---

## ADR-0018 - Testing Must Grow with the Architecture

### Decision
Testing strategy must evolve incrementally alongside infrastructure and domain growth.

### Reason
Delaying test structure until late stages would create fragile integrations and make regression control much harder once DB, chain, contracts, and real-time behavior are introduced.

### Impact
- test layers are defined early
- infrastructure work should become testable as introduced
- critical flows should gain regression protection as they stabilize

---

## ADR-0019 - Path-Based API Versioning Becomes Canonical After Phase 0.8

### Decision
The backend adopts path-based API versioning as the canonical public-contract strategy, with `/api/v1/...` defined as the authoritative route namespace for the currently stabilized authenticated surface.

### Reason
By the end of Phase 0.8, the backend already exposes a stable transport contract and a standardized error envelope, but it still does so through unversioned routes. Future work such as authorization and later domain growth needs an explicit contract-evolution model so compatibility and breaking changes are handled intentionally rather than informally.

### Impact
- the current unversioned routes remain available as legacy compatibility routes
- `v1` freezes the present success-payload and error-envelope semantics
- future breaking transport changes require a new API version
- versioning stays in the transport layer and must not leak into application orchestration
- frontend alignment can remain on the pre-versioned route set until Stage 0 completion without blocking backend evolution


## Decision — Authorization Model Before Enforcement

The backend introduces authorization in two explicit steps rather than mixing model definition with immediate endpoint blocking.

First, Phase 0.10.1 defines roles, permissions, a static role-to-permission mapping and an authorization-subject model in a dedicated core package. Later subphases will propagate that model through middleware, evaluate policies centrally and only then enforce permissions on selected endpoints.

This decision preserves Stage 0 stability by preventing ad-hoc authorization checks from appearing first in handlers or transport code without an agreed vocabulary. It also keeps authorization conceptually separate from authentication claims and from API versioning concerns.


## Decision — Authorization Context Must Be Hydrated Before Policy Evaluation

The backend must project authenticated identities into a dedicated authorization context before introducing centralized policy checks.

This means policy evaluation will consume a normalized `AuthorizationSubject` from request context rather than reading raw JWT claims directly inside handlers or policy code.

### Reason

Keeping policy code dependent on raw transport claims would blur the boundary between authentication and authorization, increase coupling to JWT structure and encourage handler-local access decisions.

### Impact

- authenticated routes now hydrate authorization context after successful authentication
- future policy evaluation can remain transport-agnostic
- endpoint enforcement can be introduced later without redesigning the authenticated request pipeline


## Decision — Policy Evaluation Must Be Centralized Before Endpoint Enforcement

Authorization decisions must be introduced through a dedicated core evaluator before protected endpoints start denying requests based on permissions.

Phase 0.10.3 therefore introduces explicit authorization action/resource vocabularies and a centralized `PolicyEvaluator` rather than allowing handlers or router-level code to translate roles into permissions independently.

### Reason

If endpoint enforcement were introduced first, authorization semantics would likely fragment across handlers and middleware, making later changes riskier and harder to test.

### Impact

- authorization questions now have a stable core API
- handler and middleware code can stay thin when enforcement starts
- future permission expansion can evolve in one place rather than across multiple endpoint implementations


## Decision — Endpoint Enforcement Must Start Only Where the Static Permission Model Already Fits

Phase 0.10.4 should begin operational authorization enforcement only on authenticated endpoints whose semantics are already represented by the current static action/resource and role-permission model.

### Reason

The current static authorization model can cleanly represent `read user`, `read settings` and `update settings`, but it does not yet fully explain self-profile mutation semantics for `PATCH /auth/me` without either broadening `RoleUser`, adding owner-aware rules or introducing richer resource scoping.

Starting enforcement only on the already-aligned endpoints proves the full authorization path without forcing an early policy decision that could break the stabilized Stage 0 application surface.

### Impact

- authorization is now enforced in real router flows rather than remaining only preparatory
- route-level denial still stays centralized through policy middleware
- endpoints requiring richer ownership semantics remain stable until the policy model is expanded intentionally


## Decision — Close Authorization Phase Through Documentation Consolidation

Once Phase 0.10.4 introduced real endpoint-level enforcement, the repository could no longer treat authorization as only a future concern in trunk documents. Phase 0.10.5 therefore exists to make the contract and implementation narrative converge before the next phase begins.

This means the repository now explicitly documents two truths at the same time:

- authorization rollout remains progressive and intentionally limited to the endpoints already represented by the static model
- authorization is already operational on that first slice, with a standardized forbidden path for denied requests

Documenting both points together is important because it prevents later work from incorrectly assuming either that authorization is still pending everywhere or that every authenticated endpoint is already under full policy control.

---
## ADR-0011 - Domain Module Pattern

### Decision
The backend adopts a Domain Module Pattern to standardize the internal structure of domain modules.

### Reason
The backend has reached a stage where authentication, authorization, API versioning, standardized error handling, and authenticated bootstrap flows are already stabilized. The next structural need is to remove internal inconsistencies between modules, formalize separation of concerns, and reduce accidental coupling before broader feature growth.

### Impact
- domain modules are organized around explicit internal layers
- HTTP, application orchestration, and domain semantics are separated
- repository concerns remain optional and only appear where they add real clarity
- module-to-module dependencies must be expressed through minimal explicit contracts when needed
- the pattern applies conservatively and does not imply a full clean-architecture rewrite
- external contracts remain unchanged while internal maintainability improves

## Decision — Auth Repository Runtime Compatibility During 0.11.4

### Decision

During 0.11.4, `internal/modules/auth/domain` became the canonical owner of wallet contracts and base wallet types, while the root `auth` wallet store implementations remain the active runtime implementations.

`internal/modules/auth/repository` is intentionally a transitional façade in this phase. It is not yet the canonical owner of the concrete PostgreSQL or memory store implementations.

### Reason

The wallet identity and wallet challenge stores are large, active runtime implementations used by the current auth and wallet flows. Moving them into `auth/repository` during 0.11.4 would create a high-risk migration unrelated to the immediate goal of aligning application, domain and HTTP ownership.

0.11.4 therefore locks the following staged decision:

- domain contracts are canonical now
- app/application/service ownership is canonical now
- root stores remain runtime-compatible implementations now
- repository migration is deferred to a dedicated future phase

### Impact

This preserves wallet bootstrap, wallet verify, wallet management, session and testing behavior while making the transitional state explicit. Future work must not duplicate or move the root store implementations without a dedicated repository migration plan.

### Status

Accepted in 0.11.4C4.2.

---

## Decision — Cross-Module Coordination Must Use Explicit Minimal Contracts During 0.11.5

### Decision

0.11.5 will consolidate cross-module coordination between `auth`, `user` and `usersettings` through explicit minimal contracts where coordination is real.

The objective is not to introduce a broad shared-kernel package or to move ownership between modules. Each module keeps its completed 0.11 ownership boundary:

- `auth` owns authentication, session, wallet bootstrap and authenticated-entry behavior
- `user` owns the user entity/profile boundary
- `usersettings` owns settings and preference semantics

### Reason

After 0.11.2, 0.11.3 and 0.11.4, the three modules have explicit internal structure. The remaining risk is not folder layout; it is accidental cross-module knowledge through concrete services, models or implementation details.

0.11.5 exists to map those dependencies, extract only the contracts that are actually needed, align call sites to those interfaces and validate runtime compatibility.

### Impact

- cross-module dependencies become intentional and reviewable
- module ownership remains stable
- public HTTP behavior remains unchanged
- later phases can evolve modules independently with less coupling pressure

### Status

Accepted as the execution rule for 0.11.5.0 through 0.11.5.4 and confirmed during 0.11.6 phase closure.



## Decision — Phase 0.11 Closure State

### Context

Phase 0.11 defined and applied the Domain Module Pattern across the current Stage 0 domain-facing modules. It also consolidated real cross-module dependencies between `auth`, `user` and `usersettings` through explicit minimal contracts.

### Decision

Phase 0.11 is closed after 0.11.6. The repository treats the Domain Module Pattern as delivered for the current module set, with `auth`, `user` and `usersettings` aligned to explicit internal boundaries and cross-module coordination constrained by minimal contracts where required.

### Consequences

- Future changes should not reintroduce concrete service coupling into `auth/app`.
- Root compatibility layers may continue to preserve public behavior and existing wiring.
- Further repository migration or read/write model separation should occur in later roadmap phases rather than being folded back into 0.11.

### Status

Accepted in 0.11.6.

---

## ADR-0012 - Read / Write Model Separation Is Internal And Compatibility-Preserving

### Context

After Phase 0.11, the backend has explicit module boundaries for `auth`, `user` and `usersettings`. The remaining architectural ambiguity is that some model structures can still represent input, output, domain state or application results depending on the call path.

### Decision

Phase 0.12 will separate read-oriented and write-oriented model responsibilities internally while preserving the existing public HTTP contract.

The accepted direction is:

- read models are response/view-oriented
- write models are input/command-oriented
- domain models remain module-owned
- mapping functions make transformations explicit

### Non-Decision

This is not an adoption of full CQRS or event sourcing. Phase 0.12 does not introduce separate persistence models, projections, asynchronous read stores or new API versions.

### Consequences

- Model intent becomes easier to review.
- Future module changes can avoid reusing response shapes as mutation inputs.
- The 0.11 contracts can evolve with clearer input/output boundaries.
- Public clients should observe no behavioral change.

### Status

Accepted in 0.12.0.


---

## ADR-0013 - Subdivided Subphases Must Start With A Documentation Lock

### Context

Phase 0.12.1 requires a detailed model audit before code changes. Some subphases may need to be subdivided into smaller execution units to avoid mixing definition, analysis and implementation.

Without a mandatory `.0` lock, a subdivided subphase could begin producing implementation or audit output before the trunk documentation records its internal scope and order.

### Decision

When a subphase is subdivided, the first internal step must be a `.0` documentation lock.

For Phase 0.12.1 this is:

- 0.12.1.0 — Definition & Documentation Lock

The `.0` step must update trunk documentation as applicable, register all sub-subphases and avoid code changes.

### Consequences

- Subdivided work remains traceable.
- Trunk documentation stays ahead of implementation.
- Audit or implementation work cannot silently expand without being recorded first.
- Documentation-only locks remain explicit and reviewable.

### Status

Accepted in 0.12.1.0.

---

## ADR-0014 - Model Audit Must Precede Read / Write Extraction

### Context

Phase 0.12 separates read-oriented, write-oriented and domain-owned model responsibilities. Performing extraction directly without an exhaustive audit would risk splitting the wrong structures, missing cross-layer coupling or changing public behavior accidentally.

### Decision

Complete 0.12.1 as an audit-only sequence before any Go code changes for Read / Write Model Separation.

The audit sequence includes:

- model inventory extraction
- model classification
- cross-layer usage analysis
- problem detection and risk mapping
- target separation definition
- audit consolidation and closure

### Consequences

- 0.12.2 can start from documented model evidence.
- Hybrid/transitional structures are identified before extraction.
- Public contract preservation remains explicit.
- Code changes are deferred until the audit is closed.

### Status

Accepted in 0.12.1.6.

---

## ADR-0015 - Read Model Extraction Must Preserve Public Response Compatibility

### Context

Phase 0.12.1 identified hybrid/transitional structures and defined target separation direction. Phase 0.12.2 begins the read-side extraction path. Without an explicit compatibility rule, introducing read models could accidentally change public payloads or move response semantics away from existing handlers.

### Decision

0.12.2 must introduce read models only as internal structural improvements. Existing public routes, status codes, response envelopes, authentication behavior, authorization behavior, error model behavior and API versioning remain unchanged.

The implementation must be driven by the 0.12.1 audit artifacts and must proceed through the locked 0.12.2 internal sequence:

- 0.12.2.0 — Definition & Documentation Lock
- 0.12.2.1 — Read Model Design
- 0.12.2.2 — Read Model Implementation
- 0.12.2.3 — Domain/Application → Read Mapping
- 0.12.2.4 — Response Alignment
- 0.12.2.5 — Validation & Compatibility
- 0.12.2.6 — Documentation & Closure

### Consequences

- Read-side extraction becomes traceable.
- Public compatibility remains the controlling constraint.
- Mapping must be explicit where read models are introduced.
- Any unsupported or speculative extraction target must be deferred.

### Status

Accepted in 0.12.2.0.


## Decision: Close Read Model Extraction With Internal Compatibility Preservation

### Context

Phase 0.12.2 introduced read models after the 0.12.1 audit identified hybrid/transitional structures and target separation needs. The implementation had to improve internal structure without changing public HTTP behavior.

### Decision

Close 0.12.2 with explicit read model packages, explicit Domain/Application → Read mapping functions, response alignment and compatibility validation. Public routes, response envelopes, JSON tags, authentication behavior, authorization behavior, standardized errors and API versioning remain unchanged.

### Consequences

- read-side output responsibilities are now explicit
- mapping is no longer implicit in the aligned paths
- handlers keep public compatibility
- write-side isolation remains deferred to 0.12.3

### Status

Accepted in 0.12.2.6.

---

## ADR-0016 - Write Model Isolation Must Preserve Public Request Compatibility

### Context

Phase 0.12.2 completed read model extraction and introduced explicit read-side projection structures. Phase 0.12.3 isolates the write side. Without a documentation lock, write model extraction could accidentally reuse read models as inputs, change accepted payload semantics or blur the new read boundary.

### Decision

0.12.3 must introduce write models only as internal structural improvements. Existing public routes, request payload semantics, authentication behavior, authorization behavior, error model behavior and API versioning remain unchanged.

The implementation must proceed through the locked 0.12.3 internal sequence:

- 0.12.3.0 — Definition & Documentation Lock
- 0.12.3.1 — Write Model Design
- 0.12.3.2 — Write Model Implementation
- 0.12.3.3 — Write → Domain Mapping
- 0.12.3.4 — Handler Alignment
- 0.12.3.5 — Validation & Compatibility
- 0.12.3.6 — Documentation & Closure

### Consequences

- Write-side extraction becomes traceable.
- Read models remain output-only and must not be reused as inputs.
- Public request compatibility remains the controlling constraint.
- Mapping must be explicit where write models are introduced.

### Status

Accepted in 0.12.3.0.

## Decision: Close Write Model Isolation Boundary

### Context

Phase 0.12.2 introduced explicit read models and established that read projections must remain output-only. Phase 0.12.3 introduced the write-side counterpart so mutation and command inputs do not reuse read projections or domain-owned structures directly.

### Decision

Close 0.12.3 with explicit write model packages, domain write input structures, Write → Domain mapping functions and handler alignment that preserves public request payload semantics.

### Consequences

- read models remain output-only
- write models represent mutation/input intent
- handlers can preserve existing HTTP contracts while routing through explicit write-side structures
- future mapping consolidation can build on both Domain → Read and Write → Domain directions

### Status

Accepted in 0.12.3.6.

## Decision: Centralize Mapping Ownership During Phase 0.12.4

### Context

Phase 0.12.2 introduced explicit read models and Domain/Application → Read mapping. Phase 0.12.3 introduced explicit write models, domain write inputs and Write → Domain mapping. This produced correct separation but left mapper ownership distributed across read model packages, write model packages and application support code.

### Decision

Phase 0.12.4 introduces a centralized module-local mapping layer under:

```text
internal/modules/<module>/mappers/
```

The layer will own explicit transformations while preserving all current public HTTP contracts.

The locked internal sequence is:

- 0.12.4.0 — Definition & Documentation Lock
- 0.12.4.1 — Mapping Layer Design
- 0.12.4.2 — Mapping Layer Implementation
- 0.12.4.3 — Mapping Consolidation
- 0.12.4.4 — Application Refactor
- 0.12.4.5 — Validation & Compatibility
- 0.12.4.6 — Documentation & Closure

### Consequences

- model packages remain focused on model definitions
- mapping ownership becomes explicit
- application code should stop accumulating transformation logic
- public endpoints and JSON contracts remain unchanged
- implementation must proceed incrementally and validate with `go test ./...`

### Status

Accepted in 0.12.4.0. Completed in 0.12.4.6.


## Decision: Close Phase 0.12.4 Mapping Layer Introduction

### Context

Phase 0.12.4 implemented the mapping ownership target defined after read model extraction and write model isolation. Mapping logic now has explicit module-local ownership instead of remaining distributed across read model, write model and application packages.

### Decision

Close Phase 0.12.4 as completed. The backend will proceed to Phase 0.12.5 — Contract Alignment using the centralized mapper packages as the baseline.

### Consequences

- Public HTTP request and response contracts remain stable.
- Read and write model packages remain focused on data shape ownership.
- Mapping packages own transformation logic.
- Future contract alignment must use the centralized mapping layer rather than adding new mapping logic to application handlers.

### Status

Accepted in 0.12.4.6.

## Decision: Start Contract Alignment After Mapping Layer Consolidation

Status: Accepted  
Phase: 0.12.5.0

### Context

Phase 0.12.2 introduced read models, Phase 0.12.3 introduced write models and Phase 0.12.4 centralized mapper ownership. The remaining structural risk is that provider contracts may still describe mixed model ownership or allow implicit transformations to remain hidden behind application boundaries.

### Decision

Start Phase 0.12.5 as a dedicated Contract Alignment sequence. The sequence must review internal provider contracts and align them with the established model directions:

- Write → Domain
- Domain/Application → Read
- centralized mapper ownership
- public HTTP/API compatibility preservation

### Consequences

Contract changes must be evidence-driven and incremental. Any implementation step must avoid public API drift, avoid endpoint changes and avoid weakening the compatibility guarantees already validated in Phase 0.12.2 through Phase 0.12.4.

## Decision: Close Contract Alignment With Public Contract Preservation

Phase: 0.12.5.6

Context:
Phase 0.12.5 introduced contract inventory, normalization design, implementation-level alignment, handler contract adjustment and compatibility validation.

Decision:
Close Contract Alignment as completed while preserving existing public HTTP routes, JSON request fields, JSON response fields and API versioning. Internal contract ownership is now clearer through centralized HTTP contract aliases and mapper-owned transformations.

Consequences:
Phase 0.12 can close with model separation, mapping ownership and contract alignment consistently documented and validated.

## Decision: Start Provider Layer Consolidation After Contract Alignment

Phase: 0.13.0

Context:
Phase 0.12 completed Read / Write Model Separation, centralized mapping ownership and internal contract alignment. Provider-like responsibilities still require a dedicated consolidation pass so that handlers and application flows rely on explicit provider boundaries.

Decision:
Start Phase 0.13 as Provider Layer Consolidation. The phase begins with documentation lock, then inventory, interface design, implementation, application integration, validation and closure.

Consequences:
Provider work must remain internal and compatibility-preserving. Public HTTP routes, request payloads, response payloads, API versioning and business behavior must remain unchanged unless a later explicitly approved phase changes them.
