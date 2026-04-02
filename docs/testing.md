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

## 🔗 Wallet-Owned Account Merge Validation (0.4.10)

### Step 1 — Authenticate first

Obtain a valid access token through dev login or wallet login.

---

### Step 2 — Create account-merge challenge

```bash
curl -s -X POST http://localhost:8080/auth/account/merge/wallet/challenge \
  -H "Authorization: Bearer <ACCESS_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"address":"0xSOURCE_PRIMARY...","chain":"scavium"}'
```

Expected:

- `200 OK`
- challenge returned
- challenge includes:
  - `purpose = account_merge`
  - `requested_by_user_id = current authenticated user`

---

### Step 3 — Verify account-merge challenge

```bash
curl -s -X POST http://localhost:8080/auth/account/merge/wallet/verify \
  -H "Authorization: Bearer <ACCESS_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"challenge_id":"...","address":"0xSOURCE_PRIMARY...","signature":"0x..."}'
```

Expected:

- `200 OK`
- `merged_wallet` returned
- `source_user_id` returned
- `target_user_id` returned
- all wallets from the source wallet-owned account now appear in the target inventory
- if the target already had a primary wallet, it remains primary

---

### Step 4 — Validate persisted merge

```sql
SELECT id, address, user_id, linked_at, is_primary
FROM auth_wallet_identities
WHERE user_id IN ('<TARGET_USER_ID>', '<SOURCE_USER_ID>')
ORDER BY user_id, is_primary DESC, linked_at ASC NULLS LAST, address ASC;
```

Expected:

- all former source-user wallets now point to `<TARGET_USER_ID>`
- `<SOURCE_USER_ID>` owns zero wallets after merge
- target primary semantics remain deterministic

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
- merged wallets included after successful 0.4.10 merge

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

### Merge source wallet is not linked
Expected:

- `409`
- `wallet_account_merge_source_not_linked`

### Merge is not required
Expected:

- `409`
- `wallet_account_merge_not_required`

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
- authenticated wallet-owned account merge contract
- atomic ownership consolidation at the store layer
- authenticated wallet detach-eligibility contract
- detached-wallet reattachment semantics after detach
- detached-wallet wallet-login rebound semantics after detach

---

## 🧭 Future Testing (Post 0.4.14)

Planned:

- detached-identity audit metadata once introduced
- primary-replacement preconditions for riskier detach execution
- cross-user ownership transfer edge cases
- post-merge user archival testing
- multi-auth merge preparation testing

---

## 🧩 Summary

Testing at Phase 0.4.14 guarantees:

- authentication correctness
- identity persistence
- ownership consistency
- authenticated wallet linking correctness
- authenticated wallet-owned account merge correctness
- authenticated primary-wallet switching correctness
- authenticated wallet detach-eligibility correctness
- detached-wallet reattachment correctness after detach
- detached-wallet wallet-login rebound correctness after detach
- API stability across login, link, merge, primary-switch, detach-check, and detach flows

## 0.4.11 — Primary Wallet Switching

### API

```bash
curl -s -X POST http://localhost:8080/auth/wallets/primary \
  -H "Authorization: Bearer <ACCESS_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "wallet_address": "0xYOUR_OWNED_SECONDARY_WALLET"
  }'
```

### Expected Result

- HTTP `200 OK`
- response contains:
  - `primary_wallet`
  - refreshed `wallets`
- the requested wallet becomes the only wallet with `is_primary = true`

### SQL Verification

```sql
SELECT id, address, user_id, linked_at, is_primary
FROM auth_wallet_identities
WHERE user_id = 'u_test_example_com'
ORDER BY is_primary DESC, linked_at ASC NULLS LAST, address ASC;
```

### Expected State

- requested wallet is first
- requested wallet has `is_primary = true`
- previous primary wallet has `is_primary = false`
- no ownership changes occurred


## 0.4.12 — Wallet Detach Eligibility

### API

```bash
curl -s -X POST http://localhost:8080/auth/wallets/detach/check \
  -H "Authorization: Bearer <ACCESS_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "wallet_address": "0xYOUR_OWNED_WALLET"
  }'
```

### Expected Result

- HTTP `200 OK`
- response contains:
  - `wallet_address`
  - `eligible`
  - `is_primary`
  - `owned_wallet_count`
  - `reasons`
- ownership remains unchanged

### Expected Cases

