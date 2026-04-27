# Phase 0.12.1.4 — Problem Detection & Risk Mapping

## Purpose

This document converts the cross-layer usage findings from Phase 0.12.1.3 into a concrete problem and risk map.

This subphase does not modify code, does not create read models, does not create write models, and does not introduce runtime behavior changes. Its only purpose is to identify the concrete architectural problems, risk level, impact area, and follow-up direction for the hybrid candidates previously identified.

---

## Source

* ZIP: `scavo-exchange-backend-0.12.1.3.zip`
* Baseline inventory: `docs/phase0_12_1_model_inventory.md`
* Baseline classification: `docs/phase0_12_1_model_classification.md`
* Baseline cross-layer analysis: `docs/phase0_12_1_cross_layer_usage_analysis.md`
* Hybrid candidates reviewed: 11
* Code changes: none

---

## Scope

### Included

* Convert cross-layer findings into concrete problem statements.
* Assign risk levels to each hybrid candidate.
* Identify impacted layers and modules.
* Separate high-priority extraction candidates from accepted runtime boundary models.
* Record follow-up direction for Phase 0.12.1.5.

### Excluded

* No code changes.
* No refactors.
* No package moves.
* No model extraction.
* No read/write model creation.
* No mapping implementation.
* No endpoint changes.
* No API contract changes.

---

## Risk Scale

| Risk | Meaning |
|---|---|
| CRITICAL | The model crosses domain, application, persistence, and transport boundaries in a way that can directly affect future read/write separation work. |
| HIGH | The model crosses important lifecycle boundaries and should be addressed explicitly before extraction work proceeds. |
| MEDIUM | The model crosses controlled runtime or infrastructure boundaries and likely needs documentation or a limited split. |
| LOW | The model is local or adapter-like and may be reclassified as INTERNAL during audit closure. |

---

## Summary

| Risk | Count |
|---|---:|
| CRITICAL | 3 |
| HIGH | 2 |
| MEDIUM | 4 |
| LOW | 2 |
| Total Reviewed | 11 |

---

## Problem Map

| Struct | Current Classification | Risk | Primary Problem | Follow-Up Direction |
|---|---|---|---|---|
| `User` | HYBRID | CRITICAL | Domain/persistence user model crosses user module, auth module, application flows, response construction, and tests. | Define target split between domain entity, provider-facing read contract, and auth/user response read models. |
| `UserSettings` | HYBRID | CRITICAL | Domain/persistence settings model crosses usersettings module, auth-facing contracts, read paths, update paths, router wiring, and tests. | Define target split between settings domain model, settings read model, and settings update/write model. |
| `WalletIdentity` | HYBRID | CRITICAL | Wallet identity model crosses domain, persistence, application flows, wallet HTTP flows, read model construction, and tests. | Keep canonical wallet identity internal to auth domain/persistence and map explicitly to wallet read models. |
| `WalletChallenge` | HYBRID | HIGH | Runtime/persistence challenge state is reused in application flows and public challenge responses. | Isolate challenge runtime/persistence state from public challenge response payloads. |
| `Envelope` | HYBRID | HIGH | A single WS protocol struct represents both inbound and outbound envelope semantics. | Decide whether to split inbound/outbound WS envelope shapes or document it as a transport-level exception. |
| `Claims` | HYBRID | MEDIUM | Token payload also acts as runtime identity context across HTTP, WS, auth, authorization, and application flows. | Keep as core auth runtime model if accepted, but prevent application read models from depending on raw token claims. |
| `AuthorizationSubject` | HYBRID | MEDIUM | Authorization runtime subject crosses policy, context, and HTTP middleware layers. | Treat as controlled authorization runtime model unless later usage expands into domain/read response flows. |
| `Session` | HYBRID | MEDIUM | WS runtime session state sits near similarly named auth session response/read concepts. | Preserve WS runtime session separately from auth session read models and avoid naming/lifecycle ambiguity. |
| `AppError` | HYBRID | MEDIUM | Runtime error object drives standardized HTTP response mapping. | Treat as accepted standardized error boundary from Phase 0.8 unless response leakage increases. |
| `AppErrorSpec` | HYBRID | LOW | Error catalog metadata is local to core error catalog definitions. | Reclassify candidate as INTERNAL/catalog metadata during audit closure if no wider usage is found. |
| `authErrorSpec` | HYBRID | LOW | Auth-local error adapter shape is contained to one adapter file. | Reclassify candidate as INTERNAL/adapter-local during audit closure if no wider usage is found. |

