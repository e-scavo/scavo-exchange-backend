# 🧭 Roadmap — SCAVO Exchange Backend

## 🏗️ Stage 0 — Foundation

### Phase 0.1 — Baseline & Documentation
- Audit current project
- Define architecture
- Create documentation
- Establish roadmap

### Phase 0.2 — Core Infrastructure
- Database
- Redis
- Config system
- Observability
- Testing base

---

### Phase 0.5 — User Interaction & Application Surface
- authenticated user profile surface
- minimal non-wallet user metadata update contract
- user-facing read models for application bootstrap
- lightweight application-facing identity summary

---

## Phase 0.6 — Authenticated Application Bootstrap Consolidation

### 0.6.1 — Bootstrap Surface Boundary Clarification ✔
Completed

### 0.6.2 — Authenticated Surface Contract Alignment ✔
Completed

This subphase aligns all authenticated endpoints under a shared normalized context, ensuring contract-level consistency without breaking existing APIs.

### 0.6.3 — Session-Ready Bootstrap Read Model ✔
Completed

This subphase introduces a unified bootstrap endpoint that aggregates all authenticated surfaces into a single response, enabling a one-call frontend initialization flow.

### 0.6.4 — Application Surface Consistency Hardening ✔
Completed

This subphase finalizes the authenticated application surface by enforcing a consistent structural contract across all endpoints, introducing canonical envelopes while preserving backward compatibility.

---

### Phase 0.6 — Completed

Phase 0.6 is now fully completed.

It delivered:

- explicit authenticated surface boundaries
- aligned shared authenticated context
- unified bootstrap read model
- structural contract hardening across endpoints

### Next

Phase 0.7 — Application Layer Foundation

---

## 🧱 Stage 0 — Foundation (Remaining Phases)

### 0.7 — Application Layer Foundation ✔ Completed
Introduce an application layer to separate HTTP from business logic.

#### 0.7.1 — Application Layer Boundary Definition ✔ Completed
Initial application layer boundary defined and bootstrap use-case migrated.

#### 0.7.2 — Authenticated Surface Use Cases Extraction ✔ Completed

#### 0.7.3 — Wallet Management Use Cases Consolidation ✔ Completed

#### 0.7.4 — Handler Simplification & Contract Preservation ✔ Completed

### 0.8 — Standardized Error Model ✔ Completed
Define a consistent error contract across all endpoints.

#### 0.8.1 — Error Contract Definition ✔ Completed
Introduced the standardized error envelope foundation and centralized HTTP error writing.

#### 0.8.2 — Error Type System Introduction ✔ Completed
Introduced centralized internal app error typing, shared auth error catalog normalization and reusable error factories without changing the 0.8.1 transport envelope.

#### 0.8.3 — Auth Surface Error Standardization ✔ Completed
Auth handlers now consume centralized app-error factories and catalog resolution directly, while preserving the 0.8.1 envelope and avoiding cyclic import reintroduction.

#### 0.8.4 — Error Mapping Hardening & Contract Tests ✔ Completed
Added contract-focused tests for centralized app-error mapping and HTTP envelope serialization so representative `code/status/category` behavior and transport fallback paths are now frozen by automated coverage.

### 0.9 — API Versioning Strategy ✔ Completed
Formalize canonical API versioning, preserve legacy route compatibility and freeze the public contract-evolution rules needed before authorization and later domain growth.

#### 0.9.1 — Versioning Policy Definition ✔ Completed
Path-based versioning is now the canonical strategy, `v1` semantics are explicitly frozen and the legacy/non-canonical compatibility rule is documented.

#### 0.9.2 — Router Versioning Foundation ✔ Completed
Materialized the canonical `/api/v1/...` route surface in the real router while preserving the current legacy endpoints and reusing the same handler/application behavior.

#### 0.9.3 — Authenticated Surface Version Freezing ✔ Completed
Bound the current auth/bootstrap/profile/settings/session/wallet surfaces explicitly to the canonical `v1` contract semantics while preserving the same business behavior across legacy and canonical entry paths.

#### 0.9.4 — Version-aware Contract Testing ✔ Completed
Extended representative contract-level tests so legacy and canonical routes are protected from silent divergence across key authenticated and auth-adjacent scenarios.

