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

### Phase 0.5 — User Interaction & Application Surface
Status: ✅ Completed

### Phase 0.6 — Authenticated Application Bootstrap Consolidation & Session-Ready Surface
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
| 0.4.32 | Wallet challenge purpose strictness closure | ✅ Completed |

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

## ✅ Phase 0.4.32 Closure Summary

Phase 0.4.32 closes the last permissive purpose-normalization gap in wallet challenges without changing stores, persistence, ownership, or lifecycle semantics.

### Delivered in 0.4.32

- strict creation-time purpose resolution with controlled defaulting only for empty purpose
- unknown purpose values are no longer normalized to `auth_bootstrap` at runtime
- wallet verify/login rejects unknown challenge purposes
- authenticated wallet link rejects unknown challenge purposes
- authenticated wallet-owned account merge rejects unknown challenge purposes
- tests proving invalid purpose values are preserved and rejected instead of silently reclassified

## ❌ Not Included in 0.4.32

- new wallet lifecycle operations
- ownership-rule changes
- primary-wallet changes
- detach-rule changes
- store schema or migration changes

## ✅ Phase 0.4.33 Closure Summary

Phase 0.4.33 formally closes Phase 0.4 at the documentation layer without changing handlers, stores, persistence, migrations, ownership rules, or wallet lifecycle behavior.

### Delivered in 0.4.33

- explicit formal closure of Phase 0.4 as a completed foundation phase
- cross-document alignment so README, phase status, handoff, flows, and phase documentation describe the same final state
- removal of the remaining placeholder planning note that left 0.4.33 undefined in phase-status tracking
- explicit transition guidance that the next work should start in a new phase instead of reopening Phase 0.4 without a real ZIP-validated gap

## ❌ Not Included in 0.4.33

- runtime handler changes
- new wallet-management endpoints
- ownership-rule changes
- primary-wallet changes
- detach-rule changes
- challenge-purpose behavior changes
- store-level or persistence changes
- schema or migration changes

### Next Expected Evolution

- start the next phase from a ZIP-validated runtime or product need outside the already closed Phase 0.4 scope
- preserve the finalized wallet-auth, wallet-link, wallet-merge, primary-switch, and detach contracts as the baseline
- avoid reopening Phase 0.4 unless a future ZIP proves a real regression or contractual documentation gap


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



## Phase 0.5 Subphase Status

| Subphase | Description | Status |
|----------|-------------|--------|
| 0.5.1 | Authenticated user profile surface | ✅ Completed |
| 0.5.2 | User metadata (non-wallet) | ✅ Completed |
| 0.5.3 | User settings contract foundation | ✅ Completed |
| 0.5.4 | User settings mutation | ✅ Completed |
| 0.5.5.1 | User settings hardening (soft contract guard layer) | ✅ Completed |
| 0.5.5.2 | User settings deep merge preservation and known-branch semantics | ✅ Completed |

---

## ✅ Phase 0.5.3 Closure Summary

Phase 0.5.3 introduces the first dedicated authenticated user settings contract without reopening identity, wallet ownership, or the minimal metadata editing semantics already established in 0.5.1 and 0.5.2.

### Delivered in 0.5.3

- authenticated `GET /auth/me/settings`
- dedicated `user_settings` persistence foundation separated from `users`
- minimal settings response contract containing:
  - `user_id`
  - `version`
  - `preferences`
- safe default resolution when no persisted settings row exists
- no implicit write or row creation during read
- dedicated `usersettings` module and service/repository foundation
- test coverage for:
  - default settings read
  - persisted settings read
  - unauthorized access
  - unavailable settings service

### Result

The backend now exposes three distinct authenticated user-facing surfaces:

- `GET /auth/me` for profile/bootstrap data
- `PATCH /auth/me` for minimal non-wallet metadata update
- `GET /auth/me/settings` for dedicated settings bootstrap

This preserves a clean separation between:

- durable user identity
- minimal profile metadata
- authenticated user settings

## ❌ Not Included in 0.5.3

- settings mutation
- typed preference fields such as theme, locale, or notifications
- settings audit/history
- wallet lifecycle interaction
- identity-model changes
- merging settings into `/auth/me`

### Next Expected Evolution

- add settings mutation only through a dedicated settings contract
- keep profile/bootstrap data and settings separated
- introduce typed settings fields only when a real application need requires them

---

## ✅ Phase 0.5.4 Closure Summary

Phase 0.5.4 introduces the first authenticated user settings mutation surface through `PATCH /auth/me/settings`, preserving the separation between profile/bootstrap metadata and dedicated user settings.

### Delivered in 0.5.4

- authenticated `PATCH /auth/me/settings`
- merge-based mutation for `preferences`
- partial updates without destructive overwrite of the entire settings object
- safe persistence through `user_settings`
- minimal validation for malformed settings payloads
- handler-level and service-level test coverage for mutation behavior

### Result

The backend now exposes a dedicated authenticated read/write settings surface while preserving the same separation principles established in 0.5.3:

- `GET /auth/me` for profile/bootstrap data
- `PATCH /auth/me` for minimal non-wallet metadata update
- `GET /auth/me/settings` for dedicated settings read
- `PATCH /auth/me/settings` for dedicated settings mutation

This keeps settings evolution out of `/auth/me` and away from the core `users` record.

## ❌ Not Included in 0.5.4

- strict schema enforcement
- typed preference catalog
- settings version negotiation
- audit/history for settings mutation
- wallet lifecycle interaction
- identity-model changes
- merging settings into `/auth/me`

### Next Expected Evolution

- introduce lightweight hardening to prevent structural drift in `preferences`
- preserve backward compatibility of the settings mutation contract
- keep settings flexible while improving consistency guarantees

---

## ✅ Phase 0.5.5.1 Closure Summary

Phase 0.5.5.1 introduces the first hardening layer over the authenticated user settings contract without converting the settings surface into a schema-heavy or rigid system.

### Delivered in 0.5.5.1

- recursive normalization of `preferences`
- rejection of `null` values, including nested values
- rejection of non JSON-compatible values
- top-level key trimming
- top-level shape compatibility protection across object / array / scalar categories
- soft observation of unknown top-level keys through warning logs
- expanded service-level and HTTP-level test coverage for the hardened contract

### Result

The backend now protects the authenticated settings contract from structural drift while preserving:

- backward compatibility
- merge-based mutation semantics
- flexibility of `preferences`
- separation between user profile and user settings

This means the settings surface remains extensible, but no longer accepts obviously unsafe or structurally incompatible mutations as freely as before.

## ❌ Not Included in 0.5.5.1

- strict JSON schema enforcement
- hard key whitelisting
- typed preference governance
- settings migration/versioning
- semantic validation of specific preference families
- changes to auth, wallet lifecycle, or durable identity behavior

### Next Expected Evolution

- introduce typed settings only when a real application need justifies them
- preserve the flexible authenticated settings contract while hardening incrementally
- continue Phase 0.5 application-surface evolution without reopening Phase 0.4



## ✅ Phase 0.5.5.2 Closure Summary

Phase 0.5.5.2 closes the main destructive-partial-update gap that still remained in authenticated user settings after 0.5.5.1.

