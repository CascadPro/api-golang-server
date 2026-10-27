package settings_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	core_postgres_pool "github.com/CascadePro/api-golang-server/internal/core/infrastructure/postgres/pool"
	users_errors "github.com/CascadePro/api-golang-server/internal/features/users/errors"
	"github.com/google/uuid"
)

func (r *Repository) CreateUserSettings(
	ctx context.Context,
	settings domain.UserSettings,
) (domain.UserSettings, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	id, err := uuid.NewV7()
	if err != nil {
		return domain.UserSettings{}, fmt.Errorf("generate ID: %w", err)
	}

	row := r.pool.QueryRow(ctx, createQuery, id, settings.UserID, settings.SessionExpireTerm)
	var model UserSettingsModel
	if err := row.Scan(
		&model.ID, &model.UserID, &model.Version, &model.SessionExpireTerm,
		&model.Session2FAEnabled, &model.Preferences, &model.Notifications,
		&model.CreatedAt, &model.UpdatedAt,
	); err != nil {
		err = core_postgres_pool.MapErrors(err)
		if errors.Is(err, core_postgres_pool.ErrViolatesForeignKey) {
			return domain.UserSettings{}, fmt.Errorf("%v: user with id=%s: %w", err, settings.UserID, users_errors.ErrUserNotFound)
		}
		if errors.Is(err, core_postgres_pool.ErrViolatesCheckConstraint) {
			return domain.UserSettings{}, fmt.Errorf("%v: user settings values: %w", err, core_errors.ErrInvalidArgument)
		}
		return domain.UserSettings{}, err
	}

	result, err := modelToDomain(model)
	if err != nil {
		return domain.UserSettings{}, fmt.Errorf("map created user settings: %w", err)
	}

	return result, nil
}
