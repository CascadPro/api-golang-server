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
	Title       string  `json:"title" validate:"required,min=1,max=100" example:"New Task"`
	Description *string `json:"description" validate:"omitempty,min=1,max=1000" example:"A simple new task!"`
	AuthorID    int     `json:"author_id" validate:"required" example:"1"`
}

type CreateTaskResponse tasks_http_dto.TaskDtoResponse

// CreateTask 	godoc
// @Summary 		Create a task
// @Description Create a new task in the system
// @Tags 				tasks
// @Accept 			json
// @Produce 		json
// @Param 			request body CreateTaskRequest true "CreateTask body request"
// @Success 		200 {object} CreateTaskResponse "Successfully created task"
// @Failure 		401 {object} core_http_response.ErrorResponse "Bad request"
// @Failure 		409 {object} core_http_response.ErrorResponse "Conflict error"
// @Failure		 	500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router		 	/tasks [post]
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
