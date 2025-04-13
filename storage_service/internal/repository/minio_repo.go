package repository

import (
	"bytes"
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"log"
)

type ILogoRepository interface {
	Upload(ctx context.Context, objectName string, data []byte, contentType string) (string, error)
	Delete(ctx context.Context, objectName string) error
}

type MinioRepository struct {
	client         *minio.Client
	bucket         string
	publicEndpoint string
}

func NewMinioRepository(client *minio.Client, publicEndpoint string) *MinioRepository {
	repo := &MinioRepository{
		client:         client,
		bucket:         "logos",
		publicEndpoint: publicEndpoint,
	}

	ctx := context.Background()
	exists, err := client.BucketExists(ctx, repo.bucket)
	if err != nil {
		log.Fatalf("Ошибка при проверке существования бакета: %v", err)
	}

	if !exists {
		err = client.MakeBucket(ctx, repo.bucket, minio.MakeBucketOptions{})
		if err != nil {
			log.Fatalf("Ошибка при создании бакета: %v", err)
		}

		policy := fmt.Sprintf(`{
			"Version": "2012-10-17",
			"Statement": [{
				"Effect": "Allow",
				"Principal": {"AWS": ["*"]},
				"Action": ["s3:GetObject"],
				"Resource": ["arn:aws:s3:::%s/*"]
			}]
		}`, repo.bucket)

		err = client.SetBucketPolicy(ctx, repo.bucket, policy)
		if err != nil {
			log.Fatalf("Ошибка при установке политики: %v", err)
		}
	}

	return repo
}

func (r *MinioRepository) Upload(ctx context.Context, objectName string, data []byte, contentType string) (string, error) {
	_, err := r.client.PutObject(ctx, r.bucket, objectName, bytes.NewReader(data), int64(len(data)), minio.PutObjectOptions{
		ContentType: contentType,
	})

	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("%s/%s/%s", r.publicEndpoint, r.bucket, objectName)
	return url, nil
}

func (r *MinioRepository) Delete(ctx context.Context, objectName string) error {
	fmt.Println(objectName)
	return r.client.RemoveObject(ctx, r.bucket, objectName, minio.RemoveObjectOptions{})
}
