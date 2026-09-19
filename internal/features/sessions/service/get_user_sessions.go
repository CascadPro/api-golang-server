package session_service

import (
	"context"
	"fmt"

	core_context "github.com/CascadePro/api-golang-server/internal/core/context"
)

func (s *Service) GetUserSessions(ctx context.Context) ([]Session, error) {
	userID, err := core_context.UserID(ctx)
	if err != nil {
		return nil, fmt.Errorf("get userID from context: %w", err)
	}

	sessionID, err := core_context.SessionID(ctx)
	if err != nil {
		return nil, fmt.Errorf("get sessionID from context: %w", err)
	}

	sessions, err := s.sessionsRedisRepo.GetUserSessions(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user sessions: %w", err)
	}

	sessionIDs := make([]string, 0, len(sessions))

	for _, session := range sessions {
		sessionIDs = append(sessionIDs, session.ID)
	}

	onlineSessions, err := s.presenceService.GetOnlineSessions(ctx, userID, sessionIDs)
	if err != nil {
		return nil, fmt.Errorf("get online sessions: %w", err)
	}

	result := make([]Session, 0, len(sessions))

	for _, session := range sessions {
		result = append(
			result,
			Session{
				Session: session,
				Online:  onlineSessions[session.ID],
			},
		)
	}

	for i := range result {
		if result[i].ID != sessionID {
			continue
		}

		if i == 0 {
			break
		}

		current := result[i]

		copy(result[1:i+1], result[0:i])

		result[0] = current

		break
	}

	return result, nil
}
