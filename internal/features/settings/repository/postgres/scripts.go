package settings_postgres_repository

import (
	"fmt"
	"strings"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	"github.com/google/uuid"
)

func getQuery(id, userID uuid.UUID) (query strings.Builder, arg uuid.UUID, errStr string, _ error) {
	query.WriteString(`
		SELECT id, user_id, version, session_expire_term
		FROM base.user_settings
	`)

	if id != uuid.Nil {
		query.WriteString("WHERE (id = $1)")
		arg = id
		errStr = fmt.Sprintf("user settings with id=%s", id)
	} else if userID != uuid.Nil {
		query.WriteString("WHERE (user_id = $1)")
		arg = userID
		errStr = fmt.Sprintf("user settings with user_id=%s", userID)
	} else {
		return query, arg, errStr, fmt.Errorf("insufficient fields (need `ID` OR `UserID`): %w", core_errors.ErrInvalidArgument)
	}

	query.WriteString("LIMIT 1;")

	return query, arg, errStr, nil
}

const createQuery = `
	INSERT INTO base.user_settings (id, user_id, session_expire_term)
	VALUES ($1, $2, $3)
	RETURNING id, user_id, version, session_expire_term;
`

func patchQuery(id uuid.UUID, settings *domain.UserSettings) (query strings.Builder, args []any) {
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

	return query, args
}
