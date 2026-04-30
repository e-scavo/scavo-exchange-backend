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

## 1.1.1.1 Ownership Mapping

Ownership classification uses only current repository evidence:

- Clear: one module and one primary layer own the use case, with dependencies acting as support.
- Partially clear: the primary owner is visible, but the use case relies on another module or compatibility surface for essential behavior.
- Dispersed: multiple files, layers or compatibility surfaces currently share meaningful ownership.
- Doubtful / requires later review: the current code supports the behavior, but its product or architectural owner is not settled by the code alone.

| ID | Current owner module | Owner layer | Owner file | Main owner function or method | Relevant dependencies | Ownership class |
| --- | --- | --- | --- | --- | --- | --- |
| UC-001 | `internal/app` | runtime composition | `internal/app/app.go` | `New` | `core/auth`, `core/db`, `core/cache`, `core/httpx`, `core/ws`, `auth`, `user`, `usersettings`, `system` | Partially clear |
| UC-002 | `internal/app` | runtime lifecycle | `internal/app/app.go` | `Start` | `core/ws.Hub`, HTTP server, logger | Clear |
| UC-003 | `internal/app` | runtime lifecycle | `internal/app/app.go` | `Stop` | HTTP server, database client, cache client | Clear |
| UC-004 | `auth` | app layer | `internal/modules/auth/app/application.go` | `Application.Login` | `auth/app.Service.LoginDev`, `core/auth.TokenService`, `auth/mappers`, `user` provider | Partially clear |
| UC-005 | `auth` | app layer | `internal/modules/auth/app/application.go` | `Application.GetMe` | `core/auth.Claims`, `user` provider, wallet identity store, profile helpers | Partially clear |
| UC-006 | `auth` | app layer | `internal/modules/auth/app/application.go` | `Application.GetSession` | `auth/app.Service.ResolveSessionClaims`, `core/auth.Claims`, `user` provider | Partially clear |
| UC-007 | `auth` and `user` | app layer | `internal/modules/auth/app/application.go`; `internal/modules/user/app/service.go` | `Application.UpdateProfile`; `Service.UpdateDisplayName` | `auth/domain.ProfileUpdateInput`, `user` provider, user repository | Dispersed |
| UC-008 | `auth` and `usersettings` | app layer | `internal/modules/auth/app/application.go`; `internal/modules/usersettings/app/service.go` | `Application.GetSettings`; `Service.GetOrDefault` | `core/auth.Claims`, usersettings provider, usersettings mapper/repository | Dispersed |
| UC-009 | `auth` and `usersettings` | app layer | `internal/modules/auth/app/application.go`; `internal/modules/usersettings/app/service.go` | `Application.UpdateSettings`; `Service.UpdatePreferences` | `auth/domain.SettingsUpdateInput`, usersettings provider, usersettings repository | Dispersed |
| UC-010 | `auth` | app service | `internal/modules/auth/app/application.go`; `internal/modules/auth/app/wallet_services.go` | `Application.CreateWalletChallenge`; `WalletChallengeService.Create` | wallet challenge store, public base URL, challenge TTL | Partially clear |
| UC-011 | `auth` with `user` dependency | app layer | `internal/modules/auth/app/application.go` | `Application.VerifyWallet` | wallet challenge service, wallet identity store, signature recovery, `user.ResolveOrCreateWalletUser`, token service | Partially clear |
| UC-012 | `auth` aggregator | app layer | `internal/modules/auth/app/application.go` | `Application.GetBootstrap` | auth session service, `user` provider, `usersettings` provider, wallet identity store, read mappers | Dispersed |
| UC-013 | `auth` | app layer | `internal/modules/auth/app/application.go` | `Application.ListWallets` | wallet identity store, auth mappers, query application helpers | Partially clear |
| UC-014 | `auth` | app layer/app service | `internal/modules/auth/app/application.go`; `internal/modules/auth/app/wallet_services.go` | `Application.CreateWalletLinkChallenge`; `WalletLinkingService.CreateChallenge` | wallet challenge service, wallet identity store | Partially clear |
| UC-015 | `auth` | app layer/app service | `internal/modules/auth/app/application.go`; `internal/modules/auth/app/wallet_services.go` | `Application.VerifyWalletLink`; `WalletLinkingService.VerifyAndLink` | challenge store, wallet identity store, signature recovery | Partially clear |
| UC-016 | `auth` | app layer/app service | `internal/modules/auth/app/application.go`; `internal/modules/auth/app/wallet_services.go` | `Application.CreateWalletAccountMergeChallenge`; `WalletAccountMergeService.CreateChallenge` | wallet challenge service, wallet identity store | Doubtful / requires later review |
| UC-017 | `auth` | app layer/app service | `internal/modules/auth/app/application.go`; `internal/modules/auth/app/wallet_services.go` | `Application.VerifyWalletAccountMerge`; `WalletAccountMergeService.VerifyAndMerge` | challenge store, wallet identity store, signature recovery | Doubtful / requires later review |
| UC-018 | `auth` | app layer/app service | `internal/modules/auth/app/application.go`; `internal/modules/auth/app/wallet_services.go` | `Application.SetPrimaryWallet`; `WalletPrimaryService.SetPrimary` | wallet identity store | Partially clear |
| UC-019 | `auth` | app layer/app service | `internal/modules/auth/app/application.go`; `internal/modules/auth/app/wallet_services.go` | `Application.CheckWalletDetach`; `WalletDetachService.CheckEligibility` | wallet identity store, detach reason constants | Partially clear |
| UC-020 | `auth` | app layer/app service | `internal/modules/auth/app/application.go`; `internal/modules/auth/app/wallet_services.go` | `Application.ExecuteWalletDetach`; `WalletDetachService.Execute` | wallet identity store, eligibility check | Partially clear |
| UC-021 | `auth` | transport | `internal/modules/auth/http_login.go` | HTTP handler methods | `AuthProvider`, `core/auth.Claims`, `core/errs`, auth/user/usersettings mappers | Clear |
| UC-022 | `auth` | transport | `internal/modules/auth/http_wallet.go` | Wallet HTTP handler methods | `AuthProvider`, `core/auth.Claims`, `core/errs`, wallet write inputs | Clear |
| UC-023 | `auth` | transport | `internal/modules/auth/http_wallet_list.go` | `HTTPHandlers.Wallets` | `AuthProvider.ListWallets`, query parsing helpers | Partially clear |
| UC-024 | `auth` | transport | `internal/modules/auth/http_bootstrap.go` | `HTTPHandlers.Bootstrap` | `AuthProvider.GetBootstrap`, `core/auth.Claims` | Clear |
| UC-025 | `auth` | WebSocket transport | `internal/modules/auth/ws_handlers.go` | `WSHandlers.whoami` | `core/ws.Client.Session` | Doubtful / requires later review |
| UC-026 | `auth` | WebSocket transport | `internal/modules/auth/ws_handlers.go` | `WSHandlers.session` | root `auth.Service`, `core/ws.Client.Session`, `core/auth.Claims` | Partially clear |
| UC-027 | `auth` | app service | `internal/modules/auth/app/service.go` | `Service.LoginDev` | `core/auth.TokenService`, `user` provider | Partially clear |
| UC-028 | `auth` | app service | `internal/modules/auth/app/service.go` | `Service.LoginWallet`; `Service.LoginWalletForUser` | `core/auth.TokenService`, wallet user fallback helpers | Partially clear |
| UC-029 | `auth` | app service | `internal/modules/auth/app/service.go` | `ResolveCurrentUser*`; `ResolveSession*` | `core/auth.TokenService`, `core/auth.Claims`, `user` provider | Partially clear |
| UC-030 | `auth` | app service/repository support | `internal/modules/auth/app/wallet_services.go` | `WalletChallengeService.Create`, `CreateWithOptions`, `Get`, `MarkUsed` | `auth/domain.WalletChallengeStore`, challenge purpose helpers | Partially clear |
| UC-031 | `auth` | app service | `internal/modules/auth/app/wallet_services.go` | `WalletLinkingService.CreateChallenge`; `VerifyAndLink` | wallet challenge service, wallet identity store, signature recovery | Clear |
| UC-032 | `auth` | app service | `internal/modules/auth/app/wallet_services.go` | `WalletAccountMergeService.CreateChallenge`; `VerifyAndMerge` | wallet challenge service, wallet identity store, signature recovery | Doubtful / requires later review |
| UC-033 | `auth` | app service | `internal/modules/auth/app/wallet_services.go` | `WalletPrimaryService.SetPrimary` | wallet identity store | Clear |
| UC-034 | `auth` | app service | `internal/modules/auth/app/wallet_services.go` | `WalletDetachService.CheckEligibility`; `Execute` | wallet identity store, detach reason constants | Partially clear |
| UC-035 | `user` | app service | `internal/modules/user/app/service.go` | `Service.ResolveOrCreateDevUser` | user domain repository, fallback user construction | Partially clear |
| UC-036 | `user` | app service | `internal/modules/user/app/service.go` | `Service.ResolveOrCreateWalletUser` | user domain repository, wallet-derived identity helpers | Partially clear |
| UC-037 | `user` | app service | `internal/modules/user/app/service.go` | `Service.GetByID` | user domain repository, fallback user construction | Partially clear |
| UC-038 | `user` | app service | `internal/modules/user/app/service.go` | `Service.UpdateDisplayName` | user domain repository, validation helpers | Clear |
| UC-039 | `usersettings` | app service | `internal/modules/usersettings/app/service.go` | `Service.GetOrDefault` | usersettings domain repository, default settings, preference normalization | Clear |
| UC-040 | `usersettings` | app service | `internal/modules/usersettings/app/service.go` | `Service.UpdatePreferences` | usersettings domain repository, normalization, merge and validation helpers | Clear |
| UC-041 | `system` | WebSocket transport | `internal/modules/system/ws_handlers.go` | `ping` | `core/ws.Dispatcher`, `core/ws.Client`, time source | Clear |

