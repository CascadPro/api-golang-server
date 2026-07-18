package tasks_transport_http

import (
	"net/http"

	"github.com/Svat-dev/golang-todo/internal/core/domain"
	core_logger "github.com/Svat-dev/golang-todo/internal/core/logger"
	core_http_request "github.com/Svat-dev/golang-todo/internal/core/transport/http/request"
	core_http_response "github.com/Svat-dev/golang-todo/internal/core/transport/http/response"
	tasks_http_dto "github.com/Svat-dev/golang-todo/internal/features/tasks/transport/http/dto"
)

type PatchTaskRequest tasks_http_dto.PatchTaskRequest

type PatchTaskResponse tasks_http_dto.TaskDtoResponse

// PatchTask 		godoc
// @Summary 		Patch a task
// @Description Patch a task in the system
// @Tags 				tasks
// @Accept 			json
// @Produce 		json
// @Param       id path int true "Task ID"
// @Param 			request body tasks_http_dto.PatchTaskRequestSwagger true "PatchTask body request"
// @Success 		200 {object} PatchTaskResponse "Successfully patched task"
// @Failure 		401 {object} core_http_response.ErrorResponse "Bad request"
// @Failure 		409 {object} core_http_response.ErrorResponse "Conflict error"
// @Failure		 	500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router		 	/tasks/{id} [patch]
func (h *TasksHttpHandler) PatchTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHttpResponseHandler(log, rw)

	taskID, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get taskID path value")

		return
	}

	var request PatchTaskRequest
	if err := core_http_request.DecodeAndValidate(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate http request")

		return
	}

	userPatch := userPatchFromRequest(request)

	userDomain, err := h.tasksService.PatchTask(ctx, taskID, userPatch)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to patch task")

		return
	}

	response := PatchTaskResponse(tasks_http_dto.TaskDtoFromDomain(userDomain))

	responseHandler.JsonResponse(response, http.StatusOK)
}

func userPatchFromRequest(request PatchTaskRequest) domain.TaskPatch {
	return domain.NewTaskPatch(
		request.Title.ToDomain(),
		request.Description.ToDomain(),
		request.Completed.ToDomain(),
	)
}
