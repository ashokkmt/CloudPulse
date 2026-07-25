package storage

import (
	"context"
	"io"
	"log/slog"
	"net/url"
	"os"
	"strings"
	"time"

	"cloudpulse/backend/internal/metrics"

	"net/http"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type StorageService struct {
	client     *minio.Client
	bucketName string
	cdnURL     string
}

func NewStorageService() *StorageService {
	endpoint := os.Getenv("SPACES_ENDPOINT")
	accessKeyID := os.Getenv("SPACES_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("SPACES_SECRET_ACCESS_KEY")
	bucketName := os.Getenv("SPACES_BUCKET")
	cdnURL := os.Getenv("CDN_URL")

	if endpoint == "" {
		endpoint = "s3.localhost:9000" // Use s3.localhost to ensure signature matches browser requests
	}
	if bucketName == "" {
		bucketName = "cloudpulse"
	}

	// For MinIO locally, useSSL is typically false. DigitalOcean Spaces uses true.
	useSSL := true
	if strings.Contains(endpoint, "localhost") || strings.Contains(endpoint, "minio") {
		useSSL = false
	}

	transport := otelhttp.NewTransport(http.DefaultTransport)
	creds := credentials.NewStaticV4(accessKeyID, secretAccessKey, "")

	client, err := minio.New(endpoint, &minio.Options{
		Creds:     creds,
		Secure:    useSSL,
		Transport: transport,
	})
	if err != nil {
		slog.Error("Error initializing storage client", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// Create bucket if it doesn't exist
	ctx := context.Background()
	exists, err := client.BucketExists(ctx, bucketName)
	if err == nil && !exists {
		err = client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: "us-east-1"})
		if err != nil {
			slog.Error("Failed to create bucket", slog.String("bucket", bucketName), slog.String("error", err.Error()))
		} else {
			slog.Info("Bucket created successfully", slog.String("bucket", bucketName))
		}
	}

	return &StorageService{
		client:     client,
		bucketName: bucketName,
		cdnURL:     cdnURL,
	}
}

func (s *StorageService) BucketName() string {
	return s.bucketName
}

func (s *StorageService) UploadFile(ctx context.Context, objectName string, reader io.Reader, objectSize int64, contentType string) (string, error) {
	start := time.Now()
	defer func() { metrics.UploadLatency.Observe(time.Since(start).Seconds()) }()

	_, err := s.client.PutObject(ctx, s.bucketName, objectName, reader, objectSize, minio.PutObjectOptions{
		ContentType: contentType,
	})

	if err != nil {
		metrics.UploadsTotal.WithLabelValues("failure").Inc()
		slog.Error("Spaces Error", slog.String("operation", "upload"), slog.String("object", objectName), slog.String("error", err.Error()))
		return "", err
	}

	metrics.UploadsTotal.WithLabelValues("success").Inc()
	slog.Info("File Upload", slog.String("object", objectName), slog.Int64("size", objectSize))
	return objectName, nil
}

func (s *StorageService) DeleteFile(ctx context.Context, objectName string) error {
	err := s.client.RemoveObject(ctx, s.bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		slog.Error("Spaces Error", slog.String("operation", "delete"), slog.String("object", objectName), slog.String("error", err.Error()))
	} else {
		slog.Info("File Delete", slog.String("object", objectName))
	}
	return err
}

func (s *StorageService) GeneratePresignedURL(ctx context.Context, objectName string) (string, error) {
	reqParams := make(url.Values)

	presignedURL, err := s.client.PresignedGetObject(ctx, s.bucketName, objectName, time.Minute*15, reqParams)
	if err != nil {
		return "", err
	}

	return presignedURL.String(), nil
}

func (s *StorageService) GetFile(ctx context.Context, objectName string) (io.ReadCloser, error) {
	object, err := s.client.GetObject(ctx, s.bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	return object, nil
}
