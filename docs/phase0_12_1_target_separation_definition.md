# Phase 0.12.1.5 — Target Separation Definition

## Purpose

This document defines the intended target separation for the hybrid and controlled boundary models identified during Phase 0.12.1.2, Phase 0.12.1.3, and Phase 0.12.1.4.

This subphase does not modify code, does not create read models, does not create write models, does not introduce mapping functions, and does not change runtime behavior. Its only purpose is to define the target shape that later implementation subphases may use.

---

## Source

* ZIP: `scavo-exchange-backend-0.12.1.4.zip`
* Baseline inventory: `docs/phase0_12_1_model_inventory.md`
* Baseline classification: `docs/phase0_12_1_model_classification.md`
* Baseline cross-layer analysis: `docs/phase0_12_1_cross_layer_usage_analysis.md`
* Baseline risk map: `docs/phase0_12_1_problem_detection_risk_mapping.md`
* Hybrid candidates reviewed: 11
* Code changes: none

---

## Scope

### Included

* Define the target separation for each CRITICAL and HIGH hybrid candidate.
* Define accepted treatment for controlled runtime and infrastructure boundary candidates.
* Record future extraction direction without implementing it.
* Preserve compatibility expectations for later phases.
* Identify which models should become read models, write models, domain models, runtime boundary models, or accepted exceptions.

### Excluded

* No code changes.
* No refactors.
* No package moves.
* No new Go structs.
* No model extraction.
* No mapping implementation.
* No endpoint changes.
* No API contract changes.
* No test rewrites.

---

## Target Separation Principles

Phase 0.12 separates model lifecycle responsibilities using the following target concepts.

| Target Type | Purpose |
|---|---|
| Domain Model | Internal business or persistence-oriented representation owned by a module. |
| Read Model | Output-oriented representation used for responses, projections, and public views. |
| Write Model | Input-oriented representation used for commands, mutations, updates, or patches. |
| Runtime Boundary Model | Controlled runtime object used across infrastructure boundaries, not a public read/write contract. |
| Accepted Exception | Explicitly documented hybrid retained for protocol or compatibility reasons. |
| Internal Metadata | Catalog, adapter, or configuration structure not participating in domain read/write lifecycle. |

---

## Summary

| Treatment | Count |
|---|---:|
| Target split required | 4 |
| Accepted transport exception / deferred split | 1 |
| Controlled runtime boundary | 4 |
| Internal metadata reclassification candidate | 2 |
| Total Reviewed | 11 |

---

## Target Separation Map

| Struct | Current Risk | Target Treatment | Future Direction |
|---|---|---|---|
| `User` | CRITICAL | Split required | Keep domain entity internal; introduce explicit user/auth read models; introduce write models only when user mutation flows exist. |
| `UserSettings` | CRITICAL | Split required | Keep settings domain model internal; introduce settings read model and settings write/update model. |
| `WalletIdentity` | CRITICAL | Split required | Keep wallet identity as internal auth domain/persistence model; map to wallet read models. |
| `WalletChallenge` | HIGH | Split required | Keep challenge runtime/persistence state internal; map to challenge response read model. |
| `Envelope` | HIGH | Accepted transport exception / deferred split | Retain as protocol-level envelope for now; document future inbound/outbound split candidate. |
| `Claims` | MEDIUM | Controlled runtime boundary | Keep as core auth token/runtime identity model; prevent read response dependency on raw token shape. |
| `AuthorizationSubject` | MEDIUM | Controlled runtime boundary | Keep as core authorization runtime subject. |
| `Session` | MEDIUM | Controlled runtime boundary | Keep as WS runtime session state; maintain separation from auth session read models. |
| `AppError` | MEDIUM | Controlled runtime boundary | Keep as standardized error boundary from Phase 0.8. |
| `AppErrorSpec` | LOW | Internal metadata | Reclassify as internal error catalog metadata during audit closure. |
| `authErrorSpec` | LOW | Internal metadata | Reclassify as auth-local adapter metadata during audit closure. |

---

## Required Target Splits

### 1. `User`

**Current state:**

`User` is the strongest cross-module hybrid. It is declared in `internal/modules/user/domain/model.go`, exposed through an alias in `internal/modules/user/model.go`, and consumed by user, auth, application, response construction, and tests.

**Target treatment:**

`User` should remain the user module domain/persistence entity.

Future phases should introduce explicit read models where user data leaves the domain boundary.

**Target model families:**

| Target | Description |
|---|---|
| `user/domain.User` | Canonical internal user domain and persistence model. |
| User provider read contract | Narrow auth-facing user representation returned by `UserProvider`-style contracts. |
| User response read model | Output shape for public user/profile/session/bootstrap responses when needed. |
| User write model | Not required yet unless user mutation flows are introduced. |

**Target flow:**

```text
user/domain.User
  -> provider-facing user read contract
  -> auth/user response read models
```

