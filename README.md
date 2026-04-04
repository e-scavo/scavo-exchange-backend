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
**Phase:** 0.5 — User Interaction & Application Surface  
**Current Subphase:** **0.5.1 — Authenticated User Profile Surface**

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

### 0.4.15 — Detached Identity Audit Readiness
wallet identities now preserve minimal detached-lifecycle audit metadata:

- `detached_at`

This allows the backend to distinguish a wallet that has never been detached from a wallet identity that was previously detached and later reused through linking or wallet-login rebound.

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
  "user": {
    "id": "..."
  }
}
```

Behavior:
- only `auth_bootstrap` wallet challenges are accepted at this endpoint
- `wallet_link` challenges are rejected with `wallet_challenge_purpose_mismatch`
- `account_merge` challenges are rejected with `wallet_challenge_purpose_mismatch`
- unknown or malformed persisted challenge purposes are rejected and are no longer treated as valid bootstrap challenges
- previously detached known wallets may still rebound through wallet login using a valid `auth_bootstrap` challenge

---

### Password Auth (dev only)

#### `POST /auth/login`

Authenticates a user with email and password for development/testing only.

Request:

```json
{
  "email": "admin@local",
  "password": "admin123"
}
```

Response:

```json
{
  "access_token": "...",
  "token_type": "Bearer",
  "user": {
    "id": "...",
    "email": "admin@local"
  }
}
```

---

### Authenticated Wallet Management

Requires Bearer JWT.

#### `POST /auth/wallets/link/challenge`

Creates a wallet-link confirmation challenge for the currently authenticated user.

Request:

```json
{
  "address": "0x...",
  "chain": "scavium"
}
```

Response:

```json
{
  "challenge_id": "...",
  "message": "..."
}
```

Behavior:
- creates a wallet challenge with purpose `wallet_link`
- records `requested_by_user_id` as the authenticated durable user
- allows later signed attachment of a secondary wallet to the authenticated account
- rejects unknown or malformed challenge purposes during downstream consumption

---

#### `POST /auth/wallets/link/verify`

Consumes a signed wallet-link challenge and attaches the wallet to the authenticated user.

Request:

```json
{
  "challenge_id": "...",
  "address": "0x...",
  "signature": "0x..."
}
```

Behavior:
- requires the authenticated user
- challenge must belong to purpose `wallet_link`
- challenge must have been requested by the same authenticated user
- wallet becomes attached to the authenticated durable user
- detached known wallets can be reattached through this flow
- unknown or malformed persisted challenge purposes are rejected rather than being normalized at runtime

---

#### `POST /auth/account/merge/wallet/challenge`

Creates a wallet-owned account merge challenge for the currently authenticated user.

Request:

```json
{
  "address": "0x...",
  "chain": "scavium"
}
```

Behavior:
- creates a wallet challenge with purpose `account_merge`
- records `requested_by_user_id` as the authenticated durable user
- prepares the authenticated user to absorb a wallet-owned source account only after that source wallet explicitly signs

---

#### `POST /auth/account/merge/wallet/verify`

Consumes a signed merge challenge and merges the source wallet-owned user into the authenticated target user.

Request:

```json
{
  "challenge_id": "...",
  "address": "0x...",
  "signature": "0x..."
}
```

Behavior:
- requires the authenticated user
- challenge must belong to purpose `account_merge`
- challenge must have been requested by the same authenticated user
- source durable wallet-owned account is merged into the authenticated target account
- wallet ownership is re-pointed safely to the target user
- unknown or malformed persisted challenge purposes are rejected rather than being normalized at runtime

---

#### `POST /auth/wallets/primary`

Promotes one owned wallet to primary.

Request:

```json
{
  "wallet_identity_id": "..."
}
```

Behavior:
- requires the authenticated user
- target wallet must already belong to that user
- exactly one owned wallet remains primary after the operation

---

#### `POST /auth/wallets/detach/check`

Evaluates whether one owned wallet is currently detachable.

Request:

```json
{
  "wallet_identity_id": "..."
}
```

Response example:

```json
{
  "eligible": false,
  "reasons": [
    "wallet_is_primary"
  ]
}
```

Behavior:
- requires the authenticated user
- does not mutate ownership
- acts as the authoritative detach-eligibility surface

---

#### `POST /auth/wallets/detach`

Detaches one owned wallet when eligibility constraints allow it.

Request:

```json
{
  "wallet_identity_id": "..."
}
```

Behavior:
- requires the authenticated user
- target wallet must be currently detachable
- clears `user_id`, `linked_at`, and `is_primary`
- stamps `detached_at`
- detached wallet remains a known wallet identity that may later reattach or rebound through wallet login

---

#### `GET /auth/wallets`

Returns the authenticated user wallet inventory.

Supported query params:

- `status=active|detached`
- `primary=true|false`
- `sort=linked_at`
- `order=asc|desc`
- `limit=<positive integer>`
- `offset=<non-negative integer>`

Response example:

```json
{
  "wallets": [
    {
      "id": "...",
      "address": "0xabc...",
      "user_id": "...",
      "linked_at": "2026-03-27T00:00:00Z",
      "detached_at": null,
      "is_primary": true,
      "status": "active",
      "can_set_primary": false,
      "can_detach": false,
      "detach_block_reasons": [
        "wallet_is_primary"
      ]
    },
    {
      "id": "...",
      "address": "0xdef...",
      "user_id": "...",
      "linked_at": "2026-03-27T00:05:00Z",
      "detached_at": null,
      "is_primary": false,
      "status": "active",
      "can_set_primary": true,
      "can_detach": true,
      "detach_block_reasons": []
    }
  ],
  "total": 2,
  "limit": 2,
  "offset": 0,
  "returned": 2,
  "has_more": false,
  "next_offset": null,
  "previous_offset": null
}
```

Behavior:
- returns only wallets currently owned by the authenticated user
- inventory fields are lifecycle-aware and advisory
- `can_set_primary`, `can_detach`, and `detach_block_reasons` are hints, not execution authority
- `POST /auth/wallets/primary`, `POST /auth/wallets/detach/check`, and `POST /auth/wallets/detach` remain authoritative
- `order` requires `sort`
- `sort=linked_at` defaults to ascending order when `order` is omitted
- offset-only requests remain valid and unbounded
- `returned`, `has_more`, `next_offset`, and `previous_offset` make bounded window navigation explicit

---

## 🧪 Minimal Validation Commands

### Health
```bash
curl -i http://localhost:8080/health
```

Expected:

```json
{"ok":true}
```

### Version
```bash
curl -i http://localhost:8080/version
```

### Dev Login
```bash
curl -i -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@local","password":"admin123"}'
```

### Wallet Challenge
```bash
curl -i -X POST http://localhost:8080/auth/wallet/challenge \
  -H "Content-Type: application/json" \
  -d '{"address":"0xabc123","chain":"scavium"}'
