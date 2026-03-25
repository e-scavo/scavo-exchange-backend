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

This flow represents the first wallet-auth bootstrap contract.

1. client calls `POST /auth/wallet/challenge`
2. request includes wallet address and optional chain value
3. backend validates wallet address format
4. backend generates a secure nonce
5. backend builds a stable signable message
6. backend stores the challenge in bootstrap memory storage
7. backend returns challenge id, nonce, message, issue time, and expiration time

---

## Flow 5 - Wallet Signature Verification and Token Issuance

1. client signs the issued wallet challenge message with the requested wallet
2. client calls `POST /auth/wallet/verify`
3. request includes challenge id, wallet address, and signature
4. backend loads the issued challenge and validates expiration and replay state
5. backend recovers the signer address from the signed message
6. backend compares the recovered address with the requested wallet address
7. backend marks the challenge as used
8. backend mints a JWT enriched with wallet-auth metadata
9. backend returns the access token plus wallet-authenticated session identity data

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

## Flow 8 - Future Wallet Signature Login

1. client requests wallet auth challenge
2. backend generates challenge
3. client signs challenge with wallet
4. backend verifies signature against claimed address
5. backend links or resolves user identity
6. backend issues authenticated session token
7. client uses token for further API and WebSocket calls

This is the preferred authentication direction for the DEX-first product model.

---

## Flow 9 - Wallet Portfolio Read

1. client requests portfolio for a linked wallet
2. backend resolves supported assets
3. backend queries SCAVIUM RPC and indexed/local metadata
4. backend aggregates native and token balances
5. backend reads allowances if needed
6. backend returns frontend-ready portfolio view

This flow becomes important once chain and asset modules are introduced.

---

## Flow 10 - DEX Quote Request

1. client requests a quote for token swap
2. backend validates asset pair and amount
3. backend discovers relevant pools
4. backend calculates route
5. backend estimates output, price impact, and fees
6. backend applies slippage-related constraints
7. backend returns quote payload

The backend does not execute the swap here. It only prepares decision-quality information.

---

## Flow 11 - DEX Swap Execution

1. client requests swap preparation data
2. backend validates route and current assumptions
3. backend checks allowance requirements
4. backend returns contract interaction parameters
5. frontend submits transaction request through user wallet
6. user signs transaction
7. transaction is broadcast to SCAVIUM
8. backend tracks transaction status
9. backend exposes pending, confirmed, or failed result states

Settlement is on-chain. The backend supports the flow but does not sign for the user.

---

## Flow 12 - Add Liquidity

1. client requests liquidity addition parameters
2. backend validates pool and token pair
3. backend calculates required counterpart amounts
4. backend checks allowance needs
5. backend returns router interaction parameters
6. user signs through wallet
7. contract mints LP position representation
8. backend tracks resulting transaction and later indexed state

---

## Flow 13 - Indexed Transaction Tracking

1. a user-originated DEX transaction is submitted
2. backend stores or registers tracking intent
3. chain integration or indexer monitors receipt status
4. backend updates transaction state
5. REST and WebSocket surfaces can expose the current lifecycle state

This flow is critical for frontend usability and operational traceability.

---

## Flow 13 - Future Hybrid Deposit Flow

This is not part of the initial implementation, but the architecture must preserve space for it.

1. user requests deposit information
2. backend generates or resolves deposit target
3. user sends funds to exchange-controlled destination
4. backend confirms receipt
5. internal ledger credits custodial balance
6. trading may later occur against internal balances

This flow is intentionally deferred.