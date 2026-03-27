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

---

## ✅ Phase 0.4.14 Closure Summary

Phase 0.4.14 closes the first post-detach lifecycle gap without introducing heavy new persistence semantics.

The backend now explicitly validates that a detached wallet identity remains reusable: it can be reattached again through the authenticated wallet-link flow, and it can also re-enter the wallet-login bootstrap flow to resolve back into a wallet-owned user identity.

### Delivered in 0.4.14

- explicit lifecycle clarification for detached wallet identities
- service-level reattachment coverage for detached wallets under authenticated linking
- service-level wallet-login rebound coverage for detached wallets under wallet bootstrap auth
- handler-level coverage proving detached wallets can be linked again by the authenticated owner
- documentation updates aligning detach execution with post-detach reuse semantics
- no schema expansion required to make current detached-wallet behavior explicit

---

## 🔍 Functional Result

The system now supports the following post-detach lifecycle sequence:

1. user detaches one already eligible owned wallet
2. backend clears `user_id`, `linked_at`, and `is_primary` from that wallet identity
3. wallet identity remains known to the backend by address and wallet identity ID
4. authenticated user can later reattach that wallet again through `POST /auth/wallets/link/challenge` + `POST /auth/wallets/link/verify`
5. detached wallet can also re-enter `POST /auth/wallet/verify` and resolve back into a wallet-owned user identity
6. no archival or deletion semantics are required for the current lifecycle model

---

## ❌ Not Included in 0.4.14

The following items remain intentionally out of scope:

- detached-identity audit columns
- archival / soft-delete markers for detached wallets
- recovery or dispute workflows around detached ownership
- automatic primary replacement for risky detach cases
- merge between wallet identities and future auth methods
- refresh tokens
- revocation flows
- persistent authenticated session storage

---

## ⏭️ Next Phase

### 0.4.15 — Detached Identity Audit Readiness

Planned focus:

- decide whether detached identities require explicit audit metadata
- evaluate whether detach history should become queryable later
- preserve current reusable detached-wallet semantics while preparing richer lifecycle observability
