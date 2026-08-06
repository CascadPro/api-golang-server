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
	core_ipinfo_client "github.com/CascadePro/api-golang-server/internal/core/repository/ipinfo/client"
	core_ipinfo_pool "github.com/CascadePro/api-golang-server/internal/core/repository/ipinfo/pool"
	core_mongo_pool "github.com/CascadePro/api-golang-server/internal/core/repository/mongo/pool"
	core_postgres_pool "github.com/CascadePro/api-golang-server/internal/core/repository/postgres/pool"
	core_postgres_token "github.com/CascadePro/api-golang-server/internal/core/repository/postgres/token"
	core_redis_pool "github.com/CascadePro/api-golang-server/internal/core/repository/redis/pool"
	core_s3_pool "github.com/CascadePro/api-golang-server/internal/core/repository/s3/pool"
	core_http_middleware "github.com/CascadePro/api-golang-server/internal/core/transport/http/middleware"
	core_http_server "github.com/CascadePro/api-golang-server/internal/core/transport/http/server"
	core_validation_init "github.com/CascadePro/api-golang-server/internal/core/validation/init"
	auth_service "github.com/CascadePro/api-golang-server/internal/features/auth/service"
	auth_transport_http "github.com/CascadePro/api-golang-server/internal/features/auth/transport/http"
	media_postgres_repository "github.com/CascadePro/api-golang-server/internal/features/media/repository/postgres"
	media_service "github.com/CascadePro/api-golang-server/internal/features/media/service"
	media_transport_http "github.com/CascadePro/api-golang-server/internal/features/media/transport/http"
	requests_mongo_repository "github.com/CascadePro/api-golang-server/internal/features/requests/repository/mongo"
	requests_service "github.com/CascadePro/api-golang-server/internal/features/requests/service"
	requests_transport_http "github.com/CascadePro/api-golang-server/internal/features/requests/transport/http"
	root_transport_http "github.com/CascadePro/api-golang-server/internal/features/root/transport/http"
	sessions_redis_repository "github.com/CascadePro/api-golang-server/internal/features/sessions/repository/redis"
	session_service "github.com/CascadePro/api-golang-server/internal/features/sessions/service"
	sessions_transport_http "github.com/CascadePro/api-golang-server/internal/features/sessions/transport/http"
	settings_postgres_repository "github.com/CascadePro/api-golang-server/internal/features/settings/repository/postgres"
	settings_service "github.com/CascadePro/api-golang-server/internal/features/settings/service"
	settings_transport_http "github.com/CascadePro/api-golang-server/internal/features/settings/transport/http"
	users_postgres_repository "github.com/CascadePro/api-golang-server/internal/features/users/repository/postgres"
	users_service "github.com/CascadePro/api-golang-server/internal/features/users/service"
	users_transport_http "github.com/CascadePro/api-golang-server/internal/features/users/transport/http"
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

	logger.Debug("application environment mode", zap.Any("mode", cfg.EnvMode))

	logger.Debug("application connection", zap.Any("connection", cfg.Connection))

	logger.Debug("application time zone", zap.Any("zone", cfg.TimeZone))

	// >>> PostgresQL connect start
	logger.Debug("initializing PostgreSQL connection pool")
	pgPool, err := core_postgres_pool.New(ctx, core_postgres_pool.NewConfigMust())
	if err != nil {
		logger.Fatal("failed to init PostgreSQL connection pool", zap.Error(err))
	}
	defer pgPool.Close()
	// <<< PostgresQL connect end

	// >>> Redis connect start
	logger.Debug("initializing Redis connection pool")
	rdbPool, err := core_redis_pool.New(ctx, core_redis_pool.NewConfigMust())
	if err != nil {
		logger.Fatal("failed to init Redis connection pool", zap.Error(err))
	}
	defer rdbPool.Close()
	// <<< Redis connect end

	// >>> MongoDB connect start
	logger.Debug("initializing MongoDB connection pool")
	mongoPool, err := core_mongo_pool.New(ctx, core_mongo_pool.NewConfigMust())
	if err != nil {
		logger.Fatal("failed to init MongoDB connection pool", zap.Error(err))
	}
	defer mongoPool.Close(ctx)
	// <<< MongoDB connect end

	// >>> S3 connect start
	logger.Debug("initializing S3 AWS connection pool")
	s3Pool, err := core_s3_pool.New(ctx, core_s3_pool.NewConfigMust())
	if err != nil {
		logger.Fatal("failed to init S3 AWS connection pool", zap.Error(err))
	}
	// <<< S3 connect end

	// >>> IP Info connect start
	logger.Debug("initializing IP Info connection pool")
	ipInfoPool, err := core_ipinfo_pool.NewConnectionPool(ctx, core_ipinfo_pool.NewConfigMust(), rdbPool)
	if err != nil {
		logger.Fatal("failed to init IP Info connection pool", zap.Error(err))
	}
	// <<< IP Info connect end

	logger.Debug("initializing app validator")
	err = core_validation_init.InitValidator()
	if err != nil {
		logger.Fatal("failed to inti validator", zap.Error(err))
	}

	tokenPostgresRepository := core_postgres_token.NewRepository(pgPool)
	ipInfoRepository := core_ipinfo_client.NewRepository(ipInfoPool)

	// Root routes
	logFeatureInit(logger, "root", "")
	rootRateLimit := core_http_middleware.NewRateLimitConfig(100, time.Minute*5)
	rootRouter := core_http_server.NewRouter("", rootRateLimit.Middleware())

	rootHttpHandler := root_transport_http.NewHttpHandler()
	rootRouter.RegisterRoutes(rootHttpHandler.Routes()...)

	// Media routes
	logFeatureInit(logger, "media", "/media")
	mediaRateLimit := core_http_middleware.NewRateLimitConfig(25, time.Minute*10)
	mediaRouter := core_http_server.NewRouter("/media", mediaRateLimit.Middleware())

	mediaPostgresRepo := media_postgres_repository.NewRepository(pgPool)
	mediaService := media_service.NewService(mediaPostgresRepo, s3Pool)
	mediaHttpHandler := media_transport_http.NewHttpHandler(mediaService)

	mediaRouter.RegisterRoutes(mediaHttpHandler.Routes()...)

	// Users route
	logFeatureInit(logger, "users", "/api/v1/users")
	usersRouter := core_http_server.NewRouter("/users", rootRateLimit.Middleware(), core_http_middleware.Authorization())

	usersPostgresRepository := users_postgres_repository.NewRepository(pgPool)
	usersService := users_service.NewService(mediaService, usersPostgresRepository)
	usersHttpHandler := users_transport_http.NewHttpHandler(usersService)

	usersRouter.RegisterRoutes(usersHttpHandler.Routes()...)

	// Settings route
	logFeatureInit(logger, "settings", "/api/v1/settings")
	settingsRouter := core_http_server.NewRouter("/settings",
		rootRateLimit.Middleware(), core_http_middleware.Authorization())

	settingsPostgresRepository := settings_postgres_repository.NewRepository(pgPool)
	settingsService := settings_service.NewService(settingsPostgresRepository)
	settingsHttpHandler := settings_transport_http.NewHttpHandler(settingsService)

	settingsRouter.RegisterRoutes(settingsHttpHandler.Routes()...)

	// Requests route
	logFeatureInit(logger, "requests", "/api/v1/requests")
	requestsRouter := core_http_server.NewRouter("/requests", rootRateLimit.Middleware())

	requestsMongoRepository := requests_mongo_repository.NewRepository(mongoPool)
	requestsService := requests_service.NewService(mediaService, usersPostgresRepository, requestsMongoRepository)
	requestsHttpHandler := requests_transport_http.NewHttpHandler(requestsService)

	requestsRouter.RegisterRoutes(requestsHttpHandler.Routes()...)

	// Sessions route
	logFeatureInit(logger, "sessions", "/api/v1/sessions")
	sessionsRouter := core_http_server.NewRouter("/sessions", core_http_middleware.Authorization())

	sessionsRedisRepo := sessions_redis_repository.NewRepository(rdbPool)
	sessionsService := session_service.NewService(sessionsRedisRepo)
	sessionsHttpHandler := sessions_transport_http.NewHttpHandler(sessionsService)

	sessionsRouter.RegisterRoutes(sessionsHttpHandler.Routes()...)

	// Auth routes
	logFeatureInit(logger, "authentication", "/api/v1/auth")
	authRouter := core_http_server.NewRouter("/auth")

	authService := auth_service.NewService(usersPostgresRepository, settingsPostgresRepository,
		tokenPostgresRepository, ipInfoRepository, sessionsRedisRepo)
	authHttpHandler := auth_transport_http.NewHttpHandler(authService)

	authRouter.RegisterRoutes(authHttpHandler.Routes()...)

	// Base app setup
	logger.Debug("initializing http server")

	httpServer := core_http_server.NewHttpServer(
		core_http_server.NewConfigMust(),
		logger,
		core_http_middleware.CORS(),
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Trace(),
		core_http_middleware.Panic(),
		core_http_middleware.IP(),
	)

	apiVersionRouter := core_http_server.NewApiVersionRouter(core_http_server.ApiVersion1)
	apiVersionRouter.RegisterRouters(sessionsRouter, authRouter, usersRouter, settingsRouter, requestsRouter)

	httpServer.RegisterRouters(rootRouter, mediaRouter)
	httpServer.RegisterApiRouters(apiVersionRouter)
	httpServer.RegisterSwagger()

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("http server run error", zap.Error(err))
	}
}

func logFeatureInit(log *core_logger.Logger, name, path string) {
	log.Debug(
		"initialing feature",
		zap.String("feature", name),
		zap.String("path", fmt.Sprintf("%s/*", path)),
	)
}
