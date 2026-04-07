# Phase 0.5 - User Interaction & Application Surface

## Subphase 0.5.2 - User Metadata (Non-Wallet)

## Objective

Introduce the first authenticated write capability over user metadata, strictly limited to a non-sensitive field (`display_name`), without impacting identity, wallet ownership, or authentication flows.

This subphase completes the minimal read/write application surface initiated in 0.5.1.

## Initial Context

After Phase 0.5.1:

* The backend exposes an authenticated read surface via `GET /auth/me`.
* User identity is wallet-first, durable, and stabilized in Phase 0.4.
* The system persists user fields such as `display_name`, `email`, and timestamps.

However:

* There was no contract to mutate user-owned metadata.
* The application surface was read-only.

## Problem Statement

The backend lacked a safe and minimal mechanism for allowing an authenticated user to update their own non-wallet metadata.

This prevented:

* Basic profile personalization.
* Evolution of user-facing features.
* Alignment with a real application surface.

## Scope

Included:

* Authenticated endpoint `PATCH /auth/me`.
* Update limited to `display_name`.
* Validation and normalization of input.
* Persistence through `user.Repository`.
* Response reuse of `GET /auth/me` contract.
* Minimal test coverage expansion.
* Hardening of request validation.

Explicitly excluded:

* Email mutation.
* Wallet mutation.
* User settings.
* Preferences.
* Profile extensions (avatar, bio, etc.).
* Audit logging.
* Multi-field updates.
* Business rules beyond validation.

## Root Cause Analysis

The backend already had a stable identity model, a persisted `User` entity, and metadata fields available.

But it lacked a write contract aligned with authentication context and a safe way to expose controlled mutation.

This gap emerged naturally after enabling `GET /auth/me`.

## Implementation Summary

Endpoint:

PATCH /auth/me

Request:

{
"display_name": "SCAVO Operator"
}

Response:

Reuses the same shape as `GET /auth/me`.

{
"user": { ... },
"profile": { ... }
}

## Validation Rules

Input normalization:

* `display_name` is trimmed.

Constraints:

* Must not be empty after trim.
* Maximum length: 120 characters (Unicode-aware).

## Error Mapping

* Missing/invalid JSON -> `400 bad_request`
* Unknown fields -> `400 bad_request`
* Trailing JSON -> `400 bad_request`
* Empty `display_name` -> `400 invalid_display_name`
* Too long `display_name` -> `400 display_name_too_long`
* Missing auth -> `401 unauthorized`
* User not found -> `404 user_not_found`

## Hardening Applied

Sentinel errors:

* `ErrEmptyUserID`
* `ErrEmptyDisplayName`
* `ErrDisplayNameTooLong`

This removes dependency on string comparisons.

Unicode-safe validation:

* Validation uses rune count instead of byte count.

Strict JSON decoding:

* `DisallowUnknownFields`
* Rejection of trailing JSON
* Body size limit (4KB)

Extended test coverage:

* Invalid JSON
* Unknown fields
* Trailing payloads
* Empty `display_name`
* Length violations
* `user_not_found`
* Unauthorized access

## Files Affected

Auth module:

* `internal/modules/auth/http_login.go`
* `internal/modules/auth/profile.go`
* `internal/modules/auth/http_handlers_test.go`

User module:

* `internal/modules/user/service.go`
* `internal/modules/user/service_test.go`
* `internal/modules/user/repository.go`
* `internal/modules/user/repository_postgres.go`

## Implementation Characteristics

* Additive
* Backward compatible
* No schema changes
* No breaking changes
* No modification of existing auth flows
* No impact on wallet lifecycle

## Validation

* `go test ./...` passes successfully
* No regressions detected
* Auth flows remain stable
* Wallet linking and verification unaffected

## Release Impact

Low risk:

* Introduces a single controlled write path.
* Does not alter existing contracts.
* Maintains backward compatibility.

## Risks

* `/auth/me` could accumulate unrelated responsibilities in future phases.
* Clients may assume broader edit capabilities than currently supported.

## What It Does Not Solve

* Email updates
* Settings contract
* User preferences
* Profile extensions (avatar, bio)
* Audit history
* Advanced lifecycle flags

## Conclusion

Phase 0.5.2 introduces the smallest safe writable surface for the authenticated user.

It completes the transition from a read-only identity surface (0.5.1) to a minimal read/write user interaction layer without reopening identity, wallet ownership, or authentication design.

This establishes the correct foundation for upcoming phases such as user settings, extended profile metadata, and application-level behavior tied to user context.

---

## Subphase 0.5.3 - User Settings Contract Foundation

## Objective

