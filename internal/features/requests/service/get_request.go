package requests_service

import (
	"context"
	"fmt"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	"github.com/google/uuid"
)

func (s *Service) GetRequest(
	ctx context.Context,
	requestID uuid.UUID,
) (domain.Request, domain.User, map[int]domain.File, error) {
	if requestID == uuid.Nil {
		return domain.Request{}, domain.User{}, nil, fmt.Errorf("`id` can't be NULL: %w", core_errors.ErrInvalidArgument)
	}

	request, err := s.requestsMongoRepo.GetRequest(ctx, requestID)
	if err != nil {
		return domain.Request{}, domain.User{}, nil, fmt.Errorf("get request from repository: %w", err)
	}

	filesMap := make(map[int]domain.File, 4)

	if request.ProjectDocID != nil {
		file, err := s.mediaService.GetFile(ctx, *request.ProjectDocID)
		if err != nil {
			return domain.Request{}, domain.User{}, nil, fmt.Errorf("get `ProjectPlan` file: %w", err)
		}

		filesMap[0] = file
	}
	if request.TechTaskDocID != nil {
		file, err := s.mediaService.GetFile(ctx, *request.TechTaskDocID)
		if err != nil {
			return domain.Request{}, domain.User{}, nil, fmt.Errorf("get `TechTask` file: %w", err)
		}

		filesMap[1] = file
	}
	if request.SpecificationDocID != nil {
		file, err := s.mediaService.GetFile(ctx, *request.SpecificationDocID)
		if err != nil {
			return domain.Request{}, domain.User{}, nil, fmt.Errorf("get `Specification` file: %w", err)
		}

		filesMap[2] = file
	}
	if request.ContractDocID != nil {
		file, err := s.mediaService.GetFile(ctx, *request.ContractDocID)
		if err != nil {
			return domain.Request{}, domain.User{}, nil, fmt.Errorf("get `Contract` file: %w", err)
		}

		filesMap[3] = file
	}

	var statusBy domain.User
	if request.StatusBy != nil && *request.StatusBy != uuid.Nil {
		statusBy, err = s.usersPostgresRepo.GetUser(ctx, domain.User{ID: *request.StatusBy})
		if err != nil {
			return domain.Request{}, domain.User{}, nil, fmt.Errorf("get `StatusBy` user: %w", err)
		}
	}

	return request, statusBy, filesMap, nil
}
