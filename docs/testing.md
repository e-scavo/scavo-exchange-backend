# Testing

## Objective

This document defines the testing direction for the SCAVO Exchange backend.

Its purpose is to ensure that validation grows together with the architecture so that infrastructure and feature work do not become fragile as the project expands.

---

## Current Validation Baseline

At the current stage, the project should already support the following validation categories:

- unit validation for domain services
- unit validation for readiness/status logic
- integration validation for PostgreSQL-backed repositories
- smoke validation for HTTP login baseline
- smoke validation for authenticated identity read baseline
- smoke validation for authenticated session read baseline
- migration status validation through the migration workflow

This is the first practical validation baseline for the project.

---

## Why Testing Starts Early

The backend will gradually introduce:

- persistence
- cache
- chain integration
- smart contracts
- indexing
- real-time updates
- background processing
- hybrid growth paths

If testing is postponed until those systems are already deeply integrated, regression control becomes much harder.

For this reason, testing direction is defined during the foundation stage.

---

## Testing Philosophy

Testing should be:

- incremental
- phase-appropriate
- behavior-oriented
- architecture-aware
- useful for regression prevention

The goal is not maximum test volume immediately.

The goal is reliable growth.

---

## Test Layer Model

The backend should evolve through multiple testing layers.

### Unit Tests

Scope:

- pure functions
- validation logic
- small isolated helpers
- deterministic service logic with mocked dependencies
- readiness logic and dependency evaluation
- auth service orchestration
- auth transport helper behavior
- auth claims context behavior
- auth session view behavior

Purpose:

- fast feedback
- low-cost regression prevention
- isolated rule validation

---

### Service Tests

Scope:

- module service behavior
- orchestration logic
- business decision boundaries
- dependency interactions through interfaces

Purpose:

- validate real business flows without requiring full transport or real infrastructure for every case

---

### Repository Tests

Scope:

- persistence behavior
- SQL mapping correctness
- transaction behavior
- query correctness
- migration-backed repository behavior

Purpose:

- validate DB-facing logic once repositories are introduced

---

### Integration Tests

Scope:

- interactions between application layers
- HTTP routes with infrastructure dependencies
- basic end-to-end module behavior inside the backend

Purpose:

- validate wiring and environment behavior

---

### Chain Integration Tests

Scope:

- RPC interactions
- contract read helpers
- gas estimation helpers
- transaction state interactions where applicable

Purpose:

- validate that blockchain-related infrastructure works against expected SCAVIUM behavior

---

### Contract/Backend End-to-End Tests

Scope:

- interaction between deployed contracts and backend logic
- quote-to-transaction flow support
- liquidity and swap-support flows later

Purpose:

- validate DEX behavior across system boundaries

---

### Smoke Tests

Scope:

- startup
- basic route availability
- minimal dependency checks
- internal environment sanity
- development login path
- authenticated identity read path
- authenticated session read path

Purpose:

- fast validation for local and internal testing environments

---

## Current Practical Commands

Recommended commands at the current stage:

    go build ./...

    go test ./...

    SCAVO_TEST_POSTGRES_URL=postgres://postgres:postgres@localhost:5432/scavo_exchange?sslmode=disable \
    go test ./internal/modules/user -run TestPostgresRepository_UpsertDevUser -v

    SCAVO_TEST_POSTGRES_URL=postgres://postgres:postgres@localhost:5432/scavo_exchange?sslmode=disable \
    go test ./internal/modules/user -run TestPostgresRepository_GetByID -v

    SCAVO_POSTGRES_URL=postgres://postgres:postgres@localhost:5432/scavo_exchange?sslmode=disable \
    ./scripts/migrate.sh status

    ./scripts/smoke_login.sh

Example authenticated identity smoke after login:

    TOKEN=$(curl -s -X POST http://localhost:8080/auth/login \
      -H 'Content-Type: application/json' \
      -d '{"email":"test@scavo.exchange","password":"dev"}' | jq -r '.access_token')

    curl -s http://localhost:8080/auth/me \
      -H "Authorization: Bearer $TOKEN"

Example authenticated session smoke after login:

    TOKEN=$(curl -s -X POST http://localhost:8080/auth/login \
      -H 'Content-Type: application/json' \
      -d '{"email":"test@scavo.exchange","password":"dev"}' | jq -r '.access_token')

    curl -s http://localhost:8080/auth/session \
      -H "Authorization: Bearer $TOKEN"

---

## Testing Growth Direction

Testing should grow in roughly this order:

1. unit and smoke validation
2. service tests
3. repository tests
4. infrastructure-aware integration tests
5. chain integration tests
6. contract/backend end-to-end tests

This order matches the project architecture growth.

---

## What Should Be Validated Early

Before heavy product features are introduced, the backend should gain validation for:

- startup behavior
- config loading
- handler wiring
- auth baseline behavior
- current-user authenticated read behavior
- current-session authenticated read behavior
- token extraction consistency
- auth claims propagation
- health endpoint behavior
- readiness behavior
- dependency failure visibility
- migration reproducibility
- first repository readiness

These validations are part of making Stage 0 useful in practice.

---

## Testability Principles

The architecture should support testing by design.

Important principles:

- services should depend on interfaces where appropriate
- handlers should remain thin
- repositories should isolate persistence details
- external integrations should be wrapped behind clients or adapters
- configuration should be injectable or overridable in controlled ways
- side effects should be explicit

These principles reduce friction when testing is implemented.

---

## Environment-Aware Testing

Some test layers should not depend on full local infrastructure.

Examples:

- unit tests should run quickly without DB or Redis
- service tests should avoid unnecessary real dependencies where mocks or test doubles are sufficient

Other test layers intentionally depend on infrastructure.

Examples:

- repository tests with PostgreSQL
- migration validation with PostgreSQL
- integration tests with running services
- chain integration tests with SCAVIUM-compatible endpoints

The project should keep these categories distinct.

---

## Smoke Validation Direction

Smoke validation is especially important for this project.

A minimal smoke layer should verify:

- app starts successfully
- config loads successfully
- /health responds
- /version responds
- auth baseline wiring works
- persistent login path works when PostgreSQL is enabled
- authenticated identity read works with a valid token
- authenticated session read works with a valid token
- WebSocket endpoint is reachable at a basic level

This is a practical baseline for local development and internal testing.

---

## Regression Coverage Direction

As flows stabilize, they should gain regression protection.

Examples later in the roadmap:

- wallet challenge generation
- wallet signature verification
- portfolio read flow
- quote generation
- transaction tracking
- liquidity support flows
- indexer synchronization logic

The project should not wait until the end to start protecting important flows.

---

## Non-Goals for the Current Stage

This stage does not yet require:

- a complete test suite
- CI-enforced full coverage thresholds
- contract end-to-end automation
- performance benchmarks
- fuzzing
- load testing

Those may come later as implementation matures.

---

## Recommended Next Step

The next recommended step is:

Phase 0.4.4 - Wallet Challenge Contract and Nonce Bootstrap

That phase should define nonce and challenge contracts for wallet-oriented authentication while keeping real signature verification for the following step.