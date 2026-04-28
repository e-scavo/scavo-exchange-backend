# Phase 0.12.2.4 — Response Alignment

## Objective

Align application response structures with the explicit read models introduced in Phase 0.12.2.2 and mapped in Phase 0.12.2.3, without changing HTTP contracts, endpoint behavior, routing, versioning, or handler ownership.

This subphase keeps the read/write model separation incremental and compatible by ensuring that application-level response payloads are now backed by output-only read models while preserving the existing JSON shape.

---

## Scope

### Included

* align auth application response types with explicit read models
* align user payload references with `UserReadModel`
* align user settings payload references with `UserSettingsReadModel`
* align wallet payload references with `AuthWalletReadModel`
* align wallet challenge payload references with `AuthWalletChallengeReadModel`
* preserve root auth compatibility adapters
* preserve existing HTTP JSON contracts

### Excluded

* handler rewrites
* endpoint changes
* versioning changes
* write model extraction
* CQRS implementation
* business logic changes

---

## Files Updated

### Code

* `internal/modules/auth/app/response_types.go`
* `internal/modules/auth/app/application.go`
* `internal/modules/auth/app/support_helpers.go`
* `internal/modules/auth/readmodels/mappers.go`
* `internal/modules/auth/application.go`

### Documentation

* `docs/phase0_12_2_response_alignment.md`

---

## Alignment Summary

### Login Response

`LoginResponse` now aliases the explicit read model:

```text
LoginResponse → AuthLoginReadModel
```

This preserves the existing JSON response shape:

* `access_token`
* `token_type`
* `expires_in`
* `user_id`

---

### User Output

Application-layer user output now uses:

```text
UserReadModel
```

Affected application response areas:

* session view embedded user
* profile view embedded user
* `/me` response user
* bootstrap user

Root auth compatibility remains preserved through adapter mapping back to the existing public root auth response types.

---

### User Settings Output

Bootstrap settings now use:

```text
UserSettingsReadModel
```

The existing JSON shape remains compatible with the previous settings view:

* `user_id`
* `version`
* `preferences`
* `created_at`
* `updated_at`

---

### Wallet Output

Application-layer wallet output now uses:

```text
AuthWalletReadModel
```

Affected response areas:

* wallet listing
* bootstrap wallets
* wallet link verification
* account merge verification
* detach execution
* primary wallet selection

Existing fields remain preserved:

* `id`
* `address`
* `user_id`
* `linked_at`
* `detached_at`
* `is_primary`
* `status`
* `can_set_primary`
* `can_detach`
* `detach_block_reasons`

---

### Wallet Challenge Output

Application-layer wallet challenge output now uses:

```text
AuthWalletChallengeReadModel
```

Affected response areas:

* wallet link challenge
* wallet link verify
* account merge challenge
* account merge verify

Existing JSON fields remain preserved.

---

## Compatibility Boundary

The root `auth` package still exposes existing public response structures expected by current handlers and tests.

To preserve compatibility, root-level adapters map application read models back to the existing root response payload types where needed.

This avoids forcing a handler rewrite during 0.12.2.4.

---

## Contract Preservation

This subphase does not change:

* endpoint paths
* HTTP methods
* request payloads
* response JSON field names
* API versioning
* error envelopes
* authentication behavior
* authorization behavior

---

## Implementation Notes

The application layer now uses explicit output projections for read-facing payloads.

Mapping remains explicit and directional:

```text
Domain → Read Model → Response Boundary
```

No read model is used as an input command.

---

## Validation Notes

The implementation is additive/alignment-focused.

Expected validation command:

```bash
go test ./...
```

The local execution environment used for generation could not complete `go test ./...` because it attempted to download the Go 1.25.0 toolchain and external network access was unavailable.

---

## Status

Phase: 0.12.2
Subphase: 0.12.2.4
Status: COMPLETED
Code Impact: RESPONSE ALIGNMENT ONLY
HTTP Contract Impact: NONE

---

## Next Step

```text
0.12.2.5 — Validation & Compatibility
```