The backend now preserves nested object branches during partial settings mutation and introduces minimal structural semantics for the already-known top-level settings families without turning the settings surface into a rigid schema-driven subsystem.

### Delivered in 0.5.5.2

- deep merge preservation for nested object branches inside `preferences`
- sibling-key preservation during partial updates of nested settings families
- recursive compatibility enforcement during nested merge operations
- minimal known-branch semantics requiring `notifications`, `preferences`, and `ui` top-level values to remain objects
- expanded service-level and HTTP-level test coverage for nested merge preservation and known-branch rejection

### Result

The authenticated settings contract now better matches real partial-update expectations:

- nested object patches no longer destructively replace entire persisted branches
- known settings namespaces keep object semantics instead of accepting arbitrary scalars or arrays
- the backend remains flexible for unknown keys while providing stronger guarantees for the families it already recognizes

This keeps the application surface incremental and backward compatible while making `PATCH /auth/me/settings` substantially safer for frontend-driven partial mutations.

## ❌ Not Included in 0.5.5.2

- strict schema enforcement for all preference families
- per-key semantic validation such as theme catalogs or locale normalization
- hard rejection of unknown top-level keys
- settings migration/versioning
- audit/history for settings mutation
- changes to auth, wallet lifecycle, or durable identity behavior

### Next Expected Evolution

- continue contract stabilization only where real application needs justify it
- introduce typed settings selectively instead of freezing the whole settings tree prematurely
- preserve Phase 0.5 application-surface scope without reopening Phase 0.4

---


## ✅ Phase 0.5.5.3 Closure Summary

Phase 0.5.5.3 closes the next contract-surface gap in authenticated user settings by exposing persisted resource timestamps through the stable `settings` view returned by both `GET /auth/me/settings` and `PATCH /auth/me/settings`.

### Delivered in 0.5.5.3

- `usersettings.View` now exposes `created_at` when persistence metadata exists
- `usersettings.View` now exposes `updated_at` when persistence metadata exists
- settings responses remain backward compatible with the existing envelope and `preferences` payload
- zero-value timestamps are omitted so default, non-persisted settings resolution does not fabricate persistence state
- HTTP-level coverage now validates timestamp exposure and omission semantics for the authenticated settings surface


---

## Phase 0.6 Subphase Status

| Subphase | Description | Status |
|----------|-------------|--------|
| 0.6.1 | Bootstrap Surface Boundary Clarification | ✅ Completed |
| 0.6.2 | Authenticated Surface Contract Alignment | ✅ Completed |
| 0.6.3 | Session-Ready Bootstrap Read Model | ⬜ Pending |
| 0.6.4 | Application Surface Consistency Hardening | ⬜ Pending |

### Phase 0.6 Goal

Phase 0.6 continues Stage 0 without introducing new business domains. Its purpose is to consolidate the already authenticated application surface into a coherent bootstrap layer for frontend consumption by clarifying the boundary between identity, session, settings, and wallet inventory reads.

### Included in Phase 0.6

- consolidation of the authenticated bootstrap surface already present in the backend
- alignment of `/auth/me`, `/auth/session`, `/auth/me/settings`, and `/auth/wallets`
- semantic and contractual clarification across related authenticated reads
- consistency hardening backed by focused tests and documentation

### Excluded from Phase 0.6

- new business endpoints
- payments, billing, trading, or domain expansion
- wallet lifecycle redesign
- settings migrations or advanced typed-settings redesign
- changes to the durable auth model already stabilized in Phase 0.4 and Phase 0.5

## ✅ Phase 0.6.1 Closure Summary

Phase 0.6.1 establishes an explicit semantic boundary across the authenticated application bootstrap surface without changing existing public payloads. The subphase locks in the intended separation of responsibilities between session context, authenticated identity bootstrap, settings retrieval, and wallet inventory reads.

### Delivered in 0.6.1

- explicit boundary coverage now distinguishes `/auth/me` from `/auth/session`
- `/auth/me` is now contractually protected as the authenticated bootstrap identity surface
- `/auth/session` is now contractually protected as the token/session context surface
- the authenticated settings surface remains anchored at `/auth/me/settings`
- wallet inventory remains anchored at `/auth/wallets`
- no handler, router, service, persistence, or payload contract changes were required
- focused HTTP tests now guard against accidental surface overlap between bootstrap identity data and token-derived session metadata

### Result

The authenticated application surface is now clearer for downstream consumers:

- `/auth/me` represents authenticated identity/bootstrap data
- `/auth/session` represents token-derived session state
- `/auth/me/settings` represents authenticated settings retrieval and mutation
- `/auth/wallets` represents wallet ownership inventory

This gives Phase 0.6 a stable starting point for the next alignment step, where response shapes and metadata consistency can be refined without ambiguity about endpoint responsibilities.

## ❌ Not Included in 0.6.1

- response shape normalization across authenticated endpoints
- metadata unification between identity and session payloads
- a new aggregate bootstrap endpoint
- frontend-specific orchestration behavior
- changes to wallet lifecycle semantics
- settings schema redesign
- business-surface expansion beyond the authenticated foundation

### Next Expected Evolution

- align related authenticated response shapes where doing so improves bootstrap consumption without breaking compatibility
- preserve the boundary clarified in 0.6.1 while improving cross-endpoint coherence in 0.6.2
- continue treating authenticated bootstrap consolidation as a Stage 0 concern rather than a new domain feature


---

## ✅ Phase 0.6.2 Closure Summary

Phase 0.6.2 aligns the authenticated bootstrap surface without breaking existing public payloads by introducing a shared internal normalization path for authenticated context and by expanding cross-endpoint contract coverage across identity, session, and wallet reads.

### Delivered in 0.6.2

- a shared authenticated-context helper now normalizes `user_id`, `email`, `wallet_id`, `wallet_address`, `auth_method`, `chain`, and `has_wallet_session`
- `/auth/me` now derives its shared authenticated context from the same logical source used by `/auth/session`
- wallet-context normalization is now explicit, so missing `wallet_address` also clears partial wallet context such as `wallet_id` and `chain`
- cross-endpoint contract coverage now validates consistency between `/auth/me` and `/auth/session` for shared authenticated fields
- bootstrap coherence coverage now validates that top-level `user` and nested `profile.user` remain aligned inside `/auth/me`
- wallet-surface coherence coverage now validates that the primary wallet exposed by `/auth/me` remains aligned with the primary wallet returned by `/auth/wallets`
- the subphase preserves backward compatibility by keeping existing JSON envelopes and field names unchanged

### Result

The authenticated bootstrap surface now remains semantically separated while also becoming more contractually aligned for downstream consumers:

- `/auth/me` and `/auth/session` share a normalized authenticated context instead of relying on separate drift-prone derivations
- partial wallet claim state is now less likely to diverge across bootstrap-related reads
- `/auth/me` and `/auth/wallets` now have stronger tested coherence around primary wallet representation
- consumers can rely on tighter internal consistency without requiring any migration of existing response contracts

This positions Phase 0.6 for the next step, where the bootstrap read model can be consolidated more explicitly without reopening boundary semantics or introducing contract instability.

