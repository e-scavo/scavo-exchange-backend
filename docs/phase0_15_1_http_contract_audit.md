# Phase 0.15.1 — HTTP Contract Audit

---

## Status

**Phase:** 0.15 — Contract Hardening & Freeze  
**Subphase:** 0.15.1 — HTTP Contract Audit  
**Status:** Completed  
**Type:** Documentation and contract inventory audit  
**Code changes in 0.15.1:** No

---

## Source Baseline

The audit was performed against the repository state after 0.15.0.

Audited code surfaces:

- `internal/core/httpx/router.go`
- `internal/core/httpx/auth.go`
- `internal/core/httpx/error.go`
- `internal/core/httpx/middleware.go`
- `internal/core/httpx/authorization_enforcement.go`
- `internal/core/status/service.go`
- `internal/modules/auth/http_login.go`
- `internal/modules/auth/http_bootstrap.go`
- `internal/modules/auth/http_wallet.go`
- `internal/modules/auth/http_wallet_list.go`
- `internal/modules/auth/http_contracts.go`
- `internal/modules/auth/ws_handlers.go`
- `internal/modules/system/ws_handlers.go`

This subphase does not change handlers, routes, status codes, response schemas or middleware behavior.

---

## Context Inherited

0.15.0 opened the Contract Hardening & Freeze phase and locked the rule that all later work must be derived from the real repository state.

Phase 0.15 inherits the Stage 0 path stabilized by previous phases:

- Phase 0.12 clarified read/write model separation.
- Phase 0.13 consolidated providers as the runtime composition boundary.
- Phase 0.14 exposed observability and diagnostics without changing business behavior.

With that baseline in place, 0.15.1 audits the HTTP surface before any error alignment, provider validation, response normalization or freeze enforcement work is attempted.

---

## Real Problem

The backend had a working HTTP route surface, but the contract inventory was still implicit.

The route registration in `router.go` defined public endpoints, legacy auth routes and `/api/v1` auth routes. The handlers returned stable responses, but there was no dedicated 0.15 contract audit table describing:

- method
- path
- route family
- authentication requirement
- success status
- response payload family
- public error envelope behavior

Without that inventory, later hardening subphases could accidentally normalize the wrong surface, miss a legacy route, or treat a behavior-compatible alias as a separate contract.

---

## Decision Taken

0.15.1 records the real HTTP contract inventory without changing runtime behavior.

The audit classifies the current route surface into four families:

1. Foundation status and diagnostics routes.
2. WebSocket upgrade route.
3. Public authentication routes.
4. Authenticated account, settings and wallet-management routes.

The auth routes are intentionally registered twice:

- legacy path without prefix
- canonical versioned path under `/api/v1`

Both surfaces currently point to the same handler set and are therefore audited as paired contracts.

---

## HTTP Router Inventory

The router is built in `internal/core/httpx/router.go`.

Global middleware and router behavior:

- CORS allows `GET`, `POST`, `PUT`, `PATCH`, `DELETE` and `OPTIONS`.
- `X-Request-Id` is accepted and exposed.
- `RequestID()` is applied globally.
- `Recoverer()` is applied globally.
- `AccessLog()` and `Timeout(30s)` are applied inside the main HTTP group.
- `/ws` is registered outside the AccessLog/Timeout group as a WebSocket route.

---

## Foundation HTTP Contracts

| Method | Path | Handler / Owner | Auth | Success Status | Success Payload Family | Notes |
|---|---|---|---|---:|---|---|
| GET | `/health` | inline router handler / `status.Service.Health()` | No | 200 | health object | Fallback payload is `{ok,status}` when status service is nil. |
| GET | `/readiness` | inline router handler / `status.Service.Readiness()` | No | 200 or 503 | readiness object | 503 is returned only when a required dependency is not ready. |
| GET | `/version` | inline router handler | No | 200 | version object | Includes `version`, `commit`, `env`. |
| GET | `/diagnostics` | inline router handler / `status.Service.Diagnostics()` | No | 200 | diagnostics object | Reports observability flags introduced by 0.14. |
| GET | `/ws` | `ws.NewHandler` | Token-based WebSocket flow | upgrade-dependent | WebSocket protocol | HTTP JSON response contract is not the primary surface. |

---

## Public Auth HTTP Contracts

The following public auth routes are registered in both legacy and canonical forms.

| Method | Legacy Path | Canonical Path | Handler | Auth | Success Status | Success Payload Family |
|---|---|---|---|---|---:|---|
| POST | `/auth/login` | `/api/v1/auth/login` | `Login` | No | 200 | login token response |
| POST | `/auth/wallet/challenge` | `/api/v1/auth/wallet/challenge` | `WalletChallenge` | No | 200 | wallet challenge response |
| POST | `/auth/wallet/verify` | `/api/v1/auth/wallet/verify` | `WalletVerify` | No | 200 | wallet verification token response |

