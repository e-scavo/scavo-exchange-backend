# SCAVO Exchange — Backend

## 🧠 Overview

SCAVO Exchange Backend is a Go-based service that provides authentication, user management, and wallet-based identity for the SCAVO ecosystem.

The backend follows a **wallet-first identity model** that progressively evolves into a **durable account architecture** suitable for exchange-grade ownership, linking, and future multi-auth identity expansion.

---

## 🏗️ Architecture Principles

- **Wallet-first authentication**
- **Durable user abstraction**
- **Stateless JWT sessions**
- **Explicit ownership persistence**
- **Incremental account consolidation**
- **Database-backed persistence with in-memory fallback**

---

## 🚧 Current Stage

**Stage:** 0 — Foundation  
**Phase:** 0.4 — Auth and User Stabilization  
**Current Subphase:** **0.4.20 — Wallet Inventory Cursorless Navigation Hints**

---

## 🔐 Authentication Model

The backend currently supports two authentication methods:

### 1. Password-based authentication (dev only)
- intended only for internal development and testing
- not meant for production operation

### 2. Wallet-based authentication (EVM)

Base wallet login flow:

1. Client requests challenge  
   `POST /auth/wallet/challenge`

2. Backend creates challenge:
   - unique ID
   - wallet address binding
   - chain binding
   - expiration timestamp
   - challenge purpose metadata

3. Client signs the challenge message

4. Client verifies challenge  
   `POST /auth/wallet/verify`

5. Backend:
   - validates challenge state
   - verifies signature
   - consumes challenge
   - resolves wallet identity
   - resolves or creates durable user
   - enforces ownership invariants
   - issues JWT

---

## 🧩 Identity Model Evolution

### Pre 0.4.7
- wallet identity was not durably linked to a platform user
- session identity and persistent identity were not unified

### 0.4.7 — Wallet ↔ User Linking
- each wallet identity is linked to a durable user
- JWT identity becomes unified around `user_id`

### 0.4.8 — Multi-Wallet Ownership Foundations
wallet identities gained ownership metadata:

- `user_id`
- `linked_at`
- `detached_at`
- `is_primary`

This allowed:

- one user → multiple wallets
- explicit primary wallet designation
- ownership persistence independent from JWT sessions

### 0.4.9 — Authenticated Wallet Linking Contract
wallet management now supports an authenticated user-driven linking flow:

- `POST /auth/wallets/link/challenge`
- `POST /auth/wallets/link/verify`

This allows a signed secondary-wallet attachment flow without creating a new session or performing account merge heuristics.

### 0.4.10 — User-Driven Wallet-Owned Account Merge Execution
wallet management now also supports an authenticated merge flow for wallet-owned accounts:

- `POST /auth/account/merge/wallet/challenge`
- `POST /auth/account/merge/wallet/verify`

This allows the current authenticated user to absorb another wallet-owned account only after the source wallet explicitly signs a merge challenge.

### 0.4.11 — Explicit Primary-Wallet Switching
wallet management now also supports an authenticated primary-wallet switch flow:

- `POST /auth/wallets/primary`

This allows the current authenticated user to explicitly select which owned wallet is primary without changing ownership.

### 0.4.12 — Wallet Detach Eligibility Contract
wallet management now also supports an authenticated detach-eligibility evaluation flow:

- `POST /auth/wallets/detach/check`

This allows the current authenticated user to ask the backend whether one owned wallet is currently safe to detach, without changing ownership and without executing unlink behavior.

### 0.4.15 — Detached Identity Audit Readiness
wallet identities now preserve minimal detached-lifecycle audit metadata:

- `detached_at`

This allows the backend to distinguish a wallet that has never been detached from a wallet identity that was previously detached and later reused through linking or wallet-login rebound.


### 0.4.13 — Protected Wallet Detach Execution
wallet management now also supports an authenticated detach execution flow for already eligible owned wallets:

- `POST /auth/wallets/detach`

