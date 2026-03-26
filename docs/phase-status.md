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

---

## ✅ Phase 0.4.11 Closure Summary

Phase 0.4.11 adds the first explicit post-merge wallet-ownership management operation on top of the linking and merge model delivered through 0.4.10.

The backend can now accept an authenticated request from an already logged-in user, validate ownership of an already linked wallet, and atomically reassign the `is_primary` flag so the requested owned wallet becomes the single primary wallet for that user.

### Delivered in 0.4.11

- authenticated primary-wallet switch endpoint: `POST /auth/wallets/primary`
- store-level primary-wallet reassignment contract
- protected rejection when the target wallet is missing
- protected rejection when the wallet does not belong to the current authenticated user
- deterministic single-primary enforcement during reassignment
- refreshed owned-wallet inventory returned after successful primary switch
- existing wallet-link and wallet-account-merge coverage preserved

---

## 🔍 Functional Result

The system now supports the following explicit primary-wallet reassignment sequence under an existing authenticated session:

1. user authenticates normally
2. user lists or already knows their owned wallets
3. user requests `POST /auth/wallets/primary` with one owned wallet address
4. backend validates that the wallet exists and belongs to the authenticated user
5. backend atomically clears the previous primary flag for that user
6. backend marks the requested wallet as the only primary wallet
7. updated owned-wallet inventory is returned with the new primary wallet first

---

## ❌ Not Included in 0.4.11

The following items remain intentionally out of scope:

- wallet unlink endpoint
- arbitrary cross-user wallet transfer
- archival / alias records for merged source users
- merge between wallet identities and future auth methods
- refresh tokens
- revocation flows
- persistent authenticated session storage

---

## ⏭️ Next Phase

### 0.4.12 — Wallet Ownership Detach Contract Preparation

Planned focus:

- introduce safe wallet detach / unlink semantics
- preserve primary-wallet invariants during future detach operations
- deepen merge-safe identity rules
- continue progression from ownership model toward account-level management
