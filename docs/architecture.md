# Architecture

## 🧠 Overview

The SCAVO Exchange Backend is designed around a **wallet-first identity architecture** that progressively evolves into a **durable account model** suitable for exchange-grade ownership and future multi-auth identity operations.

The architecture intentionally separates:

- authentication mechanism
- wallet identity representation
- durable platform user abstraction
- persisted wallet ownership
- authenticated wallet-management contracts

---

## 🧩 Core Layers

### 1. Transport Layer
- HTTP API
- JSON-based communication
- stateless request handling

### 2. Auth Layer
Located in:

- `internal/modules/auth`

Responsibilities:

- wallet challenge generation
- wallet signature verification
- JWT issuance
- identity resolution
- ownership enforcement
- authenticated wallet-linking flows

### 3. User Layer
Located in:

- `internal/modules/user`

Responsibilities:

- durable user creation
- durable user resolution
- future auth-provider expansion

### 4. Persistence Layer
- PostgreSQL (primary)
- in-memory fallback (dev/testing)

---

## 🔐 Identity Model Evolution

### Pre 0.4.6
- identity was session-oriented
- wallet state was not durable

### 0.4.6 — Wallet Identity Persistence
- wallet identity stored in `auth_wallet_identities`
- wallet address becomes a stable registry entry

### 0.4.7 — Unified Identity Model
- wallet identity linked to durable user
- `user_id` introduced
- JWT identity unified around durable platform user

### 0.4.8 — Ownership Model Introduction
wallet identity becomes a first-class ownership entity with:

- `id`
- `address`
- `user_id`
- `linked_at`
- `detached_at`
- `is_primary`

### 0.4.9 — Authenticated Wallet Linking Contract
the architecture adds a dedicated authenticated linking flow, still based on challenge + signature verification, but now explicitly bound to the current authenticated user.

This is the first backend-managed wallet operation that acts on ownership under an existing session rather than during initial login bootstrap.

### 0.4.10 — Wallet-Owned Account Merge Execution
the architecture now adds a second authenticated ownership operation that allows the current user to absorb a wallet-owned source account after an explicit merge challenge is signed by the source wallet.

---

## 🏷️ Ownership Model

### Core Rules

1. a wallet belongs to exactly one user
2. a user can own multiple wallets
3. only one wallet per user can be primary
4. wallet ownership cannot be reassigned across users
5. authenticated wallet linking adds secondary wallets only
6. 0.4.9 does not switch primary ownership during linking
7. 0.4.10 preserves the current target primary wallet when a merge occurs
8. 0.4.11 allows explicit primary-wallet reassignment only within the current owner's wallet set
9. 0.4.12 exposes detach eligibility as a guarded evaluation contract and 0.4.13 adds detach execution only for already eligible owned wallets

---

## 🏷️ Ownership Metadata

| Field | Description |
|------|-------------|
| `user_id` | owning durable user |
| `linked_at` | ownership creation timestamp |
| `is_primary` | primary-wallet flag |

---

## 🧾 Challenge Model

### Pre 0.4.9
wallet challenge was effectively used only for authentication bootstrap.

### 0.4.9
wallet challenges now include:

- `purpose`
- `requested_by_user_id`

### Challenge purposes

- `auth_bootstrap`
- `wallet_link`
- `account_merge`

This avoids reusing the same challenge semantics blindly across two very different operations.

---

## 🔄 Authentication Flow (Wallet Login)

1. client requests login challenge
2. backend persists challenge with `auth_bootstrap` purpose
3. client signs message
4. backend verifies signature
5. challenge is consumed
6. wallet identity is resolved
7. durable user is resolved or created
8. ownership is enforced
9. JWT is issued

---

## 🔄 Authenticated Wallet Linking Flow

1. user already holds valid JWT
2. client requests link challenge:
   - `POST /auth/wallets/link/challenge`
3. backend persists challenge with:
   - `purpose = wallet_link`
   - `requested_by_user_id = current user`
4. user signs with the secondary wallet
5. client submits:
   - `POST /auth/wallets/link/verify`
