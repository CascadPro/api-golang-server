package presence_service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (s *Service) Heartbeat(ctx context.Context, userID uuid.UUID, sessionID, connectionID string) error {
	err := s.presenceRedisRepo.Heartbeat(ctx, userID, sessionID, connectionID)
	if err != nil {
		return fmt.Errorf("repository heartbeat: %w", err)
	}

	return nil
}
