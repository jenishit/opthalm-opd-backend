package port

import (
	"context"
	"io"
	"time"
)

// ObjectStorage abstracts the S3-compatible object store (MinIO in this
// deployment) that clinical attachment and consent-signature images are
// kept in, so services depend on this interface rather than the minio-go
// client directly.
type ObjectStorage interface {
	// Upload streams size bytes from r to key and returns nothing — the
	// caller already knows the key it asked for. Overwrites any existing
	// object at key.
	Upload(ctx context.Context, key string, r io.Reader, size int64, contentType string) error
	// PresignedURL returns a time-limited URL the caller can hand to a
	// browser/client to fetch the object directly from the store, without
	// routing the (potentially large) file bytes through this API.
	PresignedURL(ctx context.Context, key string, expiry time.Duration) (string, error)
	// Delete removes the object at key. Not an error if it doesn't exist.
	Delete(ctx context.Context, key string) error
}
