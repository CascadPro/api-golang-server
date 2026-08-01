package users_transport_http

import (
	"net/http"

	core_context "github.com/CascadePro/api-golang-server/internal/core/context"
	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	core_logger "github.com/CascadePro/api-golang-server/internal/core/logger"
	core_http_request "github.com/CascadePro/api-golang-server/internal/core/transport/http/request"
	core_http_response "github.com/CascadePro/api-golang-server/internal/core/transport/http/response"
)

func (h *HttpHandler) UpdateAvatar(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewResponseHandler(log, rw)

	userID, err := core_context.UserIDFromContext(ctx)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get user ID from context")
		return
	}

	const maxFileSize = 10 << 20 // 10 MB
	uploadedFile, content, err := core_http_request.GetFile(rw, r, maxFileSize)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get file from request")
		return
	}
	if uploadedFile.Tag != domain.FileTagAvatars {
		responseHandler.ErrorResponse(core_errors.ErrInvalidArgument, "invalid file tag, expected 'avatars'")
		return
	}

	if err := uploadedFile.GeneratePlaceholder(content); err != nil {
		responseHandler.ErrorResponse(err, "failed to generate placeholder for file")
		return
	}

	if err := h.usersService.UpdateAvatar(ctx, userID, uploadedFile, content); err != nil {
		responseHandler.ErrorResponse(err, "failed to update user avatar")
		return
	}

	responseHandler.NoContentResponse()
}
