package users_postgres_repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	core_postgres_pool "github.com/CascadePro/api-golang-server/internal/core/infrastructure/postgres/pool"
	"github.com/google/uuid"
)

// Get user: fields: "ID" / "Email"
func (r *Repository) GetUser(ctx context.Context, user domain.User) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	var query strings.Builder
	var arg any
	var errStr string

	query.WriteString(`
		SELECT id, version, activated, role, email, password_hash, name, surname, last_name, last_active_at, avatar_file_id
		FROM base.users
	`)

	if user.ID != uuid.Nil {
		query.WriteString("WHERE (id = $1)")
		arg = user.ID
		errStr = fmt.Sprintf("user with id=%s", user.ID)
	} else if user.Email != nil {
		query.WriteString("WHERE (email = $1)")
		arg = *user.Email
		errStr = fmt.Sprintf("user with email=%s", *user.Email)
	} else {
		return domain.User{}, fmt.Errorf("insufficient fields (need `ID` OR `Email`): %w", core_errors.ErrInvalidArgument)
	}

	query.WriteString("LIMIT 1;")

	row := r.pool.QueryRow(ctx, query.String(), arg)

	var model UserModel
	if err := row.Scan(
		&model.ID,
		&model.Version,
		&model.Activated,
		&model.Role,
		&model.Email,
		&model.PasswordHash,
		&model.Name,
		&model.Surname,
		&model.LastName,
		&model.LastActiveAt,
		&model.AvatarFileID,
	); err != nil {
		err = core_postgres_pool.MapErrors(err)

		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.User{}, fmt.Errorf("%v: %s: %w", err, errStr, core_errors.ErrNotFound)
		}

		return domain.User{}, fmt.Errorf("scan error: %w", err)
	}

	return domainFromModel(model), nil
}
