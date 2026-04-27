# Phase 0.12.1.3 — Cross-Layer Usage Analysis

## Purpose

This document records cross-layer usage for the hybrid candidates identified during Phase 0.12.1.2.

This subphase does not modify code, does not create new models, and does not introduce runtime behavior changes. Its only purpose is to identify where hybrid structures cross module, transport, domain, application, persistence, and core boundaries.

---

## Source

* ZIP: `scavo-exchange-backend-0.12.1.2.zip`
* Baseline inventory: `docs/phase0_12_1_model_inventory.md`
* Baseline classification: `docs/phase0_12_1_model_classification.md`
* Scan scope: all `*.go` files under the repository, excluding `.git` and generated ZIP payloads
* Hybrid candidates reviewed: 11

---

## Scope

### Included

* Review all previously classified HYBRID candidates.
* Identify cross-layer references by package/file.
* Record layer-boundary crossings.
* Identify follow-up candidates for risk mapping in Phase 0.12.1.4.

### Excluded

* No code changes.
* No model extraction.
* No read/write model creation.
* No mapping implementation.
* No contract changes.
* No endpoint changes.

---

## Cross-Layer Usage Summary

| Struct | Declared In | References | Files | Cross-Layer Usage |
|---|---|---:|---:|---|
| `Claims` | `internal/core/auth/jwt.go` | 42 | 20 | Core auth, HTTP auth helpers, WS session, auth app/service/profile flows, tests. |
| `AuthorizationSubject` | `internal/core/authorization/subject.go` | 26 | 10 | Core authorization, context hydration, HTTP authorization middleware, tests. |
| `Session` | `internal/core/ws/session.go` | 57 | 11 | WS client/session state, WS dispatcher/handler, auth WS handlers, HTTP auth session naming overlap. |
| `Envelope` | `internal/core/ws/protocol.go` | 22 | 6 | WS protocol, dispatcher, helpers, system WS handlers, auth WS handlers. |
| `AppErrorSpec` | `internal/core/errs/catalog_auth.go` | 5 | 1 | Core error catalog specification only. |
| `AppError` | `internal/core/errs/app_error.go` | 60 | 6 | Core error value, factories, HTTP error writer, auth HTTP handling, tests. |
| `User` | `internal/modules/user/domain/model.go` + alias in `internal/modules/user/model.go` | 133 | 23 | User domain, user repository, user app service, auth services, auth responses, wallet flows, tests. |
| `WalletIdentity` | `internal/modules/auth/domain/contracts.go` + alias in `internal/modules/auth/wallet_identity_store.go` | 98 | 17 | Auth domain, wallet stores, app services, HTTP wallet list/profile/link/verify flows, tests. |
| `WalletChallenge` | `internal/modules/auth/domain/contracts.go` + alias in `internal/modules/auth/wallet_challenge.go` | 62 | 16 | Auth domain, challenge stores, app response payloads, HTTP wallet flows, router wiring/tests. |
| `authErrorSpec` | `internal/modules/auth/error_adapter.go` | 3 | 1 | Auth-local error adapter shape. |
| `UserSettings` | `internal/modules/usersettings/domain/model.go` + alias in `internal/modules/usersettings/model.go` | 96 | 13 | User settings domain/app/repository, auth HTTP/app flows, router/versioning tests. |

---

## Layer Categories Used

| Category | Meaning |
|---|---|
| Core | Shared runtime packages under `internal/core/*`. |
| Transport | HTTP and WS handlers, routers, request/response helpers, dispatchers. |
| Application | Application services and orchestration packages. |
| Domain | Domain models and domain contracts. |
| Persistence | Repository/store implementations. |
| Test | Test-only references validating runtime behavior. |

---

## Detailed Hybrid Usage

### 1. `Claims`

**Declaration:** `internal/core/auth/jwt.go:17`

**Observed files:**

