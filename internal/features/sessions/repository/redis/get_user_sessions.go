package sessions_redis_repository

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	core_logger "github.com/CascadePro/api-golang-server/internal/core/logger"
	core_redis_pool "github.com/CascadePro/api-golang-server/internal/core/repository/redis/pool"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func (r *Repository) GetUserSessions(ctx context.Context, userID uuid.UUID) ([]domain.Session, error) {
	log := core_logger.FromContext(ctx)
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	key := fmt.Sprintf("%s:%s:*", core_redis_pool.SessionFolder, userID)

	keys, err := r.pool.GetKeys(ctx, 0, key, 0)
	if err != nil {
		if errors.Is(err, core_redis_pool.ErrNoValue) {
			return nil, fmt.Errorf("keys with user_id='%s': %v: %w", userID, err, core_errors.ErrNotFound)
		}

		return nil, fmt.Errorf("redis get keys: %w", err)
	}

	pipe := r.pool.Pipeline()

	var cmds []*redis.MapStringStringCmd
	for _, key := range keys {
		cmds = append(cmds, pipe.HGetAll(ctx, key))
	}

	if _, err := pipe.Exec(ctx); err != nil {
		return nil, fmt.Errorf("redis pipe exec: %w", core_redis_pool.MapErrors(err))
	}

	var sessions []domain.Session
	for i, cmd := range cmds {
		var model SessionModel
		if err := cmd.Scan(&model); err != nil {
			return nil, fmt.Errorf("redis hash get all: at index %d: %w", i, core_redis_pool.MapErrors(err))
		}

		splittedKey := strings.Split(keys[i], ":")
		sessionID := splittedKey[len(splittedKey)-1]

		session, err := modelToDomain(sessionID, model)
		if err != nil {
			log.Error("user sessions scan failed", zap.Error(fmt.Errorf("redis scan error: iter %d: %w", i, err)))
			continue
		}

		sessions = append(sessions, session)
	}
	if len(sessions) == 0 {
		return nil, fmt.Errorf("sessions with user_id='%s': %w", userID, core_errors.ErrNotFound)
	}

	sort.Slice(sessions, func(i, j int) bool {
		return sessions[i].CreatedAt.After(sessions[j].CreatedAt)
	})

	return sessions, nil
}
