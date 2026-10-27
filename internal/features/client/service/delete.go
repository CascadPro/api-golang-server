package client_service

import (
	"context"
	"fmt"

	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	"github.com/google/uuid"
)

func (s *Service) DeleteClient(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return fmt.Errorf("`ID` can't be NULL: %w", core_errors.ErrInvalidArgument)
	}

	if err := s.clientPostgresRepo.DeleteClient(ctx, id); err != nil {
		return fmt.Errorf("delete client: %w", err)
	}

	return nil
}
