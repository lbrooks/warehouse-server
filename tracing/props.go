package tracing

import (
	"main/inventory"
	"os"
)

func getEndpoint() string {
	return inventory.GetEnvOrDefault("TRACER_ENDPOINT", "http://localhost:14268/api/traces")
}

func getService() string {
	return inventory.GetEnvOrDefault("TRACER_SERVICE", "inventory-api")
}

func getHost() string {
	name, err := os.Hostname()
	if err != nil {
		name = "unknown"
	}
	return inventory.GetEnvOrDefault("TRACER_HOST", name)
}
