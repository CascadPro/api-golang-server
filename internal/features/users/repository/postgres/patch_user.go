package users_postgres_repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	core_postgres_pool "github.com/CascadePro/api-golang-server/internal/core/infrastructure/postgres/pool"
	core_utils "github.com/CascadePro/api-golang-server/internal/core/utils"
	"github.com/google/uuid"
)

// PatchUser, fields to update: "Activated", "Email", "PasswordHash" (provides as string), "Role", "Name", "Surname",
// "LastName", "AvatarFileID", "LastActiveAt"
func (r *Repository) PatchUser(ctx context.Context, id uuid.UUID, user domain.User) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	var query strings.Builder
	var args []any

	query.WriteString(`
		UPDATE base.users
		SET version = version + 1, activated = $1
	`)
	args = append(args, user.Activated)

	if user.Email != nil {
		fmt.Fprintf(&query, ", email = $%d", len(args)+1)
		args = append(args, user.Email)
	}

	if user.PasswordHash != nil {
		passwordHash, err := core_utils.GenerateHash(*user.PasswordHash)
		if err != nil {
			return domain.User{}, fmt.Errorf("patch user: %w", err)
		}

		fmt.Fprintf(&query, ", password_hash = $%d", len(args)+1)
		args = append(args, passwordHash)
	}

	if user.Role != "" {
		fmt.Fprintf(&query, ", role = $%d", len(args)+1)
		args = append(args, user.Role)
	}

	if user.Name != "" {
		fmt.Fprintf(&query, ", name = $%d", len(args)+1)
		args = append(args, user.Name)
	}

	if user.Surname != "" {
		fmt.Fprintf(&query, ", surname = $%d", len(args)+1)
		args = append(args, user.Surname)
	}

	if user.LastName != nil {
		fmt.Fprintf(&query, ", last_name = $%d", len(args)+1)
		args = append(args, user.LastName)
	}

	if user.AvatarFileID != nil {
		fmt.Fprintf(&query, ", avatar_file_id = $%d", len(args)+1)
		args = append(args, user.AvatarFileID)
	}

	if user.LastActiveAt != (time.Time{}) {
		fmt.Fprintf(&query, ", last_active_at = $%d", len(args)+1)
		args = append(args, user.LastActiveAt)
	}

	fmt.Fprintf(&query, " WHERE (id = $%d AND version = $%d)", len(args)+1, len(args)+2)
	args = append(args, id, user.Version)

	query.WriteString(`
		RETURNING id, version, activated, role, email, name, surname, last_name, avatar_file_id, last_active_at;
	`)

	row := r.pool.QueryRow(ctx, query.String(), args...)

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
		&model.AvatarFileID,
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
