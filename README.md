# SCAVO Exchange — Backend

## 🧠 Overview

SCAVO Exchange Backend is a Go-based service that provides authentication, user management, and wallet-based identity for the SCAVO ecosystem.

The system is designed following a **wallet-first identity model**, progressively evolving into a **unified account architecture** capable of supporting exchange-grade features such as multi-wallet ownership, account consolidation, and future compliance layers.

---

## 🏗️ Architecture Principles

- **Wallet-first authentication**
- **Durable user abstraction**
- **Stateless JWT sessions**
- **Progressive identity consolidation**
- **Database-backed persistence with in-memory fallback**

---

## 🚧 Current Stage

**Stage:** 0 — Foundation  
**Phase:** 0.4 — Auth and User Stabilization  
**Current Subphase:** **0.4.8 — Account Consolidation and Multi-Wallet Ownership Foundations**

---

## 🔐 Authentication Model

The backend supports two authentication methods:

### 1. Password-based (dev only)
- Used for internal testing
- Not intended for production

### 2. Wallet-based authentication (EVM)

Flow:

1. Client requests challenge  
   `POST /auth/wallet/challenge`

2. Server generates challenge:
   - unique ID
   - message to sign
   - expiration

3. Client signs message with wallet

4. Client verifies:
   `POST /auth/wallet/verify`

5. Backend:
   - verifies signature
   - consumes challenge (one-time use)
   - resolves wallet identity
   - resolves or creates durable user
   - links wallet → user
   - issues JWT

---

## 🧩 Identity Model Evolution

### Pre 0.4.7
- Wallet identity existed independently
- No durable user linkage

### 0.4.7 — Wallet ↔ User Linking
- Each wallet identity is linked to a durable user
- JWT identity becomes unified

### 0.4.8 — Multi-Wallet Ownership Foundations

Wallet identities now support ownership metadata:

- `user_id`
- `linked_at`
- `is_primary`

This enables:

- One user → multiple wallets
- Single primary wallet designation
- Ownership persistence independent of sessions

---

## 🗄️ Persistence Model

### Tables involved

#### `auth_wallet_challenges`
- challenge lifecycle
- expiration + one-time use

#### `auth_wallet_identities`
- wallet address registry
- ownership metadata:
  - `user_id`
  - `linked_at`
  - `is_primary`

#### `users`
- durable platform users
- wallet-backed or future auth methods

---

## 🔌 API Endpoints

### Wallet Auth

#### `POST /auth/wallet/challenge`

Request:
```json
{
  "address": "0x...",
  "chain": "scavium"
}
```

---

#### `POST /auth/wallet/verify`

Request:
```json
{
  "challenge_id": "...",
  "address": "0x...",
  "signature": "0x..."
}
```

Response:
```json
{
  "access_token": "...",
  "token_type": "Bearer",
  "expires_in": 3600,
  "user_id": "...",
  "wallet_id": "...",
  "wallet_address": "0x...",
  "chain": "scavium",
  "auth_method": "wallet_evm"
}
```

---

#### `GET /auth/wallets`

Returns all wallet identities linked to the authenticated user.

Response:
```json
{
  "wallets": [
    {
      "id": "...",
      "address": "0x...",
      "user_id": "...",
      "linked_at": "...",
      "is_primary": true
    }
  ]
}
```

---

## 🧾 JWT Claims

Tokens include:

- `user_id`
- `wallet_id`
- `wallet_address`
- `auth_method`
- `exp`, `iat`, `nbf`

---

## 🧪 Testing

Basic validation:

```bash
go test ./...
```

---

## 🚧 What 0.4.8 Solves

- Durable wallet challenge storage
- Persistent wallet identity creation
- Durable wallet ↔ user linkage
- Unified wallet-backed session identity
- Multi-wallet ownership foundations
- Primary wallet semantics
- Read-only wallet listing for authenticated users
- Protection against cross-user wallet reassignment

---

## ❌ What 0.4.8 Does Not Solve Yet

- User-driven wallet linking API
- Wallet unlink API
- Account merge workflows
- Cross-user wallet transfers
- Refresh tokens
- Token revocation
- Persistent session storage

---

## 🧭 Next Phase

### 0.4.9 — User-Driven Wallet Linking Contract and Protected Account Merge Preparation

This phase will introduce:

- controlled wallet linking
- ownership validation flows
- merge-safe identity contracts
- groundwork for account consolidation

---

## 🧩 Summary

At the end of Phase 0.4.8:

- Wallet authentication is stable
- Identity model is unified
- Multi-wallet ownership is structurally supported
- The backend is ready for controlled identity expansion

---
