package app

import (
	"context"
	"fmt"

	core_postgres_outbox "github.com/CascadePro/api-golang-server/internal/core/infrastructure/postgres/outbox"
	media_postgres_repository "github.com/CascadePro/api-golang-server/internal/features/media/repository/postgres"
	worker_outbox "github.com/CascadePro/api-golang-server/internal/workers/outbox"
	outbox_media_handler "github.com/CascadePro/api-golang-server/internal/workers/outbox/handler/media"
	outbox_realtime_handler "github.com/CascadePro/api-golang-server/internal/workers/outbox/handler/realtime"
	"go.uber.org/zap"
)

func (a *App) initWorkers(ctx context.Context) error {
	repository := core_postgres_outbox.NewRepository(a.infrastructure.Postgres)

	mediaRepository := media_postgres_repository.NewRepository(a.infrastructure.Postgres)

	a.logger.Debug("initializing outbox handler", zap.String("handler", "media"))
	outboxMediaHandler := outbox_media_handler.New(a.infrastructure.S3, mediaRepository)

	a.logger.Debug("initializing outbox handler", zap.String("handler", "realtime"))
	realtimeMediaHandler := outbox_realtime_handler.New(a.infrastructure.Publisher)

	outboxWorkerHandler := worker_outbox.NewHandler(
		outboxMediaHandler,
		realtimeMediaHandler,
	)

	outboxWorker, err := worker_outbox.New(repository, outboxWorkerHandler, a.logger, worker_outbox.NewConfigMust())
	if err != nil {
		return fmt.Errorf("create outbox worker: %w", err)
	}

	a.logger.Debug("starting outbox worker", zap.String("worker", "outbox"))

	workerCtx, cancel := context.WithCancel(ctx)

	a.workerCancel = cancel
	a.workerDone = make(chan struct{})

	go func() {
		defer close(a.workerDone)
		outboxWorker.Run(workerCtx)
	}()

	return nil
}
