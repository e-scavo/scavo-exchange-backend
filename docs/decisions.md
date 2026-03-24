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