This allows the current authenticated user to detach one owned non-primary wallet only when the ownership guardrails introduced in 0.4.12 are satisfied.

### 0.4.14 — Detached Wallet Reattachment Semantics and Lifecycle Clarification
wallet lifecycle now explicitly clarifies what happens after detach, without introducing a new lifecycle table or schema state:

- detached wallet identities remain known wallet identities
- detached wallet identities retain their address and wallet identity record
- detached wallet identities clear `user_id`, `linked_at`, and `is_primary`
- detached wallet identities can be reattached through the authenticated linking flow
- detached wallet identities can also re-enter the wallet-login bootstrap flow and resolve back into a wallet-owned user

This phase formalizes that detached wallets are reusable known wallet identities rather than archived or terminal identities.

---

## 🗄️ Persistence Model

### Main tables involved

#### `auth_wallet_challenges`
stores challenge lifecycle and now also includes linking metadata:

- `purpose`
- `requested_by_user_id`

Used for:
- wallet auth bootstrap challenges
- authenticated wallet-link confirmation challenges
- authenticated wallet-owned account merge challenges

#### `auth_wallet_identities`
stores wallet registry and ownership metadata:

- `id`
- `address`
- `user_id`
- `linked_at`
- `detached_at`
- `is_primary`

#### `users`
stores durable platform users:

- wallet-backed users
- future multi-auth identities

---

## 🔌 API Endpoints

### Wallet Auth

#### `POST /auth/wallet/challenge`

Creates a login bootstrap challenge for wallet authentication.

Request:

```json
{
  "address": "0x...",
  "chain": "scavium"
}
```

---

#### `POST /auth/wallet/verify`

Verifies wallet signature and returns a JWT-backed session.

Request:

```json
{
  "challenge_id": "...",
  "address": "0x...",
  "signature": "0x..."
}
```

Response:

```json
{
  "access_token": "...",
  "token_type": "Bearer",
  "expires_in": 3600,
  "user_id": "...",
  "wallet_id": "...",
  "wallet_address": "0x...",
  "chain": "scavium",
  "auth_method": "wallet_evm"
}
```

---

### Wallet Ownership

#### `GET /auth/wallets`

Returns all wallet identities linked to the authenticated durable user.

Response:

```json
{
  "wallets": [
    {
      "id": "...",
      "address": "0x...",
      "user_id": "...",
      "linked_at": "...",
      "detached_at": null,
      "is_primary": true,
      "status": "active"
    }
  ],
  "total": 1,
  "limit": 0,
  "offset": 0
}
```

Supported optional query params:

- `status=active|detached`
- `primary=true|false`
- `sort=linked_at`
- `order=asc|desc`
- `limit=<positive integer>`
- `offset=<non-negative integer>`

---

### Authenticated Wallet Linking

#### `POST /auth/wallets/link/challenge`

Creates a wallet-linking challenge bound to the currently authenticated user.

Request:

```json
{
  "address": "0x...",
  "chain": "scavium"
}
```

Behavior:

- requires valid JWT
- challenge purpose becomes `wallet_link`
- challenge stores `requested_by_user_id`

---

#### `POST /auth/wallets/link/verify`

Verifies the linking signature and attaches the wallet to the current user as a **secondary wallet**.

Request:

```json
{
  "challenge_id": "...",
  "address": "0x...",
  "signature": "0x..."
}
```

Response:

```json
{
  "linked_wallet": {
    "id": "...",
    "address": "0x...",
    "user_id": "...",
    "linked_at": "...",
    "is_primary": false
  },
  "wallets": [
    {
      "id": "...",
      "address": "0x...",
      "user_id": "...",
      "linked_at": "...",
      "is_primary": true
    },
    {
      "id": "...",
      "address": "0x...",
      "user_id": "...",
      "linked_at": "...",
      "is_primary": false
    }
  ]
}
```

---

### Authenticated Wallet-Owned Account Merge

