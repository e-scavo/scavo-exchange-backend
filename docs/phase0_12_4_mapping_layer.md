# Phase 0.12.4 — Mapping Layer Introduction

## Subphase 0.12.4.0 — Definition & Documentation Lock

---

## Objective

Define and lock the Mapping Layer Introduction sequence before any code-level mapper consolidation starts.

Phase 0.12.4 centralizes explicit transformation ownership after the completed read-side and write-side model separation work.

The target is to move the system from distributed mapper ownership toward a dedicated module-local mapping boundary.

---

## Context

Phase 0.12.2 introduced:

* explicit read model packages
* Domain/Application → Read mapping
* response alignment while preserving HTTP contracts

Phase 0.12.3 introduced:

* explicit write model packages
* domain write input structures
* Write → Domain mapping
* handler alignment while preserving request payload semantics

The system now has both mapping directions, but ownership is still distributed across model packages and application support code.

---

## Problem Statement

The current post-0.12.3 state is structurally valid but not yet fully consolidated.

Mapping responsibility exists in multiple places:

* read model packages
* write model packages
* application support helpers
* handler-facing response assembly

This creates risks:

* duplicated transformation logic
* unclear mapper ownership
* future drift between read and write boundaries
* accidental reintroduction of hybrid behavior inside application code

---

## Target Architecture

Introduce a dedicated module-local mapping layer:

```text
internal/modules/<module>/mappers/
```

Examples:

```text
internal/modules/auth/mappers/
internal/modules/user/mappers/
internal/modules/usersettings/mappers/
```

The mapper package owns transformation logic.

Model packages own model definitions only.

Application packages orchestrate use cases and should not accumulate transformation details.

---

## Mapping Directions

The mapping layer standardizes the explicit directions introduced in previous subphases:

```text
Write → Domain
Domain/Application → Read
```

Optional internal transformations may be introduced only when they preserve existing contracts and remain module-local.

---

## Scope

### Included

* define centralized mapper ownership
* define module-local mapper package strategy
* prepare consolidation of existing read and write mappers
* preserve current runtime behavior
* preserve existing public HTTP request and response contracts

### Excluded

* public API changes
* route changes
* authentication behavior changes
* authorization behavior changes
* standardized error model changes
* CQRS
* event sourcing
* multi-tenant redesign

---

## Sub-Subphase Plan

| Sub-subphase | Name | Status |
| --- | --- | --- |
| 0.12.4.0 | Definition & Documentation Lock | Completed |
| 0.12.4.1 | Mapping Layer Design | Pending |
| 0.12.4.2 | Mapping Layer Implementation | Pending |
| 0.12.4.3 | Mapping Consolidation | Pending |
| 0.12.4.4 | Application Refactor | Pending |
| 0.12.4.5 | Validation & Compatibility | Pending |
| 0.12.4.6 | Documentation & Closure | Pending |

---

## Compatibility Rules

During 0.12.4:

* existing HTTP routes must remain unchanged
* existing JSON request payload semantics must remain unchanged
* existing JSON response payload semantics must remain unchanged
* read models must remain output-only
* write models must remain input-only
* domain models must remain domain-owned
* mapper consolidation must be incremental
* every code-level subphase must validate with `go test ./...`

---

## Implementation Constraints

The implementation must not remove current mapper behavior before equivalent centralized mapper behavior exists.

The implementation must avoid import cycles.

The implementation must preserve module boundaries established in Phase 0.11.

The implementation must preserve read/write model separation established in Phase 0.12.2 and Phase 0.12.3.

---

## Expected Outcome

After 0.12.4 is completed, the backend should have:

* a clear module-local mapping layer
* reduced transformation logic inside application support code
* read/write model packages focused on model definitions
* explicit mapper ownership
* unchanged public HTTP contracts
* unchanged API versioning
* validated runtime compatibility

---

## Status

Phase: 0.12.4  
Subphase: 0.12.4.0  
Status: Completed  
Code impact: None  
Next: 0.12.4.1 — Mapping Layer Design
