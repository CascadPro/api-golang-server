package media_service

import (
	"context"
	"encoding/json"
	"fmt"

	core_context "github.com/CascadePro/api-golang-server/internal/core/context"
	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_validation "github.com/CascadePro/api-golang-server/internal/core/validation"
)

func (s *Service) DeleteFile(ctx context.Context, fileTag domain.FileTag, fileID string) error {
	if err := core_validation.ValidateArray(domain.FileTags, fileTag); err != nil {
		return fmt.Errorf("validate file tag: %w", err)
	}

	if err := core_validation.ValidateID(fileID, domain.FileIDByteLength); err != nil {
		return fmt.Errorf("validate file id: %w", err)
	}

	tx, err := s.mediaPostgresRepo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	defer func() {
		_ = tx.Rollback(ctx)
	}()

	if err := s.mediaPostgresRepo.DeleteFileTx(ctx, tx, fileID); err != nil {
		_ = tx.Rollback(ctx)

		return fmt.Errorf("delete file from repository: %w", err)
	}

	payload, err := deleteFileEventPayload(domain.FileTagAvatars, fileID)
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

func deleteFileEventPayload(tag domain.FileTag, fileID string) ([]byte, error) {
	payload, err := json.Marshal(
		domain.NewEventTypeMediaDeleteFilePayload(tag, fileID),
	)
	if err != nil {
		return nil, fmt.Errorf("marshal outbox payload: %w", err)
	}

	return payload, nil
}
