# Phase 0.12.5 — Contract Alignment

## Subphase 0.12.5.4 — Handler Contract Adjustment

---

## Objective

Adjust handler-level contract ownership so HTTP compatibility contracts are explicit, centralized and easier to audit, without changing public routes, JSON payloads, status codes or runtime behavior.

This subphase continues the contract alignment sequence introduced in Phase 0.12.5 by separating handler contract declarations from handler execution logic.

---

## Baseline

Input baseline:

```text
scavo-exchange-backend-0.12.5.3.zip
```

Relevant previous work:

- Phase 0.12.2 — Read Model Extraction
- Phase 0.12.3 — Write Model Isolation
- Phase 0.12.4 — Mapping Layer Introduction
- Phase 0.12.5.1 — Contract Inventory & Classification
- Phase 0.12.5.2 — Contract Normalization Design
- Phase 0.12.5.3 — Contract Alignment Implementation

---

## Scope

### Included

- centralize auth HTTP compatibility contract declarations
- keep public request names stable
- keep public response names stable
- keep request aliases pointing to explicit write models
- keep response contracts pointing to read models or stable application responses
- remove contract declaration clutter from handler execution files

### Excluded

- route changes
- JSON field changes
- API version changes
- business logic changes
- repository changes
- mapper behavior changes
- validation rule changes

---

## Files Changed

### Added

```text
internal/modules/auth/http_contracts.go
```

Purpose:

- centralize HTTP compatibility contract declarations
- make request/write-model ownership explicit
- make response/read-model ownership explicit
- keep handler files focused on request execution flow

### Updated

```text
internal/modules/auth/http_login.go
internal/modules/auth/http_wallet.go
```

Purpose:

- remove inline request/response contract declarations from execution-oriented handler files
- keep handler behavior unchanged
- keep existing handler method names unchanged
- keep decode/write behavior unchanged

---

## Contract Ownership After This Subphase

Handler-facing public contract names remain stable.

Examples:

```text
LoginRequest
UpdateMeRequest
WalletChallengeRequest
WalletVerifyRequest
WalletLinkChallengeRequest
WalletDetachExecuteRequest
WalletPrimarySetRequest
```

These names are now centralized in:

```text
internal/modules/auth/http_contracts.go
```

Request contracts continue to point to explicit write models:

```text
LoginRequest → AuthLoginWriteModel
WalletVerifyRequest → AuthWalletVerifyWriteModel
WalletPrimarySetRequest → AuthWalletPrimarySetWriteModel
```

Response contracts continue to preserve public JSON shape while relying on read models or stable application responses where appropriate:

```text
MeResponse → UserReadModel
MeSettingsResponse → UserSettingsReadModel
WalletChallengeResponse → AuthWalletChallengeReadModel
WalletLinkVerifyResponse → app.WalletLinkVerifyResponse
```

---

## Runtime Behavior

No runtime behavior is intentionally changed.

The following behavior remains unchanged:

- HTTP routes
- request body decoding
- maximum request sizes
- auth claim enforcement
- error handling
- response status codes
- JSON field names
- wallet challenge and verification flow
- profile update flow
- settings update flow
- wallet link, merge, detach and primary wallet flows

---

## Compatibility Notes

This subphase uses Go type aliases and stable response structs to preserve compatibility.

This means existing callers continue to interact with the same public names and JSON contracts while the internal source of truth becomes explicit:

```text
HTTP Compatibility Contract → Write Model → Domain Input
Domain/Application Response → Read Model / Stable Response → HTTP Compatibility Contract
```

---

## Mapping Layer Relationship

This subphase does not move or duplicate mapping logic.

Mapping ownership remains under:

```text
internal/modules/auth/mappers/
internal/modules/user/mappers/
internal/modules/usersettings/mappers/
```

The new `http_contracts.go` file owns compatibility naming only.

It does not own transformation behavior.

---

## Validation Expectations

The implementation must satisfy:

```bash
gofmt ./internal/modules/auth/http_contracts.go ./internal/modules/auth/http_login.go ./internal/modules/auth/http_wallet.go
go test ./...
```

The execution environment used to generate this patch could run `gofmt`, but could not complete `go test ./...` because the local Go toolchain is older than the project requirement and network access for toolchain download is blocked.

The user-side validation must remain the authoritative test result for this subphase.

---

## Result

Handler contract declarations are now isolated from handler execution flow.

This improves:

- contract traceability
- future contract reviews
- handler readability
- alignment with explicit read/write model separation

Public behavior remains unchanged.

---

## Status

Phase: 0.12.5
Subphase: 0.12.5.4
Status: Completed
Code impact: Handler contract organization only
Runtime behavior impact: None intended
Next: 0.12.5.5 — Validation & Compatibility