### Ownership Class Summary

- Clear: UC-002, UC-003, UC-021, UC-022, UC-024, UC-031, UC-033, UC-038, UC-039, UC-040, UC-041.
- Partially clear: UC-001, UC-004, UC-005, UC-006, UC-010, UC-011, UC-013, UC-014, UC-015, UC-018, UC-019, UC-020, UC-023, UC-026, UC-027, UC-028, UC-029, UC-030, UC-034, UC-035, UC-036, UC-037.
- Dispersed: UC-007, UC-008, UC-009, UC-012.
- Doubtful / requires later review: UC-016, UC-017, UC-025, UC-032.

### Cross-Ownership Detected

- `auth` -> `user`: login, wallet bootstrap, current-user/session resolution and profile update depend on user ownership through the `auth/domain.UserProvider` contract.
- `auth` -> `usersettings`: settings read/write and bootstrap aggregation depend on the `auth/domain.UserSettingsProvider` contract.
- `auth` -> `core/auth`: login, wallet login, session and authenticated account flows depend on token and claims ownership in `internal/core/auth`.
- `auth` -> `core/ws`: WebSocket auth actions are owned by auth handlers but transported through `internal/core/ws`.
- `system` -> `core/ws`: `system.ping` is owned by the system module but exists only as a WebSocket action registered in the core dispatcher.
- `internal/app` -> all runtime modules: runtime composition owns wiring, not product behavior, while it directly constructs module providers and core services.

