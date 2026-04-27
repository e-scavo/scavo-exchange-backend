# Phase 0.12.1.2 — Model Classification (READ / WRITE / HYBRID)

## Purpose

This document records the lifecycle classification of the structs previously extracted during Phase 0.12.1.1.

This subphase classifies detected structs without modifying code, without moving files, and without introducing new runtime contracts.

---

## Source

* ZIP: `scavo-exchange-backend-0.12.1.1.zip`
* Baseline inventory: `docs/phase0_12_1_model_inventory.md`
* Scan scope: all `*.go` files under the repository, excluding `.git`
* Total struct declarations reviewed: 125

---

## Classification Rules

### READ

Structs used as output projections, response payloads, result payloads, view models, error payloads, or reporting shapes.

### WRITE

Structs used as input payloads, request payloads, query objects, options, parameters, or runtime configuration inputs.

### HYBRID

Structs that currently carry domain, persistence, context, transport, or error semantics across multiple lifecycle boundaries and therefore require later cross-layer analysis.

### INTERNAL

Structs that are not read/write data models: services, handlers, repositories, stores, cryptographic implementation state, test stubs, adapters, infrastructure clients, or application wiring containers.

> INTERNAL is recorded to keep the 125-struct inventory exhaustive while avoiding false READ/WRITE/HYBRID labeling for non-model infrastructure structs.

---

## Summary

* READ: 48
* WRITE: 18
* HYBRID: 11
* INTERNAL: 48
* Data model candidates (READ + WRITE + HYBRID): 77
* Total reviewed: 125

---

## Hybrid Candidates

These structs are the primary candidates for deeper cross-layer usage analysis in Phase 0.12.1.3.

| Module | File | Struct | Line | Reason |
|---|---|---:|---:|---|
| core | `internal/core/auth/jwt.go` | `Claims` | 17 | Context/state object crossing authentication or runtime boundaries. |
| core | `internal/core/authorization/subject.go` | `AuthorizationSubject` | 5 | Context/state object crossing authentication or runtime boundaries. |
| core | `internal/core/ws/session.go` | `Session` | 9 | Context/state object crossing authentication or runtime boundaries. |
| core | `internal/core/ws/protocol.go` | `Envelope` | 20 | Transport envelope can represent inbound and outbound protocol data. |
| core | `internal/core/errs/catalog_auth.go` | `AppErrorSpec` | 8 | Internal error model used to drive behavior and output mapping. |
| core | `internal/core/errs/app_error.go` | `AppError` | 15 | Internal error model used to drive behavior and output mapping. |
| user | `internal/modules/user/domain/model.go` | `User` | 5 | Domain/persistence model used across read and write flows. |
| auth | `internal/modules/auth/domain/contracts.go` | `WalletIdentity` | 28 | Domain/persistence model used across read and write flows. |
| auth | `internal/modules/auth/domain/contracts.go` | `WalletChallenge` | 48 | Domain/persistence model used across read and write flows. |
| auth | `internal/modules/auth/error_adapter.go` | `authErrorSpec` | 5 | Adapter shape used to map internal error behavior to contract output. |
| usersettings | `internal/modules/usersettings/domain/model.go` | `UserSettings` | 5 | Domain/persistence model used across read and write flows. |

---

## Full Classification Inventory

### Module: app

#### File: `internal/app/app.go`

| Struct | Line | Classification | Reason |
|---|---:|---|---|
| `App` | 22 | INTERNAL | Operational/infrastructure struct; not a read/write data model. |

### Module: auth

#### File: `internal/modules/auth/app/application.go`

| Struct | Line | Classification | Reason |
|---|---:|---|---|
| `Application` | 20 | INTERNAL | Operational/infrastructure struct; not a read/write data model. |

#### File: `internal/modules/auth/app/response_types.go`

