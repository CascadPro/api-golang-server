package settings_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_postgres_pool "github.com/CascadePro/api-golang-server/internal/core/infrastructure/postgres/pool"
	settings_errors "github.com/CascadePro/api-golang-server/internal/features/settings/errors"
)

// Get user settings, fields to search: "ID" / "UserID"
func (r *Repository) GetUserSettings(ctx context.Context, settings domain.UserSettings) (domain.UserSettings, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query, arg, errStr, err := getQuery(settings.ID, settings.UserID)
	if err != nil {
		return domain.UserSettings{}, err
	}

	row := r.pool.QueryRow(ctx, query.String(), arg)

	var model UserSettingsModel
	if err := row.Scan(&model.ID, &model.UserID, &model.Version, &model.SessionExpireTerm); err != nil {
		err = core_postgres_pool.MapErrors(err)

		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.UserSettings{}, fmt.Errorf("%v: %s: %w", err, errStr, settings_errors.ErrSettingsNotFound)
		}

		return domain.UserSettings{}, fmt.Errorf("scan error: %w", err)
	}

	return modelToDomain(model), nil
}
