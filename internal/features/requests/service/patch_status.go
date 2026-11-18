package requests_service

import (
	"context"
	"fmt"

	core_context "github.com/CascadePro/api-golang-server/internal/core/context"
	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (s *Service) PatchRequestStatus(ctx context.Context, requestID uuid.UUID, status domain.RequestStatus) error {
	if requestID == uuid.Nil {
		return fmt.Errorf("`ID` can't be NULL: %w", core_errors.ErrInvalidArgument)
	}

	userID, err := core_context.UserID(ctx)
	if err != nil {
		return fmt.Errorf("get user id from context: %w", err)
	}

	if err := s.requestsMongoRepo.PatchStatus(ctx, requestID, userID, status); err != nil {
		return fmt.Errorf("patch request status: %w", err)
	}

	_ = s.cache.Delete(ctx, cacheKeyStats)
	_ = s.cache.Delete(ctx, cacheKeyRequest(requestID))

	request, err := s.requestsMongoRepo.Get(ctx, requestID)
	if err != nil {
		return fmt.Errorf("get request from repository: %w", err)
	}

	data := domain.NewRealtimeEventRequestStatusChangedData(
		request.ID.String(),
		request.Status,
		userID.String(),
		request.UpdatedAt,
	)

	event := domain.NewRealtimeEvent(domain.RealtimeEventRequestStatusChanged, &userID, nil, data)

	if err := s.realtimePublisher.Publish(ctx, event); err != nil {
		s.logger.Error(
			"failed to publish realtime event",
			zap.String("type", string(event.Type)),
			zap.String("user_id", userID.String()),
			zap.Error(err),
		)
	}

	return nil
}
