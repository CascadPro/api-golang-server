package users_service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	idempotency_service "github.com/CascadePro/api-golang-server/internal/features/idempotency/service"
	users_errors "github.com/CascadePro/api-golang-server/internal/features/users/errors"
	"github.com/google/uuid"
	"go.uber.org/zap"
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
		return fmt.Errorf("user has no avatar: %w", users_errors.ErrAvatarNotFound)
	}

	tx, err := s.usersPostgresRepo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("begin avatar deletion transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	patch := domain.NewAvatarUserPatch(domain.Nullable[string]{Set: true})
	patched, err := user.ApplyPatch(patch)
	if err != nil {
		return fmt.Errorf("prepare avatar removal: %w", err)
	}

	if _, err := s.usersPostgresRepo.PatchUserTx(ctx, tx, userID, patched); err != nil {
		return fmt.Errorf("clear user avatar: %w", err)
	}

	if err := s.mediaService.DeleteFileTx(ctx, tx, domain.FileTagAvatars, *user.AvatarFileID); err != nil {
		return fmt.Errorf("schedule avatar file deletion: %w", err)
	}

	if err := idempotency_service.CompleteInTransaction(ctx, tx, http.StatusNoContent, nil, http.Header{}); err != nil {
		return fmt.Errorf("save idempotency response in transaction: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit avatar deletion transaction: %w", err)
	}

	idempotency_service.MarkTransactionCommitted(ctx)

	realtimeEvent := getRealtimeEvent(nil, &userID)

	if err := s.realtimePublisher.Publish(ctx, realtimeEvent); err != nil {
		s.logger.Error(
			"failed to publish user avatar updated event",
			zap.String("user_id", userID.String()),
			zap.Error(err),
		)
	}

	if err := s.cache.Delete(ctx, cacheKeyUser(userID)); err != nil {
		s.logger.Error(
			"failed to invalidate user cache after avatar deletion",
			zap.String("user_id", userID.String()),
			zap.Error(err),
		)
	}

	return nil
}
