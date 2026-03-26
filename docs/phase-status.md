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

---

## ✅ Phase 0.4.9 Closure Summary

Phase 0.4.9 introduces the first authenticated wallet-management contract on top of the durable identity and ownership model created in previous subphases.

The backend can now accept an authenticated request from an already logged-in user, generate a wallet-linking challenge, verify the signed challenge, and attach a secondary wallet without issuing a new login session or performing risky account merge heuristics.

### Delivered in 0.4.9

- authenticated wallet-link challenge endpoint
- authenticated wallet-link verify endpoint
- challenge purpose separation:
  - `auth_bootstrap`
  - `wallet_link`
- persisted challenge metadata:
  - `purpose`
  - `requested_by_user_id`
- protected rejection of mismatched user-bound link challenges
- protected rejection of linking a wallet already owned by another user
- secondary-wallet attachment with `is_primary = false`
- updated wallet inventory returned after successful linking

---

## 🔍 Functional Result

The system now supports the following linked-wallet sequence under an existing authenticated session:

1. user authenticates normally
2. user requests wallet-link challenge
3. challenge is persisted with `wallet_link` purpose and `requested_by_user_id`
4. user signs with the secondary wallet
5. backend verifies signature and user-bound challenge
6. wallet identity is resolved
7. ownership conflict is checked
8. wallet is attached as a non-primary wallet
9. updated owned-wallet list is returned

---

## ❌ Not Included in 0.4.9

The following items remain intentionally out of scope:

- wallet unlink endpoint
- primary-wallet switch endpoint
- cross-user wallet transfer
- automatic account-merge execution
- merge between wallet identities and future auth methods
- refresh tokens
- revocation flows
- persistent authenticated session storage

---

## ⏭️ Next Phase

### 0.4.10 — Wallet Ownership Management and Merge-Safe Identity Progression

Planned focus:

- introduce safe wallet detach / unlink semantics
- define protected primary-wallet switching
- deepen merge-safe identity rules
- continue progression from ownership model toward account-level management