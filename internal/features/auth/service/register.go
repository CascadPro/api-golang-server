package auth_service

import (
	"context"
	"fmt"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	"github.com/google/uuid"
)

func (s *Service) Register(ctx context.Context, patch domain.UserPatch, tokenString string) error {
	if _, err := uuid.Parse(tokenString); err != nil {
		return fmt.Errorf("validate token: %w", err)
	}

	tokenDto := domain.NewGetToken(tokenString, domain.TokenTypeRegister, uuid.Nil)

	tokenDomain, err := s.tokenPostgresRepo.GetToken(ctx, tokenDto)
	if err != nil {
		return fmt.Errorf("get token: %w", err)
	}

	userDomain, err := s.userPostgresRepo.GetUser(ctx, domain.User{ID: tokenDomain.UserID})
	if err != nil {
		return fmt.Errorf("get user: %w", err)
	}

	patchedUser, err := userDomain.ApplyPatch(patch)
	if err != nil {
		return fmt.Errorf("apply user patch: %w", err)
	}

	if _, err = s.userPostgresRepo.PatchUser(ctx, tokenDomain.UserID, patchedUser); err != nil {
		return fmt.Errorf("patch user: %w", err)
	}

	if err := s.tokenPostgresRepo.DeleteToken(ctx, tokenDomain.ID); err != nil {
		return fmt.Errorf("delete token: %w", err)
	}

	return nil
}
