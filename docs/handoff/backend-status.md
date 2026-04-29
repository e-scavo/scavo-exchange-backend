# Backend Status — SCAVO Exchange

## 🧠 Overview

This document represents the **current operational and architectural state** of the SCAVO Exchange Backend.

It is intended to:

* provide continuity across development phases
* allow safe context transfer between sessions
* serve as the single source of truth for backend status

---

## 📌 Current State

**Stage:** 0 — Foundation
**Phase:** 0.12 — Read / Write Model Separation
**Latest Completed Subphase:** 0.12.1.6 — Audit Consolidation & Closure
**Phase Status:** In Progress
**Next Planned Phase:** 0.12.1.1 — Model Inventory Extraction

---

## 🔐 Authentication Status

### Implemented

* password-based authentication (dev only)
* wallet-based authentication (EVM)

Wallet login flow supports:

* challenge generation
* signature verification
* challenge consumption (one-time use)
* wallet identity resolution
* durable user resolution / creation
* ownership enforcement
* JWT issuance

---

## 🧩 Identity Model Status

The system uses a **unified durable identity model**:

* wallet identities are persisted
* each wallet is linked to a durable user
* JWT reflects unified identity through `user_id`
* wallet metadata remains available for wallet-authenticated sessions

---

## 🏷️ Ownership Model Status

Ownership is a first-class persisted concept.

### Capabilities

* one user can own multiple wallets
* wallet ownership is persisted
* ownership metadata includes:

  * `user_id`
  * `linked_at`
  * `is_primary`
* primary-wallet uniqueness enforced
* wallet reassignment across users blocked

---

## 🔗 Wallet Ownership Status (0.4.14)

The backend now supports authenticated wallet-linking, authenticated wallet-owned account merge execution, authenticated primary-wallet switching, authenticated wallet detach-eligibility evaluation, authenticated detach execution for already eligible owned wallets, explicitly validated post-detach wallet reattachment semantics, and strict challenge-purpose isolation at the wallet-auth bootstrap boundary.

### Capabilities

* authenticated user can request wallet-link challenge
* challenge is persisted with:

  * `purpose = wallet_link`
  * `requested_by_user_id`
* authenticated user can verify link signature
* secondary wallet attaches to current user
* authenticated user can execute wallet-owned account merge after source-wallet signature
* authenticated user can explicitly switch the current primary wallet
* authenticated user can request detach-eligibility evaluation for one owned wallet
* authenticated user can execute detach for one already eligible owned wallet
* detached wallet identities remain reusable known identities after detach
* detached wallets can be reattached through the protected linking flow
* detached wallets can re-enter wallet-login bootstrap and resolve into wallet-owned user identities again
* explicit detach rejection reasons are returned when detach is not yet eligible
* updated wallet inventory is returned after successful linking, merge, primary switching, and detach execution

### Protections

* link challenge must belong to current authenticated user
* wrong-purpose challenge is rejected
* wallet already owned by another user is rejected
* wallet already linked to current user is rejected
* successful link does not issue a new JWT and does not implicitly merge accounts

---

## 🗄️ Persistence Status

### `auth_wallet_challenges`

* persistent
* expiration enforced
* single-use enforced
* now supports:

  * `purpose`
  * `requested_by_user_id`

### `auth_wallet_identities`

* persistent wallet registry
* ownership metadata included
* multi-wallet ownership supported

### `users`

* durable user representation
* wallet-backed users supported
* prepared for later multi-auth evolution

### `user_settings`

* dedicated authenticated user settings persistence
* separated from `users`
* keyed by `user_id`
* stores:

  * `preferences`
  * `created_at`
  * `updated_at`
* supports safe default-backed settings bootstrap without coupling settings to the core durable user record

---

## 🔌 API Status

### Auth endpoints

* `POST /auth/login`
* `POST /auth/wallet/challenge`
* `POST /auth/wallet/verify`
* `GET /auth/me`
* `PATCH /auth/me`
* `GET /auth/me/settings`
* `PATCH /auth/me/settings`
* `GET /auth/session`

`GET /auth/me` acts as the authenticated application bootstrap surface, while `PATCH /auth/me` now allows the authenticated user to update only `display_name` as minimal non-wallet profile metadata.

`GET /auth/me/settings` now provides the first dedicated authenticated settings contract, separated from `/auth/me`, with a minimal read-only response shape:

* `user_id`
* `version`
* `preferences`

The settings surface returns persisted data when available and safe defaults when no settings row exists yet, without creating or mutating settings during read.

The read shape remains additive and includes:

* `user_id`
* `auth_method`
* `wallet_id`
* `wallet_address`
* `chain`
* `primary_wallet`
* `wallets`
* `wallet_count`
* `active_wallet_count`
* `detached_wallet_count`
* `has_wallet_session`

### Ownership endpoints

* `GET /auth/wallets`

### Wallet-link endpoints

