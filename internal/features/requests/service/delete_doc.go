package requests_service

import (
	"context"
	"fmt"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	"github.com/google/uuid"
)

func (s *Service) DeleteDoc(ctx context.Context, requestID uuid.UUID, index int) error {
	if requestID == uuid.Nil {
		return fmt.Errorf("`requestID` can't be NULL: %w", core_errors.ErrInvalidArgument)
	}
	if index < 0 || index > 3 {
		return fmt.Errorf("`index` can only have values between [0,3]: %w", core_errors.ErrInvalidArgument)
	}

	request, err := s.requestsMongoRepo.GetRequest(ctx, requestID)
	if err != nil {
		return fmt.Errorf("get request from repository: %w", err)
	}

	var fileID string
	var patch domain.RequestPatch

	null := domain.NewNullable("")
	switch index {
	case 0:
		patch.ProjectDocID = null
		fileID = *request.ProjectDocID
	case 1:
		patch.TechTaskDocID = null
		fileID = *request.TechTaskDocID
	case 2:
		patch.SpecificationDocID = null
		fileID = *request.SpecificationDocID
	case 3:
		patch.ContractDocID = null
		fileID = *request.ContractDocID
	}

	if err := request.ApplyPatch(patch); err != nil {
		return fmt.Errorf("apply patch: %w", err)
	}

	if err := s.requestsMongoRepo.PatchRequest(ctx, requestID, request); err != nil {
		return fmt.Errorf("patch request: %w", err)
	}

	if err := s.mediaService.DeleteFile(ctx, domain.FileTagDocs, fileID); err != nil {
		return fmt.Errorf("delete file from repository: %w", err)
	}

	return nil
}
