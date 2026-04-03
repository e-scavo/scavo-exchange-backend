# 📊 Phase Status

## Stage 0 — Foundation

### Phase 0.1 — Initial Project Bootstrap
Status: ✅ Completed

### Phase 0.2 — Core Infrastructure
Status: ✅ Completed

### Phase 0.3 — User and Platform Base
Status: ✅ Completed

### Phase 0.4 — Auth and User Stabilization
Status: ✅ Completed

---

## Phase 0.4 Subphase Status

| Subphase | Description | Status |
|----------|-------------|--------|
| 0.4.1 | Auth base setup | ✅ Completed |
| 0.4.2 | JWT implementation and auth normalization | ✅ Completed |
| 0.4.3 | Auth endpoints stabilization | ✅ Completed |
| 0.4.4 | Wallet challenge contract and nonce bootstrap | ✅ Completed |
| 0.4.5 | Wallet signature verification and token issuance | ✅ Completed |
| 0.4.6 | Wallet identity persistence and durable challenge storage | ✅ Completed |
| 0.4.7 | Wallet ↔ user linking and unified identity model | ✅ Completed |
| 0.4.8 | Account consolidation and multi-wallet ownership foundations | ✅ Completed |
| 0.4.9 | User-driven wallet linking contract and protected account merge preparation | ✅ Completed |
| 0.4.10 | User-driven wallet-owned account merge execution | ✅ Completed |
| 0.4.11 | Primary wallet management and ownership safety hardening | ✅ Completed |
| 0.4.12 | Wallet detach contract preparation and ownership guardrails | ✅ Completed |
| 0.4.13 | Protected wallet detach execution | ✅ Completed |
| 0.4.14 | Detached wallet reattachment semantics and lifecycle clarification | ✅ Completed |
| 0.4.15 | Detached identity audit readiness | ✅ Completed |
| 0.4.16 | Wallet identity read model enrichment | ✅ Completed |
| 0.4.17 | Wallet inventory query filtering and sorting | ✅ Completed |
| 0.4.18 | Wallet inventory pagination and windowed response | ✅ Completed |
| 0.4.19 | Wallet inventory navigation metadata | ✅ Completed |
| 0.4.20 | Wallet inventory cursorless navigation hints | ✅ Completed |
| 0.4.21 | Wallet inventory query parameter contract hardening | ✅ Completed |
| 0.4.22 | Wallet inventory response contract clarification | ✅ Completed |
| 0.4.23 | Wallet inventory query examples closure | ✅ Completed |
| 0.4.24 | Wallet inventory manual validation closure | ✅ Completed |
| 0.4.25 | Wallet actionability read model preparation | ✅ Completed |
| 0.4.26 | Wallet detach check read consistency | ✅ Completed |
| 0.4.27 | Wallet primary switch read consistency | ✅ Completed |
| 0.4.28 | Wallet management read flow closure | ✅ Completed |
| 0.4.29 | Wallet detach execute read consistency | ✅ Completed |
| 0.4.30 | Wallet management contract consolidation | ✅ Completed |
| 0.4.31 | Wallet auth bootstrap purpose enforcement | ✅ Completed |

---

## ✅ Phase 0.4.15 Closure Summary

Phase 0.4.15 closes the first minimal detached-identity audit gap without introducing heavy lifecycle redesign, event sourcing, or archival semantics.

The backend now persists `detached_at` on wallet identities whenever a detach is executed. That timestamp remains present if the wallet is later reattached through authenticated linking or rebounds through wallet-login bootstrap, allowing the system to distinguish reusable previously detached identities from identities that were never detached.

### Delivered in 0.4.15

- minimal detached-wallet lifecycle metadata via `detached_at`
- PostgreSQL migration for detached-identity audit readiness
- in-memory and PostgreSQL persistence support for `detached_at`
- detach execution updated to stamp `detached_at`
- reattachment and wallet-login rebound coverage updated to prove detached metadata survives reuse
- documentation updates aligning detached-wallet reuse with minimal audit readiness

---

## 🔍 Functional Result