* `POST /auth/wallets/link/challenge`
* `POST /auth/wallets/link/verify`
* `POST /auth/account/merge/wallet/challenge`
* `POST /auth/account/merge/wallet/verify`

---

## 🧾 JWT Status

Tokens are:

* stateless
* short-lived
* self-contained

Claims include:

* `user_id`
* `wallet_id`
* `wallet_address`
* `auth_method`

Wallet linking uses the already authenticated JWT and does not mint a replacement token.

---

## ⚙️ Behavioral Guarantees

The backend guarantees:

* wallet challenges are one-time use
* wallet signatures are verified deterministically
* wallet identities are persistent
* user identity is durable
* ownership is consistent and enforced
* wallets cannot be reassigned across users
* primary wallet uniqueness is maintained
* link challenges are user-bound
* link, merge, and login challenge purposes are not interchangeable
* settings reads are side-effect free
* settings defaults are resolved without implicit persistence

---

## 📦 Wallet Inventory Read Model

`GET /auth/wallets` now returns an explicit wallet read model including:

* `id`
* `address`
* `user_id`
* `linked_at`
* `detached_at`
* `is_primary`
* `status`

This makes the public inventory contract lifecycle-aware while preserving the same ownership semantics already stabilized in previous subphases.

## 🧪 Testing Status

Validated at the design and code level through:

* `go test ./...`
* manual API testing procedures
* SQL verification queries

Coverage now includes:

* wallet auth flow
* identity linking
* ownership enforcement
* replay protection
* wallet-link challenge flow
* wallet-link verification flow
* ownership conflict rejection during link operations
* protected primary-wallet switching
* detach eligibility and execution
* detached-wallet reattachment semantics
* enriched wallet inventory read-model serialization
* wallet inventory filtering and sorting query semantics
* authenticated user settings default resolution
* authenticated user settings persisted read behavior
* authenticated user settings unauthorized access handling
* authenticated user settings mutation (merge-based)
* settings contract hardening validation
* rejection of null and invalid preference values
* shape compatibility protection

---

## ⚠️ Known Limitations

The system intentionally does **not** yet support:

* cursor pagination or broader reporting beyond the current paginated query contract
* wallet unlink operations
* cross-user wallet transfer
* merge between wallet identities and future auth methods
* refresh tokens
* token revocation
* persistent authenticated sessions
* typed/validated concrete preference fields inside the settings contract

---

## 🧭 Next Phase

### 0.5.5.1 — User Settings Hardening (Soft Contract Guard Layer)

Delivered:

* `PATCH /auth/me/settings` merge-based mutation
* recursive normalization of preferences
* rejection of null values
* rejection of invalid values
* shape compatibility protection
* soft observation of unknown keys

Expected next focus:

* introduce typed settings only when required
* preserve flexibility while improving consistency
* prepare transition toward next phase

---

## 🧾 Summary

At the end of Phase 0.4.16:

* authentication is stable
* identity is unified
* ownership is implemented and protected
* authenticated wallet linking is implemented
* authenticated wallet-owned account merge is implemented
* explicit primary-wallet switching is implemented
* wallet detach eligibility is implemented
* wallet detach execution is implemented for already eligible owned wallets
* detached wallet identities are explicitly reusable after detach
* detached wallet identities preserve minimal audit-ready lifecycle evidence through `detached_at`
* the authenticated wallet inventory endpoint now exposes an enriched lifecycle-aware read model
* the authenticated wallet inventory endpoint now supports filtering, sorting, and simple pagination semantics without changing ownership rules

Phase 0.4.21 hardens the authenticated wallet inventory query contract by making parameter combinations and defaults explicit (`order` now requires `sort`, and `sort=linked_at` defaults to ascending order when `order` is omitted). The implementation stays entirely in the handler/read-model layer and does not modify ownership, stores, or persistence.

Phase 0.4.22 closes the documentation gap around the authenticated wallet inventory response contract. The main endpoint example is now aligned with the implemented JSON fields, including returned-window metadata and bounded navigation hints, without changing domain, stores, or persistence.

Phase 0.4.23 closes the operator-facing examples layer for `GET /auth/wallets` by documenting concrete valid and invalid query patterns, plus bounded-window response examples, without changing domain, stores, persistence, or handler behavior.

Phase 0.4.24 closes the manual-validation layer around the authenticated wallet inventory endpoint. The implementation remains documentation-only, but operators now have an explicit checklist for validating base, filtered, sorted, paginated, bounded, unbounded, and invalid query scenarios against the real `GET /auth/wallets` contract.

Phase 0.4.25 prepares the authenticated wallet inventory for wallet-management consumption by exposing minimal actionability hints per listed wallet. The implementation stays in the read-model/handler layer, reuses the established detach-domain reasons, and leaves execution authority in the existing detach and primary-switch endpoints.

