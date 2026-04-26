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

### 0.9 — API Versioning Strategy ✔ Completed
Formalize canonical API versioning, preserve legacy route compatibility and freeze the public contract-evolution rules needed before authorization and later domain growth.

#### 0.9.1 — Versioning Policy Definition ✔ Completed
Path-based versioning is now the canonical strategy, `v1` semantics are explicitly frozen and the legacy/non-canonical compatibility rule is documented.

#### 0.9.2 — Router Versioning Foundation ✔ Completed
Materialized the canonical `/api/v1/...` route surface in the real router while preserving the current legacy endpoints and reusing the same handler/application behavior.

#### 0.9.3 — Authenticated Surface Version Freezing ✔ Completed
Bound the current auth/bootstrap/profile/settings/session/wallet surfaces explicitly to the canonical `v1` contract semantics while preserving the same business behavior across legacy and canonical entry paths.

#### 0.9.4 — Version-aware Contract Testing ✔ Completed
Extended representative contract-level tests so legacy and canonical routes are protected from silent divergence across key authenticated and auth-adjacent scenarios.

#### 0.9.5 — Documentation Consolidation ✔ Completed
Consolidated the trunk documentation set so roadmap, status, architecture, testing and handoff now describe the same completed Phase 0.9 state and the same frontend/backend alignment rule.

### 0.10 — Authorization Layer ✔ Completed
Introduce a structured authorization layer on top of the authenticated Stage 0 foundation without breaking the stabilized transport and application contracts.

#### 0.10.1 — Authorization Model Definition ✔ Completed
Introduced the foundational authorization primitives under `internal/core/authorization`, including roles, permissions, static role-to-permission mapping and the authorization-subject model that later subphases will propagate and evaluate.

#### 0.10.2 — Authorization Context & Middleware ✔ Completed
Attach authorization subject data to the authenticated request lifecycle and extend middleware/context propagation without changing endpoint contracts yet.

#### 0.10.3 — Policy Evaluation Layer ✔ Completed
Introduce a centralized policy-evaluation boundary so handlers can ask authorization questions without embedding permission logic ad hoc.

#### 0.10.4 — Endpoint-Level Enforcement ✔ Completed
Apply progressive authorization checks to selected authenticated endpoints through centralized policy middleware while preserving backward compatibility and the current Stage 0 contract discipline.

#### 0.10.5 — Documentation & Contract Consolidation ✔ Completed
Align the trunk documentation set to the delivered authorization behavior and close the phase with a coherent architectural and operational narrative.

### 0.11 — Domain Module Pattern ◑ In Progress
Standardize the internal structure of the current Stage 0 domain-facing modules so future growth builds on explicit HTTP / APP / DOMAIN boundaries, minimal internal contracts and reduced accidental coupling without changing the public API surface.

#### 0.11.1 — Domain Module Pattern Definition ✔ Completed
Define the standard internal module shape, layer responsibilities, dependency direction and naming conventions that the currently relevant modules will follow.

#### 0.11.2 — User Module Refactor ✔ Completed
Apply the pattern to `internal/modules/user` so user-facing transport handlers stop carrying orchestration responsibilities and the module becomes the first concrete reference implementation of the pattern.

This subphase is already delivered in repository state: the user module now exposes explicit `app`, `domain` and `repository` boundaries while preserving backward-compatible package access for existing consumers.

#### 0.11.3 — UserSettings Module Refactor ✔ Completed
Align `internal/modules/usersettings` to the same pattern while preserving its own configuration-oriented semantics and keeping it distinct from the user-entity boundary.

This subphase is now delivered in repository state: the user settings module exposes explicit `app`, `domain` and `repository` boundaries while preserving backward-compatible package access for existing consumers and keeping settings-specific validation/orchestration in the application layer.

#### 0.11.4 — Auth Module Alignment ✔ Completed
Aligned `internal/modules/auth` to the pattern conservatively while preserving the stabilized challenge, verify, wallet-management, session, bootstrap and authorization-adjacent behavior.

This subphase is delivered in repository state: `auth/app` owns application/runtime orchestration, `auth/domain` owns canonical wallet contracts, root `auth` remains a compatibility surface, wallet HTTP management is narrowed to transport/delegation, and repository migration is explicitly deferred behind a transitional façade.

#### 0.11.5 — Cross-Module Contract Consolidation ◑ In Progress
Consolidate the explicit internal contracts between `auth`, `user` and `usersettings` so coordination remains possible without reintroducing ownership confusion or concrete cross-module coupling.

Planned internal steps:

- 0.11.5.0 — Subphase Definition & Documentation Alignment
- 0.11.5.1 — Dependency Mapping
- 0.11.5.2 — Contract Extraction
- 0.11.5.3 — Interface Alignment
- 0.11.5.4 — Runtime Compatibility Validation

#### 0.11.6 — Documentation & Phase Closure ◑ Defined
Update the trunk documentation set so the adopted module pattern, its scope and its completion criteria are recorded consistently across roadmap, status, architecture, handoff and the dedicated phase document.

#### Expected outcome
By the end of Phase 0.11, the backend should keep the same external HTTP behavior while exposing a more uniform internal module structure, clearer ownership boundaries and safer foundations for later Stage 0 and post-Stage 0 work.

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