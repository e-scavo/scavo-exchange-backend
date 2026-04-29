# Phase 0.12.5 — Contract Alignment

## Subphase 0.12.5.2 — Contract Normalization Design

---

## Objective

Define the normalization strategy for all contract surfaces inventoried in Phase 0.12.5.1, without changing runtime behavior or modifying code.

This subphase defines how contracts must remain compatible while becoming structurally aligned with the explicit model lifecycle introduced across Phase 0.12:

```text
WRITE → DOMAIN → READ
```

No code changes are performed in this subphase.

---

## Source Baseline

Baseline ZIP: `scavo-exchange-backend-0.12.5.1.zip`

Direct input documents:

* `docs/phase0_12_5_contract_alignment.md`
* `docs/phase0_12_5_contract_inventory.md`
* `docs/phase0_12_4_mapping_layer.md`
* `docs/phase0_12_read_write_model_separation.md`

Relevant completed implementation phases:

* Phase 0.12.2 — Read Model Extraction
* Phase 0.12.3 — Write Model Isolation
* Phase 0.12.4 — Mapping Layer Introduction
* Phase 0.12.5.1 — Contract Inventory & Classification

---

## Design Goal

Normalize contracts without breaking compatibility.

The system must retain current HTTP names and JSON shapes while making the underlying lifecycle explicit.

The target is:

```text
HTTP Compatibility Contract → Write Model → Domain Input / Domain Model → Read Model → HTTP Compatibility Contract
```

This means:

* external request names may remain stable
* external response names may remain stable
* internal implementation must be aligned to write/read/domain contracts
* compatibility aliases are allowed where they preserve existing contracts safely

---

## Normalization Principles

### 1. Compatibility names are allowed

Existing public or handler-facing contract names may remain in place when they preserve compatibility.

Examples:

* `LoginRequest`
* `WalletVerifyRequest`
* `WalletLinkChallengeResponse`
* `WalletPrimarySetResponse`

These names can remain as aliases or wrapper contracts when they point to explicit lifecycle models.

---

### 2. HTTP request contracts must normalize to Write Models

All HTTP input contracts must eventually resolve to explicit write models.

Target pattern:

```text
HTTP Request Alias → Write Model → Domain Input
```

Example:

```text
LoginRequest → AuthLoginWriteModel → LoginInput
```

---

### 3. HTTP response contracts must normalize to Read Models

All HTTP output contracts must eventually resolve to explicit read models or stable read response wrappers.

Target pattern:

```text
Domain Model / Domain State → Read Model → HTTP Response
```

Example:

```text
User → UserReadModel → MeResponse
```

---

### 4. Domain contracts must remain internal

Domain models, domain inputs and provider interfaces must not become HTTP contracts.

Domain contracts may be used internally by:

* application services
* provider boundaries
* repositories
* mapping layer

They must not leak directly into handler response payloads unless intentionally preserved as a legacy compatibility surface.

---

### 5. Legacy aliases must be explicit

Legacy aliases are acceptable only if they are clearly intentional compatibility contracts.

They must not hide unclear lifecycle ownership.

A legacy alias must clearly point to one of:

* explicit write model
* explicit read model
* stable application read response

---

### 6. Hybrid contracts must be isolated

Any remaining hybrid contract must be explicitly scoped and isolated.

For Phase 0.12.5, hybrid contracts are not automatically refactored unless they are part of the auth/user/usersettings contract alignment scope.

---

## Contract Groups and Target Normalization

### Auth HTTP Request Contracts

Source files:

* `internal/modules/auth/http_login.go`
* `internal/modules/auth/http_wallet.go`

Current state:

* request names are preserved as HTTP-facing aliases
* underlying models are explicit write models

Target normalization:

