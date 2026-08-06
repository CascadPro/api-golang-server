package requests_service

import (
	"context"
	"fmt"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
)

func (s *Service) GetRequests(
	ctx context.Context,
	page, limit int,
	filter domain.RequestQueryFilter,
) ([]domain.Request, int64, error) {
	if page <= 0 {
		return nil, 0, fmt.Errorf("`page` can't be below or equal zero: %w", core_errors.ErrInvalidArgument)
	}
	if limit <= 0 {
		return nil, 0, fmt.Errorf("`limit` can't be below or equal zero: %w", core_errors.ErrInvalidArgument)
	}
	if err := filter.Validate(); err != nil {
		return nil, 0, fmt.Errorf("validate `filter`: %w", err)
	}

	requests, total, err := s.requestsMongoRepo.GetRequests(ctx, page, limit, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("get requests from repository: %w", err)
	}

	ceil := (total + int64(limit) - 1) / int64(limit)

	return requests, ceil, nil
}
