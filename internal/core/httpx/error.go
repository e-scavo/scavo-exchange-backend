package httpx

import (
	"net/http"

	coreerrs "github.com/e-scavo/scavo-exchange-backend/internal/core/errs"
)

func WriteError(w http.ResponseWriter, status int, err coreerrs.ResponseError) {
	WriteJSON(w, status, coreerrs.NewErrorEnvelope(err))
}

func WriteErrorMessage(w http.ResponseWriter, status int, code, message string, details map[string]any) {
	WriteError(w, status, coreerrs.NewResponseError(code, message, details))
}

func WriteAppError(w http.ResponseWriter, appErr *coreerrs.AppError) {
	if appErr == nil {
		appErr = coreerrs.InternalError(nil)
	}
	WriteError(w, appErr.Status, appErr.ToResponseError())
}
