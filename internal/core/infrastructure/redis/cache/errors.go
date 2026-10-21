package core_redis_cache

import "errors"

var (
	ErrInvalidNamespace = errors.New("cache namespace is invalid")
	ErrInvalidKey       = errors.New("cache key is empty")
	ErrInvalidTTL       = errors.New("cache TTL cannot be negative")
	ErrNilDestination   = errors.New("cache destination is nil")
	ErrNilContext       = errors.New("cache context is nil")
	ErrNilPool          = errors.New("cache Redis pool is nil")
	ErrNilLoader        = errors.New("cache loader is nil")
)
