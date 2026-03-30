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

- wallet inventory query semantics and filtering preparation on top of the enriched read model
- preserve backward compatibility of the public wallet inventory contract
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

