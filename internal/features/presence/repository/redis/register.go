package presence_redis_repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func (r *Repository) Register(ctx context.Context, userID uuid.UUID, sessionID, connectionID string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	now := time.Now()

	result, err := r.pool.Eval(
		ctx,
		registerScript,
		[]string{presenceKey(userID)},
		now.UnixMilli(),
		now.Add(LeaseTTL).UnixMilli(),
		presenceMember(sessionID, connectionID),
	)
	if err != nil {
		return false, fmt.Errorf("register presence: %w", err)
	}

	value, ok := result.(int64)
	if !ok {
		return false, fmt.Errorf("unexpected register presence result: %T", result)
	}

	return value == 1, nil
}
