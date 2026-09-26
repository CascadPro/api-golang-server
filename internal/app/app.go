package app

import (
	"context"
	"fmt"

	core_config "github.com/CascadePro/api-golang-server/internal/core/config"
	core_i18n "github.com/CascadePro/api-golang-server/internal/core/i18n"
	core_logger "github.com/CascadePro/api-golang-server/internal/core/logger"
	core_http_server "github.com/CascadePro/api-golang-server/internal/core/transport/http/server"
	core_ws_server "github.com/CascadePro/api-golang-server/internal/core/transport/ws/server"
	"go.uber.org/zap"
)

type App struct {
	cfg    *core_config.Config
	logger *core_logger.Logger

	infrastructure *Infrastructure
	features       *Features

	httpServer *core_http_server.HttpServer
	wsServer   *core_ws_server.Server

	workerCancel context.CancelFunc
	workerDone   chan struct{}
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

	logger.Debug("initializing i18n localization")
	if err := core_i18n.Init(); err != nil {
		return nil, fmt.Errorf("init localization: %w", err)
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

	a.logger.Warn("Shutting down the server...")

	if wsErr := a.wsServer.Shutdown(shutdownCtx); wsErr != nil {
		a.logger.Error("failed to shutdown websocket server", zap.Error(wsErr))

		if err == nil {
			err = fmt.Errorf("shutdown websocket server: %w", wsErr)
		}
	}

	a.logger.Warn("Shut down websocket server")

	if a.workerCancel != nil {
		a.workerCancel()
	}

	if a.workerDone != nil {
		select {
		case <-a.workerDone:
		case <-shutdownCtx.Done():
			if err == nil {
				err = fmt.Errorf("outbox worker shutdown: %w", shutdownCtx.Err())
			}
		}
	}

	a.logger.Warn("Stopped app workers")

	a.infrastructure.Close(shutdownCtx)

	a.logger.Warn("Closed app infrastructure")

	return err
}
