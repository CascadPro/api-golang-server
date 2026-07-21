package core_s3_pool

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Region          string
	Bucket          string
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	SessionToken    string

	Timeout time.Duration
}

func NewConfig() (Config, error) {
	var config Config

	if err := envconfig.Process("S3", &config); err != nil {
		return Config{}, fmt.Errorf("process envconfig: %w", err)
	}

	return config, nil
}

func NewConfigMust() Config {
	config, err := NewConfig()
	if err != nil {
		err = fmt.Errorf("get S3 AWS Connection Pool config: %w", err)
		panic(err)
	}

	return config
}
