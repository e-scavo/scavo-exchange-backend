# Architecture

## 🧠 Overview

The SCAVO Exchange Backend is designed around a **wallet-first identity architecture**, progressively evolving toward a **unified account model** capable of supporting exchange-grade features.

The system separates:

- Authentication mechanism (wallet / future methods)
- Identity representation (wallet identity vs user)
- Ownership model (introduced in 0.4.8)

---

## 🧩 Core Layers

### 1. Transport Layer

- HTTP API
- JSON-based communication
- Stateless request handling

---

### 2. Auth Layer

Located in:

- `internal/modules/auth`

Responsibilities:

- wallet challenge generation
- wallet signature verification
- JWT issuance
- identity resolution
- ownership enforcement (since 0.4.8)

---

### 3. User Layer

Located in:

- `internal/modules/user`

Responsibilities:

- durable user creation
- identity abstraction
- future multi-auth support

---

### 4. Persistence Layer

- PostgreSQL (primary)
- In-memory fallback (dev/testing)

---

## 🔐 Identity Model

### Pre 0.4.6

- Identity was session-based
- No persistence

---

### 0.4.6 — Wallet Identity Persistence

- Wallet identity stored in:
  - `auth_wallet_identities`
- Address becomes stable identifier

---

### 0.4.7 — Unified Identity Model

- Wallet identity linked to durable user
- `user_id` introduced
- JWT identity unified

---

### 0.4.8 — Ownership Model Introduction

Wallet identity evolves into a **first-class ownership entity**.

Each wallet identity includes:

- `id`
- `address`
- `user_id`
- `linked_at`
- `is_primary`

---

## 🏷️ Ownership Model

### Core Rules

1. A wallet belongs to exactly one user
2. A user can own multiple wallets
3. Only one wallet per user can be primary
4. Wallet ownership cannot be reassigned across users

---

### Ownership Metadata

| Field       | Description |
|------------|------------|
| user_id    | Owner user |
| linked_at  | Timestamp of ownership |
| is_primary | Primary wallet flag |

---

### Ownership Semantics

- Ownership is persisted at DB level
- Ownership is independent from sessions
- Ownership is enforced in store layer
- Ownership conflicts are rejected

---

## 🔄 Authentication Flow (Wallet)

1. Client requests challenge
2. Server creates challenge
3. Client signs message
4. Server verifies signature
5. Challenge is consumed
6. Wallet identity is resolved
7. User is resolved or created
8. Ownership is enforced
9. JWT is issued

---

## 🔄 Ownership Resolution Flow

When a wallet logs in:

1. Retrieve wallet identity
2. If not linked → create or link user
3. If linked:
   - validate ownership
   - load user
4. Ensure primary wallet semantics
5. Return unified identity

---

## 🔌 API Layer

### Auth endpoints

- `/auth/login`
- `/auth/wallet/challenge`
- `/auth/wallet/verify`
- `/auth/me`
- `/auth/session`
- `/auth/wallets` ← introduced in 0.4.8

---

## 🧾 JWT Design

JWT tokens are:

- Stateless
- Short-lived
- Self-contained

### Claims include:

- `user_id`
- `wallet_id`
- `wallet_address`
- `auth_method`

---

## 🗄️ Data Model

### `auth_wallet_challenges`

- challenge lifecycle
- expiration
- single-use enforcement

---

### `auth_wallet_identities`

Stores wallet registry and ownership:

- `id`
- `address`
- `user_id`
- `linked_at`
- `is_primary`

---

### `users`

- durable identity
- auth-provider agnostic
- future extensibility

---

## ⚙️ Design Decisions

### Wallet-first approach

Chosen because:

- aligns with crypto-native UX
- avoids early complexity of email/password systems
- simplifies initial identity layer

---

### Separation of identity and ownership

- wallet identity ≠ user
- ownership is explicit, not implicit
- allows future merging strategies

---

### Incremental evolution

Each subphase introduces:

- one structural improvement
- backward compatibility
- minimal breakage risk

---

## ⚠️ Constraints

- No wallet reassignment allowed
- No multi-auth merge yet
- No unlink operations
- No ownership transfer

---

## 🚧 Future Evolution (Post 0.4.8)

The system is now prepared for:

### 0.4.9

- wallet linking API
- ownership validation flows
- controlled linking operations

---

### Later phases

- account consolidation
- multi-auth merging (wallet + email)
- recovery flows
- compliance-ready identity

---

## 🧩 Summary

At the end of 0.4.8:

- identity is durable
- ownership is modeled
- wallet authentication is stable
- system is structurally ready for expansion

The architecture successfully transitions from:

**stateless auth → persistent identity → unified user → multi-wallet ownership → account-level foundation**