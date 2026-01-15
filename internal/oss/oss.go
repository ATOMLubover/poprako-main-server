package oss

type OSSClient interface {
	PresignPut(ossKey string) (string, error)
	PresignGet(ossKey string) (string, error)
}
