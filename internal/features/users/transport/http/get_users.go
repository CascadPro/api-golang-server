package users_transport_http

import (
	"fmt"
	"net/http"

	core_logger "github.com/Svat-dev/golang-todo/internal/core/logger"
	core_http_response "github.com/Svat-dev/golang-todo/internal/core/transport/http/response"
	core_http_utils "github.com/Svat-dev/golang-todo/internal/core/transport/http/utils"
	users_http_dto "github.com/Svat-dev/golang-todo/internal/features/users/transport/http/dto"
)

type GetUsersResponse []users_http_dto.UserDtoResponse

func (h *UsersHttpHandler) GetUsers(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHttpResponseHandler(log, rw)

	limit, offset, err := getQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get 'limit' / 'offset' query param")

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
	limit, err := core_http_utils.GetIntQueryParam(r, "limit")
	if err != nil {
		return nil, nil, fmt.Errorf("get 'limit' query param: %w", err)
	}

	offset, err := core_http_utils.GetIntQueryParam(r, "offset")
	if err != nil {
		return nil, nil, fmt.Errorf("get 'offset' query param: %w", err)
	}

	return limit, offset, nil
}