## ❌ Not Included in 0.6.2

- removal of duplicated identity data kept for compatibility inside `/auth/me`
- JSON field renames or envelope redesign across authenticated endpoints
- creation of a new aggregate bootstrap endpoint
- changes to `/auth/me/settings` contract shape or settings persistence semantics
- wallet lifecycle redesign, pagination redesign, or business-domain expansion

### Next Expected Evolution

- consolidate a clearer session-ready bootstrap read model on top of the now-aligned authenticated context
- preserve compatibility while reducing frontend ambiguity around which endpoint should be treated as the primary bootstrap source for each authenticated concern
- continue hardening the authenticated application surface as a Stage 0 foundation concern before new business domains are introduced

### Phase 0.6 — Subphase Status (Updated)

| Subphase | Description | Status |
|----------|------------|--------|
| 0.6.1 | Bootstrap Surface Boundary Clarification | ✔ Completed |
| 0.6.2 | Authenticated Surface Contract Alignment | ✔ Completed |
| 0.6.3 | Session-Ready Bootstrap Read Model | ✔ Completed |
| 0.6.4 | Application Surface Consistency Hardening | ⬜ Pending |

---

## ✅ Phase 0.6.3 Closure Summary

### Delivered in 0.6.3

- Introduced `GET /auth/bootstrap`
- Aggregated authenticated surfaces into a single read model
- Reused existing builders and services (session, profile, settings, wallets)
- Eliminated multi-request bootstrap requirement for frontend

### Result

The authenticated surface is now:

- boundary-defined (0.6.1)
- contract-aligned (0.6.2)
- bootstrap-ready (0.6.3)

### Next Expected Evolution

0.6.4 — Application Surface Consistency Hardening

### Phase 0.6 — Subphase Status (Final)

| Subphase | Description | Status |
|----------|------------|--------|
| 0.6.1 | Bootstrap Surface Boundary Clarification | ✔ Completed |
| 0.6.2 | Authenticated Surface Contract Alignment | ✔ Completed |
| 0.6.3 | Session-Ready Bootstrap Read Model | ✔ Completed |
| 0.6.4 | Application Surface Consistency Hardening | ✔ Completed |

---

## ✅ Phase 0.6 Closure Summary

### Delivered in Phase 0.6

- Explicit authenticated surface boundaries
- Shared authenticated context normalization
- Session-ready bootstrap read model
- Canonical wallet envelope and structural hardening
- Cross-endpoint contract consistency guarantees

### Final Result

The authenticated application surface is now:

- fully defined
- internally consistent
- contractually aligned
- bootstrap-ready for frontend consumption

### Next Phase

Phase 0.7 — Application Layer Foundation

## 🔄 Phase 0.7 — Application Layer Foundation

### Subphase Status

| Subphase | Description | Status |
|----------|------------|--------|
| 0.7.1 | Application Layer Boundary Definition | ✔ Completed |
| 0.7.2 | Authenticated Surface Use Cases Extraction | ✔ Completed |
| 0.7.3 | Wallet Management Use Cases Consolidation | ✔ Completed |
| 0.7.4 | Handler Simplification & Contract Preservation | ✔ Completed |

---

## ✅ Phase 0.7.1 Summary

### Delivered

- Introduced explicit application layer (`Application`)
- Established handler → application → services boundary
- Migrated bootstrap orchestration out of HTTP layer
- Added application-level test coverage

### Result

The backend now has a formal entry point for use-case orchestration, enabling controlled evolution of business logic without overloading HTTP handlers.

### Next

0.7.2 — Authenticated Surface Use Cases Extraction

### 0.7.2 — Authenticated Surface Use Cases Extraction ✔ Completed

This subphase completes the extraction of the authenticated surface core use cases into the application layer.

#### Delivered

- Login use-case moved to application layer
- GetMe use-case moved to application layer
- GetSession use-case moved to application layer
- Handlers reduced to transport-only responsibilities

#### Result

Authenticated surface is now fully application-driven:

- 0.7.1 → boundary defined
- 0.7.2 → core use cases extracted

#### Next

0.7.3 — Wallet Management Use Cases Consolidation

### 0.7.3 — Wallet Management Use Cases Consolidation ✔ Completed

This subphase completes the migration of wallet management flows into the application layer.

#### Delivered

- Wallet listing moved to application layer
- Wallet link (challenge/verify) moved to application layer
- Wallet account merge (challenge/verify) moved to application layer
- Set primary wallet moved to application layer
- Wallet detach (check/execute) moved to application layer
- Wallet HTTP handlers reduced to transport-only responsibilities

#### Result

Auth module is now fully application-driven:

- 0.7.1 → boundary defined
- 0.7.2 → authenticated surface extracted
- 0.7.3 → wallet management consolidated

#### Next

0.7.4 — Handler Simplification & Contract Preservation

### 0.7.4 — Handler Simplification & Contract Preservation ✔ Completed

This subphase finalizes the application-layer foundation by simplifying all authenticated handlers and removing residual transport duplication.

#### Delivered

- Unified request decoding across handlers
- Centralized authenticated claims extraction
- Standardized error JSON writing
- Simplified all auth and wallet handlers to transport-only logic
- Removed duplicated transport logic patterns

#### Result

Phase 0.7 is now fully completed:

- 0.7.1 → boundary defined
- 0.7.2 → authenticated surface extracted
- 0.7.3 → wallet management consolidated
- 0.7.4 → handlers simplified and contracts preserved

Auth module is now fully application-driven with clean transport boundaries.

#### Next

0.8 — Standardized Error Model

## ✅ Phase 0.8 — Standardized Error Model

### Subphase Status

| Subphase | Description | Status |
|----------|------------|--------|
| 0.8.1 | Error Contract Definition | ✔ Completed |
| 0.8.2 | Error Type System Introduction | ✔ Completed |
| 0.8.3 | Auth Surface Error Standardization | ✔ Completed |
| 0.8.4 | Error Mapping Hardening & Contract Tests | ✔ Completed |

---

## ✅ Phase 0.8 Summary

### Delivered

- Introduced a centralized response error contract under `internal/core/errs`
- Added reusable HTTP error envelope writing in `internal/core/httpx`
- Migrated auth handlers and auth middleware to the new error envelope
- Preserved existing handler/domain error decision points while changing only the response contract structure
- Updated auth HTTP tests to assert the new envelope shape and normalized codes
- Added a centralized internal `AppError` type with categories, wrapping helpers and conversion to response errors
- Moved the auth-local normalized error catalog into `internal/core/errs` to prepare reusable mapping across later phases
- Added shared typed factories and HTTP writing support for application errors
- Standardized auth-surface handlers on centralized app-error factories without reintroducing cyclic imports
- Added contract-focused tests under `internal/core/errs` and `internal/core/httpx` to freeze representative mapping behavior and canonical HTTP error envelope serialization

### Result

Phase 0.8 is now fully completed:

- 0.8.1 → error response contract introduced
- 0.8.2 → internal app error type system introduced
- 0.8.3 → auth surface standardized on centralized app errors
- 0.8.4 → mapping hardened and contract tests added

