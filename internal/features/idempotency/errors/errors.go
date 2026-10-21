package idempotency_error

import "errors"

var (
	ErrInvalidIdempotencyKey = errors.New("invalid idempotency key")
	ErrRequestBodyTooLarge   = errors.New("request body too large")
	ErrIdempotencyKeyReused  = errors.New("idempotency key already used")
	ErrIdempotencyInProgress = errors.New("idempotency in progress")
)