| Struct | Line | Classification | Reason |
|---|---:|---|---|
| `LoginResponse` | 11 | READ | Output/projection/result payload. |
| `SessionView` | 18 | READ | Output/projection/result payload. |
| `ProfileWalletView` | 33 | READ | Output/projection/result payload. |
| `ProfileView` | 42 | READ | Output/projection/result payload. |
| `MeResponse` | 57 | READ | Output/projection/result payload. |
| `SessionResponse` | 62 | READ | Output/projection/result payload. |
| `BootstrapWalletsView` | 66 | READ | Output/projection/result payload. |
| `BootstrapResponse` | 71 | READ | Output/projection/result payload. |
| `WalletReadModel` | 79 | READ | Output/projection/result payload. |
| `WalletsQuery` | 92 | WRITE | Input/options/parameter object used to drive an operation. |
| `WalletsResponse` | 105 | READ | Output/projection/result payload. |
| `WalletLinkChallengeResponse` | 117 | READ | Output/projection/result payload. |
| `WalletLinkVerifyResponse` | 121 | READ | Output/projection/result payload. |
| `WalletAccountMergeChallengeResponse` | 127 | READ | Output/projection/result payload. |
| `WalletAccountMergeVerifyResponse` | 131 | READ | Output/projection/result payload. |
| `WalletDetachCheckResponse` | 139 | READ | Output/projection/result payload. |
| `WalletDetachExecuteResponse` | 147 | READ | Output/projection/result payload. |
| `WalletPrimarySetResponse` | 153 | READ | Output/projection/result payload. |

#### File: `internal/modules/auth/app/service.go`

| Struct | Line | Classification | Reason |
|---|---:|---|---|
| `Service` | 24 | INTERNAL | Operational/infrastructure struct; not a read/write data model. |
| `LoginResult` | 30 | READ | Output/projection/result payload. |

#### File: `internal/modules/auth/app/support_helpers.go`

| Struct | Line | Classification | Reason |
|---|---:|---|---|
| `authenticatedContextView` | 14 | READ | Output/projection/result payload. |

#### File: `internal/modules/auth/app/wallet_crypto.go`

| Struct | Line | Classification | Reason |
|---|---:|---|---|
| `secp256k1Curve` | 16 | INTERNAL | Operational/infrastructure struct; not a read/write data model. |

#### File: `internal/modules/auth/app/wallet_services.go`

| Struct | Line | Classification | Reason |
|---|---:|---|---|
| `WalletChallengeService` | 42 | INTERNAL | Operational/infrastructure struct; not a read/write data model. |
| `WalletLinkingService` | 48 | INTERNAL | Operational/infrastructure struct; not a read/write data model. |
| `WalletLinkResult` | 53 | READ | Output/projection/result payload. |
| `WalletAccountMergeService` | 59 | INTERNAL | Operational/infrastructure struct; not a read/write data model. |
| `WalletAccountMergeResult` | 64 | READ | Output/projection/result payload. |
| `WalletPrimaryService` | 72 | INTERNAL | Operational/infrastructure struct; not a read/write data model. |
| `WalletPrimaryResult` | 76 | READ | Output/projection/result payload. |
| `WalletDetachCheckResult` | 81 | READ | Output/projection/result payload. |
| `WalletDetachExecuteResult` | 89 | READ | Output/projection/result payload. |
| `WalletDetachService` | 95 | INTERNAL | Operational/infrastructure struct; not a read/write data model. |

#### File: `internal/modules/auth/application.go`

| Struct | Line | Classification | Reason |
|---|---:|---|---|
| `Application` | 20 | INTERNAL | Operational/infrastructure struct; not a read/write data model. |
| `walletChallengeStoreAdapter` | 233 | INTERNAL | Operational/infrastructure struct; not a read/write data model. |

#### File: `internal/modules/auth/auth_context.go`

| Struct | Line | Classification | Reason |
|---|---:|---|---|
| `authenticatedContextView` | 9 | READ | Output/projection/result payload. |

