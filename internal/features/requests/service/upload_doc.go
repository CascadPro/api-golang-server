package requests_service

import (
	"context"
	"fmt"
	"io"

	core_context "github.com/CascadePro/api-golang-server/internal/core/context"
	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (s *Service) UploadDoc(
	ctx context.Context,
	requestID uuid.UUID,
	file domain.File,
	content []byte,
	index int,
) error {
	return s.uploadDocInternal(ctx, requestID, &file, nil, 0, index, content)
}

func (s *Service) UploadDocStream(
	ctx context.Context,
	requestID uuid.UUID,
	file domain.File,
	content io.Reader,
	size int64,
	index int,
) error {
	return s.uploadDocInternal(ctx, requestID, &file, content, size, index, nil)
}

func (s *Service) uploadDocInternal(
	ctx context.Context,
	requestID uuid.UUID,
	file *domain.File,
	content io.Reader,
	size int64,
	index int,
	contentBytes []byte,
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

	request, err := s.requestsMongoRepo.Get(ctx, requestID)
	if err != nil {
		return fmt.Errorf("get request from repository: %w", err)
	}

	var uploadedFile domain.File
	if content != nil {
		uploadedFile, err = s.mediaService.UploadFileStream(ctx, file, content, size)
	} else {
		uploadedFile, err = s.mediaService.UploadFile(ctx, file, contentBytes)
	}
	if err != nil {
		return fmt.Errorf("upload file to repository: %w", err)
	}

	fileID := domain.NewNullable(uploadedFile.ID)

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
		_ = s.mediaService.DeleteFile(ctx, uploadedFile.Tag, uploadedFile.ID)

		return fmt.Errorf("apply patch: %w", err)
	}

	if err := s.requestsMongoRepo.Patch(ctx, requestID, request, false); err != nil {
		_ = s.mediaService.DeleteFile(ctx, uploadedFile.Tag, uploadedFile.ID)

		return fmt.Errorf("patch request: %w", err)
	}

	if oldFileID != nil {
		_ = s.mediaService.DeleteFile(ctx, uploadedFile.Tag, *oldFileID)
	}

	_ = s.cache.Delete(ctx, cacheKeyRequest(requestID))

	userID, err := core_context.UserID(ctx)
	if err != nil {
		return fmt.Errorf("get user id from context: %w", err)
	}

	data := domain.NewRealtimeEventRequestFileData(
		uploadedFile.ID,
		request.ID.String(),
		index,
		uploadedFile.Filename,
		uploadedFile.Size,
		string(uploadedFile.MimeType),
		uploadedFile.CreatedAt,
	)

	event := domain.NewRealtimeEvent(domain.RealtimeEventRequestFileUploaded, &userID, nil, data)

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