6. backend validates:
   - challenge existence
   - challenge freshness
   - purpose correctness
   - requesting user correctness
   - signature correctness
7. backend resolves wallet identity
8. backend rejects ownership conflict if wallet belongs elsewhere
9. backend attaches wallet as non-primary
10. backend returns updated wallet inventory

---

## 🔄 Authenticated Wallet Account Merge Flow

1. user already holds valid JWT
2. client requests merge challenge:
   - `POST /auth/account/merge/wallet/challenge`
3. backend persists challenge with:
   - `purpose = account_merge`
   - `requested_by_user_id = current user`
4. source wallet signs the merge challenge
5. client submits:
   - `POST /auth/account/merge/wallet/verify`
6. backend validates:
   - challenge existence
   - challenge freshness
   - purpose correctness
   - requesting user correctness
   - signature correctness
   - source wallet ownership existence
7. backend derives the source user from wallet ownership
8. backend atomically reassigns all source-user wallets to the current target user
9. backend returns updated wallet inventory

## 🔌 API Layer

### Auth endpoints
- `/auth/login`
- `/auth/wallet/challenge`
- `/auth/wallet/verify`
- `/auth/me`
- `/auth/session`
- `PATCH /auth/me`

`/auth/me` is now the additive authenticated profile surface intended for application bootstrap, while `PATCH /auth/me` is the minimal self-scoped non-wallet metadata mutation for `display_name`. `/auth/session` remains the raw session/claims-oriented contract and `/auth/wallets` remains the fuller inventory contract.

### Wallet ownership endpoints
- `/auth/wallets`
- `/auth/wallets/link/challenge`
- `/auth/wallets/link/verify`
- `/auth/account/merge/wallet/challenge`
- `/auth/account/merge/wallet/verify`

---

## 🧾 JWT Design

JWT tokens are:

- stateless
- short-lived
- self-contained

Claims include:

- `user_id`
- `wallet_id`
- `wallet_address`
- `auth_method`

Wallet linking and wallet-owned account merge do not mint a fresh token because both operate under an already authenticated durable session.

---

## 🗄️ Data Model

### `auth_wallet_challenges`
stores challenge lifecycle and, from 0.4.9 onward, also stores operation metadata:

- `purpose`
- `requested_by_user_id`
- issued / expires / used lifecycle

### `auth_wallet_identities`
stores wallet registry and ownership:

- `id`
- `address`
- `user_id`
- `linked_at`
- `detached_at`
- `is_primary`

### `users`
stores durable user abstraction:

- login-independent identity
- wallet-backed users now
- future auth-provider aggregation later

---

## ⚙️ Design Decisions

### Wallet-first approach
chosen because it aligns with crypto-native UX and reduces early auth-surface complexity.

### Separation of identity and ownership
wallet identity is not the same as durable user identity.
Ownership is explicit rather than inferred.

### Challenge-purpose separation
0.4.9 and 0.4.10 extend the challenge system rather than introducing parallel challenge subsystems, but still keep semantic separation through `purpose`.

### Incremental evolution
each subphase introduces one structural improvement while preserving previous behavior.

---

## ⚠️ Constraints

Still intentionally not supported:

- wallet unlink
- arbitrary cross-user wallet transfer outside wallet-signed merge
- multi-auth merge resolution
- user-record archival after merge

---

## 🚧 Current Detached-Identity Audit Readiness (0.4.15)

Detached wallet identities now preserve minimal lifecycle metadata through `detached_at`.

This keeps the architecture compact while allowing the system to distinguish reusable identities that were previously detached from identities that have never gone through detach execution.


### Later phases
- account consolidation
- multi-auth identity merging
- recovery flows
- compliance-ready identity expansion

---

## 🧩 Summary

At the end of 0.4.15:

- wallet authentication is stable
- durable identity is stable
- wallet ownership is stable
- authenticated wallet linking is implemented
- wallet-owned account merge execution is implemented
- explicit primary-wallet switching is implemented
- authenticated wallet detach eligibility is implemented
- authenticated wallet detach execution is implemented for already eligible owned wallets
- detached wallet identities are explicitly reusable after detach
- detached wallet identities now preserve minimal audit-ready lifecycle metadata through `detached_at`
- the backend is structurally ready to evolve into richer detached-identity observability only if later phases require it


