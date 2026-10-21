package core_redis_rate_limit

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

// Set RATE_LIMIT_REDIS_ADDR to run these integration tests against a real Redis instance.
// Optional settings: RATE_LIMIT_REDIS_PASSWORD and RATE_LIMIT_REDIS_DB.
func integrationRedis(t *testing.T) *redis.Client {
	t.Helper()
	addr := os.Getenv("RATE_LIMIT_REDIS_ADDR")
	if addr == "" {
		t.Skip("RATE_LIMIT_REDIS_ADDR is not set")
	}
	db, _ := strconv.Atoi(os.Getenv("RATE_LIMIT_REDIS_DB"))
	client := redis.NewClient(&redis.Options{
		Addr: addr, Password: os.Getenv("RATE_LIMIT_REDIS_PASSWORD"), DB: db,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		_ = client.Close()
		t.Fatalf("connect to Redis: %v", err)
	}
	t.Cleanup(func() { _ = client.Close() })
	return client
}

func TestLuaScriptConcurrentRequestsNeverExceedLimit(t *testing.T) {
	client := integrationRedis(t)
	ctx := context.Background()
	key := fmt.Sprintf("test:rate-limit:%d", time.Now().UnixNano())
	t.Cleanup(func() { _ = client.Del(context.Background(), key).Err() })

	const limit = int64(7)
	const requests = 64
	var allowed atomic.Int64
	var wg sync.WaitGroup
	for i := 0; i < requests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			result, err := client.Eval(ctx, luaScript, []string{key}, limit, int64(5000)).Result()
			if err != nil {
				t.Errorf("run Lua script: %v", err)
				return
			}
			values, ok := result.([]interface{})
			if !ok || len(values) != 3 {
				t.Errorf("unexpected script response: %#v", result)
				return
			}
			allowedValue, ok := values[0].(int64)
			if !ok {
				t.Errorf("unexpected allowed value: %#v", values[0])
				return
			}
			if allowedValue == 1 {
				allowed.Add(1)
			}
		}()
	}
	wg.Wait()
	if got := allowed.Load(); got != limit {
		t.Fatalf("expected exactly %d requests allowed, got %d", limit, got)
	}
}

func TestLuaScriptRepairsMissingTTL(t *testing.T) {
	client := integrationRedis(t)
	ctx := context.Background()
	key := fmt.Sprintf("test:rate-limit-no-ttl:%d", time.Now().UnixNano())
	t.Cleanup(func() { _ = client.Del(context.Background(), key).Err() })

	if err := client.Set(ctx, key, "5", 0).Err(); err != nil {
		t.Fatal(err)
	}
	result, err := client.Eval(ctx, luaScript, []string{key}, int64(5), int64(2000)).Result()
	if err != nil {
		t.Fatal(err)
	}
	values, ok := result.([]interface{})
	if !ok || len(values) != 3 {
		t.Fatalf("unexpected script response: %#v", result)
	}
	if values[0] != int64(0) {
		t.Fatalf("expected request to be limited, got %#v", values[0])
	}
	ttl, err := client.PTTL(ctx, key).Result()
	if err != nil {
		t.Fatal(err)
	}
	if ttl <= 0 {
		t.Fatalf("expected repaired positive TTL, got %s", ttl)
	}
}
