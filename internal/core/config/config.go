package core_config

import (
	"fmt"
	"os"
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
	TimeZone       *time.Location
	Connection     Connection `envconfig:"CONNECTION" default:"online"`
	EnvMode        EnvMode    `envconfig:"ENV_MODE" default:"dev"`
	AllowedOrigins string     `envconfig:"ALLOWED_ORIGINS" default:""`
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