---

## Wallet Identity Read Model (0.4.16)

By the end of Phase 0.4.16, the backend explicitly separates:

- the internal wallet identity domain model
- the external authenticated wallet inventory read model

The internal model preserves ownership and lifecycle state such as:

- `user_id`
- `linked_at`
- `detached_at`
- `is_primary`

The external read model exposed by `GET /auth/wallets` now returns:

- `id`
- `address`
- `user_id`
- `linked_at`
- `detached_at`
- `is_primary`
- `status`

### Status Derivation

The `status` field is derived conservatively:

- `active` when the wallet is currently linked to a user
- `detached` when the wallet has no current owner and `detached_at` is present
- `unlinked` when the wallet has no current owner and no detach evidence

For the authenticated wallet inventory route, the expected real-world case remains `active`, because the endpoint lists wallets currently owned by the authenticated user. The value of the enrichment is that previous detach lifecycle evidence can remain visible after reattachment.

### Architectural Benefit

This keeps mutation-focused ownership logic separate from read-oriented inventory projection. It improves frontend and debugging visibility without changing ownership or lifecycle business rules.

## Phase 0.7 — Application Layer Foundation

### Architectural Objective

Phase 0.7 introduces an explicit application layer between HTTP transport and module-level execution logic, establishing a repeatable orchestration boundary for authenticated use cases.

### Boundary Introduced

The backend now follows this architectural flow for the auth module:

HTTP Handler  
→ Application Layer  
→ Services / Stores  
→ Response  

### Phase 0.7.1 — Application Layer Boundary Definition

Introduced the first explicit `Application` entry point in `internal/modules/auth`, using bootstrap as the initial pilot use case.

### Phase 0.7.2 — Authenticated Surface Use Cases Extraction

Moved the main authenticated surface use cases into application orchestration:

- login
- me
- session
- bootstrap

### Phase 0.7.3 — Wallet Management Use Cases Consolidation

Extended the same application-layer pattern to all wallet-related use cases:

- wallet list
- wallet link challenge / verify
- wallet account merge challenge / verify
- set primary wallet
- detach check / execute

### Phase 0.7.4 — Handler Simplification & Contract Preservation

Unified the transport layer by standardizing:

- request decoding
- claims extraction
- error JSON writing

This leaves handlers as minimal transport adapters while preserving all public contracts.

### Result

At the end of Phase 0.7, the auth module architecture is now consistently layered:

- Handlers → transport only
- Application → orchestration
- Services / Stores → execution


## Phase 0.9 — API Versioning Strategy

### Architectural Objective

Phase 0.9 introduces an explicit API versioning layer for the already-stabilized authenticated surface so the backend can evolve its public transport contract without rewriting the underlying application or service orchestration.

### Canonical Transport Model

The architecture now distinguishes between:

- legacy public routes, which remain backward-compatible and non-canonical
- canonical versioned routes, which will live under `/api/v1/...`

This distinction belongs to the HTTP contract layer only. It does not change the execution model introduced in Phase 0.7.

### Layering Rule

Versioning is defined as:

- HTTP route exposure and contract selection → transport concern
- use-case orchestration → application concern
- domain execution and persistence → service/store concern

This preserves the application-layer boundary while making room for controlled public API evolution.

### Architectural Benefit

With this model in place, the backend can:

- preserve current consumers on legacy routes
- freeze the current authenticated surface under `v1`
- prepare authorization and later domain growth on top of explicit contract semantics
- avoid mixing route-evolution decisions with business-flow implementation details

### Frontend Alignment Constraint

This architectural move is intentionally compatible with the current project rule that the frontend remains aligned to backend Phase 0.6 until Stage 0 completes. Versioning preparation therefore improves backend contract governance without forcing immediate frontend migration.

### Result

