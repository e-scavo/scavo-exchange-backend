# Phase 0.12.3 — Write Model Isolation

## Subphase 0.12.3.3 — Mapping Introduction (Write → Domain)

---

## Objective

Introduce explicit Write → Domain mapping boundaries for the write models created in Phase 0.12.3.2.

This subphase establishes the controlled transformation path from transport/input intent into canonical domain input structures without changing runtime handlers yet.

---

## Scope

### Included

* domain input structures for write operations
* explicit mapper functions from write models into domain input structures
* additive package-level mapping boundaries
* compatibility with existing handlers and HTTP contracts

### Excluded

* handler replacement
* request decoding changes
* HTTP contract changes
* validation behavior changes
* removal of legacy request structs

---

## Implementation Summary

The implementation adds explicit domain input structures and mapper functions in three module areas:

* `auth`
* `user`
* `usersettings`

The runtime flow remains unchanged for now.

Existing handlers continue to decode the legacy request structs until Phase 0.12.3.4 performs handler alignment.

---

## Files Added

### Auth Domain Inputs

```text
internal/modules/auth/domain/write_inputs.go
```

Introduces canonical auth-domain input structures:

* `LoginInput`
* `ProfileUpdateInput`
* `SettingsUpdateInput`
* `WalletChallengeInput`
* `WalletVerifyInput`
* `WalletDetachInput`
* `WalletPrimarySetInput`

These structures represent domain-level intent after transport decoding and before application/service execution.

---

### User Domain Inputs

```text
internal/modules/user/domain/write_inputs.go
```

Introduces:

* `UserUpdateInput`

This structure represents canonical user update intent.

---

### User Settings Domain Inputs

```text
internal/modules/usersettings/domain/write_inputs.go
```

Introduces:

* `UserSettingsUpdateInput`

This structure represents canonical user settings update intent.

---

### Auth Write Mappers

```text
internal/modules/auth/writemodels/mappers.go
```

Introduces explicit mappings from auth write models to auth domain inputs.

Mappings include:

* `AuthLoginWriteModel` → `authdomain.LoginInput`
* `AuthUpdateProfileWriteModel` → `authdomain.ProfileUpdateInput`
* `AuthUpdateSettingsWriteModel` → `authdomain.SettingsUpdateInput`
* `AuthWalletChallengeWriteModel` → `authdomain.WalletChallengeInput`
* `AuthWalletVerifyWriteModel` → `authdomain.WalletVerifyInput`
* `AuthWalletLinkChallengeWriteModel` → `authdomain.WalletChallengeInput`
* `AuthWalletLinkVerifyWriteModel` → `authdomain.WalletVerifyInput`
* `AuthWalletAccountMergeChallengeWriteModel` → `authdomain.WalletChallengeInput`
* `AuthWalletAccountMergeVerifyWriteModel` → `authdomain.WalletVerifyInput`
* `AuthWalletDetachCheckWriteModel` → `authdomain.WalletDetachInput`
* `AuthWalletDetachExecuteWriteModel` → `authdomain.WalletDetachInput`
* `AuthWalletPrimarySetWriteModel` → `authdomain.WalletPrimarySetInput`

---

### User Write Mappers

```text
internal/modules/user/writemodels/mappers.go
```

Introduces:

* `UserUpdateWriteModel` → `userdomain.UserUpdateInput`

---

### User Settings Write Mappers

```text
internal/modules/usersettings/writemodels/mappers.go
```

Introduces:

* `UserSettingsUpdateWriteModel` → `usersettingsdomain.UserSettingsUpdateInput`

---

## Mapping Direction

The explicit direction introduced by this subphase is:

```text
Write Model → Domain Input
```

This complements the read-side mapping already introduced in Phase 0.12.2:

```text
Domain → Read Model
```

Together, the model flow becomes:

```text
Write Model → Domain → Read Model
```

---

## Compatibility Notes

This subphase is intentionally additive.

It does not alter:

* request structs currently decoded by handlers
* endpoint behavior
* response bodies
* service methods
* validation errors
* API versioning

The mapper functions are available for use by the next phase:

```text
0.12.3.4 — Handler Alignment
```

---

## JSON Preference Handling

Settings-related write models preserve `json.RawMessage` at transport/input level.

Mapping functions explicitly decode the raw JSON into:

```text
map[string]any
```

This preserves current invalid-payload behavior while making the transformation boundary explicit.

---

## Architectural Impact

The write side now has a clear transformation boundary:

* transport-level input models stay in `writemodels`
* domain-level intent structures stay in `domain`
* conversion is explicit and testable

No domain package imports write models.

The dependency direction remains:

```text
writemodels → domain
```

This keeps domain packages independent from transport/input packages.

---

## Validation Notes

Expected validation after integration:

```bash
go test ./...
```

The local generation environment could apply `gofmt`, but could not complete `go test ./...` because the container does not have the required Go toolchain version available.

---

## Status

Phase: 0.12.3  
Subphase: 0.12.3.3  
Status: IMPLEMENTED  
Code Impact: ADDITIVE  
Handler Impact: NONE  
HTTP Contract Impact: NONE

---

## Next Step

```text
0.12.3.4 — Handler Alignment
```