| HTTP Contract | Normalized Target | Final Category | Treatment |
|---|---|---|---|
| `LoginRequest` | `AuthLoginWriteModel` | LEGACY / WRITE | Keep alias. |
| `UpdateMeRequest` | `AuthUpdateProfileWriteModel` | LEGACY / WRITE | Keep alias. |
| `UpdateMeSettingsRequest` | `AuthUpdateSettingsWriteModel` | LEGACY / WRITE | Keep alias. |
| `WalletChallengeRequest` | `AuthWalletChallengeWriteModel` | LEGACY / WRITE | Keep alias. |
| `WalletVerifyRequest` | `AuthWalletVerifyWriteModel` | LEGACY / WRITE | Keep alias. |
| `WalletLinkChallengeRequest` | `AuthWalletLinkChallengeWriteModel` | LEGACY / WRITE | Keep alias. |
| `WalletLinkVerifyRequest` | `AuthWalletLinkVerifyWriteModel` | LEGACY / WRITE | Keep alias. |
| `WalletAccountMergeChallengeRequest` | `AuthWalletAccountMergeChallengeWriteModel` | LEGACY / WRITE | Keep alias. |
| `WalletAccountMergeVerifyRequest` | `AuthWalletAccountMergeVerifyWriteModel` | LEGACY / WRITE | Keep alias. |
| `WalletDetachCheckRequest` | `AuthWalletDetachCheckWriteModel` | LEGACY / WRITE | Keep alias. |
| `WalletDetachExecuteRequest` | `AuthWalletDetachExecuteWriteModel` | LEGACY / WRITE | Keep alias. |
| `WalletPrimarySetRequest` | `AuthWalletPrimarySetWriteModel` | LEGACY / WRITE | Keep alias. |

Design decision:

These contracts are already safely normalized through write model aliases. Phase 0.12.5 should not rename them or alter JSON behavior.

---

### Auth HTTP Response Contracts

Source files:

* `internal/modules/auth/http_login.go`
* `internal/modules/auth/http_wallet.go`
* `internal/modules/auth/http_wallet_list.go`

Current state:

* response names are preserved
* most app-level response contracts are read-oriented
* some response names are aliases to application response contracts

Target normalization:

| HTTP Contract | Normalized Target | Final Category | Treatment |
|---|---|---|---|
| `LoginResponse` | login read response | READ | Keep response name. |
| `MeResponse` | `UserReadModel` within stable wrapper | READ | Keep wrapper. |
| `SessionResponse` | session read response | READ | Keep wrapper. |
| `MeSettingsResponse` | settings read wrapper | READ | Keep wrapper. |
| `WalletChallengeResponse` | challenge read response | READ | Keep response name. |
| `WalletVerifyResponse` | wallet verification read response | READ | Keep response name. |
| `WalletLinkChallengeResponse` | app read response | LEGACY / READ | Keep alias. |
| `WalletLinkVerifyResponse` | app read response | LEGACY / READ | Keep alias. |
| `WalletAccountMergeChallengeResponse` | app read response | LEGACY / READ | Keep alias. |
| `WalletAccountMergeVerifyResponse` | app read response | LEGACY / READ | Keep alias. |
| `WalletDetachCheckResponse` | app read response | LEGACY / READ | Keep alias. |
| `WalletDetachExecuteResponse` | app read response | LEGACY / READ | Keep alias. |
| `WalletPrimarySetResponse` | app read response | LEGACY / READ | Keep alias. |

Design decision:

Response contract names should remain stable. Normalization must happen behind the response boundary through read models and mapping layer functions.

---

### Auth Application Contracts

Source file:

* `internal/modules/auth/app/response_types.go`

Current state:

* application responses are read-oriented
* some response names are compatibility aliases to read models
* `WalletsQuery` remains an application input/query contract

Target normalization:

| Contract | Normalized Target | Final Category | Treatment |
|---|---|---|---|
| `LoginResponse` | `AuthLoginReadModel` | LEGACY / READ | Keep alias. |
| `SessionView` | stable app read view | READ | Keep. |
| `ProfileWalletView` | stable app read view | READ | Keep. |
| `ProfileView` | stable app read view | READ | Keep. |
| `MeResponse` | stable app read response | READ | Keep. |
| `SessionResponse` | stable app read response | READ | Keep. |
| `BootstrapWalletsView` | stable app read view | READ | Keep. |
| `BootstrapResponse` | stable app read response | READ | Keep. |
| `WalletReadModel` | `AuthWalletReadModel` | LEGACY / READ | Keep alias. |
| `WalletsQuery` | explicit query/write input | WRITE | Keep for now. |
| wallet operation responses | stable app read responses | READ | Keep. |

Design decision:

