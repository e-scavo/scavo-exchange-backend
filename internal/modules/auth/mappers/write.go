package mappers

import (
	"encoding/json"

	authdomain "github.com/e-scavo/scavo-exchange-backend/internal/modules/auth/domain"
	authwritemodels "github.com/e-scavo/scavo-exchange-backend/internal/modules/auth/writemodels"
)

// LoginWriteToDomainInput maps login write intent into canonical domain input.
func LoginWriteToDomainInput(model authwritemodels.AuthLoginWriteModel) authdomain.LoginInput {
	return authdomain.LoginInput{
		Email:    model.Email,
		Password: model.Password,
	}
}

// ProfileUpdateWriteToDomainInput maps authenticated profile update write intent
// into canonical domain input.
func ProfileUpdateWriteToDomainInput(model authwritemodels.AuthUpdateProfileWriteModel) authdomain.ProfileUpdateInput {
	return authdomain.ProfileUpdateInput{DisplayName: model.DisplayName}
}

// SettingsUpdateWriteToDomainInput maps authenticated settings write intent into
// canonical domain input. JSON decoding remains explicit to preserve current
// invalid-payload behavior.
func SettingsUpdateWriteToDomainInput(model authwritemodels.AuthUpdateSettingsWriteModel) (authdomain.SettingsUpdateInput, error) {
	input := authdomain.SettingsUpdateInput{}
	if len(model.Preferences) == 0 {
		return input, nil
	}
	if err := json.Unmarshal(model.Preferences, &input.Preferences); err != nil {
		return authdomain.SettingsUpdateInput{}, err
	}
	return input, nil
}

// WalletChallengeWriteToDomainInput maps wallet challenge creation intent into
// canonical domain input.
func WalletChallengeWriteToDomainInput(model authwritemodels.AuthWalletChallengeWriteModel) authdomain.WalletChallengeInput {
	return authdomain.WalletChallengeInput{
		Address: model.Address,
		Chain:   model.Chain,
	}
}

// WalletVerifyWriteToDomainInput maps wallet verification intent into canonical
// domain input.
func WalletVerifyWriteToDomainInput(model authwritemodels.AuthWalletVerifyWriteModel) authdomain.WalletVerifyInput {
	return authdomain.WalletVerifyInput{
		ChallengeID: model.ChallengeID,
		Address:     model.Address,
		Signature:   model.Signature,
	}
}

// WalletLinkChallengeWriteToDomainInput maps authenticated wallet-link challenge
// creation intent into canonical domain input.
func WalletLinkChallengeWriteToDomainInput(model authwritemodels.AuthWalletLinkChallengeWriteModel) authdomain.WalletChallengeInput {
	return authdomain.WalletChallengeInput{
		Address: model.Address,
		Chain:   model.Chain,
	}
}

// WalletLinkVerifyWriteToDomainInput maps authenticated wallet-link verification
// intent into canonical domain input.
func WalletLinkVerifyWriteToDomainInput(model authwritemodels.AuthWalletLinkVerifyWriteModel) authdomain.WalletVerifyInput {
	return authdomain.WalletVerifyInput{
		ChallengeID: model.ChallengeID,
		Address:     model.Address,
		Signature:   model.Signature,
	}
}

// WalletAccountMergeChallengeWriteToDomainInput maps authenticated account-merge
// challenge creation intent into canonical domain input.
func WalletAccountMergeChallengeWriteToDomainInput(model authwritemodels.AuthWalletAccountMergeChallengeWriteModel) authdomain.WalletChallengeInput {
	return authdomain.WalletChallengeInput{
		Address: model.Address,
		Chain:   model.Chain,
	}
}

// WalletAccountMergeVerifyWriteToDomainInput maps authenticated account-merge
// verification intent into canonical domain input.
func WalletAccountMergeVerifyWriteToDomainInput(model authwritemodels.AuthWalletAccountMergeVerifyWriteModel) authdomain.WalletVerifyInput {
	return authdomain.WalletVerifyInput{
		ChallengeID: model.ChallengeID,
		Address:     model.Address,
		Signature:   model.Signature,
	}
}

// WalletDetachCheckWriteToDomainInput maps wallet detach-check intent into
// canonical domain input.
func WalletDetachCheckWriteToDomainInput(model authwritemodels.AuthWalletDetachCheckWriteModel) authdomain.WalletDetachInput {
	return authdomain.WalletDetachInput{Address: model.Address}
}

// WalletDetachExecuteWriteToDomainInput maps wallet detach-execute intent into
// canonical domain input.
func WalletDetachExecuteWriteToDomainInput(model authwritemodels.AuthWalletDetachExecuteWriteModel) authdomain.WalletDetachInput {
	return authdomain.WalletDetachInput{Address: model.Address}
}

// WalletPrimarySetWriteToDomainInput maps primary wallet update intent into
// canonical domain input.
func WalletPrimarySetWriteToDomainInput(model authwritemodels.AuthWalletPrimarySetWriteModel) authdomain.WalletPrimarySetInput {
	return authdomain.WalletPrimarySetInput{Address: model.Address}
}
