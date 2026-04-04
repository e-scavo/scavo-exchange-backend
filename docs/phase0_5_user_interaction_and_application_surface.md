# Phase 0.5 — User Interaction & Application Surface

## Objective

Open the first authenticated application-facing read surface on top of the durable identity and wallet ownership model stabilized in Phase 0.4.

## Initial Context

Phase 0.4 closed with:

- durable user identity
- wallet-backed JWT authentication
- wallet ownership persistence
- wallet inventory reads
- strict challenge-purpose enforcement

The backend already exposed `GET /auth/me`, `GET /auth/session`, and `GET /auth/wallets`, but those endpoints still reflected internal auth concerns more than a small application bootstrap surface.

## Problem Statement

Authenticated clients still needed to compose multiple calls or infer meaning from lower-level auth responses to bootstrap the application surface for the current user.

The gap was not authentication correctness anymore. The gap was an additive, stable, authenticated profile surface that could expose:

- current durable user identity
- current auth method and wallet session context
- primary wallet summary
- lightweight owned-wallet projection
- wallet counters

without reworking the wallet inventory contract or overloading the session contract.

## Scope

Phase 0.5.1 includes:

- additive authenticated profile projection for `GET /auth/me`
- primary-wallet summary derived from existing wallet ownership state
- wallet counters derived from the existing ownership store
- handler-level coverage for wallet-backed and non-wallet-backed authenticated sessions
- documentation alignment for the new application-facing user surface

Phase 0.5.1 does not include:

- user settings
- non-wallet profile editing
- account flags
- audit/event history
- exchange-specific user logic
- session persistence or refresh tokens

## Root Cause Analysis

The backend solved identity durability before it solved client-facing identity ergonomics.

`/auth/session` already exposed authenticated claims, but it was claim/session oriented. `/auth/wallets` already exposed inventory, but it was an inventory endpoint. `/auth/me` only returned the user model.

That left no small, additive, application-facing surface that represented “who is the authenticated user right now, and what wallet context should the client immediately know about?”

## Files Affected

### Code

- `internal/modules/auth/http_login.go`
- `internal/modules/auth/profile.go`
- `internal/modules/auth/http_handlers_test.go`

### Documentation

- `README.md`
- `docs/index.md`
- `docs/architecture.md`
- `docs/testing.md`
- `docs/roadmap.md`
- `docs/phase-status.md`
- `docs/handoff/backend-status.md`
- `docs/phase0_5_user_interaction_and_application_surface.md`

## Implementation Characteristics

`GET /auth/me` remains backward compatible by preserving the existing top-level `user` field.

The endpoint now also returns an additive `profile` object with:

- `user_id`
- `auth_method`
- `wallet_id`
- `wallet_address`
- `chain`
- `primary_wallet`
- `wallets`
- `wallet_count`
- `active_wallet_count`
- `detached_wallet_count`
- `has_wallet_session`

The profile surface is derived from already stabilized primitives:

- JWT claims from the authenticated request
- durable user resolution already used by auth/session flows
- wallet ownership reads from the existing wallet identity store

No new persistence layer, ownership rule, or mutation contract was introduced.

## Validation

Validated at code level through handler coverage for:

- wallet-backed authenticated `GET /auth/me`
- password/dev authenticated `GET /auth/me` without owned wallets
- missing-claims unauthorized behavior

Full `go test ./...` execution could not be completed inside this environment because the repository requires Go `1.25.0` and the container attempted to download that toolchain from `proxy.golang.org`, which is blocked in this session.

## Release Impact

Impact is additive and low risk:

- existing `GET /auth/me` consumers can keep reading `user`
- new consumers can bootstrap directly from `profile`
- `/auth/session` and `/auth/wallets` remain available with their existing roles

## Risks

- clients might start treating `profile.wallets` as a replacement for the fuller `/auth/wallets` inventory contract
- future profile expansion could overload `/auth/me` if additive boundaries are not preserved

## What it does NOT solve

This phase does not solve:

- editable profile metadata
- user settings
- soft account flags
- audit history
- refresh tokens
- persistent sessions
- exchange domain reads

## Conclusion

Phase 0.5.1 establishes the first authenticated application surface on top of the identity work completed in Phase 0.4.

The backend now exposes a small, additive, user-facing profile read on `GET /auth/me` without reopening auth lifecycle design or changing wallet ownership rules.
