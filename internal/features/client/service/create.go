package client_service

import (
	"context"
	"fmt"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
)

func (s *Service) CreateClient(ctx context.Context, client domain.Client) (domain.Client, error) {
	if err := client.Validate(); err != nil {
		return domain.Client{}, fmt.Errorf("validate client: %w", err)
	}

	client, err := s.clientPostgresRepo.CreateClient(ctx, client)
	if err != nil {
		return domain.Client{}, fmt.Errorf("create client: %w", err)
	}

	return client, nil
}
