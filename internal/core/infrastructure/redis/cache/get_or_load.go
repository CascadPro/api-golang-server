package core_redis_cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
)

// GetOrLoad implements cache-aside with process-local request coalescing. One
// goroutine loads a missing key; concurrent callers wait for the same result.
// This prevents duplicate loads inside one process, not across multiple replicas.
func (c *RedisCache) GetOrLoad(ctx context.Context, key string, dest any, ttl time.Duration, loader func(context.Context) (any, error)) error {
	if err := validateKey(key); err != nil {
		return err
	}
	if isNilValue(dest) {
		return fmt.Errorf("cache get-or-load destination: %w", ErrNilDestination)
	}
	if loader == nil {
		return ErrNilLoader
	}
	ttl, err := c.effectiveTTL(ttl)
	if err != nil {
		return err
	}

	if err := c.Get(ctx, key, dest); err == nil {
		return nil
	} else if !errors.Is(err, core_errors.ErrNotFound) {
		return err
	}

	flightKey := c.redisKey(key)
	c.flightsMu.Lock()

	if active, ok := c.flights[flightKey]; ok {
		c.flightsMu.Unlock()

		select {
		case <-ctx.Done():
			return ctx.Err()

		case <-active.done:
			if active.err != nil {
				return active.err
			}

			if err := json.Unmarshal(active.data, dest); err != nil {
				return fmt.Errorf("decode coalesced cache result: %w", err)
			}

			return nil
		}
	}

	active := &loadCall{done: make(chan struct{})}
	c.flights[flightKey] = active
	c.flightsMu.Unlock()

	loaded, loadErr := loader(ctx)
	if loadErr == nil {
		active.data, loadErr = json.Marshal(loaded)
		if loadErr != nil {
			loadErr = fmt.Errorf("encode loaded cache value: %w", loadErr)
		}
	}
	if loadErr == nil {
		loadErr = c.Set(ctx, key, loaded, ttl)
	}
	if loadErr == nil {
		if err := json.Unmarshal(active.data, dest); err != nil {
			loadErr = fmt.Errorf("decode loaded cache value: %w", err)
		}
	}

	c.flightsMu.Lock()

	active.err = loadErr

	delete(c.flights, flightKey)

	close(active.done)

	c.flightsMu.Unlock()

	if loadErr != nil {
		return loadErr
	}

	return nil
}
