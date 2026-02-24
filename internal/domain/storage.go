package domain

import (
	"context"
	"io"
	"time"
)

// UploadInput holds the parameters for uploading a file to storage.
type UploadInput struct {
	// Key is the object path/name within the bucket (e.g., "products/uuid/image.jpg").
	Key string

	// Body is the readable stream of file content.
	Body io.Reader

	// ContentType is the MIME type of the file (e.g., "image/jpeg").
	// If empty, the implementation should default to "application/octet-stream".
	ContentType string

	// Size is the content length in bytes. Pass 0 or -1 if unknown.
	Size int64
}

// UploadOutput holds the result of a successful upload.
type UploadOutput struct {
	// Key is the final object key in storage.
	Key string

	// URL is the accessible URL for the uploaded object.
	URL string
}

// ObjectInfo holds metadata about a stored object, returned by HeadObject.
type ObjectInfo struct {
	Key           string
	ContentType   string
	ContentLength int64
}

// StorageProvider defines the interface for object storage operations.
// Implementations may target AWS S3, GCS, Azure Blob Storage, or a local filesystem.
type StorageProvider interface {
	// Upload stores a file in the storage backend and returns its location.
	Upload(ctx context.Context, input UploadInput) (*UploadOutput, error)

	// Download retrieves a file from storage as a readable stream.
	// The caller is responsible for closing the returned ReadCloser.
	Download(ctx context.Context, key string) (io.ReadCloser, error)

	// Delete removes a file from storage.
	Delete(ctx context.Context, key string) error

	// GetURL returns the URL for an object given its key.
	GetURL(key string) string

	// GeneratePresignedURL creates a time-limited signed URL for downloading an object.
	GeneratePresignedURL(ctx context.Context, key string, expiry time.Duration) (string, error)

	// GeneratePresignedUploadURL creates a time-limited signed URL for uploading
	// directly from a client (browser/mobile) without proxying through the server.
	GeneratePresignedUploadURL(ctx context.Context, key string, contentType string, expiry time.Duration) (string, error)

	// HeadObject retrieves metadata for an object without downloading its content.
	// Returns nil, nil if the object does not exist.
	HeadObject(ctx context.Context, key string) (*ObjectInfo, error)
}
