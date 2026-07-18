package tasks_transport_http

import (
	"net/http"

	core_logger "github.com/Svat-dev/golang-todo/internal/core/logger"
	core_http_request "github.com/Svat-dev/golang-todo/internal/core/transport/http/request"
	core_http_response "github.com/Svat-dev/golang-todo/internal/core/transport/http/response"
)

// DeleteTask		godoc
// @Summary 		Delete task
// @Description Delete a task in the system
// @Tags 				tasks
// @Param       id path int true "Task ID"
// @Success 		204 "No content response"
// @Failure 		401 {object} core_http_response.ErrorResponse "Bad request"
// @Failure 		409 {object} core_http_response.ErrorResponse "Conflict error"
// @Failure		 	500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router		 	/tasks/{id} [delete]
func (h *TasksHttpHandler) DeleteTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHttpResponseHandler(log, rw)

	taskID, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get taskID path value")

		return
	}

	err = h.tasksService.DeleteTask(ctx, taskID)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to delete task")
	}

	responseHandler.NoContentResponse()
}
