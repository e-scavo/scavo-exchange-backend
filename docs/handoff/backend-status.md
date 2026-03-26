# 📦 Backend Status — Handoff

## Current Implementation State

The backend is currently operating with a stable authentication base that includes:

- JWT authentication
- development login flow
- EVM wallet-based authentication
- wallet challenge persistence
- wallet identity persistence
- wallet ↔ user linkage
- WebSocket auth session propagation

---

## Stage / Phase Reference

- **Stage:** 0 — Foundation
- **Phase:** 0.4 — Auth and User Stabilization
- **Current completed subphase:** 0.4.7 — Wallet ↔ User Linking and Unified Identity Model

---

## Wallet Authentication Status

| Capability | Status |
|------------|--------|
| Challenge issuance | ✅ |
| Signature verification | ✅ |
| Replay protection | ✅ |
| Durable challenge storage | ✅ |
| Persistent wallet identity creation | ✅ |
| Persistent wallet ↔ user linkage | ✅ |
| JWT issuance with wallet metadata | ✅ |
| Unified user hydration | ✅ |
| HTTP session exposure | ✅ |
| WebSocket session exposure | ✅ |

---

## Persistence Model

### Wallet Challenges

Challenges are stored in PostgreSQL and contain:

- unique ID
- target address
- chain
- nonce
- canonical signable message
- issued timestamp
- expiration timestamp
- used timestamp

Consumption is enforced through a DB-backed flow so the same challenge cannot be reused successfully.

### Wallet Identities

Wallet identities are stored in PostgreSQL and resolved uniquely by normalized wallet address.

Each wallet identity can now be linked to a durable platform user through `auth_wallet_identities.user_id`.

### Users

Wallet-authenticated users are now provisioned in the shared `users` table.

This gives the backend a durable unified identity anchor that can be used later for:

- account ownership modeling
- cross-auth method expansion
- multi-wallet account evolution

---

## Fallback Behavior

When PostgreSQL is not enabled, the backend still supports local development through in-memory fallback stores for:

- wallet challenges
- wallet identities
- wallet-to-user resolution

This preserves dev ergonomics without blocking durable production behavior when DB is available.

---

## Current Limitations

The following items are still pending after 0.4.7:

- no user-facing wallet management API yet
- no refresh-token lifecycle
- no revocation support
- no persistent session store
- no multi-wallet ownership model
- no account merge semantics across login methods

---

## Operational Readiness

### Ready For

- stable internal development
- challenge/signature wallet auth validation
- durable wallet-backed user creation
- JWT-based unified session propagation
- horizontally scalable auth storage assumptions once DB is enabled

### Not Yet Ready For

- complete exchange account aggregation model
- production-grade account recovery flows
- wallet linking management UX/API
- advanced auth lifecycle management

---

## Recommended Next Phase

### 0.4.8 — Account Consolidation and Multi-Wallet Ownership Foundations

This should become the next focus area in order to evolve the current 1:1 wallet linkage model into broader exchange account ownership semantics.
