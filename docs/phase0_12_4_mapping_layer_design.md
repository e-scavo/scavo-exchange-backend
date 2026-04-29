# Phase 0.12.4 — Mapping Layer Introduction

## Subphase 0.12.4.1 — Mapping Layer Design

---

## Objective

Design the dedicated module-local Mapping Layer that will consolidate transformation ownership introduced progressively across Phase 0.12.2 and Phase 0.12.3.

This subphase defines the target structure, naming rules, package responsibilities, migration strategy, and compatibility constraints before any code-level mapper relocation or application refactor is performed.

No code changes are performed in this subphase.

---

## Context

Phase 0.12.2 introduced explicit Read Models and Domain/Application → Read mapping.

Phase 0.12.3 introduced explicit Write Models, domain write inputs, Write → Domain mapping, and handler alignment.

The system is currently stable and validated, but mapping ownership is still distributed across model packages:

```text
internal/modules/<module>/readmodels/mappers.go
internal/modules/<module>/writemodels/mappers.go
```

Some response assembly and helper-level mapping also remains near application orchestration code.

The target of Phase 0.12.4 is not to change runtime behavior, but to make mapping ownership explicit and centralized per module.

---

## Current State

Current mapping ownership is valid but transitional.

### Read-side mapping

Read model packages currently own Domain/Application → Read transformations.

Examples:

```text
internal/modules/auth/readmodels/mappers.go
internal/modules/user/readmodels/mappers.go
internal/modules/usersettings/readmodels/mappers.go
```

### Write-side mapping

Write model packages currently own Write → Domain transformations.

Examples:

```text
internal/modules/auth/writemodels/mappers.go
internal/modules/user/writemodels/mappers.go
internal/modules/usersettings/writemodels/mappers.go
```

### Application-side assembly

Some orchestration-level transformations still happen near application/support code because they combine domain/application values into response-facing structures.

This is acceptable before 0.12.4, but should not remain the final ownership model.

---

## Problem Statement

Mapping logic currently lives close to the model packages it transforms.

This creates several long-term risks:

* read model packages may accumulate domain knowledge
* write model packages may accumulate domain construction rules
* application services may continue collecting transformation helpers
* future modules may duplicate mapper patterns inconsistently
* hybrid behavior can reappear through convenience conversions

The system needs a clear architectural rule:

```text
Models define shape.
Mappers transform shape.
Applications orchestrate use cases.
```

---

## Target Mapping Layer

Each module may define a dedicated mapper package:

```text
internal/modules/<module>/mappers/
```

Examples:

```text
internal/modules/auth/mappers/
internal/modules/user/mappers/
internal/modules/usersettings/mappers/
```

This package becomes the single owner of module-local transformation logic.

---

## Package Responsibility Model

### domain/

Owns canonical domain state and domain input structures.

Allowed:

* domain entities
* domain services/contracts
* domain input structs
* domain behavior

Not allowed:

* HTTP response construction
* HTTP request parsing
* read model construction
* write model parsing rules

---

### readmodels/

Owns output-only model definitions.

Allowed:

* read model structs
* response projection shape
* JSON tags for stable output contracts

Not allowed:

* domain-to-read mapping logic after consolidation
* business logic
* write/input semantics
* handler orchestration

---

### writemodels/

Owns input-only model definitions.

Allowed:

* write model structs
* request intent shape
* JSON tags for stable input contracts

Not allowed:

* write-to-domain mapping logic after consolidation
* response semantics
* domain behavior
* handler orchestration

---

### mappers/

Owns explicit transformation logic.

Allowed:

* Write → Domain mapping
* Domain → Read mapping
* collection mapping
* nil/default handling required for compatibility
* module-local adapter helpers

Not allowed:

* HTTP request decoding
* HTTP response writing
* database access
* business decisions
* authorization decisions
* cross-module orchestration outside documented dependencies

---

### app/

Owns application use-case orchestration.

Allowed:

* service methods
* use-case flow
* domain service coordination
* mapper invocation

Not allowed:

* inline transformation logic that belongs in mappers
* duplicated mapper helpers
* direct reconstruction of read/write/domain shape when a mapper exists

---

## Mapping Directions

The Mapping Layer standardizes two primary directions:

```text
Write Model → Domain Input / Domain Model
Domain / Application Result → Read Model
```

The preferred flow is:

```text
HTTP Request → Write Model → Mapper → Domain/Application → Mapper → Read Model → HTTP Response
```

---

## Naming Strategy

Mapper function names must be explicit and direction-aware.

### Write → Domain

```text
<WriteModelName>To<DomainInputName>
```

Examples:

```text
AuthLoginWriteModelToLoginInput
AuthWalletVerifyWriteModelToWalletVerifyInput
UserUpdateWriteModelToUserUpdateInput
UserSettingsUpdateWriteModelToUserSettingsUpdateInput
```