Phase 0.4.26 closes the consistency gap between the enriched wallet inventory read model and `POST /auth/wallets/detach/check`. The implementation adds handler-level coverage proving that inventory-side detach hints remain semantically aligned with detach-check eligibility and reasons for single-wallet and two-wallet ownership scenarios, while leaving detach authority in the existing check and execute endpoints.

Phase 0.4.28 closes the wallet-management read flow around the authenticated inventory and the existing primary / detach actions. The implementation is documentation-only, but it corrects the README phase summary and consolidates the real inventory → actionability hint → action/check endpoint → refreshed inventory flow so client and operator guidance now matches the authenticated wallet-management surface end to end.

Phase 0.4.30 consolidates the final wallet-management contract layer. Inventory, eligibility, execution, and refreshed post-action inventory are now documented as one coherent surface so future work can move beyond Phase 0.4 without reopening already-stabilized wallet-management semantics.

Phase 0.4.31 hardens the wallet-auth bootstrap boundary so `POST /auth/wallet/verify` can no longer consume `wallet_link` or `account_merge` challenges. Detached-wallet rebound remains supported, but only through the correct bootstrap-purpose challenge.

Phase 0.4.32 closes the last permissive challenge-purpose normalization gap. Controlled challenge creation still defaults empty purpose to `auth_bootstrap`, but unknown or malformed purpose values are now preserved as invalid runtime data and rejected by wallet verify/login, authenticated wallet link, and wallet-owned account merge.

Phase 0.5.2 introduces the first authenticated user metadata mutation by allowing `PATCH /auth/me` to update only `display_name` with strict validation, while keeping wallet ownership, email identity, and settings unchanged.

Phase 0.5.3 introduces the first dedicated authenticated settings contract through `GET /auth/me/settings`, backed by separated `user_settings` persistence and safe default resolution when no settings row exists yet. This keeps authenticated profile metadata and authenticated user configuration as distinct application surfaces.

Phase 0.5.4 introduces authenticated user settings mutation through `PATCH /auth/me/settings`, enabling merge-based preference updates without destructive overwrite.

Phase 0.5.5.1 introduces contract hardening to prevent structural drift in preferences, adding normalization, validation, and shape protection while preserving flexibility.

---

## Phase 0.4 Formal Closure

Phase 0.4 is now formally closed.

All authentication, wallet identity, ownership, wallet management, and challenge-purpose contracts are fully stabilized and aligned with the current implementation.

No further work should be added to Phase 0.4 unless a future ZIP proves a real regression or missing contract.

### Next Direction

* Begin a new phase for any further backend evolution
* Preserve all Phase 0.4 contracts as baseline
* Avoid reopening Phase 0.4 without a ZIP-validated reason


### 0.5.5.3 — User Settings Contract Surface Stabilization

This subphase makes the authenticated settings resource more self-descriptive without changing endpoint paths or introducing schema-heavy settings governance.

Delivered:

* `settings.created_at` is now exposed when persisted metadata exists
* `settings.updated_at` is now exposed when persisted metadata exists
* `GET /auth/me/settings` and `PATCH /auth/me/settings` now return the same stable resource-oriented settings shape
* zero-value timestamps remain omitted so default resolution without persistence metadata does not fabricate stored state

--- 
### 0.6.1 — Bootstrap Surface Boundary Clarification

This subphase clarifies the authenticated application bootstrap boundary without modifying public endpoint shapes, handlers, services, stores, or persistence.

Delivered:

* explicit contract-level tests now distinguish the role of `GET /auth/me` from `GET /auth/session`
* `/auth/me` is now formally protected as the authenticated bootstrap identity surface
* `/auth/session` is now formally protected as the token-derived session surface
* settings and wallet inventory remain separate authenticated resources under `GET /auth/me/settings` and `GET /auth/wallets`
* the authenticated surface can now move into 0.6.x alignment work without reopening 0.5.x contracts

### Phase 0.6 Direction

Phase 0.6 consolidates the authenticated application surface into a coherent bootstrap layer for frontend consumption.

Current boundary after 0.6.1:

* `GET /auth/me` → authenticated bootstrap identity surface
* `GET /auth/session` → authenticated session / token context surface
* `GET /auth/me/settings` → authenticated settings surface
* `GET /auth/wallets` → authenticated wallet inventory surface

0.6.1 is intentionally non-invasive:

* no public payloads were changed
* no router changes were introduced
* no product handlers were modified
* no stores or persistence contracts were changed
* the clarification is enforced through explicit boundary tests and aligned documentation

### Updated Next Direction

* continue with 0.6.2 — Authenticated Surface Contract Alignment
* preserve the semantic boundary now established between bootstrap identity and session context
* avoid re-opening Phase 0.5 unless the ZIP proves a real contract regression

### 0.6.2 — Authenticated Surface Contract Alignment

This subphase consolidates the authenticated bootstrap surface by aligning the shared identity and wallet-context fields exposed across related authenticated endpoints without changing public endpoint paths or breaking existing response contracts.

Delivered:

