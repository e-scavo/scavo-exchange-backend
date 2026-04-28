package writemodels

import userdomain "github.com/e-scavo/scavo-exchange-backend/internal/modules/user/domain"

// ToDomainInput maps user update write intent into canonical domain input.
func (m UserUpdateWriteModel) ToDomainInput() userdomain.UserUpdateInput {
	return userdomain.UserUpdateInput{
		DisplayName: m.DisplayName,
	}
}
