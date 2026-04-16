# Phase 0.11 — Domain Module Pattern

## Objective

Standardize the internal structure of the current Stage 0 domain-facing modules so the backend can continue evolving on top of explicit transport, application and domain boundaries without changing the already stabilized public contract.

---

## Why this phase exists

After Phases 0.6 through 0.10, the backend already has:

- an authenticated application surface consolidated for frontend consumption
- an initial application-layer boundary
- a standardized error model
- a canonical `v1` route space with compatibility aliases
- a delivered authorization layer with progressive endpoint enforcement

What still remains inconsistent is the way the current modules are organized internally. In practice, that means the repository can still accumulate:

- handlers that carry orchestration responsibilities longer than they should
- module-local structures that differ from one domain-facing module to another
- implicit ownership boundaries between `auth`, `user` and `usersettings`
- concrete cross-module knowledge where explicit contracts would be safer

Phase 0.11 addresses that structural gap.

---

## Problem Statement

Without a clear domain module pattern, the backend becomes harder to maintain even if its public contract remains stable.

The risk is not primarily public API breakage. The risk is architectural drift:

- transport concerns leak inward
- use-case orchestration stays embedded in handlers
- domain ownership remains implicit
- module dependencies become historical rather than deliberate
- future work adds more coupling instead of building on clean structural boundaries

Phase 0.11 therefore focuses on internal consistency rather than new business behavior.

---

## Scope

### Included

- definition of a standard internal module pattern
- application of that pattern to the current Stage 0 domain-facing modules
- explicit clarification of layer responsibilities
- explicit clarification of ownership boundaries between `auth`, `user` and `usersettings`
- consolidation of minimal cross-module contracts where coordination is required
- trunk-documentation alignment for the adopted pattern

### Excluded

- new endpoints
- payload changes
- public contract changes
- redesign of challenge / verify semantics
- redesign of authorization semantics
- CQRS / event sourcing introduction
- multi-tenant architecture
- generic abstraction-heavy rewrites without immediate structural value

---

## Target Pattern

The target module shape is:

```text
internal/modules/<module>/
    http/
        handlers.go
        dto.go

    app/
        service.go
        usecases.go

    domain/
        model.go
        contracts.go

    repository/
        repository.go   (when it adds real clarity)
```

This shape is intended to standardize responsibilities rather than force identical file counts in every module.

### Layer intent

#### HTTP

Owns:

- request parsing
- transport validation
- response mapping
- transport-facing error projection

Does not own:

- use-case orchestration
- domain semantics
- persistence decisions

#### APP

Owns:

- use-case orchestration
- sequencing of module-local operations
- coordination with external dependencies through explicit boundaries

Does not own:

- raw transport concerns
- storage-specific implementation details
- domain ownership transfer from other modules

#### DOMAIN

Owns:

- module-local models
- invariants
- semantic vocabulary
- explicit contracts

Does not own:

- HTTP DTO shape
- router behavior
- persistence implementation details

#### REPOSITORY

Owns:

- persistence-facing implementation boundaries when needed

This layer is optional by design. It should exist where it improves clarity, not as a mandatory abstraction tax.

---

## Dependency Direction

The intended direction is:

```text
HTTP → APP → DOMAIN
```

Repository implementations support the application/domain boundary when needed. The point of this rule is to keep transport orchestration thin and to prevent later phases from moving business flow back into handlers.

---

## Subphases

### 0.11.1 — Domain Module Pattern Definition

Status: ✔ Completed

Define:

- the standard internal module shape
- layer responsibilities
- dependency direction
- naming and organization expectations

Result expected from 0.11.1:

- a repository-level architectural definition of the module pattern
- no runtime behavior change

---

### 0.11.2 — User Module Refactor

Status: ✔ Completed

Apply the pattern to `internal/modules/user`.

Goals:

- make `user` the first concrete reference implementation of the pattern
- move orchestration responsibilities out of handlers
- make user-specific ownership explicit through `app` and `domain`
- preserve the existing external user-related contract

Expected result:

- thinner handlers
- clearer user use cases
- explicit user-domain ownership

---

### 0.11.3 — UserSettings Module Refactor

Status: ✔ Completed

Apply the pattern to `internal/modules/usersettings`.

Goals:

- align the module structurally with the new pattern
- preserve its own configuration/settings semantics
- avoid treating settings as an informal extension of the user entity
- keep update/read semantics externally unchanged

Expected result:

- settings-specific orchestration becomes explicit
- settings ownership stays distinct from user profile ownership
- defaults / effective settings semantics have a clearer domain home when relevant

---

### 0.11.4 — Auth Module Alignment

Status: ◑ Defined

Apply the pattern conservatively to `internal/modules/auth`.

Goals:

- align auth structurally without redesigning authentication behavior
- express challenge / verify / authenticated-entry flows as use cases
- preserve the already stabilized challenge, verify and authorization-adjacent semantics
- avoid collapsing new architectural clarity back into transport code

Expected result:

- auth flows become easier to reason about internally
- completed 0.10 behavior remains externally unchanged

---

