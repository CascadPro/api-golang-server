package users_transport_http

import (
	"fmt"
	"net/http"

	core_logger "github.com/Svat-dev/golang-todo/internal/core/logger"
	core_http_request "github.com/Svat-dev/golang-todo/internal/core/transport/http/request"
	core_http_response "github.com/Svat-dev/golang-todo/internal/core/transport/http/response"
	users_http_dto "github.com/Svat-dev/golang-todo/internal/features/users/transport/http/dto"
)

type GetUsersResponse []users_http_dto.UserDtoResponse

// GetUsers 		godoc
// @Summary 		Get users
// @Description Get users from the system
// @Tags 				users
// @Produce 		json
// @Param       limit query int false "Limit of users"
// @Param       offset query int false "Offset of users"
// @Success 		200 {object} GetUsersResponse "Successfully got users"
// @Failure 		401 {object} core_http_response.ErrorResponse "Bad request"
// @Failure 		409 {object} core_http_response.ErrorResponse "Conflict error"
// @Failure		 	500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router		 	/users [get]
func (h *UsersHttpHandler) GetUsers(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHttpResponseHandler(log, rw)

	limit, offset, err := getQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get limit/offset query param")

		return
	}

	userDomains, err := h.usersService.GetUsers(ctx, limit, offset)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get users")

		return
	}

	response := GetUsersResponse(users_http_dto.UsersDtoFromDomains(userDomains))

	responseHandler.JsonResponse(response, http.StatusOK)
}

func getQueryParams(r *http.Request) (*int, *int, error) {
	const (
		limitParamKey  = "limit"
		offsetParamKey = "offset"
	)

	limit, err := core_http_request.GetIntQueryParam(r, limitParamKey)
	if err != nil {
		return nil, nil, fmt.Errorf("get '%s' query param: %w", limitParamKey, err)
	}

	offset, err := core_http_request.GetIntQueryParam(r, offsetParamKey)
	if err != nil {
		return nil, nil, fmt.Errorf("get '%s' query param: %w", offsetParamKey, err)
	}

	return limit, offset, nil
}
