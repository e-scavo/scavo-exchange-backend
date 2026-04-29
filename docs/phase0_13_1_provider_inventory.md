# Phase 0.13.1 — Provider Inventory & Classification

## Status

Phase 0.13.1 is **COMPLETED**.

This subphase is documentation-only. It inventories the current provider-like boundaries, direct access patterns and consolidation candidates found in the repository before provider interface design begins.

No Go code is changed by 0.13.1.

---

## Objective

Identify and classify the current access paths between HTTP handlers, WebSocket handlers, application services, domain contracts and repositories so Phase 0.13 can consolidate the Provider Layer without changing public behavior.

The inventory is intentionally conservative. It records the real repository state and does not invent provider contracts that are not present in the code.

---

## Source Scope

The inventory was performed against the repository state supplied as `scavo-exchange-backend-0.13.0.fix2.zip`.

Reviewed source areas:

- `internal/app/app.go`
- `internal/modules/auth/*.go`
- `internal/modules/auth/app/*.go`
- `internal/modules/auth/domain/*.go`
- `internal/modules/auth/mappers/*.go`
- `internal/modules/auth/readmodels/*.go`
- `internal/modules/auth/writemodels/*.go`
- `internal/modules/user/*.go`
- `internal/modules/user/app/*.go`
- `internal/modules/user/domain/*.go`
- `internal/modules/user/mappers/*.go`
- `internal/modules/usersettings/*.go`
- `internal/modules/usersettings/app/*.go`
- `internal/modules/usersettings/domain/*.go`
- `internal/modules/usersettings/mappers/*.go`
- `internal/core/httpx/*.go`
- `internal/core/ws/*.go`
- `internal/modules/system/*.go`

---

## Classification Legend

| Classification | Meaning |
| --- | --- |
| `PROVIDER_OK` | Existing boundary is already explicit enough for current behavior and can be preserved during provider consolidation. |
| `PROVIDER_CANDIDATE` | Existing service, application facade or interface behaves like a provider but is not yet formalized as the Provider Layer boundary. |
| `PROVIDER_MISSING` | A handler or composition point still accesses services, stores or application wiring directly where a provider boundary should eventually own the orchestration. |
| `PROVIDER_INVALID` | Logic is located in a layer that should not own it, or direct access creates an architectural boundary violation that should be fixed in later subphases. |
| `COMPATIBILITY_WIRING` | Transitional wiring required by the current architecture; it should not be removed without an explicit replacement plan. |
| `UNKNOWN` | Insufficient evidence in the repository to classify safely. |

---

## Executive Summary

The repository already contains several provider-like boundaries, but they are not yet consolidated into a single Provider Layer pattern.

Current state:

- `internal/modules/auth/app/Application` is the strongest provider candidate because HTTP handlers can call application-level methods instead of lower-level stores.
- `internal/modules/user/app/Service` and `internal/modules/usersettings/app/Service` are module application services and already act as stable cross-module contracts.
- `auth/domain.UserProvider` and `auth/domain.UserSettingsProvider` are explicit provider contracts used by the auth application layer.
- Wallet operations still show mixed patterns: some flows go through `auth/app/Application`, while some HTTP handlers instantiate wallet services directly.
- `internal/app/app.go` still performs composition-root wiring directly by constructing repositories, services and stores. This is acceptable as compatibility wiring but should be considered when Provider Layer ownership is formalized.
- No public route, payload, API version or business behavior should change as a result of this inventory.

---

## Module Inventory

### Composition Root — `internal/app`

| Area | Files | Classification | Notes |
| --- | --- | --- | --- |
| Runtime service construction | `internal/app/app.go` | `COMPATIBILITY_WIRING` | Constructs repositories, module services, stores, token service, router and WS dispatcher. This is acceptable composition-root wiring, but provider construction should eventually be centralized here or near here once interfaces are defined. |
| User service wiring | `internal/app/app.go` | `PROVIDER_CANDIDATE` | Wires `usermod.NewService(...)` and passes it to auth and router layers. |
| User settings service wiring | `internal/app/app.go` | `PROVIDER_CANDIDATE` | Wires `usersettingsmod.NewService(...)` and passes it to auth and router layers. |
| Wallet challenge and identity stores | `internal/app/app.go` | `PROVIDER_MISSING` | Stores are passed directly into HTTP handler/router construction. Later provider design should determine whether wallet store orchestration remains app-level or becomes provider-owned. |
| Router parameter surface | `internal/app/app.go`, `internal/core/httpx/router.go` | `COMPATIBILITY_WIRING` | Router receives many module dependencies directly. This must remain stable until provider integration is explicitly implemented. |

