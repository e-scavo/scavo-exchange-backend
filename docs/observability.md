# Observability

## Objective

This document defines the observability baseline for the SCAVO Exchange backend.

Its purpose is to ensure that the system remains diagnosable as it grows from a minimal bootstrap into a blockchain-integrated, persistence-backed, real-time backend.

---

## Why Observability Matters Early

The backend will progressively depend on:

- PostgreSQL
- Redis
- SCAVIUM RPC
- smart contract interactions
- event ingestion
- transaction tracking
- background jobs
- internal operational flows

If observability is deferred too long, infrastructure and product phases become harder to diagnose, validate, and operate safely.

For this reason, observability is considered part of the project foundation.

---

## Observability Baseline

The baseline observability direction includes:

- structured logs
- request correlation
- process-level health visibility
- readiness visibility
- infrastructure failure visibility
- metrics direction
- future tracing support

This baseline is sufficient for early and mid-stage backend growth.

---

## Structured Logging

Structured logging is the default logging model.

Expected goals:

- machine-readable output
- consistent fields
- clear severity levels
- meaningful operational context
- no leaked secrets

Logging should help identify:

- startup failures
- config issues
- request failures
- authentication failures
- infrastructure failures
- external dependency failures
- unexpected panics or recoveries

---

## Logging Principles

Logging should be:

- intentional
- concise
- consistent
- useful for diagnosis
- safe from secret leakage

Logging should avoid:

- excessive noise
- duplicate logs for the same failure without value
- raw secret exposure
- transport-only clutter without operational meaning

---

## Request Correlation

As the backend grows, request and operation correlation becomes increasingly important.

The observability direction should support correlation across:

- HTTP request lifecycle
- WebSocket action handling
- service orchestration
- future DB interactions
- future chain interactions
- future background jobs

This does not require full tracing yet, but the design should preserve room for correlation identifiers and structured propagation.

---

## Health and Readiness

The backend must distinguish between two operational states.

### Health
Represents whether the process is alive and can respond.

This is the minimum operational state.

### Readiness
Represents whether the system is actually ready to serve its intended workload.

Once the backend depends on infrastructure, readiness may require checks against:

- PostgreSQL
- Redis
- required config presence
- migration state later
- chain connectivity later when appropriate

A healthy process may still be unready.

This distinction is mandatory for safe infrastructure growth.

---

## Metrics Direction

Metrics are part of the official direction, even if not yet implemented.

Metrics should eventually provide visibility into areas such as:

- request volume
- request latency
- error rates
- auth failures
- WebSocket connection counts
- DB connectivity failures
- Redis connectivity failures
- chain RPC errors
- background job behavior
- transaction tracking outcomes

The initial goal is not full coverage.

The initial goal is to define a direction that keeps instrumentation coherent.

---

## Tracing Direction

Distributed tracing is not required at the current stage.

However, the system should preserve room for future tracing support, especially because later phases may involve:

- DB calls
- Redis calls
- chain RPC calls
- contract-related reads
- background workers
- multi-step transaction tracking

The architecture should avoid blocking this evolution.

---

## Observability Boundaries

Observability support should primarily live in reusable infrastructure rather than ad hoc implementation inside every module.

Preferred locations include future packages such as:

- `internal/core/observability`
- middleware-level instrumentation
- shared logger/context helpers

This avoids inconsistent instrumentation patterns across modules.

---

## Failure Visibility Expectations

The backend should eventually make the following failures visible in a clear way:

- invalid configuration
- startup failure
- dependency unavailability
- DB connection issues
- Redis connection issues
- migration issues
- RPC degradation
- contract integration failures
- job failures
- unexpected panic recovery

This does not mean every failure must be surfaced identically, but each category should become diagnosable.

---

## Operational Diagnostics Direction

Operational diagnostics should progressively support:

- quick startup validation
- infrastructure dependency checks
- environmental misconfiguration detection
- basic runtime insight
- later internal team troubleshooting

These diagnostics are especially important for internal testing environments.

---

## Non-Goals for the Current Stage

This stage does not yet require:

- full metrics implementation
- full tracing implementation
- dashboards
- alerting stack
- centralized log aggregation
- production SLO definition

Those may come later once the implementation baseline is in place.

---

## Recommended Next Step

