# Phase 0.5 - User Interaction & Application Surface

## Subphase 0.5.3 - User Settings Contract Foundation

## Objective

Introduce the first dedicated authenticated user settings contract, separated from profile metadata and durable identity, without impacting wallet ownership, authentication flows, or the minimal metadata mutation introduced in 0.5.2.

This subphase extends the application surface by creating the first formal boundary between:

- authenticated profile bootstrap
- minimal non-wallet metadata editing
- authenticated user configuration

## Initial Context

After Phase 0.5.2:

- The backend already exposes an authenticated read surface via `GET /auth/me`.
- The backend already exposes a minimal authenticated write surface via `PATCH /auth/me`, limited to `display_name`.
- User identity remains wallet-first, durable, and stabilized in Phase 0.4.
- The system persists core user fields such as `display_name`, `email`, and timestamps in `users`.

However:

- There was still no dedicated contract for user settings.
- There was still no persistence surface separated from `users` for future configuration.
- The application surface still lacked a formal separation between profile metadata and user configuration.

## Problem Statement

The backend lacked a safe and minimal mechanism for exposing authenticated user settings as a distinct application-facing surface.

Without that separation:

- `/auth/me` risked accumulating unrelated responsibilities.
- profile metadata and user configuration could collapse into one mixed contract.
- future settings evolution would likely become coupled to the `users` core record.

## Scope

Included:

- Authenticated endpoint `GET /auth/me/settings`.
- Dedicated `user_settings` persistence foundation.
- Minimal settings response contract.
- Safe default resolution when no persisted settings row exists.
- Separation of settings persistence from `users`.
- Minimal service/repository foundation for future settings evolution.
- Test coverage expansion for the new authenticated settings surface.

Explicitly excluded:

- Settings mutation.
- Concrete preference fields (theme, locale, notifications, etc.).
- Wallet mutation.
- Identity redesign.
- Audit logging/history for settings.
- Multi-surface merging of profile and settings.
- Business rules beyond safe default resolution.

## Root Cause Analysis

The backend already had a stable identity model, a persisted `User` entity, and an authenticated profile bootstrap plus minimal metadata edit surface.

But it still lacked a dedicated contract aligned with authentication context for user configuration, and it lacked a persistence boundary that kept future settings out of the `users` core record.

This gap emerged naturally after enabling:

- `GET /auth/me`
- `PATCH /auth/me`

## Implementation Summary

Endpoint:

GET /auth/me/settings

Response:

{
  "settings": {
    "user_id": "user_123",
    "version": 1,
    "preferences": {}
  }
}

Behavior:

- Reuses authenticated context.
- Returns persisted settings when present.
- Returns safe defaults when no settings row exists.
- Does not create or mutate settings during read.

## Contract Characteristics

Settings response fields:

- `user_id`
- `version`
- `preferences`

Contract guarantees:

- `version` is explicit from the first release of the settings surface.
- `preferences` is always returned as an object, never `null`.
- response remains valid even when no settings row exists yet.

## Persistence Model

New table:

user_settings

Fields:

- `user_id`
- `preferences`
- `created_at`
- `updated_at`

Characteristics:

- 1:1 relationship with durable user identity
- decoupled from `users`
- extensible through JSONB
- no forced concrete setting fields in this subphase

## Default Resolution Strategy

When no `user_settings` row exists for the authenticated user:

- no implicit insert is performed
- no side effect is triggered
- service returns a default in-memory settings object

Default behavior:

{
  "settings": {
    "user_id": "user_123",
    "version": 1,
    "preferences": {}
  }
}

This keeps the first settings contract:

- read-only
- deterministic
- backward compatible
- operationally simple

## Error Mapping

- Missing auth -> 401 unauthorized
- Settings service unavailable -> 500 auth_service_error
- Unexpected settings load failure -> 500 auth_service_error

## Hardening Applied

Dedicated module separation:

- internal/modules/usersettings/model.go
- internal/modules/usersettings/repository.go
- internal/modules/usersettings/repository_postgres.go
- internal/modules/usersettings/service.go

Contract hardening:

- explicit settings response envelope
- explicit contract version
- preferences normalization to {}

Behavior hardening:

- no implicit row creation on read
- no mutation side effect inside GET /auth/me/settings
- no coupling of settings to PATCH /auth/me

Extended test coverage:

- authenticated settings read with default resolution
- authenticated settings read with persisted settings
- unauthorized access
- unavailable settings service

## Files Affected

Migrations:

- migrations/000010_user_settings.sql

Settings module:

- internal/modules/usersettings/model.go
- internal/modules/usersettings/repository.go
- internal/modules/usersettings/repository_postgres.go
- internal/modules/usersettings/service.go
- internal/modules/usersettings/service_test.go

App / routing / auth surface:

- internal/app/app.go
- internal/core/httpx/router.go
- internal/modules/auth/http_login.go
- internal/modules/auth/http_handlers_test.go

## Implementation Characteristics

- Additive
- Backward compatible
- Minimal schema extension
- No breaking changes
- No modification of existing wallet/auth flows
- No impact on wallet lifecycle
- No reopening of Phase 0.4
- Clear separation between metadata and settings

## Validation

- go test ./... should pass successfully after the full 0.5.3 implementation is applied
- Existing auth flows should remain stable
- Existing /auth/me and /auth/session behavior should remain unchanged
- Wallet linking, merge, primary switch, detach, and inventory contracts should remain unaffected

## Release Impact

Low risk:

- Introduces a single controlled authenticated read path.
- Introduces one small dedicated persistence table.
- Does not alter existing profile or wallet contracts.
- Maintains backward compatibility.

## Risks

- Future phases could still overuse preferences as an unstructured bag if contract discipline is not preserved.
- Clients may assume settings mutation exists even though 0.5.3 is read-only.
- Future settings fields may require typed validation if product requirements become stricter.

## What It Does Not Solve

- Settings mutation
- Theme / locale / notification preferences
- Typed settings validation rules
- Settings audit history
- Advanced profile metadata
- Email updates
- Wallet lifecycle extensions
- Identity/provider expansion

## Conclusion

Phase 0.5.3 introduces the smallest safe dedicated settings surface for the authenticated user.

It extends the application-facing layer opened in 0.5.1 and 0.5.2 by separating profile metadata from user configuration without reopening identity, wallet ownership, or authentication design.

This establishes the correct foundation for upcoming phases such as writable settings, typed preferences, and richer application behavior tied to authenticated user context.