The backend now exposes a single structured JSON error envelope across auth handlers and core auth HTTP middleware, backed by centralized app-error mapping and dedicated contract-level hardening tests.

## ✅ Phase 0.9 — API Versioning Strategy

### Subphase Status

| Subphase | Description | Status |
|----------|------------|--------|
| 0.9.1 | Versioning Policy Definition | ✔ Completed |
| 0.9.2 | Router Versioning Foundation | ✔ Completed |
| 0.9.3 | Authenticated Surface Version Freezing | ✔ Completed |
| 0.9.4 | Version-aware Contract Testing | ✔ Completed |
| 0.9.5 | Documentation Consolidation | ✔ Completed |

---

## ✅ Phase 0.9 Summary

### Objective

Define the canonical API versioning policy for the post-0.8 backend while preserving the currently consumed legacy transport surface and avoiding forced frontend movement before Stage 0 closes.

### Delivered

- Declared path-based versioning as the canonical public API strategy
- Reserved `/api/v1/...` as the canonical route namespace for current authenticated and auth-adjacent surfaces
- Classified the existing unversioned endpoints as legacy, backward-compatible and non-canonical rather than obsolete or immediately removable
- Froze `v1` as the semantic home of the current success-payload and standardized error-envelope behavior
- Recorded the rule that breaking transport or error-contract changes require a new API version rather than silent mutation of the existing surface
- Linked the versioning step explicitly to the post-0.8 foundation state and to the upcoming authorization layer in 0.10
- Preserved the project constraint that the frontend remains aligned to backend Phase 0.6 until Stage 0 completion, so versioning preparation must not force immediate frontend adoption
- Materialized canonical `/api/v1/...` route exposure in the real router while preserving legacy `/auth/...` access
- Reused the same handler/application behavior across legacy and canonical route spaces instead of creating route-specific business duplication
- Kept current auth middleware and transport concerns intact across both route surfaces
- Explicitly froze the current authenticated surface as canonical `v1` semantics across bootstrap, profile, settings, session and wallet inventory/read flows
- Declared legacy and canonical authenticated entry paths to be two projections of the same `v1` contract rather than independent route contracts

- Added representative router-level contract tests covering legacy/canonical path composition and success/error parity for key authenticated/auth-adjacent routes
- Consolidated roadmap, phase-status, README, architecture, testing and handoff documents so Phase 0.9 now closes with a single coherent description of canonical `v1`, legacy compatibility and the frontend Phase 0.6 alignment rule

### Current Result

Phase 0.9 is now fully completed. The versioning foundation is both implemented in runtime code and consolidated across the trunk documentation set:

- 0.9.1 → versioning policy defined
- 0.9.2 → canonical router exposure completed
- 0.9.3 → authenticated-surface `v1` freeze completed
- 0.9.4 → version-aware contract testing completed
- 0.9.5 → trunk-documentation consolidation completed

The backend remains functionally stable on its current routes while now exposing a real canonical `v1` route surface backed by representative version-aware regression tests and by a documentation set aligned to the same transport-evolution narrative.

#### Next

0.10 — Authorization Layer


## ◑ Phase 0.10 — Authorization Layer

### Subphase Status

| Subphase | Description | Status |
|----------|------------|--------|
| 0.10.1 | Authorization Model Definition | ✔ Completed |
| 0.10.2 | Authorization Context & Middleware | ✔ Completed |
| 0.10.3 | Policy Evaluation Layer | ✔ Completed |
| 0.10.4 | Endpoint-Level Enforcement | ✔ Completed |
| 0.10.5 | Documentation & Contract Consolidation | ✔ Completed |

---

## ◑ Phase 0.10 Summary

### Objective

Introduce a structured authorization layer on top of the already-stabilized authentication, application and transport foundations without breaking current runtime behavior or frozen `v1` contracts.

### Delivered So Far

- Introduced a dedicated `internal/core/authorization` package
- Defined the initial role vocabulary with `user` and `admin`
- Defined the initial permission vocabulary for current Stage 0 authenticated surfaces
- Declared the foundational static role → permission mapping
- Introduced an `AuthorizationSubject` model separate from raw JWT claims
- Added normalization helpers so later context propagation and policy evaluation can consume stable data
- Added focused unit tests for mapping immutability, normalization and multi-role permission aggregation semantics

### Current Result

Phase 0.10 is now operating with progressive endpoint-level authorization enforcement. The backend no longer treats authorization as only preparatory infrastructure: selected authenticated endpoints already deny unauthorized access through the centralized policy layer while preserving the stabilized transport contract.

The current enforcement scope intentionally covers only the authenticated endpoints already aligned with the static role-permission model, leaving richer self-update semantics for later intentional expansion rather than premature policy broadening.

#### Next

0.11 — Domain Module Pattern


## ✅ Phase 0.10.2 — Authorization Context & Middleware

### Delivered

- introduced authorization-context helpers under `internal/core/authorization` to store and retrieve a normalized `AuthorizationSubject`
- added deterministic claim-to-subject resolution so authenticated identities can be projected into the authorization model without new persistence requirements
- introduced `HydrateAuthorization()` in `internal/core/httpx` as a dedicated middleware that runs after successful authentication
- integrated authorization hydration into the protected route stack for both legacy `/auth/...` routes and canonical `/api/v1/auth/...` routes
- added focused tests covering context storage, claim projection and middleware hydration semantics

### Result

Phase 0.10.2 is completed:

- 0.10.1 → authorization model vocabulary introduced
- 0.10.2 → authorization subject now propagates through the authenticated request lifecycle

The backend still does not enforce permissions yet, but authenticated requests now carry a first-class authorization subject in addition to authentication claims. This creates the infrastructure boundary needed for centralized policy evaluation in 0.10.3 without forcing handler-local authorization logic.


## ✅ Phase 0.10.3 — Policy Evaluation Layer

### Delivered

- introduced explicit authorization `Action` and `Resource` vocabularies under `internal/core/authorization`
- added centralized action/resource → permission resolution helpers so policy questions are derived from the same static authorization model introduced in 0.10.1
- introduced `PolicyEvaluator` with `Evaluate(...)` and `Can(...)` APIs for centralized authorization decisions
- added permission-check helpers so policy code evaluates normalized authorization subjects rather than raw transport claims
- added focused unit tests covering permission projection, positive and negative role-based decisions, and unknown action/resource denial behavior

### Result

Phase 0.10.3 is completed:

- 0.10.1 → authorization model vocabulary introduced
- 0.10.2 → authorization subject now propagates through the authenticated request lifecycle
- 0.10.3 → centralized policy evaluation is now available through a stable core API

The backend still does not enforce permissions yet, but it now has a dedicated policy boundary that can answer authorization questions without pushing role or permission logic into handlers. This sets up 0.10.4 to introduce endpoint-level enforcement progressively on top of an already-centralized evaluator.


## ✅ Phase 0.10.4 — Endpoint-Level Enforcement

### Delivered

