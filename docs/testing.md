# Testing

## 🧠 Overview

This document defines validation procedures for authentication, identity resolution, and wallet ownership within the SCAVO Exchange Backend.

Testing is structured to validate:

- functional correctness
- persistence integrity
- ownership enforcement
- API behavior

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

---

## 🔐 Wallet Authentication Validation

### Step 1 — Create challenge

```bash
curl -s -X POST http://localhost:8080/auth/wallet/challenge \
  -H "Content-Type: application/json" \
  -d '{"address":"0x...","chain":"scavium"}'
```

Expected:

- `200 OK`
- challenge returned
- contains:
  - `id`
  - `message`
  - `expires_at`

---

### Step 2 — Verify wallet

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

Repeat verification with same challenge:

Expected:

- `401 Unauthorized`
- error: `wallet_challenge_used`

---

## 🧩 Identity Validation (0.4.7)

### Verify user creation

After successful login:

```sql
SELECT * FROM users WHERE email LIKE '%wallet.scavo';
```

Expected:

- user exists
- email derived from wallet
- stable user ID

---

### Verify wallet linkage

```sql
SELECT id, address, user_id
FROM auth_wallet_identities;
```

Expected:

- wallet identity exists
- `user_id` is NOT NULL

---

## 🏷️ Ownership Validation (0.4.8)

### Ownership metadata

```sql
SELECT id, address, user_id, linked_at, is_primary
FROM auth_wallet_identities
ORDER BY created_at;
```

Expected:

- `linked_at` populated
- `is_primary = true` for first wallet
- `user_id` correctly set

---

### Primary wallet uniqueness

```sql
SELECT user_id, COUNT(*) AS primary_count
FROM auth_wallet_identities
WHERE is_primary = TRUE
GROUP BY user_id
HAVING COUNT(*) > 1;
```

Expected:

- **no rows returned**

---

### Ownership enforcement

Try to attach same wallet to another user:

Expected:

- operation rejected
- error: `wallet_identity_already_linked`

---

## 📦 Wallet Inventory API Validation

### Request

```bash
curl -s http://localhost:8080/auth/wallets \
  -H "Authorization: Bearer <ACCESS_TOKEN>"
```

---

### Expected Response

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

### Ordering validation

Verify:

- primary wallet appears first
- secondary wallets sorted by:
  - `linked_at`
  - then address

---

## 🔄 Session Validation

### `/auth/me`

```bash
curl -s http://localhost:8080/auth/me \
  -H "Authorization: Bearer <ACCESS_TOKEN>"
```

Expected:

- unified user identity
- wallet context present

---

### `/auth/session`

Expected:

- consistent claims
- matches `/auth/me`

---

## ⚠️ Error Handling Validation

### Invalid address

Expected:

- `400`
- `invalid_wallet_address`

---

### Invalid signature

Expected:

- `401`
- `invalid_wallet_signature`

---

### Challenge expired

Expected:

- `401`
- `wallet_challenge_expired`

---

### Wallet already linked

Expected:

- `409` or `error`
- `wallet_identity_already_linked`

---

## 🧪 Internal Test Coverage

Modules covered:

- `internal/modules/auth`
- `internal/modules/user`

Key validations:

- signature recovery
- challenge lifecycle
- identity linking
- ownership enforcement

---

## 🧭 Future Testing (Post 0.4.8)

Planned:

- wallet linking API tests
- unlink scenarios
- ownership transfer edge cases
- multi-auth merge testing

---

## 🧩 Summary

Testing at Phase 0.4.8 guarantees:

- authentication correctness
- identity persistence
- ownership consistency
- API stability

The backend is validated for:

**wallet login → user identity → ownership model → multi-wallet readiness**