---

## Critical Risks

### 1. `User`

**Declared in:**

* `internal/modules/user/domain/model.go`
* alias in `internal/modules/user/model.go`

**Problem:**

`User` is currently the strongest cross-module hybrid. It represents a domain/persistence object inside the user module, but it is also consumed by auth flows and used as a source for session, profile, bootstrap, wallet, and test data.

**Impacted layers:**

* User domain
* User persistence
* User application
* Auth domain contracts
* Auth application
* Auth HTTP flows
* Tests

**Risk:** CRITICAL

**Why it matters:**

Future read/write separation can become unsafe if public response payloads, provider contracts, repository results, and domain entities remain tied to the same struct. Any domain change to `User` may unintentionally affect auth-facing read behavior or tests.

**Required follow-up:**

Phase 0.12.1.5 must define a target separation that distinguishes:

* user domain entity;
* auth-facing user provider contract shape;
* user/auth response read models;
* future write model candidates, if user mutations are introduced later.

---

### 2. `UserSettings`

**Declared in:**

* `internal/modules/usersettings/domain/model.go`
* alias in `internal/modules/usersettings/model.go`

**Problem:**

`UserSettings` crosses usersettings domain, persistence, application, auth-facing contracts, HTTP flows, router wiring, and tests. It participates in both read and update contexts.

**Impacted layers:**

* UserSettings domain
* UserSettings persistence
* UserSettings application
* Auth domain contracts
* Auth HTTP/app flows
* Router wiring
* Tests

**Risk:** CRITICAL

**Why it matters:**

Settings are naturally both read and write oriented. Reusing the same model for storage, read output, and mutation/update behavior increases the risk of accidental contract coupling.

**Required follow-up:**

Phase 0.12.1.5 must define a target split between:

* settings domain model;
* settings read model;
* settings update/write model or patch command;
* explicit mapping between domain and response payloads.

---

### 3. `WalletIdentity`

**Declared in:**

* `internal/modules/auth/domain/contracts.go`
* alias in `internal/modules/auth/wallet_identity_store.go`

**Problem:**

`WalletIdentity` is used as a canonical wallet model across domain contracts, persistence stores, auth application flows, wallet HTTP flows, read response construction, and tests.

**Impacted layers:**

* Auth domain
* Auth persistence stores
* Auth application
* Auth read models
* Auth HTTP flows
* Tests

**Risk:** CRITICAL

**Why it matters:**

Wallet identity state may contain internal persistence or operational fields that should not be tied directly to public read models. If response payloads are built directly from the domain/persistence model, later changes become higher risk.

**Required follow-up:**

Phase 0.12.1.5 must define whether `WalletIdentity` remains the internal canonical auth model and how it maps into dedicated wallet read models.

---

## High Risks

### 4. `WalletChallenge`

**Declared in:**

* `internal/modules/auth/domain/contracts.go`
* alias in `internal/modules/auth/wallet_challenge.go`

**Problem:**

`WalletChallenge` carries runtime and persistence challenge state, but is also used in application flows and response payload construction.

**Impacted layers:**

* Auth domain
* Auth challenge persistence
* Auth application
* Wallet HTTP flows
* Auth read models
* Router wiring
* Tests

**Risk:** HIGH

**Why it matters:**

Challenge internals and public challenge output do not necessarily have the same lifecycle. The same challenge struct should not define both persistent/runtime behavior and public response shape unless explicitly accepted.