### Auth Module — HTTP Surface

| Area | Files | Classification | Notes |
| --- | --- | --- | --- |
| `HTTPHandlers.Application()` facade creation | `internal/modules/auth/application.go` | `PROVIDER_CANDIDATE` | Creates an auth application facade from handler dependencies. This is close to the intended provider boundary, but it is still constructed ad hoc from handler fields. |
| Login / me / session / bootstrap application calls | `internal/modules/auth/http_login.go`, `internal/modules/auth/http_bootstrap.go` | `PROVIDER_CANDIDATE` | Several handlers already call `h.Application().<method>`, preserving a clear boundary. |
| Profile update | `internal/modules/auth/http_login.go` | `PROVIDER_MISSING` | `UpdateMe` uses `h.Users.UpdateDisplayName(...)` and builds profile response in the handler. This should be evaluated for provider consolidation. |
| User settings read/update | `internal/modules/auth/http_login.go` | `PROVIDER_MISSING` | `MeSettings` and `UpdateMeSettings` call `h.UserSettings` directly and map responses in the handler. |
| Wallet challenge and verification | `internal/modules/auth/http_wallet.go` | `PROVIDER_MISSING` | Handler creates `NewWalletChallengeService`, `NewService` and `NewWalletVerificationService` directly. This is a clear provider consolidation candidate. |
| Wallet list and wallet management endpoints | `internal/modules/auth/http_wallet_list.go` and wallet-related handler files | `PROVIDER_CANDIDATE` | Several wallet flows are exposed through application-level methods, but direct store/service access patterns remain in nearby handlers. Consolidation should normalize all wallet endpoints behind one provider/application entrypoint. |
| Error adaptation and public response helpers | `internal/modules/auth/error_adapter.go`, `internal/modules/auth/http_login.go` | `COMPATIBILITY_WIRING` | Must remain compatible while provider contracts are designed. Public error JSON must not drift. |

### Auth Module — Application Layer

| Area | Files | Classification | Notes |
| --- | --- | --- | --- |
| Auth application facade | `internal/modules/auth/app/application.go` | `PROVIDER_CANDIDATE` | Centralizes login, session, bootstrap, wallet listing/linking/merge/primary/detach orchestration. This is the strongest existing provider-like object. |
| Auth service | `internal/modules/auth/app/service.go` | `PROVIDER_OK` | Provides login and session resolution behavior using token service and `UserProvider`. |
| Wallet challenge/link/merge/primary/detach services | `internal/modules/auth/app/wallet_services.go` | `PROVIDER_CANDIDATE` | Encapsulate wallet use cases but are still individually created by the application and some handlers. They should probably become internal collaborators behind a consolidated auth provider. |
| Profile/session support helpers | `internal/modules/auth/app/support_helpers.go` | `PROVIDER_CANDIDATE` | Used by application flows. Should remain app/provider internal rather than being called from transport handlers. |
| Application-level mapper usage | `internal/modules/auth/app/application.go` | `PROVIDER_OK` | Mapping is already centralized through module mappers introduced in Phase 0.12. |

### Auth Module — Domain Contracts and Stores

| Area | Files | Classification | Notes |
| --- | --- | --- | --- |
| User provider contract | `internal/modules/auth/domain/user_contract.go` | `PROVIDER_OK` | Explicit contract for auth layer dependency on user resolution. |
| User settings provider contract | `internal/modules/auth/domain/usersettings_contract.go` | `PROVIDER_OK` | Explicit contract for auth layer dependency on user settings. |
| Wallet identity store contract | `internal/modules/auth/domain/contracts.go` | `PROVIDER_OK` | Store contract is explicit, but it is persistence-oriented rather than provider-oriented. |
| Wallet challenge store contract | `internal/modules/auth/domain/contracts.go` | `PROVIDER_OK` | Store contract is explicit, but direct handler/service usage should be hidden behind provider/application boundaries later. |
| Repository adapters | `internal/modules/auth/repository/*.go`, root compatibility store files | `COMPATIBILITY_WIRING` | Existing repository/store compatibility files support runtime and tests. No provider change should remove them without explicit migration. |

