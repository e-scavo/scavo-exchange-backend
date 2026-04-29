# Phase 0.12.4 — Mapping Layer Introduction

## Subphase 0.12.4.4 — Application Refactor

---

## Objective

Refactor application-layer mapping usage so that mapping behavior introduced and consolidated in Phase 0.12.4 is consumed through the explicit module-level mapping layer.

This subphase focuses on reducing residual mapping logic inside application orchestration while preserving runtime behavior, HTTP contracts, response payloads, and existing tests.

---

## Scope

### Included

* application-layer call-site cleanup
* replacement of residual direct wallet read-model mapping with centralized mapper functions
* preservation of response actionability metadata
* preservation of existing handlers and HTTP contracts
* documentation of the refactor boundary

### Excluded

* endpoint changes
* response contract changes
* request contract changes
* removal of legacy compatibility wrappers
* business logic changes
* test contract rewrites

---

## Implementation Summary

The implementation moves remaining wallet read-model actionability mapping out of application helper logic and into the centralized auth mapping layer.

The auth mapping layer now owns:

* wallet identity to wallet read-model mapping
* wallet identity list to wallet read-model list mapping
* wallet actionability enrichment
* actionability-aware wallet list mapping

Application code now consumes these explicit mapper functions instead of carrying mapping-oriented helper logic internally.

---

## Files Added or Updated

### Updated

* `internal/modules/auth/mappers/read.go`
* `internal/modules/auth/app/application.go`
* `internal/modules/auth/app/support_helpers.go`
* `internal/modules/auth/http_wallet_list.go`

### Added

* `docs/phase0_12_4_application_refactor.md`

---

## Mapping Layer Changes

### New centralized functions

The auth mapping layer now exposes:

```go
WalletIdentitiesToActionableReadModels(wallets []*authdomain.WalletIdentity) []*authreadmodels.AuthWalletReadModel
```

and:

```go
EnrichWalletReadModelsActionability(wallets []*authreadmodels.AuthWalletReadModel) []*authreadmodels.AuthWalletReadModel
```

These functions centralize the previously application-owned response actionability mapping.

---

## Application Layer Changes

### `internal/modules/auth/app/application.go`

Application flows now call the auth mapper package directly for actionability-aware wallet read-model mapping.

This affects wallet-related response construction while preserving the same JSON contract.

### `internal/modules/auth/app/support_helpers.go`

Residual wallet mapping helpers were removed from application support helpers.

The file now delegates wallet read-model generation to:

```go
authmappers.WalletIdentitiesToActionableReadModels
```

Profile wallet view creation still remains application-local because `ProfileWalletView` is an app response shape and moving it into `mappers` would introduce an undesirable package dependency direction.

---

## Root Auth Compatibility Surface

### `internal/modules/auth/http_wallet_list.go`

The root auth compatibility surface now delegates wallet mapping behavior to the centralized auth mapper package.

Compatibility helper functions remain present as a stable local facade, but the mapping implementation is no longer duplicated there.

This preserves existing tests and callers while reducing duplicate mapping logic.

---

## Compatibility Rules Preserved

The refactor preserves:

* HTTP response field names
* JSON tags
* wallet list pagination behavior
* wallet filter behavior
* wallet sort behavior
* wallet actionability flags
* detach block reasons
* handler entry points
* existing test expectations

---

## Contract Preservation

No HTTP contract changes were introduced.

The following response surfaces remain compatible:

* login response
* session response
* bootstrap response
* `/me` response
* wallet list response
* wallet link response
* wallet merge response
* wallet primary response
* wallet detach response

---

## Mapping Direction After This Subphase

The enforced mapping flow remains:

```text
Write Model → Domain Input
Domain State → Read Model
```

This subphase improves the second half of that flow by ensuring wallet read-model actionability enrichment also belongs to the centralized mapping layer.

---

## Non-Goals

This subphase does not remove all legacy mapper functions from `readmodels/` or `writemodels/`.

Those remain temporarily available for backward compatibility until final consolidation and closure steps decide whether they should be retained, deprecated, or removed.

---

## Validation Notes

The implementation was formatted with `gofmt`.

The local assistant environment could not execute `go test ./...` because the project requires Go `1.25.0` and the sandbox cannot download the required toolchain.

The expected validation command remains:

```bash
go test ./...
```

---

## Status

Phase: 0.12.4  
Subphase: 0.12.4.4  
Status: COMPLETED  
Code Impact: APPLICATION REFACTOR  
HTTP Contract Impact: NONE  

---

## Next Step

```text
0.12.4.5 — Validation & Compatibility
```
