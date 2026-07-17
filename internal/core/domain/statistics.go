package domain

import "time"

type Statistics struct {
	TasksCreated             int
	TasksCompleted           int
	TasksCompletedRate       *float64
	TasksAverageCompleteTime *time.Duration
}

func NewStatistics(created int, completed int, rate *float64, averageTime *time.Duration) Statistics {
	return Statistics{
		TasksCreated:             created,
		TasksCompleted:           completed,
		TasksCompletedRate:       rate,
		TasksAverageCompleteTime: averageTime,
	}
}
