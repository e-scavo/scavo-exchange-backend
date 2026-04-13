package httpx

import (
	"net/http"

	coreerrs "github.com/e-scavo/scavo-exchange-backend/internal/core/errs"
)

type ErrorResponse struct {
	Error coreerrs.ResponseError `json:"error"`
}

func WriteError(w http.ResponseWriter, status int, err coreerrs.ResponseError) {
	WriteJSON(w, status, ErrorResponse{Error: err})
}

func WriteErrorMessage(w http.ResponseWriter, status int, code, message string, details map[string]any) {
	WriteError(w, status, coreerrs.NewResponseError(code, message, details))
}
