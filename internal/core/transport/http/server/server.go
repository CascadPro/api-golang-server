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

func (s *HttpServer) RegisterApiRouters(routers ...*ApiVersionRouter) {
	for _, router := range routers {
		prefix := "/api/" + string(router.apiVersion)

		s.mux.Handle(prefix+"/", http.StripPrefix(prefix, router.WithMiddleware()))
	}
}

func (s *HttpServer) Run(ctx context.Context) error {
	mux := core_http_middleware.ChainMiddleware(s.mux, s.middleware...)

	server := &http.Server{
		Addr:    s.config.Addr,
		Handler: mux,
	}

	ch := make(chan error, 1)

	go func() {
		defer close(ch)

		s.log.Warn("starting http server...", zap.String("addr", s.config.Addr))

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
		s.log.Warn("shutdown http server...")

		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			s.config.ShutdownTimeout,
		)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			_ = server.Close()

			return fmt.Errorf("shutdown http server: %w", err)
		}

		s.log.Warn("http server stopped")
	}

	return nil
}
