package internal

import (
	"context"
	"log"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/lifecycle"
)

const (
	DAY = time.Hour * 24
)

func CreateMinioBucket(ctx context.Context, bucketName string, minioClient *minio.Client) {
	// checking if bucket already exists
	bucket_exists, err := minioClient.BucketExists(ctx, bucketName)
	failOnError(err, "failed to check minio bucket existence")
	if bucket_exists {
		log.Println("minio bucket already exists; skipping")
		return
	}

	// creating bucket
	err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
	failOnError(err, "failed to create minio bucket")
	log.Println("created minio bucket")
}

func SetMinioExpiration(ctx context.Context, bucketName string, resultsExpire time.Duration, minioClient *minio.Client) {
	// creating expiration
	var expirationInDays int
	if resultsExpire < DAY {
		expirationInDays = 1
	} else {
		expirationInDays = int(resultsExpire.Hours()) / 24
	}

	config := lifecycle.NewConfiguration()
	config.Rules = []lifecycle.Rule{
		{
			ID:     "expire-bucket",
			Status: "Enabled",
			Expiration: lifecycle.Expiration{
				Days: lifecycle.ExpirationDays(expirationInDays),
			},
		},
	}

	err := minioClient.SetBucketLifecycle(ctx, bucketName, config)
	failOnError(err, "failed to set bucket policy")
	log.Println("set minio expiration")
}
