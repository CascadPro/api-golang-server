package cron_idempotency_handler

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Interval time.Duration `envconfig:"INTERVAL" default:"1h"`
	Timeout  time.Duration `envconfig:"TIMEOUT" default:"30s"`
}

func NewConfig() (Config, error) {
	var config Config

	if err := envconfig.Process("WORKER_CRON_IDEMPOTENCY", &config); err != nil {
		return Config{}, fmt.Errorf("process idempotency cron config: %w", err)
	}

	if config.Interval <= 0 {
		return Config{}, fmt.Errorf("idempotency cron interval must be positive")
	}

	if config.Timeout <= 0 {
		return Config{}, fmt.Errorf("idempotency cron timeout must be positive")
	}

	return config, nil
}

func NewConfigMust() Config {
	config, err := NewConfig()
	if err != nil {
		panic(fmt.Errorf("get idempotency cron config: %w", err))
	}

	return config
}
