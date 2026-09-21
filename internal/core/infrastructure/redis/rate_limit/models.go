package core_redis_rate_limit

import (
	"time"
)

type Result struct {
	Allowed    bool
	Remaining  int64
	RetryAfter time.Duration
}
