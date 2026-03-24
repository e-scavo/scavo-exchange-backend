# Core Flows

## Flow 1 - Current HTTP Development Login

This flow represents the current minimal login bootstrap already present in the backend.

1. client calls `POST /auth/login`
2. backend validates the development payload
3. backend issues JWT
4. client stores token
5. token can later be used in authenticated flows

This flow is temporary and will evolve in later phases.

---

## Flow 2 - Current WebSocket Session Attachment

1. client connects to `/ws`
2. backend upgrades the connection
3. backend optionally parses token from request context or handshake path
4. client is attached to hub
5. dispatcher routes incoming action messages
6. if the token is valid, session context can be attached to the client

This flow is already reflected in the current project baseline.

---

## Flow 3 - Future Wallet Signature Login

1. client requests wallet auth challenge
2. backend generates challenge
3. client signs challenge with wallet
4. backend verifies signature against claimed address
5. backend links or resolves user identity
6. backend issues authenticated session token
7. client uses token for further API and WebSocket calls

This is the preferred authentication direction for the DEX-first product model.

---

## Flow 4 - Wallet Portfolio Read

1. client requests portfolio for a linked wallet
2. backend resolves supported assets
3. backend queries SCAVIUM RPC and indexed/local metadata
4. backend aggregates native and token balances
5. backend reads allowances if needed
6. backend returns frontend-ready portfolio view

This flow becomes important once chain and asset modules are introduced.

---

## Flow 5 - DEX Quote Request

1. client requests a quote for token swap
2. backend validates asset pair and amount
3. backend discovers relevant pools
4. backend calculates route
5. backend estimates output, price impact, and fees
6. backend applies slippage-related constraints
7. backend returns quote payload

The backend does not execute the swap here. It only prepares decision-quality information.

---

## Flow 6 - DEX Swap Execution

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

## Flow 7 - Add Liquidity

1. client requests liquidity addition parameters
2. backend validates pool and token pair
3. backend calculates required counterpart amounts
4. backend checks allowance needs
5. backend returns router interaction parameters
6. user signs through wallet
7. contract mints LP position representation
8. backend tracks resulting transaction and later indexed state

---

## Flow 8 - Indexed Transaction Tracking

1. a user-originated DEX transaction is submitted
2. backend stores or registers tracking intent
3. chain integration or indexer monitors receipt status
4. backend updates transaction state
5. REST and WebSocket surfaces can expose the current lifecycle state

This flow is critical for frontend usability and operational traceability.

---

## Flow 9 - Future Hybrid Deposit Flow

This is not part of the initial implementation, but the architecture must preserve space for it.

1. user requests deposit information
2. backend generates or resolves deposit target
3. user sends funds to exchange-controlled destination
4. backend confirms receipt
5. internal ledger credits custodial balance
6. trading may later occur against internal balances

This flow is intentionally deferred.