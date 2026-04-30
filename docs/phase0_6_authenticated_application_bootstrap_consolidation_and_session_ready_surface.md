# Phase 0.6 — Authenticated Application Bootstrap Consolidation & Session-Ready Surface

---

## Objective

Consolidate the authenticated application surface introduced during Phase 0.5 into a clearer bootstrap-oriented boundary that can be consumed by the frontend without semantic ambiguity.

The purpose of Phase 0.6 is not to introduce new business domains.

The purpose is to clarify and harden how the backend exposes the already existing authenticated surface through the endpoints that currently define the application bootstrap state:

* `GET /auth/me`
* `GET /auth/session`
* `GET /auth/me/settings`
* `GET /auth/wallets`

This phase exists to ensure that authenticated application startup can rely on a coherent separation between:

* durable authenticated identity
* token/session context
* user settings
* wallet inventory

without collapsing those concerns into a single overloaded endpoint.

---

## Initial Context

By the end of Phase 0.5.5.3, the backend already provides a stable authenticated surface composed of:

* session-aware authentication
* wallet-backed user identity
* authenticated profile retrieval
* user settings read and mutation
* wallet ownership inventory
* hardened settings contract behavior
* explicit tests around settings and authenticated handlers

In practical terms, the system already exposes the following authenticated endpoints:

* `GET /auth/me`
* `PATCH /auth/me`
* `GET /auth/session`
* `GET /auth/me/settings`
* `PATCH /auth/me/settings`
* `GET /auth/wallets`

The problem at this stage is not endpoint absence.

The problem is that the authenticated application startup surface is now broad enough that semantic overlap can begin to appear unless boundaries are made explicit.

The frontend and any future consumer reaching authenticated state must be able to distinguish:

* which endpoint describes who the authenticated user is
* which endpoint describes the active session and token-derived metadata
* which endpoint describes persistent user preferences
* which endpoint exposes wallet ownership detail

Phase 0.6 begins from that exact need.

---

## Problem Statement

After Phase 0.5, the backend surface is functionally useful but semantically close enough that a consumer could blur the responsibilities of `/auth/me` and `/auth/session`.

That ambiguity creates several medium-term risks:

* frontend bootstrap logic may overuse one endpoint for responsibilities owned by another
* token metadata may be interpreted as durable user profile state
* authenticated identity may be interpreted as session state
* settings and wallet inventory may be folded into the wrong bootstrap assumptions
* future contract alignment may become harder because consumers will already rely on accidental overlap instead of deliberate responsibility

The issue is subtle but real.

`/auth/me` and `/auth/session` are both authenticated endpoints and they both expose identity-adjacent information.

However, they do not exist for the same purpose.

Without explicit clarification, the system risks evolving into a surface where:

* session fields leak into profile semantics
* profile semantics become coupled to token claims
* frontend bootstrap becomes contract-fragile
* future changes must preserve accidental overlap instead of intentional design

Phase 0.6 addresses that risk incrementally.

---

## Scope

Phase 0.6 includes:

* authenticated bootstrap boundary clarification
* role separation between `/auth/me` and `/auth/session`
* clarification of the place of `/auth/me/settings` within authenticated bootstrap
* clarification of the place of `/auth/wallets` within authenticated bootstrap
* contract-oriented validation through tests
* cumulative documentation updates that keep the project aligned with the real backend state

Phase 0.6 explicitly excludes:

* new business endpoints
* billing, payments, trading, or exchange flows
* redesign of authentication architecture
* redesign of wallet lifecycle semantics
* introduction of a new bootstrap endpoint
* typed settings redesign
* persistence schema changes
* DB migrations
* breaking contract changes

This is a consolidation phase.

It is not a feature expansion phase.

---

## Root Cause Analysis

Phase 0.5 correctly prioritized delivery of the authenticated surface itself.

That was the right choice.

Before boundary refinement could even be useful, the backend needed to have the underlying surface available:

