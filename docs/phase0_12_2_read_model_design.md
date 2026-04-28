# Phase 0.12.2 — Read Model Extraction

## Subphase 0.12.2.1 — Read Model Design

---

## Objective

Define the explicit set of **Read Models** to be introduced in the system, based on the results of Phase 0.12.1.

This subphase establishes:

* which models will be extracted
* how they will be structured
* how they relate to existing domain and hybrid models
* naming and organization conventions
* mapping strategy (conceptual only)

No code changes are performed in this subphase.

---

## Context

From Phase 0.12.1:

* Total models: 125
* Hybrid models: 11
* Cross-layer coupling detected
* Separation strategy defined:

```text
HYBRID → READ + WRITE + DOMAIN
```

---

## Scope

### Included

* definition of Read Models
* structural design
* naming conventions
* module placement
* mapping definition (conceptual)

### Excluded

* code implementation
* handler modification
* contract changes
* versioning changes

---

## Read Model Principles

A Read Model MUST:

* be used only for output
* not contain mutation intent
* not be reused as input
* be optimized for response shape
* be decoupled from domain internals

---

## Naming Convention

All Read Models must follow:

```text
<Context><Entity>ReadModel
```

Examples:

* `AuthLoginReadModel`
* `AuthWalletReadModel`
* `UserReadModel`
* `UserSettingsReadModel`

---

## Location Strategy

Read Models will be placed under:

```text
internal/modules/<module>/readmodels/
```

Example:

```text
internal/modules/auth/readmodels/
internal/modules/user/readmodels/
internal/modules/usersettings/readmodels/
```

---

## Source → Target Mapping

### AUTH MODULE

#### LoginResponse (existing HYBRID/READ-like)

→ Target:

```text
AuthLoginReadModel
```

Fields:

* user_id
* token
* expiration
* wallet_info (if present)

---

#### WalletResponse

→ Target:

```text
AuthWalletReadModel
```

Fields:

* wallet_address
* chain
* balance
* metadata

---

### USER MODULE

#### User (HYBRID / DOMAIN overlap)

→ Target:

```text
UserReadModel
```

Fields:

* id
* email
* username
* status
* created_at

---

### USER SETTINGS MODULE

#### UserSettings (HYBRID)

→ Target:

```text
UserSettingsReadModel
```

Fields:

* user_id
* preferences
* flags
* timestamps

---

## Hybrid Model Resolution

For each HYBRID:

| Original       | Target Read Model     |
| -------------- | --------------------- |
| User           | UserReadModel         |
| UserSettings   | UserSettingsReadModel |
| WalletResponse | AuthWalletReadModel   |
| LoginResponse  | AuthLoginReadModel    |

---

## Domain Relationship

Read Models:

* DO NOT replace domain models
* DO NOT contain business logic
* are projections of domain state

---

## Mapping Strategy (Conceptual)

Mapping direction:

```text
Domain → Read Model
```

Rules:

* explicit mapping only
* no implicit reuse
* no struct aliasing
* no embedding of domain structs

---

## Mapping Examples

```text
UserDomain → UserReadModel
UserSettingsDomain → UserSettingsReadModel
WalletDomain → AuthWalletReadModel
```

---

## Constraints

* MUST preserve current HTTP responses
* MUST allow gradual migration
* MUST not break compatibility
* MUST support future write model separation

---

## Risks

* incomplete field mapping
* hidden domain dependencies
* handler coupling

Mitigation:

* mapping layer introduced in 0.12.2.3
* progressive replacement

---

## Output of This Subphase

This document defines:

* all Read Models to be created
* their structure
* their naming
* their placement
* their relationship to existing models

---

## Status

Phase: 0.12.2
Subphase: 0.12.2.1
Status: COMPLETED
Code Impact: NONE

---

## Next Step

```text
0.12.2.2 — Read Model Implementation
```
