package users_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	core_postgres_pool "github.com/CascadePro/api-golang-server/internal/core/repository/postgres/pool"
	core_utils "github.com/CascadePro/api-golang-server/internal/core/utils"
	"github.com/google/uuid"
)

// PatchUser, fields to update: "Activated", "Email", "PasswordHash" (provides as string), "Role", "Name", "Surname",
// "LastName", "AvatarFileID", "LastActiveAt"
func (r *Repository) PatchUser(ctx context.Context, id uuid.UUID, user domain.User) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		UPDATE base.users
		SET version = version + 1, activated = $1, email = $2, password_hash = $3, role = $4,
			name = $5, surname = $6, last_name = $7, last_active_at = $8
		WHERE (id = $9 AND version = $10)
		RETURNING id, version, activated, role, email, name, surname, last_name, last_active_at;
	`

	passwordHash, err := core_utils.GenerateHash(*user.PasswordHash)
	if err != nil {
		return domain.User{}, fmt.Errorf("patch user: %w", err)
	}

	row := r.pool.QueryRow(ctx, query,
		user.Activated,
		user.Email,
		passwordHash,
		user.Role,
		user.Name,
		user.Surname,
		user.LastName,
		user.LastActiveAt,
		id,
		user.Version,
	)

	var model UserModel
	if err := row.Scan(
		&model.ID,
		&model.Version,
		&model.Activated,
		&model.Role,
		&model.Email,
		&model.Name,
		&model.Surname,
		&model.LastName,
		&model.LastActiveAt,
	); err != nil {
		err = core_postgres_pool.MapErrors(err)

		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.User{}, fmt.Errorf(
				"user with id=%s concurrently accessed: %w",
				id, core_errors.ErrConflict,
			)
		}

		if errors.Is(err, core_postgres_pool.ErrViolatesUniqueConstraint) {
			return domain.User{}, fmt.Errorf(
				"`Email` is already occupied by another user: %w",
				core_errors.ErrInvalidArgument,
			)
		}

		return domain.User{}, fmt.Errorf("scan error: %w", err)
	}

	return domainFromModel(model), nil
}
