# Phase 0.9 — API Versioning Strategy

## Objective

Define the official API versioning model for the SCAVO Exchange Backend after the completion of Phase 0.8, preserving the current public contract while establishing a canonical route space for future evolution.

## Initial Context

By the time Phase 0.8 closed, the backend already had:

- a stable authenticated HTTP surface
- a standardized JSON error envelope
- centralized app-error typing and factories
- contract-hardening tests for representative transport error behavior

At the same time, the backend still exposed those contracts through unversioned routes such as `/auth/login`, `/auth/bootstrap`, `/auth/me`, `/auth/session`, `/auth/me/settings` and `/auth/wallets`.

This created an implicit public API version with no explicit policy for compatibility, deprecation or future breaking changes.

## Problem Statement

Without a defined API versioning strategy, the backend would face the next Stage 0 phases with unclear public-contract rules:

- new canonical behavior would have no formal route namespace
- legacy compatibility would remain implicit
- breaking changes could be introduced accidentally
- authorization work in Phase 0.10 could mix transport evolution with permission semantics
- frontend/backend alignment would become harder to reason about while the frontend intentionally remains aligned only up to backend Phase 0.6 until Stage 0 is complete

## Scope

Phase 0.9 covers:

- API versioning policy definition
- canonical vs legacy route classification
- contract evolution rules
- preparation for canonical router exposure
- version-aware testing direction
- trunk-documentation alignment

Phase 0.9 does not yet cover:

- authorization or permission enforcement
- domain expansion
- business-flow redesign
- `v2` route publication
- removal of the current route surface

## Canonical Versioning Strategy

The backend adopts path-based versioning as its canonical public-contract model.

Canonical namespace:

```text
/api/v1/...
```

This namespace is defined as the authoritative route space for the currently stabilized authenticated surface.

## Legacy Compatibility Rule

The current unversioned routes remain available and supported as:

- legacy
- backward-compatible
- non-canonical

Examples include:

- `/auth/login`
- `/auth/bootstrap`
- `/auth/me`
- `/auth/me/settings`
- `/auth/session`
- `/auth/wallets`

These routes are not immediately deprecated out of existence. They remain part of the compatibility story while the canonical `v1` surface is introduced.

## Contract Evolution Rules

Within `v1`, the following changes are considered compatible:

- additive optional fields
- additive metadata that does not alter meaning
- additive route exposure that preserves existing semantics

The following require a new API version:

- removing fields
- changing field types
- changing the structure of success payloads
- changing the structure or semantics of the standardized error envelope
- changing route behavior in a way that alters client expectations

## Error Contract Stability

The structured error envelope finalized in Phase 0.8 is now considered part of the `v1` public API contract.

That means:

- it is frozen for `v1`
- it must remain equivalent across legacy and canonical route exposure
- later breaking changes require versioned evolution rather than silent mutation

## Frontend Alignment Constraint

The project currently keeps the frontend aligned to backend Phase 0.6 while backend Stage 0 continues to mature. Phase 0.9 must respect that project rule.

This means:

- versioning preparation must not force immediate frontend migration
- legacy routes must remain stable during the Stage 0 completion path
- canonical versioning exists to prepare safe future alignment, not to invalidate the current frontend development track

## Planned Subphases

### 0.9.1 — Versioning Policy Definition ✔ Completed

The canonical versioning model, compatibility rules and Stage 0 alignment constraints are now documented.

### 0.9.2 — Router Versioning Foundation ✔ Completed

The canonical route policy is now materialized in the real HTTP router.

Delivered in this subphase:

- canonical `/api/v1/...` auth and authenticated route exposure added in `internal/core/httpx/router.go`
- current legacy `/auth/...` routes preserved without semantic change
- shared route-registration logic introduced so versioning does not duplicate handler/business behavior
- current auth middleware protection preserved across both route spaces

This subphase intentionally changes route exposure only. It does not yet redefine or freeze the semantic equivalence set endpoint by endpoint; that remains the responsibility of 0.9.3 and 0.9.4.

### 0.9.3 — Authenticated Surface Version Freezing ✔ Completed

The current authenticated surface is now explicitly frozen as the canonical `v1` authenticated contract.

Delivered in this subphase:

- bound bootstrap, profile, settings, session and wallet inventory/read behavior to canonical `v1` semantics
- clarified that legacy `/auth/...` and canonical `/api/v1/auth/...` entry paths are transport projections of the same authenticated contract
- preserved current business behavior, payload semantics, middleware wrapping and Phase 0.8 error-envelope behavior
- left strict route-by-route transport equivalence assertions to 0.9.4 rather than mixing freeze definition with contract-test expansion

### 0.9.4 — Version-aware Contract Testing ✔ Completed

The versioning model is now backed by explicit transport-aware regression coverage so the frozen `v1` authenticated surface cannot drift silently across legacy and canonical entry paths.

Delivered in this subphase:

- added router-level tests that verify helper path composition for legacy and canonical route registration
- added protected-endpoint equivalence checks confirming both `/auth/...` and `/api/v1/auth/...` return the same standardized missing-bearer contract
- added protected-endpoint equivalence checks confirming both route spaces return the same standardized unauthorized contract for invalid tokens
- added representative public success-shape checks around wallet challenge creation to ensure canonical and legacy paths continue exposing compatible payload structure

This subphase intentionally protects representative route behavior and contract shape without expanding into a new API version or redefining the already frozen `v1` payload semantics.

### 0.9.5 — Documentation Consolidation ✔ Completed

The trunk documentation set is now aligned with the finalized versioning model, the authenticated-surface `v1` freeze and the representative transport-level parity checks introduced in 0.9.4.

Delivered in this subphase:

- updated the trunk documentation set so Phase 0.9 is consistently represented as completed rather than partially pending
- aligned roadmap, phase-status, README, architecture, testing and handoff around the same explanation of canonical `/api/v1/auth/...` and legacy `/auth/...`
- corrected state drift in the operational handoff document so the latest completed subphase and next planned phase are now accurate
- preserved the project rule that the frontend remains aligned to backend Phase 0.6 until Stage 0 closure, avoiding any suggestion that Phase 0.9 forces immediate frontend route migration

This subphase intentionally changes documentation only. It closes the phase by removing narrative drift, not by altering router behavior, payload semantics, middleware or test scope.

## Expected Outcome

At the end of Phase 0.9, the backend should have:

- an explicit API versioning policy
- a canonical `v1` route space
- preserved legacy route compatibility
- version-aware regression protection
- documentation that makes transport evolution rules explicit before Phase 0.10 begins

## Phase Closure

Phase 0.9 is now fully completed. The backend has an explicit versioning policy, a real canonical `v1` route surface, a frozen authenticated `v1` contract, representative legacy-versus-canonical transport regression protection and a consolidated trunk documentation set that records the same state consistently before Phase 0.10 begins.