* introduced a shared authenticated context normalization layer used by the authenticated surface
* aligned `/auth/me` and `/auth/session` so common identity, wallet, and chain context now derive from the same internal source
* normalized empty wallet-context handling so partial wallet metadata does not drift between authenticated surfaces
* reinforced contract consistency between `/auth/me`, `/auth/session`, and `/auth/wallets` through explicit alignment tests
* preserved compatibility by keeping the public JSON envelopes and endpoint boundaries established in 0.6.1

### Phase 0.6 Direction

Phase 0.6 continues consolidating the authenticated application surface into a coherent bootstrap layer for frontend consumption.

Current boundary after 0.6.2:

* `GET /auth/me` → authenticated bootstrap identity surface with aligned shared context
* `GET /auth/session` → authenticated session / token context surface with aligned shared context
* `GET /auth/me/settings` → authenticated settings surface
* `GET /auth/wallets` → authenticated wallet inventory surface aligned with the authenticated bootstrap wallet view

0.6.2 remains intentionally compatibility-safe:

* no public endpoint paths were changed
* no public JSON field names were renamed or removed
* no wallet lifecycle rules were changed
* no settings resource redesign was introduced
* alignment is enforced through shared internal normalization and expanded contract tests

### Updated Next Direction

* continue with 0.6.3 — Session-Ready Bootstrap Read Model
* build a frontend-ready authenticated bootstrap read model on top of the now aligned authenticated surface
* preserve the semantic boundary from 0.6.1 and the contract alignment from 0.6.2 while avoiding unnecessary expansion of Phase 0

## Phase 0.6.3 — Session-Ready Bootstrap Read Model

### Summary

This subphase introduces a unified authenticated bootstrap read model, allowing frontend clients to retrieve all required authenticated surface data in a single request.

### Key Additions

- Introduced `GET /auth/bootstrap`
- Aggregates:
  - session
  - user
  - profile
  - settings
  - wallet snapshot
- Eliminates multi-request bootstrap pattern

### Resulting System State

After 0.6.3, the authenticated surface is:

- boundary-defined (0.6.1)
- contract-aligned (0.6.2)
- bootstrap-ready (0.6.3)

### Impact

- Simplifies frontend initialization
- Reduces network roundtrips
- Ensures cross-surface consistency at read time

### Next Step

0.6.4 — Application Surface Consistency Hardening

## Phase 0.6.4 — Application Surface Consistency Hardening

### Summary

This subphase finalizes the authenticated application surface by enforcing a consistent structural contract across all related endpoints.

### Key Improvements

- Introduced canonical wallet envelope (`items`, `total`)
- Preserved legacy compatibility (`wallets`)
- Hardened `/auth/bootstrap` structure to avoid drift
- Added structural tests to freeze contract expectations

### Resulting System State

After 0.6.4, the authenticated surface is:

- boundary-defined (0.6.1)
- contract-aligned (0.6.2)
- bootstrap-ready (0.6.3)
- structurally consistent (0.6.4)

### Phase 0.6 Closure

The system now provides a stable, deterministic, and consistent authenticated surface for frontend consumption.

### Next Step

Phase 0.7 — Application Layer Foundation

## Phase 0.7.1 — Application Layer Boundary Definition

### Summary

This subphase introduces the first explicit application layer, establishing a clear boundary between HTTP transport and use-case orchestration.

### Key Changes

- Introduced `Application` layer in `internal/modules/auth`
- Moved bootstrap orchestration out of HTTP handlers
- Established handler → application → service/store flow
- Preserved all public contracts and routes

### Architectural Impact

- HTTP handlers are now transport-only components
- Application layer owns use-case orchestration
- Existing services and stores remain the execution layer

### Resulting System State

After 0.7.1, the backend evolves from:

- handler-driven orchestration

to:

- structured application-layer orchestration

### Next Step

0.7.2 — Authenticated Surface Use Cases Extraction

## Phase 0.7.2 — Authenticated Surface Use Cases Extraction

### Summary

This subphase completes the migration of the authenticated surface core use cases into the application layer, removing handler-driven orchestration for login, session, and user identity flows.

### Key Changes

- Introduced application-level use cases:
  - `Application.Login(...)`
  - `Application.GetMe(...)`
  - `Application.GetSession(...)`
- Preserved `Application.GetBootstrap(...)` from 0.7.1
- Reduced HTTP handlers to transport-only responsibilities
- Centralized authenticated flow orchestration in application layer

### Architectural Impact

- Authenticated surface is now fully application-driven
- Handlers act purely as transport adapters
- Services and stores remain execution layer components

### Resulting System State

After 0.7.2, the backend evolves from:

- partial application-layer adoption (0.7.1)

to:

- full application-layer orchestration of authenticated surface (0.7.2)

### Next Step

0.7.3 — Wallet Management Use Cases Consolidation

## Phase 0.7.3 — Wallet Management Use Cases Consolidation

