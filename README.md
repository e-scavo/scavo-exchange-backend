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
**Latest Completed Phase:** **0.14 — Observability & Diagnostics Foundation**  
**Latest Completed Subphase:** **0.15.5 — Contract Freeze Enforcement**  
**Phase Status:** **0.15.5 Completed / 0.15.6 Pending**  
**Current Phase:** **0.15 — Contract Hardening & Freeze**  
**Current Subphase:** **0.15.6 — Validation & Documentation**

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

#### `user_settings`
stores authenticated user settings separately from core durable user identity:

- `user_id`
- `preferences`
- `created_at`
- `updated_at`

Used for:
- authenticated settings contract bootstrap
- future user preference evolution without coupling settings to the `users` core record

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
- source durable wallet-owned account is merged into the authenticated target user
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

## 👤 Authenticated User Surface

Requires Bearer JWT.

#### `GET /auth/me`

Returns the authenticated durable user and the additive profile view.

Response example:

```json
{
  "user": {
    "id": "...",
    "email": "admin@local",
    "display_name": "SCAVO Operator"
  },
  "profile": {
    "auth_method": "wallet",
    "primary_wallet_address": "0xabc..."
  }
}
```

Behavior:
- acts as the authenticated profile bootstrap surface
- does not mutate user state
- remains separate from wallet-management execution endpoints

---

#### `PATCH /auth/me`

Updates minimal authenticated non-wallet metadata.

Request:

```json
{
  "display_name": "SCAVO Operator"
}
```

Behavior:
- currently only allows `display_name`
- trims input before validation
- rejects empty values after trim
- enforces a Unicode-aware maximum length
- returns the same `user` + `profile` shape as `GET /auth/me`

---

#### `GET /auth/me/settings`

Returns the authenticated user settings contract.

Response example:

```json
{
  "settings": {
    "user_id": "...",
    "version": 1,
    "preferences": {},
    "created_at": "2026-04-06T22:00:00Z",
    "updated_at": "2026-04-06T22:15:00Z"
  }
}
```

Behavior:
- is authenticated and independent from wallet lifecycle
- is separate from `GET /auth/me` so profile metadata and settings do not collapse into one surface
- resolves persisted settings when available
- returns safe defaults when no settings row exists yet
- does not create or mutate settings as a side effect of reading
- exposes `created_at` and `updated_at` when persisted settings metadata is available
- omits timestamp fields when the settings resource is still being resolved from safe defaults without persistence metadata

---

#### `GET /auth/session`

Returns the authenticated session projection resolved from JWT claims and durable identity context.

Behavior:
- is authenticated
- exposes session-oriented auth context
- does not change durable identity, profile metadata, or settings

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
- authenticated profile bootstrap via `GET /auth/me`
- minimal authenticated profile metadata update via `PATCH /auth/me`
- dedicated authenticated settings contract via `GET /auth/me/settings`
- authenticated session read surface via `GET /auth/session`

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
- settings mutation (`PATCH /auth/me/settings` or equivalent)
- concrete user preference fields beyond the minimal settings contract foundation

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
- `docs/phase0_5_user_interaction_and_application_surface.md`
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
- authenticated profile bootstrap is available
- minimal authenticated profile metadata update is available
- authenticated user settings contract bootstrap is available
- authenticated session read surface is available

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

Clarify the lifecycle meaning of detached wallet identities so detach becomes a reversible ownership change rather than an archived or terminal wallet state.

### Delivered

- explicit detached-wallet lifecycle semantics documented and enforced by the already existing ownership model
- clarification that detach keeps the wallet identity record and address known to the system
- clarification that detach clears ownership metadata while preserving future reusability
- reattachment and rebound guidance aligned with the existing wallet-link and wallet-login flows
- cross-document contract cleanup for detach lifecycle behavior

### Result

Detached wallet identities are now explicitly treated as reusable known wallet identities. The backend does not archive, destroy, or permanently blacklist detached wallets as part of detach.

---

## Phase 0.4.15 — Detached Identity Audit Readiness

### Objective

Preserve minimal detached-lifecycle audit metadata so the system can distinguish wallets that were never detached from wallets that were previously detached and later reused.

### Delivered

- `detached_at` persisted on wallet identities
- detach execution stamps `detached_at`
- wallet reattachment clears `detached_at`
- wallet-login rebound clears `detached_at`
- in-memory and PostgreSQL stores aligned with detached audit semantics
- tests covering detached lifecycle audit transitions

### Result

The backend now preserves minimal detached-lifecycle history while keeping the wallet lifecycle model simple and reversible.

---

## Phase 0.4.16 — Authenticated Wallet Inventory Base Contract

### Objective

Expose the first authenticated wallet inventory read surface without changing wallet ownership rules, linking behavior, or lifecycle semantics.

### Delivered

- `GET /auth/wallets`
- authenticated wallet inventory response scoped to the current durable user
- wallet inventory base response fields:
  - `id`
  - `address`
  - `user_id`
  - `linked_at`
  - `detached_at`
  - `is_primary`
  - `status`
- service and HTTP coverage for wallet inventory bootstrap behavior

### Result

Authenticated clients can now query their owned wallet inventory through a stable read surface while wallet management remains on the existing action endpoints.

---

## Phase 0.4.17 — Wallet Inventory Query Filtering Foundations

### Objective

Prepare the authenticated wallet inventory for real client consumption by introducing bounded query filtering without changing ownership, lifecycle, or wallet-management rules.

### Delivered

- `status=active|detached`
- `primary=true|false`
- filter validation and normalization
- HTTP coverage for valid and invalid inventory filter requests
- no changes to wallet identity ownership rules or persistence model

### Result

Wallet inventory can now be narrowed through explicit query filters while preserving the same underlying ownership and lifecycle semantics.

---

## Phase 0.4.18 — Wallet Inventory Ordering Contract

### Objective

Allow authenticated wallet inventory results to be returned in explicit and predictable order without changing ownership or lifecycle semantics.

### Delivered

- `sort=linked_at`
- `order=asc|desc`
- validation that `order` requires `sort`
- deterministic linked-at ordering behavior
- HTTP and service coverage for ordered inventory responses

### Result

Authenticated wallet inventory now supports predictable ordering, making it safer for client rendering and manual operational inspection.

---

## Phase 0.4.19 — Wallet Inventory Pagination Contract

### Objective

Make wallet inventory consumable in bounded windows without changing ownership, lifecycle, or existing management endpoints.

### Delivered

- `limit`
- `offset`
- response pagination fields:
  - `total`
  - `limit`
  - `offset`
  - `returned`
  - `has_more`
  - `next_offset`
  - `previous_offset`
- validation for invalid pagination inputs
- HTTP and service coverage for paginated inventory behavior

### Result

Authenticated wallet inventory now supports offset pagination with explicit navigation metadata.

---

## Phase 0.4.20 — Wallet Inventory Read Model Stabilization

### Objective

Stabilize the authenticated wallet inventory read contract so filtering, ordering, pagination, and lifecycle projection behave as one coherent surface.

### Delivered

- inventory read-model consistency cleanup
- test coverage for combined query scenarios
- documentation alignment across inventory filtering, ordering, and pagination semantics

### Result

`GET /auth/wallets` now behaves as a stable, composable authenticated read surface for wallet ownership inspection.

---

## Phase 0.4.21 — Wallet Inventory Operational Documentation Closure

### Objective

Close the documentation and operational guidance layer around the authenticated wallet inventory contract without changing the runtime implementation.

### Delivered

- manual validation guidance for inventory requests
- documentation alignment across README, flows, testing, and handoff status
- no runtime or persistence changes

### Result

The wallet inventory contract is now documented not only as code behavior but also as an operator-facing authenticated read surface.

---

## Phase 0.4.22 — Wallet Inventory Navigation Semantics Closure

### Objective

Clarify and document the practical meaning of pagination navigation fields so clients can consume wallet inventory windows without ambiguity.

### Delivered

- explicit semantics for:
  - `returned`
  - `has_more`
  - `next_offset`
  - `previous_offset`
- cross-document pagination navigation guidance
- no runtime or persistence changes

### Result

Wallet inventory pagination now has a fully documented navigation contract suitable for client and operator use.

---

## Phase 0.4.23 — Wallet Inventory Manual Validation Closure

### Objective

Complete the manual validation guidance for the authenticated wallet inventory contract without changing handlers, stores, or persistence.

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

Phase 0.4.23 closes the manual-validation layer around `GET /auth/wallets` by documenting how to verify the existing contract end to end without changing domain, stores, persistence, or handler behavior.

---

## Phase 0.4.24 — Wallet Inventory Contract Validation Closure

### Objective

Close the contract-validation layer around the authenticated wallet inventory surface by aligning code, tests, and documentation around the already delivered read behavior.

### Delivered

- explicit README and testing guidance for inventory validation
- no runtime or persistence changes
- cross-document cleanup for wallet inventory contract wording

### Result

The authenticated wallet inventory contract is now closed not only at the runtime layer but also at the validation and documentation layers.

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

## ✅ Phase 0.5.2 Closure Summary

Phase 0.5.2 adds the first authenticated mutation on top of the application-facing user surface opened in Phase 0.5.1.

`PATCH /auth/me` now allows the authenticated user to update only `display_name` as minimal non-wallet metadata. The response remains aligned with `GET /auth/me` by returning both:

- the updated top-level `user` payload
- the additive `profile` object introduced in 0.5.1

This keeps user metadata editing intentionally small and safe:

- no email mutation
- no wallet lifecycle changes
- no settings contract yet
- no new ownership rules

