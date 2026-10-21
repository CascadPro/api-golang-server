package core_rate_limit

import (
	"errors"
	"testing"
	"time"
)

func TestPolicyValidate(t *testing.T) {
	tests := []struct {
		name    string
		policy  Policy
		wantErr bool
	}{
		{name: "valid", policy: Policy{Name: "test", Limit: 10, Window: time.Minute}},
		{name: "zero limit", policy: Policy{Name: "test", Limit: 0, Window: time.Minute}, wantErr: true},
		{name: "negative limit", policy: Policy{Name: "test", Limit: -1, Window: time.Minute}, wantErr: true},
		{name: "zero window", policy: Policy{Name: "test", Limit: 10, Window: 0}, wantErr: true},
		{name: "sub-millisecond window", policy: Policy{Name: "test", Limit: 10, Window: time.Microsecond}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.policy.Validate()
			if tt.wantErr && !errors.Is(err, ErrInvalidPolicy) {
				t.Fatalf("expected ErrInvalidPolicy, got %v", err)
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
