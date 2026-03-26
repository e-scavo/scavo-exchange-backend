# SCAVO Exchange — Backend

## 🧠 Overview

SCAVO Exchange Backend is a Go-based service that provides authentication, user management, and wallet-based identity for the SCAVO ecosystem.

The backend follows a **wallet-first identity model** that progressively evolves into a **durable account architecture** suitable for exchange-grade ownership, linking, and future multi-auth identity expansion.

---

## 🏗️ Architecture Principles

- **Wallet-first authentication**
- **Durable user abstraction**
- **Stateless JWT sessions**
- **Explicit ownership persistence**
- **Incremental account consolidation**
- **Database-backed persistence with in-memory fallback**

---

## 🚧 Current Stage

**Stage:** 0 — Foundation  
**Phase:** 0.4 — Auth and User Stabilization  
**Current Subphase:** **0.4.9 — User-Driven Wallet Linking Contract and Protected Account Merge Preparation**

---

## 🔐 Authentication Model

The backend currently supports two authentication methods:

### 1. Password-based authentication (dev only)
- intended only for internal development and testing
- not meant for production operation

### 2. Wallet-based authentication (EVM)

Base wallet login flow:

1. Client requests challenge  
   `POST /auth/wallet/challenge`

2. Backend creates challenge:
   - unique ID
   - wallet address binding
   - chain binding
   - expiration timestamp
   - challenge purpose metadata

3. Client signs the challenge message

4. Client verifies challenge  
   `POST /auth/wallet/verify`

5. Backend:
   - validates challenge state
   - verifies signature
   - consumes challenge
   - resolves wallet identity
   - resolves or creates durable user
   - enforces ownership invariants
   - issues JWT

---

## 🧩 Identity Model Evolution

### Pre 0.4.7
- wallet identity was not durably linked to a platform user
- session identity and persistent identity were not unified

### 0.4.7 — Wallet ↔ User Linking
- each wallet identity is linked to a durable user
- JWT identity becomes unified around `user_id`

### 0.4.8 — Multi-Wallet Ownership Foundations
wallet identities gained ownership metadata:

- `user_id`
- `linked_at`
- `is_primary`

This allowed:

- one user → multiple wallets
- explicit primary wallet designation
- ownership persistence independent from JWT sessions

### 0.4.9 — Authenticated Wallet Linking Contract
wallet management now supports an authenticated user-driven linking flow:

- `POST /auth/wallets/link/challenge`
- `POST /auth/wallets/link/verify`

This allows a signed secondary-wallet attachment flow without creating a new session or performing account merge heuristics.

---

## 🗄️ Persistence Model

### Main tables involved

#### `auth_wallet_challenges`
stores challenge lifecycle and now also includes linking metadata:

- `purpose`
- `requested_by_user_id`

Used for:
- wallet auth bootstrap challenges
- authenticated wallet-link confirmation challenges

#### `auth_wallet_identities`
stores wallet registry and ownership metadata:

- `id`
- `address`
- `user_id`
- `linked_at`
- `is_primary`

#### `users`
stores durable platform users:

- wallet-backed users
- future multi-auth identities

---

## 🔌 API Endpoints

### Wallet Auth

#### `POST /auth/wallet/challenge`

Creates a login bootstrap challenge for wallet authentication.

Request:

```json
{
  "address": "0x...",
  "chain": "scavium"
}
```

---

#### `POST /auth/wallet/verify`

Verifies wallet signature and returns a JWT-backed session.

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

### Wallet Ownership

#### `GET /auth/wallets`

Returns all wallet identities linked to the authenticated durable user.

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

### Authenticated Wallet Linking

#### `POST /auth/wallets/link/challenge`

Creates a wallet-linking challenge bound to the currently authenticated user.

Request:

```json
{
  "address": "0x...",
  "chain": "scavium"
}
```

Behavior:

- requires valid JWT
- challenge purpose becomes `wallet_link`
- challenge stores `requested_by_user_id`

---

#### `POST /auth/wallets/link/verify`

Verifies the linking signature and attaches the wallet to the current user as a **secondary wallet**.

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
  "linked_wallet": {
    "id": "...",
    "address": "0x...",
    "user_id": "...",
    "linked_at": "...",
    "is_primary": false
  },
  "wallets": [
    {
      "id": "...",
      "address": "0x...",
      "user_id": "...",
      "linked_at": "...",
      "is_primary": true
    },
    {
      "id": "...",
      "address": "0x...",
      "user_id": "...",
      "linked_at": "...",
      "is_primary": false
    }
  ]
}
```

---

## 🧾 JWT Claims

JWT tokens include:

- `user_id`
- `wallet_id`
- `wallet_address`
- `auth_method`
- `exp`
- `iat`
- `nbf`

Wallet linking does **not** mint a new token. It operates under the existing authenticated session.

---

## 🧪 Testing

Run:

```bash
go test ./...
```

Focus areas added in 0.4.9:

- authenticated link challenge generation
- link verification flow
- ownership conflict rejection
- secondary wallet persistence
- wallet inventory consistency after linking

---

## 🚧 What 0.4.9 Solves

- authenticated user-driven wallet linking
- challenge purpose separation between login and linking
- challenge-to-user binding through `requested_by_user_id`
- protected secondary-wallet attachment
- prevention of cross-user wallet takeover during linking
- wallet inventory refresh after successful link verification

---

## ❌ What 0.4.9 Does Not Solve Yet

- wallet unlink API
- primary-wallet switch API
- cross-user ownership transfer
- automatic account merge workflows
- merge between wallet-backed and other auth methods
- refresh tokens
- token revocation
- persistent authenticated sessions

---

## 🧭 Next Phase

### 0.4.10 — Wallet Ownership Management and Merge-Safe Identity Progression

Next expected focus:

- unlink / detach contract design
- protected primary-wallet switching
- deeper merge-safe identity preparation
- stronger account-level ownership operations

---

## 🧩 Summary

At the end of Phase 0.4.9:

- wallet authentication remains stable
- identity remains unified
- ownership remains protected
- authenticated wallet linking is now available
- the backend is ready for the first real account-level wallet management operations