#### 0.9.5 — Documentation Consolidation ✔ Completed
Consolidated the trunk documentation set so roadmap, status, architecture, testing and handoff now describe the same completed Phase 0.9 state and the same frontend/backend alignment rule.

### 0.10 — Authorization Layer ✔ Completed
Introduce a structured authorization layer on top of the authenticated Stage 0 foundation without breaking the stabilized transport and application contracts.

#### 0.10.1 — Authorization Model Definition ✔ Completed
Introduced the foundational authorization primitives under `internal/core/authorization`, including roles, permissions, static role-to-permission mapping and the authorization-subject model that later subphases will propagate and evaluate.

#### 0.10.2 — Authorization Context & Middleware ✔ Completed
Attach authorization subject data to the authenticated request lifecycle and extend middleware/context propagation without changing endpoint contracts yet.

#### 0.10.3 — Policy Evaluation Layer ✔ Completed
Introduce a centralized policy-evaluation boundary so handlers can ask authorization questions without embedding permission logic ad hoc.

#### 0.10.4 — Endpoint-Level Enforcement ✔ Completed
Apply progressive authorization checks to selected authenticated endpoints through centralized policy middleware while preserving backward compatibility and the current Stage 0 contract discipline.

#### 0.10.5 — Documentation & Contract Consolidation ✔ Completed
Align the trunk documentation set to the delivered authorization behavior and close the phase with a coherent architectural and operational narrative.

### 0.11 — Domain Module Pattern ✔ Completed
Standardize the internal structure of the current Stage 0 domain-facing modules so future growth builds on explicit HTTP / APP / DOMAIN boundaries, minimal internal contracts and reduced accidental coupling without changing the public API surface.

#### 0.11.1 — Domain Module Pattern Definition ✔ Completed
Define the standard internal module shape, layer responsibilities, dependency direction and naming conventions that the currently relevant modules will follow.

#### 0.11.2 — User Module Refactor ✔ Completed
Apply the pattern to `internal/modules/user` so user-facing transport handlers stop carrying orchestration responsibilities and the module becomes the first concrete reference implementation of the pattern.

This subphase is already delivered in repository state: the user module now exposes explicit `app`, `domain` and `repository` boundaries while preserving backward-compatible package access for existing consumers.

#### 0.11.3 — UserSettings Module Refactor ✔ Completed
Align `internal/modules/usersettings` to the same pattern while preserving its own configuration-oriented semantics and keeping it distinct from the user-entity boundary.

This subphase is now delivered in repository state: the user settings module exposes explicit `app`, `domain` and `repository` boundaries while preserving backward-compatible package access for existing consumers and keeping settings-specific validation/orchestration in the application layer.

#### 0.11.4 — Auth Module Alignment ✔ Completed
Aligned `internal/modules/auth` to the pattern conservatively while preserving the stabilized challenge, verify, wallet-management, session, bootstrap and authorization-adjacent behavior.

This subphase is delivered in repository state: `auth/app` owns application/runtime orchestration, `auth/domain` owns canonical wallet contracts, root `auth` remains a compatibility surface, wallet HTTP management is narrowed to transport/delegation, and repository migration is explicitly deferred behind a transitional façade.

#### 0.11.5 — Cross-Module Contract Consolidation ✔ Completed
Consolidate the explicit internal contracts between `auth`, `user` and `usersettings` so coordination remains possible without reintroducing ownership confusion or concrete cross-module coupling.

Delivered internal steps:

- 0.11.5.0 — Subphase Definition & Documentation Alignment ✔ Completed
- 0.11.5.1 — Dependency Mapping ✔ Completed
- 0.11.5.2 — Contract Extraction ✔ Completed
- 0.11.5.3 — Interface Alignment ✔ Completed
- 0.11.5.4 — Runtime Compatibility Validation ✔ Completed

#### 0.11.6 — Documentation & Phase Closure ✔ Completed
Updated the trunk documentation set so the adopted module pattern, its scope and its completion criteria are recorded consistently across roadmap, status, architecture, handoff and the dedicated phase document.

#### Expected outcome
Phase 0.11 closes with the same external HTTP behavior and a more uniform internal module structure, clearer ownership boundaries and safer foundations for later Stage 0 and post-Stage 0 work.

