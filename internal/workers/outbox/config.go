package worker_outbox

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Interval  time.Duration `envconfig:"INTERVAL"   default:"5s"`
	BatchSize int           `envconfig:"BATCH_SIZE" default:"25"`
	LockTTL   time.Duration `envconfig:"LOCK_TTL"   default:"30s"`
}

func NewConfig() (Config, error) {
	var config Config

	if err := envconfig.Process("WORKER_OUTBOX", &config); err != nil {
		return Config{}, fmt.Errorf("process envconfig: %w", err)
	}

	return config, nil
}

func NewConfigMust() Config {
	config, err := NewConfig()
	if err != nil {
		err = fmt.Errorf("get Worker Outbox config: %w", err)
		panic(err)
	}

	return config
}
