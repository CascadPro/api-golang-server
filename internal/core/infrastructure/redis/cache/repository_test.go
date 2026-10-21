package core_redis_cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"path"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	core_redis_pool "github.com/CascadePro/api-golang-server/internal/core/infrastructure/redis/pool"
)

type fakeRedis struct {
	mu        sync.Mutex
	values    map[string]string
	lastKey   string
	lastTTL   time.Duration
	timeout   time.Duration
	getErr    error
	setErr    error
	deleteErr error
}

func newFakeRedis() *fakeRedis {
	return &fakeRedis{values: make(map[string]string), timeout: time.Second}
}

func (f *fakeRedis) Get(_ context.Context, key string) (string, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.lastKey = key
	if f.getErr != nil {
		return "", f.getErr
	}
	value, ok := f.values[key]
	if !ok {
		return "", core_redis_pool.ErrNoValue
	}
	return value, nil
}

func (f *fakeRedis) MGet(_ context.Context, keys ...string) ([]any, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	values := make([]any, len(keys))
	if f.getErr != nil {
		return nil, f.getErr
	}
	for i, key := range keys {
		if value, ok := f.values[key]; ok {
			values[i] = value
		}
	}
	return values, nil
}

func (f *fakeRedis) Set(_ context.Context, key string, value any, ttl time.Duration) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.lastKey = key
	f.lastTTL = ttl
	if f.setErr != nil {
		return f.setErr
	}
	encoded, err := encodeFakeValue(value)
	if err != nil {
		return err
	}
	f.values[key] = encoded
	return nil
}

func (f *fakeRedis) SetMany(_ context.Context, values map[string]any, ttl time.Duration) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.setErr != nil {
		return f.setErr
	}
	f.lastTTL = ttl
	for key, value := range values {
		encoded, err := encodeFakeValue(value)
		if err != nil {
			return err
		}
		f.values[key] = encoded
		f.lastKey = key
	}
	return nil
}

func (f *fakeRedis) Del(_ context.Context, keys ...string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.deleteErr != nil {
		return f.deleteErr
	}
	for _, key := range keys {
		f.lastKey = key
		delete(f.values, key)
	}
	return nil
}

func (f *fakeRedis) Scan(_ context.Context, _ uint64, match string, _ int64) ([]string, uint64, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	keys := make([]string, 0)
	for key := range f.values {
		matched, err := path.Match(match, key)
		if err != nil {
			return nil, 0, err
		}
		if matched {
			keys = append(keys, key)
		}
	}
	return keys, 0, nil
}

func (f *fakeRedis) OpTimeout() time.Duration { return f.timeout }

func encodeFakeValue(value any) (string, error) {
	switch v := value.(type) {
	case []byte:
		return string(v), nil
	case string:
		return v, nil
	default:
		encoded, err := json.Marshal(v)
		return string(encoded), err
	}
}

func newTestCache(t *testing.T, pool *fakeRedis, namespace string, ttl time.Duration) *RedisCache {
	t.Helper()
	cache, err := NewRedisCache(pool, namespace, ttl)
	if err != nil {
		t.Fatal(err)
	}
	return cache
}

func TestRedisCacheSetAndGetJSON(t *testing.T) {
	pool := newFakeRedis()
	cache := newTestCache(t, pool, "ipinfo", 24*time.Hour)
	type payload struct {
		IP   string `json:"ip"`
		City string `json:"city"`
	}
	want := payload{IP: "8.8.8.8", City: "Mountain View"}
	if err := cache.Set(context.Background(), "8.8.8.8", want, 0); err != nil {
		t.Fatal(err)
	}
	if pool.lastKey != "cascade__cache:ipinfo:8.8.8.8" {
		t.Fatalf("unexpected Redis key: %q", pool.lastKey)
	}
	if pool.lastTTL != 24*time.Hour {
		t.Fatalf("unexpected TTL: %s", pool.lastTTL)
	}
	var got payload
	if err := cache.Get(context.Background(), "8.8.8.8", &got); err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("got %+v, want %+v", got, want)
	}
}

func TestRedisCacheSetUsesExplicitTTL(t *testing.T) {
	pool := newFakeRedis()
	cache := newTestCache(t, pool, "test", time.Hour)
	if err := cache.Set(context.Background(), "key", "value", 5*time.Minute); err != nil {
		t.Fatal(err)
	}
	if pool.lastTTL != 5*time.Minute {
		t.Fatalf("got TTL %s, want %s", pool.lastTTL, 5*time.Minute)
	}
}

func TestRedisCacheGetMapsMiss(t *testing.T) {
	cache := newTestCache(t, newFakeRedis(), "test", time.Minute)
	var got map[string]any
	err := cache.Get(context.Background(), "missing", &got)
	if !errors.Is(err, core_errors.ErrNotFound) {
		t.Fatalf("expected core ErrNotFound, got %v", err)
	}
}

