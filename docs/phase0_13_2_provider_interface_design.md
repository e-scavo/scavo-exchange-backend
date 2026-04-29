# Phase 0.13.2 — Provider Interface Design

## Status

Phase 0.13.2 is **COMPLETED**.

This subphase is documentation-only. It translates the Phase 0.13.1 provider inventory into explicit provider interface design decisions for the implementation step.

No Go code is changed by 0.13.2.

---

## Objective

Define narrow, stable provider interfaces that consolidate access to authenticated domain/application behavior without changing public HTTP routes, request payloads, response payloads, authorization behavior, API versioning or business rules.

The design uses the real repository state recorded in `docs/phase0_13_1_provider_inventory.md` and preserves the Phase 0.12 read/write model separation.

---

## Source Scope

This design was derived from the repository state supplied as `scavo-exchange-backend-0.13.1.zip`.

Reviewed source areas:

- `internal/app/app.go`
- `internal/core/httpx/*.go`
- `internal/modules/auth/application.go`
- `internal/modules/auth/http_login.go`
- `internal/modules/auth/http_bootstrap.go`
- `internal/modules/auth/http_wallet.go`
- `internal/modules/auth/http_wallet_list.go`
- `internal/modules/auth/app/application.go`
- `internal/modules/auth/app/service.go`
- `internal/modules/auth/app/wallet_services.go`
- `internal/modules/auth/domain/*.go`
- `internal/modules/auth/mappers/*.go`
- `internal/modules/auth/readmodels/*.go`
- `internal/modules/auth/writemodels/*.go`
- `internal/modules/user/app/service.go`
- `internal/modules/user/domain/*.go`
- `internal/modules/user/mappers/*.go`
- `internal/modules/usersettings/app/service.go`
- `internal/modules/usersettings/domain/*.go`
- `internal/modules/usersettings/mappers/*.go`

---

## Design Principles

1. **No public API drift**
   - Existing routes, methods, request bodies, query parameters, response JSON and error behavior must remain unchanged.

2. **Handlers depend on provider contracts**
   - Transport handlers should call explicit provider interfaces instead of constructing lower-level services or accessing repositories/stores directly.

3. **Application services remain implementation candidates**
   - Existing `auth/app.Application`, `user/app.Service` and `usersettings/app.Service` should be reused where possible instead of replacing stable behavior.

4. **Read/write model separation remains authoritative**
   - Provider methods should return existing read models or response types and accept existing write/input models where already defined.

5. **Composition root owns construction**
   - Provider instances should be built in or near `internal/app/app.go` and injected into router/handler construction during implementation.

6. **Incremental compatibility is mandatory**
   - Root compatibility aliases and existing constructor shapes may remain during 0.13.3 to avoid broad churn.

---

## Target Provider Boundaries

### 1. Auth Session Provider

Purpose: own authentication session, current-user and bootstrap orchestration currently handled mostly by `auth/app.Application`.

Candidate interface shape:

```go
type AuthSessionProvider interface {
    Login(ctx context.Context, email, password string) (LoginResponse, error)
    GetMe(ctx context.Context, claims *coreauth.Claims) (MeResponse, error)
    GetSession(ctx context.Context, claims *coreauth.Claims) (SessionResponse, error)
    GetBootstrap(ctx context.Context, claims *coreauth.Claims) (BootstrapResponse, error)
}
```

Implementation direction for 0.13.3:

- `*auth/app.Application` should satisfy this boundary using existing methods.
- Root `internal/modules/auth.Application` may remain as compatibility adapter if needed.
- Handlers should call this boundary rather than reconstructing application facades on demand.

### 2. Authenticated Account Provider

Purpose: own authenticated account-surface operations that currently leak direct calls from auth handlers into `h.Users` and `h.UserSettings`.

Candidate interface shape:

```go
type AuthenticatedAccountProvider interface {
    UpdateProfile(ctx context.Context, claims *coreauth.Claims, input UpdateProfileInput) (MeResponse, error)
    GetSettings(ctx context.Context, claims *coreauth.Claims) (UserSettingsResponse, error)
    UpdateSettings(ctx context.Context, claims *coreauth.Claims, input UpdateSettingsInput) (UserSettingsResponse, error)
}
```

Implementation direction for 0.13.3:

- Add application-level methods only if they preserve current handler behavior exactly.
- Reuse existing user and user settings application services as internal collaborators.
- Keep current validation, defaults, mapping and standardized error semantics intact.

### 3. Auth Wallet Provider

Purpose: consolidate wallet challenge, verification, listing, merge, primary and detach orchestration behind one handler-facing boundary.

Candidate interface shape:

