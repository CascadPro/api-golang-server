package media_transport_http

import (
	"context"
	"io"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_http "github.com/CascadePro/api-golang-server/internal/core/transport/http"
	core_http_server "github.com/CascadePro/api-golang-server/internal/core/transport/http/server"
)

type HttpHandler struct {
	mediaService MediaService
}

type MediaService interface {
	GetFileContent(ctx context.Context, fileID string, w, h, quality *int) (
		file domain.File, content []byte, err error,
	)

	GetFileStream(ctx context.Context, fileID, byteRange string) (
		file domain.File, body io.ReadCloser, size int64, contentRange string, err error,
	)
}

func NewHttpHandler(mediaService MediaService) *HttpHandler {
	return &HttpHandler{
		mediaService: mediaService,
	}
}

func (h *HttpHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  core_http.MethodGet,
			Path:    "/{tag}/{filename}",
			Handler: h.GetFile,
		},
	}
}
