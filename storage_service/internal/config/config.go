package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	StorageServicePort int
	MinioPort          int
	MinioHost          string
	MinioUser          string
	MinioPassword      string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("Не удалось загрузить .env файл, будут использованы переменные окружения")
	}

	storageServicePortStr := os.Getenv("STORAGE_SERVICE_PORT")
	storageServicePort, err := strconv.Atoi(storageServicePortStr)
	if err != nil {
		return nil, fmt.Errorf("invalid STORAGE_SERVICE_PORT: %w", err)
	}

	minioPortStr := os.Getenv("MINIO_PORT")
	minioPort, err := strconv.Atoi(minioPortStr)
	if err != nil {
		return nil, fmt.Errorf("invalid MINIO_API_PORT: %w", err)
	}

	return &Config{
		StorageServicePort: storageServicePort,
		MinioUser:          os.Getenv("MINIO_USER"),
		MinioPassword:      os.Getenv("MINIO_PASSWORD"),
		MinioHost:          os.Getenv("MINIO_HOST"),
		MinioPort:          minioPort,
	}, nil
}