If future user mutation endpoints are introduced:

```text
user write model
  -> validation/application command
  -> user/domain.User
```

**Compatibility direction:**

* Existing public behavior must remain unchanged.
* Existing aliases may remain temporarily as compatibility surfaces.
* Later implementation should avoid returning or exposing domain `User` directly from transport-oriented response construction.

---

### 2. `UserSettings`

**Current state:**

`UserSettings` crosses usersettings domain, persistence, application, auth-facing contracts, HTTP flows, router wiring, and tests. It participates in read and update contexts.

**Target treatment:**

`UserSettings` should remain the usersettings domain/persistence model, while read and write responsibilities should be split explicitly.

**Target model families:**

| Target | Description |
|---|---|
| `usersettings/domain.UserSettings` | Canonical internal settings domain and persistence model. |
| User settings read model | Output-oriented settings representation for session/bootstrap/profile responses. |
| User settings write/update model | Input-oriented settings update or patch command. |
| Auth-facing settings contract | Narrow provider-facing settings view consumed by auth flows. |

**Target flow:**

```text
usersettings/domain.UserSettings
  -> settings read model
  -> auth/session/bootstrap read payloads
```

For updates:

```text
settings write/update model
  -> validation/application command
  -> usersettings/domain.UserSettings
```

**Compatibility direction:**

* Existing settings endpoint behavior must remain stable.
* Update request handling should not depend on read model shape.
* Read response construction should not require exposing the full domain/persistence model.

---

### 3. `WalletIdentity`

**Current state:**

`WalletIdentity` is declared in auth domain contracts, aliased in the auth wallet identity store package, and used by persistence stores, app services, wallet HTTP flows, response construction, and tests.

**Target treatment:**

`WalletIdentity` should remain an internal auth domain/persistence model. Public wallet output should be produced through explicit wallet read models.

**Target model families:**

| Target | Description |
|---|---|
| `auth/domain.WalletIdentity` | Canonical internal wallet identity domain/persistence model. |
| Wallet identity read model | Output-oriented wallet identity representation for wallet list/profile/session responses. |
| Wallet link/read result model | Read result shape for wallet link and verification flows if distinct from list output. |
| Wallet write model | Not required for wallet identity itself unless future mutation commands are introduced. |

**Target flow:**

```text
auth/domain.WalletIdentity
  -> wallet identity read model
  -> wallet/list/profile/session responses
```

For link or verification commands, request payloads should remain command/write-oriented and should not reuse wallet identity read models.

**Compatibility direction:**

* Existing wallet response payloads must remain stable.
* Internal wallet identity fields should not automatically become public response fields.
* Store aliases may remain temporarily for compatibility but should be treated as internal compatibility surfaces.

---

### 4. `WalletChallenge`

**Current state:**

`WalletChallenge` carries runtime and persistence challenge state and is also used in application flows and public challenge response construction.

**Target treatment:**

`WalletChallenge` should remain internal challenge runtime/persistence state. Public challenge output should be represented by a dedicated read model.

**Target model families:**

| Target | Description |
|---|---|
| `auth/domain.WalletChallenge` | Internal challenge state, nonce, expiration, ownership, and persistence lifecycle. |
| Wallet challenge response read model | Output-oriented challenge payload returned to clients. |
| Wallet challenge verify/write model | Input-oriented verification command payload, if applicable. |
| Wallet challenge link/write model | Input-oriented challenge/link command payload, if applicable. |

**Target flow:**

```text
auth/domain.WalletChallenge
  -> wallet challenge response read model
```

For verification:

```text
wallet verification write model
  -> validation/application command
  -> auth/domain.WalletChallenge lookup/state transition
```

**Compatibility direction:**

* Challenge response payloads must remain stable.
* Runtime/persistence challenge fields should not automatically leak into public output.
* Challenge input payloads should remain separate from challenge output payloads.

---

## Accepted Exception / Deferred Split

### 5. `Envelope`

**Current state:**

`Envelope` is a WS protocol shape used by the WS protocol, dispatcher, helpers, auth WS handlers, and system WS handlers. It represents both inbound and outbound protocol envelope semantics.

**Target treatment:**

`Envelope` may remain an accepted transport-level exception for now because it is a protocol envelope rather than a domain model.

However, it must be explicitly documented as a deferred split candidate.

**Deferred target options:**

| Option | Description |
|---|---|
| Keep single `Envelope` | Accepted protocol simplification if inbound/outbound semantics remain stable and symmetrical. |
| Split into inbound/outbound envelopes | Future option if request and response semantics diverge. |

**Potential future target flow:**

```text
InboundEnvelope
  -> dispatcher/handler command
  -> OutboundEnvelope
```

**Compatibility direction:**

* No immediate code change is required.
* The single-envelope design should not be used as precedent for domain read/write models.
* Any future split must preserve WS contract compatibility.

