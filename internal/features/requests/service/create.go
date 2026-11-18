package requests_service

import (
	"context"
	"fmt"

	core_context "github.com/CascadePro/api-golang-server/internal/core/context"
	"github.com/CascadePro/api-golang-server/internal/core/domain"
	"go.uber.org/zap"
)

func (s *Service) CreateRequest(ctx context.Context, request domain.Request) (domain.Request, error) {
	if err := request.Validate(); err != nil {
		return domain.Request{}, fmt.Errorf("validate request: %w", err)
	}
	if request.ClientID != nil {
		if _, err := s.clientService.Get(ctx, *request.ClientID); err != nil {
			return domain.Request{}, fmt.Errorf("get client from repository: %w", err)
		}
	}

	userID, err := core_context.UserID(ctx)
	if err != nil {
		return domain.Request{}, fmt.Errorf("get user id from context: %w", err)
	}

	requestID, err := s.requestsMongoRepo.Create(ctx, request)
	if err != nil {
		return domain.Request{}, fmt.Errorf("create request in repository: %w", err)
	}

	_ = s.cache.Delete(ctx, cacheKeyStats)

	request, err = s.requestsMongoRepo.Get(ctx, requestID)
	if err != nil {
		return domain.Request{}, fmt.Errorf("get request from repository: %w", err)
	}

	var clientIDStr *string
	if request.ClientID != nil {
		str := request.ClientID.String()
		clientIDStr = &str
	}

	data := domain.NewRealtimeEventRequestCreatedData(
		request.ID.String(),
		request.Title,
		request.Status,
		clientIDStr,
		request.CreatedAt,
	)

	event := domain.NewRealtimeEvent(domain.RealtimeEventRequestCreated, &userID, nil, data)

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
