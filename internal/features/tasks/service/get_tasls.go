package tasks_service

import (
	"context"
	"fmt"

	"github.com/Svat-dev/golang-todo/internal/core/domain"
	core_errors "github.com/Svat-dev/golang-todo/internal/core/errors"
)

func (s *TasksService) GetTasks(ctx context.Context, userID *int, limit *int, offset *int) ([]domain.Task, error) {
	if userID != nil && *userID < 0 {
		return nil, fmt.Errorf("'user_id' value must be non-negative: %w", core_errors.ErrInvalidArgument)
	}

	if limit != nil && *limit < 0 {
		return nil, fmt.Errorf("'limit' value must be non-negative: %w", core_errors.ErrInvalidArgument)
	}

	if offset != nil && *offset < 0 {
		return nil, fmt.Errorf("'offset' value must be non-negative: %w", core_errors.ErrInvalidArgument)
	}

	tasks, err := s.tasksRepository.GetTasks(ctx, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("get tasks from repository")
	}

	return tasks, err
}