* session support
* user identity retrieval
* settings access
* wallet inventory access

Only after those pieces were already present did the next architectural problem become visible.

That problem is not lack of functionality.

That problem is contract meaning.

The backend now exposes enough authenticated information that the distinction between different kinds of authenticated state must be made explicit.

The root cause of the problem is therefore natural phase progression:

1. Phase 0.5 built the authenticated surface.
2. That success created adjacent endpoints with partially related information.
3. Once adjacent surfaces exist, their semantic boundaries must be clarified before future contract alignment and frontend bootstrap standardization can proceed safely.

This makes Phase 0.6 a continuation of stabilization, not a redesign.

---

## Files Affected

Code and test impact introduced for the first implemented subphase in Phase 0.6:

* `internal/modules/auth/http_handlers_test.go`

Documentation updates required to keep the project synchronized with the implemented state:

* `README.md`
* `docs/index.md`
* `docs/phase-status.md`
* `docs/handoff/backend-status.md`
* `docs/testing.md`
* `docs/flows.md`
* `docs/phase0_6_authenticated_application_bootstrap_consolidation_and_session_ready_surface.md`

No production handler or service code was changed in the implemented portion of this phase.

---

## Authenticated Bootstrap Surface Model

Phase 0.6 formalizes the authenticated application startup surface as four complementary but distinct read surfaces.

### 1. Session Surface

Endpoint:

* `GET /auth/session`

Purpose:

Expose token-derived and session-scoped information describing the currently authenticated request context.

This includes fields such as:

* authenticated state
* token type
* subject
* issuer
* expiration metadata
* current authenticated user summary

This surface is about the active authenticated session.

It is not the main durable profile bootstrap surface.

---

### 2. Bootstrap Identity Surface

Endpoint:

* `GET /auth/me`

Purpose:

Expose the authenticated user's durable identity view and profile-oriented bootstrap state.

This includes:

* user identity summary
* authenticated profile data
* wallet-linked snapshot information
* profile-facing counters and wallet summary indicators

This surface is about who the authenticated user is inside the application model.

It is not the canonical token/session metadata surface.

---

### 3. Settings Surface

Endpoint:

* `GET /auth/me/settings`

Purpose:

Expose persisted user preference state required for authenticated application behavior.

This surface is dedicated to settings.

It is not a profile endpoint and it is not a session endpoint.

---

### 4. Wallet Inventory Surface

Endpoint:

* `GET /auth/wallets`

Purpose:

Expose the authenticated user's wallet inventory, including wallet lifecycle and inventory-oriented metadata already introduced by earlier phases.

This surface is not a substitute for `/auth/me`.

It is the detailed wallet list surface.

---

## Surface Separation Principle

The key principle introduced by Phase 0.6 is that authenticated application bootstrap is a composed concept, not a single overloaded payload.

The backend therefore keeps separate endpoints for separate responsibilities.

That separation is intentional.

The frontend bootstrap sequence may combine them, but the backend contract should not erase their meaning.

The goal is clarity, not endpoint collapse.

---

## Phase 0.6.1 — Bootstrap Surface Boundary Clarification

---

## Objective

Clarify the exact responsibility of each authenticated bootstrap-related endpoint and formally establish the boundary between `/auth/me` and `/auth/session`.

This subphase ensures that the backend has a documented and tested separation between:

* session context
* authenticated identity bootstrap
* settings bootstrap
* wallet inventory bootstrap

---

## Initial Context

Before 0.6.1, the system already exposed both `/auth/me` and `/auth/session` successfully.

Both endpoints were valid and functional.

However, because both endpoints were authenticated and both exposed user-adjacent information, the semantic distinction between them still depended too much on implicit reading of the code and tests.

That was sufficient for internal development continuity.

It was not sufficient for a clean bootstrap contract phase.

---

## Problem Statement

The backend needed an explicit contract statement establishing that `/auth/me` and `/auth/session` do not represent the same layer of application meaning.

