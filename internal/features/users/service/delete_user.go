package users_service

import (
	"context"
	"fmt"
)

func (s *UsersService) DeleteUser(ctx context.Context, id int) error {
	err := s.usersRepository.DeleteUser(ctx, id)
	if err != nil {
		return fmt.Errorf("delete user from repository: %w", err)
	}

	return nil
}