Introduce the first dedicated authenticated settings contract, separated from the authenticated profile/bootstrap surface and from the core `users` record, without impacting identity, wallet ownership, or authentication flows.

This subphase extends the application surface by creating the first explicit boundary between:

* authenticated profile/bootstrap data
* minimal non-wallet metadata editing
* authenticated user settings

## Initial Context

After Phase 0.5.2:

* The backend already exposes an authenticated read surface via `GET /auth/me`.
* The backend already exposes a minimal authenticated write surface via `PATCH /auth/me`, limited to `display_name`.
* User identity remains wallet-first, durable, and stabilized in Phase 0.4.
* The system persists core user fields such as `display_name`, `email`, and timestamps in `users`.

However:

* There was still no dedicated contract for authenticated user settings.
* There was still no persistence surface separated from `users` for future configuration.
* `/auth/me` remained at risk of absorbing unrelated responsibilities if settings were added there later.

## Problem Statement

The backend lacked a safe and minimal mechanism for exposing authenticated user settings as a distinct application-facing surface.

Without that separation:

* `/auth/me` could accumulate unrelated concerns.
* profile metadata and user configuration could collapse into one mixed contract.
* future settings evolution would likely become coupled to the `users` core record.

## Scope

Included:

* Authenticated endpoint `GET /auth/me/settings`.
* Dedicated `user_settings` persistence foundation.
* Minimal settings response contract.
* Safe default resolution when no persisted settings row exists.
* Separation of settings persistence from `users`.
* Minimal service/repository foundation for future settings evolution.
* Test coverage expansion for the new authenticated settings surface.

Explicitly excluded:

* Settings mutation.
* Concrete preference fields (theme, locale, notifications, etc.).
* Wallet mutation.
* Identity redesign.
* Audit logging/history for settings.
* Multi-surface merging of profile and settings.
* Business rules beyond safe default resolution.

## Root Cause Analysis

The backend already had a stable identity model, a persisted `User` entity, and an authenticated profile/bootstrap surface plus minimal metadata edit surface.

But it still lacked a dedicated contract aligned with authentication context for user configuration, and it lacked a persistence boundary that kept future settings out of the `users` core record.

This gap emerged naturally after enabling:

* `GET /auth/me`
* `PATCH /auth/me`

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

* Reuses authenticated context.
* Returns persisted settings when present.
* Returns safe defaults when no settings row exists.
* Does not create or mutate settings during read.

## Contract Characteristics

Settings response fields:

* `user_id`
* `version`
* `preferences`

Contract guarantees:

* `version` is explicit from the first release of the settings surface.
* `preferences` is always returned as an object, never `null`.
* Response remains valid even when no settings row exists yet.

## Persistence Model

New table:

user_settings

Fields:

* `user_id`
* `preferences`
* `created_at`
* `updated_at`

Characteristics:

* 1:1 relationship with durable user identity
* decoupled from `users`
* extensible through `JSONB`
* no forced concrete setting fields in this subphase

## Default Resolution Strategy

When no `user_settings` row exists for the authenticated user:

* no implicit insert is performed
* no side effect is triggered
* service returns a default in-memory settings object

Default behavior:

{
"settings": {
"user_id": "user_123",
"version": 1,
"preferences": {}
}
}

This keeps the first settings contract:

* read-only
* deterministic
* backward compatible
* operationally simple

## Error Mapping

* Missing auth -> `401 unauthorized`
* Settings service unavailable -> `500 auth_service_error`
* Unexpected settings load failure -> `500 auth_service_error`

## Hardening Applied

Dedicated module separation:

* `internal/modules/usersettings/model.go`
* `internal/modules/usersettings/repository.go`
* `internal/modules/usersettings/repository_postgres.go`
* `internal/modules/usersettings/service.go`

Contract hardening:

* explicit settings response envelope
* explicit contract version
* `preferences` normalization to `{}`

Behavior hardening:

* no implicit row creation on read
* no mutation side effect inside `GET /auth/me/settings`
* no coupling of settings to `PATCH /auth/me`

Extended test coverage:

* authenticated settings read with default resolution
* authenticated settings read with persisted settings
* unauthorized access
* unavailable settings service

## Files Affected

Migrations:

* `migrations/000010_user_settings.sql`

Settings module:

* `internal/modules/usersettings/model.go`
* `internal/modules/usersettings/repository.go`
* `internal/modules/usersettings/repository_postgres.go`
* `internal/modules/usersettings/service.go`
* `internal/modules/usersettings/service_test.go`

App / routing / auth surface:

* `internal/app/app.go`
* `internal/core/httpx/router.go`
* `internal/modules/auth/http_login.go`
* `internal/modules/auth/http_handlers_test.go`

