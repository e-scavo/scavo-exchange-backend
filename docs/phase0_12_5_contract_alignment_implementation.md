# Phase 0.12.5 — Contract Alignment

## Subphase 0.12.5.3 — Contract Alignment Implementation

---

## Objective

Apply the first implementation-level contract alignment changes defined by Phase 0.12.5.1 and Phase 0.12.5.2 while preserving all external API compatibility.

This subphase aligns selected HTTP-facing response contracts with the explicit Read Model layer and the centralized Mapping Layer introduced in Phase 0.12.4.

---

## Source Baseline

Baseline ZIP: `scavo-exchange-backend-0.12.5.2.zip`

Relevant prior phases:

* Phase 0.12.2 — Read Model Extraction
* Phase 0.12.3 — Write Model Isolation
* Phase 0.12.4 — Mapping Layer Introduction
* Phase 0.12.5.1 — Contract Inventory & Classification
* Phase 0.12.5.2 — Contract Normalization Design

---

## Scope

### Included

* align remaining HTTP response fields that directly exposed domain/legacy view contracts where safe
* preserve existing HTTP response names
* preserve existing JSON field names
* preserve existing handler behavior
* keep mapping responsibility inside module-level mappers
* document compatibility impact

### Excluded

* endpoint renames
* route changes
* API version changes
* broad provider interface refactors
* repository contract changes
* WebSocket protocol changes
* CQRS or event sourcing

---

## Implementation Summary

The implementation keeps existing public response contract names while normalizing selected response payload fields to explicit read models.

The following compatibility surfaces remain stable:

* `MeSettingsResponse`
* `WalletChallengeResponse`
* `WalletVerifyResponse`

The underlying response payload composition now uses read models where previously a domain model or legacy view could be exposed directly.

---

## Files Changed

### `internal/modules/auth/http_login.go`

Changes:

* `MeSettingsResponse.Settings` now uses `usersettingsreadmodels.UserSettingsReadModel`.
* settings responses are mapped through `usersettingsmappers.UserSettingsToReadModel`.
* direct handler usage of `usersettingsmod.ToView` was removed from auth HTTP settings responses.

Compatibility:

* response wrapper name remains `MeSettingsResponse`
* JSON field remains `settings`
* nested JSON fields remain compatible:
  * `user_id`
  * `version`
  * `preferences`
  * `created_at`
  * `updated_at`

---

### `internal/modules/auth/http_wallet.go`

Changes:

* `WalletChallengeResponse.Challenge` now uses `authreadmodels.AuthWalletChallengeReadModel`.
* `WalletVerifyResponse.Challenge` now uses `authreadmodels.AuthWalletChallengeReadModel`.
* `WalletVerifyResponse.User` now uses `userreadmodels.UserReadModel`.
* challenge response mapping now flows through `authmappers.WalletChallengeToReadModel`.
* wallet verification user mapping now flows through `usermappers.UserToReadModel`.

Compatibility:

* response wrapper names remain unchanged:
  * `WalletChallengeResponse`
  * `WalletVerifyResponse`
* JSON field names remain unchanged:
  * `challenge`
  * `user`
* nested challenge JSON fields remain compatible with the previous domain-backed shape
* nested user JSON fields remain compatible with the previous user-backed shape

---

## Contract Alignment Details

### Settings Response Alignment

Previous flow:

```text
usersettings domain model → legacy View → MeSettingsResponse
```

New flow:

```text
usersettings domain model → Mapping Layer → UserSettingsReadModel → MeSettingsResponse
```

This reduces direct legacy view usage in auth HTTP handlers and aligns settings output with the explicit read model lifecycle.

---

### Wallet Challenge Response Alignment

Previous flow:

```text
WalletChallenge domain model → WalletChallengeResponse
```

New flow:

```text
WalletChallenge domain model → Mapping Layer → AuthWalletChallengeReadModel → WalletChallengeResponse
```

This keeps the HTTP contract stable while preventing direct domain challenge exposure from the handler response struct.

---

### Wallet Verify User Response Alignment

Previous flow:

```text
user domain model → WalletVerifyResponse
```

New flow:

```text
user domain model → Mapping Layer → UserReadModel → WalletVerifyResponse
```

This aligns wallet authentication responses with the read model boundary already used by `/me`, bootstrap and session/profile surfaces.

---

## Compatibility Rules Preserved

The following compatibility rules remain enforced:

* existing HTTP response names are preserved
* existing JSON field names are preserved
* request decoding behavior is unchanged
* error behavior is unchanged
* auth/session behavior is unchanged
* wallet challenge creation behavior is unchanged
* wallet verification behavior is unchanged
* no route changes were introduced
* no versioning changes were introduced

---

## Relationship to Phase 0.12.5.2 Design

This implementation follows the normalization target defined in Phase 0.12.5.2:

```text
Domain State / Application State
  ↓
Mapping Layer
  ↓
Read Model / Read View
  ↓
HTTP Response Name
```

The HTTP response names remain as compatibility contracts, while payload fields now point to explicit read models where safe.

---

## Remaining Compatibility Surfaces

Some compatibility surfaces intentionally remain unchanged:

* provider interfaces still return internal module/domain data
* application service internals still use domain models where appropriate
* wallet operation response aliases remain stable
* `WalletsQuery` remains an app-level query input contract
* core WebSocket and error contracts remain out of scope

These are deferred to later contract alignment or closure validation if needed.

---

## Validation Expectations

The implementation must be validated with:

```bash
go test ./...
```

Expected result:

* all packages compile
* existing tests pass
* auth HTTP tests remain valid
* no import cycles are introduced
* no JSON compatibility changes are introduced

---

## Risk Assessment

### Low Risk

The changes are type-level response alignment changes with stable JSON tags.

### Main Risk

A test or consumer could depend on Go struct field types internally rather than JSON shape.

### Mitigation

The public JSON contract remains stable, and the changed fields use read models with matching JSON-compatible fields.

---

## Status

Phase: 0.12.5  
Subphase: 0.12.5.3  
Status: COMPLETED  
Code Impact: LIMITED / COMPATIBILITY-PRESERVING  
Contract Impact: INTERNAL ALIGNMENT ONLY  
External API Impact: NONE  

---

## Next Step

```text
0.12.5.4 — Handler Contract Adjustment
```
