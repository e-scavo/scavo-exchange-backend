# Phase 0.12.5 — Contract Alignment

## Subphase 0.12.5.0 — Definition & Documentation Lock

---

## Objective

Define and lock the complete Contract Alignment sequence before changing provider contracts or runtime wiring.

Phase 0.12.5 exists to align the internal contracts stabilized during Phase 0.11 with the read/write model separation introduced during Phase 0.12.2, Phase 0.12.3 and Phase 0.12.4.

This phase must preserve public HTTP/API behavior while improving internal contract intent.

---

## Context

The backend has completed the following structural work:

- Phase 0.12.1 — model inventory, classification, cross-layer usage analysis, risk mapping and target separation
- Phase 0.12.2 — explicit read model packages and Domain/Application → Read mapping
- Phase 0.12.3 — explicit write model packages and Write → Domain mapping
- Phase 0.12.4 — module-local centralized mapper packages and reduced application-layer mapping residuals

The remaining contract risk is that internal provider contracts may still expose mixed responsibility shapes or hide transformation intent behind application services.

---

## Scope

### Included

- inventory of internal provider contracts
- review of `UserProvider` and `UserSettingsProvider` usage
- alignment of contract inputs with write-side intent
- alignment of contract outputs with read-side intent
- preservation of mapper ownership under module-local `mappers` packages
- compatibility validation

### Excluded

- public HTTP route changes
- public request payload changes
- public response payload changes
- API versioning changes
- business logic changes
- repository storage changes
- CQRS or event sourcing

---
## Completed Subphase Sequence

Phase 0.12.5 is closed with the following completed sequence:

- 0.12.5.0 — Definition & Documentation Lock (COMPLETED)
- 0.12.5.1 — Contract Inventory & Classification (COMPLETED)
- 0.12.5.2 — Contract Normalization Design (COMPLETED)
- 0.12.5.3 — Contract Alignment Implementation (COMPLETED)
- 0.12.5.4 — Handler Contract Adjustment (COMPLETED)
- 0.12.5.5 — Validation & Compatibility (COMPLETED)
- 0.12.5.6 — Documentation & Closure (COMPLETED)

---


## Contract Alignment Principles

Provider contracts must:

- express clear input/output intent
- avoid implicit read/write ambiguity
- avoid leaking transport structs into domain-facing boundaries
- avoid using read models as write inputs
- avoid using write models as response outputs
- keep mapper ownership outside application services
- preserve runtime behavior

---

## Mapping Ownership Rule

The mapping directions remain:

```text
Write → Domain
Domain/Application → Read
```

Ownership remains under:

```text
internal/modules/<module>/mappers/
```

Provider contracts should not reintroduce mapping ownership into:

- `app/`
- `readmodels/`
- `writemodels/`
- HTTP handlers

---

## Internal Sub-Subphase Plan

### 0.12.5.0 — Definition & Documentation Lock

Status: Completed.

Purpose:

- define the Contract Alignment sequence
- update trunk documentation
- preserve compatibility rules before implementation
- avoid code changes

### 0.12.5.1 — Contract Inventory & Classification

Status: Completed.

Purpose:

- inspect all provider contracts
- identify contract inputs and outputs
- identify read/write/domain usage
- detect mixed contract responsibilities

### 0.12.5.2 — Contract Normalization Design

Status: Completed.

Purpose:

- define target contract shapes
- determine which contracts remain unchanged
- determine which contracts need explicit read/write/domain alignment
- define safe migration order

### 0.12.5.3 — Contract Alignment Implementation

Status: Completed.

Purpose:

- apply provider contract alignment incrementally
- preserve runtime behavior
- preserve public HTTP/API contracts
- keep mapping ownership centralized

### 0.12.5.4 — Handler Contract Adjustment

Status: Completed.

Purpose:

- compare aligned contracts against previous behavior
- verify no public contract drift
- verify no route or versioning change
- verify no domain behavior change

### 0.12.5.5 — Validation & Compatibility

Status: Completed.

Purpose:

- record validation results
- require `go test ./...`
- document compatibility outcome

### 0.12.5.6 — Documentation & Closure

Status: Completed.

Purpose:

- update trunk documentation cumulatively
- close Phase 0.12.5
- prepare Phase 0.12.6

---

## Expected Outcome

After Phase 0.12.5 is completed:

- provider contracts have explicit read/write/domain intent
- contract usage is aligned with centralized mapper ownership
- application services no longer obscure contract transformation semantics
- public HTTP/API behavior remains unchanged
- the backend is ready for Phase 0.12.6 documentation and phase closure

---

## Status

Phase: 0.12.5  
Subphase: 0.12.5.6  
Status: Completed  
Code impact: None in closure  
Next: 0.12.6 — Documentation & Phase Closure