```

---

## ⚙️ Environment Variables

Main configuration keys:

- `APP_ENV`
- `PORT`
- `JWT_SECRET`
- `POSTGRES_DSN`
- `REDIS_ADDR`

See `.env.sample` for the full current set.

---

## 🧱 Current Foundation Scope

What Stage 0 / Phase 0.4 currently establishes:

- service bootstrap
- health/version endpoints
- password auth (dev only)
- wallet authentication
- durable user identity
- wallet registry persistence
- wallet ↔ user ownership persistence
- authenticated wallet linking
- authenticated wallet-owned account merge
- explicit primary-wallet switching
- protected wallet detach eligibility and execution
- detached wallet lifecycle clarification
- lifecycle-aware wallet inventory projection
- wallet inventory filtering, ordering, pagination, and navigation metadata
- inventory-side actionability hints aligned with primary / detach endpoints
- strict wallet challenge purpose enforcement across bootstrap, link, and merge flows
- PostgreSQL-backed persistence with in-memory fallback for local/dev usage

---

## 🚫 Out of Scope (for now)

Not yet included:

- refresh tokens
- token revocation
- session persistence
- production-grade password auth
- wallet unlink beyond protected detach semantics
- arbitrary user-to-user ownership transfer
- admin identity tooling
- non-wallet auth providers

---

## 📂 Documentation Map

Detailed docs live under `docs/`:

- `docs/architecture.md`
- `docs/architecture-deep.md`
- `docs/development.md`
- `docs/development-environment.md`
- `docs/testing.md`
- `docs/decisions.md`
- `docs/phase-status.md`
- `docs/phase0_4_auth_and_user_stabilization.md`
- `docs/flows.md`
- `docs/handoff/backend-status.md`

---

## ✅ Current Status

At the end of the current documented state:

- auth service is stable
- durable user identity exists
- wallet ownership is explicit and persistent
- authenticated wallet linking is available
- wallet-owned account merge is available
- primary wallet management is available
- wallet detach eligibility and execution are available
- detached wallet identities are reusable known identities
- wallet inventory is lifecycle-aware and management-oriented
- wallet inventory query semantics are documented and validated
- wallet challenge purpose handling is strict across creation and consumption

---

## Phase 0.4.1 — Auth Base Setup

### Objective

Establish the initial authentication module structure with a working dev login flow, JWT issuance, and protected identity context.

### Delivered

- `POST /auth/login`
- JWT creation and validation
- HTTP auth middleware
- user module and repository abstraction
- protected context injection
- basic auth route registration

### Result

The backend has a minimal but working authentication base suitable for local development and later wallet integration.

---

## Phase 0.4.2 — JWT Implementation and Auth Normalization

### Objective

Normalize the JWT and auth transport layer so the authentication contract is explicit and reusable across future auth methods.

### Delivered

- normalized Bearer parsing
- JWT claims structure cleanup
- shared auth middleware behavior
- protected route consistency

### Result

The backend now has one consistent token transport and one protected-user context model.

---

## Phase 0.4.3 — Auth Endpoints Stabilization

### Objective

Stabilize the initial auth endpoints and close remaining gaps in HTTP responses and login behavior.

### Delivered

- predictable error handling for auth endpoints
- test alignment around login behavior
- route-level auth stabilization

### Result

The password-based dev auth flow is stable and ready to coexist with wallet auth.

---

## Phase 0.4.4 — Wallet Challenge Contract and Nonce Bootstrap

### Objective

Introduce the first wallet-auth contract by allowing the backend to mint signable login challenges.

### Delivered

- `POST /auth/wallet/challenge`
- challenge message generation
- challenge persistence abstraction
- in-memory challenge store
- expiration support
- challenge nonce bootstrap

### Result

The backend can now initiate wallet login by producing signable challenges bound to wallet and chain.

---

## Phase 0.4.5 — Wallet Signature Verification and Token Issuance

### Objective

Complete the first wallet-auth execution path by verifying signatures and issuing JWTs.

### Delivered

- `POST /auth/wallet/verify`
- EVM-style signature verification
- wallet challenge consumption
- token issuance after successful verification
- wallet-auth HTTP tests

### Result

Wallets can now authenticate directly against the backend.

---

## Phase 0.4.6 — Wallet Identity Persistence and Durable Challenge Storage

### Objective

Make wallet authentication durable by persisting known wallet identities and supporting PostgreSQL-backed challenge storage.

### Delivered

- persistent wallet identity store abstraction
- in-memory wallet identity store
- PostgreSQL wallet identity store
- PostgreSQL challenge store
- wallet auth migration
- wallet auth service wiring

### Result

Wallet authentication is no longer ephemeral. Known wallets and challenges can persist across process restarts when PostgreSQL is configured.

---

## Phase 0.4.7 — Wallet ↔ User Linking and Unified Identity Model

### Objective

Unify runtime identity and persistent ownership by explicitly linking wallet identities to durable users.

### Delivered

- wallet identities now persist `user_id`
- wallet auth resolves or creates durable users
- JWT identity is now issued around `user_id`
- wallet identity and session identity become aligned
- ownership-safe linking behavior covered in service tests

### Result

Wallet login now resolves to one durable platform identity instead of treating wallet records as isolated auth artifacts.

---

## Phase 0.4.8 — Multi-Wallet Ownership Foundations

### Objective

Prepare the backend for real wallet ownership management by allowing one durable user to own multiple wallets with explicit ownership metadata.

### Delivered

- wallet identities now persist:
  - `user_id`
  - `linked_at`
  - `is_primary`
- migration for wallet ownership metadata
- wallet auth bootstrap assigns the first wallet as primary
- wallet identity service behavior updated for explicit ownership handling
- in-memory and PostgreSQL stores aligned with the new ownership model

### Result

The backend now supports one durable user owning multiple wallets while preserving primary-wallet semantics and explicit ownership persistence.

---

## Phase 0.4.9 — Authenticated Wallet Linking Contract

### Objective

Allow an already authenticated durable user to attach an additional wallet through an explicit signed flow instead of relying only on wallet login heuristics.

### Delivered

- `POST /auth/wallets/link/challenge`
- `POST /auth/wallets/link/verify`
- wallet-link challenge purpose support
- `requested_by_user_id` persisted with linking challenges
- authenticated linking service flow
- HTTP coverage for successful wallet linking

### Result

Secondary-wallet attachment is now a protected authenticated flow that preserves durable ownership semantics without creating a new account or relying on implicit merge rules.

---

## Phase 0.4.10 — User-Driven Wallet-Owned Account Merge Execution

### Objective

Allow the authenticated user to safely absorb another wallet-owned account only after the source wallet explicitly signs a merge challenge.

### Delivered

- `POST /auth/account/merge/wallet/challenge`
- `POST /auth/account/merge/wallet/verify`
- merge challenge purpose support
- authenticated merge flow
- ownership-safe reassignment of wallet identities from source user to target user
- source user deletion after merge
- HTTP and service coverage for successful merge execution

### Result

Wallet-owned account merge is now explicit, wallet-signed, and controlled by the authenticated target user rather than inferred implicitly during login.

---

## Phase 0.4.11 — Explicit Primary-Wallet Switching

### Objective

Allow the authenticated user to explicitly promote one owned wallet to primary without changing ownership.

### Delivered

- `POST /auth/wallets/primary`
- ownership-safe primary switch flow
- service coverage for switching across two owned wallets
- HTTP validation coverage for the primary-switch request path

### Result

Users can now explicitly control which owned wallet is primary while keeping the one-user/multi-wallet ownership model intact.

---

## Phase 0.4.12 — Wallet Detach Eligibility Contract

### Objective

Introduce a safe detach-evaluation flow so one authenticated user can know whether one owned wallet may currently be detached without executing the detach.

### Delivered

- `POST /auth/wallets/detach/check`
- detach eligibility evaluation logic
- explicit block reasons:
  - `wallet_is_primary`
  - `user_would_have_no_wallets`
- service coverage for protected detach constraints
- HTTP validation coverage for detach-check

### Result

The backend can now answer whether one owned wallet is detachable while preserving ownership safety and primary-wallet invariants.

---

## Phase 0.4.13 — Protected Wallet Detach Execution

### Objective

Allow the authenticated user to detach one owned wallet only when the detach-safety rules are satisfied.

### Delivered

- `POST /auth/wallets/detach`
- detach execution flow reusing the existing safety checks
- wallet identity ownership clearing:
  - `user_id = nil`
  - `linked_at = nil`
  - `is_primary = false`
- service coverage for successful and rejected detach execution
- HTTP coverage for detach execution behavior

### Result

One owned secondary wallet can now be detached safely without breaking the invariant that the original user must retain at least one owned wallet.

---

## Phase 0.4.14 — Detached Wallet Reattachment Semantics and Lifecycle Clarification

### Objective

Clarify what a detached wallet becomes and formalize whether it can be reused later without introducing a new lifecycle table or terminal state model.

### Delivered

- detached wallets remain known wallet identities
- detached wallet records are preserved after detach
- detached wallets can be reattached through the authenticated link flow
- detached wallets can also re-enter through wallet-login bootstrap
- documentation and tests aligned with reusable detached-wallet semantics

### Result

Detached wallets are now explicitly modeled as reusable known identities rather than deleted or terminal records.

---

## Phase 0.4.15 — Detached Identity Audit Readiness

### Objective

Close the first detached-identity audit gap by preserving minimal lifecycle evidence on wallet identities without introducing heavy event history or lifecycle redesign.

### Delivered

- `detached_at` added to wallet identities
- PostgreSQL migration for detached audit readiness
- in-memory and PostgreSQL stores updated to persist `detached_at`
- detach execution now stamps `detached_at`
- reattachment and login rebound preserve detached audit metadata
- tests aligned with reusable detached-wallet semantics plus minimal audit readiness

### Result

The backend now retains minimal evidence that one wallet identity was previously detached even if that wallet is later reused.

---

## Phase 0.4.16 — Wallet Identity Read Model Enrichment

### Objective

Expose a richer lifecycle-aware wallet inventory read model so authenticated clients can reason about current wallet ownership and minimal detached-lifecycle evidence.

### Delivered

- `GET /auth/wallets`
- explicit wallet read model with:
  - `id`
  - `address`
  - `user_id`
  - `linked_at`
  - `detached_at`
  - `is_primary`
  - `status`
- lifecycle-aware inventory serialization
- tests covering active and detached-history wallet visibility

### Result

Phase 0.4.16 closes the gap between the internal wallet identity lifecycle model and the authenticated wallet inventory API contract. The backend now exposes a richer wallet inventory read model while preserving all ownership guarantees stabilized in previous subphases.

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

Phase 0.4 is now formally closed.

Expected next focus:

- start a new phase only when the next ZIP shows a real runtime, product, or documentation need outside the already stabilized Phase 0.4 scope
- preserve backward compatibility of the authenticated wallet inventory and wallet-action endpoints
- avoid reopening Phase 0.4 unless a future ZIP proves a real regression or contractual gap

---

## 🧩 Summary

At the end of Phase 0.4:

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

## Phase 0.4.17 — Wallet Inventory Query Filtering and Sorting

### Objective

Expose the authenticated wallet inventory as a small but explicitly queryable read surface without changing ownership rules or persistence behavior.

### Delivered

- optional `status=active|detached` filter on `GET /auth/wallets`
- optional `primary=true|false` filter on `GET /auth/wallets`
- optional `sort=linked_at` with `order=asc|desc`
- strict HTTP `400` handling for unsupported query values
- tests covering filtering, sorting, and invalid query parameters

### Result

Authenticated clients can now query owned-wallet inventory more precisely while the backend keeps ownership rules and wallet lifecycle behavior unchanged.

---

## Phase 0.4.18 — Wallet Inventory Pagination and Windowed Response

### Objective

Add simple pagination semantics to the authenticated wallet inventory contract without changing stores, ownership rules, or lifecycle behavior.

### Delivered

- optional `limit`
- optional `offset`
- response metadata:
  - `total`
  - `limit`
  - `offset`
- strict HTTP `400` handling for malformed pagination values
- tests covering valid and invalid paginated requests

### Result

The wallet inventory endpoint now supports bounded windows while preserving the same lifecycle-aware wallet read model introduced previously.

---

## Phase 0.4.19 — Wallet Inventory Navigation Metadata

### Objective

Complete the paginated wallet inventory contract with additive navigation metadata.

### Delivered

- additive response metadata:
  - `returned`
  - `has_more`
- deterministic post-filter/post-sort window calculations
- tests covering default, paginated, filtered, and empty-window scenarios

### Result

Wallet inventory responses now expose enough metadata for clients to reason about current page shape without introducing cursor-based complexity.

---

## Phase 0.4.20 — Wallet Inventory Cursorless Navigation Hints

### Objective

Clarify next/previous bounded-window navigation without introducing cursor semantics or changing the underlying pagination strategy.

### Delivered

- additive response metadata:
  - `next_offset`
  - `previous_offset`
- bounded-window navigation hints aligned with current `limit` / `offset`
- test coverage for navigation-hint behavior across different windows

### Result

Clients can now follow straightforward offset-based inventory navigation while the backend remains intentionally cursorless.

---

## Phase 0.4.21 — Wallet Inventory Query Parameter Contract Hardening

### Objective

Harden the query-parameter contract for authenticated wallet inventory without introducing new filters, new sorting fields, or runtime ownership changes.

### Delivered

- `order` now requires an explicit `sort`
- `sort=linked_at` defaults to ascending order when `order` is omitted
- offset-only requests remain valid and unbounded
- tests covering the hardened combinations and defaults

### Result

The wallet inventory query contract is now explicit, stricter, and less ambiguous while remaining backward-compatible for valid callers.

---

## Phase 0.4.22 — Wallet Inventory Response Contract Clarification

### Objective

Clarify the visible wallet inventory API contract so documentation matches the response behavior already implemented in previous subphases.

### Delivered

- `GET /auth/wallets` response examples now include:
  - `returned`
  - `has_more`
  - `next_offset`
  - `previous_offset`
- explicit documentation of bounded vs unbounded window semantics
- response-field contract alignment across docs

### Result

Phase 0.4.22 closes the documentation gap around the wallet inventory response contract. The endpoint behavior remains unchanged, but the visible API contract is now explicit and aligned with the implementation.

---

## Phase 0.4.23 — Wallet Inventory Query Examples Closure

### Objective

Close the examples layer around `GET /auth/wallets` so implementers can see concrete valid and invalid request patterns.

### Delivered

- base request example
- filtered request example
- sorted request example
- paginated request example
- invalid `order` without `sort` example
- bounded-window response examples aligned with the existing handler contract

### Result

Phase 0.4.23 closes the remaining examples gap around `GET /auth/wallets` by documenting concrete request and response patterns without changing domain, stores, persistence, or handler behavior.

---

## Phase 0.4.24 — Wallet Inventory Manual Validation Closure

### Objective

Close the manual-validation layer around `GET /auth/wallets` so operators have an explicit checklist for validating the already-implemented query contract end-to-end.

### Delivered

- manual validation steps for:
  - base inventory requests
  - filtered inventory requests
  - sorted inventory requests
  - paginated inventory requests
  - invalid contract requests
- explicit checks for:
  - `returned`
  - `has_more`
  - `next_offset`
  - `previous_offset`

### Result

Phase 0.4.24 closes the manual-validation layer around `GET /auth/wallets` by documenting how to verify the existing contract end-to-end without changing domain, stores, persistence, or handler behavior.

---

## Phase 0.4.25 — Wallet Actionability Read Model Preparation

### Objective

Expose minimal wallet-management actionability hints inside authenticated wallet inventory without changing detach rules, primary-switch rules, stores, or persistence.

### Delivered

- additive wallet inventory fields:
  - `can_set_primary`
  - `can_detach`
  - `detach_block_reasons`
- detach block reasons aligned with detach-check semantics:
  - `wallet_is_primary`
  - `user_would_have_no_wallets`
- tests covering inventory actionability across single-wallet and multi-wallet ownership scenarios

### Result

Authenticated clients can now observe advisory wallet-management hints directly from inventory while the authoritative execution and eligibility endpoints remain unchanged.

---

## Phase 0.4.26 — Wallet Detach Check Read Consistency

### Objective

Align inventory-side detach hints with the authoritative detach-check endpoint without changing detach rules, stores, or persistence.

### Delivered

- tests proving:
  - `can_detach=false` remains aligned with `eligible=false`
  - `can_detach=true` remains aligned with `eligible=true`
  - detach block reasons remain coherent across inventory and detach-check
- documentation alignment around advisory vs authoritative detach semantics

### Result

Inventory-side detach hints and `POST /auth/wallets/detach/check` now describe the same detach-eligibility reality while keeping authority at the detach-check endpoint.

---

## Phase 0.4.27 — Wallet Primary Switch Read Consistency

### Objective

Align inventory-side primary actionability hints with the authoritative primary-switch endpoint without changing primary-switch rules, stores, or persistence.

### Delivered

- tests proving:
  - the current primary remains `can_set_primary=false`
  - an eligible secondary wallet remains `can_set_primary=true`
  - promoted wallets become non-promotable after the switch
- documentation alignment around advisory vs authoritative primary semantics

### Result

Inventory-side primary hints and `POST /auth/wallets/primary` now describe one coherent primary-management contract while keeping authority at the execution endpoint.

---

## Phase 0.4.28 — Wallet Management Read Flow Closure

### Objective

Close the operational read flow around wallet inventory, primary switching, and detach operations without changing domain rules, stores, or persistence.

### Delivered

- wallet-management flow documentation:
  - inventory
  - actionability hint
  - action/check endpoint
  - refreshed inventory
- README header alignment with the real subphase state
- manual validation guidance for refreshed inventory after wallet-management actions

### Result

Phase 0.4.28 does not change wallet-management rules, handlers, stores, or persistence. It closes the operational read flow around the existing wallet inventory and action endpoints so client and operator guidance now matches the real authenticated wallet-management surface end to end.

---

## Phase 0.4.29 — Wallet Detach Execute Read Consistency

### Objective

Align inventory-side detach hints with the authoritative detach-execution endpoint without changing detach rules, stores, or persistence.

### Delivered

- tests proving:
  - a detachable secondary wallet can be detached successfully
  - the detach execute response remains compatible with pre-detach inventory hints
  - refreshed inventory recalculates detach hints coherently after detach
- documentation alignment around advisory inventory vs authoritative execution semantics

### Result

Inventory-side detach hints and `POST /auth/wallets/detach` now describe one coherent detach-management contract while keeping authority at the execution endpoint.

---

## Phase 0.4.30 — Wallet Management Contract Consolidation

### Objective

Consolidate the authenticated wallet-management surface into one explicit contract without changing handlers, stores, persistence, or domain rules.

### Delivered

- explicit consolidated wallet-management contract spanning:
  - `GET /auth/wallets`
  - `POST /auth/wallets/primary`
  - `POST /auth/wallets/detach/check`
  - `POST /auth/wallets/detach`
- explicit clarification that:
  - inventory is the advisory read surface
  - `detach/check` is the eligibility surface
  - `primary` and `detach` are execution surfaces
  - refreshed inventory is the post-action observable state
- unified testing guidance for inventory-driven wallet-management validation
- cross-document cleanup so README, flows, testing, handoff, and phase status describe the same wallet-management contract

### Result

Phase 0.4 now closes with wallet management described as one consolidated, inventory-driven contract without changing handlers, stores, persistence, or ownership rules.

---

## Phase 0.4.31 — Wallet Auth Bootstrap Purpose Enforcement

### Objective

Harden the wallet-auth bootstrap contract so `POST /auth/wallet/verify` only consumes `auth_bootstrap` challenges.

### Delivered

- service-level purpose enforcement in wallet verify/login
- explicit rejection of `wallet_link` challenges during wallet-auth bootstrap
- explicit rejection of `account_merge` challenges during wallet-auth bootstrap
- handler-level conflict response with `wallet_challenge_purpose_mismatch`
- test coverage proving the service and HTTP boundaries reject non-bootstrap challenges

### Result

The wallet lifecycle remains unchanged:

- authenticated linking still uses `wallet_link`
- wallet-owned account merge still uses `account_merge`
- wallet login / rebound after detach still uses `auth_bootstrap`

What changes is the contract enforcement: these challenge purposes are no longer interchangeable at the wallet-auth bootstrap boundary.

---

## Phase 0.4.32 — Wallet Challenge Purpose Strictness Closure

### Objective

Close the remaining permissive purpose-normalization gap so wallet challenges no longer silently degrade unknown or malformed purpose values into `auth_bootstrap` during runtime consumption.

### Delivered

- strict purpose resolution for controlled challenge creation
- explicit rejection of unknown challenge purposes at wallet verify/login
- explicit rejection of unknown challenge purposes at authenticated wallet link
- explicit rejection of unknown challenge purposes at authenticated wallet-owned account merge
- runtime loading now preserves unknown persisted purpose values instead of reclassifying them as bootstrap
- test coverage proving creation defaults remain controlled while unknown runtime purposes are rejected

### Result

Phase 0.4 closes with strict challenge-purpose handling across creation and consumption:

- controlled creation still defaults empty purpose to `auth_bootstrap`
- supported purposes remain `auth_bootstrap`, `wallet_link`, and `account_merge`
- unknown or malformed purposes are no longer treated as valid bootstrap challenges at runtime

---

## Phase 0.4.33 — Phase 0.4 Formal Closure

### Objective

Formally close Phase 0.4 at the documentation layer now that wallet auth, ownership, wallet management, and challenge-purpose enforcement are already stabilized in the implementation and reflected across the ZIP.

### Delivered

- explicit formal closure of Phase 0.4 as a completed foundation phase
- README alignment so the declared current subphase matches the final documented Phase 0.4 state
- explicit transition guidance that the next work must begin in a new phase unless a future ZIP proves a real Phase 0.4 regression or documentation gap
- no runtime, store, persistence, migration, or handler changes

### Result

Phase 0.4 is now formally closed. The backend keeps the already stabilized contracts for:

- wallet auth bootstrap
- wallet ↔ user linking
- wallet-owned account merge
- primary wallet switching
- detach eligibility and detach execution
- lifecycle-aware wallet inventory
- strict wallet challenge purpose handling

Any future continuation must start from a new phase rather than extending Phase 0.4 without a new ZIP-validated need.


## ✅ Phase 0.5.1 Closure Summary

Phase 0.5.1 opens the first application-facing user surface on top of the identity and ownership work completed in Phase 0.4.

`GET /auth/me` now remains backward compatible through the existing `user` field while also returning an additive `profile` object that summarizes:

- authenticated durable user identity
- auth method and current wallet-backed session context
- primary wallet summary when present
- owned wallet list projection suitable for application bootstrap
- aggregated wallet counters

This keeps `/auth/session` focused on raw authenticated session claims and `/auth/wallets` focused on the fuller inventory contract, while giving clients one small surface for authenticated application bootstrap.
