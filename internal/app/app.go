package app

import (
	"context"
	"fmt"

	core_config "github.com/CascadePro/api-golang-server/internal/core/config"
	core_i18n "github.com/CascadePro/api-golang-server/internal/core/i18n"
	core_redis_rate_limit "github.com/CascadePro/api-golang-server/internal/core/infrastructure/redis/rate_limit"
	core_logger "github.com/CascadePro/api-golang-server/internal/core/logger"
	core_rate_limit "github.com/CascadePro/api-golang-server/internal/core/rate_limit"
	core_http_server "github.com/CascadePro/api-golang-server/internal/core/transport/http/server"
	core_ws_server "github.com/CascadePro/api-golang-server/internal/core/transport/ws/server"
	"go.uber.org/zap"
)

type App struct {
	cfg    *core_config.Config
	logger *core_logger.Logger

	infrastructure *Infrastructure
	features       *Features

	rateLimiter       core_rate_limit.Limiter
	rateLimitObserver *core_rate_limit.CounterObserver
	rateLimits        rateLimitPolicies

	httpServer *core_http_server.HttpServer
	wsServer   *core_ws_server.Server

	workerCancel context.CancelFunc
	workerDone   chan struct{}
}

func New(ctx context.Context, cfg *core_config.Config, logger *core_logger.Logger) (*App, error) {
	logger.Info("application environment mode", zap.Any("mode", cfg.EnvMode))
	logger.Info("application connection", zap.Any("connection", cfg.Connection))
	logger.Info("application time zone", zap.Any("zone", cfg.TimeZone))

	if cfg.RateLimitFailureMode == "" {
		cfg.RateLimitFailureMode = string(core_rate_limit.FailureClosed)
	}

	app := &App{
		cfg:               cfg,
		logger:            logger,
		rateLimits:        defaultRateLimitPolicies(),
		rateLimitObserver: core_rate_limit.NewCounterObserver(),
	}

	if err := validateRateLimitPolicies(app.rateLimits); err != nil {
		return nil, err
	}
	if err := core_rate_limit.FailureMode(cfg.RateLimitFailureMode).Validate(); err != nil {
		return nil, fmt.Errorf("validate rate limit failure mode: %w", err)
	}

	logger.Debug("initializing i18n localization")
	if err := core_i18n.Init(); err != nil {
		return nil, fmt.Errorf("init localization: %w", err)
	}

	logger.Debug("initializing app infrastructure")
	if err := app.initInfrastructure(ctx); err != nil {
		return nil, fmt.Errorf("init infrastructure: %w", err)
	}

	logger.Debug("initializing app rate limiter")
	app.rateLimiter = core_redis_rate_limit.New(app.infrastructure.Redis)

	logger.Debug("initializing app workers")
	if err := app.initWorkers(ctx); err != nil {
		app.infrastructure.Close(ctx)
		return nil, fmt.Errorf("init workers: %w", err)
	}

	logger.Debug("initializing app features")
	if err := app.initFeatures(); err != nil {
		app.infrastructure.Close(ctx)
		return nil, fmt.Errorf("init features: %w", err)
	}

	logger.Debug("initializing app websocket server")
	if err := app.initWebsocket(ctx); err != nil {
		app.infrastructure.Close(ctx)
		return nil, fmt.Errorf("init websocket server: %w", err)
	}

	logger.Debug("initializing app http server")
	if err := app.initHttp(); err != nil {
		app.infrastructure.Close(ctx)
		return nil, fmt.Errorf("init http server: %w", err)
	}

	return app, nil
}

func (a *App) Run(ctx context.Context) error {
	err := a.httpServer.Run(ctx)

	shutdownCtx, cancel := context.WithTimeout(context.Background(), a.cfg.ShutdownTimeout)
	defer cancel()

	a.logger.Warn("shutting down application")
	if shutdownErr := a.httpServer.Shutdown(shutdownCtx); shutdownErr != nil {
		a.logger.Error("failed to shutdown HTTP server", zap.Error(shutdownErr))
		if err == nil {
			err = shutdownErr
		}
	}

	if wsErr := a.wsServer.Shutdown(shutdownCtx); wsErr != nil {
		a.logger.Error("failed to shutdown websocket server", zap.Error(wsErr))
		if err == nil {
			err = wsErr
		}
	}

	if a.workerCancel != nil {
		a.workerCancel()
	}

	if a.workerDone != nil {
		select {
		case <-a.workerDone:
		case <-shutdownCtx.Done():
			if err == nil {
				err = shutdownCtx.Err()
			}
		}
	}

	a.infrastructure.Close(shutdownCtx)

	return err
}
