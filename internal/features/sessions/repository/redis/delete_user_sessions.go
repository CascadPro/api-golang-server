package sessions_redis_repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	core_redis_pool "github.com/CascadePro/api-golang-server/internal/core/infrastructure/redis/pool"
	core_validation "github.com/CascadePro/api-golang-server/internal/core/validation"
	"github.com/google/uuid"
)

func (r *Repository) DeleteUserSessions(ctx context.Context, userID uuid.UUID, sessionID string) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	if userID == uuid.Nil {
		return fmt.Errorf("validate user id: %w", core_errors.ErrInvalidArgument)
	}
	if err := core_validation.ValidateID(sessionID, domain.SessionIDByteLength); err != nil {
		return fmt.Errorf("validate session id: %w", err)
	}

	rootKey := fmt.Sprintf("%s:%s:", core_redis_pool.SessionFolder, userID)

	keys, err := r.pool.GetKeys(ctx, 0, rootKey+"*", 0)
	if err != nil {
		if errors.Is(err, core_redis_pool.ErrNoValue) {
			return fmt.Errorf("keys with user_id='%s': %v: %w", userID, err, core_errors.ErrNotFound)
		}

		return fmt.Errorf("redis get keys: %w", err)
	}

	var idx int = -1
	for i, key := range keys {
		if strings.TrimPrefix(key, rootKey) == sessionID {
			idx = i
			break
		}
	}

	if idx != -1 {
		keys = append(keys[:idx], keys[idx+1:]...)
	}

	if len(keys) == 0 {
		return fmt.Errorf("nothing to delete: %w", core_errors.ErrNotFound)
	}

	if err := r.pool.Del(ctx, keys...); err != nil {
		return fmt.Errorf("redis del session: %w", err)
	}

	return nil
}
