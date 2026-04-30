# Phase 0.15.4 — Response Schema Normalization

---

## Status

Completed.

---

## Context inherited

0.15.4 starts after three contract-hardening baselines were already established:

- 0.15.1 documented the real HTTP route surface.
- 0.15.2 aligned the canonical public error envelope.
- 0.15.3 validated provider boundaries through compile-time assertions.

The backend therefore entered this subphase with stable routes, stable error envelopes and validated internal provider seams.

---

## Real problem

The success and error response payloads were already structurally compatible with the documented contracts, but two response-writing details still allowed minor schema drift:

1. Auth error responses written through the auth-local helper used `application/json` while the core HTTP writer used `application/json; charset=utf-8`.
2. The defensive timeout fallback JSON in the core HTTP middleware still lacked the `details` field introduced as mandatory by 0.15.2.

These differences did not change business behavior, but they left response serialization policy split across handlers and middleware.

---

## Decision taken

Normalize response serialization details without changing public routes, status codes, success payloads, error codes, provider behavior, domain logic or repository behavior.

The subphase remains intentionally narrow: it does not introduce a new envelope for successful responses and does not wrap existing payloads.

---

## Concrete change

0.15.4 applies two compatibility-preserving changes:

- `internal/modules/auth/http_login.go` now emits auth error JSON with `application/json; charset=utf-8`, matching the core HTTP JSON writer.
- `internal/core/httpx/middleware.go` now keeps the defensive timeout fallback aligned with the canonical `{error:{code,message,details}}` shape.

A regression test was added in:

- `internal/modules/auth/http_handlers_test.go`

The test confirms that auth error responses use the canonical JSON content type.

---

## What did not change

0.15.4 does not change:

- HTTP routes
- HTTP methods
- status codes
- success payload field names
- success response envelopes
- error codes
- error messages
- provider interfaces
- application logic
- domain rules
- repository behavior
- WebSocket behavior

---

## Impact observable

After 0.15.4:

- JSON response metadata is aligned between core and auth handlers.
- Defensive timeout fallback JSON remains compatible with the mandatory `details` object.
- Existing frontend-facing response payloads remain compatible.
- 0.15.5 can define freeze enforcement on top of a normalized response serialization baseline.

---

## Validation

The subphase is intended to be validated with:

```bash
go test ./...
```

No local test execution is recorded in this generated package because the assistant environment does not provide the required Go toolchain/runtime parity for this repository.

---

## Handoff to 0.15.5

The next subphase is:

**0.15.5 — Contract Freeze Enforcement**

It must define the rules that prevent future contract drift. It must use the route audit, error alignment, provider validation and response normalization artifacts as the freeze baseline.

---

## End of Document
