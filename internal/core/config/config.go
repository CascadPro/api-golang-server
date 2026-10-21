package core_config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type EnvMode string
type Connection string

const (
	EnvModeDev        = EnvMode("dev")
	EnvModeProd       = EnvMode("prod")
	ConnectionOnline  = Connection("online")
	ConnectionOffline = Connection("offline")
)

type Config struct {
	TimeZone             *time.Location
	Connection           Connection    `envconfig:"CONNECTION" default:"online"`
	EnvMode              EnvMode       `envconfig:"ENV_MODE" default:"dev"`
	AllowedOrigins       string        `envconfig:"ALLOWED_ORIGINS" default:""`
	ShutdownTimeout      time.Duration `envconfig:"SHUTDOWN_TIMEOUT" default:"1min"`
	TrustedProxies       []string      `envconfig:"TRUSTED_PROXIES" default:""`
	RateLimitFailureMode string        `envconfig:"RATE_LIMIT_FAILURE_MODE" default:"closed"`
}

func NewConfig() (*Config, error) {
	tz := os.Getenv("CORE_TIME_ZONE")
	if tz == "" {
		tz = time.UTC.String()
	}

	zone, err := time.LoadLocation(tz)
	if err != nil {
		return nil, fmt.Errorf("load time zone: %s: %w", tz, err)
	}

	var config Config
	if err := envconfig.Process("CORE", &config); err != nil {
		return &Config{}, fmt.Errorf("process envconfig: %w", err)
	}

	config.TimeZone = zone
	for i, proxy := range config.TrustedProxies {
		config.TrustedProxies[i] = strings.TrimSpace(proxy)
	}

	if config.RateLimitFailureMode == "" {
		config.RateLimitFailureMode = "closed"
	}

	if config.RateLimitFailureMode != "open" && config.RateLimitFailureMode != "closed" {
		return nil, fmt.Errorf("CORE_RATE_LIMIT_FAILURE_MODE must be open or closed")
	}

	return &config, nil
}

func NewConfigMust() *Config {
	config, err := NewConfig()
	if err != nil {
		err = fmt.Errorf("get core config: %w", err)
		panic(err)
	}

	return config
}
