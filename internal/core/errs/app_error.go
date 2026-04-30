package errs

import "errors"

type Category string

const (
	CategoryGeneric  Category = "generic"
	CategoryAuth     Category = "auth"
	CategoryWallet   Category = "wallet"
	CategorySettings Category = "settings"
	CategoryInternal Category = "internal"
)

type AppError struct {
	Code     string
	Message  string
	Status   int
	Category Category
	Details  map[string]any
	Cause    error
}

func (e *AppError) Error() string {
	if e == nil {
		return ""
	}
	if e.Message != "" {
		return e.Message
	}
	return e.Code
}

func (e *AppError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Cause
}

func (e *AppError) WithDetails(details map[string]any) *AppError {
	if e == nil {
		return nil
	}
	if len(details) == 0 {
		return e
	}

	merged := e.PublicDetails()
	for key, value := range details {
		merged[key] = value
	}

	copy := *e
	copy.Details = merged
	return &copy
}

func (e *AppError) WithContext(key string, value any) *AppError {
	if e == nil || key == "" {
		return e
	}
	return e.WithDetails(map[string]any{key: value})
}

func (e *AppError) WithRequestID(requestID string) *AppError {
	if requestID == "" {
		return e
	}
	return e.WithContext("request_id", requestID)
}

func (e *AppError) PublicDetails() map[string]any {
	if e == nil || len(e.Details) == 0 {
		return nil
	}

	details := make(map[string]any, len(e.Details))
	for key, value := range e.Details {
		details[key] = value
	}
	return details
}

func (e *AppError) ToResponseError() ResponseError {
	if e == nil {
		return NewResponseError("INTERNAL_ERROR", "internal server error", nil)
	}
	return NewResponseError(e.Code, e.Message, e.PublicDetails())
}

func New(code, message string, status int, category Category) *AppError {
	return &AppError{
		Code:     code,
		Message:  message,
		Status:   status,
		Category: category,
	}
}

func Wrap(cause error, code, message string, status int, category Category) *AppError {
	err := New(code, message, status, category)
	err.Cause = cause
	return err
}

func IsAppError(err error) bool {
	var appErr *AppError
	return errors.As(err, &appErr)
}

func AsAppError(err error) (*AppError, bool) {
	var appErr *AppError
	if !errors.As(err, &appErr) || appErr == nil {
		return nil, false
	}
	return appErr, true
}
