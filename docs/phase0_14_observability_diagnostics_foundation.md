# Phase 0.14 — Observability & Diagnostics Foundation

## Status

**Current phase:** 0.14 — Observability & Diagnostics Foundation  
**Initial subphase:** 0.14.0 — Phase Definition & Documentation Lock  
**Stage:** Stage 0 — Foundation  
**Status:** Pending implementation after documentation lock

---

## Objective

Introduce an observability and diagnostics foundation that makes the backend easier to trace, debug and operate without changing public API behavior.

The phase must improve internal visibility only. It must not introduce new business behavior, new public contracts, external metrics systems, dashboards or runtime dependencies that alter the current request/response surface.

---

## Inherited Context

Phase 0.12 clarified read/write model direction, mapper ownership and internal contract boundaries.

Phase 0.13 then consolidated the Provider Layer as the explicit runtime entry point to domain services. That phase reduced orchestration ambiguity and made handler-to-domain composition cleaner, but it intentionally did not introduce a diagnostics layer.

As a result, the backend now has a cleaner architecture and stable public behavior, but internal runtime visibility remains limited.

The 0.14 transition is therefore:

```text
Phase 0.12: clarify model direction and mapping ownership
Phase 0.13: clarify orchestration ownership and runtime composition
Phase 0.14: clarify runtime visibility and diagnostic traceability
```

---

## Problem Statement

The current backend can be tested and validated, but debugging still depends heavily on manual log inspection and local reasoning.

The main gaps are:

- requests do not have a consistent correlation model
- logs do not consistently carry request-scoped context
- internal errors can be hard to diagnose without additional metadata
- provider/application/domain flow is not observable as a coherent path
- diagnostics are not yet represented as a deliberate foundation capability

---

## Decision

Phase 0.14 defines observability as a Stage 0 foundation concern.

The system will gain request correlation, structured logging conventions, internal error context enrichment, flow tracing and a minimal diagnostics surface, while preserving the external HTTP/API contract.

The phase is explicitly not an OpenTelemetry, Prometheus, dashboard or metrics-export phase.

---

## Scope

Included:

- request correlation model
- request ID / trace propagation through context
- structured logging conventions
- error context enrichment for internal diagnosis
- flow tracing across HTTP, provider, application, domain and repository boundaries
- minimal diagnostics hooks or surface where supported by the existing architecture

Excluded:

- external observability platforms
- dashboards
- Prometheus
- OpenTelemetry
- API contract changes
- business behavior changes
- unrelated refactors

---

## Architecture Impact

The intended diagnostic path is:

```text
HTTP boundary
  → Provider boundary
    → Application service
      → Domain service
        → Repository / store boundary
```

The phase adds visibility to that path without changing ownership.

---

## Subphase Plan

### 0.14.0 — Phase Definition & Documentation Lock ✅ Completed

Lock the phase definition across trunk documentation before code changes.

Scope:

- read and validate all Markdown trunk documents
- correct the 0.14 roadmap definition
- register the 0.14 subphase plan
- preserve the Phase 0.12 → Phase 0.13 → Phase 0.14 narrative
- perform no Go code changes

### 0.14.1 — Correlation Model (Request ID / Trace) ✅ Completed

Introduce a request correlation model.

Scope:

- define request ID ownership
- propagate request ID through context
- integrate correlation at the HTTP boundary
- preserve public behavior

### 0.14.2 — Logging Standardization ✅ Completed

Standardize logging conventions for observable runtime paths.

Scope:

- define log context shape
- include correlation metadata
- apply conventions first at HTTP and provider boundaries
- avoid changing business logic

### 0.14.3 — Error Context Enrichment ✅ Completed

Improve internal diagnostic metadata while preserving the public error contract.

Scope:

- enrich internal error context
- keep public error shape stable
- avoid leaking internal implementation details to clients
- align with the standardized error model introduced earlier in Stage 0

### 0.14.4 — Flow Tracing Integration ⬜ Pending