### 0.12 — Read / Write Model Separation 🚧 In Progress
Separate read-oriented, write-oriented and domain-owned model responsibilities without introducing full CQRS, event sourcing, public API changes or business behavior changes.

#### 0.12.0 — Phase Definition & Documentation Lock ✔ Completed
Define the complete phase scope, subphase order, compatibility constraints and documentary lock before changing code.

#### 0.12.1 — Model Classification Audit ✔ Completed
Audit real repository models and classify them as READ, WRITE, DOMAIN, INFRASTRUCTURE, CONTRACT or HYBRID/TRANSITIONAL.

Internal subdivision:

- 0.12.1.0 — Definition & Documentation Lock ✔ Completed
- 0.12.1.1 — Model Inventory Extraction ✔ Completed
- 0.12.1.2 — Model Classification ✔ Completed
- 0.12.1.3 — Cross-Layer Usage Analysis ✔ Completed
- 0.12.1.4 — Problem Detection & Risk Mapping ✔ Completed
- 0.12.1.5 — Target Separation Definition ✔ Completed
- 0.12.1.6 — Audit Consolidation & Closure ✔ Completed

### 0.12.1 Completion Note

0.12.1 is completed as a documentation and audit-only subphase. It produced the inventory, classification, cross-layer usage analysis, risk map and target separation definition required before read model extraction begins.

The next planned subphase after 0.12.1 was 0.12.2 — Read Model Extraction, now completed in this repository state.

#### 0.12.2 — Read Model Extraction ✔ Completed
Extract explicit read models for output/view/response paths while preserving current public response compatibility.

Internal subdivision:

- 0.12.2.0 — Definition & Documentation Lock ✔ Completed
- 0.12.2.1 — Read Model Design ✔ Completed
- 0.12.2.2 — Read Model Implementation ✔ Completed
- 0.12.2.3 — Domain/Application → Read Mapping ✔ Completed
- 0.12.2.4 — Response Alignment ✔ Completed
- 0.12.2.5 — Validation & Compatibility ✔ Completed
- 0.12.2.6 — Documentation & Closure ✔ Completed

0.12.2 is completed with explicit read model packages, Domain/Application → Read mapping functions, response alignment, compatibility validation and closure documentation. Public HTTP contracts remain unchanged.

#### 0.12.3 — Write Model Isolation ✔ Completed
Isolate input and mutation-oriented models so write paths do not reuse read models accidentally.

Internal subdivision:

- 0.12.3.0 — Definition & Documentation Lock ✔ Completed
- 0.12.3.1 — Write Model Design ✔ Completed
- 0.12.3.2 — Write Model Implementation ✔ Completed
- 0.12.3.3 — Write → Domain Mapping ✔ Completed
- 0.12.3.4 — Handler Alignment ✔ Completed
- 0.12.3.5 — Validation & Compatibility ✔ Completed
- 0.12.3.6 — Documentation & Closure ✔ Completed

#### 0.12.4 — Mapping Layer Introduction ✔ Completed
Introduce an explicit centralized mapping layer for Domain → Read Model and Write Model → Domain transformations, consolidating mapper ownership outside read model, write model and application packages.

Internal subdivision:

- 0.12.4.0 — Definition & Documentation Lock ✔ Completed
- 0.12.4.1 — Mapping Layer Design ✔ Completed
- 0.12.4.2 — Mapping Layer Implementation ✔ Completed
- 0.12.4.3 — Mapping Consolidation ✔ Completed
- 0.12.4.4 — Application Refactor ✔ Completed
- 0.12.4.5 — Validation & Compatibility ✔ Completed
- 0.12.4.6 — Documentation & Closure ✔ Completed

#### 0.12.5 — Contract Alignment ✔ Completed
Review the 0.11 contracts and align UserProvider/UserSettingsProvider usage with the read/write separation, write/read model boundaries and centralized mapping ownership.

Internal subdivision:

- 0.12.5.0 — Definition & Documentation Lock ✔ Completed
- 0.12.5.1 — Contract Inventory & Classification ✔ Completed
- 0.12.5.2 — Contract Normalization Design ✔ Completed
- 0.12.5.3 — Contract Alignment Implementation ✔ Completed
- 0.12.5.4 — Handler Contract Adjustment ✔ Completed
- 0.12.5.5 — Validation & Compatibility ✔ Completed
- 0.12.5.6 — Documentation & Closure ✔ Completed

