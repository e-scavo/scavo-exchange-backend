package errs

type ResponseError struct {
	Code    string         `json:"code"`
	Message string         `json:"message"`
	Details map[string]any `json:"details,omitempty"`
}

type ErrorEnvelope struct {
	Error ResponseError `json:"error"`
}

func NewResponseError(code, message string, details map[string]any) ResponseError {
	err := ResponseError{
		Code:    code,
		Message: message,
	}
	if len(details) > 0 {
		err.Details = details
	}
	return err
}

func NewErrorEnvelope(err ResponseError) ErrorEnvelope {
	return ErrorEnvelope{Error: err}
}
