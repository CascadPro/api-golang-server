package tasks_http_dto

import (
	"time"

	"github.com/Svat-dev/golang-todo/internal/core/domain"
)

type TaskDtoResponse struct {
	ID          int        `json:"id"`
	Version     int        `json:"version"`
	Title       string     `json:"title"`
	Description *string    `json:"description"`
	Completed   bool       `json:"completed"`
	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at"`
	AuthorID    int        `json:"author_id"`
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