By the end of Phase 0.9, the architecture has a defined transport-evolution policy, a concrete router-level implementation of that policy, an explicitly frozen authenticated `v1` surface and representative transport regression coverage protecting key legacy-versus-canonical scenarios:

- Handlers → transport plus real legacy/canonical route exposure rules
- Application → orchestration
- Services / Stores → execution

Versioning is now implemented as router-level projection of the same stabilized handler/application behavior into two public route spaces rather than as duplicated route-specific logic.

The architectural consequence of 0.9.3 is that the authenticated surface itself is no longer merely reachable through canonical paths; it is now explicitly treated as the frozen `v1` contract surface for bootstrap, profile, settings, session and wallet inventory/read behavior. Legacy and canonical routes are therefore architectural aliases over the same stabilized authenticated boundary rather than separate transport contracts.

The final architectural consequence of 0.9.5 is not a new runtime mechanism but a closed architectural record: roadmap, status, testing and handoff now all describe the same versioning boundary and the same frontend Phase 0.6 alignment rule.

This prepares the backend for Phase 0.10 and beyond, where authorization and later cross-cutting concerns can be introduced on top of a stable application boundary, an explicit public contract-evolution model and a real canonical route surface already present in the transport layer.

## Phase 0.10.1 — Authorization Model Definition

Phase 0.10.1 introduces the first explicit authorization architecture element without yet changing the request pipeline. The architectural goal is to ensure the backend stops treating authorization as a future handler concern and instead defines a dedicated core model before request-context propagation and enforcement are added.

The new `internal/core/authorization` package establishes the foundational vocabulary of the authorization layer:

- `Role` as the coarse-grained actor classification
- `Permission` as the explicit capability vocabulary
- a static role → permission mapping as the initial policy foundation
- `AuthorizationSubject` as the normalized authorization-facing representation of the authenticated actor

This keeps the architecture layered correctly:

- authentication still resolves identity
- authorization now has an explicit model
- middleware/context integration remains deferred to 0.10.2
- policy evaluation remains deferred to 0.10.3
- endpoint enforcement remains deferred to 0.10.4

The key architectural consequence is that the backend no longer needs to jump directly from authenticated claims to endpoint decisions. There is now a dedicated model boundary between identity and future permission checks, introduced without mutating the transport contract or the current application/use-case composition.



## Phase 0.10.2 — Authorization Context & Middleware

Phase 0.10.2 moves authorization from a purely static model into the request pipeline without yet turning it into enforcement. Architecturally, this is the moment where the backend starts carrying two distinct actor representations during authenticated requests:

- authentication claims as the transport/security identity artifact
- `AuthorizationSubject` as the normalized authorization artifact

This separation matters because later policy evaluation should not need to depend directly on raw JWT claims or transport-specific details. The new middleware layer projects claims into authorization context once, near the transport boundary, and the rest of the system can later consume authorization state through a dedicated core abstraction.

The protected-route pipeline is therefore now effectively structured as:

`Request → RequireAuth → HydrateAuthorization → Handler/Application`

At this stage, the authorization middleware is intentionally non-blocking. Its role is to hydrate context, not to decide access. That preserves current runtime stability while preparing a clean handoff to centralized policy evaluation in 0.10.3.


## Phase 0.10.3 — Policy Evaluation Layer

Phase 0.10.3 introduces the first centralized decision point of the authorization architecture without yet turning authorization into transport enforcement. Up to 0.10.2, the backend had a model and a hydrated authorization subject in request context, but still lacked a stable policy boundary.

The new policy layer closes that gap by defining:

- explicit `Action` and `Resource` vocabularies as authorization-facing intent
- centralized action/resource → permission projection derived from the static role-permission model
- a `PolicyEvaluator` that answers authorization questions against a normalized `AuthorizationSubject`

Architecturally, this matters because endpoint enforcement should not be the first place where policy semantics appear. The evaluator now becomes the boundary between authorization context and transport/application enforcement:

`Request → RequireAuth → HydrateAuthorization → PolicyEvaluator → Handler/Application`

