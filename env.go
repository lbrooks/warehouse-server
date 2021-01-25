package server

import "os"

// GetEnvOrDefault wrapper to get an environmental variable or use default
func GetEnvOrDefault(key, defaultValue string) string {
	endpoint := os.Getenv(key)
	if endpoint == "" {
		return defaultValue
	}
	return endpoint
}
