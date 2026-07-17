package statistics_service

import (
	"context"
	"fmt"
	"time"

	"github.com/Svat-dev/golang-todo/internal/core/domain"
	core_errors "github.com/Svat-dev/golang-todo/internal/core/errors"
)

func (s *StatisticsService) GetStatistics(
	ctx context.Context,
	userID *int,
	from *time.Time,
	to *time.Time,
) (domain.Statistics, error) {
	if from != nil && to != nil {
		if to.Before(*from) || to.Equal(*from) {
			return domain.Statistics{}, fmt.Errorf("`to` must be after `from`: %w", core_errors.ErrInvalidArgument)
		}
	}

	tasks, err := s.statisticsRepository.GetTasks(ctx, userID, from, to)
	if err != nil {
		return domain.Statistics{}, fmt.Errorf("failed to get statistics from repository: %w", err)
	}

	stats := calcStatistics(tasks)

	return stats, nil
}

func calcStatistics(tasks []domain.Task) domain.Statistics {
	if len(tasks) == 0 {
		return domain.NewStatistics(0, 0, nil, nil)
	}

	totalTasksCreated := len(tasks)

	var totalTasksCompleted int
	var totalCompletionDuration time.Duration

	for _, task := range tasks {
		if task.Completed {
			totalTasksCompleted += 1
		}

		completionDuration := task.CompletionDuration()
		if completionDuration != nil {
			totalCompletionDuration += *completionDuration
		}
	}

	tasksCompletedRate := float64(totalTasksCompleted) / float64(totalTasksCreated) * 100

	var avgTimeDuration *time.Duration
	if totalTasksCompleted > 0 && totalCompletionDuration != 0 {
		avg := totalCompletionDuration / time.Duration(totalTasksCompleted)
		avgTimeDuration = &avg
	}

	return domain.NewStatistics(totalTasksCreated, totalTasksCompleted, &tasksCompletedRate, avgTimeDuration)
}
