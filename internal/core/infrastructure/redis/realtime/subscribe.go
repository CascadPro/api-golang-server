package core_redis_realtime

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
)

func (p *Publisher) Subscribe(ctx context.Context, handler func(domain.RealtimeEvent)) error {
	ctx, cancel := context.WithTimeout(ctx, p.pool.OpTimeout())
	defer cancel()

	pubsub := p.pool.Subscribe(ctx, ChannelKey)

	if _, err := pubsub.Receive(ctx); err != nil {
		_ = pubsub.Close()

		return fmt.Errorf("subscribe realtime channel: %w", err)
	}

	go func() {
		defer pubsub.Close()

		for {
			select {
			case <-ctx.Done():
				return

			case message, ok := <-pubsub.Channel():
				if !ok {
					return
				}

				var model EventModel
				if err := json.Unmarshal([]byte(message.Payload), &model); err != nil {
					continue
				}

				handler(modelToDomain(model))
			}
		}
	}()

	return nil
}
