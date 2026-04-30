# Phase 0.15.2 — Error Contract Alignment

---

## Status

**Phase:** 0.15 — Contract Hardening & Freeze  
**Subphase:** 0.15.2 — Error Contract Alignment  
**Status:** Completed  
**Type:** Contract hardening  
**Code changes:** Yes  

---

## Context Inherited

0.15.2 starts from the 0.15.1 HTTP Contract Audit baseline.

That audit made the HTTP surface explicit before changing any contract behavior. It recorded foundation routes, diagnostics, the WebSocket upgrade route, public auth routes, authenticated account/settings routes and authenticated wallet-management routes. It also confirmed that legacy auth paths and canonical `/api/v1` paths are active paired contracts backed by the same handler set.

The inherited error model comes from Phase 0.8 and was enriched by Phase 0.14:

- Phase 0.8 introduced the standardized public error envelope.
- Phase 0.8.2 centralized normalized application errors under `internal/core/errs`.
- Phase 0.8.3 moved auth-surface error emission to centralized app-error factories.
- Phase 0.8.4 added contract tests for representative mappings and envelope serialization.
- Phase 0.14.3 added safe diagnostic context enrichment without replacing the public error shape.
- Phase 0.15.1 documented the route surface that 0.15.2 must use as the alignment baseline.

---

## Real Problem

The backend already used the standardized public error envelope:

```json
{
  "error": {
    "code": "...",
    "message": "...",
    "details": {}
  }
}
```

However, the implementation still allowed one implicit variation: when no details were provided, `error.details` could be omitted because `ResponseError.Details` used `omitempty` and `NewResponseError()` did not initialize an empty details map.

This created a mismatch between the documented contract and the serialized contract for detail-free errors.

The risk was small but contractual:

- clients could receive `{ error: { code, message } }` for some failures
- documentation and tests described `{ error: { code, message, details } }`
- later contract freeze work would have to freeze an avoidable variation
- frontend consumers would need to treat `details` as conditionally absent instead of a stable object

---

## Decision Taken

0.15.2 aligns the canonical public error contract so `error.details` is always present as a JSON object.

This is a contract-hardening change, not a business-logic change.

The decision preserves the existing envelope, error codes, messages, status mappings, route behavior and handler decisions. It only removes the optional omission of the `details` field from serialized error payloads.

---

## Concrete Change

### Core error contract

`internal/core/errs/response_error.go` was updated so:

- `ResponseError.Details` serializes as `json:"details"` instead of `json:"details,omitempty"`
- `NewResponseError()` always initializes `Details` as a map
- input detail maps are copied before being attached to the response error

This produces stable payloads for errors with and without contextual details.

### Contract tests

`internal/core/errs/app_error_test.go` now verifies:

- detail-free response errors still include an initialized empty details object
- response error details are copied from caller-provided maps and cannot be mutated through the original input map

`internal/core/httpx/error_test.go` now verifies:

- HTTP application errors serialize `error.details` as an empty JSON object when no details exist

---

## Contract After 0.15.2

All public HTTP error responses that use the canonical error writer must preserve this envelope:

```json
{
  "error": {
    "code": "STRING",
    "message": "STRING",
    "details": {}
  }
}
```

Rules:

- `error` is the envelope root.
- `error.code` is required.
- `error.message` is required.
- `error.details` is required and must be a JSON object.
- when no public details exist, `details` must be `{}`.
- contextual details such as `request_id` remain allowed inside `details` when explicitly attached.
- internal causes are not exposed through the public error payload.

---

## Error Writers Covered

The alignment applies to the canonical error construction path used by:

- `errs.NewResponseError()`
- `errs.AppError.ToResponseError()`
- `httpx.WriteError()`
- `httpx.WriteErrorMessage()`
- `httpx.WriteAppError()`
- auth transport helpers that build envelopes through `coreerrs.NewErrorEnvelope(appErr.ToResponseError())`
- timeout middleware fallback payloads generated through the same envelope constructor

---

## Explicit Non-Changes

0.15.2 does not change:

- routes
- HTTP methods
- success payloads
- status code decisions
- error code names
- error messages
- auth behavior
- authorization behavior
- provider behavior
- application/domain behavior
- repository behavior
- WebSocket protocol behavior

No new feature is introduced.

---

## Observable Impact

After 0.15.2:

- the public error envelope is stable across detail-free and detail-rich errors
- `details` is no longer conditionally omitted
- frontend consumers can treat `error.details` as a stable object
- existing request correlation details remain supported
- the error contract is ready for later response-schema normalization and freeze enforcement

---

## Validation

Validation was attempted in the execution environment with:

```bash
GOTOOLCHAIN=local go test ./internal/core/errs ./internal/core/httpx
```

The command could not execute because the local toolchain is Go 1.23.2 while `go.mod` requires Go 1.25.0.

A full validation command must be run in the developer environment:

```bash
go test ./...
```

---

## Handoff to 0.15.3

The next subphase is:

**0.15.3 — Provider Contract Validation**

It must start from the aligned error contract defined here and from the HTTP route inventory created in 0.15.1.

It must not change the public error envelope unless a new explicitly approved contract issue is detected.

---

## End of Document
