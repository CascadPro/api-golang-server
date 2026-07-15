package core_http_server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	core_logger "github.com/Svat-dev/golang-todo/internal/core/logger"
	core_http_middleware "github.com/Svat-dev/golang-todo/internal/core/transport/http/middleware"
	"go.uber.org/zap"
)

type HttpServer struct {
	mux    *http.ServeMux
	config Config
	log    *core_logger.Logger

	middleware []core_http_middleware.Middleware
}

func NewHttpServer(config Config, log *core_logger.Logger, middleware ...core_http_middleware.Middleware) *HttpServer {
	return &HttpServer{
		mux:        http.NewServeMux(),
		config:     config,
		log:        log,
		middleware: middleware,
	}
}

func (h *HttpServer) RegisterApiRouters(routers ...*ApiVersionRouter) {
	for _, router := range routers {
		prefix := "/api/" + string(router.apiVersion)

		h.mux.Handle(prefix+"/", http.StripPrefix(prefix, router))
	}
}

func (h *HttpServer) Run(ctx context.Context) error {
	mux := core_http_middleware.ChainMiddleware(h.mux, h.middleware...)

	server := &http.Server{
		Addr:    h.config.Addr,
		Handler: mux,
	}

	ch := make(chan error, 1)

	go func() {
		defer close(ch)

		h.log.Warn("starting http server...", zap.String("addr", h.config.Addr))

		err := server.ListenAndServe()

		if !errors.Is(err, http.ErrServerClosed) {
			ch <- err
		}
	}()

	select {
	case err := <-ch:
		if err != nil {
			return fmt.Errorf("listen and serve http: %w", err)
		}
	case <-ctx.Done():
		h.log.Warn("shutdown http server...")

		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			h.config.ShutdownTimeout,
		)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			_ = server.Close()

			return fmt.Errorf("shutdown http server: %w", err)
		}

		h.log.Warn("http server stopped")
	}

	return nil
}
