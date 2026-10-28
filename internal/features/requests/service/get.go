package requests_service

import (
	"context"
	"errors"
	"fmt"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	"github.com/google/uuid"
)

func (s *Service) GetRequest(
	ctx context.Context,
	requestID uuid.UUID,
) (domain.Request, domain.User, domain.Client, map[int]domain.File, error) {
	if requestID == uuid.Nil {
		return domain.Request{}, domain.User{}, domain.Client{}, nil, fmt.Errorf("`id` can't be NULL: %w", core_errors.ErrInvalidArgument)
	}

	var cache cacheRequest
	err := s.cache.Get(ctx, cacheKeyRequest(requestID), &cache)
	if err == nil {
		return cache.Request, cache.User, cache.Client, cache.Files, nil
	}
	if !errors.Is(err, core_errors.ErrNotFound) {
		return domain.Request{}, domain.User{}, domain.Client{}, nil, fmt.Errorf("get redis cache: %w", err)
	}

	request, err := s.requestsMongoRepo.Get(ctx, requestID)
	if err != nil {
		return domain.Request{}, domain.User{}, domain.Client{}, nil, fmt.Errorf("get request from repository: %w", err)
	}

	filesMap := make(map[int]domain.File, 4)

	if request.ProjectDocID != nil {
		file, err := s.mediaService.GetFile(ctx, *request.ProjectDocID)
		if err != nil {
			return domain.Request{}, domain.User{}, domain.Client{}, nil, fmt.Errorf("get `ProjectPlan` file: %w", err)
		}

		filesMap[0] = file
	}
	if request.TechTaskDocID != nil {
		file, err := s.mediaService.GetFile(ctx, *request.TechTaskDocID)
		if err != nil {
			return domain.Request{}, domain.User{}, domain.Client{}, nil, fmt.Errorf("get `TechTask` file: %w", err)
		}

		filesMap[1] = file
	}
	if request.SpecificationDocID != nil {
		file, err := s.mediaService.GetFile(ctx, *request.SpecificationDocID)
		if err != nil {
			return domain.Request{}, domain.User{}, domain.Client{}, nil, fmt.Errorf("get `Specification` file: %w", err)
		}

		filesMap[2] = file
	}
	if request.ContractDocID != nil {
		file, err := s.mediaService.GetFile(ctx, *request.ContractDocID)
		if err != nil {
			return domain.Request{}, domain.User{}, domain.Client{}, nil, fmt.Errorf("get `Contract` file: %w", err)
		}

		filesMap[3] = file
	}

	var statusBy domain.User
	if request.StatusBy != nil && *request.StatusBy != uuid.Nil {
		statusBy, err = s.usersPostgresRepo.GetUser(ctx, domain.User{ID: *request.StatusBy})
		if err != nil {
			return domain.Request{}, domain.User{}, domain.Client{}, nil, fmt.Errorf("get `StatusBy` user: %w", err)
		}
	}

	var client domain.Client
	if request.ClientID != nil && *request.ClientID != uuid.Nil {
		client, err = s.clientService.Get(ctx, *request.ClientID)
		if err != nil {
			return domain.Request{}, domain.User{}, domain.Client{}, nil, fmt.Errorf("get `ClientID` client: %w", err)
		}
	}

	cache = cacheRequest{
		Request: request,
		User:    statusBy,
		Client:  client,
		Files:   filesMap,
	}

	if err := s.cache.Set(ctx, cacheKeyRequest(requestID), cache, cacheTTLRequest); err != nil {
		return domain.Request{}, domain.User{}, domain.Client{}, nil, fmt.Errorf("set redis cache: %w", err)
	}

	return request, statusBy, client, filesMap, nil
}
