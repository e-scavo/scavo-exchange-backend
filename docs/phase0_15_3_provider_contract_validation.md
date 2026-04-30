# Phase 0.15.3 — Provider Contract Validation

---

## Status

Completed.

---

## Context Inherited

Phase 0.15.3 starts from the validated state produced by the previous Contract Hardening & Freeze subphases.

0.15.1 established the HTTP contract baseline by auditing the registered route surface, documenting the public endpoint families, and preserving the observed route behavior without modifying code.

0.15.2 aligned the public error envelope by making `error.details` a stable object field in canonical HTTP error responses while preserving existing status codes, error codes, messages and handler decisions.

Together, those subphases made it possible to validate provider contracts without mixing in route discovery, error-shape corrections or response normalization work.

---

## Real Problem

The backend already had a clean provider-oriented architecture.

The HTTP layer consumed the auth module through a handler-facing provider boundary composed of:

- `AuthSessionProvider`
- `AuthenticatedAccountProvider`
- `AuthWalletProvider`
- `AuthProvider`

The auth application also consumed user and usersettings modules through minimal domain contracts:

- `authdomain.UserProvider`
- `authdomain.UserSettingsProvider`

However, the most important contract relationships were still validated primarily by usage.

That meant a future method signature drift could break the boundary between HTTP, auth, user and usersettings without being explicitly described as a provider-contract violation at the owning seam.

---

## Decision Taken

The decision for 0.15.3 was to validate provider contracts through compile-time assertions.

This keeps the phase aligned with Contract Hardening & Freeze:

- no new feature
- no new endpoint
- no response-shape change
- no error-envelope change
- no business-rule change
- no repository behavior change

The chosen hardening point is `internal/modules/auth/application.go`, because that file owns the public auth application wrapper consumed by HTTP handlers and wires the cross-module providers used by auth.

---

## Concrete Change

The following compile-time assertions were added:

```go
var (
	_ AuthProvider                    = (*Application)(nil)
	_ authdomain.UserProvider         = (*usermod.Service)(nil)
	_ authdomain.UserSettingsProvider = (*usersettingsmod.Service)(nil)
)
```

These assertions prove that:

- `*Application` satisfies the handler-facing `AuthProvider` contract
- `*user.Service` satisfies the minimal user provider contract consumed by auth
- `*usersettings.Service` satisfies the minimal usersettings provider contract consumed by auth

---

## Contract Impact

### HTTP to Provider

The HTTP handler boundary remains explicit.

Handlers continue to depend on `AuthProvider`, not on transport-specific implementation details.

The compile-time assertion now guarantees that the concrete auth application wrapper remains compatible with that boundary.

### Provider to Application

The auth wrapper continues to delegate into the application layer without changing method inputs, outputs or error propagation.

The assertion prevents method-signature drift from silently weakening the provider boundary.

### Auth to User

The user dependency remains restricted to the minimal `authdomain.UserProvider` interface.

The concrete user service is now compile-time validated against that contract.

### Auth to User Settings

The usersettings dependency remains restricted to the minimal `authdomain.UserSettingsProvider` interface.

The concrete usersettings service is now compile-time validated against that contract.

---

## Explicit Non-Changes

0.15.3 does not change:

- registered HTTP routes
- HTTP methods
- status codes
- response payloads
- error codes
- error messages
- error envelope shape
- authorization behavior
- token behavior
- wallet behavior
- user profile behavior
- settings behavior
- repository behavior
- database behavior

---

## Observable Impact

Provider drift now fails at compile time.

If a future change removes, renames or changes a method required by `AuthProvider`, `authdomain.UserProvider` or `authdomain.UserSettingsProvider`, the backend build fails before the drift reaches handlers, frontend consumers or production runtime behavior.

This makes the internal provider boundary explicit and enforceable while preserving the public contract baseline established in 0.15.1 and 0.15.2.

---

## Validation Notes

The implementation is intentionally narrow.

The expected validation command remains:

```bash
go test ./...
```

The change is compile-time oriented. A successful full test run validates that the asserted concrete types still satisfy the provider contracts and that no existing package-level behavior was broken.

---

## Handoff to 0.15.4

The next subphase is:

**0.15.4 — Response Schema Normalization**

It must start from these validated baselines:

- HTTP route baseline from 0.15.1
- canonical error envelope from 0.15.2
- provider-boundary compile-time validation from 0.15.3

0.15.4 must not introduce new features or business behavior. It must focus only on response schema consistency.

---

## Final State

Phase 0.15.3 leaves the backend with:

- explicit handler-facing provider contract validation
- explicit cross-module provider contract validation
- no runtime behavior change
- no public response change
- no error-contract regression
- a safer base for response schema normalization

---

## End of Document
