package tasks_transport_http

import (
	"context"
	"net/http"

	"github.com/Svat-dev/golang-todo/internal/core/domain"
	core_http_server "github.com/Svat-dev/golang-todo/internal/core/transport/http/server"
)

type TasksHttpHandler struct {
	tasksService TasksService
}

type TasksService interface {
	CreateTask(ctx context.Context, task domain.Task) (domain.Task, error)
	GetTasks(ctx context.Context, userID *int, limit *int, offset *int) ([]domain.Task, error)
	GetTask(ctx context.Context, id int) (domain.Task, error)
	DeleteTask(ctx context.Context, id int) error
	PatchTask(ctx context.Context, id int, patch domain.TaskPatch) (domain.Task, error)
}

func NewTasksHttpHandler(tasksService TasksService) *TasksHttpHandler {
	return &TasksHttpHandler{
		tasksService: tasksService,
	}
}

func (h *TasksHttpHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/tasks",
			Handler: h.CreateTask,
		},
		{
			Method:  http.MethodGet,
			Path:    "/tasks",
			Handler: h.GetTasks,
		},
		{
			Method:  http.MethodGet,
			Path:    "/tasks/{id}",
			Handler: h.GetTask,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/tasks/{id}",
			Handler: h.DeleteTask,
		},
		{
			Method:  http.MethodPatch,
			Path:    "/tasks/{id}",
			Handler: h.PatchTask,
		},
	}
}