### Summary

This subphase completes the migration of all wallet management flows into the application layer, removing handler-driven orchestration entirely from the auth module.

### Key Changes

- Introduced application-level wallet use cases:
  - `Application.ListWallets(...)`
  - `Application.CreateWalletLinkChallenge(...)`
  - `Application.VerifyWalletLink(...)`
  - `Application.CreateWalletAccountMergeChallenge(...)`
  - `Application.VerifyWalletAccountMerge(...)`
  - `Application.SetPrimaryWallet(...)`
  - `Application.CheckWalletDetach(...)`
  - `Application.ExecuteWalletDetach(...)`
- Refactored all wallet-related HTTP handlers to delegate to application layer
- Consolidated wallet listing and management under a unified execution model

### Architectural Impact

- Auth module is now fully application-driven
- HTTP handlers act strictly as transport adapters
- Application layer owns all use-case orchestration
- Services and stores remain execution layer

### Resulting System State

After 0.7.3, the backend evolves from:

- partially application-driven (0.7.2)

to:

- fully application-driven auth module (0.7.3)

### Next Step

0.7.4 — Handler Simplification & Contract Preservation

## Phase 0.7.4 — Handler Simplification & Contract Preservation

### Summary

This subphase finalizes the application-layer foundation by simplifying all authenticated handlers and eliminating residual transport duplication, without altering public contracts.

### Key Changes

- Introduced unified HTTP transport helpers:
  - request decoding
  - authenticated claims extraction
  - error JSON writing
- Refactored all auth and wallet handlers to:
  - remove duplicated transport logic
  - follow a consistent execution pattern
- Ensured consistent error handling across endpoints

### Architectural Impact

- Auth module now exhibits fully clean separation:
  - Handlers → transport only
  - Application → orchestration
  - Services/Stores → execution
- No residual business logic remains in handlers

### Resulting System State

After 0.7.4, the backend reaches:

- fully application-driven auth module
- consistent and minimal HTTP transport layer
- stable and preserved external contracts

### Next Step

0.9 — API Versioning Strategy


## Phase 0.8 — Standardized Error Model

### Summary

This phase completed the transport-level error standardization initiated after the application-layer foundation. The backend now exposes a single structured error envelope backed by centralized app-error typing, reusable factories and dedicated hardening tests.

### Key Changes

- Introduced a canonical HTTP error envelope under `internal/core/httpx`
- Added centralized `AppError` typing and shared auth/internal factory helpers under `internal/core/errs`
- Standardized auth handlers and auth middleware on the centralized error model
- Added contract-hardening tests to freeze representative mapping and serialization behavior

### Resulting System State

After 0.8, the backend reaches:

- stable authenticated transport contracts
- centralized error normalization
- explicit contract-level hardening around error mapping behavior
- no cyclic-import regression while adopting the shared error system

### Next Step

0.9 — API Versioning Strategy

## Phase 0.9 — API Versioning Strategy

### Summary

This phase opens the next Stage 0 step by defining how the post-0.8 backend will evolve its public API contract without breaking the currently available route surface.

### Key Changes

- Declares path-based versioning as the canonical API strategy
- Reserves `/api/v1/...` as the canonical route namespace for the current authenticated surface
- Classifies the existing unversioned endpoints as legacy compatibility routes
- Freezes the current success-payload semantics and 0.8 error-envelope behavior as `v1` contract semantics
- Preserves the project rule that the frontend remains aligned only up to backend Phase 0.6 until Stage 0 is fully completed

### Architectural Impact

- Transport contract evolution becomes explicit
- Legacy compatibility is preserved without duplicating business logic
- Authorization in 0.10 gains a stable versioning foundation instead of defining transport policy implicitly

### Resulting System State

After 0.9.3, the backend is:

- functionally unchanged for current legacy consumers
- actually exposing canonical `/api/v1/...` routes in the runtime router
- explicitly freezing the current authenticated surface as canonical `v1` semantics
- still using one shared handler/application surface rather than duplicated route-specific logic
- documented with an explicit compatibility model for the remaining Stage 0 work

### Next Step

0.9.4 — Version-aware Contract Testing

## Phase 0.9.2 — Router Versioning Foundation

### Summary

This subphase turns the 0.9.1 versioning policy into real transport exposure by projecting the stabilized auth/authenticated surface into the canonical `/api/v1/...` namespace.

### Key Changes

- Added canonical `/api/v1/auth/...` route registration in `internal/core/httpx/router.go`
- Preserved the existing `/auth/...` compatibility surface
- Centralized auth route registration so legacy and canonical paths reuse the same handlers
- Preserved auth middleware wrapping and existing transport concerns across both route spaces

### Architectural Impact

- Versioning is now implemented in router composition rather than only described in docs
- The backend can move into contract freezing and equivalence hardening against a real canonical runtime surface
- Frontend alignment remains unchanged because the legacy route space still exists and remains supported

