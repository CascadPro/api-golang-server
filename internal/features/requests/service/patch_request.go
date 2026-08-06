package requests_service

import (
	"context"
	"fmt"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	"github.com/google/uuid"
)

func (s *Service) PatchRequest(ctx context.Context, id uuid.UUID, patch domain.RequestPatch) (domain.Request, error) {
	if id == uuid.Nil {
		return domain.Request{}, fmt.Errorf("`id` can't be NULL: %w", core_errors.ErrInvalidArgument)
	}

	request, err := s.requestsMongoRepo.GetRequest(ctx, id)
	if err != nil {
		return domain.Request{}, fmt.Errorf("get request from repository: %w", err)
	}

	if err := request.ApplyPatch(patch); err != nil {
		return domain.Request{}, fmt.Errorf("apply patch: %w", err)
	}

	if err := s.requestsMongoRepo.PatchRequest(ctx, id, request); err != nil {
		return domain.Request{}, fmt.Errorf("patch request: %w", err)
	}

	return request, nil
}
