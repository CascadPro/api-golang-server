package app

import (
	"context"
	"fmt"

	core_ipinfo_pool "github.com/CascadePro/api-golang-server/internal/core/infrastructure/ipinfo/pool"
	core_mongo_pool "github.com/CascadePro/api-golang-server/internal/core/infrastructure/mongo/pool"
	core_postgres_pool "github.com/CascadePro/api-golang-server/internal/core/infrastructure/postgres/pool"
	core_redis_pool "github.com/CascadePro/api-golang-server/internal/core/infrastructure/redis/pool"
	core_s3_pool "github.com/CascadePro/api-golang-server/internal/core/infrastructure/s3/pool"
	core_jwt_security "github.com/CascadePro/api-golang-server/internal/core/security/jwt"
	core_validation_init "github.com/CascadePro/api-golang-server/internal/core/validation/init"
	"github.com/golang-jwt/jwt/v5"
)

type Infrastructure struct {
	Postgres    *core_postgres_pool.ConnectionPool
	Redis       *core_redis_pool.ConnectionPool
	Mongo       *core_mongo_pool.ConnectionPool
	S3          *core_s3_pool.ConnectionPool
	IPInfo      *core_ipinfo_pool.ConnectionPool
	TokenIssuer *core_jwt_security.Issuer
}

func (a *App) initInfrastructure(ctx context.Context) error {
	// >>> PostgreSQL connecting
	a.logger.Debug("initializing PostgreSQL connection pool")

	pgPool, err := core_postgres_pool.New(ctx, core_postgres_pool.NewConfigMust())
	if err != nil {
		return fmt.Errorf("initialize PostgreSQL: %w", err)
	}

	// >>> Redis connecting
	a.logger.Debug("initializing Redis connection pool")

	redisPool, err := core_redis_pool.New(ctx, core_redis_pool.NewConfigMust())
	if err != nil {
		pgPool.Close()

		return fmt.Errorf("initialize Redis: %w", err)
	}

	// >>> MongoDB connecting
	a.logger.Debug("initializing MongoDB connection pool")

	mongoPool, err := core_mongo_pool.New(ctx, core_mongo_pool.NewConfigMust())
	if err != nil {
		pgPool.Close()
		redisPool.Close()

		return fmt.Errorf("initialize MongoDB: %w", err)
	}

	// >>> S3 connecting
	a.logger.Debug("initializing S3 AWS connection pool")

	s3Pool, err := core_s3_pool.New(ctx, core_s3_pool.NewConfigMust())
	if err != nil {
		pgPool.Close()
		redisPool.Close()
		mongoPool.Close(ctx)

		return fmt.Errorf("initialize S3: %w", err)
	}

	// >>> IP Info connect start
	a.logger.Debug("initializing IP Info connection pool")

	ipInfoPool, err := core_ipinfo_pool.New(ctx, core_ipinfo_pool.NewConfigMust(), redisPool)
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

		return fmt.Errorf("initialize JWT issuer")
	}

	a.infrastructure = &Infrastructure{
		Postgres:    pgPool,
		Redis:       redisPool,
		Mongo:       mongoPool,
		S3:          s3Pool,
		IPInfo:      ipInfoPool,
		TokenIssuer: tokenIssuer,
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