### Resulting System State

After 0.9.2, the system has:

- a documented versioning policy
- a real canonical `v1` route namespace
- preserved legacy compatibility routes
- no business-flow duplication introduced by versioning

### Next Step

0.9.3 — Authenticated Surface Version Freezing

## Phase 0.9.3 — Authenticated Surface Version Freezing

### Summary

This subphase freezes the currently stabilized authenticated behavior as the canonical `v1` authenticated contract without changing business flows, payload structure or middleware behavior.

### Key Changes

- Declared bootstrap, profile, settings, session and wallet inventory/read behavior to be part of the frozen authenticated `v1` surface
- Clarified that legacy `/auth/...` and canonical `/api/v1/auth/...` entry paths are transport aliases over the same authenticated contract
- Preserved the current Phase 0.8 standardized error envelope as part of the same `v1` semantics
- Kept exhaustive route-by-route equivalence hardening deferred to 0.9.4

### Architectural Impact

- Canonical routing is now backed by an explicit surface freeze rather than only by route exposure
- Later authorization work can target a stable authenticated contract boundary instead of a merely projected router surface
- Frontend alignment remains unchanged because the legacy route surface is still supported during Stage 0

### Resulting System State

After 0.9.3, the system has:

- a documented versioning policy
- a real canonical `v1` route namespace
- a frozen authenticated `v1` contract surface
- preserved legacy compatibility routes
- no business-flow duplication introduced by versioning

### Next Step

0.9.4 — Version-aware Contract Testing

## Phase 0.9.4 — Version-aware Contract Testing

### Summary

This subphase turns the declared authenticated-surface `v1` freeze into representative automated transport regression protection so legacy and canonical route spaces do not silently drift while both remain supported.

### Key Changes

- Added router-level versioning tests under `internal/core/httpx` rather than introducing route-specific business logic
- Verified legacy/canonical helper path composition and representative parity for protected error scenarios
- Verified representative public success-shape compatibility for wallet challenge creation across `/auth/...` and `/api/v1/auth/...`
- Kept the scope focused on representative contract parity rather than attempting a new API version or a full endpoint-by-endpoint matrix in the same step

### Architectural Impact

- The canonical route surface is now not only exposed and frozen; it is also protected by automated transport-level checks
- The risk of accidental divergence between legacy and canonical Stage 0 entry paths is reduced before authorization work begins
- Frontend alignment remains unchanged because the current legacy route surface is still supported and the frontend continues to target backend Phase 0.6 semantics during Stage 0

### Resulting System State

After 0.9.4, the system has:

- a documented versioning policy
- a real canonical `v1` route namespace
- a frozen authenticated `v1` contract surface
- representative version-aware regression protection for key legacy/canonical paths
- preserved legacy compatibility routes

### Next Step

0.9.5 — Documentation Consolidation

## Phase 0.9.5 — Documentation Consolidation

### Summary

This subphase closes Phase 0.9 by consolidating the trunk documentation set around the already-delivered versioning policy, router foundation, authenticated-surface freeze and representative version-aware contract testing.

### Key Changes

- Updated trunk documents so Phase 0.9 is consistently represented as completed
- Recorded 0.9.4 as completed and 0.9.5 as the formal closure step for the phase
- Aligned README, roadmap, phase-status, architecture, testing and handoff around the same explanation of canonical `/api/v1/auth/...`, legacy `/auth/...` and the frontend Phase 0.6 alignment rule
- Corrected operational handoff state so the next planned phase is now accurately represented as 0.10 — Authorization Layer

### Architectural Impact

- The backend now has a closed and coherent architectural record for its transport-evolution model
- Authorization in 0.10 can start from an unambiguous versioning and contract-governance baseline rather than from fragmented status narratives
- Documentation no longer implies that versioning work is still partially pending when the router, freeze and representative tests are already in place

### Resulting System State

After 0.9.5, the system has:

- a completed Phase 0.9 versioning foundation
- a documented and implemented canonical `v1` route surface
- a frozen authenticated `v1` contract boundary
- representative version-aware regression protection
- a trunk documentation set aligned to the same backend/frontend evolution rule

### Next Step

0.10 — Authorization Layer



## Phase 0.10.2 — Authorization Context & Middleware

### Current State

The backend now carries authorization context through authenticated HTTP requests.

Phase 0.10.1 introduced the static authorization model (`Role`, `Permission`, role → permission mapping and `AuthorizationSubject`). Phase 0.10.2 extends that work by attaching a normalized authorization subject to the authenticated request context after JWT authentication succeeds.

### Delivered in 0.10.2

- authorization context helpers in `internal/core/authorization`
- deterministic projection from authenticated JWT claims to `AuthorizationSubject`
- dedicated `HydrateAuthorization()` middleware in `internal/core/httpx`
- middleware integration across protected legacy and canonical auth routes
- focused transport/core tests for authorization-context propagation

### Important Boundary

