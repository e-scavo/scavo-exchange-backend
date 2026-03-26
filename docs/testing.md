# Testing

## 🧠 Overview

This document defines validation procedures for authentication, identity resolution, wallet ownership, and authenticated wallet linking within the SCAVO Exchange Backend.

Testing validates:

- functional correctness
- persistence integrity
- ownership enforcement
- challenge-purpose correctness
- authenticated wallet-link behavior
- API contract stability

---

## ⚙️ General Validation

Run all tests:

```bash
go test ./...
```

Expected:

- no compilation errors
- all tests passing
- auth and user modules validated
- wallet-link flow tests passing

---

## 🔐 Wallet Authentication Validation

### Step 1 — Create login challenge

```bash
curl -s -X POST http://localhost:8080/auth/wallet/challenge \
  -H "Content-Type: application/json" \
  -d '{"address":"0x...","chain":"scavium"}'
```

Expected:

- `200 OK`
- challenge returned
- challenge purpose behaves as login bootstrap
- payload contains:
  - `id`
  - `message`
  - `expires_at`

---

### Step 2 — Verify wallet login

```bash
curl -s -X POST http://localhost:8080/auth/wallet/verify \
  -H "Content-Type: application/json" \
  -d '{"challenge_id":"...","address":"0x...","signature":"0x..."}'
```

Expected:

- `200 OK`
- valid JWT
- includes:
  - `user_id`
  - `wallet_id`
  - `wallet_address`
  - `auth_method`

---

### Step 3 — Replay protection

Repeat verification with the same login challenge.

Expected:

- `401 Unauthorized`
- error: `wallet_challenge_used`

---

## 🧩 Identity Validation (0.4.7)

### Verify user creation

After successful login:

```sql
SELECT *
FROM users
WHERE email LIKE '%wallet.scavo%';
```

Expected:

- wallet-backed durable user exists
- stable user ID
- email derived from wallet identity

---

### Verify wallet linkage

```sql
SELECT id, address, user_id
FROM auth_wallet_identities;
```

Expected:

- wallet identity exists
- `user_id` is not null

---

## 🏷️ Ownership Validation (0.4.8)

### Ownership metadata

```sql
SELECT id, address, user_id, linked_at, is_primary
FROM auth_wallet_identities
ORDER BY linked_at NULLS LAST, address;
```

Expected:

- `linked_at` populated for linked wallets
- one primary wallet for the first owned wallet
- `user_id` correctly set

---

### Primary-wallet uniqueness

```sql
SELECT user_id, COUNT(*) AS primary_count
FROM auth_wallet_identities
WHERE is_primary = TRUE
GROUP BY user_id
HAVING COUNT(*) > 1;
```

Expected:

- no rows returned

---

### Ownership enforcement

Try to attach the same wallet to another user.

Expected:

- operation rejected
- error equivalent to `wallet_identity_already_linked`

---

## 🔗 Wallet Linking Validation (0.4.9)

### Step 1 — Authenticate first

Obtain a valid access token through dev login or wallet login.

---

### Step 2 — Create wallet-link challenge

```bash
curl -s -X POST http://localhost:8080/auth/wallets/link/challenge \
  -H "Authorization: Bearer <ACCESS_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"address":"0xSECONDARY...","chain":"scavium"}'
```

Expected:

- `200 OK`
- challenge returned
- challenge includes:
  - `purpose = wallet_link`
  - `requested_by_user_id = current authenticated user`

---

### Step 3 — Verify wallet-link challenge

```bash
curl -s -X POST http://localhost:8080/auth/wallets/link/verify \
  -H "Authorization: Bearer <ACCESS_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"challenge_id":"...","address":"0xSECONDARY...","signature":"0x..."}'
```

Expected:

- `200 OK`
- `linked_wallet` returned
- linked wallet belongs to current user
- linked wallet has:
  - `is_primary = false`
  - `linked_at` populated
- `wallets` array reflects the expanded inventory

---

### Step 4 — Validate persisted link

```sql
SELECT id, address, user_id, linked_at, is_primary
FROM auth_wallet_identities
WHERE user_id = '<CURRENT_USER_ID>'
ORDER BY is_primary DESC, linked_at ASC NULLS LAST, address ASC;
```

Expected:

- original primary wallet remains primary
- new linked wallet appears as secondary

---

## 📦 Wallet Inventory API Validation

### Request

```bash
curl -s http://localhost:8080/auth/wallets \
  -H "Authorization: Bearer <ACCESS_TOKEN>"
```

Expected:

- `200 OK`
- `wallets` array returned
- primary wallet first
- newly linked wallet included after successful 0.4.9 linking

---

## 🔄 Session Validation

### `/auth/me`

```bash
curl -s http://localhost:8080/auth/me \
  -H "Authorization: Bearer <ACCESS_TOKEN>"
```

Expected:

- unified durable user identity
- wallet-backed context still valid

### `/auth/session`

Expected:

- consistent claims
- matches `/auth/me`
- no forced token refresh after wallet linking

---

## ⚠️ Error Handling Validation

### Invalid address
Expected:

- `400`
- `invalid_wallet_address`

### Invalid signature
Expected:

- `401`
- `invalid_wallet_signature`

### Challenge expired
Expected:

- `401`
- `wallet_challenge_expired`

### Wallet already linked to another user
Expected:

- `409`
- `wallet_identity_already_linked`

### Wallet already linked to current user
Expected:

- `409`
- `wallet_identity_already_linked_to_user`

### Challenge belongs to another authenticated user
Expected:

- `403`
- `wallet_link_challenge_user_mismatch`

### Wrong challenge purpose
Expected:

- `409`
- `wallet_challenge_purpose_mismatch`

---

## 🧪 Internal Test Coverage

Modules covered:

- `internal/modules/auth`
- `internal/modules/user`

Key validations now include:

- signature recovery
- challenge lifecycle
- durable identity linking
- ownership enforcement
- authenticated wallet-link contract
- wallet-link conflict rejection
- wallet inventory refresh after link

---

## 🧭 Future Testing (Post 0.4.9)

Planned:

- unlink scenarios
- primary-wallet switching scenarios
- cross-user ownership transfer edge cases
- multi-auth merge preparation testing

---

## 🧩 Summary

Testing at Phase 0.4.9 guarantees:

- authentication correctness
- identity persistence
- ownership consistency
- authenticated wallet linking correctness
- API stability across both login and link flows