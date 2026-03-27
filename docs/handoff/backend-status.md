# Backend Status — SCAVO Exchange

## 🧠 Overview

This document represents the **current operational and architectural state** of the SCAVO Exchange Backend.

It is intended to:

- provide continuity across development phases
- allow safe context transfer between sessions
- serve as the single source of truth for backend status

---

## 📌 Current State

**Stage:** 0 — Foundation  
**Phase:** 0.4 — Auth and User Stabilization  
**Latest Completed Subphase:** 0.4.13 — Protected Wallet Detach Execution

---

## 🔐 Authentication Status

### Implemented

- password-based authentication (dev only)
- wallet-based authentication (EVM)

Wallet login flow supports:

- challenge generation
- signature verification
- challenge consumption (one-time use)
- wallet identity resolution
- durable user resolution / creation
- ownership enforcement
- JWT issuance

---

## 🧩 Identity Model Status

The system uses a **unified durable identity model**:

- wallet identities are persisted
- each wallet is linked to a durable user
- JWT reflects unified identity through `user_id`
- wallet metadata remains available for wallet-authenticated sessions

---

## 🏷️ Ownership Model Status

Ownership is a first-class persisted concept.

### Capabilities

- one user can own multiple wallets
- wallet ownership is persisted
- ownership metadata includes:
  - `user_id`
  - `linked_at`
  - `is_primary`
- primary-wallet uniqueness enforced
- wallet reassignment across users blocked

---

## 🔗 Wallet Ownership Status (0.4.13)

The backend now supports authenticated wallet-linking, authenticated wallet-owned account merge execution, authenticated primary-wallet switching, authenticated wallet detach-eligibility evaluation, and authenticated detach execution for already eligible owned wallets.

### Capabilities

- authenticated user can request wallet-link challenge
- challenge is persisted with:
  - `purpose = wallet_link`
  - `requested_by_user_id`
- authenticated user can verify link signature
- secondary wallet attaches to current user
- authenticated user can execute wallet-owned account merge after source-wallet signature
- authenticated user can explicitly switch the current primary wallet
- authenticated user can request detach-eligibility evaluation for one owned wallet
- authenticated user can execute detach for one already eligible owned wallet
- explicit detach rejection reasons are returned when detach is not yet eligible
- updated wallet inventory is returned after successful linking, merge, primary switching, and detach execution

### Protections

- link challenge must belong to current authenticated user
- wrong-purpose challenge is rejected
- wallet already owned by another user is rejected
- wallet already linked to current user is rejected
- successful link does not issue a new JWT and does not implicitly merge accounts

---

## 🗄️ Persistence Status

### `auth_wallet_challenges`
- persistent
- expiration enforced
- single-use enforced
- now supports:
  - `purpose`
  - `requested_by_user_id`

### `auth_wallet_identities`
- persistent wallet registry
- ownership metadata included
- multi-wallet ownership supported

### `users`
- durable user representation
- wallet-backed users supported
- prepared for later multi-auth evolution

---

## 🔌 API Status

### Auth endpoints
- `POST /auth/login`
- `POST /auth/wallet/challenge`
- `POST /auth/wallet/verify`
- `GET /auth/me`
- `GET /auth/session`

### Ownership endpoints
- `GET /auth/wallets`

### Wallet-link endpoints
- `POST /auth/wallets/link/challenge`
- `POST /auth/wallets/link/verify`
- `POST /auth/account/merge/wallet/challenge`
- `POST /auth/account/merge/wallet/verify`

---

## 🧾 JWT Status

Tokens are:

- stateless
- short-lived
- self-contained

Claims include:

- `user_id`
- `wallet_id`
- `wallet_address`
- `auth_method`

Wallet linking uses the already authenticated JWT and does not mint a replacement token.

---

## ⚙️ Behavioral Guarantees

The backend guarantees:

- wallet challenges are one-time use
- wallet signatures are verified deterministically
- wallet identities are persistent
- user identity is durable
- ownership is consistent and enforced
- wallets cannot be reassigned across users
- primary wallet uniqueness is maintained
- link challenges are user-bound
- link and login challenge purposes are not interchangeable

---

## 🧪 Testing Status

Validated at the design and code level through:

- `go test ./...`
- manual API testing procedures
- SQL verification queries

Coverage now includes:

- wallet auth flow
- identity linking
- ownership enforcement
- replay protection
- wallet-link challenge flow
- wallet-link verification flow
- ownership conflict rejection during link operations
- protected primary-wallet switching

---

## ⚠️ Known Limitations

The system intentionally does **not** yet support:

- wallet unlink operations
- cross-user wallet transfer
- merge between wallet identities and future auth methods
- refresh tokens
- token revocation
- persistent authenticated sessions

---

## 🧭 Next Phase

### 0.4.14 — Detach Follow-Up Semantics and Source Identity Lifecycle

Expected next focus:

- define whether detached wallets should remain free identities only or bootstrap fresh wallet-only accounts later
- evaluate audit and lifecycle markers for detached identities
- preserve conservative ownership invariants while enriching detach semantics

---

## 📌 Development Guidelines

When continuing development:

- do not break ownership invariants
- do not allow wallet reassignment
- do not bypass challenge-purpose checks
- preserve backward compatibility of wallet login
- keep challenge-to-user binding explicit in wallet-management flows
- maintain documentation alignment with implementation

---

## 🧾 Summary

At the end of Phase 0.4.13:

- authentication is stable
- identity is unified
- ownership is implemented and protected
- authenticated wallet linking is implemented
- authenticated wallet-owned account merge is implemented
- explicit primary-wallet switching is implemented
- wallet detach eligibility is implemented
- wallet detach execution is implemented for already eligible owned wallets
- the backend is ready to move into controlled detach execution design
