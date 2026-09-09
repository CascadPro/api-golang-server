package core_errors

import "fmt"

type APIError struct {
	Code       Code
	Err        error
	Message    string
	StatusCode int
	Fields     map[string]string
}

