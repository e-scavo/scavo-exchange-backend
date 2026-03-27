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

## 🔗 Authenticated Wallet-Owned Account Merge Flow (0.4.10)

### Description

Allows an already authenticated user to absorb a source wallet-owned account after the source wallet explicitly signs a merge challenge.

### Step A — Merge challenge creation

1. client sends  
   `POST /auth/account/merge/wallet/challenge`

2. request includes:
   - source wallet address
   - chain

3. backend validates JWT and extracts current `user_id`

4. backend creates challenge with:
   - `purpose = account_merge`
   - `requested_by_user_id = current user`

5. backend returns challenge for signature

---

### Step B — Merge verification

1. source wallet signs the merge challenge
2. client sends  
   `POST /auth/account/merge/wallet/verify`

3. backend validates:
   - challenge existence
   - unused state
   - expiration
   - challenge purpose
   - challenge user binding
   - signature correctness
   - source wallet ownership existence

4. backend resolves source wallet identity
5. backend derives source user from wallet ownership
6. backend atomically reassigns all source-user wallets to the authenticated target user
7. backend returns:
   - merged wallet
   - source user id
   - target user id
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

### Wallet-Owned Account Merge
- wrong challenge purpose → `wallet_challenge_purpose_mismatch`
- challenge belongs to different authenticated user → `wallet_account_merge_user_mismatch`
- source wallet not linked → `wallet_account_merge_source_not_linked`
- merge not required because wallet already belongs to current user → `wallet_account_merge_not_required`

### Wallet Detach Eligibility
- invalid address → `invalid_wallet_address`
- wallet missing → `wallet_identity_not_found`
- wallet belongs to another user → `wallet_identity_not_owned_by_user`
- detach currently blocked because wallet is primary → response `eligible = false` with `wallet_is_primary`
- detach currently blocked because user would have no wallets left → response `eligible = false` with `user_would_have_no_wallets`

---

## 9. Detached Wallet Reattachment

1. authenticated user previously detaches one eligible owned secondary wallet
2. detached wallet identity remains persisted by address and wallet identity ID
3. authenticated user later calls `POST /auth/wallets/link/challenge` for that same address
4. backend creates a normal protected wallet-link challenge bound to the authenticated user
5. wallet signs and verifies through `POST /auth/wallets/link/verify`
6. backend reattaches the detached wallet as an owned secondary wallet again

#### Safety Rules

- detached-wallet reattachment still requires an authenticated session
- detached-wallet reattachment still requires wallet signature proof
- detached-wallet reattachment does not bypass current ownership checks
- detached-wallet reattachment does not implicitly restore historical primary state

### 10. Detached Wallet Wallet-Login Rebound

1. detached wallet later initiates the standard wallet-login bootstrap flow
2. backend creates a normal auth challenge through `POST /auth/wallet/challenge`
3. wallet signs and verifies through `POST /auth/wallet/verify`
4. backend resolves or recreates the wallet-owned user identity for that wallet
5. backend reattaches the wallet under that wallet-owned durable user
6. wallet becomes primary in that wallet-owned identity scope

#### Safety Rules

- wallet-login rebound follows the same wallet-auth flow as any other wallet bootstrap
- detached-wallet rebound does not restore the previous detached owner automatically
- detached-wallet rebound does not create archival or audit metadata in the current phase

---

## 9. Detached Identity Audit Readiness

1. authenticated user detaches one already eligible wallet through `POST /auth/wallets/detach`
2. backend clears ownership fields from that wallet identity
3. backend stamps `detached_at` on that wallet identity
4. detached wallet remains reusable through future protected linking or wallet-login rebound
5. later reuse keeps `detached_at` as minimal lifecycle evidence that the identity has been detached before

#### Safety Rules

- `detached_at` never changes current ownership by itself
- `detached_at` never restores previous ownership
- detached lifecycle metadata is intentionally minimal and non-destructive
- current phase does not introduce detached-history tables, event sourcing, or archival semantics

## 🧭 Future Flow Extensions (Post 0.4.15)

### Potential later evolution
- detached-by-user audit metadata
- queryable detach history
- richer detached-wallet observability

### Later phases
- multi-auth merge flow
- account consolidation flow
- recovery flow

---

## 🧩 Summary

At the end of Phase 0.4.15:

- authentication is stable
- identity is unified
- ownership is enforced
- user-driven wallet linking is implemented under authenticated control
- wallet-owned account merge execution is implemented under authenticated control
- primary-wallet switching is implemented under authenticated control
- wallet detach eligibility is implemented under authenticated control
- wallet detach execution is implemented under authenticated control for already eligible wallets
- detached wallets are explicitly reusable through reattachment or wallet-login rebound
- detached wallets now preserve minimal audit-ready lifecycle evidence through `detached_at`

The backend now transitions from:

**authentication flows → ownership flows → authenticated wallet-management flows → detached-wallet lifecycle clarification → detached-identity audit readiness**

### 6. Primary-Wallet Switching

1. authenticated user calls `POST /auth/wallets/primary`
2. backend extracts current authenticated `user_id`
3. backend normalizes and validates `wallet_address`
4. wallet identity store verifies:
   - wallet exists
   - wallet belongs to the current authenticated user
5. backend clears `is_primary` from all other wallets owned by that user
6. backend marks the requested wallet as `is_primary = true`
7. backend returns the refreshed wallet inventory:
   - primary first
   - then by `linked_at`
   - then by address

#### Safety Rules

- switching primary wallet never changes ownership
- switching primary wallet never attaches or detaches wallets
- switching primary wallet never bypasses authentication
- exactly one primary wallet remains for the user after the operation


### 7. Wallet Detach Eligibility

1. authenticated user calls `POST /auth/wallets/detach/check`
2. backend extracts current authenticated `user_id`
3. backend normalizes and validates `wallet_address`
4. wallet identity store verifies:
   - wallet exists
   - wallet belongs to the current authenticated user
5. detach service evaluates:
   - whether the wallet is currently primary
   - how many wallets the user currently owns
6. backend returns a structured detach-eligibility response
7. ownership remains unchanged regardless of the result

#### Safety Rules

- detach eligibility never detaches a wallet
- detach eligibility never changes ownership
- detach eligibility never reassigns primary automatically
- primary wallets remain non-eligible for detach execution
- single-wallet users remain non-eligible for detach execution

### 8. Wallet Detach Execution

1. authenticated user calls `POST /auth/wallets/detach`
2. backend extracts current authenticated `user_id`
3. backend normalizes and validates `wallet_address`
4. detach service reuses the eligibility rules from the detach-check contract
5. backend rejects the request if the wallet is primary or if the user would become wallet-empty
6. store clears `user_id`, `linked_at`, and `is_primary` from the detached wallet identity
7. backend returns the detached wallet snapshot plus the refreshed remaining owned-wallet inventory

#### Safety Rules

- detach execution only works for wallets already eligible under the detach-check rules
- detach execution never reassigns primary automatically
- detach execution never moves the detached wallet to a different user
- detach execution preserves the remaining primary wallet exactly as it was before the operation
