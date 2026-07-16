package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	core_config "github.com/Svat-dev/golang-todo/internal/core/config"
	core_logger "github.com/Svat-dev/golang-todo/internal/core/logger"
	core_postgres_pool "github.com/Svat-dev/golang-todo/internal/core/repository/postgres/pool"
	core_http_middleware "github.com/Svat-dev/golang-todo/internal/core/transport/http/middleware"
	core_http_server "github.com/Svat-dev/golang-todo/internal/core/transport/http/server"
	users_postgres_repository "github.com/Svat-dev/golang-todo/internal/features/users/repository/postgres"
	users_service "github.com/Svat-dev/golang-todo/internal/features/users/service"
	users_transport_http "github.com/Svat-dev/golang-todo/internal/features/users/transport/http"
	"go.uber.org/zap"
)

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

	var usersRoutes []core_http_server.Route
	{
		logger.Debug("initializing feature", zap.String("feature", "users"))

		repo := users_postgres_repository.NewUsersRepository(pool)
		service := users_service.NewUsersService(repo)
		httpTransport := users_transport_http.NewUsersHttpHandler(service)

		usersRoutes = httpTransport.Routes()
	}

	logger.Debug("initializing http server")

	httpServer := core_http_server.NewHttpServer(
		core_http_server.NewConfigMust(),
		logger,
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Trace(),
		core_http_middleware.Panic(),
	)

	apiVersionRouter := core_http_server.NewApiVersionRouter(core_http_server.ApiVersion1)
	apiVersionRouter.RegisterRouter(usersRoutes...)

	httpServer.RegisterApiRouters(apiVersionRouter)

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("http server run error", zap.Error(err))
	}
}
