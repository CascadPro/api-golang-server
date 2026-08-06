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
	Version           int64           `json:"version"                      example:"1"`
	Role              domain.UserRole `json:"role"                         example:"regular"`
	Email             string          `json:"email"                        example:"test@example.com"`
	Name              string          `json:"name"                         example:"John"`
	Surname           string          `json:"surname"                      example:"Doe"`
	LastName          *string         `json:"last_name,omitempty"          example:""`
	LastActiveAt      time.Time       `json:"last_active_at"               example:"2006-01-02T15-04-05.000000"`
	AvatarFileID      *string         `json:"avatar_file_id,omitempty"     example:"831f21c1798a7972fa9cda12dac0"`
	AvatarPlaceholder []byte          `json:"avatar_placeholder,omitempty" example:"0"`
}

// GetCurrentUser godoc
// @Summary      Get current user
// @Description  Get user data for the authenticated user
// @Tags         users
// @Produce      json
// @Success      200 {object} GetCurrentUserResponse "Current user data"
// @Failure      401 {object} core_http_response.ErrorResponse "Unauthorized"
// @Failure      404 {object} core_http_response.ErrorResponse "Not found"
// @Failure      409 {object} core_http_response.ErrorResponse "Conflict error"
// @Failure      429 {object} core_http_response.ErrorResponse "Too many requests"
// @Failure      500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router       /users/my [get]
func (h *HttpHandler) GetCurrentUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewResponseHandler(log, rw)

	userID, err := core_context.UserIDFromContext(ctx)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get userID from context")
		return
	}

	user, placeholder, err := h.usersService.GetCurrentUser(ctx, userID)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get current user")
		return
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
