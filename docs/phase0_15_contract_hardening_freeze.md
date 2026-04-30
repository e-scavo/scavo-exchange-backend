# Phase 0.15 — Contract Hardening & Freeze

---

## Status

**Phase:** 0.15 — Contract Hardening & Freeze  
**Current Subphase:** 0.15.0 — Phase Definition & Documentation Lock  
**Status:** In Progress  
**Type:** Documentation-only phase definition step  
**Code changes in 0.15.0:** No

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

- 0.15.0 — Phase Definition & Documentation Lock: **In Progress**
- 0.15.1 — HTTP Contract Audit: **Pending**
- 0.15.2 — Error Contract Alignment: **Pending**
- 0.15.3 — Provider Contract Validation: **Pending**
- 0.15.4 — Response Schema Normalization: **Pending**
- 0.15.5 — Contract Freeze Enforcement: **Pending**
- 0.15.6 — Validation & Documentation: **Pending**

---

## Handoff to 0.15.1

The next subphase must audit the real HTTP route surface from the repository.

It must not invent endpoints, routes, status codes or response shapes. Any contract inventory must be derived from the current code and validated against the existing documentation.

---

## End of Document
