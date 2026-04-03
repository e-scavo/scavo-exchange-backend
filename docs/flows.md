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



## 0.4.18 — Wallet Inventory Pagination and Windowed Response

### Objective
Add small, explicit pagination semantics to the authenticated wallet inventory endpoint without changing domain or persistence behavior.

### Scope
- add optional `limit` pagination for `GET /auth/wallets`
- add optional `offset` pagination for `GET /auth/wallets`
- expose additive response metadata: `total`, `limit`, `offset`
- keep backward compatibility with the existing filter/sort contract
- validate invalid pagination values explicitly at handler level
- extend handler-level coverage for valid and invalid pagination windows

### Delivered
- optional query params:
  - `limit=<positive integer>`
  - `offset=<non-negative integer>`
- additive wallet inventory response metadata:
  - `total`
  - `limit`
  - `offset`
- strict `400` errors for malformed pagination values
- pagination applied only after the existing wallet inventory filter/sort pipeline
- documentation updates aligning the wallet inventory endpoint with its new windowed response contract

### Lifecycle / Query Semantics Clarified
This subphase still does not widen the ownership scope of the endpoint.

`GET /auth/wallets` still lists only wallets currently owned by the authenticated user. Because of that:
- pagination operates only inside the authenticated owned-wallet inventory
- `total` reflects the filtered inventory size before the requested window is applied
- `limit=0` means no explicit page-size cap was requested
- `offset=0` remains the default starting position

### Not Introduced Here
- cursor pagination
- next-page tokens
- text search
- admin inventory views
- store-level pagination APIs
- SQL pagination
- detached-wallet history reporting
- ownership-rule changes

### Validation
Validation path for this subphase:
- `go test ./...`
- handler-level validation for `limit` and `offset` compatibility
- handler-level validation for invalid pagination parameters returning `400`
- handler-level validation for empty but valid inventory windows


## 0.4.19 — Wallet Inventory Navigation Metadata

### Updated Read Flow
Authenticated wallet inventory flow now behaves as follows:

1. authenticated user calls `GET /auth/wallets`
2. backend loads wallets currently owned by that durable user
3. backend maps identities into the lifecycle-aware wallet read model
4. optional filters are applied
5. optional ordering is applied
6. optional pagination window is applied
7. backend returns:
   - `wallets`
   - `total`
   - `limit`
   - `offset`
   - `returned`
   - `has_more`

### Response Semantics
- `returned` is the size of the current `wallets` array
- `has_more` indicates whether another wallet exists after the current window
- when `limit=0`, the request is treated as unbounded and `has_more=false`


## 0.4.21 — Wallet Inventory Query Parameter Contract Hardening

### Flow impact
`GET /auth/wallets` keeps the same authenticated ownership-scoped flow, but now makes the query contract more explicit:

- `order` is rejected unless `sort` is present
- `sort=linked_at` defaults to ascending order when `order` is omitted
- `offset` without `limit` remains valid and produces an unbounded window with no navigation hints

### Request semantics
Evaluation order remains:
1. authenticate request
2. list wallet identities for the authenticated user
3. map to wallet read model
4. apply filters
5. apply sort contract
6. apply pagination contract
7. compute navigation metadata
8. return response

### Outcome
The flow remains read-only and ownership-scoped while becoming more predictable for clients consuming wallet inventory queries.


## 0.4.22 — Wallet Inventory Response Contract Clarification

### Flow impact
The authenticated wallet inventory flow does not change behavior in 0.4.22. This subphase only makes the visible response contract explicit for operators and client implementers.

### Response semantics clarified
For `GET /auth/wallets`, the response contract is now documented as follows:

- `wallets`: current filtered, sorted, and windowed rows
- `total`: total filtered row count before windowing
- `limit`: requested page size, or `0` for unbounded requests
- `offset`: requested starting offset
- `returned`: number of rows actually returned in the current response
- `has_more`: whether another row exists after the current window
- `next_offset`: next bounded offset hint when another page exists
- `previous_offset`: previous bounded offset hint when a previous page exists

### Outcome
The flow remains read-only and ownership-scoped, while the response contract becomes fully explicit and aligned with the implemented endpoint behavior.


## 0.4.23 — Wallet Inventory Query Examples Closure

### Flow impact
The authenticated wallet inventory flow does not change behavior in 0.4.23. This subphase only adds concrete examples showing how the already-supported query contract is consumed.

### Example request patterns
- `GET /auth/wallets`
- `GET /auth/wallets?primary=true`
- `GET /auth/wallets?sort=linked_at&order=desc`
- `GET /auth/wallets?limit=2&offset=0`
- `GET /auth/wallets?order=desc` → contractual `400` because `sort` is required

### Outcome
The flow remains read-only and ownership-scoped, while the documentation now shows concrete query usage and error examples aligned with the real endpoint contract.


## Phase 0.4.24 — Wallet Inventory Manual Validation Closure

### Objective
Close the manual validation layer for `GET /auth/wallets` without changing code, stores, or persistence.

### Delivered
- consolidated manual validation coverage for base, filtered, sorted, paginated, and unbounded inventory requests
- explicit operator checks for navigation metadata (`returned`, `has_more`, `next_offset`, `previous_offset`)
- explicit invalid-query checks for the hardened query-parameter contract
- documentation-only closure of the wallet inventory manual validation layer

### Conclusion
Phase 0.4.24 leaves the wallet inventory flow behavior unchanged while making its operator-facing manual verification path explicit and complete.