Application response contracts are already read-aligned. Phase 0.12.5 should document and preserve them, not collapse them into domain models.

---

### Explicit Write Models

Source files:

* `internal/modules/auth/writemodels/write_models.go`
* `internal/modules/user/writemodels/write_models.go`
* `internal/modules/usersettings/writemodels/write_models.go`

Target normalization:

Explicit write models remain the normalized input lifecycle layer.

They must:

* remain input-only
* not be reused as responses
* map explicitly through `internal/modules/<module>/mappers`
* remain independent of HTTP response contracts

Design decision:

No renaming is required in Phase 0.12.5.

---

### Explicit Read Models

Source files:

* `internal/modules/auth/readmodels/readmodels.go`
* `internal/modules/user/readmodels/readmodels.go`
* `internal/modules/usersettings/readmodels/readmodels.go`

Target normalization:

Explicit read models remain the normalized output lifecycle layer.

They must:

* remain output-only
* not be reused as request payloads
* be populated through mapping layer functions
* preserve compatibility with current HTTP responses

Design decision:

No renaming is required in Phase 0.12.5.

---

### Domain Write Input Contracts

Source files:

* `internal/modules/auth/domain/write_inputs.go`
* `internal/modules/user/domain/write_inputs.go`
* `internal/modules/usersettings/domain/write_inputs.go`

Target normalization:

Domain input contracts are the internal write boundary after HTTP/write model mapping.

Target flow:

```text
HTTP Request → Write Model → Domain Input → Application / Domain Logic
```

Design decision:

Domain write inputs should remain internal and must not become handler request types.

---

### Domain State Contracts

Source files:

* `internal/modules/auth/domain/contracts.go`
* `internal/modules/user/domain/model.go`
* `internal/modules/usersettings/domain/model.go`

Target normalization:

Domain state contracts remain internal state representations.

They must not be treated as direct read or write HTTP contracts.

Specific treatment:

| Contract | Target Category | Treatment |
|---|---|---|
| `WalletIdentity` | INTERNAL | Keep domain state. |
| `WalletChallenge` | INTERNAL | Keep domain state. |
| `WalletChallengeOptions` | WRITE / INTERNAL | Keep internal options contract. |
| `user/domain.User` | INTERNAL | Keep domain model. |
| `usersettings/domain.UserSettings` | INTERNAL | Keep domain model. |
| `usersettings/domain.View` | LEGACY / READ | Keep temporarily as compatibility read view unless replaced later. |

Design decision:

Domain models should remain behind mapping boundaries. `usersettings/domain.View` is the main candidate for later cleanup, but should not be removed in this subphase.

---

### Provider Interfaces

Source files:

* `internal/modules/auth/domain/user_contract.go`
* `internal/modules/auth/domain/usersettings_contract.go`

Current state:

* `UserProvider` and `UserSettingsProvider` are internal auth dependencies
* they are not HTTP contracts
* they are alignment-sensitive because they can leak domain/read shapes across module boundaries

Target normalization:

Provider interfaces should expose stable internal contracts that do not force HTTP response models across module boundaries.

Recommended direction:

| Provider | Target Treatment |
|---|---|
| `UserProvider` | Continue returning auth-needed user data through a stable internal/read-compatible type. |
| `UserSettingsProvider` | Continue returning settings data through a stable internal/read-compatible type. |

Design decision:

Provider contracts should be reviewed in implementation subphases, but no breaking change should be introduced unless tests and imports confirm there is no cycle or response contract break.

---

### Repository Interfaces

Source files:

* `internal/modules/user/domain/contracts.go`
* `internal/modules/usersettings/domain/contracts.go`

Target normalization:

Repository interfaces remain internal persistence contracts.

They do not need HTTP-oriented normalization.

Design decision:

No Phase 0.12.5 implementation change is required unless provider alignment exposes a specific mismatch.

---

### Core Protocol and Error Contracts

Source files:

* `internal/core/ws/protocol.go`
* `internal/core/errs/response_error.go`

Target normalization:

| Contract | Current Classification | Target Treatment |
|---|---|---|
| `Envelope` | HYBRID | Keep out of immediate Phase 0.12.5 alignment scope. |
| `ErrPayload` | READ | Keep. |
| `ResponseError` | READ | Keep. |
| `ErrorEnvelope` | READ | Keep. |

