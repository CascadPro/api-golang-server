package settings_postgres_repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	core_postgres_pool "github.com/CascadePro/api-golang-server/internal/core/repository/postgres/pool"
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

	var query strings.Builder
	var args []any

	query.WriteString(`
		UPDATE base.user_settings
		SET version = version + 1
	`)

	if settings.SessionExpireTerm != domain.SessionExpireNil {
		fmt.Fprintf(&query, ", session_expire_term = $%d", len(args)+1)
		args = append(args, settings.SessionExpireTerm)
	}

	fmt.Fprintf(&query, " WHERE (id = $%d AND version = $%d)", len(args)+1, len(args)+2)
	args = append(args, id, settings.Version)

	query.WriteString("RETURNING id, user_id, version, session_expire_term;")

	row := r.pool.QueryRow(ctx, query.String(), args...)

	var model UserSettingsModel
	if err := row.Scan(&model.ID, &model.UserID, &model.Version, &model.SessionExpireTerm); err != nil {
		err := core_postgres_pool.MapErrors(err)

		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.UserSettings{}, fmt.Errorf(
				"user settings with id=%s concurrently accessed: %w",
				id, core_errors.ErrConflict,
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
