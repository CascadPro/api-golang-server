package users_transport_http

import (
	"net/http"
	"time"

	core_context "github.com/CascadePro/api-golang-server/internal/core/context"
	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_logger "github.com/CascadePro/api-golang-server/internal/core/logger"
	core_http_response "github.com/CascadePro/api-golang-server/internal/core/transport/http/response"
	"github.com/google/uuid"
)

type GetCurrentUserResponse struct {
	ID                uuid.UUID       `json:"id"                           example:"00000000-000000-000000-000000000000"`
	Version           int             `json:"version"                      example:"1"`
	Role              domain.UserRole `json:"role"                         example:"regular"`
	Email             string          `json:"email"                        example:"test@example.com"`
	Name              string          `json:"name"                         example:"John"`
	Surname           string          `json:"surname"                      example:"Doe"`
	LastName          *string         `json:"last_name,omitempty"          example:""`
	LastActiveAt      time.Time       `json:"last_active_at"               example:"2006-01-02T15-04-05.000000"`
	AvatarFileID      *string         `json:"avatar_file_id,omitempty"     example:"831f21c1798a7972fa9cda12dac0"`
	AvatarPlaceholder []byte          `json:"avatar_placeholder,omitempty" example:"[bytes...]"`
}

func (h *HttpHandler) GetCurrentUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewResponseHandler(log, rw)

	userID, err := core_context.UserIDFromContext(ctx)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get userID from context")
		return
	}

	user, err := h.usersService.GetCurrentUser(ctx, userID)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get current user")
		return
	}

	var placeholder []byte
	if user.AvatarFileID != nil {
		placeholder, err = h.mediaService.GetFilePlaceholder(ctx, *user.AvatarFileID)
		if err != nil {
			responseHandler.ErrorResponse(err, "failed to get avatar placeholder")
			return
		}
	}

	response := currentUserDtoFromDomain(user, placeholder)

	responseHandler.JsonResponse(response, http.StatusOK)
}

func currentUserDtoFromDomain(user domain.User, placeholder []byte) GetCurrentUserResponse {
	return GetCurrentUserResponse{
		ID:                user.ID,
		Version:           user.Version,
		Role:              user.Role,
		Email:             *user.Email,
		Name:              user.Name,
		Surname:           user.Surname,
		LastName:          user.LastName,
		LastActiveAt:      user.LastActiveAt,
		AvatarFileID:      user.AvatarFileID,
		AvatarPlaceholder: placeholder,
	}
}