#### File: `internal/modules/auth/domain/contracts.go`

| Struct | Line | Classification | Reason |
|---|---:|---|---|
| `WalletIdentity` | 28 | HYBRID | Domain/persistence model used across read and write flows. |
| `WalletChallenge` | 48 | HYBRID | Domain/persistence model used across read and write flows. |
| `WalletChallengeOptions` | 68 | WRITE | Input/options/parameter object used to drive an operation. |

#### File: `internal/modules/auth/error_adapter.go`

| Struct | Line | Classification | Reason |
|---|---:|---|---|
| `authErrorSpec` | 5 | HYBRID | Adapter shape used to map internal error behavior to contract output. |

#### File: `internal/modules/auth/http_handlers_test.go`

| Struct | Line | Classification | Reason |
|---|---:|---|---|
| `errorPayload` | 20 | READ | Test/support response shape used to assert output contract. |
| `errorBody` | 24 | READ | Test/support response shape used to assert output contract. |
| `stubUserSettingsRepo` | 30 | INTERNAL | Test/support stub, not an application read/write model. |
| `walletDetachConflictPayload` | 2121 | READ | Test/support response shape used to assert output contract. |

#### File: `internal/modules/auth/http_login.go`

| Struct | Line | Classification | Reason |
|---|---:|---|---|
| `LoginRequest` | 16 | WRITE | HTTP input/request payload. |
| `LoginResponse` | 21 | READ | Output/projection/result payload. |
| `UpdateMeRequest` | 28 | WRITE | HTTP input/request payload. |
| `MeResponse` | 32 | READ | Output/projection/result payload. |
| `SessionResponse` | 37 | READ | Output/projection/result payload. |
| `UpdateMeSettingsRequest` | 41 | WRITE | HTTP input/request payload. |
| `MeSettingsResponse` | 45 | READ | Output/projection/result payload. |
| `HTTPHandlers` | 49 | INTERNAL | Operational/infrastructure struct; not a read/write data model. |

#### File: `internal/modules/auth/http_wallet.go`

| Struct | Line | Classification | Reason |
|---|---:|---|---|
| `WalletChallengeRequest` | 18 | WRITE | HTTP input/request payload. |
| `WalletChallengeResponse` | 23 | READ | Output/projection/result payload. |
| `WalletVerifyRequest` | 27 | WRITE | HTTP input/request payload. |
| `WalletVerifyResponse` | 33 | READ | Output/projection/result payload. |
| `WalletLinkChallengeRequest` | 113 | WRITE | HTTP input/request payload. |
| `WalletLinkVerifyRequest` | 120 | WRITE | HTTP input/request payload. |
| `WalletAccountMergeVerifyRequest` | 128 | WRITE | HTTP input/request payload. |
| `WalletDetachCheckRequest` | 136 | WRITE | HTTP input/request payload. |
| `WalletDetachExecuteRequest` | 142 | WRITE | HTTP input/request payload. |
| `WalletPrimarySetRequest` | 148 | WRITE | HTTP input/request payload. |
| `WalletAccountMergeChallengeRequest` | 156 | WRITE | HTTP input/request payload. |

#### File: `internal/modules/auth/service.go`

| Struct | Line | Classification | Reason |
|---|---:|---|---|
| `Service` | 19 | INTERNAL | Operational/infrastructure struct; not a read/write data model. |
| `LoginResult` | 25 | READ | Output/projection/result payload. |

#### File: `internal/modules/auth/service_test.go`

| Struct | Line | Classification | Reason |
|---|---:|---|---|
| `stubUserRepo` | 14 | INTERNAL | Test/support stub, not an application read/write model. |

#### File: `internal/modules/auth/wallet_challenge.go`

| Struct | Line | Classification | Reason |
|---|---:|---|---|
| `WalletChallengeService` | 59 | INTERNAL | Operational/infrastructure struct; not a read/write data model. |

