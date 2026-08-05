package requests_service

import (
	"context"
	"fmt"

	core_context "github.com/CascadePro/api-golang-server/internal/core/context"
	"github.com/CascadePro/api-golang-server/internal/core/domain"
	"github.com/google/uuid"
)

func (s *Service) RejectRequest(ctx context.Context, id uuid.UUID) error {
	userID, err := core_context.UserIDFromContext(ctx)
	if err != nil {
		return fmt.Errorf("get user id from context: %w", err)
	}

	request, err := s.requestsMongoRepo.GetRequest(ctx, id)
	if err != nil {
		return fmt.Errorf("get request from repository: %w", err)
	}

	patchedRequest := domain.NewStatusPatchRequest(request.Version, domain.RequestStatusRejected, userID)

	if err := s.requestsMongoRepo.PatchRequest(ctx, id, patchedRequest); err != nil {
		return fmt.Errorf("patch request to repository: %w", err)
	}

	return nil
}
