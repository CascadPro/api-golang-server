package worker_cron

import (
	"context"
	"fmt"
	"sync"
	"time"

	core_logger "github.com/CascadePro/api-golang-server/internal/core/logger"
	"go.uber.org/zap"
)

// Handler is one independently scheduled CRON operation.
// Each handler owns its schedule and timeout, so one slow operation does not
// block other handlers from running on schedule.
type Handler interface {
	Name() string

	Interval() time.Duration

	Timeout() time.Duration

	Handle(context.Context) error
}

type Worker struct {
	handlers []Handler
	logger   *core_logger.Logger
}

func New(logger *core_logger.Logger, handlers ...Handler) (*Worker, error) {
	if logger == nil {
		return nil, fmt.Errorf("cron worker logger is nil")
	}
	if len(handlers) == 0 {
		return nil, fmt.Errorf("cron worker requires at least one handler")
	}

	seen := make(map[string]struct{}, len(handlers))
	for _, handler := range handlers {
		if handler == nil {
			return nil, fmt.Errorf("cron worker handler is nil")
		}

		if handler.Name() == "" {
			return nil, fmt.Errorf("cron worker handler name is empty")
		}

		if _, exists := seen[handler.Name()]; exists {
			return nil, fmt.Errorf("duplicate cron handler name %q", handler.Name())
		}

		seen[handler.Name()] = struct{}{}

		if handler.Interval() <= 0 {
			return nil, fmt.Errorf("cron handler %q interval must be positive", handler.Name())
		}

		if handler.Timeout() <= 0 {
			return nil, fmt.Errorf("cron handler %q timeout must be positive", handler.Name())
		}
	}

	return &Worker{
		handlers: append([]Handler(nil), handlers...),
		logger:   logger,
	}, nil
}

// Run starts each handler in its own loop. Every handler runs once immediately
// and then repeats at its own interval. A handler cannot overlap with itself.
func (w *Worker) Run(ctx context.Context) {
	var wg sync.WaitGroup
	wg.Add(len(w.handlers))

	for _, handler := range w.handlers {
		go func() {
			defer wg.Done()
			w.runHandler(ctx, handler)
		}()
	}

	wg.Wait()
}

func (w *Worker) runHandler(ctx context.Context, handler Handler) {
	ticker := time.NewTicker(handler.Interval())
	defer ticker.Stop()

	w.logger.Debug("starting cron handler", zap.String("handler", handler.Name()))
	w.execute(ctx, handler)

	for {
		select {
		case <-ctx.Done():
			w.logger.Debug("stopping cron handler", zap.String("handler", handler.Name()))
			return

		case <-ticker.C:
			w.execute(ctx, handler)
		}
	}
}

func (w *Worker) execute(ctx context.Context, handler Handler) {
	if ctx.Err() != nil {
		return
	}

	startedAt := time.Now()
	handlerCtx, cancel := context.WithTimeout(ctx, handler.Timeout())
	defer cancel()

	if err := handler.Handle(handlerCtx); err != nil {
		w.logger.Error(
			"cron handler failed",
			zap.String("handler", handler.Name()),
			zap.Duration("duration", time.Since(startedAt)),
			zap.Error(err),
		)
		return
	}

	w.logger.Debug(
		"cron handler completed",
		zap.String("handler", handler.Name()),
		zap.Duration("duration", time.Since(startedAt)),
	)
}