**Required follow-up:**

Phase 0.12.1.5 must define a target read model for challenge responses and keep challenge runtime state internal.

---

### 5. `Envelope`

**Declared in:**

* `internal/core/ws/protocol.go`

**Problem:**

`Envelope` is used by the WS protocol, dispatcher, helpers, auth WS handlers, and system WS handlers. It represents both inbound and outbound protocol message semantics.

**Impacted layers:**

* WS protocol
* WS dispatcher
* WS helpers
* Auth WS handlers
* System WS handlers

**Risk:** HIGH

**Why it matters:**

A single transport envelope used for both inbound and outbound messages can remain acceptable as a protocol simplification, but it must be explicitly documented. Otherwise it can undermine the read/write separation rule by keeping a transport-level hybrid hidden in core.

**Required follow-up:**

Phase 0.12.1.5 must decide whether `Envelope` becomes:

* an accepted transport-level exception; or
* a candidate for future split into inbound/outbound envelope shapes.

---

## Medium Risks

### 6. `Claims`

**Declared in:**

* `internal/core/auth/jwt.go`

**Problem:**

`Claims` is a JWT token payload and also acts as identity context across HTTP, WS, auth, authorization, application flows, and tests.

**Risk:** MEDIUM

**Reasoning:**

This cross-layer usage appears intentional as a runtime authentication object. The risk is not immediate response coupling, but accidental dependency on token payload fields in application read models.

**Required follow-up:**

Keep `Claims` scoped as a core auth runtime model unless 0.12.1.5 finds read response or mutation behavior depending on raw token payload shape.

---

### 7. `AuthorizationSubject`

**Declared in:**

* `internal/core/authorization/subject.go`

**Problem:**

`AuthorizationSubject` crosses authorization policy, context hydration, HTTP middleware, and tests.

**Risk:** MEDIUM

**Reasoning:**

This is likely an acceptable authorization runtime boundary model. It should not become a domain read model or an output contract.

**Required follow-up:**

Record as controlled runtime model in 0.12.1.5 unless later evidence shows it leaking into response payloads.

---

### 8. `Session`

**Declared in:**

* `internal/core/ws/session.go`

**Problem:**

`Session` is WS runtime state that carries identity and claims, but it sits near auth session response concepts such as `SessionView` and `SessionResponse`.

**Risk:** MEDIUM

**Reasoning:**

The primary issue is lifecycle and naming ambiguity, not necessarily direct read/write coupling.

**Required follow-up:**

Preserve a clear distinction between WS runtime session state and auth session read models.

---

### 9. `AppError`

**Declared in:**

* `internal/core/errs/app_error.go`

**Problem:**

`AppError` is a runtime error object that drives standardized HTTP error response mapping.

**Risk:** MEDIUM

**Reasoning:**

This is an intentional architecture result from Phase 0.8. The model bridges internal behavior and output mapping, but does so in a controlled standardized error boundary.

**Required follow-up:**

Document as an accepted error boundary model unless 0.12.1.5 determines that a stricter split between internal errors and error read models is needed.

---

## Low Risks

### 10. `AppErrorSpec`

**Declared in:**

* `internal/core/errs/catalog_auth.go`

**Problem:**

`AppErrorSpec` is currently local to core error catalog metadata.

**Risk:** LOW

**Reasoning:**

It does not cross module boundaries directly and appears to be internal catalog configuration rather than a lifecycle hybrid model.

**Required follow-up:**

Reclassify as INTERNAL/catalog metadata during audit closure unless new usage is discovered.

---

### 11. `authErrorSpec`

**Declared in:**

* `internal/modules/auth/error_adapter.go`

**Problem:**

`authErrorSpec` is local to the auth error adapter.

**Risk:** LOW

**Reasoning:**

It appears adapter-local and not part of the broader read/write model lifecycle.

**Required follow-up:**