Instrument key runtime flow transitions.

Scope:

- trace HTTP → Provider → Application → Domain → Repository movement
- log meaningful step transitions
- avoid noisy or duplicated logs
- keep trace points diagnostic rather than behavioral

### 0.14.5 — Diagnostics Surface Exposure ✅ Completed

Expose a minimal diagnostics-oriented surface where compatible with the current architecture.

Scope:

- surface observability readiness or diagnostic state
- avoid introducing external metrics systems
- avoid changing existing public API semantics
- keep the surface minimal and foundation-oriented

### 0.14.6 — Validation & Documentation ⬜ Pending

Validate compatibility and close the phase narratively.

Scope:

- run `go test ./...`
- validate that public behavior remains unchanged
- document final observability decisions
- update trunk documentation consistently
- close Phase 0.14 only after implementation validation

---

## Compatibility Rules

- No public API contract changes.
- No business logic changes.
- No provider ownership regression.
- No direct bypass of the Provider Layer.
- No external metrics dependency.
- No dashboard dependency.
- No undocumented diagnostics behavior.

---

## Documentation Rules

Every 0.14 documentation update must preserve the inherited narrative:

1. context inherited from Phase 0.12 and Phase 0.13
2. real observability problem being solved
3. decision taken
4. concrete change made
5. observable impact

Documentation must remain cumulative and trunk-safe.

---

## 0.14.1 Implementation Result — Correlation Model

### Context inherited

0.14.0 locked the corrected 0.14 direction as Observability & Diagnostics Foundation. The backend entered this subphase after Phase 0.13 had consolidated provider ownership and after the 0.14 lock had clarified that observability work must not change public contracts or business behavior.

The real code already contained an HTTP request correlation seed in `internal/core/httpx`: the request middleware accepted `X-Request-Id`, generated a UUID when absent, echoed the effective value in the response header and attached it to the request context.

### Problem

That behavior was useful but still implicit. The request ID was stored behind a private context key and read directly inside middleware-adjacent code. Without a stable accessor, later 0.14 work would either duplicate direct context access patterns or introduce a parallel correlation concept.

### Decision

0.14.1 keeps request ID ownership inside `internal/core/httpx` and exposes only a safe accessor:

```go
func RequestIDFromContext(ctx context.Context) string
```

The context key remains private. Missing or nil contexts return an empty string. This preserves encapsulation while making correlation reusable by later logging, error and tracing subphases.

### Concrete change

- `internal/core/httpx/middleware.go` now exposes `RequestIDFromContext`.
- `AccessLog` reads the request ID through the accessor.
- `Recoverer` reads the request ID through the accessor.
- `internal/core/httpx/middleware_test.go` validates inherited, generated and missing request ID cases.

### Observable impact

The backend now has a stable internal request correlation seam. Public HTTP/API contracts remain unchanged: `X-Request-Id` behavior is preserved, response payloads are unchanged and no provider/domain behavior was modified.

### Validation note

The implementation was formatted with `gofmt`. Full `go test ./...` could not be completed in this environment because the Go toolchain attempted to download Go 1.25.0 from `proxy.golang.org` and DNS/network access was unavailable.

---

## 0.14.2 Implementation Result — Logging Standardization

### Context inherited

0.14.1 established a stable request correlation seam through `httpx.RequestIDFromContext(ctx)` while preserving the existing `X-Request-Id` behavior. That gave the backend a safe internal way to read correlation metadata without exposing the private HTTP context key.

The next problem was not the absence of JSON logs: the backend already used `slog.JSONHandler` through `internal/core/logger`. The real gap was naming and reuse. HTTP access and panic logs still emitted the request identifier as `rid`, while the rest of the observability plan uses `request_id` as the canonical correlation field.

### Problem

If later subphases continued to emit request correlation manually, log records would drift between `rid`, direct context access and ad-hoc attribute construction. That would make diagnostics harder exactly when 0.14 is supposed to make request movement easier to follow.

