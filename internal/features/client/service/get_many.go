package client_service

import (
	"context"
	"fmt"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
)

type GetManyResult struct {
	Clients []domain.Client
	Length  int64
}

func (s *Service) GetMany(ctx context.Context, page, limit int64, searchQuery *string) (GetManyResult, error) {
	if page <= 0 {
		return GetManyResult{}, fmt.Errorf("`page` can't be below or equal zero: %w", core_errors.ErrInvalidArgument)
	}
	if limit <= 0 {
		return GetManyResult{}, fmt.Errorf("`limit` can't be below or equal zero: %w", core_errors.ErrInvalidArgument)
	}
	if searchQuery == nil {
		searchQuery = new(string)
	}

	clients, total, err := s.clientPostgresRepo.GetMany(ctx, page, limit, searchQuery)
	if err != nil {
		return GetManyResult{}, fmt.Errorf("get many clients from repository: %w", err)
	}

	return GetManyResult{
		Clients: clients,
		Length:  total,
	}, nil
}