* `internal/core/auth/jwt.go`
* `internal/core/auth/context.go`
* `internal/core/authorization/resolve.go`
* `internal/core/httpx/auth.go`
* `internal/core/httpx/authorization_test.go`
* `internal/core/ws/handler.go`
* `internal/core/ws/session.go`
* `internal/modules/auth/app/application.go`
* `internal/modules/auth/app/service.go`
* `internal/modules/auth/app/support_helpers.go`
* `internal/modules/auth/application.go`
* `internal/modules/auth/auth_context.go`
* `internal/modules/auth/http_login.go`
* `internal/modules/auth/http_wallet.go`
* `internal/modules/auth/profile.go`
* `internal/modules/auth/service.go`
* `internal/modules/auth/ws_handlers.go`
* related auth and authorization tests

**Cross-layer usage:**

| Layer | Usage |
|---|---|
| Core | Token creation/parsing and context propagation. |
| Authorization | Converted into `AuthorizationSubject`. |
| HTTP Transport | Retrieved from request context to authorize authenticated routes. |
| WS Transport | Stored inside WS `Session`. |
| Auth Application | Passed into current-user, session, profile, and bootstrap flows. |
| Tests | Used directly to validate authorization and auth handlers. |

**Boundary notes:**

`Claims` is currently both an authentication token payload and a runtime identity carrier. It crosses from token parsing into HTTP context, WS session state, authorization subject resolution, and auth application flows.

**Follow-up needed in 0.12.1.4:**

Determine whether `Claims` should remain a core authentication runtime object or whether authenticated application flows should depend on a narrower read/context model.

---

### 2. `AuthorizationSubject`

**Declaration:** `internal/core/authorization/subject.go:5`

**Observed files:**

* `internal/core/authorization/subject.go`
* `internal/core/authorization/context.go`
* `internal/core/authorization/policy.go`
* `internal/core/authorization/resolve.go`
* `internal/core/httpx/authorization.go`
* `internal/core/httpx/authorization_enforcement_test.go`
* `internal/core/httpx/authorization_test.go`
* authorization policy/context/subject tests

**Cross-layer usage:**

| Layer | Usage |
|---|---|
| Core Authorization | Subject normalization, policy evaluation, permission checks. |
| Core Context | Stored/retrieved from request context. |
| HTTP Transport | Hydrated into request context and consumed by authorization middleware. |
| Tests | Used to assert policy and enforcement behavior. |

**Boundary notes:**

`AuthorizationSubject` is primarily core authorization state, but it crosses into HTTP middleware and request context. The usage is controlled and remains within authorization/runtime boundaries.

**Follow-up needed in 0.12.1.4:**

Assess if this is an acceptable core runtime model rather than a read/write domain model requiring extraction.

---

### 3. `Session`

**Declaration:** `internal/core/ws/session.go:9`

**Observed files:**

* `internal/core/ws/session.go`
* `internal/core/ws/client.go`
* `internal/core/ws/dispatcher.go`
* `internal/core/ws/handler.go`
* `internal/modules/auth/ws_handlers.go`
* `internal/modules/auth/app/response_types.go`
* `internal/modules/auth/app/application.go`
* `internal/modules/auth/application.go`
* `internal/modules/auth/http_login.go`
* `internal/core/httpx/router.go`
* auth HTTP handler tests

**Cross-layer usage:**

| Layer | Usage |
|---|---|
| WS Core | Client runtime state and authenticated WS session carrier. |
| WS Transport | Dispatcher validates session state before protected operations. |
| Auth WS | Auth handlers read session to resolve user/session responses. |
| Auth HTTP/App | Name overlaps with HTTP session responses and session views, though those are distinct structs. |
| Tests | Auth session tests and handler tests reference session concepts heavily. |

**Boundary notes:**

`Session` is a WS runtime state object that carries user identity and token claims. Its cross-layer usage is mostly runtime-oriented, but it sits near response/session concepts exposed by auth HTTP and app packages.

**Follow-up needed in 0.12.1.4:**

Separate naming and lifecycle risk: WS runtime session should not be confused with auth read models such as `SessionView` or `SessionResponse`.

---

### 4. `Envelope`

**Declaration:** `internal/core/ws/protocol.go:20`

**Observed files:**

* `internal/core/ws/protocol.go`
* `internal/core/ws/dispatcher.go`
* `internal/core/ws/helpers.go`
* `internal/core/ws/handler.go`
* `internal/modules/auth/ws_handlers.go`
* `internal/modules/system/ws_handlers.go`

**Cross-layer usage:**