### Decision

0.14.2 keeps the existing logger package and JSON handler, but introduces a small standardization layer around request correlation attributes.

The canonical log field is:

```text
request_id
```

The logger package now owns the reusable key and helper methods, while HTTP remains the owner of extracting the effective request ID from context.

### Concrete change

- `internal/core/logger/logger.go` now exposes `RequestIDKey`.
- `internal/core/logger/logger.go` now exposes `AttrsWithRequestID(requestID string)`.
- `internal/core/logger/logger.go` now exposes `WithRequestID(log, requestID)` for future request-scoped logger derivation.
- `internal/core/httpx/middleware.go` now emits `request_id` instead of `rid` in access logs.
- `internal/core/httpx/middleware.go` now emits `request_id` instead of `rid` in panic logs.
- `internal/core/logger/logger_test.go` validates the standard request ID attribute behavior.

### Observable impact

Runtime logs now have a stable request correlation field name that aligns with the rest of Phase 0.14. This improves grep/searchability and prepares 0.14.3 error context enrichment and 0.14.4 flow tracing without changing HTTP responses, JSON payloads, provider behavior or business logic.

### Validation note

The implementation was formatted with `gofmt`. Full `go test ./...` could not be completed in this environment because the Go toolchain attempted to download Go 1.25.0 from `proxy.golang.org` and DNS/network access was unavailable.

---

## 0.14.3 Implementation Result — Error Context Enrichment

### Context inherited

0.14.1 created a stable request correlation seam through `httpx.RequestIDFromContext(ctx)`. 0.14.2 then standardized the runtime log correlation key as `request_id`. Together, those subphases gave the backend a consistent correlation concept, but the error model still needed a controlled way to carry diagnostic context.

The existing error layer was already stable from Phase 0.8. `AppError` carried `Code`, `Message`, `Status`, `Category`, `Details` and optional `Cause`, and the public envelope was already normalized as `{ error: { code, message, details } }`.

### Problem

The backend needed richer diagnostics without breaking that public contract. Before this subphase, enrichment was possible only through map-level `WithDetails`, and `ToResponseError()` reused the internal details map directly. That was functional, but not ideal for observability because later tracing work could accidentally mutate internal error details through response construction or introduce inconsistent keys for request correlation.

### Decision

0.14.3 keeps the existing error model and adds small, explicit enrichment helpers:

```go
func (e *AppError) WithContext(key string, value any) *AppError
func (e *AppError) WithRequestID(requestID string) *AppError
func (e *AppError) PublicDetails() map[string]any
```

The canonical request correlation key remains `request_id`. Empty keys or empty request IDs do not add diagnostic data. Public details are copied before response serialization.

### Concrete change

- `internal/core/errs/app_error.go` now exposes `WithContext`.
- `internal/core/errs/app_error.go` now exposes `WithRequestID`.
- `internal/core/errs/app_error.go` now exposes `PublicDetails`.
- `WithDetails` now merges on top of a copied public details map.
- `ToResponseError` now serializes copied details through `PublicDetails`.
- `internal/core/errs/app_error_test.go` validates non-mutating context enrichment, request ID enrichment, empty request ID behavior, public details copying and response details copy isolation.

### Observable impact

Errors can now be enriched with request correlation and controlled diagnostic metadata without changing the public JSON envelope, HTTP status behavior, provider behavior or business logic. This prepares 0.14.4 flow tracing integration by giving the flow instrumentation a safe error-context surface.

### Validation note

The implementation was formatted manually in standard Go style. Full `go test ./...` could not be completed in this environment because the Go toolchain attempted to download Go 1.25.0 from `proxy.golang.org` and DNS/network access was unavailable.

---

## 0.14.4 Implementation Result — Flow Tracing Integration

### Context inherited

0.14.1 established request correlation through `httpx.RequestIDFromContext(ctx)`. 0.14.2 standardized the canonical log key as `request_id`. 0.14.3 then added safe diagnostic enrichment to the error layer without changing the public error envelope.

