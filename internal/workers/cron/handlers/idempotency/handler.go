package cron_idempotency_handler

import (
	"context"
	"fmt"
	"time"

	core_logger "github.com/CascadePro/api-golang-server/internal/core/logger"
	"go.uber.org/zap"
)

const handlerName = "idempotency_cleanup"

type CleanupService interface {
	CleanupExpired(context.Context) (int64, error)
}

type Handler struct {
	service CleanupService
	logger  *core_logger.Logger

	interval time.Duration
	timeout  time.Duration
}

func New(service CleanupService, logger *core_logger.Logger, config Config) (*Handler, error) {
	if service == nil {
		return nil, fmt.Errorf("idempotency cleanup service is nil")
	}

	if logger == nil {
		return nil, fmt.Errorf("idempotency handler logger is nil")
	}

	if config.Interval <= 0 {
		return nil, fmt.Errorf("idempotency handler interval must be positive")
	}

	if config.Timeout <= 0 {
		return nil, fmt.Errorf("idempotency handler timeout must be positive")
	}

	return &Handler{
		service:  service,
		logger:   logger,
		interval: config.Interval,
		timeout:  config.Timeout,
	}, nil
}

func (h *Handler) Name() string {
	return handlerName
}

func (h *Handler) Interval() time.Duration {
	return h.interval
}

func (h *Handler) Timeout() time.Duration {
	return h.timeout
}

func (h *Handler) Handle(ctx context.Context) error {
	deleted, err := h.service.CleanupExpired(ctx)
	if err != nil {
		return fmt.Errorf("cleanup expired idempotency records: %w", err)
	}

	h.logger.Debug("cleaned up expired idempotency records", zap.Int64("deleted_count", deleted))

	return nil
}