The system now supports the following detached-identity lifecycle sequence:

1. user detaches one already eligible owned wallet
2. backend clears `user_id`, `linked_at`, and `is_primary` from that wallet identity
3. backend stamps `detached_at` on that wallet identity
4. wallet identity remains known to the backend by address and wallet identity ID
5. authenticated user can later reattach that wallet again through `POST /auth/wallets/link/challenge` + `POST /auth/wallets/link/verify`
6. detached wallet can also re-enter `POST /auth/wallet/verify` and resolve back into a wallet-owned user identity
7. the detached timestamp remains available as minimal lifecycle evidence even after reuse

---

## ❌ Not Included in 0.4.14

The following items remain intentionally out of scope:

- detached-by-user audit metadata
- multi-event detached lifecycle history
- archival / soft-delete markers for detached wallets
- recovery or dispute workflows around detached ownership
- automatic primary replacement for risky detach cases
- merge between wallet identities and future auth methods
- refresh tokens
- revocation flows
- persistent authenticated session storage

---

## ⏭️ Next Phase

### Next Expected Evolution

- optional `detached_by_user_id` or richer detach history if future audit scope requires it
- queryable lifecycle reporting if operational observability later needs it
- preserve current reusable detached-wallet semantics while preparing richer lifecycle observability


---

## ✅ Phase 0.4.16 Closure Summary

Phase 0.4.16 enriches the authenticated wallet inventory contract without changing ownership behavior or persistence rules.

The backend now exposes a dedicated wallet read model through `GET /auth/wallets`, including:

- `id`
- `address`
- `user_id`
- `linked_at`
- `detached_at`
- `is_primary`
- `status`

### Delivered in 0.4.16

- explicit `WalletReadModel` projection for authenticated wallet inventory
- lifecycle-aware visibility of `linked_at`, `detached_at`, and `is_primary`
- derived `status` field for wallet inventory responses
- handler-level validation for active wallet inventory serialization
- handler-level validation proving detached-then-reattached wallets still expose `detached_at`
- documentation updates aligning the public API contract with the real wallet identity lifecycle model

---

## 🔍 Functional Result

The system now supports the following authenticated inventory behavior:

1. authenticated user calls `GET /auth/wallets`
2. backend resolves all wallet identities currently owned by that durable user
3. backend projects each identity through an explicit read model
4. lifecycle-aware fields are returned, including current ownership state and minimal historical detach evidence
5. clients can observe both current ownership and preserved detached-wallet lifecycle evidence without requiring additional endpoints

---

## ❌ Not Included in 0.4.16

The following items remain intentionally out of scope:

- filtering or query parameters for wallet inventory
- pagination
- search
- admin reporting endpoints
- richer detached-identity history endpoints
- event sourcing or audit tables
- additional ownership mutations

---

## ⏭️ Next Phase

### Next Expected Evolution

- wallet inventory pagination on top of the filtered and sortable read model
- additive response metadata for windowed wallet inventory delivery
- avoid reworking ownership invariants already stabilized in Phase 0.4



## ✅ Phase 0.4.17 Closure Summary

Phase 0.4.17 turns the enriched wallet inventory read model into a small but explicitly queryable authenticated API contract.

The backend now supports optional query semantics on `GET /auth/wallets` for:

- `status`
- `primary`
- `linked_at` ordering

### Delivered in 0.4.17

- optional `status=active|detached` filter
- optional `primary=true|false` filter
- optional `sort=linked_at` with `order=asc|desc`
- strict `400` validation for unsupported query values
- handler-level test coverage for filtering, ordering, and invalid query parameters
- documentation updates aligning the authenticated wallet inventory contract with the new query semantics

---

## 🔍 Functional Result

The system now supports the following authenticated inventory behavior:

1. authenticated user calls `GET /auth/wallets` with or without optional query params
2. backend resolves wallet identities currently owned by that durable user
3. backend maps each identity into the existing lifecycle-aware wallet read model
4. optional filters are applied to that read model
5. optional `linked_at` ordering is applied only when explicitly requested
6. invalid query semantics are rejected with `400` instead of being silently accepted

