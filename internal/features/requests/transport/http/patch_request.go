package requests_transport_http

import (
	"net/http"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_logger "github.com/CascadePro/api-golang-server/internal/core/logger"
	core_http_request "github.com/CascadePro/api-golang-server/internal/core/transport/http/request"
	core_http_response "github.com/CascadePro/api-golang-server/internal/core/transport/http/response"
	request_http_dto "github.com/CascadePro/api-golang-server/internal/features/requests/transport/http/dto"
)

type PatchRequestRequest struct {
	request_http_dto.PatchRequestRequest
}

// PatchRequest godoc
// @Summary      Patch request
// @Description  Patch request by ID
// @Tags         requests
// @Accept			 json
// @Param        id path string true "Request ID"
// @Param 			 request body PatchRequestRequest true "Patch request body"
// @Success      204 "Patched request"
// @Failure      400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure      401 {object} core_http_response.ErrorResponse "Unauthorized"
// @Failure      404 {object} core_http_response.ErrorResponse "Request not found"
// @Failure 		 409 {object} core_http_response.ErrorResponse "Conflict error"
// @Failure      429 {object} core_http_response.ErrorResponse "Too many requests"
// @Failure      500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router       /requests/{id}/update [patch]
func (h *HttpHandler) PatchRequest(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewResponseHandler(log, rw)

	requestID, err := core_http_request.GetUUIDPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get request id path value")
		return
	}

	var request PatchRequestRequest
	if err := core_http_request.DecodeAndValidate(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate http request")
		return
	}

	patch := patchRequestToPatch(request)

	if _, err := h.requestsService.PatchRequest(ctx, requestID, patch); err != nil {
		responseHandler.ErrorResponse(err, "failed to patch request")
		return
	}

	responseHandler.NoContentResponse()
}

func patchRequestToPatch(request PatchRequestRequest) domain.RequestPatch {
	patch := domain.RequestPatch{
		Title:     request.Title.ToDomain(),
		ClientID:  request.ClientID.ToDomain(),
		WorkTypes: request.WorkTypes.ToDomain(),
		Geography: request.Geography.ToDomain(),
		Deadline:  request.Deadline.ToDomain(),
	}

	if request.Origin.Set && request.Origin.Value != nil {
		origin := make([]domain.RequestOrigin, len(*request.Origin.Value))
		for i, j := range *request.Origin.Value {
			origin[i] = domain.RequestOrigin{Type: j.Type, Value: j.Value}
		}
		patch.Origin = domain.NewNullable(origin)
	}

	return patch
}
