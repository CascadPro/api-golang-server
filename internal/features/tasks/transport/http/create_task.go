package tasks_transport_http

import (
	"net/http"

	"github.com/Svat-dev/golang-todo/internal/core/domain"
	core_logger "github.com/Svat-dev/golang-todo/internal/core/logger"
	core_http_request "github.com/Svat-dev/golang-todo/internal/core/transport/http/request"
	core_http_response "github.com/Svat-dev/golang-todo/internal/core/transport/http/response"
	tasks_http_dto "github.com/Svat-dev/golang-todo/internal/features/tasks/transport/http/dto"
)

type CreateTaskRequest struct {
	Title       string  `json:"title" validate:"required,min=1,max=100"`
	Description *string `json:"description" validate:"omitempty,min=1,max=1000"`
	AuthorID    int     `json:"author_id" validate:"required"`
}

type CreateTaskResponse tasks_http_dto.TaskDtoResponse

func (h *TasksHttpHandler) CreateTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHttpResponseHandler(log, rw)

	var request CreateTaskRequest
	if err := core_http_request.DecodeAndValidate(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate http request")

		return
	}

	taskDomain := domain.NewTaskUninitialized(request.Title, request.Description, request.AuthorID)

	task, err := h.tasksService.CreateTask(r.Context(), taskDomain)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to create task")

		return
	}

	response := CreateTaskResponse(tasks_http_dto.TaskDtoFromDomain(task))

	responseHandler.JsonResponse(response, http.StatusOK)
}
