# Flows

## 🧠 Overview

This document describes the operational flows of authentication, identity resolution, and wallet ownership within the SCAVO Exchange Backend.

Flows are designed to reflect the evolution from:

- stateless wallet login
- → persistent wallet identity
- → unified user identity
- → multi-wallet ownership model

---

## 🔐 Wallet Authentication Flow

### Description

This is the primary authentication mechanism using EVM-compatible wallets.

---

### Flow

1. Client requests a challenge:

   `POST /auth/wallet/challenge`

2. Backend:

   - normalizes wallet address
   - generates a unique challenge
   - stores it with expiration

3. Client signs the challenge message

4. Client sends verification request:

   `POST /auth/wallet/verify`

5. Backend:

   - validates challenge existence
   - checks expiration
   - verifies signature
   - recovers wallet address
   - ensures address matches

6. Challenge is marked as used (one-time use)

7. Wallet identity is resolved:

   - retrieved if exists
   - created if not

8. User is resolved:

   - retrieved if linked
   - created if not

9. Ownership is enforced:

   - wallet is linked to user if not already linked
   - reassignment is rejected

10. Primary wallet semantics are applied:

   - first wallet becomes primary
   - additional wallets remain non-primary (future expansion)

11. JWT is issued with unified identity

---

## 🧩 Identity Resolution Flow

### Description

Defines how the system transitions from wallet identity to durable user identity.

---

### Flow

1. Wallet identity is loaded
2. Check `user_id`:

   - if present:
     - load user
   - if missing:
     - create user
     - link identity → user

3. Ensure ownership consistency:

   - wallet cannot belong to multiple users

4. Return unified identity:

   - user + wallet context

---

## 🏷️ Wallet Ownership Flow (0.4.8)

### Description

Defines how wallet ownership is managed and enforced.

---

### Flow

1. Wallet identity exists or is created
2. System evaluates ownership:

   - if `user_id` is empty:
     - assign to user
     - set `linked_at`
     - set `is_primary = true`

   - if `user_id` exists:
     - verify ownership
     - reject if mismatch

3. Ownership metadata is persisted:

   - `user_id`
   - `linked_at`
   - `is_primary`

---

## 🔄 Ownership Enforcement Flow

### Description

Prevents invalid ownership transitions.

---

### Flow

1. Wallet identity retrieved
2. Incoming user ID compared against existing `user_id`
3. If mismatch:

   - reject operation
   - return `ErrWalletIdentityAlreadyLinked`

4. If match:

   - allow operation

---

## 📦 Authenticated Wallet Inventory Flow

### Description

Allows clients to retrieve all wallets linked to the authenticated user.

---

### Flow

1. Client sends request:

   `GET /auth/wallets`

2. Backend:

   - extracts JWT from request
   - validates token
   - extracts `user_id` from claims

3. Backend queries wallet identities:

   - filter by `user_id`
   - order:
     - primary wallet first
     - then by `linked_at`
     - then by address

4. Backend returns:

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

## 🔄 Session Flow

### Description

Defines how authenticated sessions are resolved.

---

### Flow

1. Client sends request with JWT:

   `Authorization: Bearer <token>`

2. Backend:

   - validates token
   - parses claims

3. Backend resolves identity:

   - loads user
   - attaches wallet context

4. Returns session:

   `/auth/me` or `/auth/session`

---

## ⚙️ Error Handling Flow

### Wallet Challenge

- invalid address → `invalid_wallet_address`
- expired challenge → `wallet_challenge_expired`
- reused challenge → `wallet_challenge_used`

---

### Wallet Verification

- invalid signature → `invalid_wallet_signature`
- challenge not found → `wallet_challenge_not_found`

---

### Ownership

- wallet already linked → `wallet_identity_already_linked`

---

## 🧭 Future Flow Extensions (Post 0.4.8)

### Planned in 0.4.9

- user-driven wallet linking flow
- ownership confirmation flows
- conflict resolution flows

---

### Later phases

- wallet unlink flow
- multi-auth merge flow
- account consolidation flow
- recovery flow

---

## 🧩 Summary

At the end of Phase 0.4.8:

- authentication is stable
- identity is unified
- ownership is enforced
- multi-wallet support is structurally enabled

The system is ready to transition from:

**authentication flows → ownership flows → account-level flows**