The backend now supports both profile bootstrap (`GET /auth/me`) and minimal authenticated profile metadata update (`PATCH /auth/me`) without reopening identity design from Phase 0.4.

## ✅ Phase 0.5.3 Closure Summary

Phase 0.5.3 introduces the first dedicated authenticated user settings contract without reopening identity, wallet ownership, or metadata editing semantics.

`GET /auth/me/settings` now provides a separate authenticated surface for user settings with this minimal contract:

- `user_id`
- `version`
- `preferences`

This keeps the application surface intentionally separated:

- `GET /auth/me` remains the profile bootstrap surface
- `PATCH /auth/me` remains limited to minimal non-wallet metadata
- `GET /auth/me/settings` becomes the dedicated settings read surface

The settings foundation also remains intentionally small and safe:

- settings persistence is separated from `users`
- defaults are returned when no settings row exists yet
- no implicit write happens during read
- no concrete preference fields are forced prematurely
- no settings mutation contract is introduced yet

The backend now supports a three-part authenticated user surface:

- profile bootstrap
- minimal profile metadata update
- dedicated settings contract bootstrap

This establishes the first formal boundary between durable user metadata and future user configuration.

---

## Phase 0.5.3 — User Settings Contract Foundation

This subphase introduces the first dedicated authenticated user settings contract, extending the application surface without modifying identity, wallet ownership, or authentication flows.

### What was introduced

- new authenticated endpoint: `GET /auth/me/settings`
- dedicated `user_settings` persistence, separated from `users`
- minimal settings response contract:
  - `user_id`
  - `version`
  - `preferences`

### Behavior

- settings are resolved from authenticated context
- persisted settings are returned when available
- safe defaults are returned when no settings row exists
- no implicit insert or mutation occurs during read

### Architectural impact

This establishes a strict separation between:

- identity (Phase 0.4)
- profile metadata (`/auth/me`, `PATCH /auth/me`)
- user configuration (`/auth/me/settings`)

The `/auth/me` endpoint remains focused on authenticated profile/bootstrap concerns and does not absorb settings responsibilities.

### Notes

- settings are currently read-only
- no typed preference fields are enforced yet
- future evolution should extend settings through a dedicated contract, not through `/auth/me`

---

## Phase 0.5.4 — User Settings Mutation

The backend now supports **authenticated mutation of user settings** through a controlled and safe partial update mechanism.

### New Endpoint

```
PATCH /auth/me/settings
```

### Characteristics

* Requires JWT authentication
* Fully independent from `/auth/me` (profile)
* Uses **merge-based updates** (no full overwrite)
* Preserves existing settings
* Prevents destructive writes

### Request Format

```json
{
  "preferences": {
    "key": "value"
  }
}
```

### Behavior

* Only provided keys are updated
* Existing keys remain unchanged
* New keys are added
* Null values are rejected

### Persistence

* Backed by `user_settings`
* Upsert behavior:

  * Creates record if not exists
  * Merges and updates if exists
* Maintains `created_at`
* Updates `updated_at`

### Validation

* Ensures valid JSON
* Requires `preferences` object
* Rejects invalid types
* No schema enforcement (intentionally minimal)

### Compatibility

* Fully backward compatible
* Does not modify existing endpoints
* Does not impact authentication or wallet lifecycle

### Purpose

This phase completes the **settings contract** by transitioning from:

```
READ ONLY → READ + WRITE
```

Enabling:

* Frontend persistence of user preferences
* Future feature-level customization
* Progressive enhancement without architectural changes

---


## Phase 0.5.5.3 — User Settings Contract Surface Stabilization

Phase 0.5.5.3 closes the next contract gap in authenticated user settings by making the resource surface more explicit without introducing schema-heavy governance or changing the endpoint envelope.

Delivered in this subphase:

- `usersettings.View` now exposes `created_at` and `updated_at` when persisted metadata exists
- `GET /auth/me/settings` and `PATCH /auth/me/settings` now return the same resource-oriented settings view with stable timestamp visibility
- zero-value timestamps remain omitted so default, non-persisted settings resolution does not fabricate persistence metadata
- HTTP-level tests now assert both timestamp omission and timestamp exposure behavior

This keeps the authenticated settings contract backward compatible while making the resource state more self-descriptive for frontend consumers.


---

## Phase 0.6 — Authenticated Application Bootstrap Consolidation & Session-Ready Surface

Phase 0.6 opens the next authenticated surface consolidation layer by treating the already existing authenticated endpoints as part of a single application bootstrap boundary that must remain semantically explicit, contract-safe, and frontend-ready without introducing new business domains.

This phase is intentionally limited to authenticated application bootstrap semantics and does not introduce billing, trading, payment, settlement, or any other business-facing feature surface.

The authenticated bootstrap layer covered by Phase 0.6 is composed of:

- `GET /auth/me`
- `GET /auth/session`
- `GET /auth/me/settings`
- `GET /auth/wallets`

The purpose of this phase is to ensure that these endpoints can be consumed together as a coherent authenticated application entry surface while preserving their individual responsibilities.

---

## Phase 0.6.1 — Bootstrap Surface Boundary Clarification

Phase 0.6.1 clarifies the semantic boundaries of the authenticated application bootstrap surface without modifying the production handlers, endpoint shapes, persistence model, or JWT lifecycle.

Before this subphase, the authenticated endpoints were already functional, but the separation between identity bootstrap semantics and token-derived session semantics remained mostly implicit in the code and tests.

This subphase makes that separation explicit and contractual.

### Boundary clarified in 0.6.1

#### `GET /auth/me`

`/auth/me` is the authenticated **bootstrap identity surface**.

Its responsibility is to expose the authenticated user read model together with profile-level authenticated context such as:

- durable user identity
- authenticated profile projection
- wallet ownership summary
- wallet session awareness
- profile-side wallet counters and primary wallet context

`/auth/me` is **not** the source of truth for token metadata such as issuer, subject, token type, authenticated session status, or expiration metadata.

#### `GET /auth/session`

`/auth/session` is the authenticated **session surface**.

Its responsibility is to expose token-derived session context, including:

- authenticated flag
- token type
- issuer
- subject
- expiration metadata
- session-level wallet claim context
- attached user summary for convenience

`/auth/session` is **not** the endpoint for wallet inventory, wallet counters, primary-wallet bootstrap projection, or profile-side ownership summarization.

#### `GET /auth/me/settings`

`/auth/me/settings` remains the authenticated **settings surface**.

Its responsibility is limited to user settings retrieval and mutation contract continuity.

It does not become a profile endpoint, session endpoint, or wallet inventory endpoint as part of Phase 0.6.1.

#### `GET /auth/wallets`

`/auth/wallets` remains the authenticated **wallet inventory surface**.

Its responsibility is to expose the detailed wallet ownership read model and wallet-level actionability/inventory metadata.

It is not promoted into a generic authenticated session or profile bootstrap endpoint.

### Delivered in this subphase

Delivered in Phase 0.6.1:

- explicit authenticated bootstrap boundary clarification between `/auth/me` and `/auth/session`
- explicit preservation of `/auth/me/settings` as settings-only surface
- explicit preservation of `/auth/wallets` as wallet inventory surface
- contract-level HTTP tests that assert the separation between identity bootstrap payload concerns and session payload concerns
- preservation of full backward compatibility for all authenticated endpoints

### Test-level enforcement added in 0.6.1

The subphase is enforced through authenticated HTTP handler tests that now explicitly verify:

- `/auth/me` returns `user` and `profile`
- `/auth/me` does not return `session`
- `/auth/me` profile payload does not absorb session-only fields such as `authenticated`, `token_type`, `issuer`, `subject`, or `expires_at`
- `/auth/session` returns `session`
- `/auth/session` does not return `profile`
- `/auth/session` preserves session-only fields such as `authenticated`, `token_type`, `issuer`, `subject`, and `expires_at`
- `/auth/session` does not absorb profile-side wallet inventory counters or wallet ownership summary fields

### Compatibility

Phase 0.6.1 is fully backward compatible.

No production handler behavior was redesigned in this subphase.

Specifically:

- no endpoint was renamed
- no payload envelope was changed
- no settings schema was introduced
- no wallet lifecycle behavior was modified
- no JWT issuance or validation behavior was changed
- no migration was added

### Why this subphase matters

This clarification is necessary because the frontend-authenticated application bootstrap layer must be able to distinguish:

- who the authenticated user is
- what the current token/session says
- what the persisted user settings are
- what wallets belong to the authenticated user

Without ambiguity about which endpoint owns which responsibility.

Phase 0.6.1 closes that semantic ambiguity and prepares the backend for the next consolidation steps:

- contract alignment across authenticated surfaces
- session-ready bootstrap read model clarification
- final authenticated application surface consistency hardening


---

## Phase 0.6.2 — Authenticated Surface Contract Alignment

Phase 0.6.2 consolidates the authenticated bootstrap surface by aligning the shared authenticated context exposed by the existing endpoints without introducing breaking contract changes.

While Phase 0.6.1 clarified boundary ownership between `/auth/me`, `/auth/session`, `/auth/me/settings`, and `/auth/wallets`, the construction of the authenticated context still depended on separate internal builders. That left a future drift risk between identity bootstrap data and token-derived session data.

This subphase removes that drift risk by aligning the shared authenticated context behind a common internal normalizer and by freezing cross-endpoint consistency with dedicated tests.

### Delivered in 0.6.2

Delivered in Phase 0.6.2:

