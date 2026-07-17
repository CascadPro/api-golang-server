package statistics_postgres_repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Svat-dev/golang-todo/internal/core/domain"
	core_postgres_pool "github.com/Svat-dev/golang-todo/internal/core/repository/postgres/pool"
)

func (r *StatisticsRepository) GetTasks(
	ctx context.Context,
	userID *int,
	from *time.Time,
	to *time.Time,
) ([]domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	var query strings.Builder

	query.WriteString(`
		SELECT id, completed, created_at, completed_at
		FROM todoapp.tasks
	`)

	args := []any{}
	conditions := []string{}

	if userID != nil {
		conditions = append(conditions, fmt.Sprintf("author_user_id = $%d", len(args)+1))
		args = append(args, userID)
	}

	if from != nil {
		conditions = append(conditions, fmt.Sprintf("created_at >= $%d", len(args)+1))
		args = append(args, from)
	}

	if to != nil {
		conditions = append(conditions, fmt.Sprintf("created_at < $%d", len(args)+1))
		args = append(args, to)
	}

	if len(conditions) > 0 {
		query.WriteString(" WHERE ")
		query.WriteString(strings.Join(conditions, " AND "))
	}

	query.WriteString(" ORDER BY id ASC;")

	rows, err := r.pool.Query(ctx, query.String(), args...)
	if err != nil {
		err = core_postgres_pool.MapErrors(err)

		return nil, fmt.Errorf("select tasks: %w", err)
	}
	defer rows.Close()

	var taskModels []TaskModel
	for rows.Next() {
		var taskModel TaskModel

		err := rows.Scan(&taskModel.ID, &taskModel.Completed, &taskModel.CreatedAt, &taskModel.CompletedAt)
		if err != nil {
			err = core_postgres_pool.MapErrors(err)

			return nil, fmt.Errorf("scan tasks: %w", err)
		}

		taskModels = append(taskModels, taskModel)
	}
	if err := rows.Err(); err != nil {
		err = core_postgres_pool.MapErrors(err)

		return nil, fmt.Errorf("next rows: %w", err)
	}

	taskDomains := taskDomainsFromModels(taskModels)

	return taskDomains, nil
}
