package users_service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	users_errors "github.com/CascadePro/api-golang-server/internal/features/users/errors"
	"github.com/google/uuid"
)

func (s *Service) UpdateAvatar(ctx context.Context, userID uuid.UUID, uploadedFile *domain.File, content []byte) error {
	if userID == uuid.Nil {
		return fmt.Errorf("`userID` can't be NULL: %w", core_errors.ErrInvalidArgument)
	}

	if uploadedFile == nil {
		return fmt.Errorf("`uploadedFile` can't be NULL: %w", core_errors.ErrInvalidArgument)
	}
	if err := uploadedFile.Validate(); err != nil {
		return fmt.Errorf("validate file: %w", core_errors.ErrInvalidArgument)
	}

	user, err := s.usersPostgresRepo.GetUser(ctx, domain.User{ID: userID})
	if err != nil {
		return fmt.Errorf("get user from repository: %w", err)
	}
	if !user.Activated {
		return fmt.Errorf("user is not activated: %w", users_errors.ErrUserNotActivated)
	}

	oldAvatarID := user.AvatarFileID

	// 1. Upload new file first.
	file, err := s.mediaService.UploadFile(ctx, uploadedFile, content)
	if err != nil {
		return fmt.Errorf("upload new avatar: %w", err)
	}

	// 2. Prepare domain patch.
	patch := domain.NewAvatarUserPatch(domain.NewNullable(file.ID))

	patched, err := user.ApplyPatch(patch)
	if err != nil {
		// New file is no longer referenced.
		_ = s.mediaService.DeleteFile(ctx, domain.FileTagAvatars, file.ID)

		return fmt.Errorf("apply user patch: %w", err)
	}

	// 3. Begin PostgreSQL transaction.
	tx, err := s.usersPostgresRepo.BeginTx(ctx)
	if err != nil {
		_ = s.mediaService.DeleteFile(ctx, domain.FileTagAvatars, file.ID)

		return fmt.Errorf("begin user transaction: %w", err)
	}

	defer func() {
		_ = tx.Rollback(ctx)
	}()

	payload, err := updateAvatarEventPayload(file.ID)
	if err != nil {
		_ = tx.Rollback(ctx)

		_ = s.mediaService.DeleteFile(ctx, domain.FileTagAvatars, file.ID)

		return fmt.Errorf("create outbox payload: %w", err)
	}

	event := domain.NewOutboxEvent(
		domain.EventTypeMediaAvatarProcess,
		&userID,
		payload,
	)

	if _, err := s.outboxPostgresRepo.CreateEventTx(ctx, tx, event); err != nil {
		_ = tx.Rollback(ctx)

		_ = s.mediaService.DeleteFile(ctx, domain.FileTagAvatars, file.ID)

		return fmt.Errorf("create outbox event: %w", err)
	}

	// 4. Update user inside transaction.
	if _, err := s.usersPostgresRepo.PatchUserTx(ctx, tx, userID, patched); err != nil {
		_ = tx.Rollback(ctx)

		// Compensation: new file is not referenced.
		_ = s.mediaService.DeleteFile(ctx, domain.FileTagAvatars, file.ID)

		return fmt.Errorf("update user in transaction: %w", err)
	}

	// 6. If old avatar exists, create deletion event.
	if oldAvatarID != nil {
		if err := s.mediaService.DeleteFileTx(ctx, tx, domain.FileTagAvatars, *oldAvatarID); err != nil {
			_ = tx.Rollback(ctx)

			return fmt.Errorf("delete old avatar: %w", err)
		}
	}

	// 6. Commit.
	if err := tx.Commit(ctx); err != nil {
		// DB transaction didn't commit.
		// New S3 object is potentially orphaned.
		// Cleanup is best-effort here.
		_ = s.mediaService.DeleteFile(ctx, domain.FileTagAvatars, file.ID)

		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

func updateAvatarEventPayload(fileID string) ([]byte, error) {
	payload, err := json.Marshal(
		domain.NewEventMediaAvatarProcessPayload(fileID),
	)
	if err != nil {
		return nil, fmt.Errorf("marshal outbox payload: %w", err)
	}

	return payload, nil
}
