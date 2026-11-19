package users_postgres_repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_postgres_pool "github.com/CascadePro/api-golang-server/internal/core/infrastructure/postgres/pool"
	core_utils "github.com/CascadePro/api-golang-server/internal/core/utils"
	"github.com/google/uuid"
)

// Create user, fields to create.
// If activated: "Role", "Email", "Phone" (optional), "PasswordHash" (provides as string)
// If NOT activated: "Role", "Name", "Surname", "*LastName"
func (r *Repository) CreateUser(ctx context.Context, user domain.User) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	var query strings.Builder
	var args []any

	id, err := uuid.NewV7()
	if err != nil {
		return domain.User{}, fmt.Errorf("generate ID: %w", err)
	}

	model := modelFromDomain(user)
	model.ID = id

	if user.Activated {
		query.WriteString(`
			INSERT INTO base.users (id, role, email, phone, password_hash)
			VALUES($1, $2, $3, $4, $5)
		`)

		passwordHash, err := core_utils.GenerateHash(*user.PasswordHash)
		if err != nil {
			return domain.User{}, fmt.Errorf("create user: %w", err)
		}

		args = append(args, model.ID, model.Role, model.Email, model.Phone, passwordHash)
	} else {
		query.WriteString(`
			INSERT INTO base.users (id, role, name, surname, last_name)
			VALUES($1, $2, $3, $4, $5)
		`)

		args = append(args, model.ID, model.Role, model.Name, model.Surname, model.LastName)
	}

	query.WriteString(`
		RETURNING id, version, role, email, phone, name, surname, last_name, last_active_at;
	`)

	row := r.pool.QueryRow(ctx, query.String(), args...)

	var modelResult UserModel

	if err := row.Scan(
		&modelResult.ID,
		&modelResult.Version,
		&modelResult.Role,
		&modelResult.Email,
		&modelResult.Phone,
		&modelResult.Name,
		&modelResult.Surname,
		&modelResult.LastName,
		&modelResult.LastActiveAt,
	); err != nil {
		return domain.User{}, fmt.Errorf("scan error: %w", core_postgres_pool.MapErrors(err))
	}

	return domainFromModel(modelResult), nil
}