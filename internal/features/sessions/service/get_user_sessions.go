package session_service

import (
	"context"
	"fmt"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_utils "github.com/CascadePro/api-golang-server/internal/core/utils"
)

func (s *Service) GetUserSessions(ctx context.Context) ([]domain.Session, error) {
	userID, err := core_utils.ParseUUIDFromContext(ctx, "userID")
	if err != nil {
		return nil, fmt.Errorf("parse userID from context: %w", err)
	}

	sessionID, err := core_utils.ParseIDStingFromContext(ctx, "sessionID", domain.SessionIDByteLength)
	if err != nil {
		return nil, fmt.Errorf("parse session id from context: %w", err)
	}

	sessions, err := s.sessionsRedisRepo.GetUserSessions(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user sessions: %w", err)
	}

	var idx int = -1
	for i, s := range sessions {
		if s.ID == sessionID {
			idx = i
			break
		}
	}

	if idx > 0 {
		// Current session get
		current := sessions[idx]

		// 						 " "   :  "idx"                   "idx+1"   :   " "
		// Slice from start to current, plus slice from current+1 to  end
		sessions = append(sessions[:idx], sessions[idx+1:]...)

		// Placing current session at the start
		sessions = append([]domain.Session{current}, sessions...)
	}

	return sessions, nil
}
