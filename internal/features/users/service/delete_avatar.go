package users_service

import (
	"context"
	"fmt"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	"github.com/google/uuid"
)

func (s *Service) DeleteAvatar(ctx context.Context, userID uuid.UUID) error {
	if userID == uuid.Nil {
		return fmt.Errorf("`userID` can't be NULL: %w", core_errors.ErrInvalidArgument)
	}

	user, err := s.usersPostgresRepo.GetUser(ctx, domain.User{ID: userID})
	if err != nil {
		return fmt.Errorf("get user from repository: %w", err)
	}
	if user.AvatarFileID == nil {
		return fmt.Errorf("user has no avatar: %w", core_errors.ErrNotFound)
	}

	if err := s.mediaService.DeleteFile(ctx, domain.FileTagAvatars, *user.AvatarFileID); err != nil {
		return fmt.Errorf("delete avatar file: %w", err)
	}

	return nil
}