The next recommended step is:

Phase 0.3.1 - Implementation Bootstrap for Persistence and Health Infrastructure

That phase should begin translating observability direction into concrete building blocks such as:

- health endpoint evolution
- readiness direction
- infrastructure-aware startup behavior
- reusable observability scaffolding

### Phase 0.9 Impact — API Versioning
The introduction of `/api/v1` as canonical API requires developers to consider versioned endpoints during development and testing.

---

## Phase 0.13 — Provider Layer And Observability Boundary

Phase 0.13 is Provider Layer Consolidation, not an observability implementation phase. The earlier observability direction remains valid as a future infrastructure concern, but this phase intentionally avoids adding logging, tracing, metrics or diagnostics behavior.

This distinction matters because the roadmap previously contained an observability-oriented 0.13 direction before the phase was corrected to Provider Layer Consolidation. The correction does not discard observability as a future concern; it changes the order of work. Provider boundaries first make later instrumentation safer because they identify stable orchestration entry points.

### Observability implication

Provider consolidation improves future observability readiness because orchestration now has a clearer boundary. Later diagnostics work can attach observations around provider entry points without requiring handlers to expose lower-level service/store dependency details.

### Constraint

No new observability behavior is introduced by Phase 0.13. Existing logs, middleware behavior, error responses and runtime diagnostics remain unchanged.

### Phase 0.13 subphase state

- 0.13.0 ✔ Definition & Documentation Lock
- 0.13.1 ✔ Provider Inventory & Classification
- 0.13.2 ✔ Provider Interface Design
- 0.13.3 ✔ Provider Implementation
- 0.13.4 ✔ Application Integration
- 0.13.5 ✔ Validation & Compatibility
- 0.13.6 ✔ Documentation & Closure

Phase 0.13 is complete. Any future observability phase should treat provider entry points as natural instrumentation candidates, while preserving public API behavior.



---

# Phase 0.14 — Observability & Diagnostics Foundation

## Context

The observability document existed before Phase 0.14 as a baseline description of diagnostic expectations. After Phase 0.13, the backend has a consolidated Provider Layer and a cleaner runtime composition path, but the actual execution path still needs stronger request-level visibility.

Phase 0.14 makes observability an explicit Stage 0 foundation concern.

## Problem

The current backend can validate behavior through tests and logs, but the diagnostic path is still fragmented:

- a request cannot yet be followed consistently through all internal boundaries
- log entries are not standardized around a shared correlation model
- internal errors do not consistently carry diagnostic context
- provider/application/domain transitions are not yet visible as a coherent execution flow

## Decision

Phase 0.14 introduces observability without changing behavior.

The phase will focus on request correlation, structured logging, internal error context enrichment, flow tracing and minimal diagnostics exposure. It will not introduce external metrics systems, dashboards, Prometheus, OpenTelemetry or API contract changes.

## Expected Impact

The expected outcome is a backend that keeps the same public contract while becoming easier to debug and operate internally.

The diagnostic path should make the following movement observable:

```text
HTTP → Provider → Application → Domain → Repository
```

## 0.14 Subphase Alignment

- 0.14.0 — Phase Definition & Documentation Lock ✅ Completed
- 0.14.1 — Correlation Model (Request ID / Trace) ✅ Completed
- 0.14.2 — Logging Standardization ✅ Completed
- 0.14.3 — Error Context Enrichment ✅ Completed
- 0.14.4 — Flow Tracing Integration ✅ Completed
- 0.14.5 — Diagnostics Surface Exposure ⬜ Pending
- 0.14.6 — Validation & Documentation ⬜ Pending

---

## Phase 0.14.1 Correlation Model Result

0.14.1 consolidated the correlation model that already existed at the HTTP middleware boundary instead of replacing it with a new runtime concept.

Context inherited from 0.14.0: the phase lock established that observability must improve visibility without changing public behavior. The real code already generated or propagated `X-Request-Id`, stored the value in `context.Context` and emitted it from access and panic logs, but that behavior was still implicit and locally coupled to middleware internals.

Problem addressed: request correlation could not be safely reused outside the middleware implementation because the context key remained private and there was no exported accessor. This made future logging standardization, error enrichment and flow tracing more likely to duplicate context access patterns or bypass the correlation boundary.

