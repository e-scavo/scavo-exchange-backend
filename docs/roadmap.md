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

### 0.7 — Application Layer Foundation ⬜ Pending
Introduce an application layer to separate HTTP from business logic.

### 0.8 — Standardized Error Model ⬜ Pending
Define a consistent error contract across all endpoints.

### 0.9 — API Versioning Strategy ⬜ Pending
Introduce versioning to support backward-compatible evolution.

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