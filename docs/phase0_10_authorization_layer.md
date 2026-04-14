# Phase 0.10 — Authorization Layer

## Objective

Introduce a structured authorization layer on top of the already-stabilized authentication, application and versioning foundations without breaking current runtime behavior or the canonical `v1` transport contract.

## Initial Context

By the end of Phase 0.9.5, the backend already had:

- stabilized authentication flows
- authenticated application/use-case boundaries
- standardized error handling
- canonical `/api/v1/...` routing and `v1` contract freezing
- representative version-aware transport regression protection

What it still lacked was a first-class authorization model. Authenticated users were recognized, but the backend did not yet have a dedicated internal language for permissions or differentiated roles.

## Problem Statement

Without an explicit authorization layer:

- all authenticated users are effectively treated the same
- permission semantics risk emerging ad hoc inside handlers or services
- future endpoint protection would have no stable role/permission vocabulary
- later architectural growth would mix identity, policy and enforcement concerns

Phase 0.10 exists to solve that in a layered, non-breaking way.

## Scope

Phase 0.10 covers:

- authorization model definition
- authorization context propagation
- centralized policy evaluation foundations
- progressive endpoint-level enforcement
- documentation consolidation around the new authorization layer

Phase 0.10 does not yet cover:

- complex RBAC administration surfaces
- tenant-aware authorization
- database-managed permission assignment
- external policy engines
- frontend authorization behavior

## Architectural Intent

Authorization must remain a layer distinct from authentication:

- authentication answers **who the actor is**
- authorization answers **what the actor may do**

The intended architectural progression is:

`HTTP → Auth → Authorization → Application → Domain`

This phase therefore starts by defining the authorization primitives before middleware wiring or endpoint blocking is attempted.

## Subphases

### 0.10.1 — Authorization Model Definition ✔ Completed

This subphase introduces the foundational authorization model under `internal/core/authorization`:

- initial roles: `user`, `admin`
- initial permissions for current authenticated Stage 0 resources
- static role → permission mapping helpers
- `AuthorizationSubject` as a normalized authorization-facing actor model
- focused unit tests validating normalization, immutability and permission aggregation behavior

This subphase intentionally does **not** yet:

- attach authorization state to request context
- evaluate policies at runtime
- deny requests based on permissions
- persist roles or permissions in storage

Its purpose is to establish a stable internal vocabulary so later authorization work is layered rather than ad hoc.

### 0.10.2 — Authorization Context & Middleware ✔ Completed

Phase 0.10.2 propagates authorization subject information through the authenticated request lifecycle without changing payload or enforcement behavior. The current auth context and HTTP middleware stack now hydrate an explicit authorization subject for authenticated requests while still deferring policy decisions and endpoint denial to later subphases.

### 0.10.3 — Policy Evaluation Layer ✔ Completed

This subphase introduces the centralized policy boundary that handlers and application code can later use to ask authorization questions through a stable interface rather than embedding role/permission decisions locally.

### 0.10.4 — Endpoint-Level Enforcement ✔ Completed

Once policy evaluation exists, selected endpoints should begin enforcing permissions progressively while preserving current contract discipline and minimizing risk to stabilized Stage 0 flows.

### 0.10.5 — Documentation & Contract Consolidation ✔ Completed

The final subphase should align README, roadmap, architecture, testing, status and handoff documents so the delivered authorization layer is represented consistently across the trunk documentation set.

## Current Result After 0.10.1

Phase 0.10 is now in progress. The backend still behaves exactly as it did after 0.9.5 from a client/runtime perspective, but it now includes an explicit authorization model at the core layer.

That means the project no longer needs to invent authorization semantics later inside transport code or feature handlers. The next work can proceed through context propagation, policy evaluation and endpoint enforcement on top of a real model foundation.


## Current Result After 0.10.2

The backend now has both halves of the authorization foundation in place:

- a static authorization model (`Role`, `Permission`, role mapping and `AuthorizationSubject`)
- request-context propagation of a normalized authorization subject for authenticated HTTP flows

Phase 0.10.2 specifically delivers:

- context helpers for storing and retrieving authorization subject data
- deterministic projection from JWT claims to authorization subject
- a dedicated authorization-hydration middleware in `internal/core/httpx`
- integration of that middleware into protected route registration for both legacy and canonical route spaces
- focused tests protecting the new context and middleware behavior

The runtime contract remains stable: clients still see the same success and error shapes, and the backend still does not enforce endpoint permissions. What changed is the internal request lifecycle, which now carries authorization data in a form suitable for centralized policy evaluation.

### Next

0.10.4 — Endpoint-Level Enforcement


## Current Result After 0.10.3

The backend now has a centralized policy evaluation layer on top of the authorization model and authorization-context propagation already delivered in the previous subphases.

Phase 0.10.3 specifically delivers:

- explicit `Action` and `Resource` vocabularies for authorization questions
- static action/resource → permission projection helpers
- permission evaluation against normalized `AuthorizationSubject` data
- a `PolicyEvaluator` with `Evaluate(...)` and `Can(...)` entry points
- focused tests covering allowed, denied and unknown-policy decisions

The runtime contract still remains stable: endpoint behavior, success payloads and standardized errors are unchanged, and no request is denied yet because of authorization policy. What changed is that the backend now has a real internal policy boundary instead of having to invent one at the moment enforcement begins.

### Next

0.10.4 — Endpoint-Level Enforcement


## Current Result After 0.10.4

The backend now has the first real endpoint-level authorization enforcement slice on top of the model, context-propagation and policy-evaluation work already delivered in 0.10.1 through 0.10.3.

Phase 0.10.4 specifically delivers:

- route-level permission enforcement middleware under `internal/core/httpx`
- centralized denial decisions powered by `PolicyEvaluator`
- standardized `AUTH_FORBIDDEN` responses for authorization failures
- progressive enforcement on the authenticated endpoints already aligned with the current static permission map
- legacy and canonical route coverage for the first enforced `me` and `settings` surfaces

The enforced scope currently includes:

- `GET /auth/me` and `GET /api/v1/auth/me`
- `GET /auth/me/settings` and `GET /api/v1/auth/me/settings`
- `PATCH /auth/me/settings` and `PATCH /api/v1/auth/me/settings`

`PATCH /auth/me` remains intentionally outside the first enforcement slice because its self-update semantics are not yet fully represented by the current static role model.

The runtime contract remains disciplined: success payloads and legacy/canonical route compatibility stay intact for permitted users, while selected endpoints now deny unauthorized access through the standardized authorization error path.

## Current Result After 0.10.5

Phase 0.10 is now fully consolidated. No new runtime authorization behavior is added in 0.10.5; instead, the repository documents the delivered state consistently across the trunk set and the dedicated authorization narrative.

Phase 0.10.5 specifically delivers:

- completion-state alignment across README, roadmap, index, handoff and phase-status documents
- removal of stale intermediate descriptions that still treated authorization as pre-enforcement after 0.10.4
- explicit contract documentation for the progressive first enforcement slice and the standardized forbidden path
- a single coherent explanation of how model, context, policy and selected endpoint enforcement now fit together

Phase 0.10 is completed with the following delivered progression:

- 0.10.1 → authorization model definition
- 0.10.2 → authorization context and middleware hydration
- 0.10.3 → centralized policy evaluation
- 0.10.4 → progressive endpoint-level enforcement
- 0.10.5 → documentation and contract consolidation

The next phase can now build on a closed and consistently described authorization layer rather than on partially synchronized implementation notes.

### Next

0.11 — Domain Module Pattern
