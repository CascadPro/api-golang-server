package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	core_config "github.com/Svat-dev/golang-todo/internal/core/config"
	core_logger "github.com/Svat-dev/golang-todo/internal/core/logger"
	core_postgres_pool "github.com/Svat-dev/golang-todo/internal/core/repository/postgres/pool"
	core_http_middleware "github.com/Svat-dev/golang-todo/internal/core/transport/http/middleware"
	core_http_server "github.com/Svat-dev/golang-todo/internal/core/transport/http/server"
	test_postgres_repository "github.com/Svat-dev/golang-todo/internal/features/test/repository/postgres"
	test_service "github.com/Svat-dev/golang-todo/internal/features/test/service"
	test_transport_http "github.com/Svat-dev/golang-todo/internal/features/test/transport/http"
	"go.uber.org/zap"

	_ "github.com/Svat-dev/golang-todo/docs"
)

// @title 				Test API
// @version 			1.0.0
// @description 	Test description
// @host 					127.0.0.1:8080
// @BasePath 			/api/v1
func main() {
	cfg := core_config.NewConfigMust()
	time.Local = cfg.TimeZone

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM,
	)
	defer cancel()

	logger, err := core_logger.NewLogger(core_logger.NewConfigMust())
	if err != nil {
		fmt.Println("failed to init app logger:", err)
		os.Exit(1)
	}
	defer logger.Close()

	logger.Debug("application time zone", zap.Any("zone", cfg.TimeZone))

	logger.Debug("initializing PostgreSQL connection pool")
	pool, err := core_postgres_pool.NewConnectionPool(ctx, core_postgres_pool.NewConfigMust())
	if err != nil {
		logger.Fatal("failed to init PostgreSQL connection pool", zap.Error(err))
	}
	defer pool.Close()

	const featureInitKey = "initializing feature"

	var testRoutes []core_http_server.Route
	{
		logger.Debug(featureInitKey, zap.String("feature", "statistics"))

		repo := test_postgres_repository.NewTestRepository(pool)
		service := test_service.NewTestService(repo)
		httpTransport := test_transport_http.NewTestHttpHandler(service)

		testRoutes = httpTransport.Routes()
	}

	logger.Debug("initializing http server")

	httpServer := core_http_server.NewHttpServer(
		core_http_server.NewConfigMust(),
		logger,
		core_http_middleware.CORS(),
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Trace(),
		core_http_middleware.Panic(),
	)

	apiVersionRouter := core_http_server.NewApiVersionRouter(core_http_server.ApiVersion1)
	apiVersionRouter.RegisterRouter(testRoutes...)

	httpServer.RegisterApiRouters(apiVersionRouter)
	httpServer.RegisterSwagger()

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("http server run error", zap.Error(err))
	}
}
