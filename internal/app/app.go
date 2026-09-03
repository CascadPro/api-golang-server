package app

import (
	"context"
	"fmt"

	core_config "github.com/CascadePro/api-golang-server/internal/core/config"
	core_logger "github.com/CascadePro/api-golang-server/internal/core/logger"
	core_http_server "github.com/CascadePro/api-golang-server/internal/core/transport/http/server"
	"go.uber.org/zap"
)

type App struct {
	cfg    *core_config.Config
	logger *core_logger.Logger

	infrastructure *Infrastructure
	features       *Features
	httpServer     *core_http_server.HttpServer
}

func New(
	ctx context.Context,
	cfg *core_config.Config,
	logger *core_logger.Logger,
) (*App, error) {
	logger.Info("application environment mode", zap.Any("mode", cfg.EnvMode))

	logger.Info("application connection", zap.Any("connection", cfg.Connection))

	logger.Info("application time zone", zap.Any("zone", cfg.TimeZone))

	app := &App{
		cfg:    cfg,
		logger: logger,
	}

	logger.Debug("initializing app infrastructure")
	if err := app.initInfrastructure(ctx); err != nil {
		return nil, fmt.Errorf("init infrastructure: %w", err)
	}

	logger.Debug("initializing app workers")
	if err := app.initWorkers(ctx); err != nil {
		return nil, fmt.Errorf("init workers: %w", err)
	}

	logger.Debug("initializing app features")
	if err := app.initFeatures(); err != nil {
		app.infrastructure.Close(ctx)

		return nil, fmt.Errorf("init features: %w", err)
	}

	logger.Debug("initializing app http server")
	if err := app.initHttp(); err != nil {
		app.infrastructure.Close(ctx)

		return nil, fmt.Errorf("init http server: %w", err)
	}

	return app, nil
}

func (a *App) Run(ctx context.Context) error {
	return a.httpServer.Run(ctx)
}
