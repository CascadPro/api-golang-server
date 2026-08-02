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

// Get user settings, fields to search: "ID" / "UserID"
func (r *Repository) GetUserSettings(ctx context.Context, settings domain.UserSettings) (domain.UserSettings, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	var query strings.Builder
	var arg uuid.UUID
	var errStr string

	query.WriteString(`
		SELECT id, user_id, version, session_expire_term
		FROM base.user_settings
	`)

	if settings.ID != uuid.Nil {
		query.WriteString("WHERE (id = $1)")
		arg = settings.ID
		errStr = fmt.Sprintf("user settings with id=%s", settings.ID)
	} else if settings.UserID != uuid.Nil {
		query.WriteString("WHERE (user_id = $1)")
		arg = settings.UserID
		errStr = fmt.Sprintf("user settings with user_id=%s", settings.UserID)
	} else {
		return domain.UserSettings{},
			fmt.Errorf("insufficient fields (need `ID` OR `UserID`): %w", core_errors.ErrInvalidArgument)
	}

	query.WriteString("LIMIT 1;")

	row := r.pool.QueryRow(ctx, query.String(), arg)

	var model UserSettingsModel
	if err := row.Scan(&model.ID, &model.UserID, &model.Version, &model.SessionExpireTerm); err != nil {
		err = core_postgres_pool.MapErrors(err)

		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.UserSettings{}, fmt.Errorf("%v: %s: %w", err, errStr, core_errors.ErrNotFound)
		}

		return domain.UserSettings{}, fmt.Errorf("scan error: %w", err)
	}

	return modelToDomain(model), nil
}
