# Core Flows

## Flow 1 - Current HTTP Development Login

This flow represents the current stabilized development login baseline.

1. client calls `POST /auth/login`
2. HTTP handler decodes the request
3. auth service validates the development credentials
4. auth service resolves or creates the development user
5. auth service issues JWT
6. backend returns token metadata and resolved user id
7. client stores token

This flow is still temporary and intentionally bootstrap-oriented.

---

## Flow 2 - Current Authenticated Identity Read

This flow represents the minimal authenticated REST identity path.

1. client calls `GET /auth/me`
2. client sends bearer token in `Authorization`
3. HTTP auth middleware extracts and validates the token
4. middleware injects claims into request context
5. auth handler resolves current user through the `auth` and `user` modules
6. backend returns the current authenticated user payload

---

## Flow 3 - Current Authenticated Session Read

This flow represents the current session-oriented REST path.

1. client calls `GET /auth/session`
2. client sends bearer token in `Authorization`
3. HTTP auth middleware extracts and validates the token
4. middleware injects claims into request context
5. auth handler resolves a shared session view
6. backend returns authenticated session metadata plus resolved user information

---

## Flow 4 - Wallet Challenge Bootstrap

This flow represents the wallet-auth challenge contract.

1. client calls `POST /auth/wallet/challenge`
2. request includes wallet address and optional chain value
3. backend validates wallet address format
4. backend generates a secure nonce
5. backend builds a stable signable message
6. backend stores the challenge in the configured challenge store
7. backend returns challenge id, nonce, message, issue time, and expiration time

---

## Flow 5 - Wallet Signature Verification and Unified Identity Login

1. client signs the issued wallet challenge message with the requested wallet
2. client calls `POST /auth/wallet/verify`
3. request includes challenge id, wallet address, and signature
4. backend loads the issued challenge and validates expiration and replay state
5. backend recovers the signer address from the signed message
6. backend compares the recovered address with the requested wallet address
7. backend marks the challenge as used
8. backend resolves or creates a wallet identity
9. backend resolves or creates a linked platform user
10. backend persists `auth_wallet_identities.user_id` when durable storage is available
11. backend mints a JWT enriched with unified wallet/user metadata
12. backend returns the access token plus wallet-authenticated session identity data

---

## Flow 6 - Current WebSocket Session Attachment

1. client connects to `/ws`
2. backend upgrades the connection
3. backend extracts token from `Authorization` header or `token` query parameter
4. backend validates JWT if a token is present
5. backend enriches the client session with user id, email, wallet address, auth method, chain, subject, issuer, and expiration
6. client is attached to hub
7. dispatcher routes incoming action messages

This flow is already reflected in the current project baseline.

---

## Flow 7 - Current WebSocket Auth Session Read

1. authenticated client sends `auth.session`
2. WebSocket auth guard checks that the client has a session
3. auth module resolves a shared session view from stored claims
4. backend returns session metadata and resolved user information

This flow keeps WebSocket aligned with REST session semantics.

---

## Flow 8 - Wallet Portfolio Read

1. client requests portfolio for a linked wallet
2. backend resolves supported assets
3. backend queries SCAVIUM RPC and indexed/local metadata
4. backend aggregates native and token balances
5. backend reads allowances if needed
6. backend returns frontend-ready portfolio view
