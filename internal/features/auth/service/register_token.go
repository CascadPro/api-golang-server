package auth_service

import (
	"context"
	"fmt"
	"time"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	"github.com/google/uuid"
)

func (s *Service) CreateRegisterToken(ctx context.Context, user domain.User) (token domain.Token, err error) {
	if err := user.Validate(); err != nil {
		return domain.Token{}, fmt.Errorf("validate user domain: %w", err)
	}

	userDomain, err := s.userPostgresRepo.CreateUser(ctx, user)
	if err != nil {
		return domain.Token{}, fmt.Errorf("create user: %w", err)
	}

	tokenDto := domain.NewToken(
		uuid.NewString(),
		domain.TokenTypeRegister,
		time.Now().Add(time.Hour*12),
		userDomain.ID,
	)

	tokenDomain, err := s.tokenPostgresRepo.CreateToken(ctx, tokenDto)
	if err != nil {
		return domain.Token{}, fmt.Errorf("create token: %w", err)
	}

	return tokenDomain, nil
}