At this stage, the evaluator is present but still not invoked to deny requests in router/handler flows. This keeps Phase 0.10.3 non-breaking while ensuring 0.10.4 can enforce permissions progressively through a centralized policy API instead of duplicating role/permission logic near endpoints.


## Phase 0.10.4 — Endpoint-Level Enforcement

Phase 0.10.4 is the first point where authorization becomes operational rather than purely preparatory. Up to 0.10.3, the backend could model subjects, hydrate authorization context and answer policy questions, but no protected endpoint actually denied requests because of those answers.

This subphase introduces a dedicated route-level enforcement middleware in `internal/core/httpx` that sits after authentication and authorization-context hydration, and before the selected handler execution boundary:

`Request → RequireAuth → HydrateAuthorization → RequirePermission → Handler/Application`

Architecturally, this is important because the transport layer still does not embed role logic directly. Router configuration declares the required action/resource pair, and the middleware delegates the actual decision to the centralized policy layer. That keeps enforcement thin, testable and consistent with the static permission model already defined in 0.10.1 and evaluated in 0.10.3.

The enforcement rollout is intentionally progressive. Only the authenticated endpoints already aligned with the current static permission map are enforced in 0.10.4, while endpoints that still require richer self/ownership semantics remain outside the first denial boundary. This preserves Stage 0 behavioral stability while proving the end-to-end authorization path in production code.


## Phase 0.10.5 — Documentation & Contract Consolidation

Phase 0.10.5 closes the architectural introduction of authorization by consolidating how the repository describes the new boundary. No new middleware, policy logic or enforcement path is added here. Instead, the subphase confirms that the architecture now contains all four delivered runtime layers of authorization:

- model vocabulary
- request-context hydration
- centralized policy evaluation
- progressive endpoint-level enforcement

The architectural value of 0.10.5 is that the code and the trunk documentation now match. The backend can move forward with later module work without carrying conflicting descriptions of whether authorization exists only as preparatory infrastructure or already participates in selected authenticated request decisions.


## Phase 0.11 — Domain Module Pattern

Phase 0.11 is the current architectural consolidation step after the completed authorization layer. It does not add a new public runtime capability on its own; instead, it standardizes how the current Stage 0 domain-facing modules are structured internally so later work can grow on top of a clearer module model.

The architectural target is a consistent internal module organization of the form:

```text
internal/modules/<module>/
    http/
    app/
    domain/
    repository/   (when it adds real clarity)
```

This pattern makes the current layer intent explicit:

- `http` owns request/response transport and DTO mapping
- `app` owns use-case orchestration
- `domain` owns models, invariants and module-local contracts
- `repository` owns persistence-facing implementation boundaries when the module really needs them

Architecturally, the important consequence is that the backend stops treating module organization as mostly historical or handler-driven. Instead, the current domain-facing modules are expected to follow an explicit direction of dependency:

`HTTP → APP → DOMAIN`

with repository implementations supporting the application/domain boundary rather than collapsing it.

The modules in scope remain `auth`, `user` and `usersettings`. They do not all have identical semantics, but they align to the same structural pattern:

- `auth` keeps ownership of authentication flows and authenticated-identity entry behavior
- `user` keeps ownership of the user entity/profile boundary
- `usersettings` keeps ownership of configuration and preference semantics

This matters because later Stage 0 and post-Stage 0 work should not continue to grow through implicit ownership transfer or direct concrete cross-module knowledge. Phase 0.11 therefore establishes the architectural expectation that coordination between these modules should happen through minimal explicit contracts rather than through accidental import-level coupling.

Phase 0.11 remains intentionally non-breaking. It preserves:

- the completed authenticated and authorization-enabled transport surface
- the canonical `/api/v1/...` contract state
- the existing payloads and error model
- the already stabilized challenge / verify / bootstrap semantics

### Current Repository State

0.11.1 is completed and has already established the architectural definition of the pattern.

0.11.2 is also completed and has already applied the pattern to `internal/modules/user`, making the `user` module the first concrete runtime reference implementation of the new structure.

