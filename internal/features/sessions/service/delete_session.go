package session_service

import (
	"context"
	"fmt"

	core_context "github.com/CascadePro/api-golang-server/internal/core/context"
	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_logger "github.com/CascadePro/api-golang-server/internal/core/logger"
	core_validation "github.com/CascadePro/api-golang-server/internal/core/validation"
	"go.uber.org/zap"
)

func (s *Service) DeleteSession(ctx context.Context, sessionID string) error {
	if err := core_validation.ValidateID(sessionID, domain.SessionIDByteLength); err != nil {
		return fmt.Errorf("validate session id: %w", err)
	}

	userID, err := core_context.UserID(ctx)
	if err != nil {
		return fmt.Errorf("get userID from context: %w", err)
	}

	if err := s.sessionsRedisRepo.DeleteSession(ctx, userID, sessionID); err != nil {
		return fmt.Errorf("delete session from repository: %w", err)
	}

	log := core_logger.FromContext(ctx)

	go func(log *core_logger.Logger) {
		event := domain.NewRealtimeEvent(domain.RealtimeEventSessionRevoked, nil, &sessionID, nil)
		if err := s.publisher.Publish(ctx, event); err != nil {
			log.Error("publish session revoked event", zap.Error(err))
		}
	}(log)

	return nil
}
