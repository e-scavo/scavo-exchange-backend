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