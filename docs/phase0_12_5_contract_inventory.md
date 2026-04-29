# Phase 0.12.5 — Contract Alignment

## Subphase 0.12.5.1 — Contract Inventory & Classification

---

## Objective

Inventory and classify the current contract surface after the Read / Write Model separation and Mapping Layer consolidation completed in Phase 0.12.2, Phase 0.12.3 and Phase 0.12.4.

This subphase identifies the contracts that must be preserved, normalized or reviewed before implementation-level contract alignment begins.

No code changes are performed in this subphase.

---

## Source Baseline

Baseline ZIP: `scavo-exchange-backend-0.12.5.0.zip`

Relevant completed phases:

* Phase 0.12.1 — Model Classification Audit
* Phase 0.12.2 — Read Model Extraction
* Phase 0.12.3 — Write Model Isolation
* Phase 0.12.4 — Mapping Layer Introduction
* Phase 0.12.5.0 — Contract Alignment Definition & Documentation Lock

---

## Classification Rules

Contracts are classified as:

* `READ` — output-facing contract, response model, view, projection or error payload
* `WRITE` — input-facing contract, request model, command payload or write intent
* `HYBRID` — contract still used across more than one lifecycle direction or boundary
* `LEGACY` — compatibility alias or compatibility surface intentionally kept to preserve the current HTTP contract
* `INTERNAL` — provider, repository, service, store, infrastructure or domain-only contract not directly exposed as request/response

---

## Summary

| Contract Group | Count |
|---|---:|
| HTTP request aliases / write contracts | 14 |
| HTTP response contracts / read contracts | 19 |
| Application response/view contracts | 18 |
| Explicit write models | 15 |
| Explicit read models | 5 |
| Domain provider/store/repository interfaces | 6 |
| Core protocol/error contracts | 4 |
| Domain state/input contracts | 13 |

---

## Inventory — Auth HTTP Contracts

### File: `internal/modules/auth/http_login.go`

| Contract | Kind | Classification | Notes |
|---|---|---|---|
| `LoginRequest` | alias | LEGACY / WRITE | Alias to `authwritemodels.AuthLoginWriteModel`; preserves HTTP request name. |
| `LoginResponse` | struct | READ | HTTP login response; preserves JSON contract. |
| `UpdateMeRequest` | alias | LEGACY / WRITE | Alias to `authwritemodels.AuthUpdateProfileWriteModel`. |
| `MeResponse` | struct | READ | HTTP `/me` response using read model fields. |
| `SessionResponse` | struct | READ | HTTP session response. |
| `UpdateMeSettingsRequest` | alias | LEGACY / WRITE | Alias to `authwritemodels.AuthUpdateSettingsWriteModel`. |
| `MeSettingsResponse` | struct | READ | HTTP settings response. |
| `HTTPHandlers` | struct | INTERNAL | Handler composition object, not a request/response contract. |

### File: `internal/modules/auth/http_wallet.go`

| Contract | Kind | Classification | Notes |
|---|---|---|---|
| `WalletChallengeRequest` | alias | LEGACY / WRITE | Alias to `AuthWalletChallengeWriteModel`. |
| `WalletChallengeResponse` | struct | READ | HTTP wallet challenge response. |
| `WalletVerifyRequest` | alias | LEGACY / WRITE | Alias to `AuthWalletVerifyWriteModel`. |
| `WalletVerifyResponse` | struct | READ | HTTP wallet verification response. |
| `WalletLinkChallengeRequest` | alias | LEGACY / WRITE | Alias to `AuthWalletLinkChallengeWriteModel`. |
| `WalletLinkChallengeResponse` | alias | LEGACY / READ | Alias to app response contract. |
| `WalletLinkVerifyRequest` | alias | LEGACY / WRITE | Alias to `AuthWalletLinkVerifyWriteModel`. |
| `WalletLinkVerifyResponse` | alias | LEGACY / READ | Alias to app response contract. |
| `WalletAccountMergeChallengeRequest` | alias | LEGACY / WRITE | Alias to `AuthWalletAccountMergeChallengeWriteModel`. |
| `WalletAccountMergeChallengeResponse` | alias | LEGACY / READ | Alias to app response contract. |
| `WalletAccountMergeVerifyRequest` | alias | LEGACY / WRITE | Alias to `AuthWalletAccountMergeVerifyWriteModel`. |
| `WalletAccountMergeVerifyResponse` | alias | LEGACY / READ | Alias to app response contract. |
| `WalletDetachCheckRequest` | alias | LEGACY / WRITE | Alias to `AuthWalletDetachCheckWriteModel`. |
| `WalletDetachCheckResponse` | alias | LEGACY / READ | Alias to app response contract. |
| `WalletDetachExecuteRequest` | alias | LEGACY / WRITE | Alias to `AuthWalletDetachExecuteWriteModel`. |
| `WalletDetachExecuteResponse` | alias | LEGACY / READ | Alias to app response contract. |
| `WalletPrimarySetRequest` | alias | LEGACY / WRITE | Alias to `AuthWalletPrimarySetWriteModel`. |
| `WalletPrimarySetResponse` | alias | LEGACY / READ | Alias to app response contract. |