Design decision:

Core protocol/error contracts are stable foundation contracts. They must not be changed as part of auth/user/usersettings contract normalization.

---

## Normalization Matrix

| Contract Area | Current Issue | Normalization Target | Implementation Treatment |
|---|---|---|---|
| Auth HTTP requests | legacy names | alias to write models | Preserve aliases. |
| Auth HTTP responses | stable response names | read models / read wrappers | Preserve wrappers. |
| App responses | read-oriented but mixed naming | stable read contracts | Preserve. |
| App query input | query input inside app response file | write/query input | Keep; document as write/query. |
| Write models | explicit input models | normalized write layer | Preserve. |
| Read models | explicit output models | normalized read layer | Preserve. |
| Domain inputs | internal write boundary | normalized domain input layer | Preserve. |
| Domain state | internal domain state | internal only | Preserve behind mappers. |
| Provider contracts | module boundary sensitive | stable internal contract | Review in implementation. |
| WebSocket envelope | hybrid transport envelope | out of scope | Preserve. |

---

## Implementation Guidance for 0.12.5.3+

### Required

* preserve all current HTTP names
* preserve all current JSON tags
* preserve all current routes
* preserve all current status/error behavior
* use mapping layer as the normalized transformation point
* keep compatibility aliases where they already provide safe migration

### Allowed

* adjust internal provider contracts if required and safe
* add internal helper methods for contract normalization
* document legacy compatibility aliases explicitly
* reduce direct usage of legacy domain views if a read model equivalent is available

### Prohibited

* deleting compatibility aliases
* renaming request/response structs used by handlers
* changing JSON field names
* changing endpoint paths
* changing API versioning
* changing WebSocket protocol
* introducing CQRS/event sourcing

---

## Contract Normalization Flow

### Write Flow

```text
HTTP Request Name
  ↓
Write Model
  ↓
Mapping Layer
  ↓
Domain Input
  ↓
Application / Domain Logic
```

### Read Flow

```text
Domain State / Application State
  ↓
Mapping Layer
  ↓
Read Model / Read View
  ↓
HTTP Response Name
```

---

## Specific Alignment Candidates

### Candidate 1 — `WalletsQuery`

Current classification:

* WRITE query input

Target:

* keep as app-level query input for now
* do not convert to read model
* optionally move to a query-specific file in a future cleanup phase

Phase 0.12.5 treatment:

* document and preserve

---

### Candidate 2 — `usersettings/domain.View`

Current classification:

* LEGACY / READ

Target:

* keep as compatibility read view unless replacement is safe
* prefer `UserSettingsReadModel` at new mapping boundaries

Phase 0.12.5 treatment:

* avoid direct expansion of its usage
* review whether provider outputs can align with read model semantics

---

### Candidate 3 — Provider Interfaces

Current classification:

* INTERNAL
* alignment-sensitive

Target:

* avoid leaking HTTP response contracts
* avoid forcing domain model exposure into auth response composition
* preserve import boundaries

Phase 0.12.5 treatment:

* review implementation-level usage
* adjust only if safe and test-backed

---

### Candidate 4 — Compatibility Aliases

Current classification:

* LEGACY / WRITE
* LEGACY / READ

Target:

* keep aliases as stable compatibility surface
* make underlying lifecycle explicit

Phase 0.12.5 treatment:

* preserve aliases
* document as intentional compatibility contracts

---

## Expected Output of Phase 0.12.5.3

The next implementation subphase should produce one or more of the following, only where needed:

* explicit internal normalization helpers
* provider contract alignment
* read/write compatibility documentation in code structure
* reduced direct legacy contract usage

It should not produce:

* new API contracts
* endpoint renames
* JSON changes
* broad handler rewrites

---

## Compatibility Statement

Contract normalization is an internal structural alignment.

External behavior must remain unchanged:

* same routes
* same request names
* same response JSON
* same status behavior
* same authentication and authorization behavior

---

## Status

Phase: 0.12.5  
Subphase: 0.12.5.2  
Status: COMPLETED  
Code Impact: NONE  
Next: 0.12.5.3 — Contract Alignment Implementation