#### Eligible
- target wallet belongs to current user
- target wallet is not primary
- current user owns more than one wallet
- `eligible = true`
- `reasons = []`

#### Not Eligible
- target wallet is primary → `wallet_is_primary`
- current user would have no wallets left after detach → `user_would_have_no_wallets`
- multiple reasons may be returned together

## 0.4.13 — Wallet Detach Execution

### API

```bash
curl -s -X POST http://localhost:8080/auth/wallets/detach \
  -H "Authorization: Bearer <ACCESS_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "wallet_address": "0xYOUR_ELIGIBLE_OWNED_WALLET"
  }'
```

### Expected Result

- HTTP `200 OK` for eligible owned non-primary wallets
- response contains:
  - `detached_wallet`
  - `wallets`
  - `check`
- detached wallet returns with cleared ownership fields
- remaining wallet inventory preserves the existing primary wallet

### Conflict Case

```bash
curl -s -X POST http://localhost:8080/auth/wallets/detach \
  -H "Authorization: Bearer <ACCESS_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "wallet_address": "0xYOUR_PRIMARY_OR_ONLY_WALLET"
  }'
```

- HTTP `409 Conflict` when current guardrails still mark the wallet as non-eligible
- response error: `wallet_detach_not_eligible`
- conflict payload includes the same detach-check snapshot and reasons used to reject execution


## 0.4.14 — Detached Wallet Reattachment and Rebound

### Expected Coverage

- a wallet detached through `POST /auth/wallets/detach` can later be reattached through the protected wallet-link contract
- a wallet detached from one user can later re-enter `POST /auth/wallet/verify` and rebound into its wallet-owned durable user identity
- neither path restores prior primary ownership implicitly

### Expected Result

- detached wallet keeps its durable wallet identity record
- detached wallet can be safely reused
- detached wallet lifecycle remains explicit even without new audit columns

```

FILE: docs/architecture.md
```md
# Architecture

## 🧠 Overview

The SCAVO Exchange Backend is designed around a **wallet-first identity architecture** that progressively evolves into a **durable account model** suitable for exchange-grade ownership and future multi-auth identity operations.

The architecture intentionally separates:

- authentication mechanism
- wallet identity representation
- durable platform user abstraction
- persisted wallet ownership
- authenticated wallet-management contracts

---

## 🧩 Core Layers

### 1. Transport Layer
- HTTP API
- JSON-based communication
- stateless request handling

### 2. Auth Layer
Located in:

- `internal/modules/auth`

Responsibilities:

- wallet challenge generation
- wallet signature verification
- JWT issuance
- identity resolution
- ownership enforcement
- authenticated wallet-linking flows

### 3. User Layer
Located in:

- `internal/modules/user`

Responsibilities:

- durable user creation
- durable user resolution
- future auth-provider expansion

### 4. Persistence Layer
- PostgreSQL (primary)
- in-memory fallback (dev/testing)

---

## 🔐 Identity Model Evolution

### Pre 0.4.6
- identity was session-oriented
- wallet state was not durable

### 0.4.6 — Wallet Identity Persistence
- wallet identity stored in `auth_wallet_identities`
- wallet address becomes a stable registry entry

### 0.4.7 — Unified Identity Model
- wallet identity linked to durable user
- `user_id` introduced
- JWT identity unified around durable platform user

### 0.4.8 — Ownership Model Introduction
wallet identity becomes a first-class ownership entity with:

- `id`
- `address`
- `user_id`
- `linked_at`
- `detached_at`
- `is_primary`

### 0.4.9 — Authenticated Wallet Linking Contract
the architecture adds a dedicated authenticated linking flow, still based on challenge + signature verification, but now explicitly bound to the current authenticated user.

This is the first backend-managed wallet operation that acts on ownership under an existing session rather than during initial login bootstrap.

### 0.4.10 — Wallet-Owned Account Merge Execution
the architecture now adds a second authenticated ownership operation that allows the current user to absorb a wallet-owned source account after an explicit merge challenge is signed by the source wallet.

---

## 🏷️ Ownership Model

### Core Rules

