package presence_service

import (
	"context"
	"errors"
	"fmt"

	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	presence_redis_repository "github.com/CascadePro/api-golang-server/internal/features/presence/repository/redis"
	"github.com/google/uuid"
)

type RegisterResult struct {
	WasSessionOnline bool
	WasUserOnline    bool
}

type UnregisterResult struct {
	SessionOffline bool
	UserOffline    bool
}

func (s *Service) Register(ctx context.Context, userID uuid.UUID, sessionID, connectionID string) (RegisterResult, error) {
	result, err := s.presenceRedisRepo.Register(ctx, userID, sessionID, connectionID)
	if err != nil {
		return RegisterResult{}, fmt.Errorf("register presence in repository: %w", err)
	}

	if err := s.sessionsRedisRepo.PatchLastActive(ctx, userID, sessionID); err != nil {
		if !errors.Is(err, core_errors.ErrNotFound) {
			return RegisterResult{}, fmt.Errorf("patch last active in repository: %w", err)
		}
	}

	return parseRegisterResult(result), nil
}

func (s *Service) Unregister(ctx context.Context, userID uuid.UUID, sessionID, connectionID string) (UnregisterResult, error) {
	result, err := s.presenceRedisRepo.Unregister(ctx, userID, sessionID, connectionID)
	if err != nil {
		return UnregisterResult{}, fmt.Errorf("unregister presence in repository: %w", err)
	}

	if err := s.sessionsRedisRepo.PatchLastActive(ctx, userID, sessionID); err != nil {
		if !errors.Is(err, core_errors.ErrNotFound) {
			return UnregisterResult{}, fmt.Errorf("patch last active in repository: %w", err)
		}
	}

	return parseUnregisterResult(result), nil
}

func parseRegisterResult(result presence_redis_repository.RegisterResult) RegisterResult {
	return RegisterResult{
		WasSessionOnline: result.WasSessionOnline != 0,
		WasUserOnline:    result.WasUserOnline != 0,
	}
}

func parseUnregisterResult(result presence_redis_repository.UnregisterResult) UnregisterResult {
	return UnregisterResult{
		SessionOffline: result.RemainingSession == 0,
		UserOffline:    result.RemainingUser == 0,
	}
}
