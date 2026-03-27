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

---

## ✅ Phase 0.4.13 Closure Summary

Phase 0.4.13 turns the detach-eligibility contract from 0.4.12 into a real, authenticated ownership mutation for the already safe cases.

The backend can now accept an authenticated request from an already logged-in user, validate ownership of an already linked wallet, reuse the same conservative detach guardrails, and execute the detach only when the wallet is both owned and eligible.

### Delivered in 0.4.13

- authenticated detach execution endpoint: `POST /auth/wallets/detach`
- detach service execution path that reuses the 0.4.12 eligibility rules before mutating ownership
- store-level detach contract that clears `user_id`, `linked_at`, and `is_primary` from the detached wallet identity
- updated wallet inventory response after successful detach execution
- protected conflict response when the wallet is not detach-eligible under current guardrails
- existing link, merge, primary-switch, and detach-check coverage preserved

---

## 🔍 Functional Result

The system now supports the following detach execution sequence under an existing authenticated session:

1. user authenticates normally
2. user lists or already knows their owned wallets
3. user requests `POST /auth/wallets/detach` with one owned wallet address
4. backend validates that the wallet exists and belongs to the authenticated user
5. backend reuses the detach eligibility rules from 0.4.12
6. backend rejects the request if the wallet is primary or if detach would leave the user without wallets
7. backend clears ownership metadata from the detached wallet and returns the refreshed remaining wallet inventory

---

## ❌ Not Included in 0.4.13

The following items remain intentionally out of scope:

- wallet unlink execution
- arbitrary cross-user wallet transfer
- archival / alias records for merged source users
- automatic replacement of primary during detach
- merge between wallet identities and future auth methods
- refresh tokens
- revocation flows
- persistent authenticated session storage

---

## ⏭️ Next Phase

### 0.4.14 — Detach Follow-Up Semantics and Source Identity Lifecycle

Planned focus:

- define whether detached wallets should bootstrap fresh wallet-only users automatically in future flows
- evaluate whether detached identities require archival, history, or audit-friendly lifecycle markers
- preserve strict ownership invariants while preparing richer unlink semantics
