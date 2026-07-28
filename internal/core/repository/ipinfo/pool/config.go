package core_ipinfo_pool

import (
	"fmt"
	"net/url"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Token    string        `envconfig:"TOKEN" required:"true"`
	CacheTTL time.Duration `envconfig:"CACHE_TTL" default:"24h"`
	Timeout  time.Duration `enconfig:"TIMEOUT" required:"true"`
	BaseURL  url.URL       `enconfig:"BASEURL" required:"true"`
}

func NewConfig() (Config, error) {
	var config Config

	if err := envconfig.Process("IPINFO", &config); err != nil {
		return Config{}, fmt.Errorf("process envconfig: %w", err)
	}

	return config, nil
}

func NewConfigMust() Config {
	config, err := NewConfig()
	if err != nil {
		err = fmt.Errorf("get IP Info Connection Pool config: %w", err)
		panic(err)
	}

	return config
}
