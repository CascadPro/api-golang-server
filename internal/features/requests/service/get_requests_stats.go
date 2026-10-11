package requests_service

import (
	"context"
	"fmt"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
)

func (s *Service) GetRequestsStats(ctx context.Context) (domain.RequestStats, error) {
	stats, err := s.requestsMongoRepo.GetRequestsStats(ctx)
	if err != nil {
		return domain.RequestStats{}, fmt.Errorf("get request stats from repository: %w", err)
	}

	return stats, nil
}
