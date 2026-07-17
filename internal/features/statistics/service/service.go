package statistics_service

import (
	"context"
	"time"

	"github.com/Svat-dev/golang-todo/internal/core/domain"
)

type StatisticsService struct {
	statisticsRepository StatisticsRepository
}

type StatisticsRepository interface {
	GetTasks(ctx context.Context, userID *int, from *time.Time, to *time.Time) ([]domain.Task, error)
}

func NewStatisticsService(statisticsRepository StatisticsRepository) *StatisticsService {
	return &StatisticsService{
		statisticsRepository: statisticsRepository,
	}
}
