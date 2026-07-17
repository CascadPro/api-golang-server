package statistics_transport_http

import (
	"context"
	"net/http"
	"time"

	"github.com/Svat-dev/golang-todo/internal/core/domain"
	core_http_server "github.com/Svat-dev/golang-todo/internal/core/transport/http/server"
)

type StatisticsHttpHandler struct {
	statisticsService StatisticsService
}

type StatisticsService interface {
	GetStatistics(ctx context.Context, userID *int, from *time.Time, to *time.Time) (domain.Statistics, error)
}

func NewStatisticsHttpHandler(statisticsService StatisticsService) *StatisticsHttpHandler {
	return &StatisticsHttpHandler{
		statisticsService: statisticsService,
	}
}

func (h *StatisticsHttpHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodGet,
			Path:    "/statistics",
			Handler: h.GetStatistics,
		},
	}
}
