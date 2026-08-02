package sessions_redis_repository

import (
	"context"
	"fmt"
	"time"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	core_redis_pool "github.com/CascadePro/api-golang-server/internal/core/repository/redis/pool"
	"github.com/google/uuid"
)

func (r *Repository) PatchSession(ctx context.Context, userID uuid.UUID, session domain.Session) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	if userID == uuid.Nil {
		return fmt.Errorf("validate user id: %w", core_errors.ErrInvalidArgument)
	}
	if err := session.Validate(); err != nil {
		return fmt.Errorf("validate session: %w", err)
	}

	key := fmt.Sprintf("%s:%s:%s", core_redis_pool.SessionFolder, userID, session.ID)

	model, err := domainToModel(session)
	if err != nil {
		return fmt.Errorf("session domain to model: %w", err)
	}

	pipe := r.pool.TxPipeline()

	pipe.HSet(ctx, key, HashFieldIP, model.IP.String())
	pipe.HSet(ctx, key, HashFieldExpiresAt, model.ExpiresAt)
	pipe.HSet(ctx, key, HashFieldLastActive, model.LastActiveAt)
	pipe.HSet(ctx, key, HashFieldMetadata, model.Metadata)

	pipe.Expire(ctx, key, time.Duration(session.ExpirationTime))

	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("redis pipe exec: %w", err)
	}

	return nil
}
