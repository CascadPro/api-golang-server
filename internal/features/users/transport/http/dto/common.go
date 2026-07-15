package users_http_dto

import "github.com/Svat-dev/golang-todo/internal/core/domain"

type UserDtoResponse struct {
	ID          int     `json:"id"`
	Version     int     `json:"version"`
	FullName    string  `json:"full_name"`
	PhoneNumber *string `json:"phone_number"`
}

func UserDtoFromDomain(user domain.User) UserDtoResponse {
	return UserDtoResponse{
		ID:          user.ID,
		Version:     user.Version,
		FullName:    user.FullName,
		PhoneNumber: user.PhoneNumber,
	}
}

func UsersDtoFromDomains(users []domain.User) []UserDtoResponse {
	usersDto := make([]UserDtoResponse, len(users))

	for i, user := range users {
		usersDto[i] = UserDtoFromDomain(user)
	}

	return usersDto
}
