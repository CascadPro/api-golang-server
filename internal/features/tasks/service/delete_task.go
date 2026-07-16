package tasks_service

import (
	"context"
	"fmt"
)

func (s *TasksService) DeleteTask(ctx context.Context, id int) error {
	err := s.tasksRepository.DeleteTask(ctx, id)
	if err != nil {
		return fmt.Errorf("delete task from repository: %w", err)
	}

	return nil
}
