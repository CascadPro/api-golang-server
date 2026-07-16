package tasks_postgres_repository

import (
	"time"

	"github.com/Svat-dev/golang-todo/internal/core/domain"
)

type TaskModel struct {
	ID      int
	Version int

	Title       string
	Description *string
	Completed   bool
	CreatedAt   time.Time
	CompletedAt *time.Time

	AuthorID int
}

func taskDomainFromModel(task TaskModel) domain.Task {
	return domain.Task{
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

func taskDomainsFromModels(tasks []TaskModel) []domain.Task {
	taskDomains := make([]domain.Task, len(tasks))

	for i, task := range tasks {
		taskDomains[i] = taskDomainFromModel(task)
	}

	return taskDomains
}