#### Expected outcome
Phase 0.12 closes with clearer internal model responsibilities, explicit mapping boundaries and unchanged public HTTP/API behavior.

### 0.13 — Provider Layer Consolidation ✔ Completed
Consolidate the Provider Layer as the explicit entry point to domain services after Read / Write Model Separation.

Phase 0.13 continues the Stage 0 foundation sequence by moving from model/contract clarity to orchestration clarity. The phase does not add public API behavior; it establishes a provider-owned runtime boundary so later hardening can build on a stable composition model.

#### 0.13.0 — Definition & Documentation Lock ✔ Completed
Locked the phase definition, corrected roadmap consistency and registered the subphase plan before code changes.

#### 0.13.1 — Provider Inventory & Classification ✔ Completed
Inventoried existing provider-like responsibilities, missing provider boundaries, compatibility wiring and direct access patterns without changing Go code.

#### 0.13.2 — Provider Interface Design ✔ Completed
Defined narrow provider interfaces aligned with the 0.13.1 inventory findings, read/write model separation and module ownership. No Go code or public API behavior was changed.

#### 0.13.3 — Provider Implementation ✔ Completed
Implemented provider boundaries incrementally while preserving compatibility.

#### 0.13.4 — Application Integration ✔ Completed
Integrated runtime HTTP wiring with the consolidated provider boundary while preserving public contracts and compatibility.

#### 0.13.5 — Validation & Compatibility ✔ Completed
Validate public contract preservation and runtime compatibility.
It records compatibility after provider application integration. The validation baseline confirms that the provider-oriented routing path remains build/test compatible and preserves existing public contracts.

#### 0.13.6 — Documentation & Closure ✔ Completed
Closed the phase documentation, removed generic repeated status-only updates, corrected misplaced roadmap text and restored the Phase 0.12 → Phase 0.13 narrative across trunk documentation. The closure records not only that the phase completed, but why Provider Layer Consolidation was the next architectural step after read/write model separation.

### 0.14 — Observability & Diagnostics Foundation ✅ Completed
Introduce observability and diagnostics as a Stage 0 foundation capability without changing public API behavior.

Phase 0.14 follows the Provider Layer Consolidation completed in 0.13. The provider boundary made runtime composition explicit; 0.14 now makes the same path diagnosable through request correlation, structured logging, internal error context, flow tracing and a minimal diagnostics surface.

The phase is intentionally limited to internal visibility. External metrics platforms, dashboards, Prometheus, OpenTelemetry and public API changes remain out of scope.

#### 0.14.0 — Phase Definition & Documentation Lock ✅ Completed
Lock the phase definition, correct the roadmap from the previous placeholder phase name and register the observability subphase plan before code changes.

#### 0.14.1 — Correlation Model (Request ID / Trace) ✅ Completed
Introduce request correlation and propagate request-scoped trace context through the backend flow.

#### 0.14.2 — Logging Standardization ✅ Completed
Define structured logging conventions and apply request context to HTTP and provider-level runtime paths.

#### 0.14.3 — Error Context Enrichment ✅ Completed
Add internal diagnostic metadata to errors while preserving the existing public error contract.

#### 0.14.4 — Flow Tracing Integration ✅ Completed
Instrument key flow transitions across HTTP, provider, application, domain and repository boundaries.

#### 0.14.5 — Diagnostics Surface Exposure ✅ Completed
Expose a minimal diagnostics-oriented surface or hook set compatible with the current architecture.

#### 0.14.6 — Validation & Documentation ✅ Completed
Validate compatibility with `go test ./...`, confirm no public behavior change and close the phase narratively. Completed after reconciling the Phase 0.14 documentation with the implemented request correlation, logging, error-context, flow-tracing and diagnostics surface work.

Phase 0.14 outcome: the backend remains contract-compatible while gaining an internal observability foundation: request correlation, structured request-scoped logging, safe diagnostic error context, minimal flow tracing events and `GET /diagnostics`.

---