| Layer | Usage |
|---|---|
| WS Core Protocol | Defines message envelope structure. |
| WS Dispatcher | Routes incoming envelope messages. |
| WS Helpers | Builds responses and error payloads. |
| Auth WS | Consumes envelopes for auth commands. |
| System WS | Consumes envelopes for system commands. |

**Boundary notes:**

`Envelope` is a transport protocol shape used for both inbound and outbound WS messages. This is a clear lifecycle hybrid because a single struct represents request envelope and response envelope semantics.

**Follow-up needed in 0.12.1.4:**

Mark as a high-priority candidate for eventual split into inbound/outbound WS protocol shapes or documented transport-level exception.

---

### 5. `AppErrorSpec`

**Declaration:** `internal/core/errs/catalog_auth.go:8`

**Observed files:**

* `internal/core/errs/catalog_auth.go`

**Cross-layer usage:**

| Layer | Usage |
|---|---|
| Core Errors | Error catalog specification used to define code/message/status metadata. |

**Boundary notes:**

`AppErrorSpec` is currently contained within the core error catalog. It does not cross module boundaries directly, but it shapes runtime error creation and eventual HTTP output through `AppError`.

**Follow-up needed in 0.12.1.4:**

Likely low risk. Confirm whether it should remain INTERNAL rather than HYBRID in final audit closure.

---

### 6. `AppError`

**Declaration:** `internal/core/errs/app_error.go:15`

**Observed files:**

* `internal/core/errs/app_error.go`
* `internal/core/errs/catalog_auth.go`
* `internal/core/errs/factories.go`
* `internal/core/httpx/error.go`
* `internal/modules/auth/http_login.go`
* `internal/core/errs/app_error_test.go`

**Cross-layer usage:**

| Layer | Usage |
|---|---|
| Core Errors | Runtime application error model. |
| Core Factories | Built by catalog/factory helpers. |
| HTTP Transport | Converted into standardized HTTP error responses. |
| Auth HTTP | Auth handlers use standardized error paths. |
| Tests | Validates behavior and error metadata. |

**Boundary notes:**

`AppError` is both a runtime domain-like error object and the source for standardized HTTP response generation. It bridges behavior and output mapping, but this is an intentional result of Phase 0.8.

**Follow-up needed in 0.12.1.4:**

Determine whether it remains an accepted standardized error boundary model or requires a stricter internal/error response split.

---

### 7. `User`

**Declarations:**

* `internal/modules/user/domain/model.go:5`
* alias: `internal/modules/user/model.go:5`

**Observed files:**

* `internal/modules/user/domain/model.go`
* `internal/modules/user/domain/contracts.go`
* `internal/modules/user/model.go`
* `internal/modules/user/repository/postgres.go`
* `internal/modules/user/app/service.go`
* `internal/modules/user/service_test.go`
* `internal/modules/auth/domain/user_contract.go`
* `internal/modules/auth/app/service.go`
* `internal/modules/auth/app/application.go`
* `internal/modules/auth/app/support_helpers.go`
* `internal/modules/auth/app/response_types.go`
* `internal/modules/auth/app/wallet_services.go`
* `internal/modules/auth/application.go`
* `internal/modules/auth/service.go`
* `internal/modules/auth/profile.go`
* `internal/modules/auth/http_login.go`
* `internal/modules/auth/http_wallet.go`
* `internal/modules/auth/wallet_challenge.go`
* `internal/modules/auth/wallet_verify.go`
* multiple auth tests

**Cross-layer usage:**

| Layer | Usage |
|---|---|
| User Domain | Canonical user model and contracts. |
| Persistence | Repository reads database rows into `User`. |
| User Application | Service returns domain user values. |
| Auth Domain Contract | Auth module depends on a user provider contract returning user values. |
| Auth Application | Used to build session/profile/bootstrap read models. |
| Auth HTTP | Used indirectly through auth service/application calls. |
| Tests | Used widely in auth/user service and handler tests. |

**Boundary notes:**

`User` is the strongest cross-module hybrid candidate. It is a domain/persistence model that is exposed across module boundaries and used to construct output read models. The alias in `internal/modules/user/model.go` preserves compatibility but also widens access to the domain model.

**Follow-up needed in 0.12.1.4:**

