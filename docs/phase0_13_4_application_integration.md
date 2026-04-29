# Phase 0.13.4 — Application Integration

## Status

COMPLETED

## Purpose

Integrate the Provider Layer implemented in 0.13.3 into the runtime application and HTTP routing path without changing public API behavior.

## Scope

0.13.4 focuses on runtime composition and handler wiring.

Included:

- align HTTP router construction with the consolidated auth provider boundary
- reduce production handler wiring to provider-level dependencies
- preserve fallback compatibility for transitional tests
- keep public routes, versioned routes, middleware and response contracts unchanged
- document the resulting runtime flow

Excluded:

- public API changes
- business logic changes
- authorization policy changes
- repository changes
- database or migration changes
- validation closure beyond this integration step

## Implementation Summary

Runtime auth HTTP wiring now uses the consolidated provider boundary directly.

Before this integration step, router construction still carried multiple auth implementation details into the HTTP handler surface:

- user service
- user settings service
- wallet challenge store
- wallet identity store
- challenge TTL
- public base URL

After this integration step, production router construction creates auth HTTP handlers from the provider boundary only.

The resulting production path is:

```text
HTTP → Provider → Application → Domain → Repository
```

## Code Changes

### internal/core/httpx/router.go

`RouterParams` was reduced so auth runtime wiring depends on:

- token service
- auth provider
- status service
- shared router dependencies

Auth handlers are now created through provider-level construction instead of field-by-field implementation wiring.

### internal/modules/auth/http_login.go

A provider-oriented HTTP handler constructor was introduced:

```go
func NewHTTPHandlers(provider AuthProvider) HTTPHandlers
```

The constructor records the expected production wiring shape for HTTP handlers.

Legacy fields remain on `HTTPHandlers` only for transitional test compatibility and fallback behavior. They no longer define production router wiring.

### internal/app/app.go

The application composition root still builds the concrete auth provider, but router construction now receives only the consolidated provider boundary for auth HTTP handling.

## Compatibility Notes

0.13.4 preserves:

- existing auth routes
- `/api/v1` mirrored routes
- request payloads
- response payloads
- error envelope behavior
- authorization middleware behavior
- wallet challenge and wallet management behavior
- existing tests and compatibility fallbacks

## Validation

The expected validation commands are:

```bash
make build
go test ./...
```

## Documentation Notes

All general Markdown trunk documents were reviewed under the Phase 0.13 documentation rules. Phase-specific documents were reviewed and updated only when directly related to the current phase progression.

## Result

0.13.4 completes the application integration step for the Provider Layer.

The repository advanced to:

```text
0.13.5 — Validation & Compatibility
```

0.13.5 has since completed the validation and compatibility checkpoint.