### User Module

| Area | Files | Classification | Notes |
| --- | --- | --- | --- |
| User service | `internal/modules/user/app/service.go` | `PROVIDER_OK` | Owns user lookup, dev user creation, wallet user creation and display-name update behavior. |
| Root compatibility alias | `internal/modules/user/service.go` | `COMPATIBILITY_WIRING` | Preserves existing import paths and constructor shape. |
| Domain repository contract | `internal/modules/user/domain/contracts.go` | `PROVIDER_OK` | Clear domain-owned repository contract. |
| Read/write mappers | `internal/modules/user/mappers/*.go`, `readmodels`, `writemodels` | `PROVIDER_OK` | Mapping ownership remains aligned with Phase 0.12. |
| Direct access from auth handlers | `internal/modules/auth/http_login.go` | `PROVIDER_MISSING` | The user module service is used directly by auth handlers for profile update; provider design should determine the final boundary. |

### User Settings Module

| Area | Files | Classification | Notes |
| --- | --- | --- | --- |
| User settings service | `internal/modules/usersettings/app/service.go` | `PROVIDER_OK` | Owns default retrieval, preference validation, deep merge and persistence orchestration. |
| Root compatibility alias | `internal/modules/usersettings/service.go` | `COMPATIBILITY_WIRING` | Preserves existing import paths and constructor shape. |
| Domain repository contract | `internal/modules/usersettings/domain/contracts.go` | `PROVIDER_OK` | Clear domain-owned repository contract. |
| Read/write mappers | `internal/modules/usersettings/mappers/*.go`, `readmodels`, `writemodels` | `PROVIDER_OK` | Mapping ownership remains aligned with Phase 0.12. |
| Direct access from auth handlers | `internal/modules/auth/http_login.go` | `PROVIDER_MISSING` | Auth handlers call user settings service directly for settings read/update. Provider design should normalize this. |

### System Module and Core Runtime

| Area | Files | Classification | Notes |
| --- | --- | --- | --- |
| System WebSocket ping | `internal/modules/system/ws_handlers.go` | `PROVIDER_OK` | Simple system handler with no domain provider requirement at current scope. |
| Core status service | `internal/core/status/service.go` | `PROVIDER_OK` | Core operational service remains outside Phase 0.13 provider-domain consolidation. |
| HTTP router and middleware | `internal/core/httpx/*.go` | `COMPATIBILITY_WIRING` | Router composes modules and enforces versioning/auth/authorization boundaries. Provider integration must preserve public behavior. |
| WebSocket dispatcher/session | `internal/core/ws/*.go` | `COMPATIBILITY_WIRING` | Core transport boundary should remain stable. Module providers may be injected later only through explicit integration. |

---

## Current Provider-Like Boundaries

The following existing structures can be reused or formalized during later subphases:

| Existing Boundary | Current Location | Proposed Role In Later Subphases |
| --- | --- | --- |
| `auth/app.Application` | `internal/modules/auth/app/application.go` | Primary auth provider candidate / application facade. |
| `auth/app.Service` | `internal/modules/auth/app/service.go` | Internal collaborator for login and session resolution. |
| `auth/app.WalletChallengeService` | `internal/modules/auth/app/wallet_services.go` | Internal collaborator behind auth provider. |
| `auth/app.WalletLinkingService` | `internal/modules/auth/app/wallet_services.go` | Internal collaborator behind auth provider. |
| `auth/app.WalletAccountMergeService` | `internal/modules/auth/app/wallet_services.go` | Internal collaborator behind auth provider. |
| `auth/app.WalletPrimaryService` | `internal/modules/auth/app/wallet_services.go` | Internal collaborator behind auth provider. |
| `auth/app.WalletDetachService` | `internal/modules/auth/app/wallet_services.go` | Internal collaborator behind auth provider. |
| `user/app.Service` | `internal/modules/user/app/service.go` | User provider or provider implementation candidate. |
| `usersettings/app.Service` | `internal/modules/usersettings/app/service.go` | User settings provider or provider implementation candidate. |
| `auth/domain.UserProvider` | `internal/modules/auth/domain/user_contract.go` | Existing narrow provider contract. |
| `auth/domain.UserSettingsProvider` | `internal/modules/auth/domain/usersettings_contract.go` | Existing narrow provider contract. |

