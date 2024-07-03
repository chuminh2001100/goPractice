package miniovt

import (
	"context"
	"errors"
	"fmt"

	"github.com/minio/minio-go/v7"
)

var _ Service = (*minioService)(nil)

type minioService struct {
	mClient *minio.Client
}

func New(mc *minio.Client) Service {
	return &minioService{
		mClient: mc,
	}
}

type Service interface {
	CreateBucket(ctx context.Context, bucketName string) error
}

func (m minioService) CreateBucket(ctx context.Context, bucketName string) error {
	exists, errBucketExists := m.mClient.BucketExists(ctx, bucketName)
	if errBucketExists == nil && exists {
		return errors.New("bucket already exists")
	} else if errBucketExists != nil {
		fmt.Println("Check bucket fail")
		return errBucketExists
	}
	return m.mClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
}
