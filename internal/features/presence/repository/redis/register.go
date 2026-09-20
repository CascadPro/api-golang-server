package presence_redis_repository

import (
	"context"
	"fmt"
	"time"

	core_redis_pool "github.com/CascadePro/api-golang-server/internal/core/infrastructure/redis/pool"
	"github.com/google/uuid"
)

type RegisterResult struct {
	WasSessionOnline int
	WasUserOnline    int
}

func (r *Repository) Register(ctx context.Context, userID uuid.UUID, sessionID, connectionID string) (RegisterResult, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	now := time.Now()

	member := presenceMember(sessionID, connectionID)
	sessionPrefix := fmt.Sprintf("%s:%s:", core_redis_pool.PresenceFolder, sessionID)

	result, err := r.pool.Eval(
		ctx,
		registerScript,
		[]string{presenceKey(userID)},
		now.UnixMilli(),
		now.Add(LeaseTTL).UnixMilli(),
		member,
		sessionPrefix,
	)
	if err != nil {
		return RegisterResult{}, fmt.Errorf("register presence: %w", err)
	}

	values, ok := result.([]interface{})
	if !ok {
		return RegisterResult{}, fmt.Errorf("unexpected register presence result: %T", result)
	}

	if len(values) != 2 {
		return RegisterResult{}, fmt.Errorf("unexpected register presence result length: %d", len(values))
	}

	wasSessionOnline, ok := values[0].(int64)
	if !ok {
		return RegisterResult{}, fmt.Errorf("unexpected session online result: %T", values[0])
	}

	wasUserOnline, ok := values[1].(int64)
	if !ok {
		return RegisterResult{}, fmt.Errorf("unexpected user online result: %T", values[1])
	}

	return RegisterResult{
		WasSessionOnline: int(wasSessionOnline),
		WasUserOnline:    int(wasUserOnline),
	}, nil
}
