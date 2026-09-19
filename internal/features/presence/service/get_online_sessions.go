package presence_service

import (
	"context"
	"fmt"

	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	"github.com/google/uuid"
)

func (s *Service) GetOnlineSessions(ctx context.Context, userID uuid.UUID, sessionIDs []string) (map[string]bool, error) {
	if userID == uuid.Nil {
		return nil, fmt.Errorf("userID can't be NULL: %w", core_errors.ErrInvalidArgument)
	}

	if len(sessionIDs) == 0 {
		return nil, fmt.Errorf("sessionIDs can't be empty: %w", core_errors.ErrInvalidArgument)
	}

	onlineSessions, err := s.presenceRedisRepo.GetOnlineSessions(ctx, userID, sessionIDs)
	if err != nil {
		return nil, fmt.Errorf("get from repository: %w", err)
	}

	return onlineSessions, nil
}
