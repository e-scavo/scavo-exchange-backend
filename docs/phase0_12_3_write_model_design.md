# Phase 0.12.3 — Write Model Isolation

## Subphase 0.12.3.1 — Write Model Design

---

## Objective

Define the explicit set of **Write Models** to be introduced in the system, based on the results of Phase 0.12.1 and the Read Model separation implemented in Phase 0.12.2.

This subphase establishes:

* which input models must be created
* how they are structured
* how they map to domain models
* how they replace hybrid usage
* naming and placement conventions

No code changes are performed in this subphase.

---

## Context

From Phase 0.12.1:

* Hybrid models detected: 11
* Write-capable models detected: subset of HYBRID + existing request structs

From Phase 0.12.2:

* Read Models introduced and aligned
* Domain → Read mapping established
* Responses decoupled from domain

Current problem:

```text
Input models are still mixed with domain and transport responsibilities
```

Target:

```text
WRITE → DOMAIN → READ
```

---

## Scope

### Included

* definition of Write Models
* input structure design
* naming conventions
* module placement
* conceptual mapping (Write → Domain)

### Excluded

* code implementation
* handler modification
* contract changes
* validation logic refactor

---

## Write Model Principles

A Write Model MUST:

* represent **input intent only**
* contain only required fields
* not expose domain internals
* not be reused as response
* be minimal and explicit
* be stable across layers

---

## Naming Convention

```text
<Context><Action>WriteModel
```

Examples:

* `AuthLoginWriteModel`
* `AuthWalletLinkWriteModel`
* `AuthWalletMergeWriteModel`
* `UserUpdateWriteModel`
* `UserSettingsUpdateWriteModel`

---

## Location Strategy

Write Models will be placed under:

```text
internal/modules/<module>/writemodels/
```

Examples:

```text
internal/modules/auth/writemodels/
internal/modules/user/writemodels/
internal/modules/usersettings/writemodels/
```

---

## Source → Target Identification

### AUTH MODULE

#### Existing Inputs

* Login request payload
* Wallet link / merge / detach flows
* Challenge flows

#### Write Models

```text
AuthLoginWriteModel
AuthWalletLinkWriteModel
AuthWalletMergeWriteModel
AuthWalletDetachWriteModel
AuthWalletChallengeWriteModel
```

---

### USER MODULE

#### Existing Inputs

* user updates
* partial updates (if any)

#### Write Models

```text
UserUpdateWriteModel
```

---

### USER SETTINGS MODULE

#### Existing Inputs

* preferences update
* flags update

#### Write Models

```text
UserSettingsUpdateWriteModel
```

---

## Hybrid Model Resolution

Each HYBRID model previously identified will be split into:

```text
HYBRID → WRITE + DOMAIN + READ
```

Example:

```text
User → UserUpdateWriteModel + UserDomain + UserReadModel
```

---

## Domain Relationship

Write Models:

* DO NOT replace domain models
* DO NOT contain behavior
* ARE NOT persisted directly
* must be mapped explicitly into domain entities

---

## Mapping Strategy (Conceptual)

Mapping direction:

```text
WRITE → DOMAIN
```

Rules:

* explicit mapping only
* no direct struct reuse
* no embedding domain inside write model
* validation occurs before or during mapping

---

## Mapping Examples

```text
AuthLoginWriteModel → AuthDomainLoginInput
UserUpdateWriteModel → UserDomain
UserSettingsUpdateWriteModel → UserSettingsDomain
```

---

## Constraints

* MUST preserve existing HTTP contracts
* MUST allow progressive migration
* MUST not break existing handlers
* MUST coexist with legacy request parsing

---

## Risks

* duplicated validation logic
* incomplete field coverage
* hidden coupling in handlers

Mitigation:

* mapping layer introduced in 0.12.3.3
* handler alignment deferred to 0.12.3.4

---

## Output of This Subphase

This document defines:

* all Write Models to be created
* their naming and structure
* their placement
* their mapping direction

---

## Status

Phase: 0.12.3
Subphase: 0.12.3.1
Status: COMPLETED
Code Impact: NONE

---

## Next Step

```text
0.12.3.2 — Write Model Implementation
```