- introduced route-level authorization enforcement middleware in `internal/core/httpx` backed by the centralized `PolicyEvaluator`
- added a standardized `AUTH_FORBIDDEN` error for denied authorization checks
- applied progressive enforcement to the authenticated endpoints already aligned with the current static permission map
- enforced `read user` on `GET /auth/me` and `GET /api/v1/auth/me`
- enforced `read settings` on `GET /auth/me/settings` and `GET /api/v1/auth/me/settings`
- enforced `update settings` on `PATCH /auth/me/settings` and `PATCH /api/v1/auth/me/settings`
- added focused middleware tests for allowed, denied and missing-subject authorization outcomes
- added router-level tests confirming legacy and canonical authenticated routes still succeed when the user holds the required permissions

### Result

Phase 0.10.4 is completed:

- 0.10.1 → authorization model vocabulary introduced
- 0.10.2 → authorization subject now propagates through the authenticated request lifecycle
- 0.10.3 → centralized policy evaluation is now available through a stable core API
- 0.10.4 → selected authenticated endpoints now enforce centralized authorization decisions

Authorization is no longer only preparatory infrastructure. The backend now actively denies unauthorized access on a selected subset of authenticated endpoints using the centralized policy layer, while preserving the stabilized transport contract and deliberately deferring endpoints whose semantics are not yet fully represented by the current static permission model.

## ✅ Phase 0.10.5 — Documentation & Contract Consolidation

### Delivered

- aligned Phase 0.10 completion state across README, roadmap, index, phase-status, handoff and the dedicated authorization phase document
- removed stale intermediate narrative that still described authorization as non-enforcing after endpoint-level enforcement had already been introduced in 0.10.4
- consolidated the public contract narrative so the repository now documents both the non-breaking rollout strategy and the new standardized authorization denial path
- recorded the intentionally progressive scope of the first authorization enforcement slice, including the deliberate exclusion of `PATCH /auth/me` from the current static permission map

### Result

Phase 0.10 is completed:

- 0.10.1 → authorization model vocabulary introduced
- 0.10.2 → authorization subject now propagates through the authenticated request lifecycle
- 0.10.3 → centralized policy evaluation is available through a stable core API
- 0.10.4 → selected authenticated endpoints enforce centralized authorization decisions
- 0.10.5 → repository documentation and contract state are consolidated around the delivered authorization layer

Authorization is now both implemented and consistently documented. The backend keeps its stabilized Stage 0 transport discipline while explicitly exposing the new authorization boundary, its current enforcement scope and the standardized forbidden path for covered endpoints.

#### Next

0.11 — Domain Module Pattern


## ✅ Phase 0.11 — Domain Module Pattern

Status: ✅ Completed

### Objective

Define and apply a consistent internal module pattern for the current Stage 0 domain-facing modules so the backend can continue growing on top of explicit transport, application and domain boundaries without changing the already stabilized public contract.

### Scope

Phase 0.11 is structural rather than functional. It is intended to:

- standardize the internal module layout around `http`, `app`, `domain` and optional `repository` boundaries
- reduce handler-local orchestration in the current Stage 0 modules
- make internal ownership clearer across `auth`, `user` and `usersettings`
- replace accidental concrete cross-module knowledge with minimal explicit contracts where those dependencies are real

This phase does not introduce new endpoints, payload shapes, authentication semantics or authorization behavior.

### Subphase Status

| Subphase | Description | Status |
|----------|-------------|--------|
| 0.11.1 | Domain Module Pattern Definition | ✅ Completed |
| 0.11.2 | User Module Refactor | ✅ Completed |
| 0.11.3 | UserSettings Module Refactor | ✅ Completed |
| 0.11.4 | Auth Module Alignment | ✅ Completed |
| 0.11.5 | Cross-Module Contract Consolidation | ✅ Completed |
| 0.11.6 | Documentation & Phase Closure | ✅ Completed |

### Architectural Intent

The target pattern for the current module set is:

```text
internal/modules/<module>/
    http/
    app/
    domain/
    repository/   (when it adds real clarity)
```

with the following responsibility split:

- `http` → request/response transport boundary
- `app` → use-case orchestration
- `domain` → models, invariants and internal contracts
- `repository` → persistence implementation boundary when needed

### Current Result After 0.11.1

0.11.1 is completed at the documentation and architectural-definition level. The repository now has an explicit Domain Module Pattern, explicit layer responsibility rules, explicit ownership boundaries between `auth`, `user` and `usersettings`, and explicit completion criteria for the rest of the phase.

### Current Result After 0.11.2

0.11.2 is completed in repository state.

The backend now includes the first concrete runtime application of the pattern under `internal/modules/user`, where the `user` module exposes explicit `app`, `domain` and `repository` boundaries while preserving the current external behavior and backward-compatible package access for existing consumers.

This subphase is validated by the passing repository test state, including the `internal/modules/user` package set.

### Current Result After 0.11.3

0.11.3 is completed in repository state.

The backend now also applies the Domain Module Pattern to `internal/modules/usersettings`, where the module exposes explicit `app`, `domain` and `repository` boundaries while preserving backward-compatible package access and current external settings behavior.

This subphase is validated by the passing repository test state, including the `internal/modules/usersettings` package set.

### Current Result After 0.11.4

0.11.4 is completed in repository state.

The `internal/modules/auth` package has been aligned conservatively to the Domain Module Pattern without changing public auth, wallet, bootstrap, session or authorization-adjacent behavior.

Delivered sub-steps within 0.11.4:

- 0.11.4A — boundary preparation
- 0.11.4B1 — support surface extraction
- 0.11.4B2a — app support decoupling
- 0.11.4B2b1 — app service and wallet runtime support decoupling
- 0.11.4B2b2.1 — wallet type ownership consolidation
- 0.11.4B2b2.2 — wallet read model consolidation
- 0.11.4B2b2.3 — wallet request/query compatibility alignment
- 0.11.4B2b2.4 — helper ownership consolidation cleanup
- 0.11.4B2b2.5 — application compatibility narrowing
- 0.11.4C1 — bootstrap/root DTO ownership reduction
- 0.11.4C2 — wallet response ownership reduction
- 0.11.4C3.1a — root service dependency stabilization
- 0.11.4C3.1b1 — session path narrowing
- 0.11.4C3.1b2 — login/current-user narrowing
- 0.11.4C3.2a — authenticated wallet management narrowing
- 0.11.4C3.2b1 — wallet bootstrap auth boundary clarification
- 0.11.4C4.1 — repository/runtime inventory
- 0.11.4C4.2 — repository/runtime decision lock
- 0.11.4D — final auth alignment validation
- 0.11.4E — documentation and closure

Resulting state:

- `internal/modules/auth/app` owns the runtime/application orchestration for auth, session, bootstrap and wallet-management use cases.
- `internal/modules/auth/domain` owns canonical wallet contracts and base wallet types.
- root `internal/modules/auth` remains a backward-compatible transport/runtime façade where needed.
- `internal/modules/auth/repository` is explicitly documented as a transitional façade, not the canonical implementation owner yet.
- root wallet store implementations remain accepted runtime implementations for stability.
- `http_bootstrap.go` and authenticated wallet-management handlers are narrowed to transport/delegation responsibilities.
- public payloads, routes, error semantics and test expectations remain stable.