### 0.15 — Contract Hardening & Freeze ✅ Completed
Formalize, validate and freeze the existing backend contracts after the completed observability foundation.

Phase 0.15 follows the Observability & Diagnostics Foundation completed in 0.14. Phase 0.12 clarified data direction, Phase 0.13 clarified runtime composition and Phase 0.14 made that runtime path diagnosable. Phase 0.15 now hardens the contracts exposed and consumed along that path so future evolution can remain controlled.

The phase is intentionally limited to contract validation and freeze discipline. It does not introduce new features, new routes, business logic changes, architecture changes or additional observability platforms.

#### 0.15.0 — Phase Definition & Documentation Lock ✅ Completed
Register Phase 0.15 in the trunk documentation, define scope and subphase order and lock the baseline before any code or contract audit changes.

#### 0.15.1 — HTTP Contract Audit ✅ Completed
Audit the existing HTTP route surface and validate method, path, status and payload expectations from the real repository.

0.15.1 records 39 registered HTTP route entries and 22 unique behavior contracts across foundation, diagnostics, WebSocket upgrade and auth/account/wallet surfaces. It confirms that legacy auth routes and `/api/v1` auth routes are active paired contracts backed by the same handler set.

#### 0.15.2 — Error Contract Alignment ✅ Completed
Validate and align public error envelopes around the existing standardized error model: `code`, `message` and `details`.

0.15.2 hardens the canonical error envelope by making `error.details` a required JSON object. Detail-free errors now serialize `details: {}` instead of omitting the field, while existing codes, messages, statuses and handler decisions remain unchanged.

#### 0.15.3 — Provider Contract Validation ✅ Completed

Validate internal HTTP → Provider → Application contracts and confirm that Phase 0.13 provider boundaries remain explicit and stable.

0.15.3 validates the internal provider boundary without changing runtime behavior. Compile-time assertions now prove that `*auth.Application` satisfies the handler-facing `AuthProvider` contract and that the concrete user and usersettings services satisfy the minimal cross-module provider contracts consumed by auth.

#### 0.15.4 — Response Schema Normalization ✅ Completed
Normalize response schemas where actual repository evidence shows uncontrolled drift or ambiguous historical variation.

#### 0.15.5 — Contract Freeze Enforcement ✅ Completed
Define the freeze policy that controls what can change, what must be versioned and how future contract drift is detected.

0.15.5 freezes the audited Stage 0 public HTTP surface, canonical error envelope, provider boundary assertions and response serialization policy. Contract drift is now guarded by dedicated regression tests and by explicit documentation rules for future changes.

#### 0.15.6 — Validation & Documentation ✅ Completed
Validate the full system with `go test ./...`, reconcile trunk documentation and close the phase narratively.

0.15.6 records the successful developer-environment `go test ./...` validation after 0.15.5 and reconciles the trunk documentation so Phase 0.15 closes with the correct chronological subphase order and no active 0.15 subphase.

Phase 0.15 outcome: explicit, audited and frozen contracts for the current Stage 0 backend surface, preserving public behavior while preparing controlled future evolution.

### Phase 0.15 Closure Note — Contract Hardening & Freeze

Phase 0.15 completed after the validation and documentation closure in 0.15.6.

The next roadmap phase is intentionally not defined by 0.15.6.

The final baseline records that successful response payloads remain intentionally unwrapped and compatible, JSON response metadata and defensive error fallback shape are aligned with the canonical contract introduced in 0.15.2, provider boundary expectations are compile-time guarded and representative frozen contracts are covered by regression tests.

---

## 🚀 Stage 1 — Application Surface & Product Capabilities

Stage 1 transforms the stable Stage 0 foundation into usable application-level product capabilities.

Stage 0 closed the backend foundation by stabilizing architecture, observability, diagnostics, error shape, provider boundaries, response metadata and contract freeze rules. Stage 1 must not reopen that foundation work. It builds on it by defining real application flows, account-level capabilities, authorization behavior, data interaction patterns, write flows and end-to-end validation.

### Objective

Build product-facing backend capabilities on top of the frozen Stage 0 contracts.

Stage 1 focuses on:

- real application use cases
- coherent account and identity surfaces
- wallet management usability
- permission and ownership behavior
- data interaction patterns
- mutation flows
- end-to-end product validation

