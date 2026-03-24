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