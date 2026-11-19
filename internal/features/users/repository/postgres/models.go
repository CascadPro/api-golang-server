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

	Email        *string // stored as normalized lowercase email
	Phone        *string // stored as E.164 format (+7xxxxxxxxxx)
	PasswordHash *string
	Role         domain.UserRole

	Name     string
	Surname  string
	LastName *string

	LastActiveAt time.Time
	AvatarFileID *string
}

func domainFromModel(model UserModel) domain.User {
	var email *domain.Email
	if model.Email != nil {
		e := domain.Email(*model.Email)
		email = &e
	}

	var phone *domain.PhoneNumber
	if model.Phone != nil {
		p, err := domain.PhoneNumberFromE164(*model.Phone)
		if err == nil {
			phone = &p
		}
		// If parsing fails, leave phone as nil (data corruption)
	}

	return domain.User{
		ID:           model.ID,
		Version:      model.Version,
		Activated:    model.Activated,
		Email:        email,
		Phone:        phone,
		PasswordHash: model.PasswordHash,
		Role:         model.Role,
		Name:         model.Name,
		Surname:      model.Surname,
		LastName:     model.LastName,
		LastActiveAt: model.LastActiveAt,
		AvatarFileID: model.AvatarFileID,
	}
}

func modelFromDomain(user domain.User) UserModel {
	var email *string
	if user.Email != nil {
		s := user.Email.String()
		email = &s
	}

	var phone *string
	if user.Phone != nil {
		s := user.Phone.String()
		phone = &s
	}

	return UserModel{
		ID:           user.ID,
		Version:      user.Version,
		Activated:    user.Activated,
		Email:        email,
		Phone:        phone,
		PasswordHash: user.PasswordHash,
		Role:         user.Role,
		Name:         user.Name,
		Surname:      user.Surname,
		LastName:     user.LastName,
		LastActiveAt: user.LastActiveAt,
		AvatarFileID: user.AvatarFileID,
	}
}