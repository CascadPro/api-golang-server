package core_s3_pool

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Pool interface {
	PutObject(ctx context.Context, key string, body []byte) error
	PutObjectStream(ctx context.Context, key string, body io.Reader, size int64) error

	GetObject(ctx context.Context, key string) ([]byte, error)
	GetObjectStream(ctx context.Context, key, byteRange string) (io.ReadCloser, int64, string, error)

	DeleteObject(ctx context.Context, key string) error

	OpTimeout() time.Duration
	StreamTimeout() time.Duration
}

type ConnectionPool struct {
	client        *aws.Config
	streamClient  *aws.Config
	bucket        string
	timeout       time.Duration
	streamTimeout time.Duration
}

func credentialsProvider(cfg Config) aws.CredentialsProviderFunc {
	return func(ctx context.Context) (aws.Credentials, error) {
		return aws.Credentials{
			AccessKeyID:     cfg.AccessKeyID,
			SecretAccessKey: cfg.SecretAccessKey,
			SessionToken:    cfg.SessionToken,
			Source:          "inline_config",
		}, nil
	}
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

	withEndpoint := config.WithBaseEndpoint(cfg.Endpoint)
	withRegion := config.WithRegion(cfg.Region)
	withCredentials := config.WithCredentialsProvider(credentialsProvider(cfg))

	awsCfg, err := config.LoadDefaultConfig(
		ctx,
		withEndpoint,
		withRegion,
		withCredentials,
		config.WithHTTPClient(&httpClient),
	)
	if err != nil {
		return nil, fmt.Errorf("load aws sdk config: %w", err)
	}

	streamHttpClient := http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 20,
			IdleConnTimeout:     90 * time.Second,
		},
		Timeout: 0,
	}
	if cfg.StreamTimeout > 0 {
		streamHttpClient.Timeout = cfg.StreamTimeout
	}

	streamAwsCfg, err := config.LoadDefaultConfig(
		ctx,
		withEndpoint,
		withRegion,
		withCredentials,
		config.WithHTTPClient(&streamHttpClient),
	)
	if err != nil {
		return nil, fmt.Errorf("load aws sdk streaming config: %w", err)
	}

	s3client := s3.NewFromConfig(awsCfg)
	_, err = s3client.HeadBucket(ctx, &s3.HeadBucketInput{Bucket: &cfg.Bucket})
	if err != nil {
		return nil, fmt.Errorf("s3 head bucket failed: %w", err)
	}

	return &ConnectionPool{
		client:        &awsCfg,
		streamClient:  &streamAwsCfg,
		bucket:        cfg.Bucket,
		timeout:       cfg.Timeout,
		streamTimeout: cfg.StreamTimeout,
	}, nil
}

func (p *ConnectionPool) OpTimeout() time.Duration {
	return p.timeout
}

func (p *ConnectionPool) StreamTimeout() time.Duration {
	return p.streamTimeout
}