### 0.11.5 — Cross-Module Contract Consolidation

Status: ◑ Defined

Clarify and minimize the way the modules interact.

Goals:

- reduce concrete cross-module knowledge
- define minimal explicit contracts where coordination is real
- keep ownership boundaries clear between `auth`, `user` and `usersettings`

Key rule:

- coordination is allowed
- ownership transfer is not

Expected result:

- the modules can collaborate without informal dependency growth

---

### 0.11.6 — Documentation & Phase Closure

Status: ◑ Defined

Close the phase by aligning the trunk documentation set with the adopted module pattern.

Goals:

- reflect the new pattern across roadmap, status, architecture, handoff and README where relevant
- record the phase as an internal structural consolidation step
- state clear completion criteria

Expected result:

- one coherent documentary narrative for Phase 0.11

---

## Ownership Boundaries

Phase 0.11 assumes the following ownership model:

- `auth` owns authentication flows and authenticated-entry behavior
- `user` owns the user entity/profile boundary
- `usersettings` owns settings/configuration semantics

This matters because the pattern is not only about folders. It is also about preventing later work from blurring domain ownership.

---

## Cross-Module Coordination Principle

Where module interaction is necessary, the repository should prefer minimal explicit contracts over concrete internal knowledge.

That principle is especially important between:

- `auth` and `user`
- `auth` and `usersettings`
- `user` and `usersettings`

The purpose is not abstraction for its own sake. The purpose is to ensure that real collaboration does not become accidental coupling.

---

## Compatibility Guarantee

Phase 0.11 is intentionally non-breaking.

It should preserve:

- endpoints
- payloads
- error model behavior
- canonical `/api/v1/...` contract state
- authenticated bootstrap behavior
- challenge / verify semantics
- currently delivered authorization behavior

The phase is therefore structural in impact, not functional in public behavior.

---

## Risks

The main structural risks are:

- over-abstracting modules with interfaces that do not add real value
- accidentally mixing structural refactor with functional redesign
- reintroducing ownership confusion while trying to “share” code
- keeping direct cross-module knowledge even after module-local cleanup

Mitigation for all of them is the same:

- keep the phase pragmatic
- keep the contracts minimal
- preserve ownership clarity
- preserve external behavior

---

## Completion Criteria

Phase 0.11 should only be considered complete when:

- the module pattern is explicitly defined
- `auth`, `user` and `usersettings` are aligned to it
- layer responsibilities are visible and coherent
- cross-module contracts are explicit where needed
- ownership boundaries are clearer than in the pre-0.11 state
- no external HTTP contract regression has been introduced
- trunk documentation reflects the new structural state coherently

---

## Current Repository State

At this repository point, Phase 0.11 is no longer only a defined target. The phase is **in progress**.

0.11.1 is completed at the architectural-definition level. The repository now has an explicit description of the Domain Module Pattern, explicit layer responsibilities, explicit dependency direction and explicit ownership boundaries between `auth`, `user` and `usersettings`.

0.11.2 is also completed in runtime repository state. The `internal/modules/user` package is now the first concrete implementation of the pattern, with explicit `app`, `domain` and `repository` boundaries and backward-compatible package access preserved for existing consumers.

0.11.3 is now also completed in runtime repository state. The `internal/modules/usersettings` package follows the same pattern with explicit `app`, `domain` and `repository` boundaries while preserving settings-specific semantics and backward-compatible package access for existing consumers.

## Current Result After 0.11.1

Phase 0.11.1 specifically delivered:

- the standard internal module shape for the current Stage 0 modules
- the explicit `HTTP → APP → DOMAIN` dependency direction
- the responsibility split between transport, orchestration, domain semantics and optional persistence boundaries
- the ownership model between `auth`, `user` and `usersettings`
- the explicit non-goals and completion criteria for the rest of the phase

No runtime behavior was changed by 0.11.1. Its purpose was to close the structural-definition gap before any controlled runtime migration began.

### Next

0.11.4 — Auth Module Alignment

## Current Result After 0.11.2

Phase 0.11.2 specifically delivers:

- the first concrete runtime application of the Domain Module Pattern under `internal/modules/user`
- explicit `app`, `domain` and `repository` package boundaries for the user module
- preservation of current external user-related behavior and existing consumers through backward-compatible root-package access
- repository state validated by passing tests, including the `internal/modules/user` package set

The repository is therefore no longer only describing 0.11 as a future intention. It has already delivered the first completed structural migration while keeping the Stage 0 public contract intact.

## Current Result After 0.11.3

Phase 0.11.3 specifically delivers:

- the application of the same Domain Module Pattern to `internal/modules/usersettings`
- explicit `app`, `domain` and `repository` package boundaries for the user settings module
- preservation of configuration/settings semantics without collapsing them into the user-entity boundary
- backward-compatible root-package access for existing consumers and passing repository validation for the `internal/modules/usersettings` package set

The repository now contains two completed runtime applications of the pattern (`user` and `usersettings`) while preserving the Stage 0 public contract.

### Next

0.11.4 — Auth Module Alignment
