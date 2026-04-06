# Phase 0.5.4 — User Settings Mutation

---

## Objective

Introduce a safe, minimal, and controlled mutation mechanism for user settings, enabling partial updates without compromising existing persisted preferences or architectural boundaries established in previous phases.

---

## Initial Context

By the end of Phase 0.5.3, the system already supports:

* Authenticated user identity (wallet-first model)
* Stable JWT-based session handling
* Profile read and minimal mutation:

  * GET /auth/me
  * PATCH /auth/me (display_name only)
* User settings read:

  * GET /auth/me/settings
* Persistent storage via `user_settings`
* Default resolution without side effects

However, user settings are strictly **read-only**.

---

## Problem Statement

The system lacks a mechanism to **persist user-defined preferences**.

This creates the following limitations:

* Frontend cannot store user personalization
* No persistence for UI/UX state
* No evolution path for feature-level preferences
* Settings contract is incomplete (read-only)

At the same time, introducing mutation incorrectly could:

* Overwrite existing settings unintentionally
* Break separation between profile and settings
* Introduce side effects or hidden defaults
* Add unnecessary complexity

---

## Scope

This phase introduces:

* Authenticated endpoint:

  * PATCH /auth/me/settings
* Partial update semantics (merge-based)
* Minimal validation
* Safe persistence
* Unit tests

This phase explicitly excludes:

* Full schema validation
* Settings versioning
* Complex preference systems
* Cross-field validation logic
* Default injection during write

---

## Root Cause Analysis

The absence of a mutation layer for user settings stems from:

* Intentional deferral during Phase 0.5.3
* Need to first stabilize:

  * Settings read model
  * Persistence contract
  * Default resolution behavior

Without a defined mutation contract:

* Any write implementation risks inconsistency
* Overwrites become likely
* Data integrity cannot be guaranteed

---

## Files Affected

* internal/core/httpx/router.go
* internal/modules/auth/http_login.go
* internal/modules/auth/http_handlers_test.go
* internal/modules/usersettings/service.go
* internal/modules/usersettings/service_test.go
* internal/modules/usersettings/repository.go
* internal/modules/usersettings/repository_postgres.go

---

## Implementation Characteristics

### Endpoint

PATCH /auth/me/settings

* Requires valid JWT
* Uses same authentication context as /auth/me

---

### Request Format

```json
{
  "preferences": {
    "key": "value"
  }
}
```

---

### Update Strategy

* Merge-based update (NOT replace)
* Only provided keys are updated
* Existing keys remain untouched

Example:

Input:

```json
{
  "preferences": {
    "theme": "dark"
  }
}
```

Result:

```json
{
  "theme": "dark",
  "notifications": true
}
```

---

### Validation Rules

Minimal and intentional:

* Request must be valid JSON
* `preferences` must exist
* `preferences` must be an object
* `preferences` cannot contain null values
* Unknown fields are rejected

No schema enforcement is applied.

---

### Persistence Behavior

* Upsert on `user_settings`
* If record does not exist → create
* If record exists → merge and update
* `created_at` preserved
* `updated_at` refreshed

---

### Architectural Guarantees

* No coupling with profile (/auth/me)
* No modification to auth flow
* No changes to wallet lifecycle
* No side effects during write
* No implicit defaults injection

---

## Validation

Validated via:

* Unit tests for:

  * Merge correctness
  * Persistence behavior
  * Authentication enforcement
* Manual validation using HTTP client (curl)

All tests pass:

```bash
go test ./...
```

---

## Release Impact

This phase:

* Transitions settings from read-only → read/write
* Enables frontend persistence capabilities
* Does not break any existing endpoint
* Maintains backward compatibility

---

## Risks

* Improper client usage could still overwrite nested structures if misused
* Lack of schema validation may allow inconsistent data shapes (acceptable at this stage)
* Future migrations may require normalization

---

## What it does NOT solve

* Structured preferences system
* Validation schemas
* Per-feature settings isolation
* Settings versioning
* Conflict resolution strategies

---

## Conclusion

Phase 0.5.4 successfully introduces a **minimal, safe, and extensible mutation layer** for user settings.

The implementation respects all architectural constraints established in previous phases while enabling the system to evolve toward richer user personalization features without introducing premature complexity.
