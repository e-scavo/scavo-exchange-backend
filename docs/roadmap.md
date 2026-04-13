# 🧭 Roadmap — SCAVO Exchange Backend

## 🏗️ Stage 0 — Foundation

### Phase 0.1 — Baseline & Documentation
- Audit current project
- Define architecture
- Create documentation
- Establish roadmap

### Phase 0.2 — Core Infrastructure
- Database
- Redis
- Config system
- Observability
- Testing base

---

### Phase 0.5 — User Interaction & Application Surface
- authenticated user profile surface
- minimal non-wallet user metadata update contract
- user-facing read models for application bootstrap
- lightweight application-facing identity summary

---

## Phase 0.6 — Authenticated Application Bootstrap Consolidation

### 0.6.1 — Bootstrap Surface Boundary Clarification ✔
Completed

### 0.6.2 — Authenticated Surface Contract Alignment ✔
Completed

This subphase aligns all authenticated endpoints under a shared normalized context, ensuring contract-level consistency without breaking existing APIs.

### 0.6.3 — Session-Ready Bootstrap Read Model ✔
Completed

This subphase introduces a unified bootstrap endpoint that aggregates all authenticated surfaces into a single response, enabling a one-call frontend initialization flow.

### 0.6.4 — Application Surface Consistency Hardening ✔
Completed

This subphase finalizes the authenticated application surface by enforcing a consistent structural contract across all endpoints, introducing canonical envelopes while preserving backward compatibility.

---

### Phase 0.6 — Completed

Phase 0.6 is now fully completed.

It delivered:

- explicit authenticated surface boundaries
- aligned shared authenticated context
- unified bootstrap read model
- structural contract hardening across endpoints

### Next

Phase 0.7 — Application Layer Foundation

---

## 🧱 Stage 0 — Foundation (Remaining Phases)

### 0.7 — Application Layer Foundation ✔ Completed
Introduce an application layer to separate HTTP from business logic.

#### 0.7.1 — Application Layer Boundary Definition ✔ Completed
Initial application layer boundary defined and bootstrap use-case migrated.

#### 0.7.2 — Authenticated Surface Use Cases Extraction ✔ Completed

#### 0.7.3 — Wallet Management Use Cases Consolidation ✔ Completed

#### 0.7.4 — Handler Simplification & Contract Preservation ✔ Completed

### 0.8 — Standardized Error Model ✔ Completed
Define a consistent error contract across all endpoints.

#### 0.8.1 — Error Contract Definition ✔ Completed
Introduced the standardized error envelope foundation and centralized HTTP error writing.

#### 0.8.2 — Error Type System Introduction ✔ Completed
Introduced centralized internal app error typing, shared auth error catalog normalization and reusable error factories without changing the 0.8.1 transport envelope.

#### 0.8.3 — Auth Surface Error Standardization ✔ Completed
Auth handlers now consume centralized app-error factories and catalog resolution directly, while preserving the 0.8.1 envelope and avoiding cyclic import reintroduction.

#### 0.8.4 — Error Mapping Hardening & Contract Tests ✔ Completed
Added contract-focused tests for centralized app-error mapping and HTTP envelope serialization so representative `code/status/category` behavior and transport fallback paths are now frozen by automated coverage.

### 0.9 — API Versioning Strategy ⏳ In Progress
Formalize canonical API versioning, preserve legacy route compatibility and freeze the public contract-evolution rules needed before authorization and later domain growth.

#### 0.9.1 — Versioning Policy Definition ✔ Completed
Path-based versioning is now the canonical strategy, `v1` semantics are explicitly frozen and the legacy/non-canonical compatibility rule is documented.

#### 0.9.2 — Router Versioning Foundation ✔ Completed
Materialized the canonical `/api/v1/...` route surface in the real router while preserving the current legacy endpoints and reusing the same handler/application behavior.

#### 0.9.3 — Authenticated Surface Version Freezing ⬜ Pending
Bind the current auth/bootstrap/profile/settings/session/wallet surfaces to the canonical `v1` contract without changing their business behavior.

#### 0.9.4 — Version-aware Contract Testing ⬜ Pending
Extend contract-level tests so legacy and canonical routes are protected from divergence.

#### 0.9.5 — Documentation Consolidation ⬜ Pending
Propagate the versioning model through the trunk documentation set and leave explicit frontend/backend alignment notes for the Stage 0 completion path.

### 0.10 — Authorization Layer ⬜ Pending
Add roles, permissions, and enforcement mechanisms.

### 0.11 — Domain Module Pattern ⬜ Pending
Standardize module structure for future domain expansion.

### 0.12 — Read/Write Model Separation ⬜ Pending
Formalize CQRS-lite approach across the system.

### 0.13 — Observability & Diagnostics Foundation ⬜ Pending
Introduce structured logging, tracing, and metrics.

### 0.14 — Contract Hardening & Freeze ⬜ Pending
Stabilize and freeze API contracts before feature expansion.

---

## 👤 Stage 1 — Identity & Wallets

- User model
- Authentication
- Web3 login
- Wallet linking

---

## 🔗 Stage 2 — Blockchain Integration

- RPC integration
- Balance reads
- Allowances
- Indexing

---

## 🧠 Stage 3 — DEX Contracts

- AMM design
- Factory
- Pools
- Router
- Testing

---

## ⚙️ Stage 4 — DEX Backend

- Quotes
- Routing
- Swap intents
- Liquidity APIs

---

## 🌐 Stage 5 — APIs

- REST
- WebSocket
- Frontend contract

---

## 🔄 Stage 6 — Hybrid Expansion

- Internal ledger
- P2P design

---

## 🔐 Stage 7 — Security & Ops

- Hardening
- Observability
- Compliance hooks

---

## 🧪 Stage 8 — Testing

- Integration tests
- Internal testing
- Release candidate