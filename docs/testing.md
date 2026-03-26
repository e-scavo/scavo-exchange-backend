# 🧪 Testing

## 1. Run the automated test suite

```bash
go test ./...
```

---

## 2. Request a wallet challenge

```bash
curl -s -X POST http://localhost:8080/auth/wallet/challenge \
  -H 'Content-Type: application/json' \
  -d '{"address":"0x1111111111111111111111111111111111111111","chain":"scavium"}'
```

Expected result:

- challenge payload returned
- persisted challenge when PostgreSQL is enabled

---

## 3. Sign the challenge message

Use an EVM-compatible wallet or a controlled test script to sign the returned challenge message.

Expected result:

- signature tied to the original message and wallet address

---

## 4. Verify the wallet signature

```bash
curl -s -X POST http://localhost:8080/auth/wallet/verify \
  -H 'Content-Type: application/json' \
  -d '{
    "challenge_id":"CHALLENGE_ID",
    "address":"0x1111111111111111111111111111111111111111",
    "signature":"0xSIGNATURE"
  }'
```

Expected result:

- wallet metadata included in response
- challenge marked as used
- wallet identity created or reused
- linked user created or reused
- token issued with unified identity metadata

---

## 5. Validate durable challenge state

```sql
SELECT id, address, used_at
FROM auth_wallet_challenges
ORDER BY created_at DESC
LIMIT 10;
```

Expected result:

- verified challenge shows `used_at`

---

## 6. Validate wallet identity linkage

```sql
SELECT id, address, user_id, created_at
FROM auth_wallet_identities
ORDER BY created_at DESC
LIMIT 10;
```

Expected result:

- verified wallet identity contains a non-null `user_id`

---

## 7. Validate wallet-backed user provisioning

```sql
SELECT id, email, display_name, last_login_at
FROM users
WHERE id = 'u_wallet_1111111111111111111111111111111111111111';
```

Expected result:

- user exists in `users`
- email is synthetic/internal
- display name matches the wallet address
- `last_login_at` reflects the latest successful login

---

## 8. Validate session hydration

Call:

- `GET /auth/me`
- `GET /auth/session`
- authenticated WebSocket `auth.session`

Expected result:

- user payload resolves from the shared `users` table
- wallet metadata remains present in session output
- `uid` is stable across repeated logins for the same wallet

---

## 9. Phase 0.4.7 Validation Summary

This phase is considered validated when:

- challenge persists durably when DB is enabled
- reused challenge is rejected
- expired challenge is rejected
- wallet identity is created on first successful verification
- wallet identity is reused on subsequent logins
- linked user is created on first successful verification
- linked user is reused on subsequent logins
- JWT includes wallet-related and user-related claims
- HTTP session and WebSocket session expose unified identity metadata
