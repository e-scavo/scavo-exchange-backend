# Phase 0.12 — Read / Write Model Separation

## Status

In progress.

Latest completed subphase: **0.12.0 — Phase Definition & Documentation Lock**.

---

## Objective

Separate explicitly the models used for reading data from the models used for writing data, while preserving the existing public HTTP/API behavior and the module boundaries completed in Phase 0.11.

Phase 0.12 is a structural Stage 0 phase. It improves internal clarity without introducing new business behavior.

---

## Background

Phase 0.11 completed the Domain Module Pattern across the current Stage 0 domain-facing modules:

- `auth`
- `user`
- `usersettings`

That work clarified where responsibilities live. Phase 0.12 now clarifies what responsibility each model has inside and across those boundaries.

The current repository already contains response-oriented structures, request-oriented structures, domain structures, application result structures and cross-module contracts. Some are already explicit. Others still require audit before extraction.

---

## Scope

### Included

- identify all current model structures in the real repository
- classify models as READ, WRITE, DOMAIN, CONTRACT, INFRASTRUCTURE or HYBRID / TRANSITIONAL
- extract explicit read models where required
- isolate write models where required
- introduce explicit mapping functions
- align internal contracts introduced or consolidated in 0.11
- preserve current public route and payload compatibility
- update trunk documentation cumulatively

### Excluded

- full CQRS
- event sourcing
- public API version changes
- business behavior changes
- multi-tenant redesign
- independent read/write persistence stores
- frontend migration requirements
- observability expansion

---

## Architectural Direction

The intended internal flow is:

```text
HTTP request DTO -> Write Model -> Application / Domain behavior
Application / Domain result -> Read Model -> HTTP response DTO
```

This phase does not require every package to use identical filenames or folders. It requires model responsibility to be explicit and transformations to be visible.

---

## Model Categories

### READ

A read model is used for output, views, responses or bootstrap/session/profile data returned by application use cases or HTTP handlers.

### WRITE

A write model is used for input, command-style mutation requests or application-level write operations.

### DOMAIN

A domain model is owned by a module and represents module semantics, invariants or persistent identity concepts.

### CONTRACT

A contract is an explicit interface or minimal cross-module shape used to coordinate between modules without taking ownership of another module's implementation.

### INFRASTRUCTURE

An infrastructure model supports transport, storage, tokens, configuration, logging, errors or framework integration.

### HYBRID / TRANSITIONAL

A hybrid/transitional model currently carries more than one responsibility and must be considered for later extraction, isolation or mapping.

---

## Subphases

### 0.12.0 — Phase Definition & Documentation Lock

Status: **Completed**.

Purpose:

- define the phase before code changes
- document scope and non-goals
- register all subphases
- update trunk documentation
- preserve the completed 0.11 runtime state

No code changes are included in this subphase.

### 0.12.1 — Model Classification Audit

Status: **Completed**.

Purpose:

- inspect real repository code
- identify all relevant models
- classify each model
- detect hybrid/transitional structures
- document required extraction candidates

Result:

- complete inventory extracted
- 125 model-like structs recorded
- classification completed
- 11 hybrid/transitional structures identified
- cross-layer usage analyzed
- risk map produced
- target separation direction defined
- audit closure completed without Go code changes

### 0.12.2 — Read Model Extraction

Status: **Completed**.

Internal subdivision:

- 0.12.2.0 — Definition & Documentation Lock ✔ Completed
- 0.12.2.1 — Read Model Design ✔ Completed
- 0.12.2.2 — Read Model Implementation ✔ Completed
- 0.12.2.3 — Domain/Application → Read Mapping ✔ Completed
- 0.12.2.4 — Response Alignment ✔ Completed
- 0.12.2.5 — Validation & Compatibility ✔ Completed
- 0.12.2.6 — Documentation & Closure ✔ Completed

Purpose:

- create explicit read models where required
- preserve response compatibility
- introduce Domain / Application → Read Model mappings where needed
- validate with `go test ./...`

### 0.12.3 — Write Model Isolation

Status: **Completed**.

Internal subdivision:

