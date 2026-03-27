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
**Current Subphase:** **0.4.13 — Protected Wallet Detach Execution**

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

### 0.4.10 — User-Driven Wallet-Owned Account Merge Execution
wallet management now also supports an authenticated merge flow for wallet-owned accounts:

- `POST /auth/account/merge/wallet/challenge`
- `POST /auth/account/merge/wallet/verify`

This allows the current authenticated user to absorb another wallet-owned account only after the source wallet explicitly signs a merge challenge.

### 0.4.11 — Explicit Primary-Wallet Switching
wallet management now also supports an authenticated primary-wallet switch flow:

- `POST /auth/wallets/primary`

This allows the current authenticated user to explicitly select which owned wallet is primary without changing ownership.

### 0.4.12 — Wallet Detach Eligibility Contract
wallet management now also supports an authenticated detach-eligibility evaluation flow:

- `POST /auth/wallets/detach/check`

This allows the current authenticated user to ask the backend whether one owned wallet is currently safe to detach, without changing ownership and without executing unlink behavior.

### 0.4.13 — Protected Wallet Detach Execution
wallet management now also supports an authenticated detach execution flow for already eligible owned wallets:

- `POST /auth/wallets/detach`

This allows the current authenticated user to detach one owned non-primary wallet only when the ownership guardrails introduced in 0.4.12 are satisfied.

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
- authenticated wallet-owned account merge challenges

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

### Authenticated Wallet-Owned Account Merge

#### `POST /auth/account/merge/wallet/challenge`

Creates an account-merge challenge bound to the currently authenticated user.

Request:

```json
{
  "address": "0x...",
  "chain": "scavium"
}
```

Behavior:

- requires valid JWT
- challenge purpose becomes `account_merge`
- challenge stores `requested_by_user_id`

---

#### `POST /auth/account/merge/wallet/verify`

Verifies the merge signature and consolidates all wallets from the source wallet-owned account into the current user.

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
  "merged_wallet": {
    "id": "...",
    "address": "0x...",
    "user_id": "...",
    "linked_at": "...",
    "is_primary": false
  },
  "source_user_id": "u_wallet_...",
  "target_user_id": "u_current_user",
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

JWT tokens include:

- `user_id`
- `wallet_id`
- `wallet_address`
- `auth_method`
- `exp`
- `iat`
- `nbf`

Wallet linking and wallet-owned account merge do **not** mint a new token. Both operate under the existing authenticated session.

---

## 🧪 Testing

Run:

```bash
go test ./...
```

Focus areas added in 0.4.13:

- authenticated detach-eligibility endpoint
- explicit detach rejection reasons for unsafe ownership states
- conservative protection against primary-wallet detach
- conservative protection against detaching the last owned wallet
- existing link, merge, and primary-switch coverage preserved

---

## 🚧 What 0.4.13 Solves

- authenticated user-driven wallet linking
- authenticated wallet-owned account merge execution
- protected primary-wallet switching under an authenticated user session
- authenticated wallet detach-eligibility evaluation under an authenticated user session
- challenge purpose separation between login, linking, and merge
- challenge-to-user binding through `requested_by_user_id`
- protected secondary-wallet attachment
- protected wallet-signed ownership consolidation
- deterministic single-primary wallet reassignment
- explicit detach rejection reasons for ownership-unsafe states

---

## ❌ What 0.4.13 Does Not Solve Yet

- wallet unlink API
- arbitrary cross-user ownership transfer outside wallet-signed merge
- merge between wallet-backed and other auth methods
- refresh tokens
- token revocation
- persistent authenticated sessions
- archival or alias records for merged source users

---

## 🧭 Next Phase

### 0.4.14 — Detach Follow-Up Semantics and Source Identity Lifecycle

Next expected focus:

- define whether future detach flows should bootstrap fresh wallet-only users automatically
- evaluate history, audit, or archival semantics for detached wallet identities
- preserve ownership invariants while extending unlink lifecycle semantics

---

## 🧩 Summary

At the end of Phase 0.4.13:

- wallet authentication remains stable
- identity remains unified
- ownership remains protected
- authenticated wallet linking is available
- wallet-owned account merge execution is available
- explicit primary-wallet switching is available
- wallet detach eligibility is available under authenticated control
- wallet detach execution is available for already eligible owned wallets
- the backend is ready for controlled detach execution design