Validated by a full `go test ./...` pass after the final C4.2/D state.


### Current Result After 0.11.5

0.11.5 is completed in repository state.

The backend now consolidates the contracts between the already-aligned `auth`, `user` and `usersettings` modules. The concrete result is that `auth/app` no longer depends on concrete `user.Service` or `usersettings.Service` implementations. Instead, it coordinates through explicit minimal contracts while root `auth` preserves backward-compatible wiring and public behavior.

Delivered sub-steps within 0.11.5:

- 0.11.5.0 — Subphase Definition & Documentation Alignment
- 0.11.5.1 — Dependency Mapping
- 0.11.5.2 — Contract Extraction
- 0.11.5.3 — Interface Alignment
- 0.11.5.4 — Runtime Compatibility Validation

0.11.5 did not introduce new public API behavior. Its purpose was to make existing cross-module coordination explicit and interface-based where coordination is real, while preserving the ownership boundaries completed in 0.11.1 through 0.11.4.

### Current Result After 0.11.6

0.11.6 is completed in repository state.

The trunk documentation now records Phase 0.11 as closed. Roadmap, phase status, architecture, deep architecture, handoff and README now agree that the Domain Module Pattern has been defined, applied to `user`, `usersettings` and `auth`, and finalized with cross-module contract consolidation.

### Expected Outcome

With Phase 0.11 completed, the repository now exposes:

- a consistent internal module pattern across `auth`, `user` and `usersettings`
- clearer ownership boundaries between authentication, user entity concerns and user settings concerns
- explicit cross-module contracts where coordination is required
- unchanged external API behavior compared with the completed 0.10 state

### Next

0.12 — Read/Write Model Separation

---

## 🚧 Phase 0.12 — Read / Write Model Separation

### Status

In progress.

### Latest Completed Subphase

0.12.4.0 — Mapping Layer Introduction Definition & Documentation Lock.

### Objective

Establish explicit separation between read models, write models and domain-owned models across the current Stage 0 modules, while preserving public API behavior and runtime compatibility.

### Scope

Included:

- complete model classification against the real repository
- explicit read model extraction
- explicit write model isolation
- mapping-layer introduction
- internal contract alignment after the 0.11 Domain Module Pattern
- cumulative documentation updates

Excluded:

- full CQRS
- event sourcing
- public API version changes
- business behavior changes
- multi-tenant redesign
- observability expansion

### Subphase Plan

| Subphase | Name | Status |
| --- | --- | --- |
| 0.12.0 | Phase Definition & Documentation Lock | ✅ Completed |
| 0.12.1 | Model Classification Audit | ✅ Completed |
| 0.12.2 | Read Model Extraction | ✅ Completed |
| 0.12.3 | Write Model Isolation | ✅ Completed |
| 0.12.4 | Mapping Layer Introduction | ✅ Completed |
| 0.12.5 | Contract Alignment | ✅ Completed |

### 0.12.0 Result

0.12.0 is documentation-only and is completed in this repository state. It defines the phase before implementation, records the subphase sequence and locks the non-goals that keep the phase compatibility-preserving.

No Go code is changed in 0.12.0.


### 0.12.1 Internal Sub-Subphase Plan

| Sub-subphase | Name | Status |
| --- | --- | --- |
| 0.12.1.0 | Definition & Documentation Lock | ✅ Completed |
| 0.12.1.1 | Model Inventory Extraction | ✅ Completed |
| 0.12.1.2 | Model Classification | ✅ Completed |
| 0.12.1.3 | Cross-Layer Usage Analysis | ✅ Completed |
| 0.12.1.4 | Problem Detection & Risk Mapping | ✅ Completed |
| 0.12.1.5 | Target Separation Definition | ✅ Completed |
| 0.12.1.6 | Audit Consolidation & Closure | ✅ Completed |

### 0.12.1.0 Result

0.12.1.0 is documentation-only and is completed in this repository state. It defines the Model Classification Audit methodology before inventory extraction starts. It also records the rule that when a subphase is subdivided, the `.0` step must lock trunk documentation before execution proceeds.

No Go code is changed in 0.12.1.0.

### Real Baseline For 0.12.1

The first implementation subphase must inspect the real models currently present under:

- `internal/modules/auth`
- `internal/modules/auth/app`
- `internal/modules/auth/domain`
- `internal/modules/user`
- `internal/modules/user/app`
- `internal/modules/user/domain`
- `internal/modules/usersettings`
- `internal/modules/usersettings/app`
- `internal/modules/usersettings/domain`
- relevant shared structures under `internal/core`

### 0.12.1.6 Result

0.12.1 is completed as an audit-only subphase. It produced the full model inventory, model classification, cross-layer usage analysis, problem and risk mapping, target separation definition and closure documentation required before code-level extraction begins.

The audit baseline records 125 model-like structs, with 11 hybrid/transitional structures requiring explicit separation planning. No Go code was modified during 0.12.1.

### 0.12.2 Internal Sub-Subphase Plan

| Sub-subphase | Name | Status |
| --- | --- | --- |
| 0.12.2.0 | Definition & Documentation Lock | ✅ Completed |
| 0.12.2.1 | Read Model Design | ✅ Completed |
| 0.12.2.2 | Read Model Implementation | ✅ Completed |
| 0.12.2.3 | Domain/Application → Read Mapping | ✅ Completed |
| 0.12.2.4 | Response Alignment | ✅ Completed |
| 0.12.2.5 | Validation & Compatibility | ✅ Completed |
| 0.12.2.6 | Documentation & Closure | ✅ Completed |

### 0.12.2.0 Result

0.12.2.0 is documentation-only and is completed in this repository state. It defines the Read Model Extraction sequence before code-level design starts, records the response-compatibility constraints and requires the 0.12.1 audit artifacts to drive all extraction decisions.

No Go code is changed in 0.12.2.0.

### 0.12.2.6 Result

0.12.2 is completed. It introduced explicit read model packages for `auth`, `user` and `usersettings`, explicit mapping functions for Domain/Application → Read transformations, response alignment preserving public JSON contracts and a validation record backed by `go test ./...`.

### 0.12.3 Internal Sub-Subphase Plan

| Sub-subphase | Name | Status |
| --- | --- | --- |
| 0.12.3.0 | Definition & Documentation Lock | ✅ Completed |
| 0.12.3.1 | Write Model Design | ✅ Completed |
| 0.12.3.2 | Write Model Implementation | ✅ Completed |
| 0.12.3.3 | Write → Domain Mapping | ✅ Completed |
| 0.12.3.4 | Handler Alignment | ✅ Completed |
| 0.12.3.5 | Validation & Compatibility | ✅ Completed |
| 0.12.3.6 | Documentation & Closure | ✅ Completed |

### 0.12.3.0 Result

0.12.3.0 is documentation-only and is completed in this repository state. It defines the Write Model Isolation sequence before code-level design starts, records the input-compatibility constraints and requires the completed read-side boundary from 0.12.2 to remain untouched.

No Go code is changed in 0.12.3.0.

### 0.12.3.6 Result

