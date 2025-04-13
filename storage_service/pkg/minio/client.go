package minio

import (
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
	"storage_service/internal/config"
)

func NewMinioClient(config *config.Config) *minio.Client {
	endpoint := fmt.Sprintf("%s:%d", config.MinioHost, config.MinioPort)
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.MinioUser, config.MinioPassword, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalf("failed to init minio client: %v", err)
	}
	return client
}
