package sessions_redis_repository

import (
	"context"
	"fmt"
	"time"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	core_redis_pool "github.com/CascadePro/api-golang-server/internal/core/infrastructure/redis/pool"
	core_validation "github.com/CascadePro/api-golang-server/internal/core/validation"
	"github.com/google/uuid"
)

func (r *Repository) PatchLastActive(ctx context.Context, userID uuid.UUID, sessionID string) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	if userID == uuid.Nil {
		return fmt.Errorf("validate user id: %w", core_errors.ErrInvalidArgument)
	}
	if err := core_validation.ValidateID(sessionID, domain.SessionIDByteLength); err != nil {
		return fmt.Errorf("validate session id: %w", err)
	}

	key := fmt.Sprintf("%s:%s:%s", core_redis_pool.SessionFolder, userID, sessionID)
	now := time.Now()

	if err := r.pool.HSet(ctx, key, string(HashFieldLastActive), now); err != nil {
		return fmt.Errorf("patch last active: %w", err)
	}

	return nil
}
