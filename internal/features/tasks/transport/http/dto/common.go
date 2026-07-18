package tasks_http_dto

import (
	"time"

	"github.com/Svat-dev/golang-todo/internal/core/domain"
)

type TaskDtoResponse struct {
	ID          int        `json:"id" example:"10"`
	Version     int        `json:"version" example:"3"`
	Title       string     `json:"title" example:"New Task 10"`
	Description *string    `json:"description" example:"A simple new task 10!"`
	Completed   bool       `json:"completed" example:"false"`
	CreatedAt   time.Time  `json:"created_at" example:"2026-07-16T16:11:53.256251Z"`
	CompletedAt *time.Time `json:"completed_at" example:"null"`
	AuthorID    int        `json:"author_id" example:"23"`
}

func TaskDtoFromDomain(task domain.Task) TaskDtoResponse {
	return TaskDtoResponse{
		ID:          task.ID,
		Version:     task.Version,
		Title:       task.Title,
		Description: task.Description,
		Completed:   task.Completed,
		CreatedAt:   task.CreatedAt,
		CompletedAt: task.CompletedAt,
		AuthorID:    task.AuthorID,
	}
}

func TasksDtoFromDomains(tasks []domain.Task) []TaskDtoResponse {
	tasksDto := make([]TaskDtoResponse, len(tasks))

	for i, task := range tasks {
		tasksDto[i] = TaskDtoFromDomain(task)
	}

	return tasksDto
}