### Scope

Included:

- application use-case consolidation
- account and identity product capabilities
- authorization and permission model refinement
- data interaction patterns
- mutation and write flows
- system behavior consistency
- end-to-end validation and closure

Excluded:

- reopening Stage 0 foundation architecture
- changing frozen Stage 0 contracts without explicit versioning
- blockchain capability implementation
- DEX contracts
- frontend implementation
- infrastructure overengineering

### Phase 1.1 — Application Use Case Consolidation

Define the real application use cases that Stage 1 must support.

This phase aligns the existing endpoint surface with product flows before expanding functionality.

Focus areas:

- identify real backend-supported application flows
- map existing endpoints to use cases
- detect orphaned, unclear or duplicate surfaces
- define product-level success criteria
- preserve Stage 0 contracts while clarifying intended usage

#### 1.1.0 — Phase Definition & Documentation Lock

- register Phase 1.1 scope in trunk documentation
- confirm Stage 1 starts from the frozen Stage 0 baseline
- define the audit boundary before touching code
- avoid feature implementation during phase definition

#### 1.1.1 — Endpoint Surface Inventory

- inventory existing HTTP and WebSocket surfaces
- classify endpoints by product intent
- identify current consumers and expected usage
- preserve frozen Stage 0 contracts during inventory

#### 1.1.2 — Application Flow Mapping

- map endpoints to real application flows
- define flow entry points and dependencies
- distinguish bootstrap, account, settings, wallet and system flows
- identify missing flow documentation without adding features prematurely

#### 1.1.3 — Orphan, Duplicate & Ambiguous Surface Review

- detect endpoints without a clear Stage 1 use case
- detect duplicate or overlapping surfaces
- document ambiguous ownership or naming
- propose safe consolidation paths without breaking compatibility

#### 1.1.4 — Product Success Criteria Definition

- define success criteria per supported flow
- document expected inputs, outputs and error behavior
- align criteria with frozen contracts
- prepare the basis for end-to-end validation

#### 1.1.5 — Use Case Contract Documentation

- document the approved Stage 1 application use cases
- connect each use case to existing contracts
- record excluded or deferred use cases
- prevent scope drift before later Stage 1 phases

#### 1.1.6 — Validation & Closure

- validate documentation coherence
- run the project test suite when code was touched
- close Phase 1.1 with explicit handoff to Phase 1.2

### Phase 1.2 — Account & Identity Product Capabilities

Transform the identity infrastructure into a usable product surface while preserving the frozen Stage 0 contracts.

Focus areas:

- complete session lifecycle
- coherent account surface
- usable wallet management
- real user settings behavior
- end-to-end account flows

#### 1.2.0 — Phase Definition & Documentation Lock

- register Phase 1.2 scope in trunk documentation
- preserve the Stage 0 frozen contracts before identity product work
- define account, session, wallet and settings boundaries
- avoid implementation before account flows are explicitly scoped

#### 1.2.1 — Session Lifecycle Completion

- login → session → validation → refresh → logout
- session consistency across endpoints
- expiration and renewal behavior

#### 1.2.2 — Account Surface Consolidation

- `/me`
- `/profile`
- `/settings`
- `/wallets`
- structural consistency across account endpoints
- read normalization

#### 1.2.3 — Wallet Management Usability

- wallet listing
- wallet add/remove flows
- primary wallet behavior
- ownership validation

#### 1.2.4 — User Settings Productization

- consistent read behavior
- real update behavior
- defaults
- validations

#### 1.2.5 — Account-Level Authorization Refinement

- endpoint permissions
- ownership checks
- session consistency

#### 1.2.6 — End-to-End Flow Validation

- login → bootstrap → update → logout
- real edge cases
- functional consistency

#### 1.2.7 — Documentation & Closure

- flow documentation
- final validation
- phase closure

### Phase 1.3 — Authorization & Permission Model

Refine authorization from foundational enforcement into product-level permission behavior.

Focus areas:

- roles
- permissions
- ownership checks
- endpoint-level enforcement
- session-aware permission consistency

#### 1.3.0 — Phase Definition & Documentation Lock

- register Phase 1.3 scope in trunk documentation
- define authorization boundaries before implementation
- preserve Stage 0 authentication and error contracts
- confirm which account-level capabilities require permission checks