Together, those subphases made correlation and diagnostic context available, but the runtime flow still lacked explicit lifecycle events. The backend could log an HTTP request after it finished, but it did not yet expose a minimal start/end trace that made request movement observable as a flow.

### Problem

Debugging still required reading isolated log entries and inferring flow order manually. Introducing OpenTelemetry, Prometheus, dashboards or a diagnostics API in this subphase would exceed the defined scope, but leaving the HTTP/application lifecycle without explicit flow markers would postpone the core value of Phase 0.14.

### Decision

0.14.4 introduces a minimal internal flow tracing convention over the existing `slog` JSON logger:

```text
message: flow_trace
field: flow_event
```

The logger package owns the reusable flow event key and helper. HTTP remains the owner of request correlation extraction, and application lifecycle events can use the same helper without a request ID.

### Concrete change

- `internal/core/logger/logger.go` now exposes `FlowEventKey`.
- `internal/core/logger/logger.go` now exposes `AttrsWithFlowEvent(requestID, event, attrs...)`.
- `internal/core/logger/logger_test.go` validates flow event attributes with and without request IDs.
- `internal/core/httpx/middleware.go` now emits `http_request_start` flow traces before handler execution.
- `internal/core/httpx/middleware.go` now emits `http_request_end` flow traces after handler execution.
- `internal/app/app.go` now emits `application_start` and `application_stop` flow traces.

### Observable impact

Existing JSON logs can now describe request and application movement through consistent flow events. This keeps the backend contract stable while making the HTTP lifecycle easier to follow with the previously established `request_id` correlation model.

No public API payload, error envelope, provider contract, business behavior, external metrics system, dashboard or diagnostics endpoint was introduced.

### Validation note

The code was prepared in standard Go formatting style. Full `go test ./...` could not be completed in this environment because the local Go toolchain is older than the module requirement and automatic toolchain download is unavailable.

---

## 0.14.5 Implementation Result — Diagnostics Surface Exposure

### Context inherited

0.14.1 established request correlation through `httpx.RequestIDFromContext(ctx)`. 0.14.2 standardized the runtime log key as `request_id`. 0.14.3 added safe diagnostic enrichment to the error layer. 0.14.4 then introduced flow tracing events over the existing structured logger.

Together, those subphases made observability available internally, but the backend still lacked a minimal runtime surface that could report the active observability foundation capabilities without asking operators to infer them from logs or code.

### Problem

The existing operational endpoints were intentionally narrow:

- `/health` reports service liveness.
- `/readiness` reports dependency readiness.
- `/version` reports build identity.

None of them exposed whether request correlation, structured logging, error context enrichment or flow tracing were available. Adding Prometheus, OpenTelemetry, dashboards or counters at this stage would exceed the phase scope.

### Decision

0.14.5 exposes a minimal diagnostics surface:

```http
GET /diagnostics
```

The endpoint reports service identity, environment, version, commit, timestamp and a foundation-level `observability` object.

### Concrete change

- `internal/core/status.Service.Diagnostics()` now returns the canonical diagnostics payload.
- `internal/core/status/service_test.go` validates the observability capability flags.
- `internal/core/httpx.NewRouter` now registers `GET /diagnostics` alongside the existing operational endpoints.
- `internal/core/httpx/router_diagnostics_test.go` validates the endpoint response shape.

### Observable impact

The backend now exposes a small diagnostics snapshot confirming the active 0.14 observability capabilities:

- `request_correlation`
- `structured_logging`
- `error_context_enrichment`
- `flow_tracing`

No public business API payload, authentication contract, provider interface, domain behavior, external metrics backend, OpenTelemetry integration or dashboard was introduced.

### Validation note

The code was prepared in standard Go formatting style. Full `go test ./...` could not be completed in this environment because the local Go toolchain is older than the module requirement and automatic toolchain download is unavailable.
