package sessions_redis_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	core_redis_pool "github.com/CascadePro/api-golang-server/internal/core/repository/redis/pool"
	core_validation "github.com/CascadePro/api-golang-server/internal/core/validation"
	"github.com/google/uuid"
)

func (r *Repository) GetSession(ctx context.Context, userID uuid.UUID, sessionID string) (domain.Session, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	if err := core_validation.ValidateID(sessionID, domain.SessionIDByteLength); err != nil {
		return domain.Session{}, fmt.Errorf("validate session id: %w", err)
	}

	key := fmt.Sprintf("%s:%s:%s", core_redis_pool.SessionFolder, userID, sessionID)

	var model SessionModel
	if err := r.pool.HGetAll(ctx, key).Scan(&model); err != nil {
		return domain.Session{}, fmt.Errorf("redis hash get all: %w", core_redis_pool.MapErrors(err))
	}

	session, err := modelToDomain(sessionID, model)
	if err != nil {
		err = core_redis_pool.MapErrors(err)

		if errors.Is(err, core_redis_pool.ErrNoValue) {
			return domain.Session{}, fmt.Errorf("session with id='%s': %w", sessionID, core_errors.ErrNotFound)
		}

		return domain.Session{}, fmt.Errorf("session model to domain: %w", err)
	}

	return session, nil
}
