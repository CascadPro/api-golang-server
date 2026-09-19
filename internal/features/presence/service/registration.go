package presence_service

import (
	"context"
	"fmt"

	presence_redis_repository "github.com/CascadePro/api-golang-server/internal/features/presence/repository/redis"
	"github.com/google/uuid"
)

type UnregisterResult struct {
	SessionOffline bool
	UserOffline    bool
}

func (s *Service) Register(ctx context.Context, userID uuid.UUID, sessionID, connectionID string) (bool, error) {
	ok, err := s.presenceRedisRepo.Register(ctx, userID, sessionID, connectionID)
	if err != nil {
		return false, fmt.Errorf("register presence in repository: %w", err)
	}

	if err := s.sessionsRedisRepo.PatchLastActive(ctx, userID, sessionID); err != nil {
		return false, fmt.Errorf("patch last active in repository: %w", err)
	}

	return ok, nil
}

func (s *Service) Unregister(ctx context.Context, userID uuid.UUID, sessionID, connectionID string) (UnregisterResult, error) {
	result, err := s.presenceRedisRepo.Unregister(ctx, userID, sessionID, connectionID)
	if err != nil {
		return UnregisterResult{}, fmt.Errorf("unregister presence in repository: %w", err)
	}

	if err := s.sessionsRedisRepo.PatchLastActive(ctx, userID, sessionID); err != nil {
		return UnregisterResult{}, fmt.Errorf("patch last active in repository: %w", err)
	}

	return parseUnregisterResult(result), nil
}

func parseUnregisterResult(result presence_redis_repository.UnregisterResult) UnregisterResult {
	return UnregisterResult{
		SessionOffline: result.SessionOffline,
		UserOffline:    result.UserOffline,
	}
}