func TestRedisCacheDeleteAndBatchOperations(t *testing.T) {
	pool := newFakeRedis()
	cache := newTestCache(t, pool, "test", time.Minute)
	ctx := context.Background()
	if err := cache.SetMany(ctx, map[string]any{"a": map[string]int{"v": 1}, "b": map[string]int{"v": 2}}, 0); err != nil {
		t.Fatal(err)
	}
	var a, b map[string]int
	failures, err := cache.GetMany(ctx, map[string]any{"a": &a, "b": &b, "missing": &map[string]int{}})
	if err != nil {
		t.Fatal(err)
	}
	if a["v"] != 1 || b["v"] != 2 {
		t.Fatalf("unexpected values: a=%v b=%v", a, b)
	}
	if !errors.Is(failures["missing"], core_errors.ErrNotFound) {
		t.Fatalf("expected missing-key error, got %v", failures["missing"])
	}
	if err := cache.DeleteMany(ctx, "a", "b"); err != nil {
		t.Fatal(err)
	}
	if len(pool.values) != 0 {
		t.Fatalf("expected no keys, got %v", pool.values)
	}
}

func TestRedisCacheInvalidateOnlyNamespace(t *testing.T) {
	pool := newFakeRedis()
	cache := newTestCache(t, pool, "ipinfo", time.Minute)
	pool.values["cascade__cache:ipinfo:a"] = `{"v":1}`
	pool.values["cascade__cache:ipinfo:b"] = `{"v":2}`
	pool.values["cascade__cache:other:c"] = `{"v":3}`
	if err := cache.Invalidate(context.Background()); err != nil {
		t.Fatal(err)
	}
	if _, ok := pool.values["cascade__cache:ipinfo:a"]; ok {
		t.Fatal("ipinfo key a was not invalidated")
	}
	if _, ok := pool.values["cascade__cache:ipinfo:b"]; ok {
		t.Fatal("ipinfo key b was not invalidated")
	}
	if _, ok := pool.values["cascade__cache:other:c"]; !ok {
		t.Fatal("other namespace must remain untouched")
	}
}

func TestRedisCacheGetOrLoadCoalescesConcurrentMisses(t *testing.T) {
	cache := newTestCache(t, newFakeRedis(), "test", time.Minute)
	var loads atomic.Int32
	const workers = 12
	var wg sync.WaitGroup
	errs := make(chan error, workers)
	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var dest struct {
				Value string `json:"value"`
			}
			err := cache.GetOrLoad(context.Background(), "shared", &dest, time.Minute, func(context.Context) (any, error) {
				loads.Add(1)
				time.Sleep(25 * time.Millisecond)
				return struct {
					Value string `json:"value"`
				}{Value: "loaded"}, nil
			})
			if err == nil && dest.Value != "loaded" {
				err = fmt.Errorf("unexpected loaded value: %q", dest.Value)
			}
			errs <- err
		}()
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		if err != nil {
			t.Fatal(err)
		}
	}
	if got := loads.Load(); got != 1 {
		t.Fatalf("loader called %d times, want 1", got)
	}
}

func TestRedisCacheRejectsInvalidConfigurationAndTTL(t *testing.T) {
	pool := newFakeRedis()
	if _, err := NewRedisCache(pool, " ", time.Minute); !errors.Is(err, ErrInvalidNamespace) {
		t.Fatalf("expected ErrInvalidNamespace, got %v", err)
	}
	if _, err := NewRedisCache(pool, "unsafe:*", time.Minute); !errors.Is(err, ErrInvalidNamespace) {
		t.Fatalf("expected wildcard namespace to be rejected, got %v", err)
	}
	if _, err := NewRedisCache(pool, "test", -time.Second); !errors.Is(err, ErrInvalidTTL) {
		t.Fatalf("expected ErrInvalidTTL, got %v", err)
	}
	cache := newTestCache(t, pool, "test", time.Minute)
	if err := cache.Set(context.Background(), "key", "value", -time.Second); !errors.Is(err, ErrInvalidTTL) {
		t.Fatalf("expected ErrInvalidTTL, got %v", err)
	}
}

func TestRedisCacheRejectsNilContextAndDestination(t *testing.T) {
	cache := newTestCache(t, newFakeRedis(), "test", time.Minute)
	if err := cache.Get(nil, "key", new(string)); !errors.Is(err, ErrNilContext) {
		t.Fatalf("expected ErrNilContext, got %v", err)
	}
	var dest *string
	if err := cache.Get(context.Background(), "key", dest); !errors.Is(err, ErrNilDestination) {
		t.Fatalf("expected ErrNilDestination, got %v", err)
	}
}
