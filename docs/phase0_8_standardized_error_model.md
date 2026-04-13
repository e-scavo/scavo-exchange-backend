# Phase 0.8 — Standardized Error Model

## Objective

Define and apply a uniform error model across the authenticated backend surface so errors become consistent, predictable, traceable, and easier for clients to consume.

---

## Initial Context

After Phase 0.7, the auth module already had:

- application-driven handlers
- transport simplification
- domain-semantic error conditions
- preserved success contracts

What remained inconsistent was the error contract itself.

Legacy responses still used payloads such as:

- `{"error":"unauthorized"}`
- `{"error":"bad_request"}`
- `{"error":"wallet_detach_not_eligible","check":{...}}`

This made error consumption inconsistent across endpoints and left the backend without a formal structured error envelope.

---

## Phase Breakdown

### 0.8.1 — Error Contract Definition ✔

Delivered in this subphase:

- introduced a shared response error contract under `internal/core/errs`
- introduced centralized HTTP error envelope writing under `internal/core/httpx`
- migrated auth middleware and auth handlers to emit the new envelope
- preserved existing domain decisions and transport flow while changing only the error response structure

Resulting envelope shape:

```json
{
  "error": {
    "code": "AUTH_UNAUTHORIZED",
    "message": "authentication required"
  }
}
```

When contextual data is needed, the envelope now supports:

```json
{
  "error": {
    "code": "WALLET_CANNOT_DETACH",
    "message": "wallet cannot be detached under current ownership rules",
    "details": {
      "check": {
        "eligible": false
      }
    }
  }
}
```

---

### 0.8.2 — Error Type System Introduction ✔

Delivered in this subphase:

- introduced a reusable internal `AppError` type under `internal/core/errs`
- introduced centralized error categories for generic/auth/wallet/settings/internal cases
- moved the auth normalized error catalog out of `internal/modules/auth/http_login.go` into `internal/core/errs`
- added typed helper factories and conversion helpers from internal errors to response errors
- added `WriteAppError(...)` support in `internal/core/httpx`
- preserved the 0.8.1 transport envelope while removing auth-local normalized error definitions

Result:

The backend now has a centralized internal error type system that sits underneath the 0.8.1 HTTP envelope and prepares later auth-surface-wide mapping work without changing current routes or success payloads.

---

### 0.8.3 — Auth Surface Error Standardization ✔

Delivered in this subphase:

- migrated auth transport handlers to consume centralized `core/errs` factories directly
- reduced handler-local ownership of auth/wallet/settings error codes and messages
- preserved auth-local JSON writing to avoid cyclic imports through `internal/core/httpx/router.go`
- kept routes, success payloads and business flows unchanged while standardizing auth surface error emission

Result:

The auth surface now writes the Phase 0.8 standardized error contract from centralized app-error factories instead of relying on handler-local string ownership.

---

### 0.8.4 — Error Mapping Hardening & Contract Tests ✔

Delivered in this subphase:

- added dedicated catalog and `AppError` hardening coverage in `internal/core/errs/app_error_test.go`
- added HTTP envelope and transport hardening coverage in `internal/core/httpx/error_test.go`
- froze representative `legacyKey -> code/status/category` mappings under automated tests
- froze canonical `AppError -> ResponseError -> HTTP envelope` behavior under automated tests
- added regression protection for nil/unknown error fallbacks so transport safety remains explicit

Result:

The standardized error model is now not only implemented but also contract-hardened through focused tests that protect shape, code and status behavior against regression.

---

## Architecture Introduced in 0.8.1–0.8.4

### Response Contract Layer

The backend now has an explicit reusable response error model:

- `internal/core/errs/response_error.go`

Responsibilities:

- define structured response error shape
- carry code/message/details in a transport-neutral form

---

### Internal Error Type Layer

The backend now has a reusable application-level error type system:

- `internal/core/errs/app_error.go`
- `internal/core/errs/catalog_auth.go`
- `internal/core/errs/factories.go`

Responsibilities:

- define normalized internal application errors with code/message/status/category
- centralize auth/wallet/settings legacy-key normalization in one place
- provide reusable factories and wrapping helpers for auth surface adoption and later hardening

---

### HTTP Serialization Layer

The backend now has centralized HTTP error writing:

- `internal/core/httpx/error.go`

Responsibilities:

- wrap response errors under the canonical `error` envelope
- serialize consistent JSON errors across handlers and middleware

---

### Auth Transport Adoption

The auth module now routes handler error responses through centralized app-error factories while preserving local JSON writing safety.

This includes:

- authenticated claims failure paths
- bad request decoding paths
- auth handler error responses
- wallet handler error responses
- auth middleware unauthorized/configuration failures
- panic and timeout transport responses in `httpx`
- auth-local JSON writing remaining inside `internal/modules/auth` to avoid cyclic imports while still consuming centralized app-error factories

---

## Guarantees After 0.8.4

- one structured JSON error envelope exists
- auth HTTP handlers no longer emit legacy root-level string-only errors
- error details are now nested under `error.details`
- middleware-level auth failures align with the new envelope
- auth surface handlers now consume centralized app-error factories directly
- representative catalog mappings are now frozen with dedicated tests
- canonical HTTP envelope writing is now frozen with dedicated tests
- success payloads remain unchanged
- routes remain unchanged
- domain/application logic remains unchanged

---

## Result

Phase 0.8 is now fully implemented in code and hardened with contract-focused tests.

0.8.1, 0.8.2, 0.8.3 and 0.8.4 together establish the contractual envelope, the internal typing system, the auth-surface adoption layer and the regression-resistant test hardening required to treat the standardized error model as closed.

---

## Conclusion

Phase 0.8 does not redesign business behavior.

Together its four subphases introduce the formal response contract, the centralized internal error type layer, the auth-surface-wide adoption of centralized app-error factories and the contract-hardening tests that freeze the expected mapping and transport behavior, while preserving the current auth execution model established by Phase 0.7.
