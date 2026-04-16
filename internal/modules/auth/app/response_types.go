package app

import rootauth "github.com/e-scavo/scavo-exchange-backend/internal/modules/auth"

type LoginResponse = rootauth.LoginResponse
type MeResponse = rootauth.MeResponse
type SessionResponse = rootauth.SessionResponse

type BootstrapWalletsView = rootauth.BootstrapWalletsView
type BootstrapResponse = rootauth.BootstrapResponse

type WalletReadModel = rootauth.WalletReadModel
type WalletsQuery = rootauth.WalletsQuery
type WalletsResponse = rootauth.WalletsResponse

type ProfileView = rootauth.ProfileView
type ProfileWalletView = rootauth.ProfileWalletView

type WalletLinkChallengeResponse = rootauth.WalletLinkChallengeResponse
type WalletLinkVerifyResponse = rootauth.WalletLinkVerifyResponse
type WalletAccountMergeChallengeResponse = rootauth.WalletAccountMergeChallengeResponse
type WalletAccountMergeVerifyResponse = rootauth.WalletAccountMergeVerifyResponse
type WalletPrimarySetResponse = rootauth.WalletPrimarySetResponse
type WalletDetachCheckResponse = rootauth.WalletDetachCheckResponse
type WalletDetachExecuteResponse = rootauth.WalletDetachExecuteResponse
