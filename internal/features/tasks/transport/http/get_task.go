package tasks_transport_http

import (
	"net/http"

	core_logger "github.com/Svat-dev/golang-todo/internal/core/logger"
	core_http_request "github.com/Svat-dev/golang-todo/internal/core/transport/http/request"
	core_http_response "github.com/Svat-dev/golang-todo/internal/core/transport/http/response"
	tasks_http_dto "github.com/Svat-dev/golang-todo/internal/features/tasks/transport/http/dto"
)

type GetTaskResponse tasks_http_dto.TaskDtoResponse

func (h *TasksHttpHandler) GetTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHttpResponseHandler(log, rw)

	taskID, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get taskID path value")

		return
	}

	taskDomain, err := h.tasksService.GetTask(ctx, taskID)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get task")

		return
	}

	response := GetTaskResponse(tasks_http_dto.TaskDtoFromDomain(taskDomain))

	responseHandler.JsonResponse(response, http.StatusOK)
}