0.12.3 is completed. It introduced explicit write model packages for `auth`, `user` and `usersettings`, domain write input structures, Write → Domain mapping functions, handler alignment preserving public request payload semantics and validation documentation backed by `go test ./...`.

### 0.12.4 Internal Sub-Subphase Plan

| Sub-subphase | Name | Status |
| --- | --- | --- |
| 0.12.4.0 | Definition & Documentation Lock | ✅ Completed |
| 0.12.4.1 | Mapping Layer Design | ✅ Completed |
| 0.12.4.2 | Mapping Layer Implementation | ✅ Completed |
| 0.12.4.3 | Mapping Consolidation | ✅ Completed |
| 0.12.4.4 | Application Refactor | ✅ Completed |
| 0.12.4.5 | Validation & Compatibility | ✅ Completed |
| 0.12.4.6 | Documentation & Closure | ✅ Completed |

### 0.12.4.0 Result

0.12.4.0 is documentation-only and is completed in this repository state. It defines the centralized Mapping Layer Introduction sequence before code-level design starts, records the mapper ownership target and requires read-side and write-side compatibility to remain unchanged during consolidation.

No Go code is changed in 0.12.4.0.

### Next

0.12.5 — Contract Alignment.


### 0.12.4.6 Result

0.12.4 is completed. The phase introduced module-local `mappers` packages, consolidated read-side and write-side mapping ownership, reduced residual application mapping, preserved public HTTP contracts and recorded validation with `go test ./...`.

Next planned phase: 0.12.5 — Contract Alignment.

### 0.12.5 Internal Sub-Subphase Plan

| Sub-subphase | Name | Status |
| --- | --- | --- |
| 0.12.5.0 | Definition & Documentation Lock | ✅ Completed |
| 0.12.5.1 | Contract Inventory & Classification | ✅ Completed |
| 0.12.5.2 | Contract Normalization Design | ✅ Completed |
| 0.12.5.3 | Contract Alignment Implementation | ✅ Completed |
| 0.12.5.4 | Handler Contract Adjustment | ✅ Completed |
| 0.12.5.5 | Validation & Compatibility | ✅ Completed |
| 0.12.5.6 | Documentation & Closure | ✅ Completed |

### 0.12.5.0 Result

0.12.5.0 is documentation-only and is completed in this repository state. It defines the Contract Alignment sequence before code-level contract changes start, records the provider contract scope and requires runtime compatibility with the public HTTP/API surface to remain unchanged.

No Go code is changed in 0.12.5.0.

### 0.12.5.6 Result

0.12.5 is completed. It produced the contract inventory, contract normalization design, contract alignment implementation, handler contract adjustment and validation documentation. HTTP contract aliases are centralized, handlers preserve public JSON payloads and validation is recorded with `go test ./...`.

### Next

Phase 0.12 is complete. Continue with the next roadmap-defined phase.

---

## ✅ Phase 0.13 — Provider Layer Consolidation

### Status

✅ Completed

### Objective

Consolidate the Provider Layer as the explicit, consistent and reusable entry point to domain services.

Phase 0.13 is the continuation of the Phase 0.12 stabilization line. Phase 0.12 clarified model direction and mapper ownership; Phase 0.13 clarifies orchestration ownership. The completion state therefore means more than a status transition: the backend now has a documented provider-first runtime boundary for authenticated module composition while preserving public contracts.

### Scope

Included:

- provider inventory and classification
- provider interface design
- provider implementation where required
- application and handler integration
- validation and compatibility checks
- documentation closure

Excluded:

- public API changes
- business behavior changes
- API versioning changes
- CQRS or event sourcing
- observability implementation

### Subphase Plan

| Subphase | Name | Status |
| --- | --- | --- |
| 0.13.0 | Definition & Documentation Lock | ✅ Completed |
| 0.13.1 | Provider Inventory & Classification | ✅ Completed |
| 0.13.2 | Provider Interface Design | ✅ Completed |
| 0.13.3 | Provider Implementation | ✅ Completed |
| 0.13.4 | Application Integration | ✅ Completed |
| 0.13.5 | Validation & Compatibility | ✅ Completed |
| 0.13.6 | Documentation & Closure | ✅ Completed |

### 0.13.0 Result

0.13.0 completed the documentation lock for Provider Layer Consolidation, corrected the roadmap-level Phase 0.13 definition and expanded the subphase sequence before provider inventory began.

### 0.13.1 Result

0.13.1 completed provider inventory and classification. It identified existing provider-like boundaries, missing provider boundaries, compatibility wiring and risky direct access patterns without changing Go code.

### 0.13.2 Result

0.13.2 completed provider interface design. It locked target provider interface boundaries for session/bootstrap, authenticated account/profile/settings and wallet orchestration flows while preserving read/write model separation and public HTTP/API contracts.

### 0.13.3 Result

0.13.3 implemented the first concrete provider boundary. The auth module exposes explicit provider interfaces, the composition root builds a consolidated auth provider and HTTP handlers route relevant operations through that provider boundary while preserving public behavior.

### 0.13.4 Result

0.13.4 completed application integration. Runtime HTTP wiring now consumes the consolidated provider boundary directly through router construction instead of passing lower-level user services, settings services, wallet stores and challenge configuration through production handler wiring.

### 0.13.5 Result

0.13.5 completed validation and compatibility. External validation confirmed `make build` and `go test ./...` after the 0.13.4 router compatibility fix path. Public HTTP/API behavior, versioned route behavior, error envelopes and module compatibility remain unchanged.

### 0.13.6 Result

0.13.6 completed documentation closure. The closure corrected misplaced roadmap content, removed generic repeated status-only blocks from trunk documentation, expanded the final subphase state consistently and documented the actual Provider Layer impact in architecture, flows, testing, development and observability documentation.

### Next

Phase 0.13 is complete. Continue with the corrected next roadmap-defined phase: 0.14 — Observability & Diagnostics Foundation.

---

# Phase 0.14 — Observability & Diagnostics Foundation

## Status

**Phase:** 0.14 — Observability & Diagnostics Foundation  
**Current subphase:** 0.14.6 — Validation & Documentation  
**Status:** Phase 0.14 completed

## Context

Phase 0.12 clarified read/write model direction and mapper ownership. Phase 0.13 then consolidated runtime composition around the Provider Layer. Those phases left the backend with cleaner internal ownership and stable public behavior, but they did not solve runtime visibility.

The next foundation gap is observability: the system needs request correlation, consistent logging context, internal error metadata and traceable flow movement so debugging no longer depends only on manually reading isolated logs.

## Problem

The backend remains operationally opaque in several areas:

- requests are not consistently correlated
- logs do not uniformly expose request-scoped context
- internal errors lack enough diagnostic metadata for efficient debugging
- HTTP → Provider → Application → Domain → Repository flow is not yet observable as one path

## Decision

Phase 0.14 will introduce an Observability & Diagnostics Foundation.

This decision replaces the previous placeholder roadmap label for 0.14 with the actual Phase 0.14 scope defined for this stage. The change is documentation-lock work only and does not modify Go code.

## Subphase Plan