#### File: `internal/modules/auth/wallet_challenge_store_memory.go`

| Struct | Line | Classification | Reason |
|---|---:|---|---|
| `InMemoryWalletChallengeStore` | 28 | INTERNAL | Operational/infrastructure struct; not a read/write data model. |

#### File: `internal/modules/auth/wallet_challenge_store_pg.go`

| Struct | Line | Classification | Reason |
|---|---:|---|---|
| `WalletChallengeStorePG` | 31 | INTERNAL | Operational/infrastructure struct; not a read/write data model. |

#### File: `internal/modules/auth/wallet_crypto.go`

| Struct | Line | Classification | Reason |
|---|---:|---|---|
| `secp256k1Curve` | 21 | INTERNAL | Operational/infrastructure struct; not a read/write data model. |

#### File: `internal/modules/auth/wallet_detach.go`

| Struct | Line | Classification | Reason |
|---|---:|---|---|
| `WalletDetachCheckResult` | 8 | READ | Output/projection/result payload. |
| `WalletDetachExecuteResult` | 16 | READ | Output/projection/result payload. |
| `WalletDetachService` | 22 | INTERNAL | Operational/infrastructure struct; not a read/write data model. |

#### File: `internal/modules/auth/wallet_identity_store_memory.go`

| Struct | Line | Classification | Reason |
|---|---:|---|---|
| `InMemoryWalletIdentityStore` | 29 | INTERNAL | Operational/infrastructure struct; not a read/write data model. |

#### File: `internal/modules/auth/wallet_identity_store_pg.go`

| Struct | Line | Classification | Reason |
|---|---:|---|---|
| `WalletIdentityStorePG` | 32 | INTERNAL | Operational/infrastructure struct; not a read/write data model. |

#### File: `internal/modules/auth/wallet_link.go`

| Struct | Line | Classification | Reason |
|---|---:|---|---|
| `WalletLinkingService` | 9 | INTERNAL | Operational/infrastructure struct; not a read/write data model. |
| `WalletLinkResult` | 14 | READ | Output/projection/result payload. |

#### File: `internal/modules/auth/wallet_merge.go`

| Struct | Line | Classification | Reason |
|---|---:|---|---|
| `WalletAccountMergeService` | 9 | INTERNAL | Operational/infrastructure struct; not a read/write data model. |
| `WalletAccountMergeResult` | 14 | READ | Output/projection/result payload. |

#### File: `internal/modules/auth/wallet_primary.go`

| Struct | Line | Classification | Reason |
|---|---:|---|---|
| `WalletPrimaryService` | 8 | INTERNAL | Operational/infrastructure struct; not a read/write data model. |
| `WalletPrimaryResult` | 12 | READ | Output/projection/result payload. |

#### File: `internal/modules/auth/wallet_verify.go`

| Struct | Line | Classification | Reason |
|---|---:|---|---|
| `WalletVerificationService` | 12 | INTERNAL | Operational/infrastructure struct; not a read/write data model. |

#### File: `internal/modules/auth/ws_handlers.go`

| Struct | Line | Classification | Reason |
|---|---:|---|---|
| `WSHandlers` | 9 | INTERNAL | Operational/infrastructure struct; not a read/write data model. |

### Module: core

#### File: `internal/core/auth/jwt.go`

| Struct | Line | Classification | Reason |
|---|---:|---|---|
| `TokenService` | 11 | INTERNAL | Operational/infrastructure struct; not a read/write data model. |
| `Claims` | 17 | HYBRID | Context/state object crossing authentication or runtime boundaries. |
| `MintOptions` | 27 | WRITE | Input/options/parameter object used to drive an operation. |

#### File: `internal/core/authorization/policy.go`

| Struct | Line | Classification | Reason |
|---|---:|---|---|
| `Decision` | 19 | READ | Response/reporting payload. |
| `PolicyEvaluator` | 26 | INTERNAL | Operational/infrastructure struct; not a read/write data model. |

