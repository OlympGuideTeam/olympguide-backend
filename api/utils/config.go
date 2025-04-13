package utils

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	DBHost             string
	DBPort             int
	DBUser             string
	DBPassword         string
	DBName             string
	RedisHost          string
	RedisPort          int
	RedisPassword      string
	ServerPort         int
	StorageServiceHost string
	StorageServicePort int
}

func LoadConfig() (*Config, error) {
	dbPortStr := os.Getenv("DB_PORT")
	serverPortStr := os.Getenv("API_PORT")
	redisPortStr := os.Getenv("REDIS_PORT")
	storageServicePortStr := os.Getenv("STORAGE_SERVICE_PORT")

	dbPort, err := strconv.Atoi(dbPortStr)
	if err != nil {
		return nil, fmt.Errorf("invalid DB_PORT: %w", err)
	}

	serverPort, err := strconv.Atoi(serverPortStr)
	if err != nil {
		return nil, fmt.Errorf("invalid SERVER_PORT: %w", err)
	}

	redisPort, err := strconv.Atoi(redisPortStr)
	if err != nil {
		return nil, fmt.Errorf("invalid REDIS_PORT: %w", err)
	}

	storageServicePort, err := strconv.Atoi(storageServicePortStr)
	if err != nil {
		return nil, fmt.Errorf("invalid STORAGE_SERVICE_PORT: %w", err)
	}

	cfg := &Config{
		DBHost:             os.Getenv("DB_HOST"),
		DBPort:             dbPort,
		DBUser:             os.Getenv("DB_USER"),
		DBPassword:         os.Getenv("DB_PASSWORD"),
		DBName:             os.Getenv("DB_NAME"),
		RedisHost:          os.Getenv("REDIS_HOST"),
		RedisPort:          redisPort,
		RedisPassword:      os.Getenv("REDIS_PASSWORD"),
		ServerPort:         serverPort,
		StorageServicePort: storageServicePort,
		StorageServiceHost: os.Getenv("STORAGE_SERVICE_HOST"),
	}
	return cfg, nil
}
