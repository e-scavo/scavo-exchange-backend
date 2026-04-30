# Stage 1 Phase 1.1 Use Case Inventory

This document records the 1.1.1.0 Use Case Inventory Baseline for the existing backend implementation.

It is documentary only. It does not authorize code changes, refactors, new features, endpoint changes, contract changes or Phase 1.2+ work. Stage 0 remains frozen and `docs/roadmap.md` remains the governing document.

## Scope

Analyzed areas:

- `internal/app`
- `internal/modules/auth`
- `internal/modules/user`
- `internal/modules/usersettings`
- `internal/modules/system`

The inventory records existing use cases only. Doubtful ownership, duplication and dispersion are documented for later Phase 1.1 consolidation work, not resolved here.

## Use Case Inventory

| ID | Module | Source file | Main function or method | Current responsibility | Layer | Apparent type | Future consolidation note |
| --- | --- | --- | --- | --- | --- | --- | --- |
| UC-001 | app | `internal/app/app.go` | `New` | Compose runtime dependencies, module services, handlers, routes, WebSocket dispatcher and server wiring. | app runtime | app composition | Treat as operational composition, not product capability. |
| UC-002 | app | `internal/app/app.go` | `Start` | Start WebSocket hub and HTTP server. | app runtime | app lifecycle | Operational surface only; keep separate from product use cases. |
| UC-003 | app | `internal/app/app.go` | `Stop` | Stop HTTP server and close database/cache resources. | app runtime | app lifecycle | Operational surface only; keep separate from product use cases. |
| UC-004 | auth | `internal/modules/auth/app/application.go` | `Login` | Authenticate the current development login flow and return login read model data. | app layer | app/use case | HTTP transport and service wrapper should stay aligned with one product name. |
| UC-005 | auth | `internal/modules/auth/app/application.go` | `GetMe` | Resolve the authenticated user and build the current profile view. | app layer | app/readmodel | Overlaps conceptually with WebSocket `auth.whoami`; keep both inventoried. |
| UC-006 | auth | `internal/modules/auth/app/application.go` | `GetSession` | Resolve authenticated session claims and expose session read model data. | app layer | app/readmodel | Overlaps conceptually with WebSocket `auth.session`; keep both inventoried. |
| UC-007 | auth | `internal/modules/auth/app/application.go` | `UpdateProfile` | Update the authenticated user's display name through the user module. | app layer | app/writemodel | Cross-module ownership with `user.UpdateDisplayName` should be clarified later. |
| UC-008 | auth | `internal/modules/auth/app/application.go` | `GetSettings` | Read current user settings through the usersettings module. | app layer | app/readmodel | Settings ownership is split between auth surface and usersettings module. |
| UC-009 | auth | `internal/modules/auth/app/application.go` | `UpdateSettings` | Update current user preferences through the usersettings module. | app layer | app/writemodel | Mutation behavior belongs to later consolidation, not this baseline. |
| UC-010 | auth | `internal/modules/auth/app/application.go` | `CreateWalletChallenge` | Create a wallet verification challenge for unauthenticated wallet login/bootstrap. | app layer | app/use case | Challenge creation is reused by multiple wallet flows. |
| UC-011 | auth | `internal/modules/auth/app/application.go` | `VerifyWallet` | Verify wallet signature/challenge, resolve or create wallet user identity and return wallet login data. | app layer | app/use case | Existing behavior only; account/identity product expansion remains Phase 1.2+. |
| UC-012 | auth | `internal/modules/auth/app/application.go` | `GetBootstrap` | Aggregate session, user/profile, settings and wallet data for authenticated bootstrap. | app layer | app/readmodel | Aggregates multiple module concerns and needs later flow mapping. |
| UC-013 | auth | `internal/modules/auth/app/application.go` | `ListWallets` | List the authenticated user's wallets with current query filtering, sorting and pagination behavior. | app layer | app/readmodel | Query parsing is in HTTP while result shaping is in app flow. |
| UC-014 | auth | `internal/modules/auth/app/application.go` | `CreateWalletLinkChallenge` | Create a challenge to link a wallet to the authenticated account. | app layer | app/use case | Uses shared wallet challenge mechanics. |
| UC-015 | auth | `internal/modules/auth/app/application.go` | `VerifyWalletLink` | Verify wallet link challenge and attach the wallet to the authenticated account. | app layer | app/writemodel | Existing wallet management capability only; no new behavior added. |
| UC-016 | auth | `internal/modules/auth/app/application.go` | `CreateWalletAccountMergeChallenge` | Create a challenge for existing wallet/account merge flow. | app layer | app/use case | Product ownership may belong to later account identity work. |
| UC-017 | auth | `internal/modules/auth/app/application.go` | `VerifyWalletAccountMerge` | Verify merge challenge and connect wallet identity to the current account flow. | app layer | app/writemodel | Existing behavior only; authorization and mutation rules remain deferred. |
| UC-018 | auth | `internal/modules/auth/app/application.go` | `SetPrimaryWallet` | Mark an authenticated user's wallet as primary. | app layer | app/writemodel | Wallet usability refinement remains later Stage 1 work. |
| UC-019 | auth | `internal/modules/auth/app/application.go` | `CheckWalletDetach` | Check whether a wallet can be detached from the authenticated account. | app layer | app/readmodel | Eligibility rule ownership should be mapped later. |
| UC-020 | auth | `internal/modules/auth/app/application.go` | `ExecuteWalletDetach` | Detach a wallet from the authenticated account when eligible. | app layer | app/writemodel | Existing mutation flow only; no consolidation applied in this baseline. |
| UC-021 | auth | `internal/modules/auth/http_login.go` | `Login`, `Me`, `UpdateMe`, `MeSettings`, `UpdateMeSettings`, `Session` | Expose auth/account/settings use cases through HTTP handlers. | transport | adapter | Handler names and app use-case names should remain traceable. |
| UC-022 | auth | `internal/modules/auth/http_wallet.go` | Wallet handler methods | Expose wallet challenge, verify, link, merge, primary and detach flows through HTTP handlers. | transport | adapter | Transport surface is frozen Stage 0 contract; no change in this sub-subphase. |
| UC-023 | auth | `internal/modules/auth/http_wallet_list.go` | `Wallets` | Expose wallet listing and parse query options for wallet list view. | transport | adapter/readmodel support | Query handling is partly transport-local and should be reviewed later. |
| UC-024 | auth | `internal/modules/auth/http_bootstrap.go` | `Bootstrap` | Expose authenticated bootstrap aggregation through HTTP. | transport | adapter | Keep mapped to `Application.GetBootstrap`. |
| UC-025 | auth | `internal/modules/auth/ws_handlers.go` | `whoami` | Expose authenticated identity/profile snapshot through WebSocket RPC. | transport | adapter/readmodel | Overlaps with HTTP current-user surface. |
| UC-026 | auth | `internal/modules/auth/ws_handlers.go` | `session` | Expose authenticated session snapshot through WebSocket RPC. | transport | adapter/readmodel | Overlaps with HTTP session surface. |
| UC-027 | auth | `internal/modules/auth/app/service.go` | `LoginDev` | Validate development login credentials and mint a token. | app service | app/use case support | Root compatibility service mirrors this behavior. |
| UC-028 | auth | `internal/modules/auth/app/service.go` | `LoginWallet`, `LoginWalletForUser` | Mint wallet-authenticated tokens. | app service | app/use case support | Supports wallet login and merge/link paths. |
| UC-029 | auth | `internal/modules/auth/app/service.go` | `ResolveCurrentUser`, `ResolveCurrentUserClaims`, `ResolveSession`, `ResolveSessionClaims` | Resolve authenticated user/session state from token claims. | app service | app/readmodel support | Used by HTTP, WebSocket and bootstrap paths. |
| UC-030 | auth | `internal/modules/auth/app/wallet_services.go` | `WalletChallengeService.Create`, `CreateWithOptions`, `Get`, `MarkUsed` | Manage wallet challenges used by login, linking and merge flows. | app service | app/repository-backed support | Shared challenge use should be named consistently later. |
| UC-031 | auth | `internal/modules/auth/app/wallet_services.go` | `WalletLinkingService.CreateChallenge`, `VerifyAndLink` | Create and verify authenticated wallet linking. | app service | app/writemodel support | Mirrors app-layer link methods. |
| UC-032 | auth | `internal/modules/auth/app/wallet_services.go` | `WalletAccountMergeService.CreateChallenge`, `VerifyAndMerge` | Create and verify wallet account merge. | app service | app/writemodel support | Product placement is doubtful and deferred. |
| UC-033 | auth | `internal/modules/auth/app/wallet_services.go` | `WalletPrimaryService.SetPrimary` | Set a user's primary wallet. | app service | app/writemodel support | Mirrors app-layer primary wallet method. |
| UC-034 | auth | `internal/modules/auth/app/wallet_services.go` | `WalletDetachService.CheckEligibility`, `Execute` | Check and execute wallet detachment. | app service | app/readmodel and writemodel support | Eligibility and mutation are separated but remain in one service. |
| UC-035 | user | `internal/modules/user/app/service.go` | `ResolveOrCreateDevUser` | Resolve or create the development user account. | app service | app/writemodel support | Existing auth bootstrap dependency. |
| UC-036 | user | `internal/modules/user/app/service.go` | `ResolveOrCreateWalletUser` | Resolve or create a user account associated with wallet identity. | app service | app/writemodel support | Account product semantics deferred to Phase 1.2. |
| UC-037 | user | `internal/modules/user/app/service.go` | `GetByID` | Read user by ID. | app service | app/readmodel support | Used by auth/session/profile paths. |
| UC-038 | user | `internal/modules/user/app/service.go` | `UpdateDisplayName` | Validate and persist user display name changes. | app service | app/writemodel | Exposed through auth account surface today. |
| UC-039 | usersettings | `internal/modules/usersettings/app/service.go` | `GetOrDefault` | Read persisted settings or return normalized default settings. | app service | app/readmodel | Exposed through auth settings surface today. |
| UC-040 | usersettings | `internal/modules/usersettings/app/service.go` | `UpdatePreferences` | Validate, normalize, merge and persist preference updates. | app service | app/writemodel | Mutation consolidation deferred to later Phase 1.1/1.5 work. |
| UC-041 | system | `internal/modules/system/ws_handlers.go` | `ping` | Expose a minimal WebSocket system ping response. | transport | adapter/readmodel | System behavior changes remain outside this baseline. |

