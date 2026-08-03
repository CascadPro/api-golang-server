package core_s3_pool

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Pool interface {
	PutObject(ctx context.Context, key string, body []byte) error
	GetObject(ctx context.Context, key string) ([]byte, error)
	DeleteObject(ctx context.Context, key string) error

	OpTimeout() time.Duration
}

type ConnectionPool struct {
	client  *aws.Config
	bucket  string
	timeout time.Duration
}

func New(ctx context.Context, cfg Config) (*ConnectionPool, error) {
	httpClient := http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 20,
			IdleConnTimeout:     90 * time.Second,
		},
		Timeout: 30 * time.Second,
	}

	awsCfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithBaseEndpoint(cfg.Endpoint),
		config.WithRegion(cfg.Region),
		config.WithCredentialsProvider(aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
			return aws.Credentials{
				AccessKeyID:     cfg.AccessKeyID,
				SecretAccessKey: cfg.SecretAccessKey,
				SessionToken:    cfg.SessionToken,
				Source:          "inline_config",
			}, nil
		})),
		config.WithHTTPClient(&httpClient),
	)
	if err != nil {
		return nil, fmt.Errorf("load aws sdk config: %w", err)
	}

	s3client := s3.NewFromConfig(awsCfg)
	_, err = s3client.HeadBucket(ctx, &s3.HeadBucketInput{Bucket: &cfg.Bucket})
	if err != nil {
		return nil, fmt.Errorf("s3 head bucket failed: %w", err)
	}

	return &ConnectionPool{
		client:  &awsCfg,
		bucket:  cfg.Bucket,
		timeout: cfg.Timeout,
	}, nil
}

func (p *ConnectionPool) OpTimeout() time.Duration {
	return p.timeout
}
