# Phase 0.12.2.3 — Mapping Introduction (Domain → Read)

## Objective

Introduce explicit mapping functions from canonical domain/application data into the read models created in Phase 0.12.2.2.

This subphase keeps the change additive and non-invasive:

- no handler rewiring
- no HTTP contract changes
- no route behavior changes
- no domain model replacement
- no write model changes

---

## Source ZIP

Implemented from:

`scavo-exchange-backend-0.12.2.2.zip`

---

## Scope

### Included

- auth read model mapper functions
- user read model mapper functions
- user settings read model mapper functions
- nil-safe mapping behavior
- collection mapping where needed
- documentation of the mapping boundary

### Excluded

- response wiring
- handler changes
- public payload changes
- model deletion
- domain-to-write mapping
- CQRS behavior

---

## Files Added

### Auth

`internal/modules/auth/readmodels/mappers.go`

Introduced:

- `NewAuthLoginReadModel`
- `FromWalletIdentity`
- `FromWalletIdentities`
- `FromWalletChallenge`

### User

`internal/modules/user/readmodels/mappers.go`

Introduced:

- `FromUser`

### User Settings

`internal/modules/usersettings/readmodels/mappers.go`

Introduced:

- `FromUserSettings`

---

## Mapping Boundary

The mapping direction introduced in this subphase is:

```text
Domain/Application data → Read Model
```

This keeps read models as output-only projections and prevents transport responses from depending directly on domain internals over time.

---

## Auth Mapping Details

### Login Mapping

`NewAuthLoginReadModel` creates an `AuthLoginReadModel` from resolved authentication values.

It intentionally accepts resolved scalar values instead of importing application service structs, avoiding package coupling from readmodels back into the application layer.

### Wallet Identity Mapping

`FromWalletIdentity` maps:

`auth/domain.WalletIdentity → AuthWalletReadModel`

It preserves current output-oriented wallet fields and introduces a deterministic read status:

- `active`
- `detached`
- `unknown`

### Wallet Identity Collection Mapping

`FromWalletIdentities` maps:

`[]*auth/domain.WalletIdentity → []*AuthWalletReadModel`

Nil input returns an empty slice to keep response behavior stable.

### Wallet Challenge Mapping

`FromWalletChallenge` maps:

`auth/domain.WalletChallenge → AuthWalletChallengeReadModel`

---

## User Mapping Details

`FromUser` maps:

`user/domain.User → UserReadModel`

This preserves the existing user response shape while creating a dedicated output-only projection.

---

## User Settings Mapping Details

`FromUserSettings` maps:

`usersettings/domain.UserSettings → UserSettingsReadModel`

Behavior preserved from the existing view conversion pattern:

- nil settings produce version `1` with empty preferences
- nil preferences become an empty map
- timestamps are normalized to UTC pointers when present

---

## Compatibility Notes

This subphase does not wire the new mappers into handlers or response structs yet.

Therefore:

- public responses remain unchanged
- application behavior remains unchanged
- existing tests should continue to pass
- package imports remain acyclic

---

## Validation

`gofmt` was applied to all added Go files.

`go test ./...` could not be completed in the execution environment because the local Go toolchain attempted to download Go `1.25.0`, but external network/DNS access was unavailable.

The expected local validation command remains:

```bash
go test ./...
```

---

## Status

Phase: 0.12.2

Subphase: 0.12.2.3

Status: Completed

Code impact: Additive only

Handlers modified: No

HTTP contracts modified: No

Next: 0.12.2.4 — Response Alignment
