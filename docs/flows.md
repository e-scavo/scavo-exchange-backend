# Flows

## 🧠 Overview

This document describes the operational flows of authentication, identity resolution, wallet ownership, and authenticated wallet linking within the SCAVO Exchange Backend.

The flow evolution is:

- stateless wallet login
- → persistent wallet identity
- → unified durable user identity
- → multi-wallet ownership model
- → authenticated wallet-link contract

---

## 🔐 Wallet Authentication Flow

### Description

Primary authentication mechanism using EVM-compatible wallets.

### Flow

1. client requests challenge  
   `POST /auth/wallet/challenge`

2. backend:
   - normalizes wallet address
   - generates unique challenge
   - persists challenge with `purpose = auth_bootstrap`

3. client signs challenge message

4. client sends verification request  
   `POST /auth/wallet/verify`

5. backend:
   - validates challenge existence
   - checks expiration
   - checks unused state
   - verifies signature
   - recovers wallet address
   - ensures address matches

6. challenge is marked as used

7. wallet identity is resolved:
   - loaded if exists
   - created if missing

8. durable user is resolved:
   - loaded if linked
   - created if missing

9. ownership is enforced:
   - wallet linked if previously unowned
   - ownership conflict rejected

10. JWT is issued with unified identity

---

## 🧩 Identity Resolution Flow

### Description

Defines how the system transitions from wallet identity to durable user identity.

### Flow

1. wallet identity is loaded
2. backend checks `user_id`
3. if `user_id` exists:
   - load durable user
4. if `user_id` is missing:
   - create durable wallet-backed user
   - attach wallet identity to user
5. return unified identity

---

## 🏷️ Wallet Ownership Flow (0.4.8)

### Description

Defines how wallet ownership is persisted and enforced.

### Flow

1. wallet identity exists or is created
2. system evaluates ownership:
   - if `user_id` empty:
     - assign user
     - set `linked_at`
     - set `is_primary = true` for first wallet
   - if `user_id` exists:
     - validate ownership
     - reject mismatch
3. ownership metadata persists independently of JWT session

---

## 🔄 Ownership Enforcement Flow

### Description

Prevents invalid ownership transitions.

### Flow

1. wallet identity retrieved
2. incoming user compared with persisted `user_id`
3. if persisted owner differs:
   - reject operation
   - return ownership conflict
4. if same owner:
   - allow read operations
   - prevent unsafe duplicate-link semantics where required

---

## 📦 Authenticated Wallet Inventory Flow

### Description

Returns all wallets linked to the authenticated durable user.

### Flow

1. client sends  
   `GET /auth/wallets`
2. backend validates JWT
3. backend extracts `user_id`
4. backend loads all wallet identities for that user
5. backend orders results:
   - primary first
   - then by `linked_at`
   - then by address
6. backend returns wallet list

---

## 🔗 Authenticated Wallet Linking Flow (0.4.9)

### Description

Allows an already authenticated user to attach a new secondary wallet to the current durable account.

### Step A — Link challenge creation

1. client sends  
   `POST /auth/wallets/link/challenge`

2. request includes:
   - target wallet address
   - chain

3. backend validates JWT and extracts current `user_id`

4. backend creates challenge with:
   - `purpose = wallet_link`
   - `requested_by_user_id = current user`

5. backend returns challenge for signature

---

### Step B — Link verification

1. client signs the link challenge with the target secondary wallet
2. client sends  
   `POST /auth/wallets/link/verify`

3. backend validates:
   - challenge existence
   - unused state
   - expiration
   - challenge purpose
   - challenge user binding
   - signature correctness
   - wallet address correctness

4. backend resolves wallet identity

5. backend enforces ownership:
   - reject if wallet belongs to another user
   - reject if wallet already linked to current user
   - attach as secondary if unowned

6. backend consumes challenge
7. backend returns:
   - linked wallet
   - updated wallet inventory

---

## 🔄 Session Flow

### Description

Defines how authenticated sessions are resolved.

### Flow

1. client sends request with JWT
2. backend validates token
3. backend parses claims
4. backend resolves durable user
5. backend attaches wallet context from claims
6. backend returns session through:
   - `/auth/me`
   - `/auth/session`

---

## ⚙️ Error Handling Flow

### Wallet Challenge
- invalid address → `invalid_wallet_address`
- expired challenge → `wallet_challenge_expired`
- reused challenge → `wallet_challenge_used`

### Wallet Verification
- invalid signature → `invalid_wallet_signature`
- challenge not found → `wallet_challenge_not_found`

### Wallet Linking
- wrong challenge purpose → `wallet_challenge_purpose_mismatch`
- challenge belongs to different authenticated user → `wallet_link_challenge_user_mismatch`
- wallet already linked elsewhere → `wallet_identity_already_linked`
- wallet already linked to current user → `wallet_identity_already_linked_to_user`

---

## 🧭 Future Flow Extensions (Post 0.4.9)

### Planned in 0.4.10
- wallet unlink flow
- primary-wallet switching flow
- deeper wallet-management contracts

### Later phases
- multi-auth merge flow
- account consolidation flow
- recovery flow

---

## 🧩 Summary

At the end of Phase 0.4.9:

- authentication is stable
- identity is unified
- ownership is enforced
- user-driven wallet linking is implemented under authenticated control

The backend now transitions from:

**authentication flows → ownership flows → authenticated wallet-management flows**