High-priority risk mapping candidate. Determine target split between user domain entity, auth-facing provider DTO/read contract, and response read models.

---

### 8. `WalletIdentity`

**Declarations:**

* `internal/modules/auth/domain/contracts.go:28`
* alias: `internal/modules/auth/wallet_identity_store.go:5`

**Observed files:**

* `internal/modules/auth/domain/contracts.go`
* `internal/modules/auth/wallet_identity_store.go`
* `internal/modules/auth/wallet_identity_store_pg.go`
* `internal/modules/auth/wallet_identity_store_memory.go`
* `internal/modules/auth/wallet_link.go`
* `internal/modules/auth/wallet_verify.go`
* `internal/modules/auth/wallet_detach.go`
* `internal/modules/auth/wallet_merge.go`
* `internal/modules/auth/wallet_primary.go`
* `internal/modules/auth/http_wallet_list.go`
* `internal/modules/auth/profile.go`
* `internal/modules/auth/application.go`
* `internal/modules/auth/app/application.go`
* `internal/modules/auth/app/support_helpers.go`
* `internal/modules/auth/app/wallet_services.go`
* `internal/modules/auth/app/response_types.go`
* auth HTTP handler tests

**Cross-layer usage:**

| Layer | Usage |
|---|---|
| Auth Domain | Wallet identity model and store contracts. |
| Persistence | Postgres and memory stores persist and return wallet identity values. |
| Auth Application | Used in wallet link/verify/detach/merge/primary flows. |
| Auth Read Models | Used to build profile/bootstrap/wallet response payloads. |
| Auth HTTP | Exposed through wallet endpoints and handler flows. |
| Tests | Used widely across wallet and HTTP handler tests. |

**Boundary notes:**

`WalletIdentity` crosses domain, persistence, application, and read-output construction. It is currently a canonical auth wallet model and response source.

**Follow-up needed in 0.12.1.4:**

High-priority risk mapping candidate. Determine whether wallet read models should be built only from dedicated mapping functions and whether store models should remain domain-only.

---

### 9. `WalletChallenge`

**Declarations:**

* `internal/modules/auth/domain/contracts.go:48`
* alias: `internal/modules/auth/wallet_challenge.go:50`

**Observed files:**

* `internal/modules/auth/domain/contracts.go`
* `internal/modules/auth/wallet_challenge.go`
* `internal/modules/auth/wallet_challenge_store_pg.go`
* `internal/modules/auth/wallet_challenge_store_memory.go`
* `internal/modules/auth/wallet_link.go`
* `internal/modules/auth/wallet_verify.go`
* `internal/modules/auth/wallet_merge.go`
* `internal/modules/auth/http_wallet.go`
* `internal/modules/auth/application.go`
* `internal/modules/auth/app/wallet_services.go`
* `internal/modules/auth/app/response_types.go`
* `internal/core/httpx/router.go`
* wallet challenge/link/verify/http tests

**Cross-layer usage:**

| Layer | Usage |
|---|---|
| Auth Domain | Wallet challenge model and challenge store contracts. |
| Persistence | Postgres and memory stores persist challenge state. |
| Auth Application | Used in wallet link, verify, and merge flows. |
| Auth Read Models | Embedded in challenge response payloads. |
| Auth HTTP | Returned by wallet challenge endpoints. |
| Core Router | Router wires wallet challenge routes that trigger these flows. |
| Tests | Used in challenge/link/verify and HTTP tests. |

**Boundary notes:**

`WalletChallenge` carries operational challenge state and is also embedded in read responses. This makes it a direct read/write separation target because challenge internals may not always match public response needs.

**Follow-up needed in 0.12.1.4:**

High-priority risk mapping candidate. Determine target read model for challenge responses and isolate persistence/runtime state.

---

### 10. `authErrorSpec`

**Declaration:** `internal/modules/auth/error_adapter.go:5`

**Observed files:**

* `internal/modules/auth/error_adapter.go`

**Cross-layer usage:**

| Layer | Usage |
|---|---|
| Auth Adapter | Local auth error normalization helper. |

**Boundary notes:**

`authErrorSpec` appears contained to one adapter file. It does not currently cross architectural layers as a system model.

**Follow-up needed in 0.12.1.4:**

