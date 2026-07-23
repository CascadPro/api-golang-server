package core_config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type SecretConfig struct {
	JwtSecretKey    string `envconfig:"JWT" required:"true"`
	UploadSecretKey string `envconfig:"UPLOAD" required:"true"`
}

func NewSecretConfig() (*SecretConfig, error) {
	var config SecretConfig

	if err := envconfig.Process("SECRET_KEY", &config); err != nil {
		return &SecretConfig{}, fmt.Errorf("process envconfig: %w", err)
	}

	return &config, nil
}

func NewSecretConfigMust() *SecretConfig {
	config, err := NewSecretConfig()
	if err != nil {
		err = fmt.Errorf("get core secret config: %w", err)
		panic(err)
	}

	return config
}
