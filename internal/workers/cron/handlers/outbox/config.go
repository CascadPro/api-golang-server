package cron_outbox_handler

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Interval time.Duration `envconfig:"INTERVAL" default:"3h"`
	Timeout  time.Duration `envconfig:"TIMEOUT" default:"30s"`
}

func NewConfig() (Config, error) {
	var config Config

	if err := envconfig.Process("WORKER_CRON_OUTBOX", &config); err != nil {
		return Config{}, fmt.Errorf("process outbox cron config: %w", err)
	}

	if config.Interval <= 0 {
		return Config{}, fmt.Errorf("outbox cron interval must be positive")
	}

	if config.Timeout <= 0 {
		return Config{}, fmt.Errorf("outbox cron timeout must be positive")
	}

	return config, nil
}

func NewConfigMust() Config {
	config, err := NewConfig()
	if err != nil {
		panic(fmt.Errorf("get outbox cron config: %w", err))
	}

	return config
}
