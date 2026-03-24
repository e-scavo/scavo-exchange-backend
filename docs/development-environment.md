# Development Environment

## Objective

This document defines the baseline local development environment direction for the SCAVO Exchange backend.

Its purpose is to make local setup reproducible, explicit, and ready for future implementation without introducing unnecessary operational complexity in the early phases.

---

## Environment Goal

The initial development environment must support a deterministic backend foundation.

The minimum environment target is:

- backend application
- PostgreSQL
- Redis

This is the minimum baseline required to prepare the backend for real persistence, cache-aware infrastructure, and future indexing or job-related features.

---

## Environment Philosophy

The development environment should be:

- explicit
- reproducible
- easy to onboard
- close enough to later internal testing environments
- simple enough for early development

The goal is not to mirror full production at this stage.

The goal is to avoid hidden machine-specific setup and undocumented configuration drift.

---

## Minimum Local Services

### Backend Application
The backend application remains the main executable service.

Responsibilities in local development:

- serve HTTP endpoints
- serve WebSocket endpoints
- load environment-based configuration
- later connect to PostgreSQL
- later connect to Redis
- later connect to SCAVIUM RPC

### PostgreSQL
PostgreSQL is the primary local durable datastore.

Local development responsibilities:

- persist backend durable state
- support future migrations
- support repository testing
- support deterministic schema setup

### Redis
Redis is the local ephemeral infrastructure store.

Local development responsibilities:

- support cache-oriented development
- support coordination-oriented development
- support future rate limiting or background coordination
- support experimentation without redefining persistence boundaries

---

## Configuration Direction

The local environment should be driven by explicit configuration through environment variables.

The baseline environment variable direction includes at least:

- `APP_ENV`
- `APP_NAME`
- `HTTP_ADDR`
- `JWT_SECRET`
- `JWT_ISSUER`
- `JWT_TTL_MINUTES`
- `POSTGRES_HOST`
- `POSTGRES_PORT`
- `POSTGRES_DB`
- `POSTGRES_USER`
- `POSTGRES_PASSWORD`
- `POSTGRES_SSLMODE`
- `REDIS_ADDR`
- `REDIS_PASSWORD`
- `REDIS_DB`
- `SCAVIUM_RPC_URL`
- `SCAVIUM_CHAIN_ID`
- `LOG_LEVEL`

Future phases may extend this set.

---

## Local Environment Profiles

The configuration model should leave room for distinct profiles such as:

### local
Used for active developer work.

Characteristics:

- simplified secrets
- local hostnames or container names
- verbose logging when useful
- development-safe defaults

### internal-testing
Used later for coordinated internal team testing.

Characteristics:

- more realistic configuration
- stable shared endpoints
- closer behavior to release-like environments
- controlled test data and infrastructure

### production-like
Reserved for later deployment preparation.

Characteristics:

- secured secrets handling
- stricter logging and security controls
- explicit external infrastructure integration

---

## Docker Direction

Docker-based local infrastructure is recommended.

The target development direction includes:

- one backend service
- one PostgreSQL service
- one Redis service

Optional later additions may include:

- migration runner
- seed runner
- observability helpers
- contract tooling container if needed

Docker is recommended because it improves:

- onboarding consistency
- environment reproducibility
- CI alignment later
- reduced machine-specific setup differences

---

## Manual Local Support

Docker is recommended, but manual service startup may still be supported for developers who prefer it.

However, manual setup must not become the undocumented default.

Whenever manual setup is supported, the equivalent configuration expectations should remain documented.

---

## Migration Workflow Direction

The local environment must support a migration-driven workflow.

That means:

- schema setup must become reproducible
- new environments should not require manual table creation
- schema evolution should follow versioned migration history
- local setup should eventually support applying migrations from scratch

This is required for future repository and test stability.

---

## Local Data Philosophy

Early local environments do not need large seeded datasets.

The initial local data goal is:

- valid schema
- valid connectivity
- deterministic empty-state startup
- minimal development data only when needed

Large fixture or seeded environments can come later.

---

## Health Expectations

As infrastructure is introduced, local development should make it easy to identify:

- backend startup success
- PostgreSQL connectivity success or failure
- Redis connectivity success or failure
- configuration errors
- migration errors later

This expectation aligns with the upcoming observability and readiness phase.

---

## Security Expectations for Local Development

Local development may use simplified credentials, but the model must still preserve good practices:

- no hardcoded real secrets in source
- no accidental exposure of future production secrets
- explicit environment-driven values
- clear separation between local and non-local configuration

---

## Non-Goals for the Current Stage

This document does not yet require:

- a final docker-compose file
- final migration tooling choice
- final production deployment model
- seeded blockchain test environments
- advanced observability stack
- final CI environment configuration

Those will be addressed in later phases.

---

## Recommended Next Step

The next recommended step is:

Phase 0.2.3 - Observability and Test Bootstrap

That phase should define:

- health and readiness expectations
- observability baseline
- test structure direction
- error visibility and operational diagnostics baseline