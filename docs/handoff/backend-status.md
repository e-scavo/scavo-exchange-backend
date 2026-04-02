# Backend Status — SCAVO Exchange

## 🧠 Overview

This document represents the **current operational and architectural state** of the SCAVO Exchange Backend.

It is intended to:

- provide continuity across development phases
- allow safe context transfer between sessions
- serve as the single source of truth for backend status

---

## 📌 Current State

**Stage:** 0 — Foundation  
**Phase:** 0.4 — Auth and User Stabilization  
**Latest Completed Subphase:** 0.4.23 — Wallet Inventory Query Examples Closure

---

## 🔐 Authentication Status

### Implemented

- password-based authentication (dev only)
- wallet-based authentication (EVM)

Wallet login flow supports:

- challenge generation
- signature verification
- challenge consumption (one-time use)
- wallet identity resolution
- durable user resolution / creation
- ownership enforcement
- JWT issuance

---

## 🧩 Identity Model Status

The system uses a **unified durable identity model**:

- wallet identities are persisted
- each wallet is linked to a durable user
- JWT reflects unified identity through `user_id`
- wallet metadata remains available for wallet-authenticated sessions

---

## 🏷️ Ownership Model Status

Ownership is a first-class persisted concept.

### Capabilities

- one user can own multiple wallets
- wallet ownership is persisted
- ownership metadata includes:
  - `user_id`
  - `linked_at`
  - `is_primary`
- primary-wallet uniqueness enforced
- wallet reassignment across users blocked

---

## 🔗 Wallet Ownership Status (0.4.14)

The backend now supports authenticated wallet-linking, authenticated wallet-owned account merge execution, authenticated primary-wallet switching, authenticated wallet detach-eligibility evaluation, authenticated detach execution for already eligible owned wallets, and explicitly validated post-detach wallet reattachment semantics.

### Capabilities

- authenticated user can request wallet-link challenge
- challenge is persisted with:
  - `purpose = wallet_link`
  - `requested_by_user_id`
- authenticated user can verify link signature
- secondary wallet attaches to current user
- authenticated user can execute wallet-owned account merge after source-wallet signature
- authenticated user can explicitly switch the current primary wallet
- authenticated user can request detach-eligibility evaluation for one owned wallet
- authenticated user can execute detach for one already eligible owned wallet
- detached wallet identities remain reusable known identities after detach
- detached wallets can be reattached through the protected linking flow
- detached wallets can re-enter wallet-login bootstrap and resolve into wallet-owned user identities again
- explicit detach rejection reasons are returned when detach is not yet eligible
- updated wallet inventory is returned after successful linking, merge, primary switching, and detach execution

### Protections

- link challenge must belong to current authenticated user
- wrong-purpose challenge is rejected
- wallet already owned by another user is rejected
- wallet already linked to current user is rejected
- successful link does not issue a new JWT and does not implicitly merge accounts

---

## 🗄️ Persistence Status

### `auth_wallet_challenges`
- persistent
- expiration enforced
- single-use enforced
- now supports:
  - `purpose`
  - `requested_by_user_id`

### `auth_wallet_identities`
- persistent wallet registry
- ownership metadata included
- multi-wallet ownership supported

### `users`
- durable user representation
- wallet-backed users supported
- prepared for later multi-auth evolution

---

## 🔌 API Status

### Auth endpoints
- `POST /auth/login`
- `POST /auth/wallet/challenge`
- `POST /auth/wallet/verify`
- `GET /auth/me`
- `GET /auth/session`

### Ownership endpoints
- `GET /auth/wallets`

### Wallet-link endpoints
- `POST /auth/wallets/link/challenge`
- `POST /auth/wallets/link/verify`
- `POST /auth/account/merge/wallet/challenge`
- `POST /auth/account/merge/wallet/verify`

---

## 🧾 JWT Status

Tokens are:

- stateless
- short-lived
- self-contained

Claims include:

- `user_id`
- `wallet_id`
- `wallet_address`
- `auth_method`

Wallet linking uses the already authenticated JWT and does not mint a replacement token.

---

## ⚙️ Behavioral Guarantees

The backend guarantees:

- wallet challenges are one-time use
- wallet signatures are verified deterministically
- wallet identities are persistent
- user identity is durable
- ownership is consistent and enforced
- wallets cannot be reassigned across users
- primary wallet uniqueness is maintained
- link challenges are user-bound
- link and login challenge purposes are not interchangeable

---


## 📦 Wallet Inventory Read Model

`GET /auth/wallets` now returns an explicit wallet read model including:

- `id`
- `address`
- `user_id`
- `linked_at`
- `detached_at`
- `is_primary`
- `status`