This subphase does not yet evaluate permissions or deny requests. Endpoint behavior, response payloads and error contracts remain unchanged. The backend is now prepared for centralized policy evaluation, which is the next planned step.

### Next

0.10.3 — Policy Evaluation Layer


## Phase 0.10.3 — Policy Evaluation Layer

### Current State

The backend now has a centralized authorization decision boundary on top of the model and context work delivered in 0.10.1 and 0.10.2.

Phase 0.10.3 introduces explicit action/resource vocabulary and a dedicated evaluator so authorization questions no longer need to be answered by reading raw claims or manually reproducing role-permission logic in transport code.

### Delivered in 0.10.3

- `Action` and `Resource` vocabularies in `internal/core/authorization`
- centralized action/resource → permission projection helpers
- `PolicyEvaluator` with `Evaluate(...)` and `Can(...)` APIs
- permission evaluation against normalized `AuthorizationSubject` data
- focused tests covering allowed, denied and unknown policy decisions

### Important Boundary

This subphase still does not enforce permissions at the HTTP endpoint level. Handlers, payloads and error contracts remain unchanged. What changed is that the backend now has a stable internal policy API that can be used by the next subphase without introducing handler-local authorization logic.

### Next

0.10.4 — Endpoint-Level Enforcement


## Phase 0.10.5 — Documentation & Contract Consolidation

Phase 0.10.5 does not add new authorization runtime behavior. It closes the authorization layer by consolidating the repository narrative around what was already delivered in 0.10.1 through 0.10.4.

This subphase aligns the trunk documentation so that README, roadmap, phase-status, handoff, architecture and the dedicated Phase 0.10 document all describe the same state:

- the authorization model exists
- authorization subject hydration exists in the authenticated request lifecycle
- centralized policy evaluation exists
- selected endpoints already enforce authorization decisions
- the first enforcement slice is intentionally progressive rather than universal

Phase 0.10 is therefore complete. The next architectural step should build on this now-documented authorization foundation rather than re-explaining or re-stabilizing it.

### Next

0.11 — Domain Module Pattern


## Phase 0.11 — Domain Module Pattern

Phase 0.11 is now the current structural Stage 0 phase. It follows the completed authorization layer and focuses on standardizing the internal organization of the current domain-facing modules without changing the public API contract.

### Current Structural Direction

The adopted module pattern is centered on:

- `http/` for transport handlers and DTOs
- `app/` for use-case orchestration
- `domain/` for models, invariants and internal contracts
- `repository/` for persistence-facing implementation boundaries when that separation adds real clarity

### Modules in Scope

- `internal/modules/user`
- `internal/modules/usersettings`
- `internal/modules/auth`

### Phase State

Phase 0.11 is **completed**.

Completed subphases:

- ✔ 0.11.1 — Domain Module Pattern Definition
- ✔ 0.11.2 — User Module Refactor
- ✔ 0.11.3 — UserSettings Module Refactor
- ✔ 0.11.4 — Auth Module Alignment
- ✔ 0.11.5 — Cross-Module Contract Consolidation
  - ✔ 0.11.5.0 — Subphase Definition & Documentation Alignment
  - ✔ 0.11.5.1 — Dependency Mapping
  - ✔ 0.11.5.2 — Contract Extraction
  - ✔ 0.11.5.3 — Interface Alignment
  - ✔ 0.11.5.4 — Runtime Compatibility Validation
- ✔ 0.11.6 — Documentation & Phase Closure

### Current Operational Meaning

0.11.1 established the repository-level architectural definition of the Domain Module Pattern, including responsibility split, dependency direction and ownership boundaries.

0.11.2 applied that pattern to `internal/modules/user`. The user module is now the first concrete runtime reference implementation of the pattern, exposing explicit `app`, `domain` and `repository` boundaries while preserving current runtime behavior and test compatibility.

0.11.3 applied the same pattern to `internal/modules/usersettings`, preserving settings-specific semantics and backward-compatible root-package access.

0.11.4 completed the conservative auth alignment. `auth/app` now owns runtime/application orchestration, `auth/domain` owns wallet contracts, root `auth` remains a compatibility/transport surface, and `auth/repository` is explicitly documented as a transitional façade while root stores remain active runtime implementations.

0.11.5 completed cross-module contract consolidation. `auth/app` now coordinates with `user` and `usersettings` through explicit minimal contracts while root `auth` preserves backward-compatible wiring and nil-service semantics.

### Completed 0.11.5 Work

0.11.5 — Cross-Module Contract Consolidation is completed. It mapped the real cross-module dependencies between `auth`, `user` and `usersettings`, extracted minimal contracts, aligned `auth/app` to those contracts and validated runtime compatibility without changing public HTTP behavior.

Delivered internal steps:

- 0.11.5.0 — Subphase Definition & Documentation Alignment
- 0.11.5.1 — Dependency Mapping
- 0.11.5.2 — Contract Extraction
- 0.11.5.3 — Interface Alignment
- 0.11.5.4 — Runtime Compatibility Validation

