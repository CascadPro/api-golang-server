package statistics_transport_http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Svat-dev/golang-todo/internal/core/domain"
	core_logger "github.com/Svat-dev/golang-todo/internal/core/logger"
	core_http_request "github.com/Svat-dev/golang-todo/internal/core/transport/http/request"
	core_http_response "github.com/Svat-dev/golang-todo/internal/core/transport/http/response"
)

type GetStatisticsResponse struct {
	TasksCreated             int      `json:"tasks_created"`
	TasksCompleted           int      `json:"tasks_completed"`
	TasksCompletedRate       *float64 `json:"tasks_completed_rate"`
	TasksAverageCompleteTime *string  `json:"tasks_avg_completion_time"`
}

func (h *StatisticsHttpHandler) GetStatistics(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHttpResponseHandler(log, rw)

	userID, from, to, err := getQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get user_id/from/to query param")

		return
	}

	stats, err := h.statisticsService.GetStatistics(ctx, userID, from, to)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get statistics")

		return
	}

	response := toDtoFromDomain(stats)

	responseHandler.JsonResponse(response, http.StatusOK)
}

func toDtoFromDomain(stats domain.Statistics) GetStatisticsResponse {
	var avgTime *string
	if stats.TasksAverageCompleteTime != nil {
		duration := stats.TasksAverageCompleteTime.String()
		avgTime = &duration
	}

	return GetStatisticsResponse{
		TasksCreated:             stats.TasksCreated,
		TasksCompleted:           stats.TasksCompleted,
		TasksCompletedRate:       stats.TasksCompletedRate,
		TasksAverageCompleteTime: avgTime,
	}
}

func getQueryParams(r *http.Request) (*int, *time.Time, *time.Time, error) {
	const (
		userIDParamKey = "user_id"
		fromParamKey   = "from"
		toParamKey     = "offset"
	)

	userID, err := core_http_request.GetIntQueryParam(r, userIDParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get '%s' query param: %w", userIDParamKey, err)
	}

	from, err := core_http_request.GetDateQueryParam(r, fromParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get '%s' query param: %w", fromParamKey, err)
	}

	to, err := core_http_request.GetDateQueryParam(r, toParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get '%s' query param: %w", toParamKey, err)
	}

	return userID, from, to, nil
}
