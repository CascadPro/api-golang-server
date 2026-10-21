package core_s3_pool

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
)

var (
	ErrNotFound     = errors.New("s3 object not found")
	ErrConflict     = errors.New("s3 conflict or already exists")
	ErrInvalidRange = errors.New("s3 invalid bytes range")
	ErrUnknown      = errors.New("unknown s3 error")
)

func MapErrors(err error) error {
	if err == nil {
		return nil
	}

	if nf, ok := errors.AsType[*s3types.NotFound](err); ok {
		return fmt.Errorf("%v: %w", nf.Message, ErrNotFound)
	}

	if alreadyOwned, ok := errors.AsType[*s3types.BucketAlreadyOwnedByYou](err); ok {
		return fmt.Errorf("%v: %w", alreadyOwned.Message, ErrConflict)
	}

	if apiErr, ok := errors.AsType[smithy.APIError](err); ok {
		if apiErr.ErrorCode() == "InvalidRange" || apiErr.ErrorCode() == "RequestedRangeNotSatisfiable" {
			return fmt.Errorf("%v: %w", err, ErrInvalidRange)
		}
	}

	return fmt.Errorf("s3 operation failed: %v: %w", err, ErrUnknown)
}

func (p *ConnectionPool) PutObject(ctx context.Context, key string, body []byte) error {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	_, err := s3.NewFromConfig(*p.client).PutObject(ctx, &s3.PutObjectInput{
		Bucket: &p.bucket, Key: &key, Body: bytes.NewReader(body),
	})

	return MapErrors(err)
}

func (p *ConnectionPool) GetObject(ctx context.Context, key string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	resp, err := s3.NewFromConfig(*p.client).GetObject(ctx, &s3.GetObjectInput{Bucket: &p.bucket, Key: &key})
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

type cancelReadCloser struct {
	io.ReadCloser
	cancel context.CancelFunc
}

func (r *cancelReadCloser) Close() error {
	err := r.ReadCloser.Close()
	r.cancel()
	return err
}

// GetObjectStream streams an S3 object (or a single byte range) without
// buffering the full object in the API process. The returned body owns the
// request timeout and must be closed by the caller.
func (p *ConnectionPool) GetObjectStream(ctx context.Context, key, byteRange string) (io.ReadCloser, int64, string, error) {
	requestCtx, cancel := context.WithTimeout(ctx, p.timeout)
	input := &s3.GetObjectInput{Bucket: &p.bucket, Key: &key}

	if byteRange != "" {
		input.Range = &byteRange
	}

	resp, err := s3.NewFromConfig(*p.client).GetObject(requestCtx, input)
	if err != nil {
		cancel()
		return nil, 0, "", MapErrors(err)
	}

	var size int64
	if resp.ContentLength != nil {
		size = *resp.ContentLength
	}

	var contentRange string
	if resp.ContentRange != nil {
		contentRange = *resp.ContentRange
	}

	return &cancelReadCloser{ReadCloser: resp.Body, cancel: cancel}, size, contentRange, nil
}

func (p *ConnectionPool) DeleteObject(ctx context.Context, key string) error {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	_, err := s3.NewFromConfig(*p.client).DeleteObject(ctx, &s3.DeleteObjectInput{Bucket: &p.bucket, Key: &key})

	return MapErrors(err)
}
