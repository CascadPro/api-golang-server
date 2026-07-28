package root_transport_http

import (
	"net/http"
	"time"

	core_logger "github.com/CascadePro/api-golang-server/internal/core/logger"
	core_http_response "github.com/CascadePro/api-golang-server/internal/core/transport/http/response"
	core_http_utils "github.com/CascadePro/api-golang-server/internal/core/transport/http/utils"
	"go.uber.org/zap"
)

type PingResponse struct {
	Status     string    `json:"status"      example:"Ok"`
	StatusCode int       `json:"status_code" example:"200"`
	Timestamp  time.Time `json:"timestamp"   example:"2006-01-02T15-04-05.000000"`
}

// Ping godoc
// @Summary 		Ping
// @Description Ping server
// @Tags 				root
// @Produce 		json
// @Success 		200 {object} PingResponse "Default ping response"
// @Failure 		429 {object} core_http_response.ErrorResponse "Too many requests"
// @Failure		 	500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router		 	/ping [get]
func (h *HttpHandler) Ping(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewResponseHandler(log, rw)

	ip, err := core_http_utils.ClientIP(r)
	if err != nil {
		log.Error("failed to get client IP", zap.Error(err))
	}

	log.Info("ping client IP", zap.String("ip", ip.String()))

	response := PingResponse{
		Status:     "Ok",
		StatusCode: http.StatusOK,
		Timestamp:  time.Now(),
	}

	responseHandler.JsonResponse(response, http.StatusOK)
}