#### File: `internal/core/authorization/subject.go`

| Struct | Line | Classification | Reason |
|---|---:|---|---|
| `AuthorizationSubject` | 5 | HYBRID | Context/state object crossing authentication or runtime boundaries. |

#### File: `internal/core/cache/redis.go`

| Struct | Line | Classification | Reason |
|---|---:|---|---|
| `Client` | 16 | INTERNAL | Operational/infrastructure struct; not a read/write data model. |

#### File: `internal/core/config/config.go`

| Struct | Line | Classification | Reason |
|---|---:|---|---|
| `Config` | 10 | WRITE | Runtime input/configuration model loaded into the application. |

#### File: `internal/core/db/postgres.go`

| Struct | Line | Classification | Reason |
|---|---:|---|---|
| `Client` | 17 | INTERNAL | Operational/infrastructure struct; not a read/write data model. |

#### File: `internal/core/errs/app_error.go`

| Struct | Line | Classification | Reason |
|---|---:|---|---|
| `AppError` | 15 | HYBRID | Internal error model used to drive behavior and output mapping. |

#### File: `internal/core/errs/catalog_auth.go`

| Struct | Line | Classification | Reason |
|---|---:|---|---|
| `AppErrorSpec` | 8 | HYBRID | Internal error model used to drive behavior and output mapping. |

#### File: `internal/core/errs/response_error.go`

| Struct | Line | Classification | Reason |
|---|---:|---|---|
| `ResponseError` | 3 | READ | Response/reporting payload. |
| `ErrorEnvelope` | 9 | READ | Response/reporting payload. |

#### File: `internal/core/httpx/error_test.go`

| Struct | Line | Classification | Reason |
|---|---:|---|---|
| `envelope` | 13 | READ | Test/support response shape used to assert output contract. |

#### File: `internal/core/httpx/middleware.go`

| Struct | Line | Classification | Reason |
|---|---:|---|---|
| `wrapWriter` | 60 | INTERNAL | Operational/infrastructure struct; not a read/write data model. |

#### File: `internal/core/httpx/router.go`

| Struct | Line | Classification | Reason |
|---|---:|---|---|
| `RouterParams` | 22 | WRITE | Input/options/parameter object used to wire an operation. |

#### File: `internal/core/httpx/router_versioning_test.go`

| Struct | Line | Classification | Reason |
|---|---:|---|---|
| `versionedErrorEnvelope` | 19 | READ | Test/support response shape used to assert output contract. |
| `versioningUserSettingsRepo` | 27 | INTERNAL | Test/support stub, not an application read/write model. |

#### File: `internal/core/logger/logger.go`

| Struct | Line | Classification | Reason |
|---|---:|---|---|
| `Logger` | 8 | INTERNAL | Operational/infrastructure struct; not a read/write data model. |

#### File: `internal/core/status/service.go`

| Struct | Line | Classification | Reason |
|---|---:|---|---|
| `FuncChecker` | 19 | INTERNAL | Operational/infrastructure struct; not a read/write data model. |
| `DependencyResult` | 40 | READ | Output/projection/result payload. |
| `Service` | 47 | INTERNAL | Operational/infrastructure struct; not a read/write data model. |

#### File: `internal/core/ws/client.go`

| Struct | Line | Classification | Reason |
|---|---:|---|---|
| `Client` | 11 | INTERNAL | Operational/infrastructure struct; not a read/write data model. |

#### File: `internal/core/ws/dispatcher.go`

| Struct | Line | Classification | Reason |
|---|---:|---|---|
| `Dispatcher` | 11 | INTERNAL | Operational/infrastructure struct; not a read/write data model. |

#### File: `internal/core/ws/handler.go`