- 0.12.3.0 — Definition & Documentation Lock ✔ Completed
- 0.12.3.1 — Write Model Design ✔ Completed
- 0.12.3.2 — Write Model Implementation ✔ Completed
- 0.12.3.3 — Write → Domain Mapping ✔ Completed
- 0.12.3.4 — Handler Alignment ✔ Completed
- 0.12.3.5 — Validation & Compatibility ✔ Completed
- 0.12.3.6 — Documentation & Closure ✔ Completed

Purpose:

- create explicit write models where required
- prevent response/read structures from being reused as mutation input models
- preserve handler compatibility
- validate with `go test ./...`

### 0.12.4 — Mapping Layer Introduction

Status: **In progress**.

Purpose:

- introduce explicit centralized mapping ownership
- consolidate Write → Domain transformations
- consolidate Domain/Application → Read transformations
- remove implicit conversions from application code
- reduce coupling between models and handlers

Internal sequence:

- 0.12.4.0 — Definition & Documentation Lock ✔ Completed
- 0.12.4.1 — Mapping Layer Design ✔ Completed
- 0.12.4.2 — Mapping Layer Implementation ✔ Completed
- 0.12.4.3 — Mapping Consolidation ✔ Completed
- 0.12.4.4 — Application Refactor ✔ Completed
- 0.12.4.5 — Validation & Compatibility ✔ Completed
- 0.12.4.6 — Documentation & Closure ✔ Completed

### 0.12.5 — Contract Alignment

Status: **Completed**.

Purpose:

- revisit 0.11 cross-module contracts
- align `UserProvider` and `UserSettingsProvider` usage with read/write separation
- align provider-facing inputs and outputs with centralized mapping ownership
- preserve runtime compatibility

Internal sequence:

- 0.12.5.0 — Definition & Documentation Lock ✔ Completed
- 0.12.5.1 — Contract Inventory & Classification ✔ Completed
- 0.12.5.2 — Contract Normalization Design ✔ Completed
- 0.12.5.3 — Contract Alignment Implementation ✔ Completed
- 0.12.5.4 — Handler Contract Adjustment ✔ Completed
- 0.12.5.5 — Validation & Compatibility ✔ Completed
- 0.12.5.6 — Documentation & Closure ✔ Completed

## Compatibility Rules

Phase 0.12 must preserve:

- existing public routes
- existing public payload semantics
- existing authentication behavior
- existing authorization behavior
- existing standardized error model
- existing API versioning policy
- existing frontend alignment rule during Stage 0

---

## 0.12.0 Closure Statement

0.12.0 is completed as a documentation-only subphase. The phase is now defined, its scope is locked and Phase 0.12.1 has started with 0.12.1.0 as its audit-level documentation lock. The next correct step is 0.12.1.1 — Model Inventory Extraction.

## 0.12.2.0 Closure Statement

0.12.2.0 is completed as a documentation-only subphase. Read Model Extraction is now defined, its internal sequence is locked and implementation must start with 0.12.2.1 — Read Model Design. No Go code is changed in 0.12.2.0.


## 0.12.2 Completion Summary

0.12.2 — Read Model Extraction is completed. The repository now includes explicit read model packages, explicit Domain/Application → Read mapping functions and response alignment that preserves public HTTP compatibility.

Validation was recorded in 0.12.2.5 using `go test ./...`. The next phase step is 0.12.3 — Write Model Isolation.

## 0.12.3.0 Closure Statement

0.12.3.0 is completed as a documentation-only subphase. Write Model Isolation is now defined, its internal sequence is locked and the subsequent implementation subphases have been completed.

## 0.12.3.6 Closure Statement

0.12.3 is completed. The repository now contains explicit write model packages, domain write input structures, Write → Domain mappers, handler alignment preserving public request payload semantics and validation documentation backed by `go test ./...`. The next planned step is 0.12.4 — Mapping Layer Introduction.

## 0.12.4.0 Closure Statement

0.12.4 is completed. Mapping Layer Introduction now has module-local mapper packages, consolidated mapping ownership, application-layer residual mapping reduction and validation documentation. The next phase is 0.12.5 — Contract Alignment.

## Phase 0.12.5.0 Contract Alignment Definition Lock

0.12.5 is completed. Contract inventory, contract normalization design, contract alignment implementation, handler contract adjustment, validation and documentation closure are now recorded. Public HTTP/API behavior remains unchanged and validation is backed by `go test ./...`. Phase 0.12 is complete.
