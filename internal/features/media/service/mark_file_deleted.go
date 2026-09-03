package media_service

import (
	"context"
	"fmt"
	"time"

	core_context "github.com/CascadePro/api-golang-server/internal/core/context"
	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_validation "github.com/CascadePro/api-golang-server/internal/core/validation"
)

func (s *Service) MarkFileDeleted(ctx context.Context, fileID string) error {
	if err := core_validation.ValidateID(fileID, domain.FileIDByteLength); err != nil {
		return fmt.Errorf("validate file id: %w", err)
	}

	file, _, err := s.mediaPostgresRepo.GetFile(ctx, fileID)
	if err != nil {
		return fmt.Errorf("get file from repository: %w", err)
	}

	tx, err := s.mediaPostgresRepo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	defer func() {
		_ = tx.Rollback(ctx)
	}()

	now := time.Now()
	file.Deleted = true
	file.DeletedAt = &now

	if _, err := s.mediaPostgresRepo.PatchFileTx(ctx, tx, fileID, file); err != nil {
		return fmt.Errorf("patch file in repository: %w", err)
	}

	payload, err := deleteFileEventPayload(domain.FileTagAvatars, file.ID)
	if err != nil {
		_ = tx.Rollback(ctx)

		return fmt.Errorf("create outbox payload: %w", err)
	}

	event := domain.NewOutboxEvent(domain.EventTypeMediaDeleteFile, nil, payload)
	if requestID, err := core_context.RequestIDFromContext(ctx); err == nil {
		event.AggregateID = &requestID
	}

	if _, err := s.outboxPostgresRepo.CreateEvent(ctx, tx, event); err != nil {
		_ = tx.Rollback(ctx)

		return fmt.Errorf("create outbox event: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}
