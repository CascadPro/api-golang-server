package core_http_server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/CascadePro/api-golang-server/docs"
	core_logger "github.com/CascadePro/api-golang-server/internal/core/logger"
	core_http_middleware "github.com/CascadePro/api-golang-server/internal/core/transport/http/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)

type HttpServer struct {
	mux    *http.ServeMux
	config Config
	log    *core_logger.Logger

	middleware []core_http_middleware.Middleware

	serverMu sync.Mutex
	server   *http.Server
}

func NewHttpServer(
	config Config,
	log *core_logger.Logger,
	middleware ...core_http_middleware.Middleware,
) *HttpServer {
	return &HttpServer{
		mux:        http.NewServeMux(),
		config:     config,
		log:        log,
		middleware: middleware,
	}
}

func (s *HttpServer) RegisterRouters(routers ...*Router) {
	for _, router := range routers {
		s.mux.Handle(
			router.prefix+"/",
			http.StripPrefix(router.prefix, router.WithMiddleware()),
		)
	}
}

func (s *HttpServer) RegisterApiRouters(routers ...*ApiVersionRouter) {
	for _, router := range routers {
		prefix := "/api/" + string(router.version)

		s.mux.Handle(
			prefix+"/",
			http.StripPrefix(prefix, router.WithMiddleware()),
		)
	}
}

func (s *HttpServer) RegisterSwagger() {
	s.mux.Handle(
		"/swagger/",
		httpSwagger.Handler(
			httpSwagger.URL("/swagger/doc.json"),
		),
	)

	s.mux.HandleFunc(
		"/swagger/doc.json",
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			_, _ = w.Write(
				[]byte(docs.SwaggerInfo.ReadDoc()),
			)
		},
	)
}

func (s *HttpServer) Run(ctx context.Context) error {
	mux := core_http_middleware.ChainMiddleware(
		s.mux,
		s.middleware...,
	)

	server := &http.Server{
		Addr:    s.config.Addr,
		Handler: mux,
	}

	s.serverMu.Lock()
	s.server = server
	s.serverMu.Unlock()

	defer func() {
		s.serverMu.Lock()
		s.server = nil
		s.serverMu.Unlock()
	}()

	s.log.Warn(
		"starting http server...",
		zap.String("addr", s.config.Addr),
	)

	errCh := make(chan error, 1)

	go func() {
		err := server.ListenAndServe()

		if !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}

		close(errCh)
	}()

	select {
	case err := <-errCh:
		if err != nil {
			return fmt.Errorf("listen and serve http: %w", err)
		}

		return nil

	case <-ctx.Done():
		return nil
	}
}

func (s *HttpServer) Shutdown(ctx context.Context) error {
	s.serverMu.Lock()
	server := s.server
	s.serverMu.Unlock()

	if server == nil {
		return nil
	}

	s.log.Warn("shutting down http server...")

	if err := server.Shutdown(ctx); err != nil {
		_ = server.Close()

		return fmt.Errorf("shutdown http server: %w", err)
	}

	s.log.Warn("http server stopped")

	return nil
}
