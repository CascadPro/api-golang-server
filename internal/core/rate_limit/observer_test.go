package core_rate_limit

import (
	"testing"
	"time"
)

func TestCounterObserverAggregatesByPolicyAndOutcome(t *testing.T) {
	observer := NewCounterObserver()
	observer.Observe(Observation{Policy: "http.auth", Outcome: OutcomeLimited, Duration: 10 * time.Millisecond})
	observer.Observe(Observation{Policy: "http.auth", Outcome: OutcomeLimited, Duration: 20 * time.Millisecond})
	observer.Observe(Observation{Policy: "http.auth", Outcome: OutcomeAllowed, Duration: 5 * time.Millisecond})

	snapshot := observer.Snapshot()
	if len(snapshot) != 2 {
		t.Fatalf("expected 2 metric groups, got %d", len(snapshot))
	}
	var limitedFound bool
	for _, item := range snapshot {
		if item.Policy == "http.auth" && item.Outcome == OutcomeLimited {
			limitedFound = true
			if item.Count != 2 || item.TotalDuration != 30*time.Millisecond {
				t.Fatalf("unexpected aggregate: %+v", item)
			}
		}
	}
	if !limitedFound {
		t.Fatal("limited outcome metric not found")
	}
}
