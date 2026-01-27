package services

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/jeancarloshp/calorieai/internal/domain"
	"go.opentelemetry.io/otel"
)

type S3Service struct {
	client     *s3.Client
	bucketName string
	logger     domain.Logger
}

func NewS3Service(cfg *domain.Config, logger domain.Logger) (*S3Service, error) {
	awsConfig, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(cfg.AWSS3Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.AWSAccessKeyID,
			cfg.AWSSecretAccessKey,
			"",
		)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := s3.NewFromConfig(awsConfig)

	return &S3Service{
		client:     client,
		bucketName: cfg.AWSS3BucketName,
		logger:     logger,
	}, nil
}

func (s *S3Service) UploadImage(ctx context.Context, fileContent io.Reader, fileType string) (string, error) {
	tr := otel.Tracer("services/s3_service.go")
	ctx, span := tr.Start(ctx, "UploadImage")
	defer span.End()

	fileID := uuid.New().String()
	ext := getExtensionFromMimeType(ctx, fileType)
	fileName := fmt.Sprintf("food-images/%s-%d%s", fileID, time.Now().Unix(), ext)

	buf := new(bytes.Buffer)
	size, err := buf.ReadFrom(fileContent)
	if err != nil {
		s.logger.Error("failed to read file content", "error", err)
		return "", fmt.Errorf("failed to read file content: %w", err)
	}

	_, err = s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(s.bucketName),
		Key:           aws.String(fileName),
		Body:          bytes.NewReader(buf.Bytes()),
		ContentLength: aws.Int64(size),
		ContentType:   aws.String(fileType),
	})
	if err != nil {
		s.logger.Error("failed to upload to S3", "error", err, "fileName", fileName)
		return "", fmt.Errorf("failed to upload to S3: %w", err)
	}

	url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.bucketName, s.client.Options().Region, fileName)
	s.logger.Info("image uploaded successfully", "url", url)

	return url, nil
}

func getExtensionFromMimeType(ctx context.Context, mimeType string) string {
	tr := otel.Tracer("services/s3_service.go")
	ctx, span := tr.Start(ctx, "getExtensionFromMimeType")
	defer span.End()

	extensions := map[string]string{
		"image/jpeg": ".jpg",
		"image/jpg":  ".jpg",
		"image/png":  ".png",
		"image/webp": ".webp",
		"image/heic": ".heic",
	}

	if ext, ok := extensions[mimeType]; ok {
		return ext
	}
	return filepath.Ext(mimeType)
}

func (s *S3Service) GenerateUploadPresignedURL(ctx context.Context, firebaseUID string, contentType string, fileSize int64) (string, string, error) {
	tr := otel.Tracer("services/s3_service.go")
	ctx, span := tr.Start(ctx, "GenerateUploadPresignedURL")
	defer span.End()

	const (
		maxFileSize     = 5 * 1024 * 1024
		presignDuration = 5 * time.Minute
	)

	if fileSize > maxFileSize {
		return "", "", fmt.Errorf("file size exceeds maximum allowed size of %d bytes", maxFileSize)
	}

	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
		"image/webp": true,
	}

	contentType = strings.ToLower(strings.TrimSpace(contentType))
	if !allowedTypes[contentType] {
		return "", "", fmt.Errorf("invalid content type: %s. Allowed types: jpeg, jpg, png, webp", contentType)
	}

	ext := getExtensionFromMimeType(ctx, contentType)
	if ext == "" {
		ext = ".jpg"
	}

	fileName := fmt.Sprintf("avatars/%s/avatar%s", firebaseUID, ext)

	presignClient := s3.NewPresignClient(s.client)

	putObjectInput := &s3.PutObjectInput{
		Bucket:        aws.String(s.bucketName),
		Key:           aws.String(fileName),
		ContentType:   aws.String(contentType),
		ContentLength: aws.Int64(fileSize),
	}

	presignedReq, err := presignClient.PresignPutObject(ctx, putObjectInput, func(opts *s3.PresignOptions) {
		opts.Expires = presignDuration
	})
	if err != nil {
		s.logger.Error("failed to generate presigned URL", "error", err, "userID", firebaseUID)
		return "", "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	avatarPath := fmt.Sprintf("/%s", fileName)

	s.logger.Info("presigned URL generated successfully",
		"userID", firebaseUID,
		"fileName", fileName,
		"expiresIn", presignDuration.String(),
	)

	return presignedReq.URL, avatarPath, nil
}

func (s *S3Service) GenerateFoodImageUploadPresignedURL(ctx context.Context, userID string, contentType string, fileSize int64) (string, string, error) {
	tr := otel.Tracer("services/s3_service.go")
	ctx, span := tr.Start(ctx, "GenerateFoodImageUploadPresignedURL")
	defer span.End()

	const (
		maxFileSize     = 5 * 1024 * 1024
		presignDuration = 5 * time.Minute
	)

	if fileSize > maxFileSize {
		return "", "", fmt.Errorf("file size exceeds maximum allowed size of %d bytes", maxFileSize)
	}

	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
		"image/webp": true,
	}

	contentType = strings.ToLower(strings.TrimSpace(contentType))
	if !allowedTypes[contentType] {
		return "", "", fmt.Errorf("invalid content type: %s. Allowed types: jpeg, jpg, png, webp", contentType)
	}

	ext := getExtensionFromMimeType(ctx, contentType)
	if ext == "" {
		ext = ".jpg"
	}

	imageID := uuid.New().String()
	fileName := fmt.Sprintf("users/%s/food_images/%s%s", userID, imageID, ext)

	presignClient := s3.NewPresignClient(s.client)

	putObjectInput := &s3.PutObjectInput{
		Bucket:        aws.String(s.bucketName),
		Key:           aws.String(fileName),
		ContentType:   aws.String(contentType),
		ContentLength: aws.Int64(fileSize),
	}

	presignedReq, err := presignClient.PresignPutObject(ctx, putObjectInput, func(opts *s3.PresignOptions) {
		opts.Expires = presignDuration
	})
	if err != nil {
		s.logger.Error("failed to generate presigned URL", "error", err, "userID", userID)
		return "", "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	imagePath := fmt.Sprintf("/%s", fileName)

	s.logger.Info("presigned URL generated successfully",
		"userID", userID,
		"fileName", fileName,
		"expiresIn", presignDuration.String(),
	)

	return presignedReq.URL, imagePath, nil
}

func (s *S3Service) DownloadImage(ctx context.Context, imagePath string) ([]byte, error) {
	tr := otel.Tracer("services/s3_service.go")
	ctx, span := tr.Start(ctx, "DownloadImage")
	defer span.End()

	// Always remove leading slash if present
	imagePath = strings.TrimPrefix(imagePath, "/")

	getObjectInput := &s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(imagePath),
	}

	result, err := s.client.GetObject(ctx, getObjectInput)
	if err != nil {
		s.logger.Error("failed to download image from S3", "error", err, "imagePath", imagePath)
		return nil, fmt.Errorf("failed to download image: %w", err)
	}
	defer result.Body.Close()

	imageBytes, err := io.ReadAll(result.Body)
	if err != nil {
		s.logger.Error("failed to read image data", "error", err, "imagePath", imagePath)
		return nil, fmt.Errorf("failed to read image data: %w", err)
	}

	s.logger.Info("image downloaded successfully", "imagePath", imagePath, "size", len(imageBytes))

	return imageBytes, nil
}
