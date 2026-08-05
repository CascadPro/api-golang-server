package requests_service

import (
	"context"
	"fmt"

	core_context "github.com/CascadePro/api-golang-server/internal/core/context"
	"github.com/CascadePro/api-golang-server/internal/core/domain"
	"github.com/google/uuid"
)

func (s *Service) ApproveRequest(ctx context.Context, id uuid.UUID) error {
	userID, err := core_context.UserIDFromContext(ctx)
	if err != nil {
		return fmt.Errorf("get user id from context: %w", err)
	}

	request, err := s.requestsMongoRepo.GetRequest(ctx, id)
	if err != nil {
		return fmt.Errorf("get request from repository: %w", err)
	}

	status := domain.RequestStatusApproved
	request.Status = status

	if err = request.Validate(); err != nil {
		status = domain.RequestStatusWaiting
	}

	patch := domain.NewStatusPatchRequest(request.Version, status, userID)
	if err := s.requestsMongoRepo.PatchRequest(ctx, request.ID, patch); err != nil {
		return fmt.Errorf("patch request to repository: %w", err)
	}

	return err
}