Reclassify as INTERNAL/adapter-local during audit closure unless new usage is discovered.

---

## Cross-Cutting Problems Detected

### Public Domain Aliases

Aliases such as:

* `internal/modules/user/model.go`
* `internal/modules/usersettings/model.go`
* `internal/modules/auth/wallet_identity_store.go`
* `internal/modules/auth/wallet_challenge.go`

preserve compatibility but also expose domain models outside their declaring domain packages.

**Risk:** HIGH

**Follow-up:**

Define whether aliases remain compatibility surfaces or should be phased out after read/write model separation.

---

### Domain Models Used as Response Sources

Several domain/persistence models are used to build read responses:

* `User`
* `UserSettings`
* `WalletIdentity`
* `WalletChallenge`

**Risk:** CRITICAL

**Follow-up:**

Introduce explicit mapping targets in later phases instead of relying on direct field reuse.

---

### Runtime Context Objects Crossing Transport Boundaries

Core runtime models cross HTTP and WS contexts:

* `Claims`
* `AuthorizationSubject`
* `Session`

**Risk:** MEDIUM

**Follow-up:**

Keep these as controlled runtime objects and prevent them from becoming response contracts.

---

### Transport-Level Hybrid Shape

`Envelope` is both inbound and outbound WS protocol shape.

**Risk:** HIGH

**Follow-up:**

Decide whether this is an accepted protocol simplification or a future inbound/outbound split candidate.

---

### Error Boundary Hybrid

`AppError` bridges internal behavior and standardized HTTP output mapping.

**Risk:** MEDIUM

**Follow-up:**

Keep as accepted standardized error boundary unless future changes require explicit error read models.

---

## Extraction Priority Candidates

The following models should drive the target separation definition in Phase 0.12.1.5.

| Priority | Struct | Risk | Target Direction |
|---:|---|---|---|
| 1 | `User` | CRITICAL | Domain entity + provider read contract + response read models. |
| 2 | `UserSettings` | CRITICAL | Domain model + read model + update/write model. |
| 3 | `WalletIdentity` | CRITICAL | Internal wallet identity + explicit wallet read models. |
| 4 | `WalletChallenge` | HIGH | Internal challenge state + challenge response read model. |
| 5 | `Envelope` | HIGH | Transport exception or inbound/outbound protocol split. |

---

## Accepted / Controlled Boundary Candidates

These models should likely remain as controlled runtime or infrastructure boundary objects, pending final target definition.

| Struct | Suggested Treatment |
|---|---|
| `Claims` | Keep as core auth runtime identity/token payload model. |
| `AuthorizationSubject` | Keep as core authorization runtime subject. |
| `Session` | Keep as WS runtime session state, separate from auth read models. |
| `AppError` | Keep as standardized error boundary model from Phase 0.8. |
| `AppErrorSpec` | Reclassify as internal catalog metadata. |
| `authErrorSpec` | Reclassify as auth-local adapter metadata. |

---

## Risk Notes for Later Implementation Phases

Later implementation phases must be especially careful with tests because several test suites depend on current hybrid shapes.

Potentially sensitive areas include:

* auth HTTP handler tests;
* wallet challenge/link/verify tests;
* user service tests;
* user settings service tests;
* router/versioning tests;
* authorization enforcement tests;
* WS handler/session tests.

No test changes are introduced in this subphase.

---

## Non-Goals Confirmed

This subphase intentionally did not:

* modify Go code;
* create read models;
* create write models;
* move packages;
* change repository contracts;
* change provider contracts;
* change HTTP routes;
* change WS protocols;
* change error responses;
* change tests;
* decide final extraction implementation.

---

## Status

Phase: 0.12 — Read / Write Model Separation  
Subphase: 0.12.1 — Model Classification Audit  
Step: 0.12.1.4 — Problem Detection & Risk Mapping  
Status: COMPLETED  
Code Impact: NONE  
Next Step: 0.12.1.5 — Target Separation Definition
