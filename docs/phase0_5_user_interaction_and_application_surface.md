# Phase 0.5 - User Interaction & Application Surface

## Subphase 0.5.2 - User Metadata (Non-Wallet)

## Objective

Introduce the first authenticated write capability over user metadata, strictly limited to a non-sensitive field (`display_name`), without impacting identity, wallet ownership, or authentication flows.

This subphase completes the minimal read/write application surface initiated in 0.5.1.

## Initial Context

After Phase 0.5.1:

- The backend exposes an authenticated read surface via `GET /auth/me`.
- User identity is wallet-first, durable, and stabilized in Phase 0.4.
- The system persists user fields such as `display_name`, `email`, and timestamps.

However:

- There was no contract to mutate user-owned metadata.
- The application surface was read-only.

## Problem Statement

The backend lacked a safe and minimal mechanism for allowing an authenticated user to update their own non-wallet metadata.

This prevented:

- Basic profile personalization.
- Evolution of user-facing features.
- Alignment with a real application surface.

## Scope

Included:

- Authenticated endpoint `PATCH /auth/me`.
- Update limited to `display_name`.
- Validation and normalization of input.
- Persistence through `user.Repository`.
- Response reuse of `GET /auth/me` contract.
- Minimal test coverage expansion.
- Hardening of request validation.

Explicitly excluded:

- Email mutation.
- Wallet mutation.
- User settings.
- Preferences.
- Profile extensions (avatar, bio, etc.).
- Audit logging.
- Multi-field updates.
- Business rules beyond validation.

## Root Cause Analysis

The backend already had a stable identity model, a persisted `User` entity, and metadata fields available.

But it lacked a write contract aligned with authentication context and a safe way to expose controlled mutation.

This gap emerged naturally after enabling `GET /auth/me`.

## Implementation Summary

Endpoint:

```http
PATCH /auth/me
```

Request:

```json
{
  "display_name": "SCAVO Operator"
}
```

Response:

Reuses the same shape as `GET /auth/me`.

```json
{
  "user": { ... },
  "profile": { ... }
}
```

## Validation Rules

Input normalization:

- `display_name` is trimmed.

Constraints:

- Must not be empty after trim.
- Maximum length: 120 characters (Unicode-aware).

## Error Mapping

- Missing/invalid JSON -> `400 bad_request`
- Unknown fields -> `400 bad_request`
- Trailing JSON -> `400 bad_request`
- Empty `display_name` -> `400 invalid_display_name`
- Too long `display_name` -> `400 display_name_too_long`
- Missing auth -> `401 unauthorized`
- User not found -> `404 user_not_found`

## Hardening Applied

Sentinel errors:

- `ErrEmptyUserID`
- `ErrEmptyDisplayName`
- `ErrDisplayNameTooLong`

This removes dependency on string comparisons.

Unicode-safe validation:

- Validation uses rune count instead of byte count.

Strict JSON decoding:

- `DisallowUnknownFields`
- Rejection of trailing JSON
- Body size limit (4KB)

Extended test coverage:

- Invalid JSON
- Unknown fields
- Trailing payloads
- Empty `display_name`
- Length violations
- `user_not_found`
- Unauthorized access

## Files Affected

Auth module:

- `internal/modules/auth/http_login.go`
- `internal/modules/auth/profile.go`
- `internal/modules/auth/http_handlers_test.go`

User module:

- `internal/modules/user/service.go`
- `internal/modules/user/service_test.go`
- `internal/modules/user/repository.go`
- `internal/modules/user/repository_postgres.go`

## Implementation Characteristics

- Additive
- Backward compatible
- No schema changes
- No breaking changes
- No modification of existing auth flows
- No impact on wallet lifecycle

## Validation

- `go test ./...` passes successfully
- No regressions detected
- Auth flows remain stable
- Wallet linking and verification unaffected

## Release Impact

Low risk:

- Introduces a single controlled write path.
- Does not alter existing contracts.
- Maintains backward compatibility.

## Risks

- `/auth/me` could accumulate unrelated responsibilities in future phases.
- Clients may assume broader edit capabilities than currently supported.

## What It Does Not Solve

- Email updates
- Settings contract
- User preferences
- Profile extensions (avatar, bio)
- Audit history
- Advanced lifecycle flags

## Conclusion

Phase 0.5.2 introduces the smallest safe writable surface for the authenticated user.

It completes the transition from a read-only identity surface (0.5.1) to a minimal read/write user interaction layer without reopening identity, wallet ownership, or authentication design.

This establishes the correct foundation for upcoming phases such as user settings, extended profile metadata, and application-level behavior tied to user context.