- introduction of a shared authenticated context normalizer for profile-side and session-side authenticated projections
- explicit normalization of wallet-derived authenticated context so partial wallet claim state cannot leak inconsistent `wallet_id` or `chain` values when `wallet_address` is absent
- alignment of `/auth/me` and `/auth/session` around the same internal source of truth for shared authenticated context fields
- contract-level test coverage that freezes cross-endpoint consistency between `/auth/me`, `/auth/session`, and `/auth/wallets`
- preservation of existing public JSON envelopes and endpoint compatibility

### Internal alignment introduced in 0.6.2

The shared authenticated context now centrally derives and aligns:

- `user_id`
- `email`
- `auth_method`
- `wallet_id`
- `wallet_address`
- `chain`
- `has_wallet_session`

This shared derivation is now used by both the profile-side bootstrap projection and the session-side JWT/session projection.

### `/auth/me` and `/auth/session` alignment

Phase 0.6.2 keeps the responsibilities established in 0.6.1 unchanged, but aligns their shared authenticated fields so that both endpoints remain semantically separated while still projecting the same authenticated context when they originate from the same claims.

The alignment now ensures consistency for:

- authenticated user identity
- auth method
- wallet session context
- wallet address context
- wallet chain context
- attached user summary consistency

### `/auth/me` and `/auth/wallets` alignment

The authenticated bootstrap profile and the wallet inventory surface remain separate endpoints, but 0.6.2 hardens their relationship by validating that the primary wallet projected by `/auth/me` stays aligned with the primary wallet represented in `/auth/wallets`.

This protects the bootstrap surface from future drift in:

- primary wallet identity
- wallet address projection
- primary wallet status
- linked and detached timestamps
- primary designation semantics

### Test-level enforcement added in 0.6.2

The subphase is now enforced with contract-oriented authenticated handler tests that verify:

- `/auth/me` and `/auth/session` expose aligned shared authenticated context fields when built from the same claims
- `/auth/me` top-level `user` and profile-side `user` remain internally consistent
- `/auth/me` primary wallet projection stays aligned with `/auth/wallets` inventory output
- wallet-session normalization remains consistent when wallet-address context is absent

### Compatibility

Phase 0.6.2 remains fully backward compatible.

Specifically:

- no endpoint was renamed
- no response envelope was removed
- no public field was renamed
- no authenticated route was moved
- no settings contract was changed
- no wallet lifecycle operation was changed
- no pagination contract was changed

### Why this subphase matters

Phase 0.6.2 turns the authenticated bootstrap boundary defined in 0.6.1 into a contract-aligned authenticated surface.

This matters because frontend bootstrap consumers must be able to rely not only on endpoint responsibility boundaries, but also on stable cross-endpoint consistency for the authenticated context that those endpoints share.

Phase 0.6.2 therefore prepares the backend for the next step of Phase 0.6:

- session-ready bootstrap read model consolidation
- authenticated bootstrap consumption without ambiguity across surfaces
- final authenticated application surface hardening

## Phase 0.6.3 — Session-Ready Bootstrap Read Model

### Objective

Provide a single endpoint that aggregates the authenticated application surface into a frontend-ready bootstrap read model.

### Implementation

- Introduced `GET /auth/bootstrap`
- Aggregates:
  - session
  - user
  - profile
  - settings
  - wallet snapshot
- Reuses existing logic:
  - session builder
  - profile builder
  - settings service
  - wallet store

### Guarantees

- No breaking changes
- No endpoint replacement
- No contract modification of existing endpoints
- Pure aggregation layer

### Result

Authenticated surface is now:

- boundary-defined (0.6.1)
- contract-aligned (0.6.2)
- bootstrap-ready (0.6.3)

### Next

0.6.4 — Application Surface Consistency Hardening

## Phase 0.6.4 — Application Surface Consistency Hardening

### Objective

Harden the authenticated application surface by freezing a consistent structural contract across the bootstrap and wallet surfaces without breaking existing consumers.

### Implementation

- Added canonical `items` support to `/auth/wallets` while preserving legacy `wallets`
- Hardened `/auth/bootstrap` wallet block to expose only the canonical shape
- Added structural hardening tests for:
  - `/auth/wallets`
  - `/auth/bootstrap`

### Guarantees

- No breaking path changes
- No removal of legacy compatibility fields from `/auth/wallets`
- No business logic expansion
- No lifecycle changes
- Canonical wallet envelope consistency across authenticated surfaces

### Result

Authenticated surface is now:

- boundary-defined (0.6.1)
- contract-aligned (0.6.2)
- bootstrap-ready (0.6.3)
- consistency-hardened (0.6.4)

### Phase 0.6 Closure

Phase 0.6 is now complete.

It delivered:

- explicit authenticated surface boundaries
- aligned shared authenticated context
- session-ready bootstrap aggregation
- structural hardening of the authenticated application surface

### Next

Phase 1 — next foundation evolution beyond authenticated bootstrap consolidation

## Phase 0.7.1 — Application Layer Boundary Definition

### Objective

Introduce the first explicit application-layer boundary so HTTP handlers stop owning full orchestration of authenticated use cases.

### Implementation

- Introduced `internal/modules/auth/application.go`
- Added an explicit `Application` entry point for authenticated use cases
- Moved bootstrap orchestration into `Application.GetBootstrap(...)`
- Reduced `/auth/bootstrap` handler responsibilities to:
  - request context extraction
  - application invocation
  - HTTP error/status mapping
  - response writing

### Guarantees

- No public contract changes
- No route changes
- No business-domain expansion
- No broad refactor of all handlers
- Existing builders, services, and stores remain the underlying execution units

### Result

Stage 0 foundation now progresses beyond authenticated surface hardening into application-structure enablement.

Current application evolution state:

- 0.6.1 → authenticated surface boundaries
- 0.6.2 → authenticated contract alignment
- 0.6.3 → bootstrap read model
- 0.6.4 → structural hardening
- 0.7.1 → first explicit application-layer boundary

### Next

0.7.2 — Authenticated Surface Use Cases Extraction

## Phase 0.7.2 — Authenticated Surface Use Cases Extraction

### Objective

Extend the new application layer so the main authenticated surface no longer depends on handler-driven orchestration.

### Implementation

- Extended `internal/modules/auth/application.go`
- Added explicit authenticated use cases for:
  - `Login(...)`
  - `GetMe(...)`
  - `GetSession(...)`
- Preserved `GetBootstrap(...)` introduced in 0.7.1
- Reduced authenticated handlers to:
  - request parsing / context extraction
  - application invocation
  - HTTP status / error mapping
  - response writing

### Guarantees

- No public contract changes
- No route changes
- No wallet-management migration yet
- No global error model introduced
- Existing services, builders, and stores remain the execution layer

### Result

The authenticated surface is now application-driven for its core use cases:

- 0.7.1 → application-layer boundary introduced
- 0.7.2 → authenticated surface use cases extracted

This establishes the correct base for wallet-specific extraction in 0.7.3.

### Next

0.7.3 — Wallet Management Use Cases Consolidation

## Phase 0.7.3 — Wallet Management Use Cases Consolidation

### Objective

Complete the application-layer migration by consolidating all wallet management flows into explicit application use cases.

### Implementation

- Extended `internal/modules/auth/application.go` with wallet use cases:
  - `ListWallets(...)`
  - `CreateWalletLinkChallenge(...)`
  - `VerifyWalletLink(...)`
  - `CreateWalletAccountMergeChallenge(...)`
  - `VerifyWalletAccountMerge(...)`
  - `SetPrimaryWallet(...)`
  - `CheckWalletDetach(...)`
  - `ExecuteWalletDetach(...)`
- Refactored wallet HTTP handlers to:
  - delegate orchestration to application layer
  - remain transport-only
- Aligned wallet listing (`/auth/wallets`) with application layer

### Guarantees

- No public contract changes
- No route changes
- No payload changes
- No restructuring of modules
- Services remain execution layer

### Result

Wallet management is now fully application-driven:

- 0.7.1 → application boundary introduced
- 0.7.2 → authenticated surface extracted
- 0.7.3 → wallet management consolidated

The auth module now follows a consistent architecture across all endpoints.

### Next

0.7.4 — Handler Simplification & Contract Preservation

## Phase 0.7.4 — Handler Simplification & Contract Preservation

### Objective

Finalize the application-layer foundation by simplifying all authenticated handlers and preserving public contracts while removing residual transport duplication.

### Implementation

- Introduced shared HTTP transport helpers for:
  - request decoding
  - authenticated claims extraction
  - error JSON writing
- Refactored authenticated handlers to follow a unified flow:
  - parse / extract context
  - invoke application layer
  - map error to HTTP
  - write response
- Simplified:
  - login handlers
  - me / session / profile update handlers
  - bootstrap handler
  - wallet handlers
  - wallet list handler

### Guarantees

- No public contract changes
- No route changes
- No payload changes
- No application-layer regression
- No service/store refactor

### Result

Phase 0.7 is now functionally complete:

- 0.7.1 → application boundary introduced
- 0.7.2 → authenticated surface extracted
- 0.7.3 → wallet management consolidated
- 0.7.4 → handlers simplified and contracts preserved

The auth module now exposes a fully application-driven execution model with transport-only handlers.

### Next

0.8 — Standardized Error Model

## Phase 0.8.1 — Error Contract Definition

### Objective

Introduce the first real implementation step of the standardized error model by defining a shared error envelope and centralizing HTTP error serialization.

### Implementation

- Added shared response error contract under `internal/core/errs`
- Added centralized error envelope serialization under `internal/core/httpx`
- Migrated auth HTTP handlers to the new envelope while preserving current control flow
- Migrated auth middleware and recoverer/timeout transport responses away from legacy string-only error payloads
- Updated auth HTTP tests to assert structured error payloads

### Guarantees

- No route changes
- No success payload changes
- No domain/application refactor yet
- Error details now travel under `error.details`
- Auth and middleware error responses now share one structured envelope

