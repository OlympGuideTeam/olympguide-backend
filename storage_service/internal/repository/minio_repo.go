package repository

import (
	"bytes"
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"log"
	"storage_service/internal/constants"
)

type ILogoRepository interface {
	Upload(ctx context.Context, objectName string, data []byte, contentType string, ext string) (string, error)
	GetFileExtension(ctx context.Context, objectName string) (string, error)
	Delete(ctx context.Context, objectName string) error
}

type MinioRepository struct {
	client *minio.Client
	bucket string
}

func NewMinioRepository(client *minio.Client) *MinioRepository {
	repo := &MinioRepository{
		client: client,
		bucket: "logos",
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
		p := `{
    "Version": "2012-10-17",
    "Statement": [{
        "Effect": "Allow",
        "Principal": {"AWS": ["*"]},
        "Action": ["s3:GetObject"],
        "Resource": ["arn:aws:s3:::%s/*"]
    }]
}`
		err = client.SetBucketPolicy(ctx, repo.bucket, p)
		if err != nil {
			log.Fatalf("Ошибка при установке публичной политики: %v", err)
		}
	}

	return repo
}

func (r *MinioRepository) Upload(ctx context.Context, objectName string, data []byte, contentType string, ext string) (string, error) {
	metaData := map[string]string{
		"file_extension": ext,
	}

	_, err := r.client.PutObject(ctx, r.bucket, objectName, bytes.NewReader(data), int64(len(data)), minio.PutObjectOptions{
		ContentType:  contentType,
		UserMetadata: metaData,
	})

	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("http://%s/%s/%s", r.client.CredContext().Endpoint, r.bucket, objectName)
	return url, nil
}

func (r *MinioRepository) GetFileExtension(ctx context.Context, objectName string) (string, error) {
	objInfo, err := r.client.StatObject(ctx, r.bucket, objectName, minio.StatObjectOptions{})
	if err != nil {
		return "", err
	}
	ext, ok := objInfo.UserMetadata["file_extension"]
	if !ok {
		return "", fmt.Errorf(constants.FailedGetExtensionErr)
	}
	return ext, nil
}

func (r *MinioRepository) Delete(ctx context.Context, objectName string) error {
	return r.client.RemoveObject(ctx, r.bucket, objectName, minio.RemoveObjectOptions{})
}
