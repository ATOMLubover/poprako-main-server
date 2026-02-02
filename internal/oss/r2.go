package oss

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type r2Client struct {
	client        *s3.Client
	presignClient *s3.PresignClient

	bucketName string
	accountID  string
	customDomain string
}

func NewR2Client() OSSClient {
	accountID := os.Getenv("R2_ACCOUNT_ID")
	if accountID == "" {
		panic("R2_ACCOUNT_ID environment variable is not set")
	}

	accessKeyID := os.Getenv("R2_ACCESS_KEY_ID")
	if accessKeyID == "" {
		panic("R2_ACCESS_KEY_ID environment variable is not set")
	}

	secretKeyID := os.Getenv("R2_SECRET_ACCESS_KEY")
	if secretKeyID == "" {
		panic("R2_SECRET_ACCESS_KEY environment variable is not set")
	}

	region := os.Getenv("R2_REGION")
	if region == "" {
		region = "auto"
	}

	bucketName := os.Getenv("R2_BUCKET_NAME")
	if bucketName == "" {
		panic("R2_BUCKET_NAME environment variable is not set")
	}

	customDomain := os.Getenv("R2_CUSTOM_DOMAIN")

	// Use account ID based endpoint directly
	r2Endpoint := fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountID)

	// Build AWS SDK config with custom endpoint and static credentials
	ctx := context.TODO()

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyID, secretKeyID, "")),
	)
	if err != nil {
		panic(fmt.Sprintf("unable to load SDK config, %v", err))
	}

	// Create R2 service client
	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(r2Endpoint)
	})

	presignClient := s3.NewPresignClient(client)

	return &r2Client{
		client:        client,
		presignClient: presignClient,
		bucketName:    bucketName,
		accountID:     accountID,
		customDomain:  customDomain,
	}
}

func (r2 *r2Client) PresignPut(ossKey string) (string, error) {
	const exp = 10 * time.Minute // 10 minutes

	// Auto-detect content type for common image extensions and include it
	// in the signed request. Clients must use the same Content-Type when
	// uploading or the signature will not match.
	contentType := detectImageContentType(ossKey)

	input := &s3.PutObjectInput{
		Bucket: aws.String(r2.bucketName),
		Key:    aws.String(ossKey),
	}
	if contentType != "" {
		input.ContentType = aws.String(contentType)
	}

	// Debug: Log input details
	// fmt.Printf("[DEBUG] PresignPut input: %+v\n", input)

	req, err := r2.presignClient.PresignPutObject(context.TODO(), input, s3.WithPresignExpires(exp))
	if err != nil {
		return "", fmt.Errorf("failed to presign put object: %w", err)
	}

	// Debug: Log generated presigned URL
	// fmt.Printf("[DEBUG] Generated presigned URL: %s\n", req.URL)

	return req.URL, nil
}

// detectImageContentType returns a reasonable Content-Type for common
// image file extensions. Returns empty string if unknown (caller may omit).
func detectImageContentType(key string) string {
	ext := strings.ToLower(filepath.Ext(key))
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	case ".svg":
		return "image/svg+xml"
	case ".avif":
		return "image/avif"
	case ".bmp":
		return "image/bmp"
	case ".tif", ".tiff":
		return "image/tiff"
	default:
		return ""
	}
}

func (r2 *r2Client) PresignGet(ossKey string) (string, error) {
	if r2.customDomain != "" {
		// Use custom domain if available
		return fmt.Sprintf("https://%s/%s", r2.customDomain, ossKey), nil
	}

	// Use bucket.accountID format to match the PUT URL format
	encodedKey := url.PathEscape(ossKey)
	
	return fmt.Sprintf("https://%s.%s.r2.cloudflarestorage.com/%s", r2.bucketName, r2.accountID, encodedKey), nil
}

func (r2 *r2Client) DeleteObject(ctx context.Context, ossKey string) error {
	const maxRetries = 3
	const retryDelay = 500 * time.Millisecond // 500ms linear backoff

	var lastErr error
	for attempt := 1; attempt <= maxRetries; attempt++ {
		input := &s3.DeleteObjectInput{
			Bucket: aws.String(r2.bucketName),
			Key:    aws.String(ossKey),
		}

		_, err := r2.client.DeleteObject(ctx, input)
		if err == nil {
			return nil // Success
		}

		// Check if object doesn't exist (NoSuchKey) - treat as success
		var nsk *types.NoSuchKey
		if errors.As(err, &nsk) {
			return nil // Object already deleted, idempotent success
		}

		// For other errors, prepare for retry
		lastErr = err

		// If not the last attempt, wait before retrying
		if attempt < maxRetries {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(retryDelay):
				// Continue to next retry
			}
		}
	}

	return fmt.Errorf("failed to delete object after %d attempts: %w", maxRetries, lastErr)
}
