package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	core_config "github.com/CascadePro/api-golang-server/internal/core/config"
	core_logger "github.com/CascadePro/api-golang-server/internal/core/logger"
	core_postgres_pool "github.com/CascadePro/api-golang-server/internal/core/repository/postgres/pool"
	core_redis_pool "github.com/CascadePro/api-golang-server/internal/core/repository/redis/pool"
	core_s3_pool "github.com/CascadePro/api-golang-server/internal/core/repository/s3/pool"
	core_http_middleware "github.com/CascadePro/api-golang-server/internal/core/transport/http/middleware"
	core_http_server "github.com/CascadePro/api-golang-server/internal/core/transport/http/server"
	test_postgres_repository "github.com/CascadePro/api-golang-server/internal/features/test/repository/postgres"
	test_service "github.com/CascadePro/api-golang-server/internal/features/test/service"
	test_transport_http "github.com/CascadePro/api-golang-server/internal/features/test/transport/http"
	"go.uber.org/zap"

	_ "github.com/CascadePro/api-golang-server/docs"
)

// @title 				Cascade Pro App API
// @version 			1.0.0
// @description 	Cascade App Pro API
// @host 					127.0.0.1:8000
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

	// >>> PostgresQL connect start
	logger.Debug("initializing PostgreSQL connection pool")
	pgPool, err := core_postgres_pool.NewConnectionPool(ctx, core_postgres_pool.NewConfigMust())
	if err != nil {
		logger.Fatal("failed to init PostgreSQL connection pool", zap.Error(err))
	}
	defer pgPool.Close()
	// <<< PostgresQL connect end

	// >>> Redis connect start
	logger.Debug("initializing Redis connection pool")
	rdbPool, err := core_redis_pool.NewConnectionPool(ctx, core_redis_pool.NewConfigMust())
	if err != nil {
		logger.Fatal("failed to init Redis connection pool", zap.Error(err))
	}
	defer rdbPool.Close()
	// <<< Redis connect end

	// >>> S3 connect start
	logger.Debug("initializing S3 AWS connection pool")
	_, err = core_s3_pool.NewConnectionPool(ctx, core_s3_pool.NewConfigMust())
	if err != nil {
		logger.Fatal("failed to init S3 AWS connection pool", zap.Error(err))
	}
	// <<< S3 connect end

	const featureInitKey = "initializing feature"

	var testRoutes []core_http_server.Route
	{
		logger.Debug(featureInitKey, zap.String("feature", "statistics"))

		repo := test_postgres_repository.NewTestRepository(pgPool)
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
