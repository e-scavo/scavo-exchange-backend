# Phase 0.4 — Auth and User Stabilization

## 🧠 Objective

Stabilize authentication, user identity, and wallet-based login flows, transitioning from ephemeral wallet identity toward a durable, unified identity model suitable for future exchange-grade features.

---

## 📌 Initial Context

At the beginning of Phase 0.4:

- authentication was partially implemented
- wallet login existed but lacked persistence
- no durable relationship existed between wallets and platform users
- identity was fragmented across sessions

This phase progressively transforms the system into a consistent identity layer.

---

## 🚧 Problem Statement

The system required:

- deterministic authentication flows
- persistent identity representation
- wallet-based login suitable for production evolution
- a unified identity model compatible with multiple auth methods
- explicit wallet ownership semantics
- a safe bridge toward user-managed account expansion

---

## 🔍 Scope

Phase 0.4 focuses on:

- authentication stabilization
- wallet login correctness
- identity persistence
- user abstraction
- session unification
- ownership foundations
- authenticated wallet-management primitives

---

## 🧩 Subphases Breakdown

### 0.4.1 — Auth Baseline Stabilization

#### Implemented
- initial auth service structure
- token generation baseline
- basic login handling

#### Result
- system capable of issuing JWT tokens
- identity consistency still incomplete

---

### 0.4.2 — Token Service Stabilization

#### Implemented
- token service refactor
- claim normalization
- expiration handling

#### Result
- reliable JWT issuance
- improved token parsing consistency

---

### 0.4.3 — Session Model Stabilization

#### Implemented
- session abstraction
- `/auth/me` and `/auth/session` endpoints
- claims hydration

#### Result
- session identity accessible across requests
- still not durably tied to persistent entities

---

### 0.4.4 — Wallet Challenge Flow

#### Implemented
- wallet challenge creation
- message signing model
- expiration control

#### Result
- secure wallet-authentication entry point

---

### 0.4.5 — Wallet Verification Baseline

#### Implemented
- signature verification
- address recovery
- challenge validation

#### Result
- functional wallet login
- still stateless from the identity-model perspective

---

### 0.4.6 — Wallet Identity Persistence

#### Implemented
- `auth_wallet_identities` table
- wallet identity storage
- durable challenge store

#### Result
- wallet identity persisted
- no durable user linkage yet

---

### 0.4.7 — Wallet ↔ User Linking and Unified Identity Model

#### Implemented
- durable user creation for wallet login
- `auth_wallet_identities.user_id`
- wallet identity linked to platform user
- unified JWT identity model

#### Result
- wallet login produces a durable user
- `/auth/me` resolves unified identity
- system transitions from wallet identity → user identity

---

### 0.4.8 — Account Consolidation and Multi-Wallet Ownership Foundations

#### Implemented
- removal of the 1:1 wallet-user restriction
- ownership metadata introduced:
  - `linked_at`
  - `is_primary`
- safe attachment semantics preventing reassignment to another user
- authenticated read-only wallet listing via `GET /auth/wallets`

#### Result
- one durable platform user can own multiple wallets
- primary wallet concept is established
- ownership becomes a first-class persisted concern

---

### 0.4.9 — User-Driven Wallet Linking Contract and Protected Account Merge Preparation

#### Implemented
- authenticated linking challenge flow:
  - `POST /auth/wallets/link/challenge`
- authenticated linking verification flow:
  - `POST /auth/wallets/link/verify`
- challenge metadata extensions:
  - `purpose`
  - `requested_by_user_id`
- challenge-purpose separation between:
  - login bootstrap
  - wallet linking
- user-bound challenge validation for linking flows
- protected rejection of linking wallets already owned by another user
- protected rejection of relinking a wallet already owned by the current user
- secondary-wallet attach behavior with `is_primary = false`
- updated wallet inventory response after successful linking

#### Result
- the backend now supports the first controlled wallet-management operation under an authenticated user session
- the system advances from ownership persistence toward account-level wallet control without introducing risky merge automation

---

### 0.4.10 — User-Driven Wallet-Owned Account Merge Execution

#### Implemented
- authenticated account-merge challenge flow:
  - `POST /auth/account/merge/wallet/challenge`
- authenticated account-merge verification flow:
  - `POST /auth/account/merge/wallet/verify`
- challenge-purpose expansion with `account_merge`
- source-wallet signature requirement before merge execution
- protected rejection of merge attempts against unlinked wallets
- protected rejection when the wallet already belongs to the current user
- store-level atomic wallet-ownership consolidation from source user to target user
- deterministic preservation of the target account primary wallet when one already exists
- merged wallet inventory response after successful consolidation