#### `POST /auth/account/merge/wallet/challenge`

Creates an account-merge challenge bound to the currently authenticated user.

Request:

```json
{
  "address": "0x...",
  "chain": "scavium"
}
```

Behavior:

- requires valid JWT
- challenge purpose becomes `account_merge`
- challenge stores `requested_by_user_id`

---

#### `POST /auth/account/merge/wallet/verify`

Verifies the merge signature and consolidates all wallets from the source wallet-owned account into the current user.

Request:

```json
{
  "challenge_id": "...",
  "address": "0x...",
  "signature": "0x..."
}
```

Response:

```json
{
  "merged_wallet": {
    "id": "...",
    "address": "0x...",
    "user_id": "...",
    "linked_at": "...",
    "is_primary": false
  },
  "source_user_id": "u_wallet_...",
  "target_user_id": "u_current_user",
  "wallets": [
    {
      "id": "...",
      "address": "0x...",
      "user_id": "...",
      "linked_at": "...",
      "is_primary": true
    }
  ]
}
```

---

## 🧾 JWT Claims

JWT tokens include:

- `user_id`
- `wallet_id`
- `wallet_address`
- `auth_method`
- `exp`
- `iat`
- `nbf`

Wallet linking and wallet-owned account merge do **not** mint a new token. Both operate under the existing authenticated session.

---

## 🧪 Testing

Run:

```bash
go test ./...
```

Focus areas added in 0.4.16:

- wallet inventory read-model enrichment
- explicit lifecycle-aware wallet inventory serialization
- `status` derivation for owned-wallet inventory responses
- visibility of `detached_at` after detach + reattach
- existing link, merge, primary-switch, and detach coverage preserved

---

## 🚧 What 0.4.16 Solves

- authenticated user-driven wallet linking
- authenticated wallet-owned account merge execution
- protected primary-wallet switching under an authenticated user session
- authenticated wallet detach-eligibility evaluation under an authenticated user session
- challenge purpose separation between login, linking, and merge
- challenge-to-user binding through `requested_by_user_id`
- protected secondary-wallet attachment
- protected wallet-signed ownership consolidation
- deterministic single-primary wallet reassignment
- explicit detach rejection reasons for ownership-unsafe states

---

## ❌ What 0.4.16 Does Not Solve Yet

- wallet unlink API
- arbitrary cross-user ownership transfer outside wallet-signed merge
- merge between wallet-backed and other auth methods
- refresh tokens
- token revocation
- persistent authenticated sessions
- archival or alias records for merged source users

---

## 🧭 Next Phase

### 0.4.21 — Wallet Inventory Response Contract Hardening

Next expected focus:

- preserve backward compatibility of the wallet inventory navigation contract
- only add further inventory semantics if a concrete client need appears
- keep all further enhancements read-only unless the ZIP proves otherwise
- avoid reworking ownership invariants already stabilized in Phase 0.4

---

## 🧩 Summary

At the end of Phase 0.4.16:

- wallet authentication remains stable
- identity remains unified
- ownership remains protected
- authenticated wallet linking is available
- wallet-owned account merge execution is available
- explicit primary-wallet switching is available
- wallet detach eligibility is available under authenticated control
- wallet detach execution is available for already eligible owned wallets
- detached wallet identities are explicitly reusable known identities
- detached wallets can be reattached via protected linking or via wallet-login bootstrap rebound
- detached wallet identities preserve minimal audit-ready lifecycle evidence through `detached_at`
- `GET /auth/wallets` now exposes an enriched lifecycle-aware wallet inventory read model


---

## Phase 0.4.16 — Wallet Identity Read Model Enrichment

### Objective

Expose a richer and more explicit wallet inventory read model through `GET /auth/wallets`, so authenticated clients can observe current ownership plus minimal lifecycle evidence already maintained by the backend.

### Initial Context

