# Phase 0.15 — Contract Hardening & Freeze

---

## Status

**Phase:** 0.15 — Contract Hardening & Freeze  
**Current Subphase:** 0.15.6 — Validation & Documentation  
**Status:** 0.15.5 Completed / 0.15.6 Pending  
**Type:** Contract hardening and freeze documentation  
**Code changes in 0.15.1:** No

---

## Context Inherited

Phase 0.15 starts from the validated Phase 0.14.6.fix3 baseline.

The inherited Stage 0 sequence is important:

- Phase 0.12 separated read and write model responsibilities.
- Phase 0.13 consolidated provider ownership and runtime composition.
- Phase 0.14 introduced request correlation, structured logging, safe diagnostic error context, minimal flow tracing and the `/diagnostics` surface.

The backend is now functionally stable, structurally clearer and observable enough to inspect runtime movement across the current HTTP → Provider → Application → Domain → Repository path.

---

## Real Problem

The remaining Stage 0 risk is contractual drift.

The system has working contracts, but several of them are still implicit or documented through behavior rather than through a dedicated hardening and freeze policy. That creates risk for later frontend integration, future backend extension and production stabilization.

The relevant risks are:

- HTTP endpoints exist and are routed, but their status and response expectations still need a dedicated contract audit.
- Error responses have a standardized model from Phase 0.8, but Phase 0.15 must verify that the shape remains aligned across the current route surface.
- Provider contracts were consolidated in Phase 0.13, but their current inputs and outputs still need explicit validation against the application-facing responsibilities they now own.
- Response schemas may contain historical variation that should either be normalized or intentionally preserved and documented.
- There is no explicit freeze policy defining what may change, what must be versioned and what must remain stable after the contract baseline is accepted.

---

## Decision Taken

Phase 0.15 introduces **Contract Hardening & Freeze** as the next Stage 0 phase after Observability & Diagnostics Foundation.

This decision keeps the phase focused on stabilizing the existing system rather than expanding it.

Phase 0.15 does not introduce new business behavior, new architecture or new observability scope. It formalizes, audits and freezes the contracts already present in the repository.

---

## Objective

Formalize, validate and freeze the current backend contracts so future evolution can happen through controlled, explicit and version-aware changes.

The target contract families are:

- public HTTP route contracts
- public response payload contracts
- public error envelope contracts
- internal provider contracts
- internal application-facing contract expectations
- contract freeze and evolution rules

---

## Scope

### Included

- HTTP contract audit
- response schema verification
- error contract alignment
- provider contract validation
- schema normalization where needed
- explicit freeze policy
- documentation reconciliation
- validation with `go test ./...` at phase closure

### Excluded

- new features
- new public routes
- business logic changes
- architecture changes
- new metrics platforms
- dashboarding
- OpenTelemetry or Prometheus work
- observability expansion beyond the Phase 0.14 baseline

---

## Phase Narrative

The Stage 0 narrative now becomes:

```text
Phase 0.12: clarify model direction and mapping ownership
Phase 0.13: clarify orchestration ownership and runtime composition
Phase 0.14: clarify runtime visibility and diagnostic traceability
Phase 0.15: clarify, harden and freeze system contracts
```

This is the correct next step because observability makes the runtime path inspectable, but inspectability alone does not prevent accidental contract drift. Phase 0.15 uses the stabilized and observable baseline to audit what exists and define what must remain stable.

---

## Architectural Impact

The architectural path remains unchanged:

```text
HTTP → Provider → Application → Domain → Repository
```

Phase 0.15 adds contractual discipline to that path:

- HTTP remains the public transport boundary.
- Providers remain the runtime composition boundary.
- Application services remain use-case owners.
- Domain and repository behavior remain below the application boundary.
- Response and error contracts become explicit stability targets.

No layer is moved or re-owned by this phase.

---

## Subphases

### 0.15.0 — Phase Definition & Documentation Lock

Register Phase 0.15 in the trunk documentation, define the phase scope and lock the documentation baseline before any contract implementation or audit work begins.

Rules:

- documentation only
- no Go code changes
- no contract changes
- no endpoint changes
- no behavior changes

### 0.15.1 — HTTP Contract Audit

Audit the currently registered HTTP endpoints and document their expected method, path, status behavior and payload families.

The purpose is detection and documentation first. Any correction must be justified by the observed contract and must not introduce new behavior.

