package users_service

import (
	"context"
	"errors"
	"fmt"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	media_errors "github.com/CascadePro/api-golang-server/internal/features/media/errors"
	users_errors "github.com/CascadePro/api-golang-server/internal/features/users/errors"
	"github.com/google/uuid"
)

func (s *Service) GetCurrentUser(ctx context.Context, userID uuid.UUID) (domain.User, []byte, error) {
	if userID == uuid.Nil {
		return domain.User{}, nil, fmt.Errorf("`id` can't be NULL: %w", core_errors.ErrInvalidArgument)
	}

	var cache cacheUser
	err := s.cache.Get(ctx, cacheKeyUser(userID), &cache)
	if err == nil {
		return cache.User, cache.Placeholder, nil
	}
	if !errors.Is(err, core_errors.ErrNotFound) {
		return domain.User{}, nil, fmt.Errorf("get redis cache: %w", err)
	}

	user, err := s.usersPostgresRepo.GetUser(ctx, domain.User{ID: userID})
	if err != nil {
		return domain.User{}, nil, fmt.Errorf("get user from repository: %w", err)
	}
	if !user.Activated {
		return domain.User{}, nil, fmt.Errorf("you must activate your account before usage: %w", users_errors.ErrUserNotActivated)
	}

	var placeholder []byte

	if user.AvatarFileID != nil {
		placeholder, err = s.mediaService.GetFilePlaceholder(ctx, *user.AvatarFileID)
		if err != nil {
			if !errors.Is(err, media_errors.ErrFilePlaceholderNotFound) {
				return domain.User{}, nil, fmt.Errorf("get avatar placeholder: %w", err)
			}
		}
	}

	cache = cacheUser{
		User:        user,
		Placeholder: placeholder,
	}

	if err := s.cache.Set(ctx, cacheKeyUser(userID), cache, cacheTTLUser); err != nil {
		return domain.User{}, nil, fmt.Errorf("set redis cache: %w", err)
	}

	return user, placeholder, nil
}