## Duplications Or Dispersion Detected

- `internal/modules/auth/application.go` keeps a root compatibility `Application` wrapper around `internal/modules/auth/app.Application` with a similar method surface.
- `internal/modules/auth/service.go` keeps a root compatibility service around `internal/modules/auth/app.Service` with overlapping login/session methods.
- Wallet challenge, identity and service concepts appear both in root `internal/modules/auth` files and under `internal/modules/auth/app` or `internal/modules/auth/repository`.
- User and usersettings services expose root compatibility aliases while the active implementation lives under each module `app` package.
- Settings use cases are implemented in `usersettings` but exposed through the `auth` HTTP/account surface.
- Wallet list query behavior is split between HTTP query parsing and application-layer list processing.
- HTTP and WebSocket current-user/session surfaces overlap in product intent.

## Doubtful Cases

- `internal/app.New`, `Start` and `Stop` are operational runtime use cases, not direct product capabilities.
- WebSocket `auth.whoami` and `auth.session` may be alternate transports for the same product intent as HTTP current-user and session reads.
- Wallet account merge is already implemented, but its product ownership may belong to later account and identity work.
- Shared wallet challenge behavior supports login, link and merge flows; the final use-case taxonomy should decide whether it remains shared or flow-specific.
- Repository-backed wallet stores and mappers support use cases but should not be counted as product use cases unless a later subphase explicitly classifies persistence operations.
- Settings read/write ownership is split between `auth` exposure and `usersettings` implementation and should be mapped before consolidation.

## Baseline Closure

1.1.1.0 establishes the initial existing-use-case baseline only.

No code, tests or configuration changes are authorized by this document. No duplicate, doubtful or dispersed behavior is corrected in this sub-subphase.