## Phase 0.8.2 — Error Type System Introduction

### Objective

Introduce a reusable internal application error type system so the backend no longer depends on auth-local string catalogs to express normalized error codes, messages, status and categories.

### Implementation

- Added `AppError`, `Category`, wrapping helpers and conversion helpers under `internal/core/errs`
- Moved the auth/wallet/settings legacy-to-normalized error catalog out of `internal/modules/auth/http_login.go` into `internal/core/errs`
- Added reusable typed factories for common auth/internal/settings error cases
- Added HTTP writing support for `AppError` in `internal/core/httpx`
- Rebased auth transport helpers and core auth middleware on the new internal error type system while preserving the 0.8.1 envelope

### Guarantees

- No route changes
- No success payload changes
- No auth surface-wide semantic remapping yet
- Existing 0.8.1 envelope remains stable
- Internal error normalization is now centralized and reusable across subsequent subphases

## Phase 0.8.3 — Auth Surface Error Standardization

### Objective

Apply the 0.8.2 internal app error type system across the auth transport surface so handlers consume centralized typed errors instead of auth-local legacy string mappings.

### Implementation

- Migrated `internal/modules/auth/http_login.go`, `http_bootstrap.go`, `http_wallet.go` and `http_wallet_list.go` to write structured transport errors from centralized `core/errs` app errors
- Preserved auth-local JSON writing to avoid reintroducing cyclic imports with `internal/core/httpx/router.go`
- Reduced the auth-local legacy adapter to compatibility/testing support instead of runtime ownership
- Standardized wallet/auth/settings handler error emission on shared factories and centralized catalog resolution

### Guarantees

- No route changes
- No success payload changes
- No business-flow redesign
- No reintroduction of `httpx -> auth -> httpx` cycles
- Auth surface now depends on centralized app-error factories instead of local string contracts

### Result

Phase 0.8 is now fully completed across its four implementation layers:

- 0.8.1 → error contract introduced
- 0.8.2 → internal app error type system introduced
- 0.8.3 → auth surface standardized on app errors
- 0.8.4 → error mapping hardened and contract tests added

### Validation

- Added contract-oriented test coverage in `internal/core/errs/app_error_test.go`
- Added HTTP-envelope hardening coverage in `internal/core/httpx/error_test.go`
- Preserved green `go test ./...` validation after hardening

### Guarantees After Phase 0.8

- Canonical `error` envelope is now frozen and reused across middleware and auth handlers
- Centralized `AppError`/catalog/factory flow is now hardened by tests
- Representative code/status/category mappings are now explicitly validated
- No success payload contract changed during the phase
- No route contract changed during the phase
- No `httpx -> auth -> httpx` cyclic import was reintroduced

### Next

Phase 0.8 is closed. The next roadmap step should start from the post-0.8 foundation state rather than from a pending error-model subphase.

## Phase 0.9.1 — Versioning Policy Definition

### Objective

Define the official API versioning policy for the post-0.8 backend so public contract evolution becomes explicit before authorization and later domain growth are introduced.

### Why This Phase Follows 0.8

Phase 0.8 closed the standardized error model and froze the authenticated transport error envelope, but the backend still exposes its public surface through unversioned routes. Once a stable transport contract exists, the next architectural need is to define how that contract evolves without breaking current consumers.

### Policy Defined

The backend now treats path-based versioning as the canonical public strategy:

- canonical route namespace: `/api/v1/...`
- current unversioned endpoints: legacy, backward-compatible, non-canonical
- `v1` semantics: current success payloads plus the Phase 0.8 standardized error envelope

### Guarantees

- No business-flow redesign
- No immediate route removal
- No success payload mutation
- No breaking error-envelope change inside `v1`
- No forced frontend migration while the frontend remains aligned to backend Phase 0.6 until Stage 0 is complete

### Immediate Architectural Consequence

Versioning is now explicitly a transport-layer concern. The application layer introduced in Phase 0.7 remains the orchestration boundary, and future authorization work in 0.10 must build on top of this versioning policy rather than define transport evolution ad hoc.

### Next

Phase 0.9.2 should introduce the canonical router exposure under `/api/v1/...` while preserving the current legacy route surface.

## Phase 0.9.2 — Router Versioning Foundation

### Objective

Materialize the canonical `v1` transport surface in the real HTTP router without changing the current business flows, payload semantics or legacy route availability.

### Why This Subphase Follows 0.9.1

Phase 0.9.1 defined the policy, but the router still exposed only the legacy unversioned route space. The next natural step is to project the already-stabilized auth and authenticated surfaces into the canonical `/api/v1/...` namespace while preserving compatibility for existing consumers.

### Implementation

- Reworked `internal/core/httpx/router.go` so auth and authenticated route registration now happens through one shared route-registration helper
- Kept the current legacy transport surface unchanged
- Added the canonical `/api/v1/...` transport surface on top of the same handlers
- Preserved current middleware behavior, including auth protection, access logging, timeout wrapping and structured error handling
- Avoided route-level business duplication by making versioning a pure transport registration concern

### Routes Now Exposed Canonically

The backend now exposes canonical `v1` routes for the currently stabilized auth and auth-adjacent surface, including:

- `/api/v1/auth/login`
- `/api/v1/auth/wallet/challenge`
- `/api/v1/auth/wallet/verify`
- `/api/v1/auth/wallets/link/challenge`
- `/api/v1/auth/wallets/link/verify`
- `/api/v1/auth/account/merge/wallet/challenge`
- `/api/v1/auth/account/merge/wallet/verify`
- `/api/v1/auth/bootstrap`
- `/api/v1/auth/me`
- `/api/v1/auth/me/settings`
- `/api/v1/auth/session`
- `/api/v1/auth/wallets`
- `/api/v1/auth/wallets/detach/check`
- `/api/v1/auth/wallets/detach`
- `/api/v1/auth/wallets/primary`

### Guarantees

- No legacy route was removed
- No success payload changed
- No error-envelope contract changed
- No application/service orchestration changed
- No frontend migration is required while the frontend remains aligned to backend Phase 0.6 during Stage 0

### Result

The backend now has two public transport entry surfaces for the same stabilized behavior:

- legacy compatibility routes under `/auth/...`
- canonical `v1` routes under `/api/v1/auth/...`

This completes the router-level materialization of the versioning model defined in 0.9.1 and prepares the contract-freezing work of 0.9.3.

### Next

Phase 0.9.4 formalizes version-aware contract tests so legacy and canonical `v1` entry paths remain protected from silent divergence. The next step is 0.9.5 — Documentation Consolidation.

## Phase 0.9.3 — Authenticated Surface Version Freezing

### Objective

Bind the currently exposed authenticated surface explicitly to canonical `v1` semantics without changing business behavior, payload structure or middleware rules.

### Why This Subphase Follows 0.9.2

Phase 0.9.2 made the canonical route space real in the runtime router, but route exposure alone does not define which transport surface is actually frozen as `v1`. The next step is to declare that the existing authenticated behavior reached through those routes is now the canonical `v1` authenticated contract rather than a provisional projection.

### Canonical `v1` Authenticated Surface

The backend now treats the following authenticated routes as the frozen `v1` authenticated surface, exposed both through legacy compatibility paths and canonical versioned paths:

- `GET /auth/bootstrap` ↔ `GET /api/v1/auth/bootstrap`
- `GET /auth/me` ↔ `GET /api/v1/auth/me`
- `PATCH /auth/me` ↔ `PATCH /api/v1/auth/me`
- `GET /auth/me/settings` ↔ `GET /api/v1/auth/me/settings`
- `PATCH /auth/me/settings` ↔ `PATCH /api/v1/auth/me/settings`
- `GET /auth/session` ↔ `GET /api/v1/auth/session`
- `GET /auth/wallets` ↔ `GET /api/v1/auth/wallets`
- authenticated wallet-management flows already routed through `/auth/...` and `/api/v1/auth/...` using the same handlers and middleware

### Freeze Semantics

For the authenticated `v1` surface, the backend now makes these guarantees explicit:

- legacy and canonical entry paths are two transport projections of the same authenticated contract
- business behavior remains identical across both route spaces
- success payload semantics remain unchanged
- the standardized Phase 0.8 error envelope remains part of the same `v1` contract
- future breaking changes to this authenticated surface require a new API version instead of silent mutation of the current one

### Non-Goals

This subphase does not yet add route-by-route equivalence assertions in automated coverage. That stricter transport regression protection remains the responsibility of 0.9.4.

### Result

After 0.9.3, the backend no longer only exposes a canonical route space; it also explicitly freezes which authenticated surface that route space represents for `v1`.

### Next

Phase 0.9.4 should formalize version-aware contract tests so legacy and canonical `v1` entry paths remain protected from silent divergence.

## Phase 0.9.4 — Version-aware Contract Testing

### Objective

Turn the `v1` freeze defined in 0.9.3 into executable regression protection by asserting that legacy `/auth/...` and canonical `/api/v1/auth/...` entry paths remain behaviorally aligned where the current Stage 0 surface is expected to match.

### Why This Subphase Follows 0.9.3

Phase 0.9.3 defined which authenticated surface is frozen as canonical `v1`, but a declared freeze is still vulnerable to silent drift unless the transport layer is covered by route-aware regression tests. The next step is therefore to express that equivalence in automated coverage rather than only in architectural and phase documentation.

### Delivered

Phase 0.9.4 adds explicit version-aware transport tests around the router foundation introduced in 0.9.2:

