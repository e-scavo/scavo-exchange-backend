# Phase 0.12.1.1 — Model Inventory Extraction

## Purpose

This document records the complete struct inventory extracted from the real ZIP used as the source of truth for Phase 0.12.1.1.

This subphase only inventories detected Go structs. It does not classify models as READ, WRITE, or HYBRID, and it does not define refactoring actions.

---

## Source

* ZIP: `scavo-exchange-backend-0.12.1.0_completed.zip`
* Scan scope: all `*.go` files under the repository, excluding `.git`
* Detection rule: `type <Name> struct {` and one-line `type <Name> struct{...}` declarations

---

## Summary

* Go files with struct declarations: 59
* Total struct declarations detected: 125
* Production struct declarations: 115
* Test/support struct declarations: 10

---

## Inventory

### Module: app

#### File: `internal/app/app.go`

* Lines: 202
* Struct declarations: 1

- `App` — line 22

---

### Module: core

#### File: `internal/core/auth/jwt.go`

* Lines: 99
* Struct declarations: 3

- `TokenService` — line 11
- `Claims` — line 17
- `MintOptions` — line 27

#### File: `internal/core/authorization/policy.go`

* Lines: 100
* Struct declarations: 2

- `Decision` — line 19
- `PolicyEvaluator` — line 26

#### File: `internal/core/authorization/subject.go`

* Lines: 34
* Struct declarations: 1

- `AuthorizationSubject` — line 5

#### File: `internal/core/cache/redis.go`

* Lines: 61
* Struct declarations: 1

- `Client` — line 16

#### File: `internal/core/config/config.go`

* Lines: 134
* Struct declarations: 1

- `Config` — line 10

#### File: `internal/core/db/postgres.go`

* Lines: 68
* Struct declarations: 1

- `Client` — line 17

#### File: `internal/core/errs/app_error.go`

* Lines: 95
* Struct declarations: 1

- `AppError` — line 15

#### File: `internal/core/errs/catalog_auth.go`

* Lines: 108
* Struct declarations: 1

- `AppErrorSpec` — line 8

#### File: `internal/core/errs/response_error.go`

* Lines: 26
* Struct declarations: 2

- `ResponseError` — line 3
- `ErrorEnvelope` — line 9

#### File: `internal/core/httpx/error_test.go`

* Lines: 100
* Struct declarations: 1

- `envelope` — line 13

#### File: `internal/core/httpx/middleware.go`

* Lines: 128
* Struct declarations: 1

- `wrapWriter` — line 60

#### File: `internal/core/httpx/router.go`

* Lines: 144
* Struct declarations: 1

- `RouterParams` — line 22

#### File: `internal/core/httpx/router_versioning_test.go`

* Lines: 215
* Struct declarations: 2

- `versionedErrorEnvelope` — line 19
- `versioningUserSettingsRepo` — line 27

#### File: `internal/core/logger/logger.go`

* Lines: 18
* Struct declarations: 1

- `Logger` — line 8

#### File: `internal/core/status/service.go`

* Lines: 139
* Struct declarations: 3

- `FuncChecker` — line 19
- `DependencyResult` — line 40
- `Service` — line 47

#### File: `internal/core/ws/client.go`

* Lines: 61
* Struct declarations: 1

- `Client` — line 11

#### File: `internal/core/ws/dispatcher.go`

* Lines: 75
* Struct declarations: 1

- `Dispatcher` — line 11

#### File: `internal/core/ws/handler.go`

* Lines: 140
* Struct declarations: 2

- `HandlerParams` — line 16
- `Handler` — line 23

#### File: `internal/core/ws/hub.go`

* Lines: 54
* Struct declarations: 1

- `Hub` — line 9

#### File: `internal/core/ws/protocol.go`

* Lines: 31
* Struct declarations: 2

- `Envelope` — line 20
- `ErrPayload` — line 28

#### File: `internal/core/ws/session.go`

* Lines: 32
* Struct declarations: 1

- `Session` — line 9

---

### Module: auth

#### File: `internal/modules/auth/app/application.go`

* Lines: 283
* Struct declarations: 1

- `Application` — line 20

#### File: `internal/modules/auth/app/response_types.go`

