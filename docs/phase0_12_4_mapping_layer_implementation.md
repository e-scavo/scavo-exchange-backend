# Phase 0.12.4.2 — Mapping Layer Implementation

## Status

Completed.

## Objective

Introduce the explicit module-level mapping layer defined in Phase 0.12.4.0 and designed in Phase 0.12.4.1.

This subphase implements the physical package structure for centralized mapping without changing HTTP contracts, endpoint behavior, or handler wiring.

## Scope

### Included

- creation of module-level `mappers` packages
- explicit `Domain → Read` mapping functions
- explicit `Write → Domain` mapping functions
- preservation of existing read/write model packages
- preservation of current handler behavior
- additive-only code changes

### Excluded

- handler refactor
- application service refactor
- removal of existing mapper functions
- HTTP contract changes
- response shape changes
- business logic changes

## Implemented Packages

### Auth

```text
internal/modules/auth/mappers/
```

Files:

```text
internal/modules/auth/mappers/read.go
internal/modules/auth/mappers/write.go
```

Responsibilities:

- map authentication domain state into auth read models
- map wallet domain state into wallet read models
- map wallet challenge domain state into challenge read models
- map auth write models into canonical domain inputs

### User

```text
internal/modules/user/mappers/
```

Files:

```text
internal/modules/user/mappers/read.go
internal/modules/user/mappers/write.go
```

Responsibilities:

- map user domain state into user read models
- map user write models into canonical user domain inputs

### User Settings

```text
internal/modules/usersettings/mappers/
```

Files:

```text
internal/modules/usersettings/mappers/read.go
internal/modules/usersettings/mappers/write.go
```

Responsibilities:

- map user settings domain state into settings read models
- map user settings write models into canonical settings domain inputs

## Implemented Mapping Directions

### Domain → Read

```text
auth/domain.WalletIdentity → auth/readmodels.AuthWalletReadModel
auth/domain.WalletChallenge → auth/readmodels.AuthWalletChallengeReadModel
user/domain.User → user/readmodels.UserReadModel
usersettings/domain.UserSettings → usersettings/readmodels.UserSettingsReadModel
```

### Write → Domain

```text
auth/writemodels.AuthLoginWriteModel → auth/domain.LoginInput
auth/writemodels.AuthUpdateProfileWriteModel → auth/domain.ProfileUpdateInput
auth/writemodels.AuthUpdateSettingsWriteModel → auth/domain.SettingsUpdateInput
auth/writemodels.AuthWalletChallengeWriteModel → auth/domain.WalletChallengeInput
auth/writemodels.AuthWalletVerifyWriteModel → auth/domain.WalletVerifyInput
auth/writemodels.AuthWalletLinkChallengeWriteModel → auth/domain.WalletChallengeInput
auth/writemodels.AuthWalletLinkVerifyWriteModel → auth/domain.WalletVerifyInput
auth/writemodels.AuthWalletAccountMergeChallengeWriteModel → auth/domain.WalletChallengeInput
auth/writemodels.AuthWalletAccountMergeVerifyWriteModel → auth/domain.WalletVerifyInput
auth/writemodels.AuthWalletDetachCheckWriteModel → auth/domain.WalletDetachInput
auth/writemodels.AuthWalletDetachExecuteWriteModel → auth/domain.WalletDetachInput
auth/writemodels.AuthWalletPrimarySetWriteModel → auth/domain.WalletPrimarySetInput
user/writemodels.UserUpdateWriteModel → user/domain.UserUpdateInput
usersettings/writemodels.UserSettingsUpdateWriteModel → usersettings/domain.UserSettingsUpdateInput
```

## Compatibility Strategy

This subphase intentionally keeps existing mapper functions under:

```text
internal/modules/<module>/readmodels/
internal/modules/<module>/writemodels/
```

Those functions remain temporarily available to avoid a broad behavior-changing refactor in the implementation step.

The new canonical mapping layer is introduced under:

```text
internal/modules/<module>/mappers/
```

Future subphases will migrate call sites progressively.

## Import Boundary

The new mapping packages are allowed to import:

- domain models
- read models
- write models
- standard library helpers required for decoding

They must not import:

- HTTP handlers
- application services
- repositories
- server/bootstrap packages

## HTTP Contract Preservation

No HTTP structs, JSON tags, endpoint paths, status codes, or response envelopes were changed in this subphase.

## Handler Preservation

Handlers were not modified in this subphase.

Current handler call sites continue using the compatibility functions introduced in earlier Phase 0.12.2 and Phase 0.12.3 work.

## Application Layer Preservation

Application orchestration remains unchanged in this subphase.

The new mapping layer is available for consolidation but is not forced into runtime call paths yet.

## Rationale

This avoids combining three concerns in a single subphase:

1. creating the mapping layer
2. migrating call sites
3. removing old mapping helpers

Keeping this subphase additive reduces risk and preserves compile-time stability.

## Validation Expectations

After applying this subphase, the project must continue to pass:

```bash
go test ./...
```

The new packages should appear as packages with no test files until dedicated mapper tests are introduced in a later phase.

## Files Added

```text
internal/modules/auth/mappers/read.go
internal/modules/auth/mappers/write.go
internal/modules/user/mappers/read.go
internal/modules/user/mappers/write.go
internal/modules/usersettings/mappers/read.go
internal/modules/usersettings/mappers/write.go
docs/phase0_12_4_mapping_layer_implementation.md
```

## Code Impact

Additive only.

## Contract Impact

None.

## Runtime Impact

None expected, because runtime call sites are not migrated in this subphase.

## Next Step

```text
0.12.4.3 — Mapping Consolidation
```

The next subphase should progressively migrate existing call sites from scattered mapper helpers into the centralized mapping layer while preserving compatibility.