### 0.15.2 — Error Contract Alignment

Validate that public error responses remain aligned with the existing standardized error model:

- `code`
- `message`
- `details`

The goal is to remove or document uncontrolled variation while preserving the existing public contract.

### 0.15.3 — Provider Contract Validation

Validate the internal contracts between HTTP handlers, providers and application services.

The goal is to confirm that provider boundaries introduced by Phase 0.13 are explicit, stable and not bypassed accidentally.

### 0.15.4 — Response Schema Normalization

Normalize response schemas only where the repository shows contract drift or unclear historical variation.

The goal is consistent response shape without changing business meaning.

### 0.15.5 — Contract Freeze Enforcement

Define the freeze policy for current contracts.

The policy must clarify:

- which contracts are frozen
- what requires versioning
- what can be extended compatibly
- what cannot be changed silently
- how future contract drift is detected

### 0.15.6 — Validation & Documentation

Close Phase 0.15 through validation and documentation reconciliation.

Expected closure actions:

- run `go test ./...`
- reconcile trunk documentation
- record exact contract baseline
- close the phase narratively
- confirm no uncontrolled contract drift remains

---

## 0.15.0 Concrete Change

This subphase introduces the Phase 0.15 documentation baseline only.

The concrete documentation change is:

- add this dedicated Phase 0.15 document
- update the current stage metadata from “awaiting Phase 0.15 definition” to active Phase 0.15.0
- register the full 0.15 subphase plan in the roadmap and documentation index
- update handoff state so the next step is 0.15.1 HTTP Contract Audit

No source code is changed by 0.15.0.

---

## Observable Impact

After 0.15.0:

- the repository has a formal Phase 0.15 definition
- Phase 0.15 is no longer pending definition
- the active subphase is 0.15.0
- the next implementation step is clearly 0.15.1
- contract work is constrained before any code change occurs

---

## Freeze Principle

Once Phase 0.15 closes, contracts may evolve only through explicit, documented and validated change.

Silent drift is not allowed.

---

## Current Phase State

- 0.15.0 — Phase Definition & Documentation Lock: **Completed**
- 0.15.1 — HTTP Contract Audit: **Completed**
- 0.15.2 — Error Contract Alignment: **Completed**
- 0.15.3 — Provider Contract Validation: **Completed**
- 0.15.4 — Response Schema Normalization: **Completed**
- 0.15.5 — Contract Freeze Enforcement: **Completed**
- 0.15.6 — Validation & Documentation: **Pending**

---

## 0.15.1 Concrete Change

0.15.1 audits the real HTTP route surface from the repository and records the route, method, authentication, success status, response family and error-envelope baseline.

The audit confirms:

- `GET /health`, `GET /readiness`, `GET /version`, `GET /diagnostics` and `GET /ws` are the foundation/runtime routes.
- auth routes are registered both as legacy paths and canonical `/api/v1` paths.
- both auth route surfaces point to the same handler set.
- authenticated routes use `RequireAuth(tokens, false)` and authorization hydration.
- selected user and settings routes also enforce explicit read/update permissions.
- handler errors use the standardized `{error:{code,message,details}}` envelope.

No source code is changed by 0.15.1.

---

## 0.15.1 Observable Impact

After 0.15.1:

- the HTTP route surface is no longer implicit
- legacy and canonical route pairing is documented
- success status families are recorded
- response families are identified
- non-blocking risks are preserved for 0.15.2 through 0.15.5
- the next correct step is Error Contract Alignment

---

## Handoff to 0.15.2

The public error envelope has now been validated and aligned in 0.15.2. The next subphase must validate provider contracts against the current route and error-contract baselines.

It must use `docs/phase0_15_1_http_contract_audit.md` and `docs/phase0_15_2_error_contract_alignment.md` as baselines and must not invent provider responsibilities or response shapes.

---


## 0.15.2 Concrete Change

0.15.2 aligns the public error envelope with the documented `{error:{code,message,details}}` contract.

The subphase confirms that the remaining divergence was not an error-code or handler-decision issue. The divergence was structural: `details` could be omitted for detail-free errors.

The implementation now requires `error.details` to serialize as a JSON object. Empty public details serialize as `{}`.

Changed code:

