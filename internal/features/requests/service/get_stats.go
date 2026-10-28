package requests_service

import (
	"context"
	"errors"
	"fmt"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
)

func (s *Service) GetRequestsStats(ctx context.Context) (domain.RequestStats, error) {
	var stats domain.RequestStats

	err := s.cache.Get(ctx, cacheKeyStats, &stats)
	if err != nil && !errors.Is(err, core_errors.ErrNotFound) {
		return domain.RequestStats{}, fmt.Errorf("get redis cache: %w", err)
	}

	if !errors.Is(err, core_errors.ErrNotFound) {
		if err := stats.Validate(); err == nil {
			return stats, nil
		}
	}

	stats, err = s.requestsMongoRepo.GetStats(ctx)
	if err != nil {
		return domain.RequestStats{}, fmt.Errorf("get request stats from repository: %w", err)
	}

	if err := s.cache.Set(ctx, cacheKeyStats, stats, cacheTTLStats); err != nil {
		return domain.RequestStats{}, fmt.Errorf("set redis cache: %w", err)
	}

	return stats, nil
}
