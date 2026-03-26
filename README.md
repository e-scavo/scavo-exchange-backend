# SCAVO Exchange — Backend

## 🧠 Overview

SCAVO Exchange Backend is a modular Go-based service that provides:

- JWT-based authentication
- EVM wallet authentication
- WebSocket real-time communication
- PostgreSQL-backed persistence
- Redis integration support
- Modular services for auth, user, system, and infrastructure concerns

---

## 🏗️ Current Stage

**Stage 0 — Foundation**  
**Phase 0.4 — Auth and User Stabilization**  
**Current completed subphase: 0.4.7 — Wallet ↔ User Linking and Unified Identity Model**

---

## 🔐 Authentication

### Supported Methods

| Method | Description |
|--------|-------------|
| `password_dev` | Development-only email/password login |
| `wallet_evm` | Wallet signature-based login using EVM-compatible signing |

---

## 🔗 Wallet Authentication Flow

The backend now supports a wallet login flow that resolves both a durable wallet identity and a durable platform user.

### Flow Summary

1. Client requests a wallet challenge
2. Backend generates a one-time challenge message
3. Client signs the challenge using an EVM wallet
4. Backend verifies the signature
5. Backend consumes the challenge
6. Backend resolves or creates a persistent wallet identity
7. Backend resolves or creates a linked platform user
8. Backend issues a JWT access token with unified identity metadata

---

## 🚀 HTTP Endpoints

### Auth

#### `POST /auth/login`
Development login endpoint.

#### `POST /auth/wallet/challenge`
Creates a wallet signing challenge.

#### `POST /auth/wallet/verify`
Verifies a wallet signature, links the wallet identity to a durable user, and returns a JWT token.

#### `GET /auth/me`
Returns the authenticated user context.

#### `GET /auth/session`
Returns normalized session metadata derived from JWT claims.

---

## 🔌 WebSocket

### Entry Point

#### `GET /ws`

The WebSocket layer supports authenticated sessions through the same JWT token used by the HTTP API.

When a wallet-authenticated token is provided, the session includes:

- `user_id`
- `email`
- `wallet_id`
- `wallet_address`
- `auth_method`
- `chain`
- `subject`
- `issuer`
- `expires_at`

---

## 🗄️ Persistence

### Wallet Challenges

Wallet challenges are durably stored in PostgreSQL.

#### Table
`auth_wallet_challenges`

#### Stored fields

- `id`
- `address`
- `chain`
- `nonce`
- `message`
- `issued_at`
- `expires_at`
- `used_at`
- `created_at`

### Wallet Identities

Wallet identities are persisted and uniquely resolved per wallet address.

#### Table
`auth_wallet_identities`

#### Stored fields

- `id`
- `address`
- `user_id`
- `created_at`

### Users

Wallet-authenticated sessions are now backed by the same `users` table already used by the development login flow.

Wallet-backed users are provisioned automatically with:

- deterministic internal user id
- synthetic internal email
- wallet address as display name
- updated `last_login_at` on successful login

---

## 🔑 JWT Claims

Wallet-authenticated sessions may include the following claims:

| Claim | Description |
|-------|-------------|
| `uid` | Internal user identifier |
| `email` | Linked platform user email |
| `wallet_id` | Persistent wallet identity ID |
| `wallet_address` | EVM wallet address |
| `auth_method` | Authentication method (`wallet_evm`, `password_dev`) |
| `chain` | Logical blockchain identifier |

---

## 🧪 Testing

Run the full backend test suite with:

```bash
go test ./...
```

Run the wallet authentication flow manually with:

1. request a wallet challenge
2. sign the returned message using a test wallet
3. submit the signature to `POST /auth/wallet/verify`
4. inspect `auth_wallet_identities.user_id` and `users` to confirm the link

---

## ⚙️ Runtime Requirements

- Go 1.25+
- PostgreSQL recommended for durable auth persistence
- Redis optional
- Compatible EVM wallet/client for wallet-auth testing

---

## 📦 Current Status

| Stage | Phase | Subphase | Status |
|------|-------|----------|--------|
| 0 | 0.4 | 0.4.7 | ✅ Completed |

---

## 🚧 What 0.4.7 Solves

- Durable wallet challenge storage
- Atomic challenge consumption
- Persistent wallet identity creation
- Wallet-to-user linking through durable persistence
- Unified user resolution for `GET /auth/me` and `GET /auth/session`
- Wallet metadata propagation into JWT/session layers
- In-memory fallback for non-DB development environments

---

## ❌ What 0.4.7 Does Not Solve Yet

- User-managed wallet linking API
- Multi-wallet account aggregation
- Refresh tokens
- Token revocation
- Persistent session management
- Account merge flows across auth methods

---

## ⏭️ Next Planned Phase

**Phase 0.4.8 — Account Consolidation and Multi-Wallet Ownership Foundations**
