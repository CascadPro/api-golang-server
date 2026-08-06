package requests_service

import (
	"context"
	"fmt"

	core_context "github.com/CascadePro/api-golang-server/internal/core/context"
	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	"github.com/google/uuid"
)

func (s *Service) PatchRequestStatus(ctx context.Context, id uuid.UUID, status domain.RequestStatus) error {
	if id == uuid.Nil {
		return fmt.Errorf("`ID` can't be NULL: %w", core_errors.ErrInvalidArgument)
	}

	userID, err := core_context.UserIDFromContext(ctx)
	if err != nil {
		return fmt.Errorf("get user id from context: %w", err)
	}

	if err := s.requestsMongoRepo.PatchRequestStatus(ctx, id, userID, status); err != nil {
		return fmt.Errorf("patch request status: %w", err)
	}

	return nil
}
