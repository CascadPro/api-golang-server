package request_http_dto

import (
	"fmt"
	"net/http"

	core_http_request "github.com/CascadePro/api-golang-server/internal/core/transport/http/request"
	"github.com/google/uuid"
)

func GetRequestDocPathValues(r *http.Request) (uuid.UUID, int, error) {
	var min, max = 0, 3

	requestID, err := core_http_request.GetUUIDPathValue(r, "id")
	if err != nil {
		return uuid.Nil, -1, fmt.Errorf("get `id` path value: %w", err)
	}

	index, err := core_http_request.GetIntPathValue(r, "index", &min, &max)
	if err != nil {
		return uuid.Nil, -1, fmt.Errorf("get `index` path value: %w", err)
	}

	return requestID, index, nil
}