By the end of 0.4.15, the backend already supported wallet authentication, authenticated wallet linking, wallet-owned account merge, primary-wallet switching, detach eligibility and execution, detached-wallet reattachment semantics, and minimal detached-identity audit readiness through `detached_at`.

The internal wallet identity model had already become more lifecycle-aware than the public wallet inventory response.

### Problem Statement

The backend already preserved wallet lifecycle fields such as:

- `linked_at`
- `detached_at`
- `is_primary`

However, `GET /auth/wallets` had not yet been explicitly upgraded into a lifecycle-aware read model contract. This created a visibility gap between internal state and client-facing inventory data.

### Scope

Included:

- explicit wallet read-model mapping
- exposure of:
  - `id`
  - `address`
  - `user_id`
  - `linked_at`
  - `detached_at`
  - `is_primary`
  - derived `status`
- handler-level validation for active wallet inventory and detached-then-reattached wallet visibility

Excluded:

- ownership rule changes
- schema changes
- detach/reattach business-rule changes
- filtering, pagination, or reporting expansion

### Root Cause Analysis

The root issue was not missing domain behavior but an outdated API projection. The internal model already tracked richer lifecycle state, while the public inventory endpoint still behaved like a simpler list projection.

### Files Affected

- `internal/modules/auth/http_wallet_list.go`
- `internal/modules/auth/http_handlers_test.go`
- `README.md`
- `docs/phase-status.md`
- `docs/handoff/backend-status.md`
- `docs/phase0_4_auth_and_user_stabilization.md`
- `docs/architecture.md`
- `docs/flows.md`
- `docs/testing.md`

### Implementation Characteristics

The wallet inventory handler now maps wallet identities into an explicit `WalletReadModel` exposing:

- `id`
- `address`
- `user_id`
- `linked_at`
- `detached_at`
- `is_primary`
- `status`

The derived `status` semantics are conservative:

- `active` when the wallet is currently linked to a user
- `detached` when there is no current owner and `detached_at` is present
- `unlinked` when there is no current owner and no detach evidence

For `GET /auth/wallets`, the operational case remains `active`, because the route lists wallets currently owned by the authenticated user. The main value of the enrichment is that previously detached lifecycle evidence remains visible after reattachment.

### Validation

Validated with:

```bash
go test ./...
```

Result:

- `internal/modules/auth` OK
- no visible regressions in the rest of the backend tree

### Release Impact

This subphase is additive and read-oriented. It improves client visibility and debugging without changing authentication, ownership, linking, merge, primary-switch, or detach rules.

### Risks

Low risk. The change is limited to response projection and handler-level contract clarity.

### What it does NOT solve

This subphase does not add:

- filtering
- pagination
- search
- admin reporting
- richer detached-identity history endpoints
- advanced lifecycle analytics

### Conclusion

Phase 0.4.16 closes the gap between the internal wallet identity lifecycle model and the authenticated wallet inventory API contract. The backend now exposes a richer wallet inventory read model while preserving all ownership guarantees stabilized in previous subphases.


## Phase 0.4.17 — Wallet Inventory Query Filtering and Sorting

### Objective

Make `GET /auth/wallets` operationally more useful for clients by adding small, explicit query semantics on top of the lifecycle-aware read model introduced in 0.4.16.

### Initial Context

By the end of 0.4.16, the backend already exposed an explicit wallet read model including:

- `linked_at`
- `detached_at`
- `is_primary`
- `status`

However, the endpoint still behaved like a fixed inventory listing. Clients could observe richer lifecycle state, but could not yet request even basic filtered or ordered views of that same read model.

### Problem Statement

The real gap was not in domain behavior or persistence. The gap was in inventory query semantics.

`GET /auth/wallets` returned the enriched lifecycle-aware projection, but it did not yet support:

- filtering by `status`
- filtering by `primary`
- explicit ordering by `linked_at`

This made the public API less useful for account-management and wallet-inventory UIs even though the underlying read model was already available.

### Scope

Included:

- optional `status` filter with supported values:
  - `active`
  - `detached`
- optional `primary` filter with supported values:
  - `true`
  - `false`
- optional sorting contract:
  - `sort=linked_at`
  - `order=asc|desc`
- strict query validation with explicit `400` errors for unsupported values
- handler-level test coverage for filtering, sorting, and invalid query contracts

Excluded:

- ownership rule changes
- store changes
- SQL query changes
- pagination
- search
- admin reporting
- broader detached-history APIs

### Root Cause Analysis

The root issue remained in the HTTP read layer. The lifecycle-aware model already existed, but the endpoint had no safe query contract for clients to consume that richer projection in a structured way.

### Files Affected

- `internal/modules/auth/http_wallet_list.go`
- `internal/modules/auth/http_handlers_test.go`
- `README.md`
- `docs/phase-status.md`
- `docs/handoff/backend-status.md`
- `docs/phase0_4_auth_and_user_stabilization.md`
- `docs/flows.md`
- `docs/testing.md`

### Implementation Characteristics

The implementation remains read-only and handler-local. It does not change domain or persistence behavior.

`GET /auth/wallets` now accepts these optional query parameters:

- `status=active|detached`
- `primary=true|false`
- `sort=linked_at`
- `order=asc|desc`

Compatibility is preserved:

- when no query parameters are provided, the endpoint keeps the existing inventory behavior
- the existing store-defined default ordering remains the default path
- filtering and sorting are applied only after the owned-wallet inventory has already been resolved

Validation is strict:

- invalid `status` returns `400` with `invalid_status`
- invalid `primary` returns `400` with `invalid_primary`
- invalid `sort` returns `400` with `invalid_sort`
- invalid `order` returns `400` with `invalid_order`
- `order` without a supported `sort` also returns `400` with `invalid_sort`

A practical note remains important: because `GET /auth/wallets` lists wallets currently owned by the authenticated user, `status=detached` is expected to return an empty result under the current contract. That is intentional and keeps the semantics explicit without widening the endpoint scope.

### Validation

Validation path for this subphase:

```
go test ./...
```

Focused coverage added in handler tests for:

- backward-compatible inventory listing without query params
- `primary=true`
- `primary=false`
- `status=active`
- `status=detached` returning an empty result under the current owned-wallet contract
- `sort=linked_at&order=desc`
- invalid query parameters returning `400`

### Release Impact

This subphase is additive and read-oriented. It improves client query semantics without changing wallet auth, ownership, linking, merge, primary switch, detach execution, or detached-wallet reuse rules.

### Risks

Low risk. The change is restricted to the authenticated wallet inventory handler and its tests.

Main guarded risks:

- accidental breakage of existing default ordering
- ambiguous invalid query behavior
- accidental persistence or domain coupling for a read-only enhancement

### What it does NOT solve

This subphase does not add:

- pagination
- search
- detached-wallet history listing
- admin inventory views
- new ownership operations
- domain redesign

### Conclusion

Phase 0.4.17 makes the lifecycle-aware wallet inventory endpoint actually queryable in a small, controlled, backward-compatible way. The backend now exposes a safer and more useful inventory contract while keeping all Phase 0.4 ownership and lifecycle invariants intact.


## Phase 0.4.18 — Wallet Inventory Pagination and Windowed Response

### Objective

Add simple, explicit pagination to `GET /auth/wallets` on top of the filtered and sortable lifecycle-aware inventory contract introduced in 0.4.17.

### Initial Context

By the end of 0.4.17, the authenticated wallet inventory endpoint already supported:

- lifecycle-aware wallet read-model projection
- `status` filtering
- `primary` filtering
- `linked_at` ordering

However, clients still had no way to request a bounded result window or receive explicit metadata about the size of the filtered inventory.

### Problem Statement

The remaining gap was not in ownership or persistence. The gap was in windowed inventory delivery.

`GET /auth/wallets` could already return a filtered and sorted inventory, but it still lacked:

