package users_postgres_repository

import (
	"time"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	"github.com/google/uuid"
)

type UserModel struct {
	ID        uuid.UUID
	Version   int64
	Activated bool

	Email        *string
	PasswordHash *string
	Role         domain.UserRole

	Name     string
	Surname  string
	LastName *string

	LastActiveAt time.Time
	AvatarFileID *string
}

func domainFromModel(model UserModel) domain.User {
	return domain.User{
		ID:           model.ID,
		Version:      model.Version,
		Activated:    model.Activated,
		Email:        model.Email,
		PasswordHash: model.PasswordHash,
		Role:         model.Role,
		Name:         model.Name,
		Surname:      model.Surname,
		LastName:     model.LastName,
		LastActiveAt: model.LastActiveAt,
		AvatarFileID: model.AvatarFileID,
	}
}
