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
**Phase:** 0.6 — Authenticated Application Bootstrap Consolidation & Session-Ready Surface
**Latest Completed Subphase:** 0.6.1 — Bootstrap Surface Boundary Clarification

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
