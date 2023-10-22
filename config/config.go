package config

import "os"

var (
	StorageType string
)

func init() {
	StorageType = getEnv("STORAGE_TYPE", "filesystem")
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
