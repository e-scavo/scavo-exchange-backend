# Phase 0.7 — Application Layer Foundation

## Objective

Establish a clear and consistent application-layer boundary across the backend, separating HTTP transport concerns from business orchestration logic.

---

## Initial Context

Before Phase 0.7, the backend exhibited a mixed orchestration model:

- HTTP handlers partially controlled execution flow
- Services were invoked directly from handlers
- Response assembly was inconsistent
- Transport and domain concerns were coupled

This resulted in:

- duplicated orchestration logic
- inconsistent handler responsibilities
- reduced maintainability
- limited scalability for future phases

---

## Architectural Goal

Introduce a strict layered execution model:

HTTP Handler  
→ Application Layer  
→ Services / Stores  
→ Response  

---

## Phase Breakdown

### 0.7.1 — Application Layer Boundary Definition ✔

- Introduced `Application` struct
- Bootstrap used as initial use-case pilot
- Defined first explicit orchestration boundary

---

### 0.7.2 — Authenticated Surface Use Cases Extraction ✔

Moved core authenticated flows into application layer:

- login
- me
- session
- bootstrap

Handlers became delegators.

---

### 0.7.3 — Wallet Management Use Cases Consolidation ✔

Extended application-layer pattern to wallet flows:

- wallet list
- wallet link (challenge / verify)
- wallet account merge (challenge / verify)
- set primary wallet
- detach (check / execute)

Removed remaining handler-driven orchestration.

---

### 0.7.4 — Handler Simplification & Contract Preservation ✔

- Unified transport helpers:
  - request decoding
  - claims extraction
  - error JSON writing
- Simplified all handlers
- Removed duplication
- Preserved all public contracts

---

## Final Architecture

After Phase 0.7:

### Transport Layer (HTTP Handlers)

Responsibilities:

- decode request
- extract claims
- invoke application
- map errors to HTTP
- return JSON

Constraints:

- no business logic
- no orchestration

---

### Application Layer

Responsibilities:

- define use cases
- orchestrate services
- centralize execution flow
- ensure deterministic behavior

Constraints:

- no HTTP concerns
- no low-level DB logic

---

### Services / Stores

Responsibilities:

- execution logic
- persistence
- validations

---

## Guarantees

- strict separation of concerns
- deterministic execution paths
- no contract changes across refactors
- unified handler structure
- extensibility for future phases

---

## Result

Phase 0.7 delivers a fully application-driven backend for the auth module:

- all use cases centralized
- handlers reduced to transport
- architecture consistent across endpoints

---

## What This Enables

Phase 0.7 unlocks the next architectural steps:

- 0.8 — Standardized Error Model
- 0.9 — API Versioning Strategy
- 0.10 — Authorization Layer

---

## Conclusion

Phase 0.7 is now fully completed.

The backend transitions from a partially structured system into a clean, layered architecture with a well-defined execution model.

This establishes a solid foundation for all subsequent phases.