* Lines: 156
* Struct declarations: 18

- `LoginResponse` — line 11
- `SessionView` — line 18
- `ProfileWalletView` — line 33
- `ProfileView` — line 42
- `MeResponse` — line 57
- `SessionResponse` — line 62
- `BootstrapWalletsView` — line 66
- `BootstrapResponse` — line 71
- `WalletReadModel` — line 79
- `WalletsQuery` — line 92
- `WalletsResponse` — line 105
- `WalletLinkChallengeResponse` — line 117
- `WalletLinkVerifyResponse` — line 121
- `WalletAccountMergeChallengeResponse` — line 127
- `WalletAccountMergeVerifyResponse` — line 131
- `WalletDetachCheckResponse` — line 139
- `WalletDetachExecuteResponse` — line 147
- `WalletPrimarySetResponse` — line 153

#### File: `internal/modules/auth/app/service.go`

* Lines: 226
* Struct declarations: 2

- `Service` — line 24
- `LoginResult` — line 30

#### File: `internal/modules/auth/app/support_helpers.go`

* Lines: 436
* Struct declarations: 1

- `authenticatedContextView` — line 14

#### File: `internal/modules/auth/app/wallet_crypto.go`

* Lines: 187
* Struct declarations: 1

- `secp256k1Curve` — line 16

#### File: `internal/modules/auth/app/wallet_services.go`

* Lines: 538
* Struct declarations: 10

- `WalletChallengeService` — line 42
- `WalletLinkingService` — line 48
- `WalletLinkResult` — line 53
- `WalletAccountMergeService` — line 59
- `WalletAccountMergeResult` — line 64
- `WalletPrimaryService` — line 72
- `WalletPrimaryResult` — line 76
- `WalletDetachCheckResult` — line 81
- `WalletDetachExecuteResult` — line 89
- `WalletDetachService` — line 95

#### File: `internal/modules/auth/application.go`

* Lines: 489
* Struct declarations: 2

- `Application` — line 20
- `walletChallengeStoreAdapter` — line 233

#### File: `internal/modules/auth/auth_context.go`

* Lines: 43
* Struct declarations: 1

- `authenticatedContextView` — line 9

#### File: `internal/modules/auth/domain/contracts.go`

* Lines: 71
* Struct declarations: 3

- `WalletIdentity` — line 28
- `WalletChallenge` — line 48
- `WalletChallengeOptions` — line 68

#### File: `internal/modules/auth/error_adapter.go`

* Lines: 16
* Struct declarations: 1

- `authErrorSpec` — line 5

#### File: `internal/modules/auth/http_handlers_test.go`

* Lines: 3799
* Struct declarations: 4

- `errorPayload` — line 20
- `errorBody` — line 24
- `stubUserSettingsRepo` — line 30
- `walletDetachConflictPayload` — line 2121

#### File: `internal/modules/auth/http_login.go`

* Lines: 320
* Struct declarations: 8

- `LoginRequest` — line 16
- `LoginResponse` — line 21
- `UpdateMeRequest` — line 28
- `MeResponse` — line 32
- `SessionResponse` — line 37
- `UpdateMeSettingsRequest` — line 41
- `MeSettingsResponse` — line 45
- `HTTPHandlers` — line 49

#### File: `internal/modules/auth/http_wallet.go`

* Lines: 432
* Struct declarations: 11

- `WalletChallengeRequest` — line 18
- `WalletChallengeResponse` — line 23
- `WalletVerifyRequest` — line 27
- `WalletVerifyResponse` — line 33
- `WalletLinkChallengeRequest` — line 113
- `WalletLinkVerifyRequest` — line 120
- `WalletAccountMergeVerifyRequest` — line 128
- `WalletDetachCheckRequest` — line 136
- `WalletDetachExecuteRequest` — line 142
- `WalletPrimarySetRequest` — line 148
- `WalletAccountMergeChallengeRequest` — line 156

#### File: `internal/modules/auth/service.go`

* Lines: 168
* Struct declarations: 2

- `Service` — line 19
- `LoginResult` — line 25

#### File: `internal/modules/auth/service_test.go`

* Lines: 211
* Struct declarations: 1