---

## Missing Provider Boundaries

The following gaps should feed 0.13.2 Provider Interface Design:

1. **Auth HTTP provider boundary**
   - Current handlers still own dependency fields and create app/service collaborators on demand.
   - Target direction: handlers should depend on an explicit provider/application contract rather than concrete services and stores.

2. **Wallet orchestration provider boundary**
   - Wallet challenge/verify flows are still service-constructed in HTTP handlers in at least one path.
   - Target direction: wallet flows should be consistently routed through the same application/provider boundary.

3. **Profile and settings provider boundary**
   - Profile update and settings read/update currently call user and user settings services directly from auth handlers.
   - Target direction: these authenticated account-surface flows should be owned by an auth/application provider contract or an explicitly named account provider.

4. **Composition-root provider construction**
   - `internal/app/app.go` wires stores and services directly into the router.
   - Target direction: provider construction should be explicit and stable once interfaces are designed.

---

## Invalid Or Risky Patterns

No public-contract-breaking provider violation was found during inventory.

However, the following patterns are risky and should be addressed incrementally:

| Pattern | Classification | Reason |
| --- | --- | --- |
| Handler-level construction of wallet services | `PROVIDER_MISSING` | Orchestration sits in transport layer instead of consolidated provider/application boundary. |
| Handler-level direct calls to `h.Users` and `h.UserSettings` | `PROVIDER_MISSING` | Authenticated surface orchestration is split between handlers and app facade. |
| Large handler dependency struct | `PROVIDER_CANDIDATE` / `COMPATIBILITY_WIRING` | The struct currently acts as a dependency bag; later subphases should reduce direct dependency exposure. |
| Store-level dependencies passed to router/handlers | `COMPATIBILITY_WIRING` | Safe for now, but provider construction should eventually own these lower-level dependencies. |

---

## Recommended Priority For 0.13.2

0.13.2 should design provider interfaces in this order:

1. **Auth provider interface**
   - Covers login, me/session/bootstrap, profile update, settings read/update and wallet operations.
   - Must preserve existing request/response contracts.

2. **User provider boundary**
   - Either reuse `user/app.Service` directly under an interface or define a narrow provider contract where cross-module usage requires it.

3. **User settings provider boundary**
   - Either reuse `usersettings/app.Service` directly under an interface or define a narrow provider contract for authenticated settings flows.

4. **Wallet provider sub-boundary**
   - Decide whether wallet use cases remain internal to auth provider or receive a dedicated internal provider interface.

5. **Composition wiring plan**
   - Define where providers are constructed and injected before implementation begins.

---

## Compatibility Constraints

The following must remain unchanged during subsequent provider work unless a later phase explicitly changes them:

- public HTTP routes
- public WebSocket actions
- request payload shapes
- response payload shapes
- error response contract
- API versioning behavior
- authentication and authorization middleware behavior
- database schema and migrations
- read/write model mapping output

---

## Validation Notes

`go test ./...` was attempted after this documentation-only inventory, but the local container could not download the required Go toolchain because outbound access to `proxy.golang.org` failed.

Observed failure reason:

```text
go: downloading go1.25.0 (linux/amd64)
go: download go1.25.0: golang.org/toolchain@v0.0.1-go1.25.0.linux-amd64: Get "https://proxy.golang.org/golang.org/toolchain/@v/v0.0.1-go1.25.0.linux-amd64.zip": dial tcp: lookup proxy.golang.org ... connection refused
```

No Go source files were modified by this subphase.

---

## Result

0.13.1 is completed as a documentation-only provider inventory and classification subphase.

The repository is ready for:

```text
0.13.2 — Provider Interface Design
```