Observed request behavior:

- JSON request bodies are decoded with `DisallowUnknownFields()`.
- Duplicate JSON objects or malformed JSON are rejected.
- Public auth body limits are 4 KiB for login and wallet challenge, and 8 KiB for wallet verify.
- Bad request decoding uses the standardized error envelope through `BadRequestInvalidPayload()`.

---

## Authenticated Auth / Account Contracts

The following authenticated routes are registered in both legacy and canonical forms.

| Method | Legacy Path | Canonical Path | Handler | Middleware | Success Status | Success Payload Family |
|---|---|---|---|---|---:|---|
| GET | `/auth/bootstrap` | `/api/v1/auth/bootstrap` | `Bootstrap` | auth + authorization hydration | 200 | bootstrap response |
| GET | `/auth/me` | `/api/v1/auth/me` | `Me` | auth + authorization hydration + user read permission | 200 | user/profile response |
| PATCH | `/auth/me` | `/api/v1/auth/me` | `UpdateMe` | auth + authorization hydration | 200 | user/profile response |
| GET | `/auth/me/settings` | `/api/v1/auth/me/settings` | `MeSettings` | auth + authorization hydration + settings read permission | 200 | settings response |
| PATCH | `/auth/me/settings` | `/api/v1/auth/me/settings` | `UpdateMeSettings` | auth + authorization hydration + settings update permission | 200 | settings response |
| GET | `/auth/session` | `/api/v1/auth/session` | `Session` | auth + authorization hydration | 200 | session response |

Observed request behavior:

- `RequireAuth(tokens, false)` requires a bearer token.
- Authenticated handlers also call `requireClaims()` and return `AUTH_UNAUTHORIZED` if claims are missing or invalid.
- Profile and settings updates decode JSON bodies with unknown-field rejection.
- Settings update rejects empty or invalid `preferences` payloads.

---

## Authenticated Wallet Management Contracts

The following authenticated wallet-management routes are registered in both legacy and canonical forms.

| Method | Legacy Path | Canonical Path | Handler | Middleware | Success Status | Success Payload Family |
|---|---|---|---|---|---:|---|
| GET | `/auth/wallets` | `/api/v1/auth/wallets` | `Wallets` | auth + authorization hydration | 200 | wallet list response |
| POST | `/auth/wallets/link/challenge` | `/api/v1/auth/wallets/link/challenge` | `WalletLinkChallenge` | auth + authorization hydration | 200 | wallet link challenge response |
| POST | `/auth/wallets/link/verify` | `/api/v1/auth/wallets/link/verify` | `WalletLinkVerify` | auth + authorization hydration | 200 | wallet link verification response |
| POST | `/auth/account/merge/wallet/challenge` | `/api/v1/auth/account/merge/wallet/challenge` | `WalletAccountMergeChallenge` | auth + authorization hydration | 200 | wallet account merge challenge response |
| POST | `/auth/account/merge/wallet/verify` | `/api/v1/auth/account/merge/wallet/verify` | `WalletAccountMergeVerify` | auth + authorization hydration | 200 | wallet account merge verification response |
| POST | `/auth/wallets/detach/check` | `/api/v1/auth/wallets/detach/check` | `WalletDetachCheck` | auth + authorization hydration | 200 | wallet detach check response |
| POST | `/auth/wallets/detach` | `/api/v1/auth/wallets/detach` | `WalletDetach` | auth + authorization hydration | 200 | wallet detach execution response |
| POST | `/auth/wallets/primary` | `/api/v1/auth/wallets/primary` | `WalletSetPrimary` | auth + authorization hydration | 200 | wallet primary set response |

Observed query contract for `GET /auth/wallets` and `GET /api/v1/auth/wallets`:

- `status`: `active` or `detached`
- `primary`: `true` or `false`
- `sort`: currently `linked_at`
- `order`: `asc` or `desc`
- `limit`: positive integer
- `offset`: non-negative integer
- `order` without `sort` is rejected
- `sort` without `order` defaults order to `asc`

---

## WebSocket Action Inventory

The HTTP route surface includes `GET /ws`, and module registration exposes the following WebSocket actions:

| Action | Owner | Auth Requirement | Response Family |
|---|---|---|---|
| `system.ping` | `internal/modules/system` | No action-level auth wrapper in registration | ping response |
| `auth.whoami` | `internal/modules/auth` | `ws.RequireAuth` | auth identity response |
| `auth.session` | `internal/modules/auth` | `ws.RequireAuth` | auth session response |