| Struct | Line | Classification | Reason |
|---|---:|---|---|
| `HandlerParams` | 16 | WRITE | Input/options/parameter object used to wire an operation. |
| `Handler` | 23 | INTERNAL | Operational/infrastructure struct; not a read/write data model. |

#### File: `internal/core/ws/hub.go`

| Struct | Line | Classification | Reason |
|---|---:|---|---|
| `Hub` | 9 | INTERNAL | Operational/infrastructure struct; not a read/write data model. |

#### File: `internal/core/ws/protocol.go`

| Struct | Line | Classification | Reason |
|---|---:|---|---|
| `Envelope` | 20 | HYBRID | Transport envelope can represent inbound and outbound protocol data. |
| `ErrPayload` | 28 | READ | Response/reporting payload. |

#### File: `internal/core/ws/session.go`

| Struct | Line | Classification | Reason |
|---|---:|---|---|
| `Session` | 9 | HYBRID | Context/state object crossing authentication or runtime boundaries. |

### Module: thirdparty

#### File: `internal/thirdparty/sha3local/sha3.go`

| Struct | Line | Classification | Reason |
|---|---:|---|---|
| `state` | 23 | INTERNAL | Operational/infrastructure struct; not a read/write data model. |

#### File: `internal/thirdparty/sha3local/shake.go`

| Struct | Line | Classification | Reason |
|---|---:|---|---|
| `cshakeState` | 42 | INTERNAL | Operational/infrastructure struct; not a read/write data model. |

### Module: user

#### File: `internal/modules/user/app/service.go`

| Struct | Line | Classification | Reason |
|---|---:|---|---|
| `Service` | 20 | INTERNAL | Operational/infrastructure struct; not a read/write data model. |

#### File: `internal/modules/user/domain/model.go`

| Struct | Line | Classification | Reason |
|---|---:|---|---|
| `User` | 5 | HYBRID | Domain/persistence model used across read and write flows. |

#### File: `internal/modules/user/repository/postgres.go`

| Struct | Line | Classification | Reason |
|---|---:|---|---|
| `PostgresRepository` | 17 | INTERNAL | Operational/infrastructure struct; not a read/write data model. |

#### File: `internal/modules/user/service_test.go`

| Struct | Line | Classification | Reason |
|---|---:|---|---|
| `stubRepo` | 10 | INTERNAL | Test/support stub, not an application read/write model. |

### Module: usersettings

#### File: `internal/modules/usersettings/app/service.go`

| Struct | Line | Classification | Reason |
|---|---:|---|---|
| `Service` | 29 | INTERNAL | Operational/infrastructure struct; not a read/write data model. |

#### File: `internal/modules/usersettings/domain/model.go`

| Struct | Line | Classification | Reason |
|---|---:|---|---|
| `UserSettings` | 5 | HYBRID | Domain/persistence model used across read and write flows. |
| `View` | 12 | READ | Output/projection/result payload. |

#### File: `internal/modules/usersettings/repository/postgres.go`

| Struct | Line | Classification | Reason |
|---|---:|---|---|
| `PostgresRepository` | 13 | INTERNAL | Operational/infrastructure struct; not a read/write data model. |

#### File: `internal/modules/usersettings/service_test.go`

| Struct | Line | Classification | Reason |
|---|---:|---|---|
| `stubRepository` | 10 | INTERNAL | Test/support stub, not an application read/write model. |

---

## Notes

* This classification is intentionally structural and conservative.
* No source code was changed during this subphase.
* No read model extraction, write model isolation, or mapping layer implementation is performed here.
* Cross-layer usage analysis remains reserved for Phase 0.12.1.3.
* Problem and risk mapping remains reserved for Phase 0.12.1.4.
* Target separation definition remains reserved for Phase 0.12.1.5.

---

## Status

* Phase: 0.12.1
* Subphase: 0.12.1.2
* Status: completed
* Code changes: none
* Documentation output: `docs/phase0_12_1_model_classification.md`
