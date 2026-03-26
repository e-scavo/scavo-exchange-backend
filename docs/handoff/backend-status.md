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
**Latest Completed Subphase:** 0.4.8 — Account Consolidation and Multi-Wallet Ownership Foundations  

---

## 🔐 Authentication Status

### Implemented

- Password-based authentication (dev only)
- Wallet-based authentication (EVM)

Wallet flow:

- challenge generation
- signature verification
- challenge consumption (one-time use)
- wallet identity resolution
- user resolution / creation
- ownership enforcement
- JWT issuance

---

## 🧩 Identity Model Status

The system now uses a **unified identity model**:

- wallet identities are persisted
- each wallet is linked to a durable user
- JWT reflects unified identity

---

## 🏷️ Ownership Model Status (0.4.8)

Ownership is now a **first-class concept**.

### Capabilities

- one user can own multiple wallets
- wallet ownership is persisted
- ownership metadata:
  - `user_id`
  - `linked_at`
  - `is_primary`
- primary wallet designation enforced
- ownership reassignment is blocked

---

## 🗄️ Persistence Status

### Tables

#### `auth_wallet_challenges`
- persistent
- expiration enforced
- single-use enforced

---

#### `auth_wallet_identities`
- persistent wallet registry
- ownership metadata included

---

#### `users`
- durable user representation
- wallet-backed users supported

---

## 🔌 API Status

### Auth endpoints

- `POST /auth/login`
- `POST /auth/wallet/challenge`
- `POST /auth/wallet/verify`
- `GET /auth/me`
- `GET /auth/session`
- `GET /auth/wallets` ← introduced in 0.4.8

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

---

## 🧪 Testing Status

Validated through:

- `go test ./...`
- manual API testing
- SQL verification queries

Coverage includes:

- wallet auth flow
- identity linking
- ownership enforcement
- replay protection

---

## ⚠️ Known Limitations

The system intentionally does NOT yet support:

- user-driven wallet linking API
- wallet unlink operations
- cross-user wallet transfer
- auth method merging (wallet + email)
- refresh tokens
- token revocation
- persistent sessions

---

## 🧭 Next Phase

### 0.4.9 — User-Driven Wallet Linking Contract and Protected Account Merge Preparation

This phase will introduce:

- wallet linking endpoints
- ownership validation flows
- controlled linking operations
- merge-safe identity model
- preparation for account consolidation

---

## 🧩 Long-Term Direction

Future evolution includes:

- multi-auth identity system
- account consolidation
- recovery flows
- compliance-ready identity layer

---

## 📌 Development Guidelines

When continuing development:

- do not break ownership invariants
- do not allow wallet reassignment
- preserve backward compatibility
- evolve identity model incrementally
- maintain documentation alignment with implementation

---

## 🧾 Summary

At the end of Phase 0.4.8:

- authentication is stable
- identity is unified
- ownership is implemented
- multi-wallet support is structurally enabled

The backend is now ready to transition into:

➡️ **controlled identity expansion and account-level features**