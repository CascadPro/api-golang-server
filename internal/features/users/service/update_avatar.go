package users_service

import (
	"context"
	"fmt"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	"github.com/google/uuid"
)

func (s *Service) UpdateAvatar(ctx context.Context, userID uuid.UUID, file *domain.File, content []byte) error {
	if userID == uuid.Nil {
		return fmt.Errorf("`userID` can't be NULL: %w", core_errors.ErrInvalidArgument)
	}
	if err := file.Validate(); err != nil {
		return fmt.Errorf("validate file: %w", core_errors.ErrInvalidArgument)
	}

	user, err := s.usersPostgresRepo.GetUser(ctx, domain.User{ID: userID})
	if err != nil {
		return fmt.Errorf("get user from repository: %w", err)
	}
	if !user.Activated {
		return fmt.Errorf("user is not activated: %w", core_errors.ErrConflict)
	}

	if user.AvatarFileID != nil {
		if err := s.mediaService.DeleteFile(ctx, domain.FileTagAvatars, *user.AvatarFileID); err != nil {
			return fmt.Errorf("delete existing avatar: %w", err)
		}
	}

	patch := domain.NewAvatarUserPatch(domain.NewNullable(file.ID))
	patched, err := user.ApplyPatch(patch)
	if err != nil {
		return fmt.Errorf("apply user patch: %w", err)
	}

	if _, err := s.usersPostgresRepo.PatchUser(ctx, userID, patched); err != nil {
		return fmt.Errorf("update user in repository: %w", err)
	}

	return nil
}
