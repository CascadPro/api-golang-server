package tasks_transport_http

import (
	"fmt"
	"net/http"

	core_logger "github.com/Svat-dev/golang-todo/internal/core/logger"
	core_http_request "github.com/Svat-dev/golang-todo/internal/core/transport/http/request"
	core_http_response "github.com/Svat-dev/golang-todo/internal/core/transport/http/response"
	tasks_http_dto "github.com/Svat-dev/golang-todo/internal/features/tasks/transport/http/dto"
)

type GetTasksResponse []tasks_http_dto.TaskDtoResponse

// GetTasks 		godoc
// @Summary 		Get tasks
// @Description Get tasks from the system
// @Tags 				tasks
// @Produce 		json
// @Param       user_id query int false "User ID to get tasks by user"
// @Param       limit query int false "Limit of tasks"
// @Param       offset query int false "Offset of tasks"
// @Success 		200 {object} GetTasksResponse "Successfully got tasks"
// @Failure 		401 {object} core_http_response.ErrorResponse "Bad request"
// @Failure 		409 {object} core_http_response.ErrorResponse "Conflict error"
// @Failure		 	500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router		 	/tasks [get]
func (h *TasksHttpHandler) GetTasks(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHttpResponseHandler(log, rw)

	userID, limit, offset, err := getQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get user_id/limit/offset query param")

		return
	}

	taskDomains, err := h.tasksService.GetTasks(ctx, userID, limit, offset)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get tasks")

		return
	}

	response := GetTasksResponse(tasks_http_dto.TasksDtoFromDomains(taskDomains))

	responseHandler.JsonResponse(response, http.StatusOK)
}

func getQueryParams(r *http.Request) (*int, *int, *int, error) {
	const (
		userIDParamKey = "user_id"
		limitParamKey  = "limit"
		offsetParamKey = "offset"
	)

	userID, err := core_http_request.GetIntQueryParam(r, userIDParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get '%s' query param: %w", userIDParamKey, err)
	}

	limit, err := core_http_request.GetIntQueryParam(r, limitParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get '%s' query param: %w", limitParamKey, err)
	}

	offset, err := core_http_request.GetIntQueryParam(r, offsetParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get '%s' query param: %w", offsetParamKey, err)
	}

	return userID, limit, offset, nil
}