### Consolidation Risks For Later Phase 1.1 Work

- Root compatibility surfaces in `auth`, `user` and `usersettings` can obscure the active owner when reading call sites.
- `auth` currently exposes account/profile/settings surfaces whose business behavior is owned partly by `user` and `usersettings`.
- Bootstrap aggregates session, profile, settings and wallets; future consolidation must avoid turning the aggregator into the owner of those domains.
- Wallet challenge ownership is shared across auth bootstrap, wallet linking and account merge flows; renaming or splitting it later could break traceability if not mapped first.
- WebSocket current-user/session actions may diverge from HTTP current-user/session behavior if ownership is not normalized in later mapping.
- `internal/app` owns construction and lifecycle only; using it as a product owner would blur the Phase 1.1 boundary.

## 1.1.1.2 Use Case Duplication & Dispersion Review

This review is documentary only. It verifies duplication, overlap, dispersion and ownership ambiguity against the current repository state without correcting or refactoring any code.

Classification key:

- A) Duplication clear
- B) Partial overlap
- C) Structural dispersion
- D) Ambiguous ownership
- E) No problem / aligned

| Finding | Use case(s) | Files involved | Functions / methods involved | Type | Risk | Future recommendation | Suggested resolution phase / subphase |
| --- | --- | --- | --- | --- | --- | --- | --- |
| DDR-001 | UC-004 to UC-020, UC-024 | `internal/modules/auth/application.go`; `internal/modules/auth/app/application.go` | root `Application.Login`, `GetMe`, `GetSession`, `UpdateProfile`, `GetSettings`, `UpdateSettings`, wallet methods and `GetBootstrap`; app `Application` methods with the same use-case surface | A) Duplication clear | The root compatibility wrapper mirrors the app-layer application surface and can make the active use-case owner harder to identify. | Keep root wrapper as compatibility only until a later consolidation can make the app/provider boundary the single documented owner or explicitly document the wrapper as transport compatibility. | Phase 1.1.5 Use Case Contract Documentation, or Phase 1.2.2 Account Surface Consolidation if contract impact is discovered. |
| DDR-002 | UC-027, UC-028, UC-029 | `internal/modules/auth/service.go`; `internal/modules/auth/app/service.go` | root `Service.LoginDev`, `LoginWallet`, `LoginWalletForUser`, `ResolveCurrentUser*`, `ResolveSession*`; app `Service` methods with the same names; root `walletUser` compatibility helper | A) Duplication clear | Login/session behavior is delegated, but duplicated service names and helper names can cause future changes to land in the compatibility layer by mistake. | Preserve behavior for now; later document or remove compatibility entry points only after call sites and tests are mapped. | Phase 1.1.5 Use Case Contract Documentation; implementation cleanup only in a later roadmap-approved consolidation subphase. |
| DDR-003 | UC-013, UC-023 | `internal/modules/auth/http_wallet_list.go`; `internal/modules/auth/app/support_helpers.go`; `internal/modules/auth/app/application.go` | root `filterWalletReadModels`, `sortWalletReadModels`, `paginateWalletReadModels`, `applyWalletsQuery`, `buildWalletsResponse`; app helpers with the same function names; `Application.ListWallets` | A) Duplication clear | Wallet list filtering, sorting, pagination and response metadata exist in both root auth and auth/app helper sets. This is the clearest current behavior-drift risk because query behavior could be changed in one copy only. | Consolidate wallet-list query execution into one owner after the use-case contract is approved. Keep HTTP parsing separate if it remains a transport concern. | Phase 1.4.1 Read Surface Inventory and Phase 1.4.2-1.4.4 pagination/filtering/sorting definition. |
| DDR-004 | UC-007, UC-038 | `internal/modules/auth/http_login.go`; `internal/modules/auth/app/application.go`; `internal/modules/user/app/service.go`; `internal/modules/auth/domain/user_contract.go` | `HTTPHandlers.UpdateMe`; `Application.UpdateProfile`; `UserProvider.UpdateDisplayName`; `user/app.Service.UpdateDisplayName` | B) Partial overlap / D) Ambiguous ownership | The endpoint and app use-case are exposed as auth/profile behavior, while validation and persistence belong to the user module. This is intentional today but ambiguous as a product owner boundary. | Keep `user` as behavior owner and `auth` as current authenticated exposure until account/profile product naming is settled. | Phase 1.2.2 Account Surface Consolidation; Phase 1.5.3 Update Flow Consolidation for mutation semantics. |
| DDR-005 | UC-008, UC-009, UC-039, UC-040 | `internal/modules/auth/http_login.go`; `internal/modules/auth/app/application.go`; `internal/modules/usersettings/app/service.go`; `internal/modules/auth/domain/usersettings_contract.go` | `HTTPHandlers.MeSettings`, `UpdateMeSettings`; `Application.GetSettings`, `UpdateSettings`; `UserSettingsProvider.GetOrDefault`, `UpdatePreferences`; `usersettings/app.Service.GetOrDefault`, `UpdatePreferences` | B) Partial overlap / D) Ambiguous ownership | Settings are product-exposed through auth/account routes while the settings module owns defaults, normalization, merge and validation. Error mapping also lives in the auth handler. | Treat `usersettings` as settings behavior owner and `auth` as authenticated route owner until settings productization defines the final surface. | Phase 1.2.4 User Settings Productization; Phase 1.5.3 Update Flow Consolidation. |
| DDR-006 | UC-012 | `internal/modules/auth/app/application.go`; `internal/modules/auth/app/service.go`; `internal/modules/user/app/service.go`; `internal/modules/usersettings/app/service.go`; `internal/modules/auth/app/support_helpers.go` | `Application.GetBootstrap`; `Service.ResolveCurrentUserClaims`; `buildSessionViewWithUser`; `buildProfileViewWithUser`; `UserSettingsProvider.GetOrDefault`; wallet list helper calls | C) Structural dispersion | Bootstrap aggregates session, user/profile, settings and wallet state. The aggregator is useful, but future changes could accidentally move domain ownership into bootstrap. | Keep bootstrap as composition/read aggregation only; document downstream owners per field before expanding account capabilities. | Phase 1.1.2 Application Flow Mapping; Phase 1.1.5 Use Case Contract Documentation. |
| DDR-007 | UC-005, UC-006, UC-025, UC-026 | `internal/modules/auth/http_login.go`; `internal/modules/auth/ws_handlers.go`; `internal/modules/auth/service.go`; `internal/core/ws/session.go`; `internal/core/auth/jwt.go` | `HTTPHandlers.Me`, `Session`; `WSHandlers.whoami`, `session`; `Service.ResolveSessionClaims`; `ws.Session`; `core/auth.Claims` | B) Partial overlap / C) Structural dispersion | HTTP and WebSocket expose similar current-user/session intent with different payload construction paths. WebSocket `whoami` builds a map directly from `ws.Session`, while HTTP uses application/provider responses. | In later flow mapping, decide whether WS actions are transport equivalents, reduced diagnostics, or separate realtime session use cases. | Phase 1.1.2 Application Flow Mapping; Phase 1.6.1 Cross-Endpoint Behavior Audit. |
| DDR-008 | UC-010, UC-014, UC-016, UC-030, UC-031, UC-032 | `internal/modules/auth/app/application.go`; `internal/modules/auth/app/wallet_services.go` | `CreateWalletChallenge`; `CreateWalletLinkChallenge`; `CreateWalletAccountMergeChallenge`; `WalletChallengeService.Create`, `CreateWithOptions`; `WalletLinkingService.CreateChallenge`; `WalletAccountMergeService.CreateChallenge` | B) Partial overlap | Challenge creation is shared across auth bootstrap, wallet linking and account merge with purpose-specific options. This is aligned structurally, but the final product taxonomy is not yet settled. | Keep shared mechanics; later name the shared capability explicitly and map each product flow to its challenge purpose. | Phase 1.1.2 Application Flow Mapping; Phase 1.2.3 Wallet Management Usability. |
| DDR-009 | UC-001, UC-002, UC-003 | `internal/app/app.go`; `internal/core/httpx/router.go`; module registration files | `app.New`, `Start`, `Stop`; `httpx.NewRouter`; `registerAuthRoutes`; `auth.RegisterWS`; `system.Register` | E) No problem / aligned | Runtime composition touches all modules but does not currently duplicate product behavior. Risk exists only if future documentation treats runtime wiring as product ownership. | Keep classified as operational composition/lifecycle, not product capability. | Keep in Phase 1.1 documentation; no consolidation required unless runtime wiring starts carrying business rules. |
| DDR-010 | UC-035 to UC-040 | `internal/modules/user/service.go`; `internal/modules/user/app/service.go`; `internal/modules/usersettings/service.go`; `internal/modules/usersettings/app/service.go` | root `type Service = userapp.Service`; root `NewService`; root `type Service = usersettingsapp.Service`; root `NewService`; app service methods | E) No problem / aligned | Root user and usersettings services are aliases/factories rather than independent method implementations. They are compatibility surfaces, but no duplicated business method bodies were found in the reviewed files. | Keep documented as compatibility aliases; do not classify as active duplication unless later code adds behavior to root wrappers. | No immediate resolution required; revisit during Phase 1.1.5 if contracts need explicit wording. |

### 1.1.1.2 Summary

- Clear duplications found: root/app `auth.Application`, root/app `auth.Service`, and wallet-list query helper behavior in root `auth` plus `auth/app`.
- Partial overlaps found: profile update, settings read/write, shared wallet challenge creation, and HTTP versus WebSocket current-user/session surfaces.
- Structural dispersion found: bootstrap aggregation and HTTP/WebSocket session/current-user behavior.
- Ambiguous ownership found: account/profile exposure through `auth` with `user` behavior ownership, and settings exposure through `auth` with `usersettings` behavior ownership.
- Aligned/no-problem cases found: `internal/app` runtime composition/lifecycle and root `user` / `usersettings` service aliases.

No duplication or dispersion was corrected in 1.1.1.2. No Go source, tests or configuration changes are authorized by this review.

## Baseline Closure

1.1.1.0 establishes the initial existing-use-case baseline only.

No code, tests or configuration changes are authorized by this document. No duplicate, doubtful or dispersed behavior is corrected in this sub-subphase.
