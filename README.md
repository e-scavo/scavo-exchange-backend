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
**Current completed subphase: 0.4.6 — Wallet Identity Persistence and Durable Challenge Storage**

---

## 🔐 Authentication

### Supported Methods

| Method        | Description |
|---------------|-------------|
| `password_dev` | Development-only email/password login |
| `wallet_evm`   | Wallet signature-based login using EVM-compatible signing |

---

## 🔗 Wallet Authentication Flow

The backend currently supports a wallet login flow based on challenge issuance and signature verification.

### Flow Summary

1. Client requests a wallet challenge
2. Backend generates a one-time challenge message
3. Client signs the challenge using an EVM wallet
4. Backend verifies the signature
5. Backend consumes the challenge
6. Backend resolves or creates a persistent wallet identity
7. Backend issues a JWT access token

---

## 🚀 HTTP Endpoints

### Auth

#### `POST /auth/login`
Development login endpoint.

#### `POST /auth/wallet/challenge`
Creates a wallet signing challenge.

#### `POST /auth/wallet/verify`
Verifies a wallet signature and returns a JWT token.

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

Wallet challenges are now durably stored in PostgreSQL.

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

Wallet identities are now persisted and uniquely resolved per wallet address.

#### Table
`auth_wallet_identities`

#### Stored fields

- `id`
- `address`
- `created_at`

---

## 🔑 JWT Claims

Wallet-authenticated sessions may include the following claims:

| Claim | Description |
|-------|-------------|
| `uid` | Internal user identifier |
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
| 0 | 0.4 | 0.4.6 | ✅ Completed |

---

## 🚧 What 0.4.6 Solves

- Durable wallet challenge storage
- Atomic challenge consumption
- Persistent wallet identity creation
- Wallet metadata propagation into JWT/session layers
- In-memory fallback for non-DB development environments

---

## ❌ What 0.4.6 Does Not Solve Yet

- Wallet ↔ user linking
- Multi-wallet account model
- Refresh tokens
- Token revocation
- Persistent session management

---

## ⏭️ Next Planned Phase

**Phase 0.4.7 — Wallet ↔ User Linking and Unified Identity Model**