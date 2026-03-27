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

---

## ✅ Phase 0.4.12 Closure Summary

Phase 0.4.12 adds the first explicit detach-preparation contract on top of the ownership model delivered through 0.4.11.

The backend can now accept an authenticated request from an already logged-in user, validate ownership of an already linked wallet, and evaluate whether that wallet is currently safe to detach without executing unlink behavior.

### Delivered in 0.4.12

- authenticated detach-eligibility endpoint: `POST /auth/wallets/detach/check`
- conservative eligibility contract for future unlink work
- explicit response fields for `eligible`, `is_primary`, `owned_wallet_count`, and `reasons`
- protected rejection when the target wallet is missing
- protected rejection when the wallet does not belong to the current authenticated user
- explicit non-eligibility when the wallet is primary
- explicit non-eligibility when detach would leave the user without any wallets
- existing link, merge, and primary-switch coverage preserved

---

## 🔍 Functional Result

The system now supports the following detach-preparation sequence under an existing authenticated session:

1. user authenticates normally
2. user lists or already knows their owned wallets
3. user requests `POST /auth/wallets/detach/check` with one owned wallet address
4. backend validates that the wallet exists and belongs to the authenticated user
5. backend evaluates whether the wallet is currently primary
6. backend evaluates how many wallets the user currently owns
7. backend returns an eligibility decision plus explicit reasons when detach is not yet safe

---

## ❌ Not Included in 0.4.12

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

### 0.4.13 — Wallet Detach Execution Design

Planned focus:

- transform detach eligibility into a real detach contract
- define how primary replacement must occur before detach execution
- preserve strict ownership invariants during unlink execution
- continue progression from ownership guardrails toward controlled wallet detachment
