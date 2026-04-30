# Phase 0.15.5 — Contract Freeze Enforcement

---

## Status

**Phase:** 0.15 — Contract Hardening & Freeze  
**Subphase:** 0.15.5 — Contract Freeze Enforcement  
**Status:** Completed  
**Code changes:** yes  
**Behavior changes:** no  

---

## Context Inherited

0.15.5 starts after four contract-hardening steps:

- 0.15.1 audited the real HTTP route surface.
- 0.15.2 aligned the canonical error envelope.
- 0.15.3 validated provider contracts with compile-time assertions.
- 0.15.4 normalized response serialization metadata and defensive fallback shape.

The backend therefore enters this subphase with explicit evidence for the current public and internal contracts.

---

## Real Problem

The repository had contract baselines, but not yet a freeze rule.

Without explicit freeze enforcement, future work could silently change:

- route status behavior
- response shape
- error envelope fields
- JSON content type policy
- provider boundary assumptions
- documentation state

This would reintroduce drift after the audit and normalization work already completed in 0.15.1 through 0.15.4.

---

## Decision Taken

Freeze the current Stage 0 contract surface.

The freeze does not introduce a new API version, a new response envelope, new routes or new behavior.

It defines what is protected and adds regression coverage so representative contract drift fails during tests.

---

## Frozen Contract Surface

The following contracts are frozen for the current Stage 0 backend baseline:

### HTTP route surface

The audited route set from 0.15.1 is the current public surface.

Adding, removing or changing route behavior requires explicit documentation and, where externally visible, versioning analysis.

### Canonical error envelope

Public HTTP errors must preserve:

```json
{
  "error": {
    "code": "...",
    "message": "...",
    "details": {}
  }
}
```

`details` remains a JSON object and must not be omitted.

### JSON response metadata

JSON responses must preserve:

```text
application/json; charset=utf-8
```

### Provider boundary contracts

The compile-time assertions introduced in 0.15.3 are part of the freeze baseline.

Provider drift must fail compilation instead of surfacing as handler drift.

### Documentation synchronization

Any contract-affecting change must update the relevant trunk documentation and dedicated phase/subphase document.

---

## Allowed Future Evolution

Future work may evolve contracts only when it is explicit.

Allowed paths:

- additive compatible fields
- documented versioned routes
- documented error codes
- documented response-family expansion
- test updates tied to an intentional contract decision

---

## Forbidden Silent Drift

The following is not allowed without an explicit phase/subphase decision:

- removing public response fields
- omitting `error.details`
- changing canonical error field names
- changing auth legacy and `/api/v1` parity accidentally
- changing JSON content type policy accidentally
- weakening provider compile-time assertions
- documenting a contract that does not match code

---

## Concrete Change

0.15.5 introduces a dedicated regression test file:

```text
internal/core/httpx/contract_freeze_test.go
```

The tests guard:

- core status route JSON content type and required response keys
- protected auth route canonical error envelope
- legacy and `/api/v1` protected auth parity for the frozen missing-bearer contract
- mandatory `details` object on canonical auth errors

---

## Impact Observable

After 0.15.5:

- frozen contracts have explicit policy
- representative drift is test-detectable
- documentation explains future evolution rules
- no runtime behavior changes were introduced
- Phase 0.15 is ready for 0.15.6 validation and documentation closure

---

## Handoff to 0.15.6

0.15.6 must validate the full repository with `go test ./...`, reconcile all trunk documentation and close Phase 0.15 narratively.

It must not introduce new features or unrelated refactors.
