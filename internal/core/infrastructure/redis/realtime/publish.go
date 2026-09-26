package core_redis_realtime

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
)

func (p *Publisher) Publish(ctx context.Context, event domain.RealtimeEvent) error {
	ctx, cancel := context.WithTimeout(ctx, p.pool.OpTimeout())
	defer cancel()

	if err := event.Validate(); err != nil {
		return fmt.Errorf("validate realtime event: %w", err)
	}

	payload, err := json.Marshal(domainToModel(event))
	if err != nil {
		return fmt.Errorf("marshal realtime event: %w", err)
	}

	if err := p.pool.Publish(ctx, ChannelKey, payload); err != nil {
		return fmt.Errorf("publish realtime event: %w", err)
	}

	return nil
}
