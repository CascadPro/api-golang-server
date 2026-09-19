package presence_redis_repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func (r *Repository) Heartbeat(ctx context.Context, userID uuid.UUID, sessionID, connectionID string) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	now := time.Now()

	result, err := r.pool.Eval(
		ctx,
		heartbeatScript,
		[]string{presenceKey(userID)},
		now.UnixMilli(),
		now.Add(LeaseTTL).UnixMilli(),
		presenceMember(sessionID, connectionID),
	)
	if err != nil {
		return fmt.Errorf("heartbeat presence: %w", err)
	}

	value, ok := result.(int64)
	if !ok || value != 1 {
		return fmt.Errorf("presence lease expired")
	}

	return nil
}
