package sessions_redis_repository

import (
	"context"
	"fmt"
	"time"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_redis_pool "github.com/CascadePro/api-golang-server/internal/core/repository/redis/pool"
	core_utils "github.com/CascadePro/api-golang-server/internal/core/utils"
	"github.com/google/uuid"
)

func (r *Repository) CreateSession(ctx context.Context, userID uuid.UUID, session domain.Session) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	sessionID, err := core_utils.GenerateID(domain.SessionIDByteLength)
	if err != nil {
		return "", fmt.Errorf("generate id: %w", err)
	}

	key := fmt.Sprintf("%s:%s:%s", core_redis_pool.SessionFolder, userID, sessionID)

	model, err := domainToModel(session)
	if err != nil {
		return "", fmt.Errorf("session domain to model: %w", err)
	}

	pipe := r.pool.Pipeline()

	pipe.HSet(ctx, key, "created_at", model.CreatedAt)
	pipe.HSet(ctx, key, "expires_at", model.ExpiresAt)
	pipe.HSet(ctx, key, "last_active_at", model.LastActiveAt)
	pipe.HSet(ctx, key, "ip", model.IP.String())
	pipe.HSet(ctx, key, "metadata", model.Metadata)

	pipe.Expire(ctx, key, time.Duration(session.ExpirationTime))

	if _, err := pipe.Exec(ctx); err != nil {
		return "", fmt.Errorf("redis pipe exec: %w", err)
	}

	return sessionID, nil
}