0.15.1 records these actions because `/ws` is part of the HTTP router. Detailed WebSocket protocol freeze is not expanded in this subphase; later contract work may reference this inventory if response schema normalization touches WebSocket-facing payloads.

---

## Error Contract Observations

The current HTTP surface consistently routes handler and middleware failures through the standardized error envelope:

```json
{
  "error": {
    "code": "...",
    "message": "...",
    "details": {}
  }
}
```

Observed error writers:

- `httpx.WriteAppError()` for core HTTP middleware.
- `httpx.WriteError()` and `httpx.WriteErrorMessage()` for generic core HTTP error helpers.
- `auth.writeAppErrorJSON()` for auth handlers.
- `auth.writeErrorJSON()` for legacy auth-key compatibility.

Observed status families:

- `400` for malformed or invalid request payload/query contracts.
- `401` for missing, invalid or expired authentication credentials.
- `403` for authorization policy denial.
- `404` for not-found user/wallet/challenge cases.
- `409` for wallet conflict or ownership-rule failures.
- `500` for internal/auth/provider/service failures.
- `503` for timeout or readiness dependency failure.

No non-envelope auth error response was identified in the audited HTTP handlers.

---

## Response Contract Observations

The success response surface is already separated into named handler-facing contract types or application response aliases:

- `LoginResponse`
- `MeResponse`
- `SessionResponse`
- `MeSettingsResponse`
- `WalletChallengeResponse`
- `WalletVerifyResponse`
- `WalletLinkChallengeResponse`
- `WalletLinkVerifyResponse`
- `WalletAccountMergeChallengeResponse`
- `WalletAccountMergeVerifyResponse`
- `WalletDetachCheckResponse`
- `WalletDetachExecuteResponse`
- `WalletPrimarySetResponse`
- `WalletsResponse`
- `BootstrapResponse`

Foundation routes use explicit map payloads or `status.Service` payload builders rather than auth module contract structs. That is acceptable for 0.15.1, but it is a normalization candidate for 0.15.4 if the phase decides that foundation surfaces should also have named response contract types.

---

## Detected Contract Risks

0.15.1 did not detect a blocking route inconsistency that requires code changes before 0.15.2.

The audit did identify non-blocking risks that must be preserved for later subphases:

1. Legacy and `/api/v1` auth routes are both active. This is intentional and tested, but the freeze policy must define whether both are frozen or whether legacy paths are compatibility aliases.
2. Foundation status routes return map-based payloads rather than named response structs. This is stable behavior today, but 0.15.4 should decide whether to leave this documented or introduce explicit response schema types.
3. The `/ws` route is an HTTP route but its actual contract is action-based WebSocket messaging. It should remain separated from REST response normalization unless explicitly targeted.
4. `GET /readiness` legitimately has two success/failure status outcomes: `200` when ready or optional dependencies are down, and `503` when a required dependency is not ready.
5. CORS allows additional methods globally even though the current route surface uses only `GET`, `POST` and `PATCH`. This is CORS policy, not an endpoint inventory.

---

## Contract Audit Result

0.15.1 establishes the HTTP route baseline for Phase 0.15.

Audited registered HTTP paths:

- Foundation and diagnostics paths: 5
- Legacy auth/account/wallet paths: 17
- Canonical `/api/v1` auth/account/wallet paths: 17

Total registered HTTP route entries audited: 39

Unique behavior contracts:

- Foundation and diagnostics contracts: 5
- Auth/account/wallet behavior contracts: 17

Total unique behavior contracts: 22

No new endpoint, route, response shape, status code, middleware or business behavior was introduced.

---

## Validation

A validation attempt was made with:

```bash
go test ./...
```

The command could not complete in the execution environment because the Go toolchain attempted to download `go1.25.0` from `proxy.golang.org`, but DNS/network access was unavailable.

This is recorded as an environment limitation, not as a code failure, because 0.15.1 does not modify Go source code.

---

## Concrete Change

0.15.1 introduces this HTTP contract audit document and reconciles trunk documentation so the current phase state advances from 0.15.0 to 0.15.1 completed.

No source code was changed.

---

## Observable Impact

After 0.15.1:

- the real HTTP route surface is explicitly inventoried
- legacy and `/api/v1` route pairing is documented
- success status families are documented
- response payload families are identified
- error envelope behavior is confirmed as the current handler contract
- later subphases have a concrete route baseline to validate against

---

## Handoff to 0.15.2

The next subphase is:

**0.15.2 — Error Contract Alignment**

It must use this HTTP route inventory as the public surface baseline and validate error contracts without inventing new endpoints or changing business behavior.

---

## End of Document
