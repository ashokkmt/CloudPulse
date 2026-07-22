package storage

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/url"
	"os"
	"time"

	"cloudpulse/backend/internal/metrics"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"net/http"
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
		endpoint = "localhost:9000" // default minio local
	}
	if bucketName == "" {
		bucketName = "cloudpulse"
	}

	// For MinIO locally, useSSL is typically false. DigitalOcean Spaces uses true.
	// We infer based on endpoint containing localhost
	useSSL := true
	if endpoint == "localhost:9000" || endpoint == "minio:9000" {
		useSSL = false
	}
	
	transport := otelhttp.NewTransport(http.DefaultTransport)

	client, err := minio.New(endpoint, &minio.Options{
		Creds:     credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
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
			// Set policy for public read (for attachments and thumbnails)
			policy := fmt.Sprintf(`{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"AWS":["*"]},"Action":["s3:GetObject"],"Resource":["arn:aws:s3:::%s/*"]}]}`, bucketName)
			_ = client.SetBucketPolicy(ctx, bucketName, policy)
		}
	}

	return &StorageService{
		client:     client,
		bucketName: bucketName,
		cdnURL:     cdnURL,
	}
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
	return s.GetObjectURL(objectName), nil
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

func (s *StorageService) GetObjectURL(objectName string) string {
	if s.cdnURL != "" {
		// e.g. http://localhost:9000/cloudpulse/attachment.jpg
		u, _ := url.Parse(s.cdnURL)
		u.Path = u.Path + "/" + objectName
		return u.String()
	}
	// Fallback to direct S3 URL
	scheme := "https"
	if s.client.EndpointURL().Scheme != "" {
		scheme = s.client.EndpointURL().Scheme
	} else if s.client.EndpointURL().Host == "minio:9000" || s.client.EndpointURL().Host == "localhost:9000" {
		scheme = "http"
	}
	return fmt.Sprintf("%s://%s/%s/%s", scheme, s.client.EndpointURL().Host, s.bucketName, objectName)
}

func (s *StorageService) GetFile(ctx context.Context, objectName string) (io.ReadCloser, error) {
	object, err := s.client.GetObject(ctx, s.bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	return object, nil
}