0.11.3 is now completed as well and applies the same structural discipline to `internal/modules/usersettings`, preserving its settings-specific semantics while exposing explicit `app`, `domain` and `repository` boundaries.

0.11.4 is now completed for `internal/modules/auth`. The auth module is aligned conservatively rather than rewritten:

- `auth/app` owns application/runtime orchestration for auth, session, bootstrap and authenticated wallet-management use cases
- `auth/domain` owns canonical wallet contracts and base wallet types
- root `auth` preserves compatibility, transport and intentionally deferred runtime-store surfaces
- `auth/repository` is documented as a transitional façade, not the active implementation owner yet
- wallet bootstrap auth remains public and stable, while authenticated wallet-management handlers are narrowed to transport/delegation

Phase 0.11 is therefore **completed**. Its final work consolidated cross-module contracts and closed the phase documentation without changing public routes, payloads, authentication semantics or settings behavior.

### 0.11.5 Contract Consolidation Result

0.11.5 formalizes and applies the expectation that `auth`, `user` and `usersettings` coordinate through explicit, minimal contracts rather than through accidental knowledge of each other's implementation details. The completed internal steps were dependency mapping, contract extraction, interface alignment and runtime compatibility validation. The result is that `auth/app` depends on `UserProvider` and `UserSettingsProvider` contracts while root `auth` preserves compatibility wiring. This remains a structural result only: it does not change public routes, payloads, authentication semantics or settings behavior.

### 0.11.6 Phase Closure Result

0.11.6 closes the documentation state for the phase. The roadmap, phase-status, README, handoff and architecture documents now record Phase 0.11 as completed and preserve the boundary decisions needed for the next roadmap phase.

---

## Phase 0.12 — Read / Write Model Separation

Phase 0.12 extends the completed Domain Module Pattern by clarifying model responsibility inside each module.

The architectural rule introduced by this phase is that model shape must reflect intent:

- read models are used for output, views and response-oriented application data
- write models are used for input, commands and mutation-oriented application data
- domain models remain owned by the module domain boundary
- mapping functions make transformations explicit

This is not a full CQRS adoption. The backend is not introducing separate read/write persistence stores, event sourcing or asynchronous projections. The separation is structural and internal, intended to reduce ambiguity in the current codebase while preserving the public HTTP contract.

### Relationship To 0.11

0.11 answered where module responsibilities live: `http`, `app`, `domain` and `repository` where useful.

0.12 answers what model responsibility a structure has once it crosses those boundaries.

The expected direction is:

```text
HTTP request DTO -> Write Model -> Domain / Application behavior
Domain / Application result -> Read Model -> HTTP response DTO
```

The exact extraction points must be determined by 0.12.1 against the real code, not by assumption.

### Compatibility Rule

0.12 must not change public routes, public payload semantics, authentication behavior, authorization behavior, error envelopes or API versioning. Any model split must be backed by explicit mapping so current behavior remains stable.

### 0.12.0 Result

0.12.0 establishes this architectural definition and locks the phase documentation before code changes begin.


### 0.12.1 Audit Definition Lock

0.12.1 is subdivided because the audit must be exhaustive before read/write extraction begins. The first internal step, 0.12.1.0, is documentation-only and records the audit method, classification vocabulary and sub-subphase order.

The audit introduces a formal model lifecycle classification layer for review purposes:

- READ
- WRITE
- DOMAIN
- CONTRACT
- INFRASTRUCTURE
- HYBRID / TRANSITIONAL

This does not change runtime architecture by itself. It creates the documentary and analytical basis for later code changes in 0.12.2 through 0.12.5.

### 0.12.1 Model Classification Audit Result

0.12.1 is completed as an audit-only architectural step. The repository now has a complete model inventory, explicit classification results, cross-layer usage analysis, risk mapping and target separation direction.

The audit establishes the next structural rule for implementation work:

```text
HYBRID / TRANSITIONAL -> READ + WRITE + DOMAIN where required
```

The implementation of that target is deferred to 0.12.2 and later. Public HTTP behavior, authentication, authorization, error envelopes and API versioning remain unchanged.