---

## Inventory — Auth Application Contracts

### File: `internal/modules/auth/app/response_types.go`

| Contract | Kind | Classification | Notes |
|---|---|---|---|
| `LoginResponse` | alias | LEGACY / READ | Alias to `authreadmodels.AuthLoginReadModel`; preserves app response name. |
| `SessionView` | struct | READ | Session projection. |
| `ProfileWalletView` | struct | READ | Wallet projection inside profile. |
| `ProfileView` | struct | READ | Authenticated profile view. |
| `MeResponse` | struct | READ | Application-level `/me` response. |
| `SessionResponse` | struct | READ | Application session response. |
| `BootstrapWalletsView` | struct | READ | Bootstrap wallet summary view. |
| `BootstrapResponse` | struct | READ | Authenticated bootstrap read response. |
| `WalletReadModel` | alias | LEGACY / READ | Alias to `authreadmodels.AuthWalletReadModel`. |
| `WalletsQuery` | struct | WRITE | Query input for wallet listing/filtering. |
| `WalletsResponse` | struct | READ | Wallet list response. |
| `WalletLinkChallengeResponse` | struct | READ | Link challenge response. |
| `WalletLinkVerifyResponse` | struct | READ | Link verification response. |
| `WalletAccountMergeChallengeResponse` | struct | READ | Merge challenge response. |
| `WalletAccountMergeVerifyResponse` | struct | READ | Merge verification response. |
| `WalletDetachCheckResponse` | struct | READ | Detach pre-check response. |
| `WalletDetachExecuteResponse` | struct | READ | Detach execution response. |
| `WalletPrimarySetResponse` | struct | READ | Primary wallet update response. |

---

## Inventory — Explicit Read Models

### File: `internal/modules/auth/readmodels/readmodels.go`

| Contract | Classification | Notes |
|---|---|---|
| `AuthLoginReadModel` | READ | Explicit login output model. |
| `AuthWalletReadModel` | READ | Explicit wallet output model. |
| `AuthWalletChallengeReadModel` | READ | Explicit wallet challenge output model. |

### File: `internal/modules/user/readmodels/readmodels.go`

| Contract | Classification | Notes |
|---|---|---|
| `UserReadModel` | READ | Explicit user output model. |

### File: `internal/modules/usersettings/readmodels/readmodels.go`

| Contract | Classification | Notes |
|---|---|---|
| `UserSettingsReadModel` | READ | Explicit user settings output model. |

---

## Inventory — Explicit Write Models

### File: `internal/modules/auth/writemodels/write_models.go`

| Contract | Classification | Notes |
|---|---|---|
| `AuthLoginWriteModel` | WRITE | Login input. |
| `AuthUpdateProfileWriteModel` | WRITE | Profile update input. |
| `AuthUpdateSettingsWriteModel` | WRITE | Settings update input. |
| `AuthWalletChallengeWriteModel` | WRITE | Wallet challenge input. |
| `AuthWalletVerifyWriteModel` | WRITE | Wallet verify input. |
| `AuthWalletLinkChallengeWriteModel` | WRITE | Wallet link challenge input. |
| `AuthWalletLinkVerifyWriteModel` | WRITE | Wallet link verification input. |
| `AuthWalletAccountMergeChallengeWriteModel` | WRITE | Account merge challenge input. |
| `AuthWalletAccountMergeVerifyWriteModel` | WRITE | Account merge verification input. |
| `AuthWalletDetachCheckWriteModel` | WRITE | Wallet detach check input. |
| `AuthWalletDetachExecuteWriteModel` | WRITE | Wallet detach execution input. |
| `AuthWalletPrimarySetWriteModel` | WRITE | Primary wallet set input. |

### File: `internal/modules/user/writemodels/write_models.go`

| Contract | Classification | Notes |
|---|---|---|
| `UserUpdateWriteModel` | WRITE | User update input. |

### File: `internal/modules/usersettings/writemodels/write_models.go`

| Contract | Classification | Notes |
|---|---|---|
| `UserSettingsUpdateWriteModel` | WRITE | User settings update input. |

---

## Inventory — Domain Write Input Contracts

### File: `internal/modules/auth/domain/write_inputs.go`

| Contract | Classification | Notes |
|---|---|---|
| `LoginInput` | WRITE | Domain-level login input. |
| `ProfileUpdateInput` | WRITE | Domain-level profile update input. |
| `SettingsUpdateInput` | WRITE | Domain-level settings update input. |
| `WalletChallengeInput` | WRITE | Domain-level wallet challenge input. |
| `WalletVerifyInput` | WRITE | Domain-level wallet verification input. |
| `WalletDetachInput` | WRITE | Domain-level detach input. |
| `WalletPrimarySetInput` | WRITE | Domain-level primary wallet input. |

### File: `internal/modules/user/domain/write_inputs.go`

