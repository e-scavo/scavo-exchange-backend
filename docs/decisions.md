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