- verifies helper path composition for legacy and canonical route registration
- verifies that protected legacy and canonical authenticated endpoints return the same standardized missing-bearer error contract
- verifies that protected legacy and canonical authenticated endpoints return the same standardized unauthorized error contract for invalid tokens
- verifies that representative public wallet-challenge endpoints exposed through `/auth/...` and `/api/v1/auth/...` preserve a compatible success payload shape

### Scope Boundary

This subphase focuses on transport-level contract equivalence for the currently frozen `v1` authenticated surface. It does not introduce `v2`, redefine payload semantics, or duplicate business behavior into route-specific handlers.

### Result

After 0.9.4, the versioning model is not only declared and routed; it is also protected by concrete regression tests that guard the coexistence of legacy and canonical `v1` paths.

### Next

Phase 0.9.5 should consolidate the complete trunk documentation set around the finalized 0.9 versioning model and testing guarantees.

## Phase 0.9.5 — Documentation Consolidation

### Objective

Close Phase 0.9 by aligning the full trunk documentation set with the now-real versioning model, the authenticated-surface `v1` freeze and the representative legacy-versus-canonical contract testing introduced through 0.9.4.

### Why This Subphase Follows 0.9.4

By the end of 0.9.4, the backend already had all three technical pieces of the versioning foundation in place:

- a documented versioning policy (0.9.1)
- a real canonical `/api/v1/...` router surface (0.9.2)
- an explicitly frozen authenticated `v1` contract plus representative transport regression tests (0.9.3 and 0.9.4)

What still remained was the trunk-documentation responsibility: removing the state drift that had accumulated while the subphases were landing and leaving the repository with a single coherent narrative before Phase 0.10 begins.

### Delivered

Phase 0.9.5 consolidates the documentation set around the finalized versioning model:

- marks Phase 0.9 as completed across roadmap, phase-status, index and top-level project status surfaces
- records 0.9.4 as completed and 0.9.5 as the documentation-consolidation closure step rather than leaving either one pending
- aligns the phase-specific, architectural, testing and handoff documents around the same canonical explanation of legacy `/auth/...` versus canonical `/api/v1/auth/...`
- preserves the project rule that the frontend remains aligned to backend Phase 0.6 until Stage 0 is fully closed, so the documentation does not imply forced frontend adoption of the new route space
- leaves Phase 0.10 clearly framed as the next planned step rather than mixing authorization concerns into the versioning closure

### Scope Boundary

This subphase does not change runtime code, router behavior, payload semantics, middleware logic or test intent. It closes Phase 0.9 by eliminating documentation drift and by making the already-delivered versioning foundation readable as one coherent trunk narrative.

### Result

After 0.9.5, Phase 0.9 is fully closed:

- `v1` versioning policy is defined
- canonical routing is implemented
- the authenticated `v1` surface is frozen
- representative version-aware transport regression tests are in place
- the trunk documentation set now describes the same state consistently

### Next

Phase 0.10 should build the authorization layer on top of this now-completed versioning and contract-governance foundation.

## Phase 0.10.1 — Authorization Model Definition

### Objective

Introduce the first real authorization artifact set without changing request behavior yet, so the backend moves from an authentication-only foundation to an explicit authorization model foundation.

### Why This Subphase Follows 0.9.5

Phase 0.9 deliberately closed the transport/versioning boundary before any permission semantics were introduced. With canonical `v1` routing, standardized errors and the authenticated application surface already stabilized, the next safe move is not immediate endpoint blocking but a static authorization model that defines the language of later enforcement.

The backend already knows who the authenticated subject is through JWT claims and auth middleware, but it still lacks a first-class representation of:

- roles
- permissions
- role-to-permission mapping
- an authorization subject model distinct from transport claims

Adding those primitives first keeps the next subphases honest: middleware/context propagation in 0.10.2 and policy evaluation in 0.10.3 can build on explicit structures instead of inventing ad-hoc handler-level checks.

### Delivered

Phase 0.10.1 introduces a new core authorization package under `internal/core/authorization` and keeps the scope intentionally static:

- defines the initial role model with `user` and `admin`
- defines the initial permission vocabulary for current Stage 0 authenticated surfaces
- defines the foundational role → permission mapping
- introduces an `AuthorizationSubject` structure that is separate from raw JWT claims
- adds normalization and immutability-oriented helpers so later subphases can consume a stable model
- adds focused unit tests for normalization, mapping immutability and permission aggregation semantics

This subphase intentionally does **not** yet:

- attach authorization data to HTTP requests
- evaluate policies at runtime
- deny endpoint access
- persist roles or permissions in the database
- mutate existing authenticated response contracts

### Current Result

After 0.10.1, the backend still behaves exactly as before at runtime, but it no longer lacks an explicit authorization vocabulary. The system now has a dedicated core model describing the authorization subject and the static permission space that later subphases will propagate and enforce.

Architecturally, the backend has now crossed from:

`authentication only`

into:

`authentication + explicit authorization model foundation`

without prematurely coupling handler logic, transport behavior or persistence to enforcement concerns.

#### Next

0.10.2 — Authorization Context & Middleware




## Phase 0.10.2 — Authorization Context & Middleware

### Objective

Propagate the authorization subject through the authenticated HTTP request lifecycle without changing endpoint success payloads, error contracts or enforcement behavior.

### Delivered

Phase 0.10.2 builds on the static authorization model introduced in 0.10.1 and wires it into the transport layer in a non-breaking way:

- adds authorization-context helpers in `internal/core/authorization` for storing and retrieving an `AuthorizationSubject` from `context.Context`
- adds a deterministic subject resolver from authenticated JWT claims to the authorization model
- introduces a dedicated HTTP middleware in `internal/core/httpx` that hydrates authorization data only after authentication has already succeeded
- integrates that middleware into the authenticated route pipeline for both legacy `/auth/...` and canonical `/api/v1/auth/...` surfaces
- adds focused tests covering authorization-context storage, claim-to-subject resolution and middleware hydration behavior

### Runtime Result

After 0.10.2, authenticated requests still behave exactly as before from a client perspective, but the backend request lifecycle now carries both:

- authentication claims
- a normalized authorization subject

This keeps the architecture aligned with the intended layering:

`HTTP → Auth → Authorization → Application → Domain`

while still deferring policy checks and endpoint denial to later subphases.

### Next

0.10.3 — Policy Evaluation Layer

## Phase 0.10.3 — Policy Evaluation Layer

### Objective

Introduce a centralized authorization decision boundary so the backend can answer permission questions through a stable core API before any endpoint begins denying requests.

### Why this subphase exists

By the end of 0.10.2, authenticated requests already carry a normalized `AuthorizationSubject`. What still remained missing was the policy boundary itself: a reusable place where the backend can decide whether a subject may perform an action on a resource without embedding that logic directly in handlers or tying it to JWT claims.

### What 0.10.3 introduces

Phase 0.10.3 adds that missing central layer under `internal/core/authorization`:

- explicit authorization `Action` vocabulary
- explicit authorization `Resource` vocabulary
- static action/resource → permission projection
- centralized permission resolution helpers
- `PolicyEvaluator` with `Evaluate(...)` and `Can(...)` APIs
- focused tests covering permission projection and role-based authorization decisions

### Runtime effect

Phase 0.10.3 remains non-breaking and non-enforcing:

- no endpoint behavior changes yet
- no new transport contract is introduced
- no handler starts denying requests yet
- the evaluator exists so 0.10.4 can enforce permissions progressively without inventing policy rules locally

### Result

After 0.10.3, the authorization layer now consists of three explicit steps:

- 0.10.1 → model vocabulary (`Role`, `Permission`, `AuthorizationSubject`)
- 0.10.2 → request-context hydration (`HydrateAuthorization`)
- 0.10.3 → centralized policy evaluation (`PolicyEvaluator`)

The backend still behaves exactly as before from the client perspective, but it no longer lacks a central place to answer authorization questions. Endpoint-level enforcement can now begin in 0.10.4 on top of a real policy layer instead of ad-hoc handler checks.

### Next

0.10.4 — Endpoint-Level Enforcement


## Phase 0.10.4 — Endpoint-Level Enforcement

Introduce progressive endpoint-level authorization enforcement on the authenticated Stage 0 surfaces already represented by the static authorization model and centralized policy evaluator.

### Why this subphase exists

By the end of 0.10.3, the backend could already answer authorization questions through a stable internal policy API. What was still missing was operational enforcement: selected endpoints still authenticated users, but did not yet deny unauthorized access through the new authorization layer.

### What 0.10.4 introduces

Phase 0.10.4 makes authorization operational on a first controlled slice by introducing:

- route-level `RequirePermission(...)` middleware in `internal/core/httpx`
- standardized `AUTH_FORBIDDEN` responses for denied authorization checks
- progressive enforcement on selected authenticated `me` and `settings` endpoints
- focused tests confirming allowed and denied outcomes under the centralized policy layer

### Runtime effect

Phase 0.10.4 is intentionally progressive rather than universal:

- `GET /auth/me` and `GET /api/v1/auth/me` now enforce `read user`
- `GET /auth/me/settings` and `GET /api/v1/auth/me/settings` now enforce `read settings`
- `PATCH /auth/me/settings` and `PATCH /api/v1/auth/me/settings` now enforce `update settings`
- `PATCH /auth/me` remains outside the first enforcement slice because its self-update semantics are not yet fully represented by the current static role model

### Result

After 0.10.4, the authorization layer no longer exists only as internal preparation. The backend now actively denies unauthorized access on a selected subset of authenticated endpoints while preserving the stabilized transport contract for allowed requests.

### Next

0.10.5 — Documentation & Contract Consolidation

## Phase 0.10.5 — Documentation & Contract Consolidation

