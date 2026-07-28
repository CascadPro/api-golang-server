package core_postgres_user

import (
	"time"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	"github.com/google/uuid"
)

type UserModel struct {
	ID        uuid.UUID
	Version   int
	Activated bool

	Email        *string
	PasswordHash *string
	Role         domain.UserRole

	Name     string
	Surname  string
	LastName *string

	LastActiveAt time.Time
}
