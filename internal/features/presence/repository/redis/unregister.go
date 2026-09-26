package presence_redis_repository

import (
	"context"
	"fmt"
	"time"

	core_redis_pool "github.com/CascadePro/api-golang-server/internal/core/infrastructure/redis/pool"
	"github.com/google/uuid"
)

type UnregisterResult struct {
	RemainingSession int
	RemainingUser    int
}

func (r *Repository) Unregister(ctx context.Context, userID uuid.UUID, sessionID, connectionID string) (UnregisterResult, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	member := presenceMember(sessionID, connectionID)

	sessionPrefix := fmt.Sprintf("%s:%s:", core_redis_pool.PresenceFolder, sessionID)

	result, err := r.pool.Eval(
		ctx,
		unregisterScript,
		[]string{presenceKey(userID)},
		time.Now().UnixMilli(),
		member,
		sessionPrefix,
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
		RemainingSession: int(remainingSession),
		RemainingUser:    int(remainingUser),
	}, nil
}
