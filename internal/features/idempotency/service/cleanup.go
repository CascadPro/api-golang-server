package idempotency_service

import "context"

func (s *Service) CleanupExpired(ctx context.Context) (int64, error) {
	return s.repository.CleanupExpired(ctx)
}
