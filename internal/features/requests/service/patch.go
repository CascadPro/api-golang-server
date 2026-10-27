package requests_service

import (
	"context"
	"fmt"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	"github.com/google/uuid"
)

func (s *Service) PatchRequest(ctx context.Context, requestID uuid.UUID, patch domain.RequestPatch) (domain.Request, error) {
	if requestID == uuid.Nil {
		return domain.Request{}, fmt.Errorf("`id` can't be NULL: %w", core_errors.ErrInvalidArgument)
	}

	request, err := s.requestsMongoRepo.GetRequest(ctx, requestID)
	if err != nil {
		return domain.Request{}, fmt.Errorf("get request from repository: %w", err)
	}

	if err := request.ApplyPatch(patch); err != nil {
		return domain.Request{}, fmt.Errorf("apply patch: %w", err)
	}

	if request.ClientID != nil {
		if _, err := s.clientService.GetClient(ctx, *request.ClientID); err != nil {
			return domain.Request{}, fmt.Errorf("get client from repository: %w", err)
		}
	}

	if err := s.requestsMongoRepo.PatchRequest(ctx, requestID, request); err != nil {
		return domain.Request{}, fmt.Errorf("patch request: %w", err)
	}

	_ = s.cache.Delete(ctx, cacheKeyRequest(requestID))

	return request, nil
}
