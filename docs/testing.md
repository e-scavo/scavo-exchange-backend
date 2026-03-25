# 🧪 Testing

## Full Test Suite

Run the complete backend test suite with:

```bash
go test ./...
```

---

## Wallet Authentication Manual Flow

### 1. Request a Wallet Challenge

```bash
curl -s -X POST http://localhost:8080/auth/wallet/challenge \
  -H 'Content-Type: application/json' \
  -d '{"address":"0xYOURADDRESS","chain":"scavium"}'
```

Expected result:

- HTTP 200
- challenge payload returned
- persisted challenge when PostgreSQL is enabled

---

### 2. Sign the Challenge Message

Use an EVM-compatible wallet or a controlled test script to sign the returned challenge message.

Expected result:

- valid hex signature
- signature tied to the original message and wallet address

---

### 3. Verify the Signature

```bash
curl -s -X POST http://localhost:8080/auth/wallet/verify \
  -H 'Content-Type: application/json' \
  -d '{
    "challenge_id":"CHALLENGE_ID",
    "address":"0xYOURADDRESS",
    "signature":"0xYOUR_SIGNATURE"
  }'
```

Expected result:

- HTTP 200
- access token returned
- wallet metadata included in response
- challenge marked as used
- wallet identity created or reused

---

## Database Validation

### Wallet Challenges

```sql
SELECT id, address, chain, issued_at, expires_at, used_at
FROM auth_wallet_challenges
ORDER BY created_at DESC;
```

### Wallet Identities

```sql
SELECT id, address, created_at
FROM auth_wallet_identities
ORDER BY created_at DESC;
```

---

## Expected Functional Outcomes for 0.4.6

- challenge persists durably when DB is enabled
- reused challenge is rejected
- expired challenge is rejected
- wallet identity is created on first successful verification
- wallet identity is reused on subsequent logins
- JWT includes wallet-related claims
- HTTP session and WebSocket session expose wallet metadata

---

## Recommended Validation Sequence

1. run `go test ./...`
2. request a wallet challenge
3. sign the challenge with a test wallet
4. verify the signature
5. inspect DB rows
6. call `/auth/session`
7. connect to `/ws` using the issued JWT