- `internal/core/errs/response_error.go`
- `internal/core/errs/app_error_test.go`
- `internal/core/httpx/error_test.go`
- `internal/core/httpx/router_versioning_test.go`
- `internal/modules/auth/http_handlers_test.go`

No route, status code, error code, message, provider contract, domain behavior or success payload changed.

---

## 0.15.2 Observable Impact

After 0.15.2:

- canonical HTTP error responses always expose `error.details`
- clients can treat `details` as a stable object
- diagnostic details remain safe and copied
- the error envelope is ready for provider validation, response normalization and freeze enforcement

---

## Handoff to 0.15.3

The next subphase is Provider Contract Validation.

It must verify provider inputs, outputs and error propagation responsibilities without changing the public error envelope aligned in 0.15.2.

## End of Document

## 0.15.3 Concrete Change

0.15.3 validates provider contracts through compile-time assertions at the auth boundary.

Context inherited: 0.15.1 defined the HTTP route baseline and 0.15.2 aligned the public error envelope. Provider validation therefore did not need to change external behavior; it needed to protect the internal seams that feed those external contracts.

Real problem: provider interfaces already existed and were typed, but their satisfaction was mostly implicit. The handler-facing auth provider contract and the cross-module user/usersettings contracts could drift through future edits unless the build explicitly rejected mismatches.

Decision taken: keep runtime behavior unchanged and add compile-time assertions to the existing auth boundary.

Concrete change: `internal/modules/auth/application.go` now asserts that `*Application` implements `AuthProvider`, `*user.Service` implements `authdomain.UserProvider`, and `*usersettings.Service` implements `authdomain.UserSettingsProvider`.

Impact observable: no HTTP response changed, but future provider drift becomes a compile-time failure. This preserves the contract freeze path without introducing new behavior.

## Handoff to 0.15.4

0.15.4 must continue with Response Schema Normalization.

It must use the following validated baselines:

- 0.15.1 HTTP route audit
- 0.15.2 canonical error envelope alignment
- 0.15.3 provider boundary compile-time validation

It must not introduce new features, new routes, new business rules or unrelated architecture changes.

## 0.15.4 Concrete Change

0.15.4 normalizes response serialization policy without changing external behavior.

Context inherited: 0.15.1 established the real route surface, 0.15.2 stabilized the public error envelope and 0.15.3 validated provider seams.

Real problem: payload structures were compatible, but response serialization still had two narrow divergences. Auth error responses used a different JSON content type than the core writer, and the defensive timeout fallback JSON did not include the mandatory `details` object.

Decision taken: preserve all public payloads and normalize only serialization details.

Concrete change: `internal/modules/auth/http_login.go` now writes auth error JSON as `application/json; charset=utf-8`; `internal/core/httpx/middleware.go` keeps the defensive fallback aligned with `{error:{code,message,details}}`; `internal/modules/auth/http_handlers_test.go` validates the auth error content type.

Impact observable: existing clients keep the same payload shapes, status codes and error codes, while response serialization is now consistent enough to support freeze enforcement.

## Handoff to 0.15.5

0.15.5 must define Contract Freeze Enforcement.

It must use the validated baselines from 0.15.1, 0.15.2, 0.15.3 and 0.15.4. It must document and enforce how contracts can evolve without allowing silent drift.

## 0.15.5 Concrete Change

0.15.5 defines and enforces the Stage 0 contract freeze.

Context inherited: 0.15.1 audited the HTTP route surface, 0.15.2 aligned the canonical error envelope, 0.15.3 validated provider contracts and 0.15.4 normalized response serialization.

Real problem: the backend had explicit contract baselines but no final policy that explained what is frozen, what may evolve only through versioning and what must fail fast during future drift.

Decision taken: freeze the current contract surface without changing runtime behavior.

Concrete change: the freeze policy is documented in `docs/phase0_15_5_contract_freeze_enforcement.md`, and representative frozen contracts are guarded by `internal/core/httpx/contract_freeze_test.go`.

Impact observable: public JSON endpoints, protected auth error envelopes, JSON metadata and provider assertions now have a clear freeze baseline. Future contract evolution must be deliberate and documented.

## Handoff to 0.15.6

0.15.6 must validate the full system and reconcile all trunk documentation for Phase 0.15 closure.

It must include the local `go test ./...` evidence supplied after applying 0.15.5 and close Phase 0.15 without introducing new behavior.