Without that clarification, future consumers could make assumptions such as:

* reading token expiration from `/auth/me`
* treating `/auth/session` as the primary profile surface
* expecting wallet inventory summary in `/auth/session`
* using `/auth/me` as a session metadata endpoint

Even if the code did not currently behave that way, the lack of explicit contract boundary left room for drift.

That drift is exactly what this subphase prevents.

---

## Scope

Phase 0.6.1 introduces:

* explicit semantic clarification of bootstrap surface roles
* boundary-oriented tests for `/auth/me`
* boundary-oriented tests for `/auth/session`
* documentation updates reflecting those semantics

Phase 0.6.1 explicitly excludes:

* public payload redesign
* handler logic changes
* service changes
* repository changes
* routing changes
* persistence changes
* session architecture redesign

---

## Implementation Summary

The implemented change for 0.6.1 was intentionally minimal and safe.

Only tests were added in order to freeze the intended meaning of the existing authenticated surfaces without changing the public contract.

The new tests were added to:

* `internal/modules/auth/http_handlers_test.go`

The specific tests introduced are:

* `TestHTTPHandlers_Me_SurfaceBoundary`
* `TestHTTPHandlers_Session_SurfaceBoundary`

These tests establish that:

### `/auth/me`

Must expose:

* `user`
* `profile`

Must not expose:

* `session`

The `profile` object must remain profile/bootstrap-oriented and must not become a session metadata container.

The tests explicitly protect against the accidental appearance of session-oriented fields inside the profile surface, including:

* `authenticated`
* `token_type`
* `issuer`
* `subject`
* `expires_at`

At the same time, they confirm that the profile surface continues to expose its expected profile and wallet-summary fields.

### `/auth/session`

Must expose:

* `session`

Must not expose:

* `profile`

The `session` object must remain session-oriented and continue to provide token/session fields such as:

* `authenticated`
* `token_type`
* `user_id`
* `subject`
* `issuer`
* `expires_at`
* `user`

The tests also protect against accidental expansion of `/auth/session` into wallet inventory or profile summary semantics by rejecting fields such as:

* `wallet_count`
* `active_wallet_count`
* `detached_wallet_count`
* `wallets`
* `primary_wallet`
* `has_wallet_session`

---

## Why Tests Only Were Correct Here

0.6.1 did not require product code changes because the backend was already behaving according to the intended responsibility split.

The problem was not broken runtime behavior.

The problem was insufficiently explicit contract separation.

In that situation, the safest implementation is to freeze the boundary using tests and aligned documentation rather than modifying functioning production code without necessity.

That keeps the system stable while still making the intended contract explicit.

---

## Validation

The implemented test change was validated with:

```bash
go test ./...
```

Observed result for the relevant package:

```text
ok      github.com/e-scavo/scavo-exchange-backend/internal/modules/auth 4.591s
```

This confirms that the new semantic boundary tests integrate cleanly with the existing auth handler suite.

---

## Release Impact

Phase 0.6.1 has no public breaking change impact.

Specifically:

* no endpoint path changed
* no request payload changed
* no response payload changed
* no auth flow changed
* no wallet lifecycle rule changed
* no settings persistence rule changed

The release impact is contract clarification, not runtime redesign.

---

## Risks

This subphase intentionally keeps risk low, but a few considerations remain important:

* future handler changes must continue respecting the now-explicit boundary
* future frontend work must consume `/auth/me` and `/auth/session` according to their clarified roles
* later subphases in 0.6 must not accidentally collapse the distinction that 0.6.1 just formalized

The main risk would be documentation drift if later work changed handlers without updating these clarified semantics.

That is why cumulative documentation alignment is part of the phase deliverable.

---

## What It Does Not Solve

Phase 0.6.1 does not yet solve:

* cross-endpoint naming alignment
* shape consistency improvements between authenticated surfaces
* a unified session-ready bootstrap read model for frontend consumption
* broader cross-surface validation beyond `/auth/me` and `/auth/session`
* final closure of Phase 0

