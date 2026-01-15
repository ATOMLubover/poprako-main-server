package oss

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type r2Client struct {
	client        *s3.Client
	presignClient *s3.PresignClient

	bucketName   string
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

	customDomain := os.Getenv("R2_CUSTOM_DEMAIN")
	if customDomain == "" {
		customDomain = accountID + ".r2.cloudflarestorage.com"
	}

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
		customDomain:  customDomain,
	}
}

func (r2 *r2Client) PresignPut(ossKey string) (string, error) {
	const exp = 10 * 60 * time.Minute // 10 minutes

	req, err := r2.presignClient.PresignPutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(r2.bucketName),
		Key:    aws.String(ossKey),
	}, s3.WithPresignExpires(exp))
	if err != nil {
		return "", fmt.Errorf("failed to presign put object: %w", err)
	}

	return req.URL, nil
}

func (r2 *r2Client) PresignGet(ossKey string) (string, error) {
	if r2.customDomain != "" {
		return fmt.Sprintf("https://%s/%s", r2.customDomain, ossKey), nil
	}

	// TODO: Implement presigned URL generation if needed
	return "", errors.New("PresignGet not implemented")
}
