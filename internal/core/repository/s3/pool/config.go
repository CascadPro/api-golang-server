package core_s3_pool

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Region          string `envconfig:"S3_REGION" required:"true"`
	Bucket          string `envconfig:"S3_BUCKET" required:"true"`
	Endpoint        string `envconfig:"S3_ENDPOINT" required:"true"`
	AccessKeyID     string `envconfig:"S3_ACCESS_KEY_ID" required:"true"`
	SecretAccessKey string `envconfig:"S3_SECRET_ACCESS_KEY" required:"true"`
	SessionToken    string `envconfig:"S3_SESSION_TOKEN" required:"false"`

	Timeout time.Duration `envconfig:"S3_TIMEOUT" required:"true"`
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