Those responsibilities belong to the remaining planned subphases:

* 0.6.2 — Authenticated Surface Contract Alignment
* 0.6.3 — Session-Ready Bootstrap Read Model
* 0.6.4 — Application Surface Consistency Hardening

---

## Continuity With Earlier Phases

Phase 0.6 is directly built on the results of:

* Phase 0.5 — User Interaction & Application Surface
* Phase 0.5.4 — User Settings Mutation
* Phase 0.5.5 — User Settings Hardening & Contract Stabilization

That continuity matters.

0.6.1 does not replace or reopen those phases.

Instead, it uses the stable authenticated surface they produced and begins the next layer of consolidation required for frontend-ready authenticated bootstrap clarity.

---

## Historical Phase 0.6 Status

Phase 0.6 is historically completed and has been superseded by later structural and architectural phases from 0.7 through 0.14.

### Completed

* 0.6.1 — Bootstrap Surface Boundary Clarification ✔
* 0.6.2 — Authenticated Surface Contract Alignment ✔
* 0.6.3 — Session-Ready Bootstrap Read Model ✔
* 0.6.4 — Application Surface Consistency Hardening ✔

---

## Conclusion

Phase 0.6 begins the transition from merely having an authenticated application surface to having a clearly structured authenticated bootstrap surface.

The first implemented subphase, 0.6.1, deliberately avoids unnecessary runtime changes and instead formalizes the distinction between session context and authenticated identity bootstrap by adding explicit test-backed boundaries.

That is the correct move at this stage because the backend already behaves correctly enough to support the intended separation.

The missing piece was contract clarity.

With 0.6.1 complete:

* `/auth/session` is explicitly protected as the session surface
* `/auth/me` is explicitly protected as the bootstrap identity surface
* `/auth/me/settings` remains the settings surface
* `/auth/wallets` remains the wallet inventory surface

This created the right base for the later subphases of Phase 0.6, where contract alignment and session-ready bootstrap composition proceeded without semantic ambiguity.

---

---

## 0.6.2 — Authenticated Surface Contract Alignment

### Objective

Align all authenticated surfaces under a shared normalized context to eliminate drift and ensure consistent contract behavior across endpoints.

### Implementation

- Introduced shared `auth_context` normalization layer
- Refactored `/auth/me` and `/auth/session` to derive common fields from a single source
- Standardized wallet context propagation (id, address, chain)
- Added cross-endpoint contract validation tests

### Guarantees

- `/auth/me` and `/auth/session` now share consistent identity and wallet context
- `/auth/me` and `/auth/wallets` maintain aligned primary wallet representation
- No breaking changes introduced

### Result

Authenticated surface is now:

- boundary-defined (0.6.1)
- contract-aligned (0.6.2)

This enables safe evolution toward a unified bootstrap read model.

### Next

0.6.3 — Session-Ready Bootstrap Read Model

---

## 0.6.3 — Session-Ready Bootstrap Read Model

### Objective

Provide a unified authenticated bootstrap read model to eliminate multi-request initialization and ensure consistent frontend state.

### Implementation

- Introduced `GET /auth/bootstrap`
- Aggregates:
  - session
  - user
  - profile
  - settings
  - wallet snapshot
- Reuses existing builders and services

### Guarantees

- No breaking changes
- No modification of existing endpoint contracts
- Pure aggregation layer
- Consistent data derived from a unified authenticated context

### Result

Authenticated surface is now:

- boundary-defined (0.6.1)
- contract-aligned (0.6.2)
- bootstrap-ready (0.6.3)

### Next

0.6.4 — Application Surface Consistency Hardening

---

## Phase Closure Note

Phase 0.6 is historically completed and has been superseded by later structural and architectural phases from 0.7 through 0.14.

All Phase 0.6 subphases must be considered closed regardless of their original draft status at the time of writing. This note exists to prevent historical drift and ensure consistency across the documentation corpus.
