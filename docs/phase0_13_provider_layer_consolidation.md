# Phase 0.13 — Provider Layer Consolidation

## Status

Phase 0.13 is **IN PROGRESS**.

Phase 0.13.0 — Definition & Documentation Lock is **COMPLETED** as the documentation-only entry point for the phase.

No Go code is changed by 0.13.0.

---

## Objective

Consolidate the Provider Layer as the explicit and consistent entry point to domain services.

The phase prepares the backend for clearer orchestration boundaries by ensuring handlers and application-level flows depend on stable provider contracts rather than scattered direct access to lower-level domain or repository logic.

---

## Why This Phase Exists

Phase 0.12 completed Read / Write Model Separation, centralized mapping ownership, and aligned internal contracts.

The remaining structural risk is that provider usage and provider responsibilities may still be distributed across handlers, application services, module services, or compatibility wiring.

Provider Layer Consolidation exists to make the boundary explicit before later Stage 0 hardening and feature expansion.

---

## Included Scope

- inventory of current provider-like responsibilities
- classification of existing providers, missing providers and invalid direct access patterns
- provider interface design per module where required
- implementation of provider boundaries without public API drift
- application and handler integration with consolidated providers
- validation of compatibility through existing test commands
- accumulated documentation updates

---

## Excluded Scope

- public HTTP route changes
- public request or response payload changes
- API versioning changes
- business behavior changes
- CQRS or event sourcing
- observability implementation
- Stage 1 feature expansion

---

## Target Architecture Direction

```text
HTTP → Provider → Application → Domain → Repository
```

The direction is structural. It does not imply that every module must receive a new provider immediately if the current boundary is already explicit and compatible.

---

## Subphase Plan

| Subphase | Name | Status |
| --- | --- | --- |
| 0.13.0 | Definition & Documentation Lock | COMPLETED |
| 0.13.1 | Provider Inventory & Classification | PENDING |
| 0.13.2 | Provider Interface Design | PENDING |
| 0.13.3 | Provider Implementation | PENDING |
| 0.13.4 | Application Integration | PENDING |
| 0.13.5 | Validation & Compatibility | PENDING |
| 0.13.6 | Documentation & Closure | PENDING |

---

## 0.13.0 — Definition & Documentation Lock

### Purpose

Document and lock the Phase 0.13 scope before any implementation begins.

### Required Actions

- verify all general Markdown documentation for phase consistency
- correct the roadmap-level Phase 0.13 definition
- register the Provider Layer Consolidation phase and subphase plan
- preserve Phase 0.12 closure state
- avoid Go code changes

### Compatibility Rule

0.13.0 did not alter runtime behavior, public contracts, routes, middleware, repositories, handlers or application services.

### Result

The phase definition is locked, the subphase plan is expanded consistently and the next step is 0.13.1 — Provider Inventory & Classification.

---

## 0.13.1 — Provider Inventory & Classification

Identify current provider boundaries and provider-like responsibilities across modules.

Classify findings as:

- existing provider
- provider candidate
- missing provider
- invalid direct access
- compatibility wiring that must remain stable

---

## 0.13.2 — Provider Interface Design

Define explicit provider interfaces where required.

Interfaces must remain narrow, module-owned and aligned with the read/write model separation introduced during Phase 0.12.

---

## 0.13.3 — Provider Implementation

Implement provider boundaries incrementally.

Implementation must preserve public behavior and avoid unrelated domain or API changes.

---

## 0.13.4 — Application Integration

Align handlers and application services with consolidated provider boundaries.

Direct access patterns identified during inventory should be replaced only when the provider contract is defined and compatible.

---

## 0.13.5 — Validation & Compatibility

Validate runtime compatibility with existing test commands.

Expected validation command:

```bash
go test ./...
```

Any additional validation must be documented with the specific subphase that introduced the need.

---

## 0.13.6 — Documentation & Closure

Close the phase by updating all required documentation, recording the final provider boundary state, and confirming compatibility preservation.

---

## Documentation Rule For This Phase

All general `.md` files are treated as trunk documentation and must be reviewed for phase consistency.

Phase-specific documents are reviewed and only changed when they contain direct references, next-phase expectations or inconsistencies that affect the current phase.
