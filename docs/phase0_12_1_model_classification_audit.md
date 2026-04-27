# Phase 0.12.1 — Model Classification Audit

## Subphase 0.12.1.0 — Definition & Documentation Lock

### Status

Completed.

### Purpose

0.12.1.0 defines and locks the complete audit structure for Phase 0.12.1 before any model classification work or code-level change is performed.

The goal is to prevent the Model Classification Audit from becoming an informal or partial inventory. Because Phase 0.12 depends on accurate read/write/domain responsibility boundaries, the audit must be explicit, exhaustive and traceable.

### Relationship To Phase 0.12

Phase 0.12 separates read-oriented, write-oriented and domain-owned model responsibilities across the current Stage 0 backend modules.

0.12.1 is the first implementation-oriented subphase of that phase, but it must start with a documentation lock because its work may need to be subdivided. The `.0` sub-subphase records the subdivision before any audit output is produced.

### Scope

Included:

- define the audit methodology for model classification
- register all 0.12.1 sub-subphases
- define READ, WRITE, DOMAIN, CONTRACT, INFRASTRUCTURE and HYBRID / TRANSITIONAL classification criteria
- define required inventory fields
- define rules for cross-layer usage review
- define the expected final audit output
- preserve the no-code-change rule for the definition lock

Excluded:

- code changes
- struct creation
- model renaming
- package moves
- handler rewiring
- mapping implementation
- public API changes
- business behavior changes

### Sub-Subphase Plan

| Sub-subphase | Name | Status | Purpose |
| --- | --- | --- | --- |
| 0.12.1.0 | Definition & Documentation Lock | ✅ Completed | Lock the audit method and subdivision before code review output is produced. |
| 0.12.1.1 | Model Inventory Extraction | ⬜ Pending | Extract the complete list of real model-like structures from the repository. |
| 0.12.1.2 | Model Classification | ⬜ Pending | Classify each structure as READ, WRITE, DOMAIN, CONTRACT, INFRASTRUCTURE or HYBRID / TRANSITIONAL. |
| 0.12.1.3 | Cross-Layer Usage Analysis | ⬜ Pending | Detect structures crossing HTTP, application, domain, repository, core or transport boundaries. |
| 0.12.1.4 | Problem Detection & Risk Mapping | ⬜ Pending | Identify hybrid responsibility, duplication, leakage and migration risk. |
| 0.12.1.5 | Target Separation Definition | ⬜ Pending | Define the target separation direction for each hybrid or transitional structure without implementing it. |
| 0.12.1.6 | Audit Consolidation & Closure | ⬜ Pending | Consolidate the final audit and prepare the repository for 0.12.2. |

### Classification Criteria

#### READ

A READ model is a structure used exclusively for output-oriented data.

Typical usage:

- HTTP responses
- application read results
- views or projections
- stable client-facing response shapes

READ models must not be used as mutation commands.

#### WRITE

A WRITE model is a structure used exclusively for input-oriented data.

Typical usage:

- HTTP request payloads
- command inputs
- mutation parameters
- intent-specific application inputs

WRITE models must not be reused as response/view structures.

#### DOMAIN

A DOMAIN model is owned by a module domain boundary and represents module-level semantics, invariants or internal state.

DOMAIN models may participate in behavior but should not become transport DTOs by accident.

#### CONTRACT

A CONTRACT model or interface is a deliberately minimal cross-module boundary shape.

Typical usage:

- provider interfaces
- minimal identity views
- small dependency inversion types

CONTRACT structures should remain narrow and explicit.

#### INFRASTRUCTURE

An INFRASTRUCTURE structure supports persistence, transport, protocol, framework or operational concerns.

Typical usage:

- WebSocket protocol envelopes
- error response envelopes
- repository persistence helpers
- token/session/config support shapes

Infrastructure structures are not automatically read/write models, but they may still expose hybrid responsibility if reused across layers.

#### HYBRID / TRANSITIONAL

A HYBRID / TRANSITIONAL structure currently carries more than one responsibility.

Common examples:

- one struct used for both request and response
- transport shape reused as domain state
- domain model returned directly through HTTP
- application result shape reused as mutation input
- core protocol structure used as both input envelope and output envelope

Hybrid classification does not imply immediate removal. It identifies the structure as a candidate for later separation.

### Required Inventory Fields

The 0.12.1 audit must record, at minimum:

| Field | Meaning |
| --- | --- |
| Module | Owning module or core area. |
| File | Exact file path. |
| Structure | Struct, interface or model-like type name. |
| Current responsibility | How it is currently used. |
| Classification | READ, WRITE, DOMAIN, CONTRACT, INFRASTRUCTURE or HYBRID / TRANSITIONAL. |
| Usage notes | Relevant boundaries crossed or behavior observed. |
| Risk | Low, Medium or High migration risk. |
| Target direction | Keep as-is, extract read model, extract write model, isolate domain model or introduce mapping. |

### Audit Method

The audit must be performed against the real repository code only.

Required inspection areas:

- `internal/modules/auth`
- `internal/modules/auth/app`
- `internal/modules/auth/domain`
- `internal/modules/user`
- `internal/modules/user/app`
- `internal/modules/user/domain`
- `internal/modules/usersettings`
- `internal/modules/usersettings/app`
- `internal/modules/usersettings/domain`
- relevant shared structures under `internal/core`

Additional files may be included if the real code shows that they define model-like structures or transport contracts.

### Execution Rules

0.12.1.0 has no code impact.

The following are prohibited in this sub-subphase:

- editing Go files
- creating new Go models
- moving packages
- changing handler behavior
- changing public routes
- changing public payloads
- changing standardized error behavior
- changing API versioning behavior

### Expected Output Of 0.12.1

The final audit must produce a complete model classification report that can guide later 0.12 work.

The audit output must identify:

- all discovered model-like structures
- their current classification
- cross-layer usage issues
- hybrid/transitional structures
- recommended separation direction
- risks that must be considered before extraction

### Closure Statement

0.12.1.0 is completed when the documentation records the subdivision and methodology for 0.12.1. The next correct step is 0.12.1.1 — Model Inventory Extraction.