- 0.14.0 — Phase Definition & Documentation Lock ✅ Completed
- 0.14.1 — Correlation Model (Request ID / Trace) ✅ Completed
- 0.14.2 — Logging Standardization ✅ Completed
- 0.14.3 — Error Context Enrichment ✅ Completed
- 0.14.4 — Flow Tracing Integration ✅ Completed
- 0.14.5 — Diagnostics Surface Exposure ✅ Completed
- 0.14.6 — Validation & Documentation ✅ Completed

## Observable Impact

After the documentation lock, all trunk references to the next phase must point to `0.14 — Observability & Diagnostics Foundation`, and the implementation phase can proceed without ambiguity about scope or sequencing.

### 0.14.1 Result

0.14.1 completed the first implementation step of the Observability & Diagnostics Foundation. The existing HTTP request ID behavior was formalized into a reusable correlation accessor without changing the public request/response contract.

The backend continues to honor `X-Request-Id` when supplied by a client, generates one when missing, mirrors the effective value in the response header and stores it in request context. `AccessLog` and `Recoverer` now consume that value through `httpx.RequestIDFromContext`, preserving the private context key while allowing future logging, error and flow tracing work to reuse the model safely.

Dedicated tests were added for inherited request IDs, generated request IDs and empty correlation lookup for missing or nil context. Full `go test ./...` could not be executed in this environment because the Go toolchain attempted to download Go 1.25.0 from `proxy.golang.org` and DNS/network access was unavailable.

### 0.14.2 Result

0.14.2 completed the logging standardization step of the Observability & Diagnostics Foundation. The backend already used structured JSON logging through `slog`, so this subphase did not replace the logger implementation.

The inherited correlation model from 0.14.1 is now reflected in log attributes through the canonical `request_id` key. `internal/core/logger` owns the standard key and helper methods, while `internal/core/httpx` remains responsible for reading the effective request ID from request context.

HTTP access logs and panic recovery logs now emit `request_id` instead of the previous local `rid` attribute. This keeps runtime behavior stable while giving later 0.14 subphases a consistent field for error context enrichment and flow tracing.

Dedicated logger tests were added for request ID attribute construction and request-scoped logger derivation behavior. Full `go test ./...` could not be executed in this environment because the Go toolchain attempted to download Go 1.25.0 from `proxy.golang.org` and DNS/network access was unavailable.

### 0.14.3 Result

0.14.3 completed the error context enrichment step of the Observability & Diagnostics Foundation. The backend already had a stable public error envelope from Phase 0.8, including the `details` field, so this subphase did not replace the error model or alter response shape.

The inherited correlation and logging work from 0.14.1 and 0.14.2 is now reflected in the error layer through safe enrichment helpers on `errs.AppError`. The helpers allow internal code to attach contextual metadata such as `request_id` through a controlled API while preserving immutable-style error handling.

`ToResponseError()` now serializes a copied public details map instead of exposing the internal map directly. This prevents accidental mutation leakage between internal error state and response construction, while preserving the existing public `{ error: { code, message, details } }` contract.

Dedicated error tests were added for context enrichment, request ID enrichment, empty request ID behavior, public details copy behavior and response details copy behavior. Full `go test ./...` could not be executed in this environment because the Go toolchain attempted to download Go 1.25.0 from `proxy.golang.org` and DNS/network access was unavailable.

### 0.14.4 Result

0.14.4 completed the flow tracing integration step of the Observability & Diagnostics Foundation.

Context inherited from 0.14.3: the backend already had request correlation, standardized `request_id` logging and safe error-context enrichment. The remaining gap was explicit flow movement: logs could describe outcomes, but not consistently mark request and application lifecycle transitions.

Problem addressed: operators needed a minimal way to follow HTTP request execution and application lifecycle movement without introducing external tracing infrastructure or changing public contracts.

Decision taken: flow tracing is represented through existing structured JSON logs using message `flow_trace` and canonical field `flow_event`. The logger package owns the reusable key and helper, while HTTP keeps responsibility for extracting the request ID from context.

Concrete change: `internal/core/httpx.AccessLog` now emits `http_request_start` and `http_request_end`; `internal/app.App` emits `application_start` and `application_stop`; logger tests validate flow event attributes.

Observable impact: request and lifecycle movement is now visible through existing logs, correlated by `request_id` where available. Public API behavior, response payloads, provider behavior and business logic remain unchanged.

### 0.14.5 Result

0.14.5 completed the diagnostics surface exposure step of the Observability & Diagnostics Foundation.

Context inherited from 0.14.4: the backend already had request correlation, standardized structured logging, error context enrichment and flow tracing events. Those capabilities were visible in logs, but there was no minimal runtime surface to report whether the observability foundation itself was active.

Problem addressed: operators needed a lightweight way to inspect observability readiness without introducing Prometheus, OpenTelemetry, dashboards or business-contract changes.

Decision taken: expose a minimal `GET /diagnostics` surface through the existing status/httpx path. The payload reports service identity, runtime environment/version/commit and the enabled foundation capabilities: request correlation, structured logging, error context enrichment and flow tracing.

Observable impact: `/diagnostics` now provides a stable foundation-level diagnostic snapshot while existing `/health`, `/readiness`, `/version`, business routes, error envelopes, providers and domain behavior remain unchanged.

### 0.14.6 Result

0.14.6 completed the validation and documentation closure for the Observability & Diagnostics Foundation.

Context inherited from 0.14.5: the backend already had request correlation, standardized structured logging, safe error context enrichment, flow tracing events and a minimal `/diagnostics` surface. The remaining risk was documentary drift: the implementation was correct, but the phase needed final reconciliation so the trunk documents, handoff state, roadmap and phase narrative described the same system.

Problem addressed: Phase 0.14 intentionally used a single phase document instead of creating one dedicated document per subphase. That avoided duplication while the implementation was evolving, but it required a deeper closure pass to preserve the narrative standard established in Phase 0.12 and reinforced in Phase 0.13.

Decision taken: close Phase 0.14 through documentation reconciliation rather than new code. The closure pass preserves the existing documents, updates only the affected trunk references and records the complete subphase narrative in the phase document and observability document.

Concrete change: the phase status, roadmap, index, handoff, README, observability model and Phase 0.14 document now agree that all subphases from 0.14.0 through 0.14.6 are completed.

Validation recorded:

```bash
go test ./...
```

The command completed successfully in the local project environment after 0.14.5 was applied.

Observable impact: Phase 0.14 is closed with coherent documentation and no additional code changes. The backend remains behavior-compatible while exposing its internal observability foundation through request IDs, structured logs, enriched errors, flow events and `/diagnostics`.

---

## Phase 0.14 Final State

Phase 0.14 — Observability & Diagnostics Foundation is complete.

The completed capability set is:

- request correlation through the HTTP boundary
- canonical `request_id` structured logging
- safe error-context enrichment
- minimal flow tracing events
- diagnostics surface exposure through `GET /diagnostics`
- final validation and documentation reconciliation

No business API contract, authentication contract, provider contract, domain behavior, public error envelope, Prometheus integration, OpenTelemetry integration or dashboard system was introduced.