Decision taken: `internal/core/httpx` now owns a small public accessor, `RequestIDFromContext(ctx context.Context) string`, while keeping the context key private. The middleware continues to accept an incoming `X-Request-Id` header or generate one when absent, echoes the effective value to the response header and attaches it to the request context.

Concrete change: `AccessLog` and `Recoverer` now read the request identifier through the accessor instead of reaching directly into the context key. Dedicated middleware tests validate inherited request IDs, generated request IDs and empty results for missing or nil contexts.

Observable impact: the backend now has a stable internal request correlation seam for the next 0.14 subphases. Public API contracts, response payloads and business behavior remain unchanged.

---

## Phase 0.14.2 Logging Standardization Result

0.14.2 standardized the request correlation field used by runtime logs without changing the existing logger backend or public behavior.

Context inherited from 0.14.1: request correlation is owned by `internal/core/httpx`, and the supported read path is `httpx.RequestIDFromContext(ctx)`. The logger already emitted JSON through `slog`, so this subphase did not replace the logging backend.

Problem addressed: HTTP access and panic logs still used `rid`, while Phase 0.14 needs a stable cross-cutting key that future diagnostics can reuse consistently.

Decision taken: `request_id` is the canonical log attribute key. The logger package now exposes the reusable key and helpers, while HTTP middleware continues to extract the effective request ID from context.

Concrete change: `AccessLog` and `Recoverer` now build attributes with `logger.AttrsWithRequestID(rid)`, and `internal/core/logger/logger_test.go` validates request ID attribute behavior.

Observable impact: logs produced by the HTTP boundary now use `request_id`, making correlation consistent for later error context enrichment and flow tracing. No HTTP contract, response payload or business behavior changed.

---

## Phase 0.14.3 Error Context Enrichment Result

0.14.3 enriched the internal error model with safe diagnostic context while preserving the public error contract established earlier in Stage 0.

Context inherited from 0.14.2: request correlation now uses the canonical `request_id` key in logs, and the logger package owns reusable helpers for request correlation attributes. The next observability gap was the error layer, where diagnostic metadata needed a controlled path before flow tracing could be added.

Problem addressed: `AppError` already carried `Details`, but callers had only coarse map-based enrichment through `WithDetails`, and `ToResponseError()` exposed the same map reference to response construction. That made future diagnostics more likely to mutate shared state accidentally or add ad-hoc metadata without a standard key.

Decision taken: `internal/core/errs` now exposes small, explicit helpers for contextual enrichment. `WithContext` adds one diagnostic detail, `WithRequestID` applies the canonical `request_id` key, and `PublicDetails` returns a safe copy for public serialization.

Concrete change: `ToResponseError()` now uses `PublicDetails()` so public response construction receives a copied details map. Dedicated tests validate request ID enrichment, non-mutating context enrichment and response detail copy behavior.

Observable impact: error diagnostics can now carry request correlation and other controlled metadata safely, while the public error envelope remains unchanged.

---

## Phase 0.14.4 Flow Tracing Integration Result

0.14.4 integrates a minimal flow tracing convention into the existing structured logging path without introducing a tracing backend or changing public contracts.

Context inherited from 0.14.3: request correlation is available through the canonical `request_id` field, structured logging is centralized through `internal/core/logger`, and errors can carry copied diagnostic context safely. The remaining gap was that runtime movement still appeared mostly as isolated log records instead of lifecycle events that describe where a request or process is in the backend flow.

Problem addressed: the backend needed explicit start/end markers for observable flow movement, but adding external tracing, metrics or diagnostics endpoints at this point would be premature and outside 0.14.4 scope.

Decision taken: flow tracing is represented as structured log records with message `flow_trace` and canonical field `flow_event`. The logger package owns `FlowEventKey` and `AttrsWithFlowEvent`, while HTTP continues to own request correlation extraction.

Concrete change: HTTP access logging now emits `http_request_start` before handler execution and `http_request_end` after handler execution, preserving method, path, remote address, status, bytes and duration context. Application lifecycle logs now emit `application_start` and `application_stop` through the same flow tracing convention.

Observable impact: operators can now follow request and lifecycle movement through consistent flow events using existing JSON logs. No HTTP response, public API contract, provider behavior, business logic, metrics backend or diagnostics endpoint changed.