Close Phase 0.10 by consolidating the repository narrative and contract documentation around the authorization layer already delivered in 0.10.1 through 0.10.4.

### Why this subphase exists

Once 0.10.4 introduced real endpoint-level enforcement, the repository could no longer leave trunk documents in mixed intermediate states. The last step of the phase is therefore documentary and contractual: bring README, roadmap, handoff, architecture, status and the dedicated phase document into one coherent explanation of the delivered authorization layer.

### What 0.10.5 introduces

Phase 0.10.5 adds no new runtime behavior. Instead, it consolidates the project narrative by:

- marking Phase 0.10 as completed across the trunk documentation set
- removing stale text that still described authorization as entirely non-enforcing after 0.10.4
- documenting the progressive scope of the first enforcement slice
- aligning the public contract narrative around the standardized `AUTH_FORBIDDEN` path for covered endpoints

### Result

After 0.10.5, the authorization layer is both implemented and consistently documented:

- 0.10.1 → authorization model definition
- 0.10.2 → authorization context and middleware hydration
- 0.10.3 → centralized policy evaluation
- 0.10.4 → progressive endpoint-level enforcement
- 0.10.5 → documentation and contract consolidation

Phase 0.10 is complete. The next phase can build on a closed and consistently described authorization foundation instead of revisiting stale intermediate state.

### Next

0.11 — Domain Module Pattern


## Phase 0.11 — Domain Module Pattern

Phase 0.11 was the structural Stage 0 step that standardized how domain-oriented modules are organized internally without changing the public transport contract, the authenticated bootstrap semantics or the already delivered authorization behavior.

### Why this phase exists

After Phases 0.6 through 0.10, the backend already has a stabilized authenticated surface, an explicit application-layer foundation, a standardized error model, a canonical `v1` route space and a progressively enforced authorization boundary. What it still lacked before 0.11 was a uniform internal module pattern across the current Stage 0 domain-facing modules.

Without that pattern, the codebase risks keeping:

- handlers that accumulate orchestration responsibilities
- implicit domain boundaries
- module-local structures that are inconsistent from one module to another
- cross-module dependencies expressed through concrete knowledge rather than explicit contracts

Phase 0.11 addresses that structural gap without reworking business behavior.

### What 0.11 defines

Phase 0.11 introduces the repository-level definition of a **Domain Module Pattern** for the current Stage 0 modules, centered on the following internal shape:

```text
internal/modules/<module>/
    http/
        handlers.go
        dto.go

    app/
        service.go
        usecases.go

    domain/
        model.go
        contracts.go

    repository/
        repository.go   (when it adds real clarity)
```

The architectural intention is explicit separation of responsibilities:

- **HTTP** → transport parsing, request validation, response mapping
- **APP** → use-case orchestration and flow coordination
- **DOMAIN** → models, invariants and internal contracts
- **REPOSITORY** → persistence-facing implementation boundaries when needed

### Phase 0.11 subphases

- **0.11.1 — Domain Module Pattern Definition** ✔ Completed
- **0.11.2 — User Module Refactor** ✔ Completed
- **0.11.3 — UserSettings Module Refactor** ✔ Completed
- **0.11.4 — Auth Module Alignment** ✔ Completed
- **0.11.5 — Cross-Module Contract Consolidation** ✔ Completed
  - **0.11.5.0 — Subphase Definition & Documentation Alignment** ✔ Completed
  - **0.11.5.1 — Dependency Mapping** ✔ Completed
  - **0.11.5.2 — Contract Extraction** ✔ Completed
  - **0.11.5.3 — Interface Alignment** ✔ Completed
  - **0.11.5.4 — Runtime Compatibility Validation** ✔ Completed
- **0.11.6 — Documentation & Phase Closure** ✔ Completed

### Intended scope

Phase 0.11 is intentionally structural and non-functional:

- reorganize internal module boundaries
- make use-case orchestration explicit
- formalize domain-facing contracts
- reduce accidental coupling across `auth`, `user` and `usersettings`
- preserve complete compatibility at the HTTP/API level

This phase does **not** introduce:

- new endpoints
- new business features
- public payload changes
- auth/authorization redesign
- CQRS/event sourcing or multi-tenant architecture

### Current result

0.11.1 established the architectural definition, scope, dependency direction and ownership model of the Domain Module Pattern across the repository documentation.

0.11.2 applied that pattern to `internal/modules/user`, making `user` the first concrete reference implementation of the new structure while preserving the current runtime contract. The refactor is validated in repository state by passing module and repository tests.

0.11.3 applied the same pattern to `internal/modules/usersettings`, preserving settings-specific semantics while exposing explicit `app`, `domain` and `repository` boundaries.

0.11.4 completed the conservative auth-module alignment. The auth module now has `auth/app` as the runtime/application owner, `auth/domain` as the canonical contract boundary, root `auth` as a compatibility surface, and `auth/repository` as a documented transitional façade. Public wallet bootstrap, authenticated wallet management, session, bootstrap and login behavior remain externally unchanged.

0.11.5 completed cross-module contract consolidation. `auth/app` now depends on explicit minimal contracts (`UserProvider` and `UserSettingsProvider`) instead of concrete `user` and `usersettings` services, while root `auth` preserves compatibility wiring and public behavior.

0.11.6 closes the phase documentation. The trunk documentation now records Phase 0.11 as complete and keeps roadmap, status, architecture and handoff documents aligned with the delivered module pattern.

### Expected outcome

With Phase 0.11 implemented and closed, the backend now has:

- a uniform internal module pattern for the currently relevant Stage 0 modules
- thinner handlers and clearer use-case boundaries
- explicit ownership between `auth`, `user` and `usersettings`
- cross-module dependencies expressed through minimal contracts instead of incidental imports
- the same externally observable behavior as the completed 0.10 surface

### Status

Phase 0.11 is **completed**. The repository now has the Domain Module Pattern defined and applied across `user`, `usersettings` and `auth`, with cross-module coordination expressed through explicit minimal contracts where coordination is required.

### Next

0.12 — Read/Write Model Separation

---

## Phase 0.12 — Read / Write Model Separation

Phase 0.12 was the structural Stage 0 step after the completed Domain Module Pattern. It formalized an explicit separation between models used to read data from the system and models used to write or mutate data into the system.

The phase is intentionally internal and compatibility-preserving. It does not introduce a public API version change, does not change business behavior and does not adopt full CQRS or event sourcing.

### Objective

Separate the current model surface into explicit responsibilities:

- **Read models** for outputs, views and response-oriented structures
- **Write models** for inputs, commands and mutation-oriented structures
- **Domain models** for module-owned semantics and invariants
- **Mapping functions** for explicit transformations between those responsibilities

### Why this phase follows 0.11

Phase 0.11 made module ownership explicit across `auth`, `user` and `usersettings`. Once module boundaries are explicit, the next ambiguity is inside those modules: several structures currently act as request DTOs, response DTOs, application result shapes or domain structures depending on where they are used.

Phase 0.12 addresses that ambiguity without reopening the public contract.

### Included scope

- identify all current model structures in the real repository
- classify them as read, write, domain, infrastructure, contract or hybrid/transitional
- extract explicit read models where response-oriented structures are mixed with other concerns
- isolate write models for mutation/input paths
- introduce explicit mapping functions where transformations are currently implicit
- align existing internal contracts after 0.11 without changing runtime compatibility
- update trunk documentation cumulatively

### Excluded scope

- no full CQRS implementation
- no event sourcing
- no API version change
- no business behavior change
- no multi-tenant model redesign
- no public payload migration during 0.12.0

### Phase 0.12 subphases

- **0.12.0 — Phase Definition & Documentation Lock** ✔ Completed
- **0.12.1 — Model Classification Audit** ✔ Completed
  - **0.12.1.0 — Definition & Documentation Lock** ✔ Completed
  - **0.12.1.1 — Model Inventory Extraction** ✔ Completed
  - **0.12.1.2 — Model Classification** ✔ Completed
  - **0.12.1.3 — Cross-Layer Usage Analysis** ✔ Completed
  - **0.12.1.4 — Problem Detection & Risk Mapping** ✔ Completed
  - **0.12.1.5 — Target Separation Definition** ✔ Completed
  - **0.12.1.6 — Audit Consolidation & Closure** ✔ Completed
- **0.12.2 — Read Model Extraction** ✔ Completed
  - **0.12.2.0 — Definition & Documentation Lock** ✔ Completed
  - **0.12.2.1 — Read Model Design** ✔ Completed
  - **0.12.2.2 — Read Model Implementation** ✔ Completed
  - **0.12.2.3 — Domain/Application → Read Mapping** ✔ Completed
  - **0.12.2.4 — Response Alignment** ✔ Completed
  - **0.12.2.5 — Validation & Compatibility** ✔ Completed
  - **0.12.2.6 — Documentation & Closure** ✔ Completed
- **0.12.3 — Write Model Isolation** ✔ Completed
  - **0.12.3.0 — Definition & Documentation Lock** ✔ Completed
  - **0.12.3.1 — Write Model Design** ✔ Completed
  - **0.12.3.2 — Write Model Implementation** ✔ Completed
  - **0.12.3.3 — Write → Domain Mapping** ✔ Completed
  - **0.12.3.4 — Handler Alignment** ✔ Completed
  - **0.12.3.5 — Validation & Compatibility** ✔ Completed
  - **0.12.3.6 — Documentation & Closure** ✔ Completed
