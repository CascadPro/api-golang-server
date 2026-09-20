package session_service

import (
	"context"
	"fmt"

	core_context "github.com/CascadePro/api-golang-server/internal/core/context"
	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_logger "github.com/CascadePro/api-golang-server/internal/core/logger"
	"go.uber.org/zap"
)

func (s *Service) DeleteUserSessions(ctx context.Context) error {
	sessionID, err := core_context.SessionID(ctx)
	if err != nil {
		return fmt.Errorf("get sessionID from context: %w", err)
	}

	userID, err := core_context.UserID(ctx)
	if err != nil {
		return fmt.Errorf("get userID from context: %w", err)
	}

	if err := s.sessionsRedisRepo.DeleteUserSessions(ctx, userID, sessionID); err != nil {
		return fmt.Errorf("delete user sessions from repository: %w", err)
	}

	log := core_logger.FromContext(ctx)
	publishCtx := context.WithoutCancel(ctx)

	go func(log *core_logger.Logger) {
		data := domain.NewRealtimeEventSessionRevokeData(sessionID, domain.SessionRevokeDataUserRevoked)

		event := domain.NewRealtimeEvent(domain.RealtimeEventSessionsRevoked, &userID, nil, data)

		if err := s.publisher.Publish(publishCtx, event); err != nil {
			log.Error("publish session revoked event", zap.Error(err))
		}
	}(log)

	return nil
}