Likely low risk. Confirm whether this should be downgraded from HYBRID to INTERNAL/adapter-local in audit closure.

---

### 11. `UserSettings`

**Declarations:**

* `internal/modules/usersettings/domain/model.go:5`
* alias: `internal/modules/usersettings/model.go:5`

**Observed files:**

* `internal/modules/usersettings/domain/model.go`
* `internal/modules/usersettings/domain/contracts.go`
* `internal/modules/usersettings/model.go`
* `internal/modules/usersettings/app/service.go`
* `internal/modules/usersettings/repository/postgres.go`
* `internal/modules/usersettings/service_test.go`
* `internal/modules/auth/domain/usersettings_contract.go`
* `internal/modules/auth/http_login.go`
* `internal/modules/auth/application.go`
* `internal/modules/auth/app/application.go`
* `internal/core/httpx/router.go`
* `internal/core/httpx/router_versioning_test.go`
* auth HTTP handler tests

**Cross-layer usage:**

| Layer | Usage |
|---|---|
| UserSettings Domain | Canonical settings model and contracts. |
| Persistence | Repository returns settings values. |
| UserSettings Application | Service get/update methods return domain settings. |
| Auth Domain Contract | Auth module depends on settings provider contract. |
| Auth HTTP/App | Settings are read and updated through auth endpoints. |
| Router | UserSettings service injected into route wiring. |
| Tests | Used across user settings service, auth HTTP, and versioning tests. |

**Boundary notes:**

`UserSettings` crosses domain, persistence, application, auth-facing contracts, and HTTP flows. The alias in `internal/modules/usersettings/model.go` keeps public module compatibility but also exposes the domain model directly.

**Follow-up needed in 0.12.1.4:**

High-priority risk mapping candidate. Determine target split between settings domain model, update/write patch command, and response/read view.

---

## Boundary Findings

### Strongest Cross-Layer Hybrids

These models cross the most important read/write and module boundaries:

| Priority | Struct | Reason |
|---:|---|---|
| 1 | `User` | Domain/persistence model used directly by user and auth modules and used to build read responses. |
| 2 | `UserSettings` | Domain/persistence model used by usersettings and auth flows, including update and read paths. |
| 3 | `WalletIdentity` | Domain/persistence wallet identity used across store, application, response, and HTTP flows. |
| 4 | `WalletChallenge` | Runtime/persistence challenge state embedded in public challenge responses. |
| 5 | `Envelope` | Same WS protocol struct represents inbound and outbound message lifecycle. |

### Controlled Runtime Hybrids

These models cross runtime boundaries but may be intentional core runtime objects rather than read/write domain models:

| Struct | Reason |
|---|---|
| `Claims` | Auth token payload used as request/session identity context. |
| `AuthorizationSubject` | Authorization runtime subject used by policy and middleware. |
| `Session` | WS runtime session state. |
| `AppError` | Standardized runtime error mapped to HTTP output. |

### Low-Risk / Local Hybrids

These appear local enough that later audit closure may reclassify them as INTERNAL or adapter-local:

| Struct | Reason |
|---|---|
| `AppErrorSpec` | Contained to core error catalog metadata. |
| `authErrorSpec` | Contained to auth error adapter. |

---

## Cross-Layer Risk Signals for Phase 0.12.1.4

The following risk signals were detected:

* Public module aliases expose domain models outside their domain packages.
* Domain/persistence models are used to construct read responses.
* Runtime/auth context objects cross HTTP, WS, authorization, and application boundaries.
* Some transport structs represent both inbound and outbound lifecycle semantics.
* Error models bridge internal behavior and response output mapping.
* Tests depend directly on current hybrid shapes, which will matter during later extraction phases.

---

## Non-Goals Confirmed

This subphase intentionally did not:

* change Go code;
* create read models;
* create write models;
* move packages;
* change HTTP or WS contracts;
* change tests;
* run a refactor;
* decide final extraction targets.

---

## Status

Phase: 0.12 — Read / Write Model Separation  
Subphase: 0.12.1 — Model Classification Audit  
Step: 0.12.1.3 — Cross-Layer Usage Analysis  
Status: COMPLETED  
Code Impact: NONE  
Next Step: 0.12.1.4 — Problem Detection & Risk Mapping
