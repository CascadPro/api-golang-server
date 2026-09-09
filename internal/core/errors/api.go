package core_errors

import "fmt"

type APIError struct {
	Code       Code
	Err        error
	Message    string
	StatusCode int
	Fields     map[string]string
}

func (e *APIError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}

	return string(e.Code)
}

func (e *APIError) Unwrap() error {
	return e.Err
}

func NewAPIError(
	code Code,
	statusCode int,
	message string,
	err error,
) *APIError {
	return &APIError{
		Code:       code,
		StatusCode: statusCode,
		Message:    message,
		Err:        err,
	}
}

func (e *APIError) WithField(field, message string) *APIError {
	if e.Fields == nil {
		e.Fields = make(map[string]string)
	}

	e.Fields[field] = message

	return e
}

func (e *APIError) WithFields(fields map[string]string) *APIError {
	e.Fields = fields

	return e
}

func (e *APIError) Is(target error) bool {
	apiErr, ok := target.(*APIError)
	if !ok {
		return false
	}

	return e.Code == apiErr.Code
}

func (e *APIError) String() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}