#### Result
- the backend now supports explicit execution of a wallet-owned account merge under authenticated user control
- the previous 0.4.9 preparation step is converted into a real, constrained merge operation without weakening ownership rules

---

### 0.4.11 — Primary Wallet Management and Ownership Safety Hardening

#### Implemented
- authenticated primary-wallet switch flow:
  - `POST /auth/wallets/primary`
- store-level `SetPrimary(...)` contract
- protected rejection when the wallet is missing
- protected rejection when the wallet does not belong to the current authenticated user
- deterministic single-primary reassignment within the owned-wallet set
- refreshed wallet inventory response after successful primary switching

#### Result
- the backend now supports the first explicit post-merge wallet-ownership management action
- ownership can be reorganized safely without changing wallet attachment or merge history

---

### 0.4.12 — Wallet Detach Contract Preparation and Ownership Guardrails

#### Implemented
- authenticated detach-eligibility evaluation flow:
  - `POST /auth/wallets/detach/check`
- detach-check response contract with:
  - `eligible`
  - `is_primary`
  - `owned_wallet_count`
  - `reasons`
- protected rejection when the wallet is missing
- protected rejection when the wallet does not belong to the current authenticated user
- conservative non-eligibility when the wallet is the current primary wallet
- conservative non-eligibility when detach would leave the user without any wallets
- explicit reasoning contract for future unlink-safe product work

#### Result
- the backend now supports detach-preparation under authenticated control without introducing destructive ownership changes
- future wallet detach execution can be designed against an already enforced eligibility contract instead of relying on implicit assumptions

---

## 🧱 Root Cause Analysis

The initial architecture lacked:

- persistent identity boundaries
- clear ownership semantics
- separation between wallet identity and user identity
- any authenticated contract for user-managed wallet expansion

Each subphase incrementally addressed one structural gap while preserving backward compatibility.

---

## 📂 Files Affected

### Core modules
- `internal/modules/auth/*`
- `internal/modules/user/*`
- `internal/core/auth/*`

### Persistence
- `auth_wallet_challenges`
- `auth_wallet_identities`
- `users`

### HTTP layer
- wallet challenge handlers
- wallet verify handlers
- wallet inventory endpoint
- authenticated wallet-link handlers
- authenticated wallet-account-merge handlers

---

## ⚙️ Implementation Characteristics

- backward-compatible with previous wallet login flow
- incremental persistence evolution
- stateless sessions with durable backing state
- in-memory fallback preserved
- challenge-purpose separation introduced without forking the entire challenge subsystem
- merge execution remains explicit and wallet-signed
- ownership rules remain enforced at the store layer
- link contract remains explicitly authenticated

---

## 🧪 Validation

### Code-level

```bash
go test ./...
```

### Behavioral
- wallet login creates or resolves durable user
- wallet identity is persisted
- ownership is persisted
- `/auth/me` resolves unified identity
- `/auth/wallets` returns owned wallets
- `/auth/wallets/link/challenge` creates user-bound link challenge
- `/auth/wallets/link/verify` attaches a new secondary wallet
- `/auth/account/merge/wallet/challenge` creates a user-bound merge challenge
- `/auth/account/merge/wallet/verify` consolidates wallet ownership from the source wallet-owned account

---

## 📈 Release Impact

- enables authenticated wallet-owned account merge execution without destabilizing login
- keeps ownership model strict while expanding functionality
- converts merge preparation into a real explicit flow
- establishes safer preconditions for later unlink and primary-switch work

---

## ⚠️ Risks

- challenge-purpose validation must remain strict
- user-bound link challenge checks must not be bypassed
- future unlink / transfer flows must preserve current ownership invariants
- later merge flows must not weaken the explicitness introduced here
- explicit primary-wallet reassignment must preserve single-primary invariants

---

## ❌ What This Phase Does NOT Solve

- wallet unlink
- cross-user wallet transfer
- arbitrary cross-user transfer outside wallet-signed merge execution
- merged-source user archival or aliasing
- token revocation
- refresh tokens
- persistent sessions

---

## 🧭 Conclusion

Phase 0.4 now establishes a strong identity and wallet-ownership foundation.

With 0.4.12:

- wallet authentication is stable
- identity is durable
- ownership is persisted
- authenticated wallet linking is available
- wallet-owned account merge execution is available
- explicit primary-wallet switching is available
- detach eligibility can be evaluated safely before unlink execution
- the backend is prepared for controlled wallet-detach design

Next expected evolution:

➡️ **0.4.13 — Wallet Detach Execution Design**
