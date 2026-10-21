package core_rate_limit

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	ErrInvalidPolicy = errors.New("invalid rate limit policy")
	ErrInvalidKey    = errors.New("invalid rate limit key")
)

type Policy struct {
	Name   string
	Limit  int64
	Window time.Duration
}

func (p Policy) Validate() error {
	if p.Name == "" {
		return fmt.Errorf("%w: name must not be empty", ErrInvalidPolicy)
	}
	if p.Limit <= 0 {
		return fmt.Errorf("%w: limit must be greater than zero", ErrInvalidPolicy)
	}
	if p.Window < time.Millisecond || p.Window.Milliseconds() <= 0 {
		return fmt.Errorf("%w: window must be at least 1ms", ErrInvalidPolicy)
	}
	return nil
}

type Decision struct {
	Allowed    bool
	Limit      int64
	Remaining  int64
	RetryAfter time.Duration
}

type Limiter interface {
	Allow(ctx context.Context, key string, policy Policy) (Decision, error)
}

type FailureMode string

const (
	FailureClosed FailureMode = "closed"
	FailureOpen   FailureMode = "open"
)

func (m FailureMode) Validate() error {
	if m != FailureClosed && m != FailureOpen {
		return fmt.Errorf("invalid rate limit failure mode %q: use %q or %q", m, FailureClosed, FailureOpen)
	}

	return nil
}

type Outcome string

const (
	OutcomeAllowed Outcome = "allowed"
	OutcomeLimited Outcome = "limited"
	OutcomeError   Outcome = "error"
	OutcomeBypass  Outcome = "bypass"
)

type Observation struct {
	Policy   string
	Outcome  Outcome
	Duration time.Duration
}

type Observer interface {
	Observe(Observation)
}

type NoopObserver struct{}

func (NoopObserver) Observe(Observation) {}

type CounterObserver struct {
	mu       sync.RWMutex
	counters map[string]uint64
	total    map[string]time.Duration
}

func NewCounterObserver() *CounterObserver {
	return &CounterObserver{
		counters: make(map[string]uint64),
		total:    make(map[string]time.Duration),
	}
}

func (o *CounterObserver) Observe(observation Observation) {
	if o == nil {
		return
	}

	key := observation.Policy + "|" + string(observation.Outcome)

	o.mu.Lock()

	o.counters[key]++

	o.total[key] += observation.Duration

	o.mu.Unlock()
}

type MetricSnapshot struct {
	Policy        string
	Outcome       Outcome
	Count         uint64
	TotalDuration time.Duration
}

func (o *CounterObserver) Snapshot() []MetricSnapshot {
	if o == nil {
		return nil
	}
	o.mu.RLock()
	defer o.mu.RUnlock()

	out := make([]MetricSnapshot, 0, len(o.counters))
	for key, count := range o.counters {
		separator := -1
		for i := 0; i < len(key); i++ {
			if key[i] == '|' {
				separator = i
				break
			}
		}

		if separator < 0 {
			continue
		}

		out = append(out, MetricSnapshot{
			Policy:        key[:separator],
			Outcome:       Outcome(key[separator+1:]),
			Count:         count,
			TotalDuration: o.total[key],
		})
	}

	return out
}

type Options struct {
	FailureMode FailureMode
	Observer    Observer
}
