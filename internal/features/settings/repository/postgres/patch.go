package settings_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	core_postgres_pool "github.com/CascadePro/api-golang-server/internal/core/infrastructure/postgres/pool"
	settings_errors "github.com/CascadePro/api-golang-server/internal/features/settings/errors"
	"github.com/google/uuid"
)

func (r *Repository) PatchUserSettings(
	ctx context.Context,
	id uuid.UUID,
	patch domain.UserSettingsPatch,
) (domain.UserSettings, error) {
	if id == uuid.Nil {
		return domain.UserSettings{}, fmt.Errorf("`id` can't be NULL: %w", core_errors.ErrInvalidArgument)
	}
	if err := patch.Validate(); err != nil {
		return domain.UserSettings{}, fmt.Errorf("validate user settings patch: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query, args := patchQuery(id, &patch)
	row := r.pool.QueryRow(ctx, query.String(), args...)

	var model UserSettingsModel
	if err := row.Scan(
		&model.ID, &model.UserID, &model.Version, &model.SessionExpireTerm,
		&model.Session2FAEnabled, &model.Preferences, &model.Notifications,
		&model.CreatedAt, &model.UpdatedAt,
	); err != nil {
		err = core_postgres_pool.MapErrors(err)
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			var existingID uuid.UUID

			existsErr := r.pool.QueryRow(ctx, settingsExistsQuery, id).Scan(&existingID)
			existsErr = core_postgres_pool.MapErrors(existsErr)

			if errors.Is(existsErr, core_postgres_pool.ErrNoRows) {
				return domain.UserSettings{}, fmt.Errorf("user settings with id=%s: %w", id, settings_errors.ErrSettingsNotFound)
			}

			if existsErr != nil {
				return domain.UserSettings{}, fmt.Errorf("check user settings after patch conflict: %w", existsErr)
			}

			return domain.UserSettings{}, fmt.Errorf("user settings version is stale: %w", core_errors.ErrConflict)
		}

		if errors.Is(err, core_postgres_pool.ErrViolatesCheckConstraint) {
			return domain.UserSettings{}, fmt.Errorf("%v: user settings values: %w", err, core_errors.ErrInvalidArgument)
		}

		return domain.UserSettings{}, fmt.Errorf("scan patched user settings: %w", err)
	}

	result, err := modelToDomain(model)
	if err != nil {
		return domain.UserSettings{}, fmt.Errorf("map patched user settings: %w", err)
	}

	return result, nil
}
