package requests_service

import (
	"context"
	"fmt"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
)

func (s *Service) CreateRequest(ctx context.Context, request domain.Request) (domain.Request, error) {
	if err := request.Validate(); err != nil {
		return domain.Request{}, fmt.Errorf("validate request: %w", err)
	}
	if request.ClientID != nil {
		if _, err := s.clientService.GetClient(ctx, *request.ClientID); err != nil {
			return domain.Request{}, fmt.Errorf("get client from repository: %w", err)
		}
	}

	requestID, err := s.requestsMongoRepo.CreateRequest(ctx, request)
	if err != nil {
		return domain.Request{}, fmt.Errorf("create request in repository: %w", err)
	}

	request, err = s.requestsMongoRepo.GetRequest(ctx, requestID)
	if err != nil {
		return domain.Request{}, fmt.Errorf("get request from repository: %w", err)
	}

	return request, nil
}
