package statistics_postgres_repository

import (
	"time"

	"github.com/Svat-dev/golang-todo/internal/core/domain"
)

type TaskModel struct {
	ID          int
	Completed   bool
	CreatedAt   time.Time
	CompletedAt *time.Time
}

func taskDomainFromModel(task TaskModel) domain.Task {
	return domain.Task{
		ID:          task.ID,
		Completed:   task.Completed,
		CreatedAt:   task.CreatedAt,
		CompletedAt: task.CompletedAt,
	}
}

func taskDomainsFromModels(tasks []TaskModel) []domain.Task {
	taskDomains := make([]domain.Task, len(tasks))

	for i, task := range tasks {
		taskDomains[i] = taskDomainFromModel(task)
	}

	return taskDomains
}