#### 1.3.1 — Permission Surface Inventory

- inventory endpoints requiring authenticated access
- identify current authorization checks
- classify account, wallet, settings and system operations
- detect unauthenticated or under-specified permission boundaries

#### 1.3.2 — Role & Capability Mapping

- define product-level roles or capabilities where applicable
- map capabilities to endpoint groups
- document default access expectations
- avoid introducing roles that are not required by real flows

#### 1.3.3 — Ownership Rule Definition

- define ownership rules for account-level resources
- document wallet ownership expectations
- align settings mutations with authenticated account context
- prevent cross-account access drift

#### 1.3.4 — Enforcement Alignment

- align endpoint checks with documented permission rules
- preserve existing public contracts
- normalize authorization failure behavior through the frozen error envelope
- avoid business logic changes outside permission enforcement

#### 1.3.5 — Authorization Regression Coverage

- add or update tests for permission boundaries
- cover denied, allowed and ownership-sensitive paths
- confirm error shape consistency
- preserve compatibility with existing consumers

#### 1.3.6 — Validation & Closure

- run `go test ./...`
- reconcile trunk documentation
- close Phase 1.3 with handoff to data interaction patterns

### Phase 1.4 — Data Interaction Patterns

Standardize the way product surfaces expose data interaction.

Focus areas:

- pagination
- filtering
- sorting
- response consistency
- read model behavior

#### 1.4.0 — Phase Definition & Documentation Lock

- register Phase 1.4 scope in trunk documentation
- define data interaction boundaries before code changes
- preserve existing response compatibility
- distinguish read patterns from write flows

#### 1.4.1 — Read Surface Inventory

- inventory list and read endpoints
- identify collection responses and single-resource responses
- document current query parameters
- detect undocumented or inconsistent read behavior

#### 1.4.2 — Pagination Contract Definition

- define supported pagination inputs
- define metadata expectations where applicable
- document defaults and limits
- preserve compatibility for existing non-paginated responses

#### 1.4.3 — Filtering Contract Definition

- identify supported filters per surface
- document unsupported or deferred filters
- validate error behavior for invalid filters
- avoid adding generic filtering without product need

#### 1.4.4 — Sorting Contract Definition

- identify sortable fields per collection surface
- document default ordering
- validate unsupported sorting behavior
- preserve deterministic output where required

#### 1.4.5 — Read Model Consistency Validation

- validate read model shape consistency
- align data interaction responses with Stage 0 response policy
- add tests where contract drift is likely
- document compatibility constraints

#### 1.4.6 — Validation & Closure

- run `go test ./...`
- reconcile trunk documentation
- close Phase 1.4 with handoff to mutation and write flows

### Phase 1.5 — Mutation & Write Flows

Complete coherent write behavior for product-facing capabilities.

Focus areas:

- create flows
- update flows
- delete flows
- validation behavior
- write model consistency

#### 1.5.0 — Phase Definition & Documentation Lock

- register Phase 1.5 scope in trunk documentation
- define mutation boundaries before implementation
- preserve Stage 0 error and response contracts
- distinguish product mutations from infrastructure changes

#### 1.5.1 — Write Surface Inventory

- inventory create, update and delete-capable surfaces
- document current write models
- identify missing validation expectations
- detect unsafe or ambiguous mutation behavior

#### 1.5.2 — Create Flow Consolidation

- define create request expectations
- document validation and conflict behavior
- align success responses with frozen response policy
- preserve ownership and authorization expectations

#### 1.5.3 — Update Flow Consolidation

- define update semantics per resource
- document partial versus full update behavior where applicable
- validate immutable fields
- align errors for invalid or unauthorized updates

#### 1.5.4 — Delete / Disable Flow Consolidation

- define delete, disable or detach semantics per resource
- document soft-delete versus hard-delete expectations where applicable
- validate ownership checks
- align not-found and conflict behavior

#### 1.5.5 — Write Model Validation Coverage

- add or update tests for mutation contracts
- cover validation, conflict, ownership and authorization cases
- confirm canonical error envelope consistency
- preserve existing public behavior unless explicitly versioned

#### 1.5.6 — Validation & Closure

