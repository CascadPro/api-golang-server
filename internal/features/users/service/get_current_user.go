package users_service

import (
	"context"
	"fmt"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	"github.com/google/uuid"
)

func (s *Service) GetCurrentUser(ctx context.Context, id uuid.UUID) (domain.User, error) {
	if id == uuid.Nil {
		return domain.User{}, fmt.Errorf("`id` can't be NULL: %w", core_errors.ErrInvalidArgument)
	}

	user, err := s.usersPostgresRepo.GetUser(ctx, domain.User{ID: id})
	if err != nil {
		return domain.User{}, fmt.Errorf("get user from repository: %w", err)
	}

	if !user.Activated {
		return domain.User{}, fmt.Errorf("you must activate your account before usage: %w", core_errors.ErrConflict)
	}

	return user, nil
}
