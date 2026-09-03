package app

import (
	"context"
	"fmt"

	core_postgres_outbox "github.com/CascadePro/api-golang-server/internal/core/infrastructure/postgres/outbox"
	worker_outbox "github.com/CascadePro/api-golang-server/internal/workers/outbox"
	outbox_media_handler "github.com/CascadePro/api-golang-server/internal/workers/outbox/handler/media"
	"go.uber.org/zap"
)

func (a *App) initWorkers(ctx context.Context) error {
	repository := core_postgres_outbox.NewRepository(a.infrastructure.Postgres)

	a.logger.Debug("initializing outbox handler", zap.String("handler", "media"))
	outboxMediaHandler := outbox_media_handler.New(a.infrastructure.S3)

	outboxWorkerHandler := worker_outbox.NewHandler(outboxMediaHandler)

	outboxWorker, err := worker_outbox.New(repository, outboxWorkerHandler, a.logger)
	if err != nil {
		return fmt.Errorf("create outbox worker: %w", err)
	}

	a.logger.Debug("starting outbox worker", zap.String("worker", "outbox"))
	go outboxWorker.Run(ctx)

	return nil
}