1. a wallet belongs to exactly one user
2. a user can own multiple wallets
3. only one wallet per user can be primary
4. wallet ownership cannot be reassigned across users
5. authenticated wallet linking adds secondary wallets only
6. 0.4.9 does not switch primary ownership during linking
7. 0.4.10 preserves the current target primary wallet when a merge occurs
8. 0.4.11 allows explicit primary-wallet reassignment only within the current owner's wallet set
9. 0.4.12 exposes detach eligibility as a guarded evaluation contract and 0.4.13 adds detach execution only for already eligible owned wallets
10. 0.4.14 clarifies that detached wallet identities remain reusable known identities that can be reattached or rebound later without new schema state

---

## 🏷️ Ownership Metadata

| Field | Description |
|------|-------------|
| `user_id` | owning durable user |
| `linked_at` | ownership creation timestamp |
| `is_primary` | primary-wallet flag |

---

## 🧾 Challenge Model

### Pre 0.4.9
wallet challenge was effectively used only for authentication bootstrap.

### 0.4.9
wallet challenges now include:

- `purpose`
- `requested_by_user_id`

### Challenge purposes

- `auth_bootstrap`
- `wallet_link`
- `account_merge`

This avoids reusing the same challenge semantics blindly across two very different operations.

---

## 🔄 Authentication Flow (Wallet Login)

1. client requests login challenge
2. backend persists challenge with `auth_bootstrap` purpose
3. client signs message
4. backend verifies signature
5. challenge is consumed
6. wallet identity is resolved
7. durable user is resolved or created
8. ownership is enforced
9. JWT is issued

---

## 🔄 Authenticated Wallet Linking Flow

1. user already holds valid JWT
2. client requests link challenge:
   - `POST /auth/wallets/link/challenge`
3. backend persists challenge with:
   - `purpose = wallet_link`
   - `requested_by_user_id = current user`
4. user signs with the secondary wallet
5. client submits:
   - `POST /auth/wallets/link/verify`
6. backend validates:
   - challenge existence
   - challenge freshness
   - purpose correctness
   - requesting user correctness
   - signature correctness
7. backend resolves wallet identity
8. backend rejects ownership conflict if wallet belongs elsewhere
9. backend attaches wallet as non-primary
10. backend returns updated wallet inventory

---

## 🔄 Authenticated Wallet Account Merge Flow

1. user already holds valid JWT
2. client requests merge challenge:
   - `POST /auth/account/merge/wallet/challenge`
3. backend persists challenge with:
   - `purpose = account_merge`
   - `requested_by_user_id = current user`
4. source wallet signs the merge challenge
5. client submits:
   - `POST /auth/account/merge/wallet/verify`
6. backend validates:
   - challenge existence
   - challenge freshness
   - purpose correctness
   - requesting user correctness
   - signature correctness
   - source wallet ownership existence
7. backend derives the source user from wallet ownership
8. backend atomically reassigns all source-user wallets to the current target user
9. backend returns updated wallet inventory

## 🔌 API Layer

### Auth endpoints
- `/auth/login`
- `/auth/wallet/challenge`
- `/auth/wallet/verify`
- `/auth/me`
- `/auth/session`

### Wallet ownership endpoints
- `/auth/wallets`
- `/auth/wallets/link/challenge`
- `/auth/wallets/link/verify`
- `/auth/account/merge/wallet/challenge`
- `/auth/account/merge/wallet/verify`

---

## 🧾 JWT Design

JWT tokens are:

- stateless
- short-lived
- self-contained

Claims include:

- `user_id`
- `wallet_id`
- `wallet_address`
- `auth_method`

Wallet linking and wallet-owned account merge do not mint a fresh token because both operate under an already authenticated durable session.

---

## 🗄️ Data Model

### `auth_wallet_challenges`
stores challenge lifecycle and, from 0.4.9 onward, also stores operation metadata:

- `purpose`
- `requested_by_user_id`
- issued / expires / used lifecycle

### `auth_wallet_identities`
stores wallet registry and ownership:

- `id`
- `address`
- `user_id`
- `linked_at`
- `detached_at`
- `is_primary`

### `users`
stores durable user abstraction:

- login-independent identity
- wallet-backed users now
- future auth-provider aggregation later

---

## ⚙️ Design Decisions

### Wallet-first approach
chosen because it aligns with crypto-native UX and reduces early auth-surface complexity.

### Separation of identity and ownership
wallet identity is not the same as durable user identity.
Ownership is explicit rather than inferred.

### Challenge-purpose separation
0.4.9 and 0.4.10 extend the challenge system rather than introducing parallel challenge subsystems, but still keep semantic separation through `purpose`.