### Next

0.12 — Read/Write Model Separation

---

## Phase 0.12 — Read / Write Model Separation

Phase 0.12 is now the active Stage 0 structural phase.

### Current State

0.12.0 — Phase Definition & Documentation Lock is completed.

0.12.1 — Model Classification Audit is completed through 0.12.1.6.

0.12.2.6 — Read Model Extraction Documentation & Closure is completed.

### Operational Meaning

The backend is still externally aligned with the completed 0.11 runtime behavior. No code has been changed for 0.12.0. The repository now has a documented plan to separate read models, write models and domain-owned structures without changing public API behavior.

### Subphase State

- ✔ 0.12.0 — Phase Definition & Documentation Lock
- ✔ 0.12.1 — Model Classification Audit
  - ✔ 0.12.1.0 — Definition & Documentation Lock
  - ✔ 0.12.1.1 — Model Inventory Extraction
  - ✔ 0.12.1.2 — Model Classification
  - ✔ 0.12.1.3 — Cross-Layer Usage Analysis
  - ✔ 0.12.1.4 — Problem Detection & Risk Mapping
  - ✔ 0.12.1.5 — Target Separation Definition
  - ✔ 0.12.1.6 — Audit Consolidation & Closure
- ✔ 0.12.2 — Read Model Extraction
  - ✔ 0.12.2.0 — Definition & Documentation Lock
  - ✔ 0.12.2.1 — Read Model Design
  - ✔ 0.12.2.2 — Read Model Implementation
  - ✔ 0.12.2.3 — Domain/Application → Read Mapping
  - ✔ 0.12.2.4 — Response Alignment
  - ✔ 0.12.2.5 — Validation & Compatibility
  - ✔ 0.12.2.6 — Documentation & Closure
- ✔ 0.12.3 — Write Model Isolation
  - ✔ 0.12.3.0 — Definition & Documentation Lock
  - ✔ 0.12.3.1 — Write Model Design
  - ✔ 0.12.3.2 — Write Model Implementation
  - ✔ 0.12.3.3 — Write → Domain Mapping
  - ✔ 0.12.3.4 — Handler Alignment
  - ✔ 0.12.3.5 — Validation & Compatibility
  - ✔ 0.12.3.6 — Documentation & Closure
- 🚧 0.12.4 — Mapping Layer Introduction
  - ✔ 0.12.4.0 — Definition & Documentation Lock
  - ⬜ 0.12.4.1 — Mapping Layer Design
  - ⬜ 0.12.4.2 — Mapping Layer Implementation
  - ⬜ 0.12.4.3 — Mapping Consolidation
  - ⬜ 0.12.4.4 — Application Refactor
  - ⬜ 0.12.4.5 — Validation & Compatibility
  - ⬜ 0.12.4.6 — Documentation & Closure
- ⬜ 0.12.5 — Contract Alignment
- ⬜ 0.12.6 — Documentation & Phase Closure

### Next Required Action

Start 0.12.4 — Mapping Layer Introduction using the completed 0.12.1 audit artifacts, completed 0.12.2 read model extraction and completed 0.12.3 write model isolation as the baseline. Both read-side response compatibility and write-side input compatibility are preserved and validated.

### 0.12.1 Audit Closure Summary

0.12.1 is completed. The repository now contains the complete audit chain:

- `docs/phase0_12_1_model_inventory.md`
- `docs/phase0_12_1_model_classification.md`
- `docs/phase0_12_1_cross_layer_usage_analysis.md`
- `docs/phase0_12_1_problem_detection_risk_mapping.md`
- `docs/phase0_12_1_target_separation_definition.md`
- `docs/phase0_12_1_audit_consolidation_closure.md`

The next phase must not start by guessing model ownership. It must use these audit artifacts to drive read model extraction.

### 0.12.2.6 Read Model Extraction Closure Summary

0.12.2 is completed. It introduced explicit read model packages, Domain/Application → Read mapping functions, response alignment and validation documentation. Public response compatibility remains preserved.

The completed step is 0.12.4.0 — Mapping Layer Introduction Definition & Documentation Lock. The next planned step is 0.12.4.1 — Mapping Layer Design. It must design centralized mapper ownership for the explicit mapping directions introduced in 0.12.2 and 0.12.3 without changing public HTTP contracts.

### 0.12.3.0 Write Model Isolation Definition Lock

0.12.3 is completed. It introduced write model packages, domain write inputs, Write → Domain mapping, handler alignment and validation documentation. The next step is 0.12.4 — Mapping Layer Introduction.

### 0.12.4.0 Mapping Layer Introduction Definition Lock

0.12.4.0 is completed as a documentation-only subphase. It defines the centralized mapper package direction, locks the 0.12.4 internal sequence and preserves the compatibility rule that public request and response contracts must remain unchanged during mapper consolidation.
