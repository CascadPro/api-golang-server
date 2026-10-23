package app

import (
	"context"
	"fmt"
	"time"

	core_ipinfo_pool "github.com/CascadePro/api-golang-server/internal/core/infrastructure/ipinfo/pool"
	core_mongo_pool "github.com/CascadePro/api-golang-server/internal/core/infrastructure/mongo/pool"
	core_postgres_pool "github.com/CascadePro/api-golang-server/internal/core/infrastructure/postgres/pool"
	core_redis_cache "github.com/CascadePro/api-golang-server/internal/core/infrastructure/redis/cache"
	core_redis_pool "github.com/CascadePro/api-golang-server/internal/core/infrastructure/redis/pool"
	core_redis_realtime "github.com/CascadePro/api-golang-server/internal/core/infrastructure/redis/realtime"
	core_s3_pool "github.com/CascadePro/api-golang-server/internal/core/infrastructure/s3/pool"
	core_jwt_security "github.com/CascadePro/api-golang-server/internal/core/security/jwt"
	core_ws_middleware "github.com/CascadePro/api-golang-server/internal/core/transport/ws/middleware"
	core_validation_init "github.com/CascadePro/api-golang-server/internal/core/validation/init"
	idempotency_postgres_repository "github.com/CascadePro/api-golang-server/internal/features/idempotency/repository/postgres"
	idempotency_service "github.com/CascadePro/api-golang-server/internal/features/idempotency/service"
	"github.com/golang-jwt/jwt/v5"
)

type Infrastructure struct {
	Postgres        *core_postgres_pool.ConnectionPool
	Redis           *core_redis_pool.ConnectionPool
	IPInfoCache     core_redis_cache.Cache
	Mongo           *core_mongo_pool.ConnectionPool
	S3              *core_s3_pool.ConnectionPool
	IPInfo          *core_ipinfo_pool.ConnectionPool
	TokenIssuer     *core_jwt_security.Issuer
	Publisher       *core_redis_realtime.Publisher
	Idempotency     *idempotency_service.Service
	WsAuthenticator *core_ws_middleware.Authenticator
}

func (a *App) initInfrastructure(ctx context.Context) error {
	a.logger.Debug("initializing PostgreSQL connection pool")
	pgPool, err := core_postgres_pool.New(ctx, core_postgres_pool.NewConfigMust())
	if err != nil {
		return fmt.Errorf("initialize PostgreSQL: %w", err)
	}

	a.logger.Debug("initializing Redis connection pool")
	redisPool, err := core_redis_pool.New(ctx, core_redis_pool.NewConfigMust())
	if err != nil {
		pgPool.Close()
		return fmt.Errorf("initialize Redis: %w", err)
	}

	a.logger.Debug("initializing MongoDB connection pool")
	mongoPool, err := core_mongo_pool.New(ctx, core_mongo_pool.NewConfigMust())
	if err != nil {
		pgPool.Close()
		redisPool.Close()
		return fmt.Errorf("initialize MongoDB: %w", err)
	}

	a.logger.Debug("initializing S3 AWS connection pool")
	s3Pool, err := core_s3_pool.New(ctx, core_s3_pool.NewConfigMust())
	if err != nil {
		pgPool.Close()
		redisPool.Close()
		mongoPool.Close(ctx)
		return fmt.Errorf("initialize S3: %w", err)
	}

	a.logger.Debug("initializing IPInfo cache")
	ipInfoConfig := core_ipinfo_pool.NewConfigMust()
	cacheRepo, err := core_redis_cache.NewRedisCache(redisPool, "ipinfo", ipInfoConfig.CacheTTL)
	if err != nil {
		pgPool.Close()
		redisPool.Close()
		mongoPool.Close(ctx)
		return fmt.Errorf("initialize cache: %w", err)
	}

	a.logger.Debug("initializing IPInfo connection pool")
	ipInfoPool, err := core_ipinfo_pool.New(ipInfoConfig, cacheRepo)
	if err != nil {
		pgPool.Close()
		redisPool.Close()
		mongoPool.Close(ctx)
		return fmt.Errorf("initialize IPInfo: %w", err)
	}

	a.logger.Debug("initializing app validator")
	if err := core_validation_init.InitValidator(); err != nil {
		pgPool.Close()
		redisPool.Close()
		mongoPool.Close(ctx)
		return fmt.Errorf("initialize validator: %w", err)
	}

	a.logger.Debug("initializing app JWT issuer")
	tokenIssuer, err := core_jwt_security.NewIssuer(jwt.SigningMethodHS512)
	if err != nil {
		pgPool.Close()
		redisPool.Close()
		mongoPool.Close(ctx)
		return fmt.Errorf("initialize JWT issuer: %w", err)
	}

	publisher := core_redis_realtime.NewPublisher(redisPool)

	idempotencyRepository := idempotency_postgres_repository.NewRepository(pgPool)
	idempotency := idempotency_service.New(idempotencyRepository, 24*time.Hour)

	a.infrastructure = &Infrastructure{
		Postgres:    pgPool,
		Redis:       redisPool,
		Mongo:       mongoPool,
		S3:          s3Pool,
		IPInfo:      ipInfoPool,
		TokenIssuer: tokenIssuer,
		Publisher:   publisher,
		Idempotency: idempotency,
	}
	return nil
}

func (i *Infrastructure) Close(ctx context.Context) {
	if i.Postgres != nil {
		i.Postgres.Close()
	}
	if i.Redis != nil {
		i.Redis.Close()
	}
	if i.Mongo != nil {
		i.Mongo.Close(ctx)
	}
}
