package core_errors

import "errors"

var (
	ErrInvalidArgument  = errors.New("invalid argument")
	ErrEmptyRequestBody = errors.New("request body is empty")
	ErrUnauthorized     = errors.New("unauthorized")
	ErrForbidden        = errors.New("don't have access")
	ErrNotFound         = errors.New("not found")
	ErrConflict         = errors.New("conflict")
	ErrTooManyRequests  = errors.New("too many requests")
)
