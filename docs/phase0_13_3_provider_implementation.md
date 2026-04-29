# Phase 0.13.3 — Provider Implementation

## Status

Phase 0.13.3 is **COMPLETED**.

This subphase implements the first concrete Provider Layer boundary for the auth surface while preserving public HTTP behavior and compatibility with the existing application/domain/repository structure.

---

## Objective

Implement explicit provider contracts and route auth HTTP handlers through those contracts instead of keeping orchestration spread directly across handler-level dependencies.

The implementation follows the design locked in `docs/phase0_13_2_provider_interface_design.md` and the inventory from `docs/phase0_13_1_provider_inventory.md`.

---

## Scope

Included:

- auth provider interface definitions
- composite auth provider contract
- provider-backed handler dependency
- composition-root provider construction
- account profile provider methods
- user settings provider methods
- public wallet bootstrap provider methods
- authenticated wallet provider routing
- cross-module contract expansion only where required by real provider implementation

Excluded:

- public API route changes
- request/response payload changes
- authorization middleware changes
- API versioning changes
- repository schema changes
- business rule changes

---

## Implemented Provider Boundary

The auth module now exposes explicit provider-facing contracts:

- `AuthSessionProvider`
- `AuthenticatedAccountProvider`
- `AuthWalletProvider`
- `AuthProvider`

`AuthProvider` is the composite handler-facing boundary used by HTTP handlers.

The concrete compatibility implementation is the existing root auth application facade, backed by `internal/modules/auth/app.Application`.

---

## Files Changed

### Code

- `internal/app/app.go`
- `internal/core/httpx/router.go`
- `internal/modules/auth/provider.go`
- `internal/modules/auth/application.go`
- `internal/modules/auth/http_login.go`
- `internal/modules/auth/http_bootstrap.go`
- `internal/modules/auth/http_wallet.go`
- `internal/modules/auth/http_wallet_list.go`
- `internal/modules/auth/app/application.go`
- `internal/modules/auth/app/response_types.go`
- `internal/modules/auth/domain/user_contract.go`
- `internal/modules/auth/domain/usersettings_contract.go`

### Documentation

- `README.md`
- `docs/index.md`
- `docs/roadmap.md`
- `docs/phase-status.md`
- `docs/handoff/backend-status.md`
- `docs/phase0_13_provider_layer_consolidation.md`
- `docs/phase0_13_3_provider_implementation.md`

---

## Provider Interfaces

The new provider contract is intentionally handler-facing and narrow.

### Session Boundary

Owns:

- login
- current user view
- session view
- bootstrap view

### Authenticated Account Boundary

Owns:

- profile update
- settings retrieval
- settings update

### Wallet Boundary

Owns:

- public wallet challenge
- public wallet verification
- wallet listing
- wallet linking
- account merge
- primary wallet selection
- wallet detach check
- wallet detach execution

---

## Composition Root Wiring

`internal/app/app.go` now constructs the auth provider once from the already existing services and stores.

`internal/core/httpx/router.go` receives that provider through router params and injects it into `auth.HTTPHandlers`.

This reduces ad-hoc handler application construction while preserving compatibility fields for existing tests and transitional wiring.

---

## Handler Alignment

Auth HTTP handlers now route provider-owned operations through `h.AuthProvider()`.

This applies to:

- login
- me
- update me
- settings read
- settings update
- session
- bootstrap
- wallet challenge
- wallet verify
- wallet listing
- wallet link
- account merge
- primary wallet
- wallet detach

Transport handlers remain responsible for:

- request decoding
- body-size limits
- auth claim extraction
- status code selection
- response writing
- public error-envelope selection

Providers own orchestration and application flow delegation.

---

## Cross-Module Contract Adjustment

The auth domain cross-module contracts were expanded only where implementation required real provider ownership.

`UserProvider` now includes:

- wallet-user resolution
- display-name update

`UserSettingsProvider` now includes:

- settings update

This avoids reintroducing direct handler-level access to user and user settings services while preserving module boundaries.

---

## Compatibility

0.13.3 does not change:

- public route paths
- HTTP methods
- request JSON contracts
- response JSON contracts
- middleware ordering
- API versioning
- standardized error envelope shape
- repository schema
- persistence behavior

---

## Validation

`gofmt` was applied to all changed Go files.

`go test ./...` could not be completed in this execution environment because the local Go tool attempted to download Go 1.25.0 from `proxy.golang.org`, but outbound DNS/network access was unavailable.

The expected user-side validation command remains:

```bash
go test ./...
```

---

## Result

0.13.3 establishes the first concrete Provider Layer implementation for the auth surface and prepares the repository for:

```text
0.13.4 — Application Integration
```

0.13.4 should focus on reducing remaining transitional compatibility wiring and confirming that application/provider integration is clean across handlers, router construction and tests.
