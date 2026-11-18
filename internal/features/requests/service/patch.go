package requests_service

import (
	"context"
	"fmt"
	"time"

	core_context "github.com/CascadePro/api-golang-server/internal/core/context"
	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (s *Service) PatchRequest(ctx context.Context, requestID uuid.UUID, patch domain.RequestPatch) (domain.Request, error) {
	if requestID == uuid.Nil {
		return domain.Request{}, fmt.Errorf("`id` can't be NULL: %w", core_errors.ErrInvalidArgument)
	}

	userID, err := core_context.UserID(ctx)
	if err != nil {
		return domain.Request{}, fmt.Errorf("get user id from context: %w", err)
	}

	request, err := s.requestsMongoRepo.Get(ctx, requestID)
	if err != nil {
		return domain.Request{}, fmt.Errorf("get request from repository: %w", err)
	}

	originSet := patch.Origin.Set

	if err := request.ApplyPatch(patch); err != nil {
		return domain.Request{}, fmt.Errorf("apply patch: %w", err)
	}

	if request.ClientID != nil {
		if _, err := s.clientService.Get(ctx, *request.ClientID); err != nil {
			return domain.Request{}, fmt.Errorf("get client from repository: %w", err)
		}
	}

	if err := s.requestsMongoRepo.Patch(ctx, requestID, request, originSet); err != nil {
		return domain.Request{}, fmt.Errorf("patch request: %w", err)
	}

	_ = s.cache.Delete(ctx, cacheKeyRequest(requestID))

	var clientIDStr *string
	if request.ClientID != nil {
		str := request.ClientID.String()
		clientIDStr = &str
	}

	var deadline *time.Time
	if request.Deadline != nil {
		deadline = request.Deadline
	}

	var status *domain.RequestStatus
	if request.Status != "" {
		status = &request.Status
	}

	data := domain.RealtimeEventRequestUpdatedData{
		ID:        request.ID.String(),
		Title:     &request.Title,
		Status:    status,
		ClientID:  clientIDStr,
		Deadline:  deadline,
		WorkTypes: request.WorkTypes,
		Geography: request.Geography,
		Origin:    request.Origin,
		UpdatedAt: request.UpdatedAt,
	}

	event := domain.NewRealtimeEvent(domain.RealtimeEventRequestUpdated, &userID, nil, data)

	if err := s.realtimePublisher.Publish(ctx, event); err != nil {
		s.logger.Error(
			"failed to publish realtime event",
			zap.String("type", string(event.Type)),
			zap.String("user_id", userID.String()),
			zap.Error(err),
		)
	}

	return request, nil
}
