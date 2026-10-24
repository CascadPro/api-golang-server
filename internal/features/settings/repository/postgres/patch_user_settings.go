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

// PatchUserSettings, fields to update: "SessionExpireTerm"
func (r *Repository) PatchUserSettings(
	ctx context.Context,
	id uuid.UUID,
	settings domain.UserSettings,
) (domain.UserSettings, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query, args := patchQuery(id, &settings)

	row := r.pool.QueryRow(ctx, query.String(), args...)

	var model UserSettingsModel
	if err := row.Scan(&model.ID, &model.UserID, &model.Version, &model.SessionExpireTerm); err != nil {
		err := core_postgres_pool.MapErrors(err)

		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.UserSettings{}, fmt.Errorf(
				"user settings with id=%s concurrently accessed: %w",
				id, settings_errors.ErrSettingsNotFound,
			)
		}

		if errors.Is(err, core_postgres_pool.ErrViolatesCheckConstraint) {
			return domain.UserSettings{}, fmt.Errorf(
				"%v: user settings values: %w",
				err, core_errors.ErrInvalidArgument,
			)
		}

		return domain.UserSettings{}, fmt.Errorf("scan error: %w", err)
	}

	return modelToDomain(model), nil
}