- optional `limit`
- optional `offset`
- explicit response metadata describing the filtered result size and the requested window

This kept the API less practical for inventory UIs that need deterministic partial rendering while preserving the current authenticated ownership scope.

### Scope

Included:

- optional `limit` query param
- optional `offset` query param
- strict validation for invalid pagination inputs
- pagination applied only after filtering and sorting
- additive response metadata:
  - `total`
  - `limit`
  - `offset`
- handler-level tests for valid and invalid pagination scenarios

Excluded:

- cursor pagination
- `has_more` or next-page tokens
- store changes
- SQL pagination
- detached-wallet history APIs
- ownership-rule changes
- new mutation flows

### Root Cause Analysis

The lifecycle-aware read model and basic query semantics already existed. The missing piece was a small read-layer contract for bounded inventory retrieval. This is a handler concern, not a domain or persistence concern.

### Files Affected

- `internal/modules/auth/http_wallet_list.go`
- `internal/modules/auth/http_handlers_test.go`
- `README.md`
- `docs/phase-status.md`
- `docs/handoff/backend-status.md`
- `docs/phase0_4_auth_and_user_stabilization.md`
- `docs/flows.md`
- `docs/testing.md`

### Implementation Characteristics

`GET /auth/wallets` now remains authenticated and ownership-scoped, but also accepts:

- `limit=<positive integer>`
- `offset=<non-negative integer>`

The response now includes additive pagination metadata:

- `wallets`
- `total`
- `limit`
- `offset`

Important behavior:

- filtering still happens first
- sorting still happens second
- pagination happens only after filtering and sorting
- `total` reflects the number of wallets after filtering and sorting, before windowing
- `limit=0` means no explicit page-size cap was requested
- `offset=0` remains the default starting position

Validation is strict:

- invalid `limit` returns `400` with `invalid_limit`
- invalid `offset` returns `400` with `invalid_offset`

This keeps the contract explicit and avoids silently accepting malformed query values.

### Validation

Validation path for this subphase:

```
go test ./...
```

Focused handler coverage now includes:

- backward-compatible wallet inventory listing with metadata defaults
- `limit` only
- `offset` only
- `limit + offset`
- valid empty window when the requested offset exceeds the filtered inventory length
- invalid `limit` values returning `400`
- invalid `offset` values returning `400`

### Release Impact

This subphase is additive and read-only. It improves the authenticated wallet inventory API contract without changing ownership, linking, merge, primary switching, detach execution, or detached-wallet reuse semantics.

### Risks

Low risk. The change remains constrained to the wallet inventory handler, its response contract, and handler-level tests.

Main guarded risks:

- accidental response-contract breakage for existing clients
- incorrect pagination ordering if applied before filtering/sorting
- silent acceptance of malformed pagination values

### What it does NOT solve

This subphase does not add:

- cursor pagination
- text search
- detached-wallet history reporting
- admin inventory views
- store-level pagination
- ownership-rule changes
- new wallet mutation endpoints

### Conclusion

Phase 0.4.18 makes the lifecycle-aware wallet inventory endpoint windowable while preserving its authenticated ownership scope and all Phase 0.4 invariants. The backend now supports small, explicit pagination semantics on top of the queryable wallet inventory contract.


## Phase 0.4.19 — Wallet Inventory Navigation Metadata

### Objective

Add minimal navigation metadata to the paginated `GET /auth/wallets` response so clients can consume windowed inventory results without inferring navigation state on their own.

### Initial Context

By the end of 0.4.18, the authenticated wallet inventory endpoint already supported lifecycle-aware projection, filtering, sorting, and explicit windowing through `limit` and `offset`.

The response already exposed:

- `wallets`
- `total`
- `limit`
- `offset`

However, clients still had to infer whether more results existed and how many items were effectively returned by the current window.

### Problem Statement

The remaining gap was not in ownership, persistence, or query parsing. The gap was in response navigation semantics.

