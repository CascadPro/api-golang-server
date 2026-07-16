package tasks_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Svat-dev/golang-todo/internal/core/domain"
	core_errors "github.com/Svat-dev/golang-todo/internal/core/errors"
	core_postgres_pool "github.com/Svat-dev/golang-todo/internal/core/repository/postgres/pool"
)

func (r *TasksRepository) CreateTask(ctx context.Context, task domain.Task) (domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		INSERT INTO todoapp.tasks (title, description, completed, created_at, completed_at, author_user_id)
		VALUES($1, $2, $3, $4, $5, $6)
		RETURNING id, version, title, description, completed, created_at, completed_at, author_user_id;
	`

	row := r.pool.QueryRow(ctx, query,
		task.Title,
		task.Description,
		task.Completed,
		task.CreatedAt,
		task.CompletedAt,
		task.AuthorID,
	)

	var taskModel TaskModel

	err := row.Scan(
		&taskModel.ID,
		&taskModel.Version,
		&taskModel.Title,
		&taskModel.Description,
		&taskModel.Completed,
		&taskModel.CreatedAt,
		&taskModel.CompletedAt,
		&taskModel.AuthorID,
	)
	if err != nil {
		err = core_postgres_pool.MapErrors(err)
		if errors.Is(err, core_postgres_pool.ErrViolatesForeignKey) {
			return domain.Task{}, fmt.Errorf("%v: user with id=%d: %w", err, task.AuthorID, core_errors.ErrNotFound)
		}

		return domain.Task{}, fmt.Errorf("scan err: %w", err)
	}

	taskDomain := taskDomainFromModel(taskModel)

	return taskDomain, nil
}
