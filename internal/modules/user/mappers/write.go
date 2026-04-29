package mappers

import (
	userdomain "github.com/e-scavo/scavo-exchange-backend/internal/modules/user/domain"
	userwritemodels "github.com/e-scavo/scavo-exchange-backend/internal/modules/user/writemodels"
)

// UserUpdateWriteToDomainInput maps user update write intent into canonical
// domain input.
func UserUpdateWriteToDomainInput(model userwritemodels.UserUpdateWriteModel) userdomain.UserUpdateInput {
	return userdomain.UserUpdateInput{DisplayName: model.DisplayName}
}
