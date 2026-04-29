# Phase 0.12.4 — Mapping Layer Introduction

## Subphase 0.12.4.3 — Mapping Consolidation

---

## Objective

Consolidate runtime mapping usage through the centralized module-level mapping layer introduced in Phase 0.12.4.2.

This subphase redirects application and HTTP call sites away from model-owned mapper functions and toward the explicit module mapping packages:

```text
internal/modules/auth/mappers
internal/modules/user/mappers
internal/modules/usersettings/mappers
```

The goal is to reduce duplicated mapping responsibility while preserving all existing HTTP contracts and runtime behavior.

---

## Scope

### Included

* align application-layer read projections with centralized mappers
* align authenticated and public auth handlers with centralized write mappers
* preserve existing request and response structs
* preserve existing JSON contracts
* keep legacy mapper methods available temporarily for compatibility

### Excluded

* removing legacy mapper files
* deleting compatibility methods
* changing endpoint contracts
* changing request payloads
* changing response payloads
* changing business logic

---

## Consolidation Strategy

The consolidation is intentionally progressive.

Phase 0.12.4.2 created the centralized mapping packages but left previous mapper call sites untouched.

Phase 0.12.4.3 changes call sites to use the centralized packages directly:

```text
readmodels.FromX(...)      → mappers.XToReadModel(...)
writemodel.ToDomainInput() → mappers.XWriteToDomainInput(...)
```

This keeps mapping ownership in the mapping layer without removing compatibility functions yet.

---

## Application Layer Alignment

The auth application layer now uses centralized read mapping helpers for:

* login response projection
* bootstrap user projection
* bootstrap settings projection
* wallet challenge projection
* wallet identity projection
* wallet collection projection through existing helper boundaries

The application layer continues to own orchestration only.

It no longer directly depends on read model mapper functions for newly aligned projections.

---

## Handler Alignment

Auth HTTP handlers now use centralized write mapping helpers for:

* login input
* profile update input
* settings update input
* wallet challenge input
* wallet verification input
* wallet link challenge input
* wallet link verification input
* account merge challenge input
* account merge verification input
* wallet detach check input
* wallet detach execute input
* primary wallet update input

Handlers still decode the same request payloads and still call the same application or service methods.

Only the mapping call path changed.

---

## Compatibility

Compatibility is preserved because:

* request type aliases remain unchanged
* response structs remain unchanged
* JSON tags remain unchanged
* endpoint wiring remains unchanged
* legacy mapper methods remain available
* only mapper invocation ownership changed

---

## Files Updated

### Go Files

```text
internal/modules/auth/app/application.go
internal/modules/auth/app/support_helpers.go
internal/modules/auth/http_login.go
internal/modules/auth/http_wallet.go
```

### Documentation

```text
docs/phase0_12_4_mapping_consolidation.md
```

---

## Validation Notes

The change was formatted with `gofmt`.

The execution environment used for generation could not complete `go test ./...` because the project requires Go `1.25.0` and the toolchain download is blocked in the sandbox environment.

Expected local validation command:

```bash
go test ./...
```

---

## Runtime Contract

No endpoint contract is changed in this subphase.

The following remain stable:

* login payloads
* authenticated profile update payloads
* settings update payloads
* wallet challenge payloads
* wallet verify payloads
* wallet management payloads
* bootstrap responses
* session/profile responses
* wallet responses

---

## Status

Phase: 0.12.4  
Subphase: 0.12.4.3  
Status: COMPLETED  
Code Impact: CONSOLIDATION ONLY  
Contract Impact: NONE  

---

## Next Step

```text
0.12.4.4 — Application Refactor
```