A paginated inventory response without explicit navigation metadata forces each client to replicate simple backend logic to determine:

- how many items the current window actually returned
- whether another page exists after the current window

### Scope

Included:

- additive response metadata:
  - `returned`
  - `has_more`
- deterministic calculation after filtering, sorting, and pagination
- handler-level tests covering default, paginated, empty-window, and filtered-window navigation behavior

Excluded:

- cursor pagination
- `next_offset` or `previous_offset`
- new filters
- store changes
- SQL pagination
- ownership changes

### Root Cause Analysis

Phase 0.4.18 made the wallet inventory endpoint windowable, but the response contract still stopped short of describing the returned window itself. This remained a read-model concern and did not require any domain or persistence expansion.

### Files Affected

- `internal/modules/auth/http_wallet_list.go`
- `internal/modules/auth/http_handlers_test.go`
- `README.md`
- `docs/phase-status.md`
- `docs/handoff/backend-status.md`
- `docs/phase0_4_auth_and_user_stabilization.md`
- `docs/flows.md`
- `docs/testing.md`

### Implementation Characteristics

`GET /auth/wallets` now keeps all existing query semantics and response fields, and adds:

- `returned`
- `has_more`

Important behavior:

- `returned` equals the actual number of wallets present in `wallets` for the current response
- `has_more` is calculated only after filtering, sorting, and applying the requested window
- when `limit=0`, the response is treated as unbounded and `has_more=false`
- the change is additive and backward compatible

### Validation

Validation path for this subphase:

```
go test ./...
```

Focused handler coverage now includes:

- default wallet inventory response with navigation metadata
- paginated windows with `has_more=true`
- final windows with `has_more=false`
- empty valid windows with `returned=0`
- filtered and sorted windows with correct navigation metadata

### Release Impact

This subphase is additive and read-only. It improves the wallet inventory response contract without changing ownership, linking, detach, merge, or primary-wallet semantics.

### Risks

Low risk. The change remains constrained to the wallet inventory handler, response payload, and handler-level tests.

Main guarded risks:

- incorrect `has_more` calculation
- accidental response-contract incompatibility
- ambiguous behavior for unbounded (`limit=0`) requests

### What it does NOT solve

This subphase does not add:

- cursor pagination
- next-page tokens
- search
- detached-wallet history reporting
- admin inventory views
- store-level pagination
- ownership-rule changes

### Conclusion

Phase 0.4.19 completes the basic navigation contract of the authenticated wallet inventory endpoint. The backend now exposes not only filtered, ordered, and windowed wallet inventory responses, but also explicit metadata describing the returned window itself.


## Phase 0.4.20 — Wallet Inventory Cursorless Navigation Hints

### Objective
Add additive cursorless navigation hints to `GET /auth/wallets` so offset-based clients can move forward and backward without recalculating hints outside the API.

### Scope
- add `next_offset` to the wallet inventory response
- add `previous_offset` to the wallet inventory response
- compute navigation hints after filtering, sorting, and pagination
- preserve backward compatibility of the existing wallet inventory contract
- extend handler-level coverage for forward and backward navigation hints

### Delivered
- additive response fields:
  - `next_offset`
  - `previous_offset`
- deterministic offset-based navigation hints for first, intermediate, last, and empty windows
- explicit `nil` navigation hints for unbounded (`limit=0`) responses
- tests covering bounded and unbounded windows, including empty pages and filtered windows

### Validation
- `go test ./...`
- handler-level validation for `next_offset` and `previous_offset` under bounded wallet inventory requests
- handler-level validation for empty and filtered windows with navigation hints

### What it does NOT solve
- cursor pagination
- continuation tokens
- next/previous URLs
- store-level pagination
- additional filters or search
- ownership-rule changes

### Conclusion
Phase 0.4.20 keeps the wallet inventory endpoint read-only and ownership-scoped while completing the basic offset-based navigation contract with explicit forward and backward hints.