## Phase 0.4.25 — Wallet Actionability Read Model Preparation

### Flow impact
The authenticated wallet inventory flow now exposes advisory actionability hints for each listed wallet while keeping all real execution authority in the existing primary-switch and detach endpoints.

### Added read-model hints
- `can_set_primary`: `true` when the wallet is active and not already primary
- `can_detach`: `true` when the wallet is active, not primary, and the user would still own another active wallet
- `detach_block_reasons`: advisory block reasons derived from the detach-domain constants

### Outcome
The inventory flow remains read-only and ownership-scoped, but it now gives client code enough per-wallet actionability context to prepare wallet-management UI affordances without inferring detach or primary-switch rules on its own.


## Phase 0.4.26 — Wallet Detach Check Read Consistency

### Flow impact
The inventory flow and the detach-check flow keep their existing behavior in 0.4.26, but their relationship is now explicitly covered and documented.

### Consistency semantics
For the same authenticated user and wallet inventory:

- `can_detach=true` in `GET /auth/wallets` should remain compatible with `eligible=true` in `POST /auth/wallets/detach/check`
- `can_detach=false` should remain compatible with `eligible=false` under the same detach-domain reasons
- inventory-side `detach_block_reasons` remain advisory, while `detach/check` stays authoritative before execution

### Outcome
The backend now protects the semantic alignment between wallet inventory actionability hints and detach-check eligibility without changing detach-domain rules or persistence behavior.


## Phase 0.4.27 — Wallet Primary Switch Read Consistency

### Flow impact
The inventory flow and the primary-switch flow keep their existing behavior in 0.4.27, but their relationship is now explicitly covered and documented.

### Consistency semantics
For the same authenticated user and wallet inventory:

- `can_set_primary=false` should remain compatible with the currently primary wallet in `GET /auth/wallets`
- `can_set_primary=true` should remain compatible with a valid candidate for `POST /auth/wallets/primary`
- after a successful primary switch, the promoted wallet should become `is_primary=true` and `can_set_primary=false` in the refreshed inventory
- the former primary should become `is_primary=false` and `can_set_primary=true` in the refreshed inventory

### Outcome
The backend now protects the semantic alignment between wallet inventory primary-actionability hints and primary-switch execution without changing primary-domain rules or persistence behavior.

## Phase 0.4.28 — Wallet Management Read Flow Closure

### Flow impact
The authenticated wallet-management surface is now documented as one continuous read flow instead of separate isolated capabilities.

### Read-flow semantics
For the same authenticated user:

- `GET /auth/wallets` is the inventory entry point and the only place where advisory actionability hints should be read
- `can_set_primary=true` should be interpreted as a valid invitation to call `POST /auth/wallets/primary` for that wallet
- `can_detach=true` should be interpreted as a valid invitation to call `POST /auth/wallets/detach/check` and, if still eligible, `POST /auth/wallets/detach`
- after either action path succeeds, the inventory should be refreshed and treated as the new source of truth for `is_primary`, `status`, `detached_at`, and updated actionability hints
- inventory hints remain advisory; primary and detach endpoints remain authoritative for execution

### Outcome
The backend keeps the same runtime behavior as 0.4.27, but the wallet-management path is now documented as a complete inventory-driven read flow from initial listing to refreshed post-action inventory.



## Phase 0.4.29 — Wallet Detach Execute Read Consistency

### Flow impact
The authenticated wallet inventory and detach-execution flow keep their existing runtime behavior in 0.4.29, but their relationship is now explicitly covered and documented.

### Consistency semantics
For the same authenticated user and wallet inventory:

- `can_detach=true` should remain compatible with a wallet that can be detached successfully through `POST /auth/wallets/detach`
- after a successful detach, the refreshed inventory should no longer expose that wallet as attached to the authenticated user
- the remaining attached wallet set should recalculate `can_detach` and `detach_block_reasons` coherently for the new post-detach state
- inventory hints remain advisory; detach execution remains authoritative for the actual detach outcome

### Outcome
The backend now protects the semantic alignment between wallet inventory detach-actionability hints, detach execution, and refreshed inventory without changing detach-domain rules or persistence behavior.


## Phase 0.4.30 — Wallet Management Contract Consolidation

The authenticated wallet inventory and wallet action endpoints keep their existing runtime behavior in 0.4.30, but their relationship is now explicitly consolidated as one contract:

- `GET /auth/wallets` is the advisory read surface
- `POST /auth/wallets/primary` is the primary-switch execution surface
- `POST /auth/wallets/detach/check` is the explicit detach-eligibility surface
- `POST /auth/wallets/detach` is the detach execution surface
- refreshed `GET /auth/wallets` is the post-action observable state

This phase does not change handlers, stores, persistence, or domain rules. It closes the wallet-management flow as one consolidated inventory-driven contract.


## Phase 0.4.31 — Wallet Auth Bootstrap Purpose Enforcement

### Implemented
- wallet verify/login now enforces `purpose = auth_bootstrap`
- `wallet_link` and `account_merge` challenges are explicitly rejected by wallet-auth bootstrap
- HTTP wallet verify now surfaces `wallet_challenge_purpose_mismatch` for wrong-purpose challenges
- tests now cover service-level and handler-level purpose mismatch cases

### Result
- wallet login, detach rebound, authenticated linking, and wallet-owned account merge keep their existing runtime semantics
- the remaining contract gap is closed: challenge purposes are no longer reusable across different wallet lifecycle entry points
- Phase 0.4 closes with stricter purpose isolation and no ownership or persistence redesign
