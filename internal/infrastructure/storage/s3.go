package storage

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/config"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/internal/domain"
)

type s3Storage struct {
	client  *s3.Client
	presign *s3.PresignClient
	bucket  string
	baseURL string
}

// NewS3Storage creates a new StorageProvider backed by AWS S3.
func NewS3Storage(cfg config.StorageConfig) (domain.StorageProvider, error) {
	optFns := []func(*awsconfig.LoadOptions) error{
		awsconfig.WithRegion(cfg.S3Region),
	}

	// Use static credentials if provided, otherwise fall back to the
	// default credential chain (IAM role, env vars, shared config).
	if cfg.S3AccessKey != "" && cfg.S3SecretKey != "" {
		optFns = append(optFns, awsconfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(cfg.S3AccessKey, cfg.S3SecretKey, ""),
		))
	}

	awsCfg, err := awsconfig.LoadDefaultConfig(context.Background(), optFns...)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	var s3OptFns []func(*s3.Options)

	if cfg.S3Endpoint != "" {
		s3OptFns = append(s3OptFns, func(o *s3.Options) {
			o.BaseEndpoint = aws.String(cfg.S3Endpoint)
			o.UsePathStyle = cfg.S3ForcePathStyle
		})
	}

	client := s3.NewFromConfig(awsCfg, s3OptFns...)
	presign := s3.NewPresignClient(client)

	baseURL := cfg.BaseURL
	if baseURL == "" {
		if cfg.S3Endpoint != "" {
			baseURL = fmt.Sprintf("%s/%s", cfg.S3Endpoint, cfg.S3Bucket)
		} else {
			baseURL = fmt.Sprintf("https://%s.s3.%s.amazonaws.com", cfg.S3Bucket, cfg.S3Region)
		}
	}

	return &s3Storage{
		client:  client,
		presign: presign,
		bucket:  cfg.S3Bucket,
		baseURL: baseURL,
	}, nil
}

func (s *s3Storage) Upload(ctx context.Context, input domain.UploadInput) (*domain.UploadOutput, error) {
	contentType := input.ContentType
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	putInput := &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(input.Key),
		Body:        input.Body,
		ContentType: aws.String(contentType),
	}

	if input.Size > 0 {
		putInput.ContentLength = aws.Int64(input.Size)
	}

	if _, err := s.client.PutObject(ctx, putInput); err != nil {
		return nil, fmt.Errorf("failed to upload object %q: %w", input.Key, err)
	}

	return &domain.UploadOutput{
		Key: input.Key,
		URL: s.GetURL(input.Key),
	}, nil
}

func (s *s3Storage) Download(ctx context.Context, key string) (io.ReadCloser, error) {
	output, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to download object %q: %w", key, err)
	}

	return output.Body, nil
}

func (s *s3Storage) Delete(ctx context.Context, key string) error {
	if _, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}); err != nil {
		return fmt.Errorf("failed to delete object %q: %w", key, err)
	}

	return nil
}

func (s *s3Storage) GetURL(key string) string {
	return fmt.Sprintf("%s/%s", s.baseURL, key)
}

func (s *s3Storage) GeneratePresignedURL(ctx context.Context, key string, expiry time.Duration) (string, error) {
	req, err := s.presign.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(expiry))
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned download URL for %q: %w", key, err)
	}

	return req.URL, nil
}

func (s *s3Storage) GeneratePresignedUploadURL(ctx context.Context, key string, contentType string, expiry time.Duration) (string, error) {
	input := &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}
	if contentType != "" {
		input.ContentType = aws.String(contentType)
	}

	req, err := s.presign.PresignPutObject(ctx, input, s3.WithPresignExpires(expiry))
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned upload URL for %q: %w", key, err)
	}

	return req.URL, nil
}
