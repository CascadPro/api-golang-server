package core_errors

type ValidationError struct {
	Fields map[string]string
}

func NewValidationError(
	fields map[string]string,
) *ValidationError {
	return &ValidationError{
		Fields: fields,
	}
}

func (e *ValidationError) Error() string {
	return "validation failed"
}
