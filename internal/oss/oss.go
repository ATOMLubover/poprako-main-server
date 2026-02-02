package oss

import "context"

type OSSClient interface {
	PresignPut(ossKey string) (string, error)
	PresignGet(ossKey string) (string, error)
	DeleteObject(ctx context.Context, ossKey string) error
}
