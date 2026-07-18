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
	statistics_postgres_repository "github.com/Svat-dev/golang-todo/internal/features/statistics/repository/postgres"
	statistics_service "github.com/Svat-dev/golang-todo/internal/features/statistics/service"
	statistics_transport_http "github.com/Svat-dev/golang-todo/internal/features/statistics/transport/http"
	tasks_postgres_repository "github.com/Svat-dev/golang-todo/internal/features/tasks/repository/postgres"
	tasks_service "github.com/Svat-dev/golang-todo/internal/features/tasks/service"
	tasks_transport_http "github.com/Svat-dev/golang-todo/internal/features/tasks/transport/http"
	users_postgres_repository "github.com/Svat-dev/golang-todo/internal/features/users/repository/postgres"
	users_service "github.com/Svat-dev/golang-todo/internal/features/users/service"
	users_transport_http "github.com/Svat-dev/golang-todo/internal/features/users/transport/http"
	"go.uber.org/zap"

	_ "github.com/Svat-dev/golang-todo/docs"
)

// @title 				Golang ToDo API
// @version 			1.0
// @description 	ToDo Application REST-API scheme
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

	var usersRoutes []core_http_server.Route
	{
		logger.Debug(featureInitKey, zap.String("feature", "users"))

		repo := users_postgres_repository.NewUsersRepository(pool)
		service := users_service.NewUsersService(repo)
		httpTransport := users_transport_http.NewUsersHttpHandler(service)

		usersRoutes = httpTransport.Routes()
	}

	var tasksRoutes []core_http_server.Route
	{
		logger.Debug(featureInitKey, zap.String("feature", "tasks"))

		repo := tasks_postgres_repository.NewTasksRepository(pool)
		service := tasks_service.NewTasksService(repo)
		httpTransport := tasks_transport_http.NewTasksHttpHandler(service)

		tasksRoutes = httpTransport.Routes()
	}

	var statisticsRoutes []core_http_server.Route
	{
		logger.Debug(featureInitKey, zap.String("feature", "statistics"))

		repo := statistics_postgres_repository.NewStatisticsRepository(pool)
		service := statistics_service.NewStatisticsService(repo)
		httpTransport := statistics_transport_http.NewStatisticsHttpHandler(service)

		statisticsRoutes = httpTransport.Routes()
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
	apiVersionRouter.RegisterRouter(tasksRoutes...)
	apiVersionRouter.RegisterRouter(statisticsRoutes...)

	httpServer.RegisterApiRouters(apiVersionRouter)
	httpServer.RegisterSwagger()

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("http server run error", zap.Error(err))
	}
}
