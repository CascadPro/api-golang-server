package worker_outbox

import (
	"context"
	"fmt"
	"time"

	core_postgres_outbox "github.com/CascadePro/api-golang-server/internal/core/infrastructure/postgres/outbox"
	core_logger "github.com/CascadePro/api-golang-server/internal/core/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Worker struct {
	repository core_postgres_outbox.RepositoryMethods
	handler    HandlerMethods
	logger     *core_logger.Logger

	workerID  uuid.UUID
	interval  time.Duration
	batchSize int
	lockTTL   time.Duration
}

func New(
	repository core_postgres_outbox.RepositoryMethods,
	handler HandlerMethods,
	logger *core_logger.Logger,
	cfg Config,
) (*Worker, error) {
	workerID, err := uuid.NewV7()
	if err != nil {
		return nil, fmt.Errorf("generate ID: %w", err)
	}

	return &Worker{
		repository: repository,
		handler:    handler,
		logger:     logger,

		workerID:  workerID,
		interval:  cfg.Interval,
		batchSize: cfg.BatchSize,
		lockTTL:   cfg.LockTTL,
	}, nil
}

func (w *Worker) Run(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	w.process(ctx)

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
	events, err := w.repository.GetPendingEvents(ctx, w.workerID, w.batchSize, w.lockTTL)
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

			if err := w.repository.MarkEventFailed(ctx, event.ID, w.workerID, err); err != nil {
				w.logger.Error(
					"failed to mark outbox event failed",
					zap.String("event_id", event.ID.String()),
					zap.Error(err),
				)
			}

			continue
		}

		if err := w.repository.MarkEventProcessed(ctx, event.ID, w.workerID); err != nil {
			w.logger.Error(
				"failed to mark outbox event processed",
				zap.String("event_id", event.ID.String()),
				zap.Error(err),
			)
		}
	}
}
