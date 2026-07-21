package test_transport_http

import (
	"net/http"

	core_http_server "github.com/CascadePro/api-golang-server/internal/core/transport/http/server"
)

type TestHttpHandler struct {
	testService TestService
}

type TestService interface {
}

func NewTestHttpHandler(testService TestService) *TestHttpHandler {
	return &TestHttpHandler{
		testService: testService,
	}
}

func (h *TestHttpHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodGet,
			Path:    "/test",
			Handler: /* Some handler */,
		},
	}
}
