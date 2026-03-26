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

---

## ✅ Phase 0.4.10 Closure Summary

Phase 0.4.10 executes the first protected wallet-owned account merge on top of the wallet-linking and ownership model delivered in 0.4.9.

The backend can now accept an authenticated request from an already logged-in user, generate an account-merge challenge for a wallet that already belongs to another wallet-owned user, verify the signed challenge, and atomically reassign all wallets from that source wallet-owned account into the current authenticated user.

### Delivered in 0.4.10

- authenticated account-merge challenge endpoint
- authenticated account-merge verify endpoint
- challenge purpose expansion:
  - `auth_bootstrap`
  - `wallet_link`
  - `account_merge`
- protected rejection when merge challenge is signed by the wrong wallet
- protected rejection when the source wallet is unlinked
- protected rejection when merge is not required because the wallet already belongs to the current user
- atomic wallet-ownership consolidation from source user to current authenticated user
- deterministic primary-wallet preservation rules during merge
- updated wallet inventory returned after successful merge

---

## 🔍 Functional Result

The system now supports the following wallet-owned account merge sequence under an existing authenticated session:

1. user authenticates normally
2. user requests account-merge challenge for a wallet already linked elsewhere
3. challenge is persisted with `account_merge` purpose and `requested_by_user_id`
4. source wallet signs the merge challenge
5. backend verifies signature and user-bound challenge
6. source wallet identity is resolved
7. source user is derived from wallet ownership
8. all source-user wallets are reassigned to the authenticated target user
9. updated owned-wallet list is returned

---

## ❌ Not Included in 0.4.10

The following items remain intentionally out of scope:

- wallet unlink endpoint
- primary-wallet switch endpoint
- arbitrary cross-user wallet transfer
- archival / alias records for merged source users
- merge between wallet identities and future auth methods
- refresh tokens
- revocation flows
- persistent authenticated session storage

---

## ⏭️ Next Phase

### 0.4.11 — Wallet Ownership Management and Primary-Control Progression

Planned focus:

- introduce safe wallet detach / unlink semantics
- define protected primary-wallet switching
- deepen merge-safe identity rules
- continue progression from ownership model toward account-level management