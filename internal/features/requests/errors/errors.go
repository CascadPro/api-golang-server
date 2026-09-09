package requests_errors

import "errors"

var (
	ErrRequestNotFound       = errors.New("request not found")
	ErrRequestAlreadyExists  = errors.New("request already exists")
	ErrRequestAccessDenied   = errors.New("request access denied")
	ErrRequestAlreadyHandled = errors.New("request already handled")
)
