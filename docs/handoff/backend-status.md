# 📦 Backend Status — Handoff

## Current Implementation State

The backend is currently operating with a stable authentication base that includes:

- JWT authentication
- development login flow
- EVM wallet-based authentication
- wallet challenge persistence
- wallet identity persistence
- WebSocket auth session propagation

---

## Stage / Phase Reference

- **Stage:** 0 — Foundation
- **Phase:** 0.4 — Auth and User Stabilization
- **Current completed subphase:** 0.4.6 — Wallet Identity Persistence and Durable Challenge Storage

---

## Wallet Authentication Status

| Capability | Status |
|------------|--------|
| Challenge issuance | ✅ |
| Signature verification | ✅ |
| Replay protection | ✅ |
| Durable challenge storage | ✅ |
| Persistent wallet identity creation | ✅ |
| JWT issuance with wallet metadata | ✅ |
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

This gives the backend a durable wallet-auth anchor that can be used later for:

- wallet ↔ user linking
- account ownership modeling
- multi-wallet account expansion

---

## Fallback Behavior

When PostgreSQL is not enabled, the backend still supports local development through in-memory fallback stores for:

- wallet challenges
- wallet identities

This preserves dev ergonomics without blocking durable production behavior when DB is available.

---

## Current Limitations

The following items are still pending after 0.4.6:

- no wallet ↔ user linking yet
- no refresh-token lifecycle
- no revocation support
- no persistent session store
- no unified account model
- no multi-wallet ownership model

---

## Operational Readiness

### Ready For

- stable internal development
- challenge/signature wallet auth validation
- JWT-based session propagation
- horizontally scalable auth storage assumptions once DB is enabled

### Not Yet Ready For

- complete exchange account ownership model
- production-grade account recovery flows
- wallet linking management UX/API
- advanced auth lifecycle management

---

## Recommended Next Phase

### 0.4.7 — Wallet ↔ User Linking and Unified Identity Model

This should become the next focus area in order to connect wallet identities with durable platform users and prepare the backend for exchange account semantics.