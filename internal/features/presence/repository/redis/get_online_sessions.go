package presence_redis_repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func (r *Repository) GetOnlineSessions(ctx context.Context, userID uuid.UUID, sessionIDs []string) (map[string]bool, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	result := make(map[string]bool, len(sessionIDs))

	for _, sessionID := range sessionIDs {
		result[sessionID] = false
	}

	if len(sessionIDs) == 0 {
		return result, nil
	}

	raw, err := r.pool.Eval(
		ctx,
		getOnlineSessionsScript,
		[]string{presenceKey(userID)},
		time.Now().UnixMilli(),
	)
	if err != nil {
		return nil, fmt.Errorf("get online sessions: %w", err)
	}

	values, ok := raw.([]any)
	if !ok {
		return nil, fmt.Errorf("unexpected online sessions result: %T", raw)
	}

	for _, value := range values {
		member, ok := value.(string)
		if !ok {
			continue
		}

		sessionID, _, ok := parsePresenceMember(member)
		if !ok {
			continue
		}

		if _, exists := result[sessionID]; exists {
			result[sessionID] = true
		}
	}

	return result, nil
}
