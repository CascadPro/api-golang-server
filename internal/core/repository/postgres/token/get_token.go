package core_postgres_token

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

// Get token, fields: "ID" / "UserID" + "Type" / "UserID" + "Token" / "Token" + "Type"
func (r *Repository) GetToken(ctx context.Context, token domain.Token) (domain.Token, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	var query strings.Builder
	var args []any
	var errStr string

	query.WriteString(`
		SELECT id, version, token, type, expires_at, user_id
		FROM base.tokens
	`)

	if token.ID != uuid.Nil {
		query.WriteString("WHERE (id = $1)")
		args = append(args, token.ID)
		errStr = fmt.Sprintf("token with id=%s", token.ID)
	} else if token.UserID != uuid.Nil && token.Type != domain.TokenTypeNil {
		query.WriteString("WHERE (user_id = $1 AND type = $2)")
		args = append(args, token.UserID, token.Type)
		errStr = fmt.Sprintf("token with user_id=%s AND type=%s", token.UserID, token.Type)
	} else if token.UserID != uuid.Nil && token.Token != "" {
		query.WriteString("WHERE (user_id = $1 AND token = $2)")
		args = append(args, token.UserID, token.Token)
		errStr = fmt.Sprintf("token with user_id=%s AND token=%s", token.UserID, token.Token)
	} else if token.Token != "" && token.Type != domain.TokenTypeNil {
		query.WriteString("WHERE (token = $1 AND type = $2)")
		args = append(args, token.Token, token.Type)
		errStr = fmt.Sprintf("token with token=%s AND type=%s", token.Token, token.Type)
	} else {
		return domain.Token{}, fmt.Errorf(
			"insufficient fields (need `ID` OR `UserID`+`Type` OR `UserID`+`Token` OR `Token`+`Type`): %w",
			core_errors.ErrInvalidArgument,
		)
	}

	query.WriteString("LIMIT 1;")

	row := r.pool.QueryRow(ctx, query.String(), args...)

	var model TokenModel
	if err := row.Scan(
		&model.ID,
		&model.Version,
		&model.Token,
		&model.Type,
		&model.ExpiresAt,
		&model.UserID,
	); err != nil {
		err = core_postgres_pool.MapErrors(err)

		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.Token{}, fmt.Errorf("%v: %s: %w", err, errStr, core_errors.ErrNotFound)
		}

		return domain.Token{}, fmt.Errorf("scan error: %w", err)
	}

	return domainFromModel(model), nil
}
