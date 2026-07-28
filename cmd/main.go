package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	core_config "github.com/CascadePro/api-golang-server/internal/core/config"
	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_logger "github.com/CascadePro/api-golang-server/internal/core/logger"
	core_ipinfo_client "github.com/CascadePro/api-golang-server/internal/core/repository/ipinfo/client"
	core_ipinfo_pool "github.com/CascadePro/api-golang-server/internal/core/repository/ipinfo/pool"
	core_postgres_pool "github.com/CascadePro/api-golang-server/internal/core/repository/postgres/pool"
	core_postgres_token "github.com/CascadePro/api-golang-server/internal/core/repository/postgres/token"
	core_postgres_user "github.com/CascadePro/api-golang-server/internal/core/repository/postgres/user"
	core_redis_pool "github.com/CascadePro/api-golang-server/internal/core/repository/redis/pool"
	core_s3_pool "github.com/CascadePro/api-golang-server/internal/core/repository/s3/pool"
	core_http_middleware "github.com/CascadePro/api-golang-server/internal/core/transport/http/middleware"
	core_http_server "github.com/CascadePro/api-golang-server/internal/core/transport/http/server"
	core_validation "github.com/CascadePro/api-golang-server/internal/core/validation"
	auth_service "github.com/CascadePro/api-golang-server/internal/features/auth/service"
	auth_transport_http "github.com/CascadePro/api-golang-server/internal/features/auth/transport/http"
	root_transport_http "github.com/CascadePro/api-golang-server/internal/features/root/transport/http"
	sessions_redis_repository "github.com/CascadePro/api-golang-server/internal/features/sessions/repository/redis"
	session_service "github.com/CascadePro/api-golang-server/internal/features/sessions/service"
	sessions_transport_http "github.com/CascadePro/api-golang-server/internal/features/sessions/transport/http"
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

	// >>> IP Info connect start
	logger.Debug("initializing IP Info connection pool")
	ipInfoPool, err := core_ipinfo_pool.NewConnectionPool(ctx, core_ipinfo_pool.NewConfigMust(), rdbPool)
	if err != nil {
		logger.Fatal("failed to init IP Info connection pool", zap.Error(err))
	}
	// <<< IP Info connect end

	logger.Debug("initializing app validator")
	err = core_validation.InitValidator(domain.Roles, domain.SessionIDByteLength)
	if err != nil {
		logger.Fatal("failed to inti validator", zap.Error(err))
	}

	const featureInitMsg = "initialing feature"

	userPostgresRepository := core_postgres_user.NewRepository(pgPool)
	tokenPostgresRepository := core_postgres_token.NewRepository(pgPool)
	ipInfoRepository := core_ipinfo_client.NewRepository(ipInfoPool)

	// Root routes
	logger.Debug(featureInitMsg, zap.String("feature", "root"), zap.String("path", "/*"))
	rootRateLimit := core_http_middleware.NewRateLimitConfig(100, time.Minute*5)
	rootRouter := core_http_server.NewRouter("", rootRateLimit.Middleware())

	rootHttpHandler := root_transport_http.NewHttpHandler()
	rootRouter.RegisterRoutes(rootHttpHandler.Routes()...)

	// Sessions route
	logger.Debug(featureInitMsg, zap.String("feature", "sessions"), zap.String("path", "/sessions/*"))
	sessionsRouter := core_http_server.NewRouter("/sessions", core_http_middleware.Authorization())

	sessionsRedisRepo := sessions_redis_repository.NewRepository(rdbPool)
	sessionsService := session_service.NewService(sessionsRedisRepo)
	sessionsHttpHandler := sessions_transport_http.NewHttpHandler(sessionsService)

	sessionsRouter.RegisterRoutes(sessionsHttpHandler.Routes()...)

	// Auth routes
	logger.Debug(featureInitMsg, zap.String("feature", "auth"), zap.String("path", "/auth/*"))
	authRouter := core_http_server.NewRouter("/auth")

	authService := auth_service.NewService(userPostgresRepository, tokenPostgresRepository,
		ipInfoRepository, sessionsRedisRepo)
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
	)

	apiVersionRouter := core_http_server.NewApiVersionRouter(core_http_server.ApiVersion1)
	apiVersionRouter.RegisterRouters(sessionsRouter, authRouter)

	httpServer.RegisterRouters(rootRouter)
	httpServer.RegisterApiRouters(apiVersionRouter)
	httpServer.RegisterSwagger()

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("http server run error", zap.Error(err))
	}
}