```go
type AuthWalletProvider interface {
    ListWallets(ctx context.Context, userID string, query WalletsQuery) (WalletsResponse, error)
    CreateWalletLinkChallenge(ctx context.Context, userID, address, chain string) (WalletLinkChallengeResponse, error)
    VerifyWalletLink(ctx context.Context, userID, challengeID, address, signature string) (WalletLinkVerifyResponse, error)
    CreateWalletAccountMergeChallenge(ctx context.Context, userID, address, chain string) (WalletAccountMergeChallengeResponse, error)
    VerifyWalletAccountMerge(ctx context.Context, userID, challengeID, address, signature string) (WalletAccountMergeVerifyResponse, error)
    SetPrimaryWallet(ctx context.Context, userID, address string) (WalletPrimarySetResponse, error)
    CheckWalletDetach(ctx context.Context, userID, address string) (WalletDetachCheckResponse, error)
    ExecuteWalletDetach(ctx context.Context, userID, address string) (WalletDetachExecuteResponse, error)
}
```

Implementation direction for 0.13.3:

- Prefer `auth/app.Application` as the provider implementation because most methods already exist there.
- Move remaining handler-level service construction behind this boundary.
- Preserve all wallet response read models introduced or stabilized before Phase 0.13.

### 4. Composite Auth Provider

Purpose: provide a single injection point for auth HTTP handlers while keeping smaller interfaces available for tests and incremental implementation.

Candidate interface shape:

```go
type AuthProvider interface {
    AuthSessionProvider
    AuthenticatedAccountProvider
    AuthWalletProvider
}
```

Implementation direction for 0.13.3:

- `HTTPHandlers` should eventually depend on the composite provider or on the narrow provider needed by each handler group.
- The implementation may be phased: first make the current application facade satisfy the target methods, then reduce direct handler dependency exposure.

### 5. User Provider Boundary

Purpose: preserve a narrow cross-module contract for user resolution and profile mutation without exposing repository details.

Current reusable boundary:

- `auth/domain.UserProvider` already covers current auth needs.
- `user/app.Service` is the current implementation candidate.

Design decision:

- Do not introduce a parallel user provider unless implementation reveals a concrete need.
- Prefer reusing the existing contract and keeping it narrow.
- Profile update should be reached through authenticated account provider methods from handlers.

### 6. User Settings Provider Boundary

Purpose: preserve a narrow cross-module contract for user settings retrieval and mutation without exposing repository details.

Current reusable boundary:

- `auth/domain.UserSettingsProvider` already covers settings retrieval/update requirements.
- `usersettings/app.Service` is the current implementation candidate.

Design decision:

- Do not introduce a parallel user settings provider unless implementation reveals a concrete need.
- Auth handlers should not call user settings service directly after provider implementation.
- Settings response mapping must remain aligned with existing read models.

---

## Interface Placement Decision

Recommended placement for 0.13.3:

| Interface | Preferred Location | Reason |
| --- | --- | --- |
| `AuthSessionProvider` | `internal/modules/auth/app` or root auth adapter package | Owns auth application response types and avoids repository exposure. |
| `AuthenticatedAccountProvider` | `internal/modules/auth/app` or root auth adapter package | Owns authenticated account surface currently exposed through auth handlers. |
| `AuthWalletProvider` | `internal/modules/auth/app` or root auth adapter package | Owns wallet orchestration and existing wallet response types. |
| `AuthProvider` | Same package as narrow provider interfaces | Composite handler-facing contract. |
| user/user settings provider contracts | Existing `internal/modules/auth/domain` contracts | Already used as cross-module contracts and should remain narrow. |

The exact Go package should be selected during 0.13.3 to avoid import cycles. The design constraint is more important than the filename: handlers must depend on provider contracts, while providers may depend on application services and domain contracts.

---

## Mapping And Model Rules

Provider interfaces must not bypass Phase 0.12 mapping ownership.

Required rules:

- Provider methods that return public HTTP responses must return existing response/read model types.
- Provider implementations must continue using module mappers for read/write transformations.
- Handlers should remain responsible for transport parsing, query extraction, status codes and response writing.
- Providers should own orchestration, validation delegation and application flow sequencing.
- Repositories and stores must remain hidden from handlers after provider integration.

---

## Error Handling Rules

Provider implementation must preserve the standardized error model from Phase 0.8 and the authorization/versioning behavior from Phases 0.9 and 0.10.

Required rules:

- Do not introduce new public error shapes.
- Keep existing `normalizeApplicationError` / error adapter behavior until explicitly replaced.
- Preserve HTTP status code behavior for existing endpoints.
- Provider errors should be application/domain errors that existing adapters can translate.

---

## 0.13.3 Implementation Handoff

0.13.3 should implement this design incrementally in this order:

1. Add provider interfaces in a package that avoids import cycles.
2. Make the existing auth application facade satisfy the session and wallet provider boundaries.
3. Add authenticated account provider methods for profile/settings flows while preserving current behavior.
4. Adjust auth handlers to depend on provider boundaries instead of direct lower-level service/store access.
5. Keep composition-root wiring compatible during the transition.
6. Run `go test ./...` and update documentation with exact changed files and validation output.

---

## Non-Goals

0.13.2 and the following implementation must not introduce:

- new public endpoints
- route renames
- API version changes
- new authorization rules
- repository rewrites
- database schema changes
- CQRS/event sourcing
- behavior changes in login, bootstrap, wallet, profile or settings flows

---

## Result

0.13.2 locks the provider interface design needed to begin implementation.

The repository is ready for:

```text
0.13.3 — Provider Implementation
```
