# Phase 0.5 — User Interaction & Application Surface

## Objective

Extend the authenticated user surface opened in Phase 0.5.1 with the smallest safe non-wallet metadata mutation for the authenticated user.

## Initial Context

Phase 0.5.1 already delivered:

- additive authenticated `GET /auth/me`
- durable user resolution through the existing auth/session path
- wallet-aware profile projection for application bootstrap

At that point the backend could expose authenticated identity, but it still could not mutate even the smallest self-scoped user metadata field through an authenticated contract.

## Problem Statement

The backend already persisted `display_name` inside the durable `users` model, but there was no authenticated contract to update it.

That left the application surface asymmetric:

- the current user could be read through `GET /auth/me`
- the current user could not update basic non-wallet metadata

The missing piece was not identity, wallets, or ownership. It was a minimal self-scoped user metadata update contract.

## Scope

Phase 0.5.2 includes:

- authenticated `PATCH /auth/me`
- update of `display_name` only
- repository/service support for persisting the new display name
- handler-level validation for trimmed non-empty values and a bounded maximum length
- response alignment with the existing `MeResponse`
- documentation and testing updates for the new mutation contract

Phase 0.5.2 does not include:

- email updates
- wallet metadata changes
- user settings
- soft account flags
- audit/event history
- exchange-specific user fields

## Root Cause Analysis

The backend had already stabilized the identity model and opened a read surface, but it still lacked a minimal writable self profile contract.

`display_name` was the only clearly non-wallet field already present in the durable user model and safe to expose for authenticated editing without reopening auth design.

## Files Affected

### Code

- `internal/core/httpx/router.go`
- `internal/modules/auth/http_login.go`
- `internal/modules/auth/profile.go`
- `internal/modules/auth/http_handlers_test.go`
- `internal/modules/user/repository.go`
- `internal/modules/user/repository_postgres.go`
- `internal/modules/user/repository_postgres_test.go`
- `internal/modules/user/service.go`
- `internal/modules/user/service_test.go`

### Documentation

- `README.md`
- `docs/roadmap.md`
- `docs/phase-status.md`
- `docs/testing.md`
- `docs/handoff/backend-status.md`
- `docs/phase0_5_user_interaction_and_application_surface.md`

## Implementation Characteristics

`PATCH /auth/me` accepts:

- `display_name`

Validation is intentionally small and explicit:

- request body must decode correctly
- authenticated user context must exist
- `display_name` is trimmed before validation
- empty trimmed value is rejected
- maximum length is limited to `120`

The response remains aligned with `GET /auth/me` by returning:

- top-level `user`
- additive `profile`

Wallet session context, wallet counters, and primary wallet projection remain derived from the existing ownership store and are not changed by this mutation.

## Validation

Validated at code level through:

- service fallback update behavior
- repository-backed service update behavior
- invalid empty display name rejection
- too-long display name rejection
- handler success path for authenticated `PATCH /auth/me`
- handler invalid body rejection
- handler unauthorized behavior
- handler validation failures
- postgres repository update test guarded by `SCAVO_TEST_POSTGRES_URL`

## Release Impact

Impact is additive and low risk:

- authenticated clients can now update a minimal profile field without a new resource tree
- `GET /auth/me` remains backward compatible
- wallet identity and ownership contracts remain unchanged

## Risks

- future profile editing could overload `/auth/me` if unrelated metadata keeps accumulating there
- clients might assume email or settings are also mutable even though this phase intentionally restricts updates to `display_name`

## What it does NOT solve

This phase does not solve:

- email mutation
- settings contract
- user preferences
- profile avatars
- audit history
- richer user lifecycle flags

## Conclusion

Phase 0.5.2 closes the smallest safe writable gap in the authenticated application surface.

The backend now supports both reading the authenticated user bootstrap surface and updating the authenticated user `display_name` without reopening wallet lifecycle or identity ownership design.