| Contract | Classification | Notes |
|---|---|---|
| `UserUpdateInput` | WRITE | Domain-level user update input. |

### File: `internal/modules/usersettings/domain/write_inputs.go`

| Contract | Classification | Notes |
|---|---|---|
| `UserSettingsUpdateInput` | WRITE | Domain-level settings update input. |

---

## Inventory — Domain State Contracts

### File: `internal/modules/auth/domain/contracts.go`

| Contract | Kind | Classification | Notes |
|---|---|---|---|
| `WalletIdentity` | struct | INTERNAL | Domain wallet identity state. |
| `WalletIdentityStore` | interface | INTERNAL | Wallet identity persistence contract. |
| `WalletChallenge` | struct | INTERNAL | Domain wallet challenge state. |
| `WalletChallengeStore` | interface | INTERNAL | Wallet challenge persistence contract. |
| `WalletChallengeOptions` | struct | WRITE / INTERNAL | Domain options used to create challenges. |

### File: `internal/modules/user/domain/model.go`

| Contract | Classification | Notes |
|---|---|---|
| `User` | INTERNAL | Domain user model, no longer intended as direct response/write contract. |

### File: `internal/modules/usersettings/domain/model.go`

| Contract | Classification | Notes |
|---|---|---|
| `UserSettings` | INTERNAL | Domain settings model. |
| `View` | LEGACY / READ | Existing output-compatible settings view retained for compatibility. |

---

## Inventory — Domain Provider and Repository Contracts

### File: `internal/modules/auth/domain/user_contract.go`

| Contract | Classification | Notes |
|---|---|---|
| `UserProvider` | INTERNAL | Auth-to-user provider boundary. Contract alignment candidate. |

### File: `internal/modules/auth/domain/usersettings_contract.go`

| Contract | Classification | Notes |
|---|---|---|
| `UserSettingsProvider` | INTERNAL | Auth-to-usersettings provider boundary. Contract alignment candidate. |

### File: `internal/modules/user/domain/contracts.go`

| Contract | Classification | Notes |
|---|---|---|
| `Repository` | INTERNAL | User persistence boundary. |

### File: `internal/modules/usersettings/domain/contracts.go`

| Contract | Classification | Notes |
|---|---|---|
| `Repository` | INTERNAL | User settings persistence boundary. |

---

## Inventory — Core Protocol and Error Contracts

### File: `internal/core/ws/protocol.go`

| Contract | Classification | Notes |
|---|---|---|
| `Envelope` | HYBRID | WebSocket transport envelope; contains both inbound and outbound envelope semantics. |
| `ErrPayload` | READ | WebSocket error payload. |

### File: `internal/core/errs/response_error.go`

| Contract | Classification | Notes |
|---|---|---|
| `ResponseError` | READ | Standardized error response payload. |
| `ErrorEnvelope` | READ | Standardized error envelope. |

---

## Contract Alignment Findings

### Finding 1 — HTTP aliases are intentional compatibility contracts

Many HTTP request names are now aliases to explicit write models.

This is correct and should be preserved during alignment unless a future versioned contract intentionally replaces them.

### Finding 2 — Application response aliases are intentional compatibility contracts

Application response names alias read models where possible while preserving historical naming and JSON contracts.

### Finding 3 — Provider interfaces remain internal but are alignment-sensitive

`UserProvider` and `UserSettingsProvider` are not HTTP contracts, but they are critical internal contracts because auth depends on them.

They must be aligned carefully in 0.12.5 without introducing import cycles.

### Finding 4 — WebSocket envelope remains hybrid

`internal/core/ws.Envelope` is still classified as HYBRID because it represents a transport envelope for more than one direction.

It should not be changed in Phase 0.12.5 unless contract alignment explicitly scopes it.

### Finding 5 — Legacy names remain useful

Compatibility aliases currently provide a safe transition path:

* public/internal call sites keep existing names
* implementation points to read/write models
* JSON contracts remain stable

---

## Contract Alignment Candidates for 0.12.5.2+

| Area | Candidate | Recommended Treatment |
|---|---|---|
| Auth HTTP requests | request aliases | Keep aliases; document as compatibility layer. |
| Auth app responses | response aliases | Keep aliases; document as read contract compatibility. |
| Auth provider interfaces | `UserProvider`, `UserSettingsProvider` | Review return/input types for read/write separation consistency. |
| App query contracts | `WalletsQuery` | Review as WRITE/query input contract. |
| Usersettings view | `domain.View` | Review whether it remains legacy read view or maps to `UserSettingsReadModel`. |
| WebSocket envelope | `Envelope` | Keep out of immediate refactor unless scoped later. |

---

## Out of Scope for This Subphase

This subphase does not:

* change code
* remove aliases
* modify handlers
* change JSON fields
* change provider interfaces
* change WebSocket protocol
* introduce new API versions

---

## Status

Phase: 0.12.5  
Subphase: 0.12.5.1  
Status: COMPLETED  
Code Impact: NONE  
Next: 0.12.5.2 — Contract Normalization Design
