package core_redis_pool

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Host     string `envconfig:"HOST" required:"true"`
	Port     string `envconfig:"PORT" default:"6379"`
	Password string `envconfig:"PASSWORD" required:"false"`
	Database int    `envconfig:"DB" required:"true"`

	PoolSize     int           `envconfig:"POOL_SIZE" required:"true"`
	MinIdleConns int           `envconfig:"MIN_IDLE_CONNS" required:"true"`
	Timeout      time.Duration `envconfig:"TIMEOUT" required:"true"`
}

func NewConfig() (Config, error) {
	var config Config

	if err := envconfig.Process("REDIS", &config); err != nil {
		return Config{}, fmt.Errorf("process envconfig: %w", err)
	}

	return config, nil
}

func NewConfigMust() Config {
	config, err := NewConfig()
	if err != nil {
		err = fmt.Errorf("get Redis Connection Pool config: %w", err)
		panic(err)
	}

	return config
}
