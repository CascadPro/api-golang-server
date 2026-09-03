package requests_service

import (
	"context"
	"fmt"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	"github.com/google/uuid"
)

func (s *Service) UploadDoc(
	ctx context.Context,
	requestID uuid.UUID,
	file domain.File,
	content []byte,
	index int,
) error {
	if requestID == uuid.Nil {
		return fmt.Errorf("`requestID` can't be NULL: %w", core_errors.ErrInvalidArgument)
	}
	if index < 0 || index > 3 {
		return fmt.Errorf("`index` can only have values between [0,3]: %w", core_errors.ErrInvalidArgument)
	}
	if err := file.Validate(); err != nil {
		return fmt.Errorf("validate file: %w", core_errors.ErrInvalidArgument)
	}

	request, err := s.requestsMongoRepo.GetRequest(ctx, requestID)
	if err != nil {
		return fmt.Errorf("get request from repository: %w", err)
	}

	file, err = s.mediaService.UploadFile(ctx, &file, content)
	if err != nil {
		return fmt.Errorf("upload file to repository: %w", err)
	}

	fileID := domain.NewNullable(file.ID)

	var patch domain.RequestPatch
	var oldFileID *string

	switch index {
	case 0:
		patch.ProjectDocID = fileID
		oldFileID = request.ProjectDocID
	case 1:
		patch.TechTaskDocID = fileID
		oldFileID = request.TechTaskDocID
	case 2:
		patch.SpecificationDocID = fileID
		oldFileID = request.SpecificationDocID
	case 3:
		patch.ContractDocID = fileID
		oldFileID = request.ContractDocID
	}

	if err := request.ApplyPatch(patch); err != nil {
		_ = s.mediaService.DeleteFile(ctx, file.Tag, file.ID)

		return fmt.Errorf("apply patch: %w", err)
	}

	if err := s.requestsMongoRepo.PatchRequest(ctx, requestID, request); err != nil {
		_ = s.mediaService.DeleteFile(ctx, file.Tag, file.ID)

		return fmt.Errorf("patch request: %w", err)
	}

	if oldFileID != nil {
		_ = s.mediaService.DeleteFile(ctx, file.Tag, *oldFileID)
	}

	return nil
}
