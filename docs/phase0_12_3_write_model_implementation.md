# Phase 0.12.3 — Write Model Isolation

## Subphase 0.12.3.2 — Write Model Implementation

---

## Objective

Introduce explicit Write Model packages and structs for the current input surfaces without changing runtime behavior.

This subphase implements the write model definitions designed in `0.12.3.1` as additive code only.

---

## Scope

### Included

* creation of write model packages
* creation of input-only structs
* preservation of existing JSON field names
* preservation of existing handlers and request structs
* preparation for explicit Write → Domain mapping

### Excluded

* handler replacement
* request decoding changes
* validation changes
* business logic changes
* endpoint contract changes
* removal of legacy request structs

---

## Implementation Summary

The following packages were introduced:

```text
internal/modules/auth/writemodels
internal/modules/user/writemodels
internal/modules/usersettings/writemodels
```

These packages are intentionally not wired into handlers yet.

They provide the structural foundation required by later subphases:

* `0.12.3.3 — Mapping Introduction (Write → Domain)`
* `0.12.3.4 — Handler Alignment`

---

## Added Files

### Auth Write Models

```text
internal/modules/auth/writemodels/write_models.go
```

Introduced input-only models for:

* login
* authenticated profile update
* authenticated settings update
* public wallet challenge
* public wallet verification
* wallet link challenge
* wallet link verification
* account merge challenge
* account merge verification
* wallet detach check
* wallet detach execution
* primary wallet update

---

### User Write Models

```text
internal/modules/user/writemodels/write_models.go
```

Introduced input-only model for:

* user profile update

---

### User Settings Write Models

```text
internal/modules/usersettings/writemodels/write_models.go
```

Introduced input-only model for:

* user settings update

---

## Write Models Introduced

### Auth Module

```text
AuthLoginWriteModel
AuthUpdateProfileWriteModel
AuthUpdateSettingsWriteModel
AuthWalletChallengeWriteModel
AuthWalletVerifyWriteModel
AuthWalletLinkChallengeWriteModel
AuthWalletLinkVerifyWriteModel
AuthWalletAccountMergeChallengeWriteModel
AuthWalletAccountMergeVerifyWriteModel
AuthWalletDetachCheckWriteModel
AuthWalletDetachExecuteWriteModel
AuthWalletPrimarySetWriteModel
```

---

### User Module

```text
UserUpdateWriteModel
```

---

### User Settings Module

```text
UserSettingsUpdateWriteModel
```

---

## Compatibility Rules Applied

### Existing HTTP Contracts Preserved

All JSON tags match the existing request payloads.

No endpoint payload shape was changed.

---

### Existing Handlers Preserved

No handler was modified in this subphase.

Existing request structs remain active.

---

### Existing Runtime Behavior Preserved

Because the new write models are additive and not wired yet, runtime behavior remains unchanged.

---

## Design Notes

### Raw JSON Preservation

Settings update write models use `json.RawMessage` for preferences.

This preserves compatibility with the current transport layer and defers semantic conversion to the mapping phase.

---

### Intent-Oriented Naming

Write model names follow operation intent rather than persistence shape.

This avoids leaking domain or database structure into input contracts.

---

### No Domain Embedding

Write models do not embed domain structs.

They are standalone input-intent models.

---

## Architectural Effect

Before this subphase, write intent was represented directly by handler request structs.

After this subphase, explicit write model packages exist and can be used as the migration target for handler alignment.

The target flow remains:

```text
HTTP Request → Write Model → Domain → Read Model → HTTP Response
```

Only the Write Model definition step is implemented here.

---

## Validation

### gofmt

`gofmt` was applied to all new Go files.

### go test

`go test ./...` could not be completed in the generation environment because the project requires Go `1.25.0` and the local environment could not download the toolchain.

The user environment must validate with:

```bash
go test ./...
```

---

## Code Impact

Code impact is additive only.

No existing Go file was modified.

---

## Status

Phase: `0.12.3`  
Subphase: `0.12.3.2`  
Status: `COMPLETED`  
Code Impact: `ADDITIVE ONLY`  
Runtime Behavior Change: `NO`  
HTTP Contract Change: `NO`

---

## Next Step

```text
0.12.3.3 — Mapping Introduction (Write → Domain)
```
