# Architecture

## Current Architectural Direction

The backend follows a **modular monolith** architecture.

This means:

- a single deployable backend application
- clear internal module boundaries
- shared core infrastructure packages
- independent domain growth inside the same repository
- the ability to extract services later only if justified

This is appropriate for the current maturity stage of the project.

---

## Current High-Level Layers

The project is currently structured around these layers:

### 1. Application Layer

Responsible for wiring the system together.

Examples:

- startup
- shutdown
- config loading
- dependency initialization
- module registration
- router creation

### 2. Core Infrastructure Layer

Shared technical capabilities used across modules.

Examples:

- logging
- config
- HTTP helpers
- WebSocket transport
- auth/token utilities
- PostgreSQL client
- Redis client
- readiness/health status

### 3. Module Layer

Business capabilities live inside modules.

Current modules:

- `system`
- `auth`
- `user`

This keeps business logic away from transport and infrastructure details.

---

## Current Transport Style

The backend supports two transport modes:

- REST HTTP
- WebSocket

REST is used for:

- health and readiness
- version information
- login and wallet-auth bootstrap
- authenticated session inspection

WebSocket is used for:

- persistent bidirectional communication
- session-aware actions
- future real-time exchange features

---

## Current HTTP Layer Responsibilities

The HTTP layer is responsible for:

- route registration
- middleware chaining
- auth extraction and validation
- handler dispatch
- JSON response behavior
- readiness exposure

Current implementation already includes:

- `/health`
- `/readiness`
- `/version`
- `/auth/login`
- `/auth/me`
- `/auth/session`
- `/auth/wallet/challenge`
- `/auth/wallet/verify`
- `/ws`

---

## Application Composition Layer

This layer wires the application together.

Responsibilities:

- config loading
- logger initialization
- token service initialization
- WebSocket hub creation
- dispatcher registration
- module registration
- server boot and shutdown lifecycle
- repository/service wiring across modules

Current implementation:

- `internal/app`

---

## Core Package Direction

Current core packages and their roles:

### `internal/core/config`
Configuration loading and normalization.

### `internal/core/logger`
Structured logging for application and request-level events.

### `internal/core/httpx`
HTTP router setup, middleware, auth middleware, and response helpers.

### `internal/core/ws`
WebSocket handler, dispatcher, client session attachment, and hub behavior.

### `internal/core/auth`
JWT minting/parsing, auth claims transport helpers, and shared claims context.

### `internal/core/db`
PostgreSQL bootstrap and connection health logic.

### `internal/core/cache`
Redis bootstrap and connectivity validation.

### `internal/core/status`
Health/readiness model and dependency checks.

---

## Module Direction

### `system`
Lightweight baseline endpoints and system-level transport behavior.

### `auth`
Authentication orchestration, token issuance, current-session resolution, wallet challenge issuance, wallet signature verification, and wallet-auth bootstrap flow handling.

### `user`
User model, repository, and service logic used by auth flows.

---

## Current Persistence Direction

Persistence is currently mixed by maturity level.

### Already persisted

- users through PostgreSQL

### Still bootstrap / non-durable

- wallet auth challenges
- wallet-auth identity bootstrap linkage
- refresh tokens
- session persistence

This is acceptable for the current stage because wallet-auth is still being introduced in controlled steps.

---

## Current Auth Direction

The project currently supports two authentication shapes:

### Development login

- email + fixed development password
- JWT issuance
- persisted or fallback user resolution

### Wallet-auth bootstrap login

- challenge issuance through `POST /auth/wallet/challenge`
- EVM-style signed message verification through `POST /auth/wallet/verify`
- one-time challenge consumption with replay rejection
- wallet-auth JWT enriched with wallet metadata
- fallback wallet identity view without durable wallet persistence yet

This staged approach keeps implementation safe while gradually moving toward production-grade wallet authentication.

---

## Current Session Model

A valid JWT currently provides enough information to resolve:

- user id
- email when present
- wallet address when present
- auth method
- chain value
- issuer
- subject
- expiration metadata

The same claims model now feeds both:

- REST authenticated endpoints
- WebSocket session attachment

This reduces duplicated auth-state logic across transports.

---

## Readiness Philosophy

Readiness is dependency-aware.

The backend can be configured to require:

- PostgreSQL
- Redis

This allows the project to remain flexible in local development while still supporting stricter deployment expectations later.

---

## Why This Architecture Still Fits

At the current stage, this architecture remains the correct choice because it provides:

- fast iteration speed
- low operational complexity
- strong internal structure
- room for safe growth
- good testability potential

There is no current need to split the backend into microservices.

---

## Recommended Architectural Next Step

The next architectural step should focus on improving durability rather than increasing distribution complexity.

Priority areas:

- durable wallet challenge persistence
- wallet identity persistence direction
- stronger session lifecycle management
- continued auth and transport hardening