- `stubUserRepo` — line 14

#### File: `internal/modules/auth/wallet_challenge.go`

* Lines: 251
* Struct declarations: 1

- `WalletChallengeService` — line 59

#### File: `internal/modules/auth/wallet_challenge_store_memory.go`

* Lines: 115
* Struct declarations: 1

- `InMemoryWalletChallengeStore` — line 28

#### File: `internal/modules/auth/wallet_challenge_store_pg.go`

* Lines: 232
* Struct declarations: 1

- `WalletChallengeStorePG` — line 31

#### File: `internal/modules/auth/wallet_crypto.go`

* Lines: 362
* Struct declarations: 1

- `secp256k1Curve` — line 21

#### File: `internal/modules/auth/wallet_detach.go`

* Lines: 104
* Struct declarations: 3

- `WalletDetachCheckResult` — line 8
- `WalletDetachExecuteResult` — line 16
- `WalletDetachService` — line 22

#### File: `internal/modules/auth/wallet_identity_store_memory.go`

* Lines: 328
* Struct declarations: 1

- `InMemoryWalletIdentityStore` — line 29

#### File: `internal/modules/auth/wallet_identity_store_pg.go`

* Lines: 572
* Struct declarations: 1

- `WalletIdentityStorePG` — line 32

#### File: `internal/modules/auth/wallet_link.go`

* Lines: 117
* Struct declarations: 2

- `WalletLinkingService` — line 9
- `WalletLinkResult` — line 14

#### File: `internal/modules/auth/wallet_merge.go`

* Lines: 133
* Struct declarations: 2

- `WalletAccountMergeService` — line 9
- `WalletAccountMergeResult` — line 14

#### File: `internal/modules/auth/wallet_primary.go`

* Lines: 49
* Struct declarations: 2

- `WalletPrimaryService` — line 8
- `WalletPrimaryResult` — line 12

#### File: `internal/modules/auth/wallet_verify.go`

* Lines: 125
* Struct declarations: 1

- `WalletVerificationService` — line 12

#### File: `internal/modules/auth/ws_handlers.go`

* Lines: 103
* Struct declarations: 1

- `WSHandlers` — line 9

---

### Module: user

#### File: `internal/modules/user/app/service.go`

* Lines: 153
* Struct declarations: 1

- `Service` — line 20

#### File: `internal/modules/user/domain/model.go`

* Lines: 12
* Struct declarations: 1

- `User` — line 5

#### File: `internal/modules/user/repository/postgres.go`

* Lines: 198
* Struct declarations: 1

- `PostgresRepository` — line 17

#### File: `internal/modules/user/service_test.go`

* Lines: 241
* Struct declarations: 1

- `stubRepo` — line 10

---

### Module: usersettings

#### File: `internal/modules/usersettings/app/service.go`

* Lines: 359
* Struct declarations: 1

- `Service` — line 29

#### File: `internal/modules/usersettings/domain/model.go`

* Lines: 56
* Struct declarations: 2

- `UserSettings` — line 5
- `View` — line 12

#### File: `internal/modules/usersettings/repository/postgres.go`

* Lines: 120
* Struct declarations: 1

- `PostgresRepository` — line 13

#### File: `internal/modules/usersettings/service_test.go`

* Lines: 556
* Struct declarations: 1

- `stubRepository` — line 10

---

### Module: thirdparty

#### File: `internal/thirdparty/sha3local/sha3.go`

* Lines: 185
* Struct declarations: 1

- `state` — line 23

#### File: `internal/thirdparty/sha3local/shake.go`

* Lines: 174
* Struct declarations: 1

- `cshakeState` — line 42

---

## Notes

* This inventory intentionally avoids READ / WRITE / HYBRID classification.
* Classification is reserved for Phase 0.12.1.2.
* Cross-layer usage analysis is reserved for Phase 0.12.1.3.
* Problem and risk mapping is reserved for Phase 0.12.1.4.
* Target separation definition is reserved for Phase 0.12.1.5.

---

## Status

* Phase: 0.12.1
* Subphase: 0.12.1.1
* Status: completed
* Code changes: none
* Documentation output: `docs/phase0_12_1_model_inventory.md`