---

## ❌ Not Included in 0.4.17

The following items remain intentionally out of scope:

- pagination
- text search
- detached-wallet history endpoints
- admin inventory views
- store-level filtering or ordering changes
- ownership-rule changes
- new wallet mutation endpoints

---

## ⏭️ Next Phase

### Next Expected Evolution

- wallet inventory pagination only if a real client need appears
- additional low-risk query semantics only if they remain read-only and backward compatible
- preserve the current ownership and lifecycle invariants already stabilized in Phase 0.4



## ✅ Phase 0.4.18 Closure Summary

Phase 0.4.18 adds simple windowing semantics to the authenticated wallet inventory API contract.

The backend now supports optional pagination on `GET /auth/wallets` through:

- `limit`
- `offset`
- additive response metadata (`total`, `limit`, `offset`)

### Delivered in 0.4.18

- optional `limit=<positive integer>`
- optional `offset=<non-negative integer>`
- strict `400` validation for malformed pagination values
- pagination applied only after wallet inventory filtering and sorting
- response metadata exposing filtered total and requested window parameters
- handler-level test coverage for valid and invalid pagination behavior

---

## 🔍 Functional Result

The system now supports the following authenticated inventory behavior:

1. authenticated user calls `GET /auth/wallets` with optional filters, sorting, and pagination
2. backend resolves wallet identities currently owned by that durable user
3. backend maps each identity into the lifecycle-aware wallet read model
4. optional filters are applied
5. optional ordering is applied
6. optional pagination window is applied
7. backend returns `wallets`, `total`, `limit`, and `offset`

---

## ❌ Not Included in 0.4.18

The following items remain intentionally out of scope:

- cursor pagination
- next-page tokens
- text search
- detached-wallet history endpoints
- admin inventory views
- store-level pagination or query expansion
- ownership-rule changes
- new wallet mutation endpoints

---

## ⏭️ Next Phase

### Next Expected Evolution

- only add further wallet inventory query semantics if a concrete client need appears
- preserve backward compatibility of the paginated wallet inventory contract
- keep future inventory enhancements read-only unless the ZIP proves otherwise


## ✅ Phase 0.4.19 Closure Summary

Phase 0.4.19 completes the wallet inventory response contract with additive navigation metadata.

The backend now exposes the following authenticated inventory metadata on `GET /auth/wallets`:

- `total`
- `limit`
- `offset`
- `returned`
- `has_more`

### Delivered in 0.4.19

- additive `returned` field describing the current window size
- additive `has_more` field describing whether a next page exists
- deterministic calculation after filtering, sorting, and pagination
- preserved backward compatibility of the wallet inventory contract
- handler-level coverage for default, paginated, empty-window, and filtered-window navigation scenarios

---

## ❌ Not Included in 0.4.19

The following items remain intentionally out of scope:

- cursor pagination
- next-page tokens
- `next_offset` / `previous_offset`
- text search
- detached-wallet history endpoints
- store-level query expansion
- ownership-rule changes

---

## ⏭️ Next Phase

### 0.4.23 — Wallet Inventory Query Examples Closure

Expected next focus:

- preserve backward compatibility of the clarified wallet inventory response contract
- only extend inventory semantics when a concrete client need appears
- keep future inventory work read-only unless the ZIP proves otherwise


## ✅ Phase 0.4.21 Closure Summary

Phase 0.4.21 hardens the wallet inventory query-parameter contract without adding new inventory features or touching ownership persistence.

### Delivered in 0.4.21

- `order` now requires an explicit `sort`
- `sort=linked_at` now defaults explicitly to ascending order when `order` is omitted
- offset-only requests remain valid and unbounded
- handler-level tests cover the hardened contract combinations and defaults

## ❌ Not Included in 0.4.21

- new filters
- new sort fields
- cursor pagination
- continuation tokens
- store-level pagination
- ownership-rule changes

### 0.4.22 — Wallet Inventory Response Contract Clarification

Expected next focus:

- preserve backward compatibility of the wallet inventory response contract
- only extend inventory semantics when a concrete client need appears
- keep future inventory work read-only unless the ZIP proves otherwise


## ✅ Phase 0.4.22 Closure Summary

Phase 0.4.22 clarifies the wallet inventory response contract so the operator-facing endpoint documentation matches the JSON behavior already implemented in prior subphases.

### Delivered in 0.4.22

- the main `GET /auth/wallets` README response example now includes `returned` and `has_more`
- response field semantics are explicitly documented for bounded and unbounded inventory requests
- navigation hints (`next_offset`, `previous_offset`) are documented as bounded-window metadata
- phase and handoff documentation now reflect that the response contract is explicitly clarified

## ❌ Not Included in 0.4.22

- new endpoint behavior
- new filters
- new sort fields
- cursor pagination
- store-level pagination
- ownership-rule changes

## ✅ Phase 0.4.24 Closure Summary

Phase 0.4.24 closes the manual-validation layer for `GET /auth/wallets` so operators have an explicit checklist for validating the already-implemented query contract end-to-end.

### Delivered in 0.4.24

- consolidated manual validation scenarios for base, filtered, sorted, and paginated wallet inventory requests
- explicit manual checks for bounded vs unbounded window behavior
- explicit manual checks for `returned`, `has_more`, `next_offset`, and `previous_offset`
- invalid-query manual checks for contractual errors such as `order` without `sort`

## ❌ Not Included in 0.4.24

- new endpoint behavior
- new filters
- new sort fields
- cursor pagination
- store-level pagination
- ownership-rule changes

## ✅ Phase 0.4.25 Closure Summary

Phase 0.4.25 prepares the authenticated wallet inventory for wallet-management consumption by exposing minimal actionability hints per listed wallet without changing stores, persistence, or execution authority.

### Delivered in 0.4.25

- additive wallet inventory fields: `can_set_primary`, `can_detach`, and `detach_block_reasons`
- detach block reasons aligned with the existing detach-domain reasons (`wallet_is_primary`, `user_would_have_no_wallets`)
- handler-level validation for single-wallet and two-wallet inventory scenarios
- inventory-side actionability semantics kept explicitly advisory, with execution authority left to the existing action endpoints

## ❌ Not Included in 0.4.25

- new wallet-management endpoints
- changes to detach or primary-switch execution behavior
- new query parameters
- store-level actionability persistence
- ownership-rule changes

## ✅ Phase 0.4.26 Closure Summary

Phase 0.4.26 closes the consistency gap between wallet inventory actionability hints and `POST /auth/wallets/detach/check` without changing detach rules, stores, or persistence.

### Delivered in 0.4.26

- handler-level consistency coverage for single-primary and two-wallet inventories
- explicit validation that `can_detach=false` remains compatible with `eligible=false` under the same detach reasons
- explicit validation that detachable secondary wallets stay aligned with `eligible=true` and empty detach reasons
- documentation that keeps inventory-side hints advisory while preserving detach-check authority

## ❌ Not Included in 0.4.26

- new detach rules
- new wallet-management endpoints
- changes to detach execution behavior
- new inventory query parameters
- ownership-rule changes

## ✅ Phase 0.4.27 Closure Summary

Phase 0.4.27 closes the consistency gap between wallet inventory primary-actionability hints and `POST /auth/wallets/primary` without changing primary-switch rules, stores, or persistence.

### Delivered in 0.4.27

- handler-level consistency coverage for a two-wallet inventory before and after primary switching
- explicit validation that the current primary stays `can_set_primary=false` before the switch
- explicit validation that a secondary wallet exposed as `can_set_primary=true` can be promoted and then becomes non-promotable after the switch
- documentation that keeps inventory-side primary hints advisory while preserving primary-switch authority

## ❌ Not Included in 0.4.27

- new primary-switch rules
- new wallet-management endpoints
- new inventory fields
- ownership-rule changes
- store-level actionability persistence

## ✅ Phase 0.4.28 Closure Summary

