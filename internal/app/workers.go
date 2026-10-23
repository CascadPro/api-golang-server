package app

import (
	"context"
	"fmt"
	"sync"

	core_postgres_outbox "github.com/CascadePro/api-golang-server/internal/core/infrastructure/postgres/outbox"
	media_postgres_repository "github.com/CascadePro/api-golang-server/internal/features/media/repository/postgres"
	worker_cron "github.com/CascadePro/api-golang-server/internal/workers/cron"
	cron_idempotency_handler "github.com/CascadePro/api-golang-server/internal/workers/cron/handlers/idempotency"
	cron_outbox_handler "github.com/CascadePro/api-golang-server/internal/workers/cron/handlers/outbox"
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

	outboxWorkerHandler := worker_outbox.NewHandler(outboxMediaHandler, realtimeMediaHandler)
	outboxWorker, err := worker_outbox.New(repository, outboxWorkerHandler, a.logger, worker_outbox.NewConfigMust())
	if err != nil {
		return fmt.Errorf("create outbox worker: %w", err)
	}

	cronIdempotencyHandler, err := cron_idempotency_handler.New(a.infrastructure.Idempotency, a.logger, cron_idempotency_handler.NewConfigMust())
	if err != nil {
		return fmt.Errorf("create idempotency cron handler: %w", err)
	}

	cronOutboxHandler, err := cron_outbox_handler.New(repository, a.logger, cron_outbox_handler.NewConfigMust())
	if err != nil {
		return fmt.Errorf("create outbox cron handler: %w", err)
	}

	cronWorker, err := worker_cron.New(a.logger, cronIdempotencyHandler, cronOutboxHandler)
	if err != nil {
		return fmt.Errorf("create cron worker: %w", err)
	}

	workerCtx, cancel := context.WithCancel(ctx)

	a.workerCancel = cancel
	a.workerDone = make(chan struct{})

	var workers sync.WaitGroup
	workers.Add(2)

	a.logger.Debug("starting outbox worker", zap.String("worker", "outbox"))
	go func() {
		defer workers.Done()
		outboxWorker.Run(workerCtx)
	}()

	a.logger.Debug("starting cron worker", zap.String("worker", "cron"))
	go func() {
		defer workers.Done()
		cronWorker.Run(workerCtx)
	}()

	go func() {
		workers.Wait()
		close(a.workerDone)
	}()

	return nil
}
