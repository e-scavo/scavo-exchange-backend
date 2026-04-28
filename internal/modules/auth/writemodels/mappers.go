package writemodels

import (
	"encoding/json"

	authdomain "github.com/e-scavo/scavo-exchange-backend/internal/modules/auth/domain"
)

// ToDomainInput maps login write intent into canonical domain input.
func (m AuthLoginWriteModel) ToDomainInput() authdomain.LoginInput {
	return authdomain.LoginInput{
		Email:    m.Email,
		Password: m.Password,
	}
}

// ToDomainInput maps authenticated profile update write intent into canonical
// domain input.
func (m AuthUpdateProfileWriteModel) ToDomainInput() authdomain.ProfileUpdateInput {
	return authdomain.ProfileUpdateInput{
		DisplayName: m.DisplayName,
	}
}

// ToDomainInput maps authenticated settings write intent into canonical domain
// input. JSON decoding is explicit so callers can preserve the current invalid
// payload behavior during handler alignment.
func (m AuthUpdateSettingsWriteModel) ToDomainInput() (authdomain.SettingsUpdateInput, error) {
	input := authdomain.SettingsUpdateInput{}
	if len(m.Preferences) == 0 {
		return input, nil
	}
	if err := json.Unmarshal(m.Preferences, &input.Preferences); err != nil {
		return authdomain.SettingsUpdateInput{}, err
	}
	return input, nil
}

// ToDomainInput maps public wallet challenge write intent into canonical domain
// input.
func (m AuthWalletChallengeWriteModel) ToDomainInput() authdomain.WalletChallengeInput {
	return authdomain.WalletChallengeInput{
		Address: m.Address,
		Chain:   m.Chain,
	}
}

// ToDomainInput maps public wallet verification write intent into canonical
// domain input.
func (m AuthWalletVerifyWriteModel) ToDomainInput() authdomain.WalletVerifyInput {
	return authdomain.WalletVerifyInput{
		ChallengeID: m.ChallengeID,
		Address:     m.Address,
		Signature:   m.Signature,
	}
}

// ToDomainInput maps authenticated wallet-link challenge write intent into
// canonical domain input.
func (m AuthWalletLinkChallengeWriteModel) ToDomainInput() authdomain.WalletChallengeInput {
	return authdomain.WalletChallengeInput{
		Address: m.Address,
		Chain:   m.Chain,
	}
}

// ToDomainInput maps authenticated wallet-link verification write intent into
// canonical domain input.
func (m AuthWalletLinkVerifyWriteModel) ToDomainInput() authdomain.WalletVerifyInput {
	return authdomain.WalletVerifyInput{
		ChallengeID: m.ChallengeID,
		Address:     m.Address,
		Signature:   m.Signature,
	}
}

// ToDomainInput maps authenticated account-merge challenge write intent into
// canonical domain input.
func (m AuthWalletAccountMergeChallengeWriteModel) ToDomainInput() authdomain.WalletChallengeInput {
	return authdomain.WalletChallengeInput{
		Address: m.Address,
		Chain:   m.Chain,
	}
}

// ToDomainInput maps authenticated account-merge verification write intent into
// canonical domain input.
func (m AuthWalletAccountMergeVerifyWriteModel) ToDomainInput() authdomain.WalletVerifyInput {
	return authdomain.WalletVerifyInput{
		ChallengeID: m.ChallengeID,
		Address:     m.Address,
		Signature:   m.Signature,
	}
}

// ToDomainInput maps wallet detach-check write intent into canonical domain
// input.
func (m AuthWalletDetachCheckWriteModel) ToDomainInput() authdomain.WalletDetachInput {
	return authdomain.WalletDetachInput{Address: m.Address}
}

// ToDomainInput maps wallet detach-execute write intent into canonical domain
// input.
func (m AuthWalletDetachExecuteWriteModel) ToDomainInput() authdomain.WalletDetachInput {
	return authdomain.WalletDetachInput{Address: m.Address}
}

// ToDomainInput maps primary wallet update write intent into canonical domain
// input.
func (m AuthWalletPrimarySetWriteModel) ToDomainInput() authdomain.WalletPrimarySetInput {
	return authdomain.WalletPrimarySetInput{Address: m.Address}
}
