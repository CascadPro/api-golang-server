package client_service

import (
	"context"
	"fmt"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	"github.com/google/uuid"
)

func (s *Service) GetClient(ctx context.Context, id uuid.UUID) (domain.Client, error) {
	if id == uuid.Nil {
		return domain.Client{}, fmt.Errorf("`id` can't be NULL: %w", core_errors.ErrInvalidArgument)
	}

	client, err := s.clientPostgresRepo.GetClient(ctx, id)
	if err != nil {
		return domain.Client{}, fmt.Errorf("get client from repository: %w", err)
	}

	return client, nil
}