### Incremental evolution
each subphase introduces one structural improvement while preserving previous behavior.

---

## ⚠️ Constraints

Still intentionally not supported:

- wallet unlink
- arbitrary cross-user wallet transfer outside wallet-signed merge
- multi-auth merge resolution
- user-record archival after merge

---

## 🚧 Current Lifecycle Audit Readiness (0.4.15)

Detached wallet identities now preserve minimal lifecycle metadata:

- `detached_at`

This timestamp is stamped during detach execution and intentionally survives later reattachment or wallet-login rebound so the backend can distinguish previously detached reusable identities from identities that have never been detached.


### Later phases
- account consolidation
- multi-auth identity merging
- recovery flows
- compliance-ready identity expansion

---

## 🧩 Summary

At the end of 0.4.15:

- wallet authentication is stable
- durable identity is stable
- wallet ownership is stable
- authenticated wallet linking is implemented
- wallet-owned account merge execution is implemented
- explicit primary-wallet switching is implemented
- authenticated wallet detach eligibility is implemented
- authenticated wallet detach execution is implemented for already eligible owned wallets
- detached wallet identities are explicitly reusable after detach
- detached wallet identities now preserve minimal audit-ready lifecycle metadata through `detached_at`
- the backend is structurally ready to move from ownership persistence into richer detached-identity observability only if future phases require it


---

## Phase 0.4.16 Testing Notes

### Goal

Validate the enriched wallet inventory read model exposed by `GET /auth/wallets`.

### Coverage Added

Handler-level validation covers:

- successful wallet inventory response
- explicit active-wallet status projection
- primary wallet visibility
- `linked_at` visibility for owned wallets
- absence of `detached_at` for active wallets that were never detached
- preservation of `detached_at` after detach + reattach

### Validation Command

```bash
go test ./...
```

### Expected Result

- `internal/modules/auth` passes
- no regressions appear in the rest of the backend tree


## Phase 0.4.17 Testing Notes

### Goal

Validate wallet inventory query filtering and sorting on top of the lifecycle-aware `GET /auth/wallets` read model.

### Coverage Added

Handler-level validation covers:

- backward-compatible wallet inventory response without query params
- `primary=true` returning only primary wallets
- `primary=false` returning only non-primary wallets
- `status=active` returning currently owned wallets
- `status=detached` returning an empty result under the current owned-wallet route contract
- `sort=linked_at&order=desc` returning wallets in explicit descending linked order
- invalid query params returning `400` with explicit error codes

### Validation Command

```
go test ./...
```

### Manual API Checks

```
curl -s http://localhost:8080/auth/wallets \
  -H "Authorization: Bearer $TOKEN"

curl -s "http://localhost:8080/auth/wallets?primary=true" \
  -H "Authorization: Bearer $TOKEN"

curl -s "http://localhost:8080/auth/wallets?primary=false" \
  -H "Authorization: Bearer $TOKEN"

curl -s "http://localhost:8080/auth/wallets?status=active" \
  -H "Authorization: Bearer $TOKEN"

curl -s "http://localhost:8080/auth/wallets?status=detached" \
  -H "Authorization: Bearer $TOKEN"

curl -s "http://localhost:8080/auth/wallets?sort=linked_at&order=desc" \
  -H "Authorization: Bearer $TOKEN"
```

### Expected Result

- `internal/modules/auth` passes with the new handler coverage
- default inventory behavior remains backward compatible
- explicit filtering and sorting behave deterministically
- unsupported query values are rejected with `400`



## Phase 0.4.18 Testing Notes

### Goal

Validate wallet inventory pagination on top of the lifecycle-aware, filterable, and sortable `GET /auth/wallets` read model.

### Coverage Added

Handler-level validation covers:

- backward-compatible wallet inventory response with metadata defaults
- `limit` only
- `offset` only
- `limit + offset`
- empty but valid inventory window when offset exceeds the filtered inventory size
- invalid `limit` returning `400` with explicit error codes
- invalid `offset` returning `400` with explicit error codes

### Validation Command

```
go test ./...
```

### Manual API Checks

