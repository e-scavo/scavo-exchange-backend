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
**Phase:** 0.6 — Authenticated Application Bootstrap Consolidation & Session-Ready Surface  
**Current Subphase:** **0.6.2 — Authenticated Surface Contract Alignment**

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