Phase 0.4.28 closes the wallet-management read flow around the authenticated inventory and the existing primary / detach actions without changing domain rules, stores, or persistence.

### Delivered in 0.4.28

- main README header corrected so the declared current subphase matches the actual state already reflected across the ZIP
- explicit documentation of the end-to-end wallet-management flow: inventory → actionability hint → action/check endpoint → refreshed inventory
- operator guidance clarifying that inventory hints remain advisory while action and check endpoints remain authoritative
- manual validation guidance covering refreshed inventory expectations after primary switching and detach execution

## ❌ Not Included in 0.4.28

- new wallet-management endpoints
- new inventory fields
- changes to primary-switch or detach rules
- ownership-rule changes
- store-level or persistence changes

## ✅ Phase 0.4.29 Closure Summary

Phase 0.4.29 closes the consistency gap between authenticated wallet inventory detach hints and `POST /auth/wallets/detach` without changing domain rules, stores, or persistence.

### Delivered in 0.4.29

- handler-level coverage proving that a secondary wallet exposed as detachable can be detached successfully
- explicit validation that the detach execute response stays compatible with the pre-detach inventory hints and eligibility snapshot
- explicit validation that refreshed inventory removes the detached wallet from the attached inventory and recalculates detach hints coherently for the remaining wallet
- documentation clarifying that inventory-side detach hints remain advisory while detach execution remains authoritative

## ❌ Not Included in 0.4.29

- new detach rules
- new wallet-management endpoints
- new inventory fields
- ownership-rule changes
- store-level or persistence changes

## ✅ Phase 0.4.30 Closure Summary

Phase 0.4.30 consolidates the authenticated wallet-management surfaces into one explicit contract without changing handlers, stores, persistence, or domain rules.

### Delivered in 0.4.30

- consolidated wallet-management contract across inventory, primary switch, detach eligibility, and detach execution
- explicit documentation of advisory versus authoritative wallet-management surfaces
- unified operator/testing guidance for the inventory → action/check → refreshed inventory cycle
- cross-document alignment so handoff, flows, README, and testing describe the same final wallet-management model

## ❌ Not Included in 0.4.30

- new wallet-management endpoints
- new inventory fields
- changes to detach or primary rules
- store-level or persistence changes

## ✅ Phase 0.4.31 Closure Summary

Phase 0.4.31 closes the remaining challenge-purpose enforcement gap at the wallet-auth bootstrap boundary without changing ownership semantics, stores, or persistence.

### Delivered in 0.4.31

- service-level enforcement that `POST /auth/wallet/verify` accepts only `auth_bootstrap` challenges
- explicit rejection of `wallet_link` challenges in wallet-auth bootstrap
- explicit rejection of `account_merge` challenges in wallet-auth bootstrap
- handler-level `wallet_challenge_purpose_mismatch` response for purpose violations
- test coverage proving non-bootstrap challenge purposes cannot be reused in wallet login

## ❌ Not Included in 0.4.31

- new wallet-management endpoints
- ownership-rule changes
- primary-wallet changes
- detach-rule changes
- store-level or persistence changes

### 0.4.32 — To Be Defined Against Real ZIP

Expected next focus:

- only continue Phase 0.4 if the next ZIP shows a concrete remaining lifecycle gap
- preserve strict challenge-purpose isolation across wallet login and wallet-management flows
- keep detached-wallet reuse semantics backward-compatible and ownership-safe


## ✅ Phase 0.4.23 Closure Summary

Phase 0.4.23 closes the concrete examples layer for `GET /auth/wallets` so operators and client implementers can see valid and invalid request patterns alongside bounded-window response examples.

### Delivered in 0.4.23

- request examples for base, filtered, sorted, and paginated wallet inventory queries
- an explicit invalid example for `order` without `sort`
- response examples showing bounded-window metadata and navigation hints
- accumulated documentation aligned with the real handler contract

## ❌ Not Included in 0.4.23

- new endpoint behavior
- new filters
- new sort fields
- cursor pagination
- store-level pagination
- ownership-rule changes