This makes the public inventory contract lifecycle-aware while preserving the same ownership semantics already stabilized in previous subphases.

## 🧪 Testing Status

Validated at the design and code level through:

- `go test ./...`
- manual API testing procedures
- SQL verification queries

Coverage now includes:

- wallet auth flow
- identity linking
- ownership enforcement
- replay protection
- wallet-link challenge flow
- wallet-link verification flow
- ownership conflict rejection during link operations
- protected primary-wallet switching
- detach eligibility and execution
- detached-wallet reattachment semantics
- enriched wallet inventory read-model serialization
- wallet inventory filtering and sorting query semantics

---

## ⚠️ Known Limitations

The system intentionally does **not** yet support:

- cursor pagination or broader reporting beyond the current paginated query contract
- wallet unlink operations
- cross-user wallet transfer
- merge between wallet identities and future auth methods
- refresh tokens
- token revocation
- persistent authenticated sessions

---

## 🧭 Next Phase

### 0.4.29 — Wallet Detach Execute Read Consistency

Delivered:

- handler-level consistency coverage now proves that a wallet exposed as detachable in inventory can be detached successfully through `POST /auth/wallets/detach`
- detach execution payload and refreshed inventory are now explicitly protected against semantic drift relative to pre-detach `can_detach` / `detach_block_reasons`
- documentation clarifies the inventory → detach execute → refreshed inventory contract while keeping inventory hints advisory and detach execution authoritative
- no changes to runtime detach rules, stores, or persistence

Expected next focus:

- continue incremental wallet-management readiness only if the next ZIP shows a real need
- preserve backward compatibility of `GET /auth/wallets`
- keep future work ownership-safe and read-model focused

---

## 📌 Development Guidelines

When continuing development:

- do not break ownership invariants
- do not allow wallet reassignment
- do not bypass challenge-purpose checks
- preserve backward compatibility of wallet login
- keep challenge-to-user binding explicit in wallet-management flows
- maintain documentation alignment with implementation

---

## 🧾 Summary

At the end of Phase 0.4.16:

- authentication is stable
- identity is unified
- ownership is implemented and protected
- authenticated wallet linking is implemented
- authenticated wallet-owned account merge is implemented
- explicit primary-wallet switching is implemented
- wallet detach eligibility is implemented
- wallet detach execution is implemented for already eligible owned wallets
- detached wallet identities are explicitly reusable after detach
- detached wallet identities preserve minimal audit-ready lifecycle evidence through `detached_at`
- the authenticated wallet inventory endpoint now exposes an enriched lifecycle-aware read model
- the authenticated wallet inventory endpoint now supports filtering, sorting, and simple pagination semantics without changing ownership rules



Phase 0.4.21 hardens the authenticated wallet inventory query contract by making parameter combinations and defaults explicit (`order` now requires `sort`, and `sort=linked_at` defaults to ascending order when `order` is omitted). The implementation stays entirely in the handler/read-model layer and does not modify ownership, stores, or persistence.


Phase 0.4.22 closes the documentation gap around the authenticated wallet inventory response contract. The main endpoint example is now aligned with the implemented JSON fields, including returned-window metadata and bounded navigation hints, without changing domain, stores, or persistence.


Phase 0.4.23 closes the operator-facing examples layer for `GET /auth/wallets` by documenting concrete valid and invalid query patterns, plus bounded-window response examples, without changing domain, stores, persistence, or handler behavior.


Phase 0.4.24 closes the manual-validation layer around the authenticated wallet inventory endpoint. The implementation remains documentation-only, but operators now have an explicit checklist for validating base, filtered, sorted, paginated, bounded, unbounded, and invalid query scenarios against the real `GET /auth/wallets` contract.



Phase 0.4.25 prepares the authenticated wallet inventory for wallet-management consumption by exposing minimal actionability hints per listed wallet. The implementation stays in the read-model/handler layer, reuses the established detach-domain reasons, and leaves execution authority in the existing detach and primary-switch endpoints.



Phase 0.4.26 closes the consistency gap between the enriched wallet inventory read model and `POST /auth/wallets/detach/check`. The implementation adds handler-level coverage proving that inventory-side detach hints remain semantically aligned with detach-check eligibility and reasons for single-wallet and two-wallet ownership scenarios, while leaving detach authority in the existing check and execute endpoints.


Phase 0.4.28 closes the wallet-management read flow around the authenticated inventory and the existing primary / detach actions. The implementation is documentation-only, but it corrects the README phase summary and consolidates the real inventory → actionability hint → action/check endpoint → refreshed inventory flow so client and operator guidance now matches the authenticated wallet-management surface end to end.
