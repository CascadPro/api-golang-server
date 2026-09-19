package presence_redis_repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type UnregisterResult struct {
	SessionOffline bool
	UserOffline    bool
}

func (r *Repository) Unregister(ctx context.Context, userID uuid.UUID, sessionID, connectionID string) (UnregisterResult, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	result, err := r.pool.Eval(
		ctx,
		unregisterScript,
		[]string{presenceKey(userID)},
		sessionID,
		connectionID,
		time.Now().UnixMilli(),
	)
	if err != nil {
		return UnregisterResult{}, fmt.Errorf("unregister presence: %w", err)
	}

	values, ok := result.([]interface{})
	if !ok {
		return UnregisterResult{}, fmt.Errorf("unexpected unregister presence result: %T", result)
	}

	if len(values) != 2 {
		return UnregisterResult{}, fmt.Errorf("unexpected unregister presence result length: %d", len(values))
	}

	remainingSession, ok := values[0].(int64)
	if !ok {
		return UnregisterResult{}, fmt.Errorf("unexpected remaining session result: %T", values[0])
	}

	remainingUser, ok := values[1].(int64)
	if !ok {
		return UnregisterResult{}, fmt.Errorf("unexpected remaining user result: %T", values[1])
	}

	return UnregisterResult{
		SessionOffline: remainingSession == 0,
		UserOffline:    remainingUser == 0,
	}, nil
}