- run `go test ./...`
- reconcile trunk documentation
- close Phase 1.5 with handoff to system behavior consistency

### Phase 1.6 — System Behavior Consistency

Validate consistent behavior across product-level flows and edge cases.

Focus areas:

- idempotency where applicable
- edge cases
- conflict behavior
- cross-endpoint consistency
- contract-preserving behavior under failure

#### 1.6.0 — Phase Definition & Documentation Lock

- register Phase 1.6 scope in trunk documentation
- define behavior consistency boundaries before code changes
- preserve Stage 0 frozen contracts
- focus on behavior guarantees rather than new features

#### 1.6.1 — Cross-Endpoint Behavior Audit

- compare behavior across related endpoints
- identify inconsistent success, empty, conflict and not-found handling
- document expected behavior per flow
- avoid changing business semantics without explicit rationale

#### 1.6.2 — Idempotency & Retry Semantics Review

- identify operations where idempotency matters
- document retry-safe and non-retry-safe behavior
- validate duplicate request outcomes
- preserve compatibility for existing clients

#### 1.6.3 — Conflict & Edge Case Alignment

- define conflict behavior for account, wallet, settings and data operations
- validate expired, missing, duplicate and invalid resource scenarios
- align edge-case errors with canonical error envelope
- document known deferred behavior

#### 1.6.4 — Failure Mode Consistency

- validate defensive fallback behavior
- confirm request correlation survives failure paths
- align timeout and internal failure responses
- preserve diagnostics without leaking internals

#### 1.6.5 — Behavior Regression Coverage

- add tests for cross-endpoint behavior guarantees
- cover edge cases and failure modes
- ensure Stage 0 freeze tests remain valid
- prevent product-level behavior drift

#### 1.6.6 — Validation & Closure

- run `go test ./...`
- reconcile trunk documentation
- close Phase 1.6 with handoff to Stage 1 end-to-end validation

### Phase 1.7 — End-to-End Validation & Closure

Validate Stage 1 as a complete product-capability layer.

Focus areas:

- full flow validation
- functional tests
- documentation reconciliation
- Stage 1 closure readiness

#### 1.7.0 — Phase Definition & Documentation Lock

- register Phase 1.7 scope in trunk documentation
- define final Stage 1 validation boundaries
- confirm all prior Stage 1 phases are represented coherently
- avoid introducing new product scope during closure

#### 1.7.1 — End-to-End Flow Matrix

- document complete Stage 1 flow matrix
- cover login, bootstrap, account, settings, wallet, authorization and mutation flows
- define expected success and failure outcomes
- connect each flow to test and documentation evidence

#### 1.7.2 — Functional Validation Pass

- execute full backend test suite
- validate representative end-to-end flows
- confirm frozen Stage 0 contracts remain intact
- record validation evidence in trunk documentation

#### 1.7.3 — Documentation Reconciliation

- reconcile all trunk Markdown documents
- ensure Stage 1 phases and subphases are ordered consistently
- remove stale future references
- preserve historical narrative without rewriting prior stages

#### 1.7.4 — Stage 1 Closure Assessment

- verify Stage 1 objectives were met
- document remaining deferred work
- confirm no unfinished Stage 1 subphase remains active
- prepare the system for the next stage definition

#### 1.7.5 — Final Validation & Handoff

- run final validation commands
- close Stage 1 narratively
- prepare handoff for the next stage without defining it prematurely
- create a stable historical baseline after acceptance

---

## 🔗 Stage 2 — Blockchain Integration

- RPC integration
- Balance reads
- Allowances
- Indexing

---

## 🧠 Stage 3 — DEX Contracts

- AMM design
- Factory
- Pools
- Router
- Testing

---

## ⚙️ Stage 4 — DEX Backend

- Quotes
- Routing
- Swap intents
- Liquidity APIs

---

## 🌐 Stage 5 — APIs

- REST
- WebSocket
- Frontend contract

---

## 🔄 Stage 6 — Hybrid Expansion

- Internal ledger
- P2P design

---

## 🔐 Stage 7 — Security & Ops

- Hardening
- Observability
- Compliance hooks

---

## 🧪 Stage 8 — Testing

- Integration tests
- Internal testing
- Release candidate

