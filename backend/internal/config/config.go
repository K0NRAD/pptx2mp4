package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port             string
	StoragePath      string
	MaxFileSize      int64
	CleanupInterval  time.Duration
	AllowedOrigins   []string
	BasePath         string
	LogLevel         string
	LogFormat        string
}

func LoadConfig() *Config {
	return &Config{
		Port:             getEnv("PORT", "8080"),
		StoragePath:      getEnv("STORAGE_PATH", "./storage"),
		MaxFileSize:      getEnvAsInt64("MAX_FILE_SIZE", 100*1024*1024),
		CleanupInterval:  getEnvAsDuration("CLEANUP_INTERVAL", time.Hour),
		AllowedOrigins:   getEnvAsSlice("ALLOWED_ORIGINS", []string{"http://localhost:3000"}),
		BasePath:         getEnv("BASE_PATH", "/pptx2mp4"),
		LogLevel:         getEnv("LOG_LEVEL", "info"),
		LogFormat:        getEnv("LOG_FORMAT", "json"),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvAsInt64(key string, defaultValue int64) int64 {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.ParseInt(valueStr, 10, 64)
	if err != nil {
		return defaultValue
	}

	return value
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := time.ParseDuration(valueStr)
	if err != nil {
		return defaultValue
	}

	return value
}

func getEnvAsSlice(key string, defaultValue []string) []string {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	return []string{valueStr}
}
