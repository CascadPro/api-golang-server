package root_transport_http

import (
	"net/http"
	"time"

	core_logger "github.com/CascadePro/api-golang-server/internal/core/logger"
	core_http_response "github.com/CascadePro/api-golang-server/internal/core/transport/http/response"
)

type RootResponse struct {
	Status      string    `json:"status"      example:"Ok"`
	StatusCode  int       `json:"status_code" example:"200"`
	Message     string    `json:"message"     example:"Greeting message"`
	Description string    `json:"description" example:"Server API description"`
	Docs        string    `json:"docs"        example:"Link to server API docs"`
	Timestamp   time.Time `json:"timestamp"   example:"2006-01-02T15-04-05.000000"`
}

// Root godoc
// @Summary 		Root
// @Description Root path of the server. (Title screen)
// @Tags 				root
// @Produce 		json
// @Success 		200 {object} RootResponse "Default root response"
// @Failure 		429 {object} core_http_response.ErrorResponse "Too many requests"
// @Failure		 	500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router		 	/ [get]
func (h *HttpHandler) Root(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewResponseHandler(log, rw)

	response := RootResponse{
		Status:      "Ok",
		StatusCode:  http.StatusOK,
		Message:     "You've reached Cascade Pro API server",
		Description: "It's server for Cascade Pro App, which you can download on mobile or desktop platforms",
		Docs:        "/swagger/index.html",
		Timestamp:   time.Now(),
	}

	responseHandler.JsonResponse(response, http.StatusOK)
}
