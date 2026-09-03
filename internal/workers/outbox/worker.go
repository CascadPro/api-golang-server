package worker_outbox

import (
	"context"
	"time"

	core_postgres_outbox "github.com/CascadePro/api-golang-server/internal/core/infrastructure/postgres/outbox"
	core_logger "github.com/CascadePro/api-golang-server/internal/core/logger"
	"go.uber.org/zap"
)

type Worker struct {
	repository core_postgres_outbox.RepositoryMethods
	handler    HandlerMethods
	logger     *core_logger.Logger

	interval  time.Duration
	batchSize int
}

func New(repository core_postgres_outbox.RepositoryMethods, handler HandlerMethods, logger *core_logger.Logger) *Worker {
	return &Worker{
		repository: repository,
		handler:    handler,
		logger:     logger,

		interval:  5 * time.Second,
		batchSize: 50,
	}
}

func (w *Worker) Run(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return

		case <-ticker.C:
			w.process(ctx)
		}
	}
}

func (w *Worker) process(ctx context.Context) {
	events, err := w.repository.GetPendingEvents(ctx, w.batchSize)
	if err != nil {
		w.logger.Error("failed to get pending outbox events", zap.Error(err))

		return
	}

	for _, event := range events {
		if err := w.handler.Handle(ctx, event); err != nil {
			w.logger.Error(
				"failed to process outbox event",
				zap.String("event_id", event.ID.String()),
				zap.String("event_type", string(event.Type)),
				zap.Error(err),
			)

			if event.Attempts >= 3 {
				_ = w.repository.MarkEventFailed(ctx, event.ID, err)
			} else {
				_ = w.repository.IncrementAttempt(ctx, event.ID)
			}

			continue
		}

		if err := w.repository.MarkEventProcessed(ctx, event.ID); err != nil {
			w.logger.Error(
				"failed to mark outbox event processed",
				zap.String("event_id", event.ID.String()),
				zap.Error(err),
			)
		}
	}
}
