---

## 🔐 Wallet Authentication Architecture

### Overview

The wallet authentication subsystem is now composed of five main layers:

1. wallet challenge issuance
2. wallet challenge persistence
3. wallet signature verification
4. wallet identity resolution
5. linked platform user resolution and token issuance

---

### Main Components

#### `WalletChallengeService`
Responsible for:

- generating wallet challenges
- building the signable message
- storing the challenge through the configured store
- retrieving challenges by ID
- consuming challenges once verified

#### `WalletVerificationService`
Responsible for:

- loading the persisted challenge
- validating the signed message
- recovering the wallet address from the signature
- comparing the recovered address against the requested address
- marking the challenge as used
- resolving or creating a persistent wallet identity
- resolving or creating a linked platform user
- delegating final token issuance to the auth service

#### `WalletChallengeStore`
Abstract storage contract used by the challenge service.

Current implementations:

- in-memory challenge store
- PostgreSQL challenge store

#### `WalletIdentityStore`
Abstract identity persistence contract used by the verification service.

Current implementations:

- in-memory wallet identity store
- PostgreSQL wallet identity store

The wallet identity store now also persists the durable `user_id` linkage when available.

#### `User Service`
Responsible for:

- resolving or creating development users
- resolving or creating wallet-backed users
- hydrating durable users by id for session reads

#### `TokenService`
Responsible for minting and parsing JWT tokens, including wallet-specific claims such as:

- `uid`
- `email`
- `wallet_id`
- `wallet_address`
- `auth_method`
- `chain`

---

### Durable Flow

The current durable wallet-auth flow is:

1. challenge is created
2. challenge is stored
3. user signs the message
4. signature is verified
5. challenge is consumed atomically
6. wallet identity is resolved or created
7. linked user is resolved or created
8. wallet identity is attached to the linked user
9. JWT is minted with unified identity metadata

---

### Security Guarantees

The current architecture enforces:

- single-use challenge semantics
- expiration-based challenge invalidation
- address recovery from signature
- replay protection after successful challenge consumption
- normalized wallet identity persistence
- deterministic wallet-backed user provisioning
- durable wallet ↔ user linkage when PostgreSQL is enabled

---

### Architectural Boundaries

Phase 0.4.7 intentionally stops before introducing:

- multi-wallet account aggregation
- user-managed linking and unlinking APIs
- refresh token lifecycle
- session persistence
- revocation infrastructure
- account merge workflows across auth methods

These concerns are deferred to upcoming auth/account phases.
