# Architecture

## 🧠 Overview

The SCAVO Exchange Backend is designed around a **wallet-first identity architecture** that progressively evolves into a **durable account model** suitable for exchange-grade ownership and future multi-auth identity operations.

The architecture intentionally separates:

- authentication mechanism
- wallet identity representation
- durable platform user abstraction
- persisted wallet ownership
- authenticated wallet-management contracts

---

## 🧩 Core Layers

### 1. Transport Layer
- HTTP API
- JSON-based communication
- stateless request handling

### 2. Auth Layer
Located in:

- `internal/modules/auth`

Responsibilities:

- wallet challenge generation
- wallet signature verification
- JWT issuance
- identity resolution
- ownership enforcement
- authenticated wallet-linking flows

### 3. User Layer
Located in:

- `internal/modules/user`

Responsibilities:

- durable user creation
- durable user resolution
- future auth-provider expansion

### 4. Persistence Layer
- PostgreSQL (primary)
- in-memory fallback (dev/testing)

---

## 🔐 Identity Model Evolution

### Pre 0.4.6
- identity was session-oriented
- wallet state was not durable

### 0.4.6 — Wallet Identity Persistence
- wallet identity stored in `auth_wallet_identities`
- wallet address becomes a stable registry entry

### 0.4.7 — Unified Identity Model
- wallet identity linked to durable user
- `user_id` introduced
- JWT identity unified around durable platform user

### 0.4.8 — Ownership Model Introduction
wallet identity becomes a first-class ownership entity with:

- `id`
- `address`
- `user_id`
- `linked_at`
- `is_primary`

### 0.4.9 — Authenticated Wallet Linking Contract
the architecture adds a dedicated authenticated linking flow, still based on challenge + signature verification, but now explicitly bound to the current authenticated user.

This is the first backend-managed wallet operation that acts on ownership under an existing session rather than during initial login bootstrap.

---

## 🏷️ Ownership Model

### Core Rules

1. a wallet belongs to exactly one user
2. a user can own multiple wallets
3. only one wallet per user can be primary
4. wallet ownership cannot be reassigned across users
5. authenticated wallet linking adds secondary wallets only
6. 0.4.9 does not switch primary ownership

---

## 🏷️ Ownership Metadata

| Field | Description |
|------|-------------|
| `user_id` | owning durable user |
| `linked_at` | ownership creation timestamp |
| `is_primary` | primary-wallet flag |

---

## 🧾 Challenge Model

### Pre 0.4.9
wallet challenge was effectively used only for authentication bootstrap.

### 0.4.9
wallet challenges now include:

- `purpose`
- `requested_by_user_id`

### Challenge purposes

- `auth_bootstrap`
- `wallet_link`

This avoids reusing the same challenge semantics blindly across two very different operations.

---

## 🔄 Authentication Flow (Wallet Login)

1. client requests login challenge
2. backend persists challenge with `auth_bootstrap` purpose
3. client signs message
4. backend verifies signature
5. challenge is consumed
6. wallet identity is resolved
7. durable user is resolved or created
8. ownership is enforced
9. JWT is issued

---

## 🔄 Authenticated Wallet Linking Flow

1. user already holds valid JWT
2. client requests link challenge:
   - `POST /auth/wallets/link/challenge`
3. backend persists challenge with:
   - `purpose = wallet_link`
   - `requested_by_user_id = current user`
4. user signs with the secondary wallet
5. client submits:
   - `POST /auth/wallets/link/verify`
6. backend validates:
   - challenge existence
   - challenge freshness
   - purpose correctness
   - requesting user correctness
   - signature correctness
7. backend resolves wallet identity
8. backend rejects ownership conflict if wallet belongs elsewhere
9. backend attaches wallet as non-primary
10. backend returns updated wallet inventory

---

## 🔌 API Layer

### Auth endpoints
- `/auth/login`
- `/auth/wallet/challenge`
- `/auth/wallet/verify`
- `/auth/me`
- `/auth/session`

### Wallet ownership endpoints
- `/auth/wallets`
- `/auth/wallets/link/challenge`
- `/auth/wallets/link/verify`

---

## 🧾 JWT Design

JWT tokens are:

- stateless
- short-lived
- self-contained

Claims include:

- `user_id`
- `wallet_id`
- `wallet_address`
- `auth_method`

Wallet linking does not mint a fresh token because it operates under an already authenticated durable session.

---

## 🗄️ Data Model

### `auth_wallet_challenges`
stores challenge lifecycle and, from 0.4.9 onward, also stores operation metadata:

- `purpose`
- `requested_by_user_id`
- issued / expires / used lifecycle

### `auth_wallet_identities`
stores wallet registry and ownership:

- `id`
- `address`
- `user_id`
- `linked_at`
- `is_primary`

### `users`
stores durable user abstraction:

- login-independent identity
- wallet-backed users now
- future auth-provider aggregation later

---

## ⚙️ Design Decisions

### Wallet-first approach
chosen because it aligns with crypto-native UX and reduces early auth-surface complexity.

### Separation of identity and ownership
wallet identity is not the same as durable user identity.
Ownership is explicit rather than inferred.

### Challenge-purpose separation
0.4.9 extends the challenge system rather than introducing a second parallel challenge subsystem, but still keeps semantic separation through `purpose`.

### Incremental evolution
each subphase introduces one structural improvement while preserving previous behavior.

---

## ⚠️ Constraints

Still intentionally not supported:

- wallet unlink
- primary-wallet switching
- cross-user wallet transfer
- automatic merge execution
- multi-auth merge resolution

---

## 🚧 Future Evolution (Post 0.4.9)

### 0.4.10
- wallet detach / unlink rules
- primary-wallet switching contract
- deeper merge-safe identity progression

### Later phases
- account consolidation
- multi-auth identity merging
- recovery flows
- compliance-ready identity expansion

---

## 🧩 Summary

At the end of 0.4.9:

- wallet authentication is stable
- durable identity is stable
- wallet ownership is stable
- authenticated wallet linking is implemented
- the backend is structurally ready to move from ownership persistence into true account-level wallet management