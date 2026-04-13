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

### 0.9.2 — Router Versioning Foundation ⬜ Pending

Introduce canonical `/api/v1/...` routing while preserving legacy route access.

### 0.9.3 — Authenticated Surface Version Freezing ⬜ Pending

Map the current authenticated surfaces onto canonical `v1` semantics without changing business behavior.

### 0.9.4 — Version-aware Contract Testing ⬜ Pending

Extend contract tests to protect canonical/legacy equivalence and forbid silent drift.

### 0.9.5 — Documentation Consolidation ⬜ Pending

Align the trunk document set with the finalized versioning model and the backend/frontend evolution rule.

## Expected Outcome

At the end of Phase 0.9, the backend should have:

- an explicit API versioning policy
- a canonical `v1` route space
- preserved legacy route compatibility
- version-aware regression protection
- documentation that makes transport evolution rules explicit before Phase 0.10 begins