- **0.12.4 — Mapping Layer Introduction** ✔ Completed
  - **0.12.4.0 — Definition & Documentation Lock** ✔ Completed
  - **0.12.4.1 — Mapping Layer Design** ✔ Completed
  - **0.12.4.2 — Mapping Layer Implementation** ✔ Completed
  - **0.12.4.3 — Mapping Consolidation** ✔ Completed
  - **0.12.4.4 — Application Refactor** ✔ Completed
  - **0.12.4.5 — Validation & Compatibility** ✔ Completed
  - **0.12.4.6 — Documentation & Closure** ✔ Completed
- **0.12.5 — Contract Alignment** ✔ Completed
  - **0.12.5.0 — Definition & Documentation Lock** ✔ Completed
  - **0.12.5.1 — Contract Inventory & Classification** ✔ Completed
  - **0.12.5.2 — Contract Normalization Design** ✔ Completed
  - **0.12.5.3 — Contract Alignment Implementation** ✔ Completed
  - **0.12.5.4 — Handler Contract Adjustment** ✔ Completed
  - **0.12.5.5 — Validation & Compatibility** ✔ Completed
  - **0.12.5.6 — Documentation & Closure** ✔ Completed
### Current result

0.12.0 establishes the documentary lock for the phase. The repository documentation now defines the scope, non-goals, subphase order and compatibility requirements for Read / Write Model Separation before any code is changed.

0.12.1.0 extends that documentary discipline to the Model Classification Audit itself. Because 0.12.1 is subdivided into smaller audit steps, the `.0` step records the audit method, classification criteria and sub-subphase order before inventory extraction begins.

0.12.1.1 through 0.12.1.6 complete the audit sequence. The repository now contains a full model inventory, model classification, cross-layer usage analysis, problem/risk mapping and target separation definition. The audit reviewed 125 model-like structs and identified 11 hybrid/transitional structures that must guide the extraction work in 0.12.2 and later.

### 0.12.2.0 result

0.12.2.0 locks the Read Model Extraction plan before any Go code changes. It records the internal sub-subphase sequence, the compatibility constraints for response-preserving extraction and the rule that 0.12.2 must use the 0.12.1 audit artifacts as evidence.

0.12.2 is completed. The repository now contains explicit read model packages for `auth`, `user` and `usersettings`, explicit Domain/Application → Read mapping functions, response alignment that preserves public JSON contracts and a validation record backed by `go test ./...`.

### 0.12.3.0 result

0.12.3.0 locks the Write Model Isolation plan before any Go code changes. It records the internal sub-subphase sequence, the compatibility constraints for input-preserving isolation and the rule that write models must not reuse read models introduced in 0.12.2.

0.12.3 is completed. The repository now contains explicit write model packages for `auth`, `user` and `usersettings`, write-side domain input structures, explicit Write → Domain mapping functions, handler alignment preserving public request payload semantics and a validation record backed by `go test ./...`.

### Next

0.12.4 — Mapping Layer Introduction

### 0.12.4.0 result

0.12.4.0 is completed as a documentation-only subphase. It locks the Mapping Layer Introduction plan before any Go code changes, defines the centralized mapper package direction and records that existing read/write mapper behavior must remain compatible until consolidation is implemented.


## Phase 0.12.4 Closure

Phase 0.12.4 — Mapping Layer Introduction is completed. The backend now contains module-local mapper packages, consolidated read/write transformation ownership, reduced application-layer mapping residuals and validation documentation backed by `go test ./...`. The next planned phase is 0.12.5 — Contract Alignment.

## Phase 0.12.5.0 — Contract Alignment Definition & Documentation Lock

Phase 0.12.5.0 starts the Contract Alignment sequence after the Mapping Layer Introduction closure. This documentation-only lock defines the contract alignment plan before any provider or runtime contract is changed.

The phase focuses on reviewing and aligning internal provider contracts, especially the contracts introduced and stabilized during Phase 0.11, with the read/write model separation and centralized mapping ownership introduced during Phase 0.12.2 through Phase 0.12.4.

No Go code is changed in 0.12.5.0.

## Phase 0.12.5 Closure

Phase 0.12.5 — Contract Alignment is completed. The backend now contains contract inventory documentation, contract normalization design, contract alignment implementation, centralized HTTP contract aliases, handler contract adjustment and validation documentation backed by `go test ./...`.

The completed subphase sequence is:

- **0.12.5.0 — Definition & Documentation Lock** ✔ Completed
- **0.12.5.1 — Contract Inventory & Classification** ✔ Completed
- **0.12.5.2 — Contract Normalization Design** ✔ Completed
- **0.12.5.3 — Contract Alignment Implementation** ✔ Completed
- **0.12.5.4 — Handler Contract Adjustment** ✔ Completed
- **0.12.5.5 — Validation & Compatibility** ✔ Completed
- **0.12.5.6 — Documentation & Closure** ✔ Completed

Phase 0.12 is complete. Next work must start from the next roadmap-defined phase definition.

## Phase 0.13 — Provider Layer Consolidation

Phase 0.13 starts after the completed Read / Write Model Separation phase.

The phase consolidates the Provider Layer as the explicit entry point to domain services, using the model separation, mapper ownership and contract alignment completed in Phase 0.12 as its baseline.

Phase 0.12 left the system with clearer read/write model ownership, centralized mapper responsibilities and safer internal contracts. Phase 0.13 continues that story by addressing the remaining runtime composition problem: handlers and router construction could still see too many lower-level services and stores even though the public behavior was already stable. The provider layer is therefore introduced as a structural consolidation step, not as a new feature surface.

The narrative transition is:

```text
Phase 0.12: clarify model direction and mapping ownership
Phase 0.13: clarify orchestration ownership and runtime composition
```

### Objective

Establish a consistent provider boundary so handlers and application flows do not rely on scattered direct access to domain or repository responsibilities.

### Included scope

- provider inventory and classification
- provider interface design
- provider implementation where required
- application and handler integration
- compatibility validation
- accumulated documentation closure

### Excluded scope

- public API changes
- business behavior changes
- API versioning changes
- CQRS or event sourcing
- observability implementation

### Phase 0.13 subphases

- **0.13.0 — Definition & Documentation Lock** ✔ Completed
- **0.13.1 — Provider Inventory & Classification** ✔ Completed
- **0.13.2 — Provider Interface Design** ✔ Completed
- **0.13.3 — Provider Implementation** ✔ Completed
- **0.13.4 — Application Integration** ✔ Completed
- **0.13.5 — Validation & Compatibility** ✔ Completed
- **0.13.6 — Documentation & Closure** ✔ Completed

### Current status

Phase 0.13 is completed. Provider Layer Consolidation is now the documented runtime boundary between HTTP handlers and module application/domain responsibilities.

The concrete outcome is that the backend can now explain how a request crosses the authenticated surface without treating handlers, application services and stores as one mixed composition area. The public API remains unchanged, but the internal ownership story is stronger: transport stays at the edge, orchestration moves behind provider contracts, use cases remain in application services and domain/repository responsibilities stay below that boundary.

### Phase 0.13.3 result

Provider implementation introduced explicit auth provider interfaces and composition-root wiring for the consolidated auth provider. Auth HTTP handlers now depend on the provider boundary for session, profile, settings and wallet orchestration while preserving public API contracts, route registration, authorization middleware and standardized error responses.


### Phase 0.13.4 result

Application integration aligned runtime HTTP wiring with the Provider Layer boundary. Router construction now creates auth HTTP handlers from the consolidated provider only, removing direct runtime handler wiring of user services, user settings services, wallet challenge stores, wallet identity stores, challenge TTL and public base URL. Compatibility fallback fields remain available only for transitional tests and do not define production wiring.

### Phase 0.13.5 result

Phase 0.13.5 completed the validation and compatibility checkpoint for Provider Layer Consolidation. The validation baseline confirms that the 0.13.4 provider integration path builds and passes `go test ./...` in the project environment after the router test compatibility fix. This subphase does not introduce production code changes.


### Phase 0.13.6 result

Phase 0.13.6 closed Provider Layer Consolidation documentation. The closure removed generic repeated phase-status blocks, corrected misplaced roadmap content, expanded the final subphase state consistently and evolved trunk documentation with the actual provider-layer impact instead of only marking state transitions. No Go code or public API behavior was changed by this closure step.



---

## Phase 0.14 — Observability & Diagnostics Foundation

Phase 0.14 starts after the completed Provider Layer Consolidation phase.

Phase 0.13 made runtime composition explicit by consolidating the Provider Layer as the controlled entry point into application and domain services. That consolidation reduced orchestration ambiguity, but it did not make the runtime path fully diagnosable.

Phase 0.14 continues the foundation sequence by addressing the next operational gap: observability.

The narrative transition is:

```text
Phase 0.12: clarify model direction and mapping ownership
Phase 0.13: clarify orchestration ownership and runtime composition
Phase 0.14: clarify runtime visibility and diagnostic traceability
```

### Objective

Introduce observability and diagnostics as an internal foundation capability without changing public API behavior.

### Problem

The backend remains difficult to debug when a request crosses multiple internal layers because:

- request correlation is not yet standardized
- logs do not consistently carry request-scoped context
- internal errors do not consistently expose diagnostic metadata
- the HTTP → Provider → Application → Domain → Repository flow is not yet traceable as a coherent path

### Decision

Phase 0.14 introduces an Observability & Diagnostics Foundation.

The phase is intentionally limited to internal visibility. It does not introduce external metrics platforms, dashboards, Prometheus, OpenTelemetry or API contract changes.

### Scope

Included:

- request correlation
- request ID / trace propagation
- structured logging standardization
- internal error context enrichment
- flow tracing integration
- minimal diagnostics surface exposure

Excluded:

- public API behavior changes
- business logic changes
- external metrics backends
- dashboards
- Prometheus
- OpenTelemetry

### Subphases

