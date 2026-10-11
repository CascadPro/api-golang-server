package requests_transport_http

import (
	"net/http"

	core_context "github.com/CascadePro/api-golang-server/internal/core/context"
	core_logger "github.com/CascadePro/api-golang-server/internal/core/logger"
	core_http_response "github.com/CascadePro/api-golang-server/internal/core/transport/http/response"
)

type GetRequestStatsResponse struct {
	Total    int64            `json:"total" example:"42"`
	ByStatus map[string]int64 `json:"by_status"`
}

// GetRequestStats godoc
// @Summary      Get request statistics
// @Description  Request statistics by status
// @Tags         requests
// @Produce      json
// @Success      200 {object} GetRequestStatsResponse "Statistics"
// @Failure      401 {object} core_http_response.ErrorResponse "Unauthorized"
// @Failure      429 {object} core_http_response.ErrorResponse "Too many requests"
// @Failure      500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router       /requests/stats [get]
func (h *HttpHandler) GetRequestStats(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	locale := core_context.Locale(ctx)

	responseHandler := core_http_response.NewResponseHandler(log, locale, rw)

	stats, err := h.requestsService.GetRequestsStats(ctx)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get request stats")
		return
	}

	response := GetRequestStatsResponse{
		Total:    stats.Total,
		ByStatus: make(map[string]int64, len(stats.ByStatus)),
	}

	for status, count := range stats.ByStatus {
		response.ByStatus[string(status)] = count
	}

	responseHandler.JsonResponse(response, http.StatusOK)
}
