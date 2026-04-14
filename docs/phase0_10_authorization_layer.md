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

### 0.10.2 — Authorization Context & Middleware ⬜ Pending

The next subphase should propagate authorization subject information through the authenticated request lifecycle. This is where the current auth context and HTTP middleware stack should gain explicit authorization-aware context wiring without yet introducing broad endpoint denial logic.

### 0.10.3 — Policy Evaluation Layer ⬜ Pending

This subphase should introduce a centralized policy boundary so handlers and application code can ask authorization questions through a stable interface rather than embedding role/permission decisions locally.

### 0.10.4 — Endpoint-Level Enforcement ⬜ Pending

Once policy evaluation exists, selected endpoints should begin enforcing permissions progressively while preserving current contract discipline and minimizing risk to stabilized Stage 0 flows.

### 0.10.5 — Documentation & Contract Consolidation ⬜ Pending

The final subphase should align README, roadmap, architecture, testing, status and handoff documents so the delivered authorization layer is represented consistently across the trunk documentation set.

## Current Result After 0.10.1

Phase 0.10 is now in progress. The backend still behaves exactly as it did after 0.9.5 from a client/runtime perspective, but it now includes an explicit authorization model at the core layer.

That means the project no longer needs to invent authorization semantics later inside transport code or feature handlers. The next work can proceed through context propagation, policy evaluation and endpoint enforcement on top of a real model foundation.