---

## Controlled Runtime Boundary Models

The following models cross layers but should not be treated as ordinary domain read/write hybrids during this phase. They are runtime or infrastructure boundary models and may remain in place with explicit constraints.

### 6. `Claims`

**Target treatment:** controlled core auth runtime boundary.

`Claims` should remain the JWT token payload and runtime identity carrier inside core auth flows.

**Constraint:**

Application read models and response payloads should not depend directly on raw token claim shape.

**Allowed flow:**

```text
JWT Claims
  -> authenticated runtime context
  -> narrower application/domain lookup inputs
```

---

### 7. `AuthorizationSubject`

**Target treatment:** controlled core authorization runtime boundary.

`AuthorizationSubject` should remain the normalized authorization subject used by authorization policies, context hydration, and HTTP middleware.

**Constraint:**

It must not become a domain read model or public output contract.

**Allowed flow:**

```text
Claims/context
  -> AuthorizationSubject
  -> policy evaluation
```

---

### 8. `Session`

**Target treatment:** controlled WS runtime session state.

`Session` should remain the WS client/session runtime object.

**Constraint:**

It must remain conceptually separate from auth session response/read models such as session view or session response payloads.

**Allowed flow:**

```text
WS connection
  -> WS runtime Session
  -> handler context
```

---

### 9. `AppError`

**Target treatment:** accepted standardized error boundary model from Phase 0.8.

`AppError` may continue bridging internal error semantics and standardized HTTP error response mapping.

**Constraint:**

Transport response formatting should remain centralized through the established error writer/mapper path.

**Allowed flow:**

```text
internal error condition
  -> AppError
  -> standardized HTTP error response
```

---

## Internal Metadata Reclassification Candidates

### 10. `AppErrorSpec`

**Target treatment:** internal error catalog metadata.

`AppErrorSpec` should be treated as catalog metadata, not as a lifecycle read/write model.

**Closure action:**

Phase 0.12.1.6 should reclassify this candidate as INTERNAL/catalog metadata unless new evidence appears.

---

### 11. `authErrorSpec`

**Target treatment:** auth-local adapter metadata.

`authErrorSpec` should be treated as local adapter metadata, not as a lifecycle read/write model.

**Closure action:**

Phase 0.12.1.6 should reclassify this candidate as INTERNAL/adapter-local unless new evidence appears.

---

## Alias Treatment

Several compatibility aliases were identified during previous audit subphases.

Examples include:

* `internal/modules/user/model.go`
* `internal/modules/usersettings/model.go`
* `internal/modules/auth/wallet_identity_store.go`
* `internal/modules/auth/wallet_challenge.go`

### Target Direction

Aliases may remain temporarily if required for compatibility, but they must be treated as compatibility surfaces rather than new ownership boundaries.

Later implementation subphases should avoid adding new dependencies to aliases when a clearer domain, read, write, or provider-facing contract is available.

---

## Implementation Guidance for Later Phases

Later implementation phases should follow this order:

1. Preserve runtime behavior and public response compatibility.
2. Introduce explicit read models close to response construction boundaries.
3. Introduce explicit write models only where mutation/input flows exist.
4. Add mapping functions instead of direct cross-layer field reuse.
5. Reduce dependence on compatibility aliases only after stable mapping exists.
6. Keep controlled runtime models separate from domain read/write models.

---

## Test Sensitivity Notes

Later phases must be careful with tests around:

* auth HTTP handlers;
* wallet challenge, link, verify, list, and profile flows;
* user service behavior;
* usersettings service behavior;
* router and versioning behavior;
* authorization enforcement behavior;
* WS session and handler behavior.

Tests currently validate behavior through existing hybrid shapes. Refactors must preserve behavior while changing internal model ownership and mapping paths.

---

## Output of This Subphase

Phase 0.12.1.5 defines target separation but does not implement it.

The target split required for later phases is:

```text
User
  -> domain entity
  -> provider-facing read contract
  -> response read models

UserSettings
  -> domain model
  -> settings read model
  -> settings write/update model

WalletIdentity
  -> internal auth domain/persistence model
  -> wallet read models

WalletChallenge
  -> internal challenge runtime/persistence model
  -> challenge response read model
  -> verification/link write models where applicable
```

The accepted or controlled models are:

```text
Envelope
  -> accepted transport exception / deferred inbound-outbound split

Claims
  -> controlled auth runtime boundary

AuthorizationSubject
  -> controlled authorization runtime boundary

Session
  -> controlled WS runtime boundary

AppError
  -> accepted standardized error boundary

AppErrorSpec
  -> internal catalog metadata

authErrorSpec
  -> auth-local adapter metadata
```

---

## Status

Subphase: 0.12.1.5 — Target Separation Definition  
Status: Completed  
Code changes: none  
Next: 0.12.1.6 — Audit Consolidation & Closure
