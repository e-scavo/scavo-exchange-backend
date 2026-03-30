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

With 0.4.15:

- wallet authentication is stable
- identity is durable
- ownership is persisted
- authenticated wallet linking is available
- wallet-owned account merge execution is available
- explicit primary-wallet switching is available
- detach eligibility can be evaluated safely before unlink execution
- detach execution is implemented for already eligible owned wallets
- detached wallets are explicitly treated as reusable known identities
- detached wallets can be reattached either through protected linking or through wallet-login bootstrap rebound
- detached identities now preserve minimal audit-ready lifecycle metadata through `detached_at`

Next expected evolution:

➡️ **0.4.16 — Detached Identity Extended Audit Semantics (only if later needed)**


---

## 0.4.14 — Detached Wallet Reattachment Semantics and Lifecycle Clarification

### Objective
Clarify the post-detach lifecycle of wallet identities without introducing premature schema expansion or audit complexity.

### Scope
- confirm that detached wallet identities remain reusable known identities
- validate that authenticated wallet-linking can reattach a previously detached wallet
- validate that wallet-login bootstrap can rebound a detached wallet into a wallet-owned user identity
- align documentation with the actual current lifecycle behavior

### Delivered
- service-level reattachment test coverage after detach
- service-level wallet-login rebound coverage after detach
- handler-level coverage for detached-wallet reattachment under the authenticated link flow
- explicit documentation that detached wallets are not deleted or archived in the current phase

### Lifecycle Semantics Clarified
After detach, the wallet identity:
- keeps its durable wallet identity record
- keeps its normalized address
- clears `user_id`
- clears `linked_at`
- clears `is_primary`
- remains reusable for future attachment

This means the current system treats detached wallets as **known reusable identities**, not as terminal, archived, or deleted entities.

### Not Introduced Here
- `detached_at`
- `detached_by_user_id`
- lifecycle audit tables
- event sourcing
- archival or soft-delete semantics

### Validation
Validated through:
- detached-wallet reattachment tests under `WalletLinkingService`
- detached-wallet wallet-login rebound tests under `WalletVerificationService`
- HTTP handler coverage for reattaching a detached wallet through the authenticated link contract


---

## 0.4.15 — Detached Identity Audit Readiness

### Objective
Add minimal persisted lifecycle metadata for detached wallet identities without introducing heavy audit tables, event sourcing, or archival semantics.

### Scope
- persist minimal detached-wallet lifecycle metadata
- ensure detach execution stamps that metadata
- ensure reattachment and wallet-login rebound preserve that metadata
- align documentation with the new audit-ready lifecycle contract

### Delivered
- `detached_at` added to `WalletIdentity`
- PostgreSQL migration adding `detached_at` to `auth_wallet_identities`
- in-memory and PostgreSQL store support for reading and writing `detached_at`
- detach execution updated to stamp `detached_at`
- test coverage proving detached metadata survives later reattachment and wallet-login rebound
- documentation updates aligning detached-wallet reuse with minimal audit readiness

### Lifecycle Semantics Clarified
A detached wallet identity now preserves minimal lifecycle evidence:
- `user_id` is cleared on detach
- `linked_at` is cleared on detach
- `is_primary` is cleared on detach
- `detached_at` is stamped on detach
- `detached_at` survives later reattachment or wallet-login rebound

This means the system can now distinguish:
- a wallet identity that has never been detached
- a wallet identity that was previously detached and later reused

### Not Introduced Here
- `detached_by_user_id`
- queryable detached-history tables
- event sourcing
- archival or soft-delete semantics
- multi-event lifecycle reporting

### Validation
Validated through:
- detach execution coverage returning detached wallet metadata with `detached_at`
- service-level reattachment tests proving `detached_at` survives authenticated relinking
- service-level wallet-login rebound tests proving `detached_at` survives bootstrap reuse
- documentation alignment across README and docs/*


---

## 0.4.16 — Wallet Identity Read Model Enrichment

### Objective
Expose the real wallet identity lifecycle state more explicitly through the authenticated wallet inventory endpoint.

### Scope
- introduce an explicit wallet read model for `GET /auth/wallets`
- expose lifecycle-aware fields already present in the backend model
- add a conservative derived `status` field
- validate active-wallet serialization and detached-then-reattached visibility
- align documentation with the enriched public contract

### Delivered
- explicit `WalletReadModel` for authenticated wallet listing
- public response fields:
  - `id`
  - `address`
  - `user_id`
  - `linked_at`
  - `detached_at`
  - `is_primary`
  - `status`
- handler-level coverage for active wallet inventory
- handler-level coverage proving `detached_at` remains visible after detach + reattach
- documentation updates aligning the API contract with the lifecycle-aware wallet identity model

### Lifecycle Semantics Clarified
The authenticated wallet inventory now exposes both:
- current ownership state
- minimal historical detach evidence

The derived `status` field is conservative:
- `active` when the wallet is currently linked to a user
- `detached` when no owner exists and `detached_at` is present
- `unlinked` when neither ownership nor detach evidence exists

For `GET /auth/wallets`, the practical operational case remains `active`, because the endpoint lists wallets currently owned by the authenticated user.

### Not Introduced Here
- filtering or query parameters for wallet inventory
- pagination
- search
- admin inventory views
- richer detached-history endpoints
- ownership-rule changes

### Validation
Validated through:
- full `go test ./...`
- handler-level wallet inventory serialization checks
- explicit validation that detached-then-reattached wallets still expose `detached_at`
- documentation alignment across README and docs/*



## 0.4.17 — Wallet Inventory Query Filtering and Sorting

### Objective
Add small, explicit query semantics to the authenticated wallet inventory endpoint without changing domain or persistence behavior.

### Scope
- add optional `status` filtering for `GET /auth/wallets`
- add optional `primary` filtering for `GET /auth/wallets`
- add optional `linked_at` sorting with `asc|desc` ordering
- keep backward compatibility when query params are omitted
- validate invalid query values explicitly at handler level
- extend handler-level coverage for the new query contract

### Delivered
- optional query params:
  - `status=active|detached`
  - `primary=true|false`
  - `sort=linked_at`
  - `order=asc|desc`
- explicit `400` errors for invalid query combinations or unsupported values
- filtering and sorting applied only on the lifecycle-aware read model already exposed by the handler
- documentation updates aligning the wallet inventory endpoint with the new query semantics

### Lifecycle / Query Semantics Clarified
This subphase does not widen the ownership scope of the endpoint.

`GET /auth/wallets` still lists only wallets currently owned by the authenticated user. Because of that:
- `status=active` is the normal operational case
- `status=detached` is a valid filter but is expected to produce an empty result under the current route contract
- `primary=true|false` operates only within the authenticated owned-wallet inventory
- explicit `sort=linked_at` overrides only the response ordering, not ownership semantics

### Not Introduced Here
- pagination
- text search
- admin inventory views
- store-level query APIs
- SQL filtering or ordering changes
- detached-wallet history reporting
- ownership-rule changes

### Validation
Validation path for this subphase:
- `go test ./...`
- handler-level validation for filter/sort compatibility
- handler-level validation for invalid query parameters returning `400`