## Implementation Characteristics

* Additive
* Backward compatible
* Minimal schema extension
* No breaking changes
* No modification of existing wallet/auth flows
* No impact on wallet lifecycle
* No reopening of Phase 0.4
* Clear separation between metadata and settings

## Validation

* `go test ./...` passes successfully
* Existing auth flows remain stable
* Existing `/auth/me` and `/auth/session` behavior remains unchanged
* Wallet linking, merge, primary switch, detach, and inventory contracts remain unaffected

## Release Impact

Low risk:

* Introduces a single controlled authenticated read path.
* Introduces one small dedicated persistence table.
* Does not alter existing profile or wallet contracts.
* Maintains backward compatibility.

## Risks

* Future phases could still overuse `preferences` as an unstructured bag if contract discipline is not preserved.
* Clients may assume settings mutation exists even though 0.5.3 is read-only.
* Future settings fields may require typed validation if product requirements become stricter.

## What It Does Not Solve

* Settings mutation
* Theme / locale / notification preferences
* Typed settings validation rules
* Settings audit history
* Advanced profile metadata
* Email updates
* Wallet lifecycle extensions
* Identity/provider expansion

## Conclusion

Phase 0.5.3 introduces the smallest safe dedicated settings surface for the authenticated user.

It extends the application-facing layer opened in 0.5.1 and 0.5.2 by separating profile metadata from user configuration without reopening identity, wallet ownership, or authentication design.

This establishes the correct foundation for upcoming phases such as writable settings, typed preferences, and richer application behavior tied to authenticated user context.

---

## Subphase 0.5.4 - User Settings Mutation

## Objective

Introduce the first authenticated write capability over user settings, extending the settings contract introduced in 0.5.3 while preserving separation from profile metadata and identity.

## Initial Context

After Phase 0.5.3:

* The backend exposes `GET /auth/me/settings`.
* Settings are persisted in `user_settings`.
* The contract is read-only.

## Problem Statement

The backend lacked a mutation surface for user settings, preventing:

* user configuration persistence
* evolution of application behavior tied to preferences

## Scope

Included:

* `PATCH /auth/me/settings`
* merge-based updates for `preferences`
* partial mutation without overwrite
* persistence through `user_settings`

Explicitly excluded:

* schema enforcement
* typed preferences
* audit/history
* identity or wallet changes

## Implementation Summary

Endpoint:

PATCH /auth/me/settings

Behavior:

* merges incoming `preferences`
* preserves existing keys
* no destructive overwrite

## Conclusion

Phase 0.5.4 extends the settings surface to support mutation while preserving flexibility and backward compatibility.

---

## Subphase 0.5.5.1 - User Settings Hardening

## Objective

Harden the settings mutation contract to prevent structural drift without introducing rigid schema enforcement.

## Scope

Included:

* recursive normalization
* rejection of null values
* rejection of invalid values
* key trimming
* shape compatibility protection

## Conclusion

Phase 0.5.5.1 stabilizes the settings contract while preserving flexibility, enabling safe future evolution of user preferences.


## Subphase 0.5.5.2 - User Settings Deep Merge Preservation & Known-Branch Semantics

## Objective

Preserve nested settings branches during partial mutation and introduce minimal structural semantics for already-known settings namespaces without converting the settings contract into a rigid schema.

## Scope

Included:

* deep merge preservation for nested objects
* sibling-key preservation during partial nested updates
* recursive shape compatibility enforcement during deep merge
* object-only semantics for known top-level branches: `notifications`, `preferences`, `ui`

Explicitly excluded:

* strict schema enforcement
* typed settings catalogs
* migration/versioning
* hard rejection of unknown top-level keys

## Conclusion

Phase 0.5.5.2 makes `PATCH /auth/me/settings` behave like a true non-destructive partial update for nested object settings while keeping the overall settings surface flexible and application-oriented.


## Subphase 0.5.5.3 - User Settings Contract Surface Stabilization

After 0.5.5.2, the settings mutation path already behaved as a non-destructive deep-merge surface for nested object branches. The next remaining stabilization gap was not internal merge safety, but the visibility of the resource state returned to authenticated clients.

This subphase keeps the existing authenticated settings endpoints unchanged:

* `GET /auth/me/settings`
* `PATCH /auth/me/settings`

The returned `settings` resource now exposes persisted timestamp metadata when available:

* `created_at`
* `updated_at`

Those fields are omitted when the settings view is being resolved from safe defaults and no persisted timestamps exist yet.

This makes the settings contract more explicit and self-descriptive for frontend consumers without introducing schema-heavy settings governance, optimistic locking, or version negotiation.
