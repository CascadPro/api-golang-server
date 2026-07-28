package sessions_redis_repository

import (
	"context"
	"errors"
	"fmt"

	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	core_redis_pool "github.com/CascadePro/api-golang-server/internal/core/repository/redis/pool"
	"github.com/google/uuid"
)

func (r *Repository) DeleteSession(ctx context.Context, userID uuid.UUID, sessionID string) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	key := fmt.Sprintf("%s:%s:%s", core_redis_pool.SessionFolder, userID, sessionID)

	if err := r.pool.Del(ctx, key); err != nil {
		if errors.Is(err, core_redis_pool.ErrNoValue) {
			return fmt.Errorf("session with id='%s': %w", sessionID, core_errors.ErrNotFound)
		}

		return fmt.Errorf("del session: %w", err)
	}

	return nil
}
