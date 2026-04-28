# Phase 0.12.3 — Write Model Isolation

## Subphase 0.12.3.4 — Handler Alignment

---

## Objective

Align HTTP handlers with the explicit Write Model layer introduced in Phase 0.12.3.2 and mapped in Phase 0.12.3.3.

This subphase introduces controlled handler-level usage of Write Models while preserving existing HTTP contracts, response shapes, endpoint behavior, validation paths, and runtime compatibility.

---

## Scope

### Included

* align auth login request decoding with `AuthLoginWriteModel`
* align authenticated profile update decoding with `AuthUpdateProfileWriteModel`
* align authenticated settings update decoding with `AuthUpdateSettingsWriteModel`
* align wallet challenge and verification request decoding with wallet Write Models
* align authenticated wallet management request decoding with wallet Write Models
* use explicit `ToDomainInput` mapping before calling existing application/service boundaries
* preserve existing request aliases and JSON field names
* preserve current error handling behavior

### Excluded

* endpoint changes
* HTTP contract changes
* response shape changes
* application service signature changes
* domain service behavior changes
* repository changes
* validation rule changes
* removal of legacy compatibility aliases

---

## Handler Alignment Strategy

The alignment is intentionally conservative.

Existing request type names are preserved as aliases to their corresponding Write Models.

This allows current handlers and tests to keep their public request names while moving the underlying input representation to the explicit write-model layer.

Pattern:

```go
LegacyRequestName = writemodels.ExplicitWriteModel
```

Each handler now performs:

```text
HTTP JSON payload → Write Model → Domain Input → existing application/service call
```

This keeps handler behavior stable while removing implicit direct usage of hybrid request structures.

---

## Files Updated

### `internal/modules/auth/http_login.go`

Aligned:

* `LoginRequest`
* `UpdateMeRequest`
* `UpdateMeSettingsRequest`

The handlers now decode into write models and call `ToDomainInput()` before delegating to existing services.

No endpoint behavior was changed.

---

### `internal/modules/auth/http_wallet.go`

Aligned:

* `WalletChallengeRequest`
* `WalletVerifyRequest`
* `WalletLinkChallengeRequest`
* `WalletLinkVerifyRequest`
* `WalletAccountMergeChallengeRequest`
* `WalletAccountMergeVerifyRequest`
* `WalletDetachCheckRequest`
* `WalletDetachExecuteRequest`
* `WalletPrimarySetRequest`

All wallet input payloads now enter the system through explicit Write Models.

The existing application methods remain unchanged.

---

## Compatibility Notes

The following remain unchanged:

* route paths
* method names
* JSON field names
* max body limits
* decode behavior
* unknown-field rejection
* error envelopes
* application service signatures
* response models
* tests expectations

---

## Write Model Usage

### Login

```text
LoginRequest → AuthLoginWriteModel → LoginInput → Application.Login
```

### Profile Update

```text
UpdateMeRequest → AuthUpdateProfileWriteModel → ProfileUpdateInput → Users.UpdateDisplayName
```

### Settings Update

```text
UpdateMeSettingsRequest → AuthUpdateSettingsWriteModel → SettingsUpdateInput → UserSettings.UpdatePreferences
```

### Wallet Challenge

```text
WalletChallengeRequest → AuthWalletChallengeWriteModel → WalletChallengeInput → WalletChallengeService.Create
```

### Wallet Verification

```text
WalletVerifyRequest → AuthWalletVerifyWriteModel → WalletVerifyInput → WalletVerificationService.VerifyAndLogin
```

### Authenticated Wallet Management

```text
WalletLinkChallengeRequest → AuthWalletLinkChallengeWriteModel → WalletChallengeInput
WalletLinkVerifyRequest → AuthWalletLinkVerifyWriteModel → WalletVerifyInput
WalletAccountMergeChallengeRequest → AuthWalletAccountMergeChallengeWriteModel → WalletChallengeInput
WalletAccountMergeVerifyRequest → AuthWalletAccountMergeVerifyWriteModel → WalletVerifyInput
WalletDetachCheckRequest → AuthWalletDetachCheckWriteModel → WalletDetachInput
WalletDetachExecuteRequest → AuthWalletDetachExecuteWriteModel → WalletDetachInput
WalletPrimarySetRequest → AuthWalletPrimarySetWriteModel → WalletPrimarySetInput
```

---

## Risk Handling

### Risk: breaking existing request names

Mitigation:

* existing request names are preserved as aliases

### Risk: changing JSON contracts

Mitigation:

* write models preserve existing JSON tags

### Risk: changing validation behavior

Mitigation:

* existing service validation is still used
* payload decode behavior remains unchanged

### Risk: premature application refactor

Mitigation:

* application method signatures remain unchanged
* handlers only map to canonical input values before delegation

---

## Validation Expectations

Required validation after applying this subphase:

```bash
go test ./...
```

Expected result:

* all existing tests pass
* no import cycles
* no handler contract regressions
* `writemodels` packages remain compile-safe

---

## Status

Phase: 0.12.3  
Subphase: 0.12.3.4  
Status: COMPLETED  
Code Impact: YES  
Contract Impact: NONE  
Handler Alignment: DONE  

---

## Next Step

```text
0.12.3.5 — Validation & Compatibility
```
