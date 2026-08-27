package core_s3_pool

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
)

var (
	ErrNotFound = errors.New("s3 object not found")
	ErrConflict = errors.New("s3 conflict or already exists")
	ErrUnknown  = errors.New("unknown s3 error")
)

func MapErrors(err error) error {
	if err == nil {
		return nil
	}

	var nf *s3types.NotFound
	if errors.As(err, &nf) {
		return fmt.Errorf("%v: %w", nf.Message, ErrNotFound)
	}

	var ardOwU *s3types.BucketAlreadyOwnedByYou
	if errors.As(err, &ardOwU) {
		return fmt.Errorf("%v: %w", ardOwU.Message, ErrConflict)
	}

	return fmt.Errorf("s3 operation failed: %v: %w", err, ErrUnknown)
}

func (p *ConnectionPool) PutObject(ctx context.Context, key string, body []byte) error {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	_, err := s3.NewFromConfig(*p.client).PutObject(ctx, &s3.PutObjectInput{
		Bucket: &p.bucket,
		Key:    &key,
		Body:   bytes.NewReader(body),
	})

	return MapErrors(err)
}

func (p *ConnectionPool) GetObject(ctx context.Context, key string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	resp, err := s3.NewFromConfig(*p.client).GetObject(ctx, &s3.GetObjectInput{
		Bucket: &p.bucket,
		Key:    &key,
	})
	if err != nil {
		return nil, MapErrors(err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read s3 response body: %w", err)
	}

	return data, nil
}

func (p *ConnectionPool) DeleteObject(ctx context.Context, key string) error {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	_, err := s3.NewFromConfig(*p.client).DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: &p.bucket,
		Key:    &key,
	})

	return MapErrors(err)
}
