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

### 0.8.2 — Error Type System Introduction ⬜

Planned:

- internal typed error system
- initial error code catalog
- separation between generic and domain errors

---

### 0.8.3 — Auth Surface Error Standardization ⬜

Planned:

- endpoint-by-endpoint standardization of auth surface error mapping
- consistent HTTP status semantics
- consistent cross-endpoint error behavior

---

### 0.8.4 — Error Mapping Hardening & Contract Tests ⬜

Planned:

- error contract tests
- mapping hardening
- shape/code/status regression protection

---

## Architecture Introduced in 0.8.1

### Response Contract Layer

The backend now has an explicit reusable response error model:

- `internal/core/errs/response_error.go`

Responsibilities:

- define structured response error shape
- carry code/message/details in a transport-neutral form

---

### HTTP Serialization Layer

The backend now has centralized HTTP error writing:

- `internal/core/httpx/error.go`

Responsibilities:

- wrap response errors under the canonical `error` envelope
- serialize consistent JSON errors across handlers and middleware

---

### Auth Transport Adoption

The auth module now routes handler error responses through the centralized error envelope writer.

This includes:

- authenticated claims failure paths
- bad request decoding paths
- auth handler error responses
- wallet handler error responses
- auth middleware unauthorized/configuration failures
- panic and timeout transport responses in `httpx`

---

## Guarantees After 0.8.1

- one structured JSON error envelope exists
- auth HTTP handlers no longer emit legacy root-level string-only errors
- error details are now nested under `error.details`
- middleware-level auth failures align with the new envelope
- success payloads remain unchanged
- routes remain unchanged
- domain/application logic remains unchanged

---

## Result

Phase 0.8 has now started concretely in code.

0.8.1 establishes the contractual base required for the remaining subphases:

- 0.8.2 will define the internal type system
- 0.8.3 will standardize error mapping semantics across auth surface
- 0.8.4 will harden and freeze the contract with tests

---

## Conclusion

Phase 0.8.1 does not redesign business behavior.

It introduces the formal response contract required for all later error standardization work, while preserving the current auth execution model established by Phase 0.7.
