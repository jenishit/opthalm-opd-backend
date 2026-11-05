// Package minio adapts the MinIO Go SDK (S3-compatible) to port.ObjectStorage
// for storing clinical images (fundus photos, OCT/visual-field exports
// pulled off diagnostic devices) and consent signature images.
package minio

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Config struct {
	Endpoint   string
	AccessKey  string
	SecretKey  string
	UseSSL     bool
	Bucket     string
	PublicBase string
}

type Store struct {
	client     *minio.Client
	bucket     string
	publicBase string
}

// New connects to the MinIO server and ensures the configured bucket
// exists, creating it if this is a fresh deployment.
func New(ctx context.Context, cfg Config) (*Store, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("minio: connect: %w", err)
	}

	exists, err := client.BucketExists(ctx, cfg.Bucket)
	if err != nil {
		return nil, fmt.Errorf("minio: check bucket: %w", err)
	}
	if !exists {
		if err := client.MakeBucket(ctx, cfg.Bucket, minio.MakeBucketOptions{}); err != nil {
			return nil, fmt.Errorf("minio: create bucket: %w", err)
		}
	}

	return &Store{client: client, bucket: cfg.Bucket, publicBase: cfg.PublicBase}, nil
}

func (s *Store) Upload(ctx context.Context, key string, r io.Reader, size int64, contentType string) error {
	_, err := s.client.PutObject(ctx, s.bucket, key, r, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return fmt.Errorf("minio: put object %q: %w", key, err)
	}
	return nil
}

// PresignedURL returns a GET URL valid for expiry. If PublicBase was
// configured (e.g. a reverse proxy serving the bucket under a stable,
// cacheable path), that's used instead since it doesn't expire — presigned
// URLs are the fallback for deployments without one.
func (s *Store) PresignedURL(ctx context.Context, key string, expiry time.Duration) (string, error) {
	if s.publicBase != "" {
		return s.publicBase + "/" + key, nil
	}

	url, err := s.client.PresignedGetObject(ctx, s.bucket, key, expiry, nil)
	if err != nil {
		return "", fmt.Errorf("minio: presign %q: %w", key, err)
	}
	return url.String(), nil
}

func (s *Store) Delete(ctx context.Context, key string) error {
	if err := s.client.RemoveObject(ctx, s.bucket, key, minio.RemoveObjectOptions{}); err != nil {
		return fmt.Errorf("minio: remove object %q: %w", key, err)
	}
	return nil
}