### Domain → Read

```text
<DomainModelName>To<ReadModelName>
```

Examples:

```text
WalletIdentityToAuthWalletReadModel
WalletChallengeToAuthWalletChallengeReadModel
UserToUserReadModel
UserSettingsToUserSettingsReadModel
```

### Collections

Collection mappers should be plural and explicit:

```text
WalletIdentitiesToAuthWalletReadModels
```

---

## Method vs Function Rule

Current write-side mapping uses methods like:

```text
writeModel.ToDomainInput()
```

The target mapping layer should move toward package-level functions:

```text
mappers.AuthLoginWriteModelToLoginInput(model)
```

Rationale:

* keeps model packages shape-only
* prevents models from owning transformation behavior
* reduces coupling between input models and domain packages
* makes mapping ownership searchable and explicit

During migration, compatibility wrappers may remain temporarily if needed, but final ownership belongs to `mappers/`.

---

## Module Dependency Rules

The mapper package may import:

* its module domain package
* its module readmodels package
* its module writemodels package
* standard library packages required for conversion

The mapper package must not import:

* HTTP handler packages
* repository packages
* app packages, unless explicitly required and documented for application-result projection
* unrelated module internals without a clear contract

Preferred dependency direction:

```text
app → mappers → domain/readmodels/writemodels
```

Avoid:

```text
domain → mappers
readmodels → mappers
writemodels → mappers
mappers → app
```

---

## Module-Level Design

### Auth Module

Target package:

```text
internal/modules/auth/mappers/
```

Initial ownership:

* login write input mapping
* profile update write input mapping
* settings update write input mapping
* wallet challenge write input mapping
* wallet verification write input mapping
* wallet detach write input mapping
* wallet primary-set write input mapping
* wallet identity read projection
* wallet challenge read projection
* login read projection

Current sources to consolidate:

```text
internal/modules/auth/readmodels/mappers.go
internal/modules/auth/writemodels/mappers.go
application/support helper mapping where applicable
```

---

### User Module

Target package:

```text
internal/modules/user/mappers/
```

Initial ownership:

* user domain → user read model
* user update write model → user update domain input

Current sources to consolidate:

```text
internal/modules/user/readmodels/mappers.go
internal/modules/user/writemodels/mappers.go
```

---

### User Settings Module

Target package:

```text
internal/modules/usersettings/mappers/
```

Initial ownership:

* user settings domain → user settings read model
* user settings write model → user settings update domain input
* preference JSON decoding compatibility
* default preference map stabilization

Current sources to consolidate:

```text
internal/modules/usersettings/readmodels/mappers.go
internal/modules/usersettings/writemodels/mappers.go
```

---

## Compatibility Rules

The mapping layer must preserve:

* current route behavior
* current request JSON payload semantics
* current response JSON payload semantics
* current error behavior for invalid input
* current nil/default handling
* current test expectations

Any mapper relocation must be behavior-preserving.

---

## Migration Strategy

The migration must be incremental.

### Step 1 — Create mapper packages

Introduce `mappers/` packages without deleting existing mapper methods/functions immediately if compatibility requires a staged move.

### Step 2 — Move read-side transformations

Move Domain → Read logic into `mappers/`.

### Step 3 — Move write-side transformations

Move Write → Domain logic into `mappers/`.

### Step 4 — Update applications/handlers

Update callers to depend on `mappers/` instead of model-owned mapping helpers.

### Step 5 — Remove or deprecate old mapper ownership

Once all callers are aligned, model packages should return to model definitions only.

---

## Non-Goals

This subphase does not introduce:

* new endpoint behavior
* new business logic
* new validation semantics
* route changes
* API version changes
* CQRS
* event sourcing
* multi-tenant behavior

---

## Risks

### Import cycles

The new mapper packages must avoid importing application packages unless absolutely necessary.

Mitigation:

* keep dependencies directed toward model/domain packages
* keep app as caller, not dependency

### Contract drift

Moving mapping logic may accidentally alter response defaults.

Mitigation:

* preserve existing nil/default behavior exactly
* validate with `go test ./...`

### Duplicate mapper survival

Old and new mapper functions may coexist temporarily.

Mitigation:

* consolidate in 0.12.4.3
* remove duplication only after caller migration is complete

---

## Expected Output of This Subphase

This subphase produces a design-level contract for the Mapping Layer.

It defines:

* package location
* responsibility boundaries
* naming rules
* dependency rules
* module-level targets
* migration sequence
* compatibility constraints

No code changes are introduced here.

---

## Status

Phase: 0.12.4  
Subphase: 0.12.4.1  
Status: Completed  
Code impact: None  
Next: 0.12.4.2 — Mapping Layer Implementation
