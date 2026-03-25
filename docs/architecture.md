---

## 🔐 Wallet Authentication Architecture

### Overview

The wallet authentication subsystem is now composed of four main layers:

1. wallet challenge issuance
2. wallet challenge persistence
3. wallet signature verification
4. wallet identity resolution and token issuance

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

#### `TokenService`
Responsible for minting and parsing JWT tokens, including wallet-specific claims such as:

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
7. JWT is minted with wallet metadata

---

### Security Guarantees

The current architecture enforces:

- single-use challenge semantics
- expiration-based challenge invalidation
- address recovery from signature
- replay protection after successful challenge consumption
- normalized wallet identity persistence

---

### Architectural Boundaries

Phase 0.4.6 intentionally stops before introducing:

- wallet ↔ user ownership linking
- account aggregation
- refresh token lifecycle
- session persistence
- revocation infrastructure

These concerns are deferred to upcoming auth/account phases.