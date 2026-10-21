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
	Timeout  time.Duration `envconfig:"TIMEOUT" required:"true"`
	BaseURL  string        `envconfig:"BASEURL" required:"true"`
}

func NewConfig() (Config, error) {
	var config Config

	if err := envconfig.Process("IPINFO", &config); err != nil {
		return Config{}, fmt.Errorf("process IPInfo config: %w", err)
	}

	if config.CacheTTL < 0 {
		return Config{}, fmt.Errorf("IPInfo cache TTL cannot be negative")
	}

	if config.Timeout <= 0 {
		return Config{}, fmt.Errorf("IPInfo timeout must be greater than zero")
	}

	baseURL, err := url.Parse(config.BaseURL)
	if err != nil || (baseURL.Scheme != "http" && baseURL.Scheme != "https") || baseURL.Host == "" {
		return Config{}, fmt.Errorf("IPInfo BASEURL must be an absolute HTTP(S) URL")
	}

	return config, nil
}

func NewConfigMust() Config {
	config, err := NewConfig()
	if err != nil {
		panic(fmt.Errorf("get IPInfo connection pool config: %w", err))
	}

	return config
}