```
curl -s http://localhost:8080/auth/wallets \
  -H "Authorization: Bearer $TOKEN"

curl -s "http://localhost:8080/auth/wallets?limit=2" \
  -H "Authorization: Bearer $TOKEN"

curl -s "http://localhost:8080/auth/wallets?offset=1" \
  -H "Authorization: Bearer $TOKEN"

curl -s "http://localhost:8080/auth/wallets?limit=1&offset=1" \
  -H "Authorization: Bearer $TOKEN"

curl -s "http://localhost:8080/auth/wallets?status=active&sort=linked_at&order=desc&limit=2&offset=0" \
  -H "Authorization: Bearer $TOKEN"
```


## Phase 0.4.19 Testing Notes

### Goal
Validate wallet inventory navigation metadata on top of the existing lifecycle-aware, filterable, sortable, and paginated `GET /auth/wallets` read model.

### Coverage Added
Handler-level validation covers:

- default wallet inventory response with `returned` and `has_more`
- paginated window with `has_more=true`
- final paginated window with `has_more=false`
- empty valid window with `returned=0`
- filtered and sorted paginated window with correct navigation metadata

### Validation Command

```
go test ./...
```

### Manual API Checks

```
curl -s http://localhost:8080/auth/wallets \
  -H "Authorization: Bearer $TOKEN"

curl -s "http://localhost:8080/auth/wallets?limit=2" \
  -H "Authorization: Bearer $TOKEN"

curl -s "http://localhost:8080/auth/wallets?limit=1&offset=1" \
  -H "Authorization: Bearer $TOKEN"

curl -s "http://localhost:8080/auth/wallets?primary=false&sort=linked_at&order=desc&limit=2&offset=1" \
  -H "Authorization: Bearer $TOKEN"
```


## Phase 0.4.21 Testing Notes

### Goal
Validate hardened query-parameter rules for the lifecycle-aware, filterable, sortable, and paginated `GET /auth/wallets` read model.

### Coverage Added
Handler-level validation covers:

- `sort=linked_at` without `order` defaulting to ascending behavior
- `order` without `sort` returning `invalid_order_requires_sort`
- offset-only requests remaining valid and unbounded
- existing invalid parameter validation remaining stable

### Validation Command

```
go test ./...
```


## Phase 0.4.22 Testing Notes

### Goal
Validate that wallet inventory response examples and documentation are aligned with the implemented `GET /auth/wallets` contract.

### Coverage Added
Documentation and contract review cover:

- response examples including `returned` and `has_more`
- documented semantics for `next_offset` and `previous_offset`
- explicit unbounded (`limit=0`) vs bounded response behavior
- no behavioral regression expected because the subphase is documentation-only

### Validation Command

```
go test ./...
```


## Phase 0.4.23 Testing Notes

### Goal
Validate that the operator-facing query examples for `GET /auth/wallets` are aligned with the implemented handler contract and its existing validation rules.

### Coverage Added
Documentation and manual validation guidance now cover:

- base inventory request examples
- filtered examples such as `primary=true`
- sorted examples such as `sort=linked_at&order=desc`
- bounded pagination examples using `limit` and `offset`
- explicit invalid example: `order` without `sort`
- no behavioral regression expected because the subphase is documentation-only

### Validation Command

```
go test ./...
```


## Phase 0.4.24 Testing Notes

### Goal
Validate that the manual validation guidance for `GET /auth/wallets` covers the real handler contract end-to-end.

### Coverage Added
Documentation and manual validation guidance now cover:

- base authenticated inventory request
- filtered requests using `primary` and `status`
- sorted requests using `sort=linked_at&order=desc`
- bounded pagination requests using `limit` and `offset`
- unbounded offset-only requests
- expected manual interpretation of `returned`, `has_more`, `next_offset`, and `previous_offset`
- explicit invalid-query checks such as `order` without `sort`, invalid `status`, invalid `primary`, and invalid `limit`
- no behavioral regression expected because the subphase is documentation-only

### Validation Command

```
go test ./...
```


## Phase 0.4.25 Testing Notes

### Goal
Validate that the authenticated wallet inventory exposes actionability hints consistent with the existing wallet detach and primary-switch rules.

### Coverage Added
Handler-level coverage now verifies:

- a single primary wallet returns `can_set_primary=false`, `can_detach=false`, and both detach block reasons
- a two-wallet inventory marks the primary wallet as non-actionable for detach and marks the secondary wallet as detachable and promotable
- detach block reasons remain aligned with the existing domain constants

### Validation Command

```
go test ./...
```
