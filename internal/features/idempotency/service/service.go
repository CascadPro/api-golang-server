package idempotency_service

import (
	"time"

	idempotency_postgres_repository "github.com/CascadePro/api-golang-server/internal/features/idempotency/repository/postgres"
)

const (
	HeaderKey      = "Idempotency-Key"
	HeaderReplayed = "Idempotency-Replayed"
	defaultTTL     = 24 * time.Hour
)

type Service struct {
	repository *idempotency_postgres_repository.Repository
	ttl        time.Duration
}

func New(repository *idempotency_postgres_repository.Repository, ttl time.Duration) *Service {
	if ttl <= 0 {
		ttl = defaultTTL
	}

	return &Service{
		repository: repository,
		ttl:        ttl,
	}
}