### 0.12.2 Read Model Extraction Definition Lock

0.12.2 begins the first implementation-oriented part of Phase 0.12, but 0.12.2.0 remains documentation-only. It locks the read model extraction sequence before any Go structs, mappings or handlers are changed.

The extraction target is limited to read-oriented output/view/response shapes. The public HTTP contract must remain compatible; therefore read models may be introduced internally only when their mapping preserves the existing externally observable response semantics.

0.12.2 must use the completed 0.12.1 audit artifacts as evidence. Hybrid/transitional models are not split by assumption; each read model candidate must be tied back to the inventory, classification, cross-layer analysis, risk map and target separation definition.


### 0.12.2 Read Model Extraction Closure

0.12.2 completes the read-side extraction path. The repository now contains explicit read model packages under the relevant modules, mapping functions that transform module-owned data into read models and response alignment that keeps current public JSON contracts stable.

The read model boundary is internal. It improves module clarity without changing routes, API versioning, authentication behavior, authorization behavior or standardized error envelopes. Write model isolation is intentionally deferred to 0.12.3.

### 0.12.3 Write Model Isolation Definition Lock

0.12.3 begins the write-side counterpart to the completed read model extraction. The write boundary is limited to input, command and mutation-oriented structures. It must not reuse read models as accepted input shapes and must not change public request payload semantics.

The intended direction is:

```text
HTTP request DTO -> Write Model -> Domain / Application behavior
```

0.12.3 is completed. Write model packages now represent mutation/input intent separately from read projections, and handler alignment uses write-side mapping while preserving the public request contract.

### 0.12.3 Write Model Isolation Closure

The write-side boundary is now explicit for the current authentication-oriented flows. Public handlers continue to accept the same JSON payload semantics, while internal processing can route through write models and domain write inputs. This keeps the read model boundary output-only and prevents read structures from becoming mutation inputs.

### 0.12.4 Mapping Layer Introduction Definition Lock

0.12.4 introduces a centralized mapping layer as an internal architectural boundary. The target location is `internal/modules/<module>/mappers/`, separate from `readmodels`, `writemodels`, `domain` and `app`.

The layer owns explicit transformations in both directions already introduced by 0.12.2 and 0.12.3:

- Write → Domain
- Domain/Application → Read

The definition lock does not move code yet. It records that future implementation must preserve public HTTP request and response contracts while reducing mapper ownership spread across read model, write model and application packages.


### 0.12.4 Mapping Layer Closure

Phase 0.12.4 completed the transition from distributed mapper ownership to module-local mapper packages under `internal/modules/<module>/mappers/`. The architecture now treats mapping as an explicit boundary between write models, domain inputs, domain state and read models. Public HTTP contracts remain unchanged.

## Phase 0.12.5 Contract Alignment Boundary

Phase 0.12.5 aligns internal provider contracts with the model separation introduced during Phase 0.12. The alignment applies to provider-facing boundaries and must not change public HTTP routes, request payloads, response payloads or API versioning.

The expected direction is:

- provider contracts remain explicit and stable
- read-side contract outputs use read models or read-oriented views intentionally
- write-side contract inputs use write models or domain input objects intentionally
- mapper ownership remains centralized under module-local `mappers` packages
- application services avoid implicit transformation ownership

## Phase 0.12.5 Contract Alignment Result

Phase 0.12.5 completed the alignment between HTTP handler contracts, internal write models, domain inputs, read models and centralized mapping ownership. Public HTTP routes, request payload semantics, response payload semantics and API versioning remain unchanged.

## Phase 0.13 Provider Layer Consolidation Boundary

Phase 0.13 follows the completed Read / Write Model Separation phase.

The architectural target is to consolidate provider boundaries as the explicit entry point between transport/application orchestration and domain services.

The intended direction is:

```text
HTTP → Provider → Application → Domain → Repository
```

This phase must preserve the public HTTP/API surface and must not introduce business behavior changes. Provider consolidation is an internal boundary clarification built on top of the Phase 0.12 mapper and contract alignment results.