- **0.14.0 — Phase Definition & Documentation Lock** ✅ Completed
- **0.14.1 — Correlation Model (Request ID / Trace)** ✅ Completed
- **0.14.2 — Logging Standardization** ✅ Completed
- **0.14.3 — Error Context Enrichment** ✅ Completed
- **0.14.4 — Flow Tracing Integration** ✅ Completed
- **0.14.5 — Diagnostics Surface Exposure** ✅ Completed
- **0.14.6 — Validation & Documentation** ✅ Completed

### Compatibility

Phase 0.14 must preserve:

- existing public HTTP/API contracts
- provider ownership established in Phase 0.13
- read/write model separation established in Phase 0.12
- existing behavior of authenticated application surfaces
- existing test compatibility


### Phase 0.14 Outcome

Phase 0.14 is completed as the Observability & Diagnostics Foundation for Stage 0.

The backend now preserves the public behavior stabilized by previous phases while making the internal runtime path observable. Requests can be correlated through the existing `X-Request-Id` boundary, runtime logs use the canonical `request_id` field, application errors can carry copied diagnostic context safely, HTTP/application lifecycle movement emits minimal flow events, and `/diagnostics` exposes a lightweight snapshot of the active observability capabilities.

This closes the gap left after Phase 0.13: the Provider Layer made runtime composition explicit, and Phase 0.14 makes that composition easier to inspect without changing business contracts, provider contracts, response envelopes or API behavior.

Final validation recorded for the phase:

```bash
go test ./...
```

The command completed successfully in the local project environment after 0.14.5 was applied.

---

## Phase 0.15 — Contract Hardening & Freeze

Phase 0.15 starts after the completed Observability & Diagnostics Foundation baseline.

Phase 0.12 clarified read/write model direction and mapping ownership. Phase 0.13 clarified provider-owned runtime composition. Phase 0.14 made that runtime path observable through request correlation, structured logging, safe diagnostic error context, minimal flow tracing and `/diagnostics`.

Phase 0.15 addresses the next Stage 0 risk: contracts are working, but they must now be audited, formalized and frozen so future backend and frontend evolution cannot drift silently.

### Objective

Formalize, validate and freeze the existing backend contracts.

The phase covers:

- HTTP contracts
- response schemas
- public error envelopes
- provider contracts
- freeze and versioning rules

### Scope

Included:

- HTTP contract audit
- error contract alignment
- provider contract validation
- response schema normalization
- contract freeze enforcement
- final validation and documentation reconciliation

Excluded:

- new features
- new routes
- business logic changes
- architecture changes
- additional observability scope
- external metrics or dashboard platforms

### Subphases

- **0.15.0 — Phase Definition & Documentation Lock** ✅ Completed
- **0.15.1 — HTTP Contract Audit** ✅ Completed
- **0.15.2 — Error Contract Alignment** ✅ Completed
- **0.15.3 — Provider Contract Validation** ✅ Completed
- **0.15.4 — Response Schema Normalization** ✅ Completed
- **0.15.5 — Contract Freeze Enforcement** ✅ Completed
- **0.15.6 — Validation & Documentation** ⏳ Pending

### 0.15.0 Result

0.15.0 opens the Contract Hardening & Freeze phase as a documentation-only lock.

Context inherited: the backend is stable after Phase 0.14.6.fix3, with provider composition explicit and runtime observability available.

Problem addressed: Phase 0.15 was pending definition, which prevented controlled contract audit work from starting. Without a phase lock, endpoint, error, provider and response contract work could begin without a shared scope.

Decision taken: define Phase 0.15 before any code changes. The phase is constrained to validating and freezing existing contracts, not adding behavior.

Concrete change: README, roadmap, phase status, documentation index, handoff and the new dedicated Phase 0.15 document now agree on the active phase, subphase order and boundaries.

Observable impact: the next correct step after 0.15.0 was 0.15.1 — HTTP Contract Audit, based only on the real route surface and existing contracts in the repository.



### 0.15.1 Result

0.15.1 completes the HTTP Contract Audit without changing Go source code.

Context inherited from 0.15.0: Phase 0.15 had been opened as Contract Hardening & Freeze after the completed Phase 0.14 observability baseline. Contract work needed to start with the real route surface, not with assumed endpoints or undocumented response expectations.

Problem addressed: HTTP contracts were present in router and handler code, but not yet captured as a dedicated 0.15 route inventory. That made later error alignment and response normalization vulnerable to missing legacy auth paths, `/api/v1` aliases, foundation endpoints or WebSocket-specific behavior.

Decision taken: document the current HTTP surface first. The audit records foundation routes, `/ws`, public auth routes, authenticated account/settings routes and authenticated wallet-management routes. It confirms that legacy auth paths and `/api/v1` auth paths are active paired contracts backed by the same handlers.

Concrete change: `docs/phase0_15_1_http_contract_audit.md` now records the HTTP inventory, status behavior, response families, error-envelope observations and non-blocking contract risks. Trunk documentation marks 0.15.0 and 0.15.1 as completed.

Observable impact: 39 registered HTTP route entries and 22 unique behavior contracts are now explicitly documented. The next correct step is 0.15.2 — Error Contract Alignment, using this audit as the public HTTP baseline.

### 0.15.2 Result

0.15.2 completes Error Contract Alignment with a narrow contract-hardening change.

Context inherited from 0.15.1: the HTTP route surface was explicitly audited and the public error-envelope baseline was documented before error alignment started.

Problem addressed: the standardized error model required `{error:{code,message,details}}`, but detail-free errors could omit `details` because the response error field used `omitempty` and the constructor did not initialize an empty details object.

Decision taken: keep the existing error envelope, error codes, messages, statuses and handler decisions, but make `error.details` a required JSON object. Detail-free errors now serialize `details: {}`.

Concrete change: `internal/core/errs/response_error.go` now always initializes and serializes `details`; `internal/core/errs/app_error_test.go` and `internal/core/httpx/error_test.go` add contract coverage for empty details and details-map copy behavior.

Observable impact: frontend consumers can rely on `error.details` being present as an object on canonical HTTP error responses. The next correct step is 0.15.3 — Provider Contract Validation.


### 0.15.3 Result

0.15.3 completes Provider Contract Validation with a focused internal-boundary hardening step.

Inherited context: 0.15.1 established the HTTP route baseline and 0.15.2 aligned the public error envelope so provider validation could focus on the internal contract seam instead of response-shape ambiguity.

Real problem: the auth provider boundary already existed through `AuthProvider`, `AuthSessionProvider`, `AuthenticatedAccountProvider` and `AuthWalletProvider`, and the auth application delegated into typed application/domain services. However, the contract was enforced mostly by usage. Cross-module provider dependencies were explicit in domain interfaces, but the build did not yet include a local compile-time guard proving that the concrete application and service types still satisfied those provider contracts.

Decision taken: add compile-time provider contract assertions without changing runtime behavior, route registration, handler logic, error mapping, response schemas or business rules.

Concrete change: `internal/modules/auth/application.go` now asserts that `*Application` satisfies `AuthProvider`, `*user.Service` satisfies `authdomain.UserProvider`, and `*usersettings.Service` satisfies `authdomain.UserSettingsProvider`.

Observable impact: future drift in HTTP-to-provider or auth-to-cross-module provider contracts now fails at compile time instead of surfacing later through handler behavior or frontend-facing responses. 0.15.4 has now normalized response serialization metadata and defensive fallback shape; the next correct step is 0.15.5 — Contract Freeze Enforcement.

### Phase 0.15.4 Result — Response Schema Normalization

0.15.4 completed Response Schema Normalization with a narrow compatibility-preserving serialization update.

Context inherited: 0.15.1 documented the HTTP route baseline, 0.15.2 aligned the public error envelope and 0.15.3 validated provider boundaries at compile time.

Real problem: response payload shapes were already compatible, but response serialization policy still had minor drift. Auth error responses used `application/json` while the core JSON writer used `application/json; charset=utf-8`, and the defensive timeout fallback JSON still lacked the mandatory `details` object.

Decision taken: normalize serialization details only. Do not introduce a success envelope, do not change any route, status code, payload field, error code, provider contract, domain rule or repository behavior.

Concrete change: auth error responses now use `application/json; charset=utf-8`, and the defensive timeout fallback keeps `{error:{code,message,details}}` aligned with the 0.15.2 contract.

Observable impact: JSON response metadata is aligned and defensive fallback shape is compatible with the frozen error envelope. The next correct step is 0.15.5 — Contract Freeze Enforcement.

### Phase 0.15.5 Result — Contract Freeze Enforcement

0.15.5 completed Contract Freeze Enforcement with explicit freeze rules and regression coverage for the currently audited Stage 0 contracts.

Context inherited: 0.15.1 documented the real HTTP route surface, 0.15.2 made the public error envelope stable, 0.15.3 added compile-time provider assertions and 0.15.4 normalized response serialization metadata.

Real problem: the contracts were audited and normalized, but future work still needed a concrete enforcement baseline. Without freeze rules and tests, route response shape, error envelope fields, content type policy or provider assumptions could drift silently.

Decision taken: freeze the current public and internal contract surface without adding behavior. Freeze enforcement is represented by documentation policy and regression tests that fail if core JSON endpoints or canonical auth error envelopes lose required fields or metadata.

Concrete change: `docs/phase0_15_5_contract_freeze_enforcement.md` defines the freeze policy, and `internal/core/httpx/contract_freeze_test.go` guards core status JSON responses plus protected auth canonical error envelopes.

Observable impact: future changes to frozen contract shape must be explicit, versioned or intentionally updated in tests and documentation. The next correct step is 0.15.6 — Validation & Documentation.
