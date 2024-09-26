package env

import (
	"os"
	"strconv"
	"time"

	"github.com/hannesdejager/utxo-tracker/internal/app/config"
)

// HTTPConfig loads HTTP server config.
func HTTPConfig() config.HTTPServer {
	idle := asIntOrDef("HTTP_IDLE_TIMEOUT", 30)
	grace := asIntOrDef("HTTP_SHUTDOWN_GRACE_PERIOD", 30)
	return config.HTTPServer{
		Port:                asIntOrDef("HTTP_PORT", 8080),
		IdleTimeout:         time.Duration(idle) * time.Second,
		ShutdownGracePeriod: time.Duration(grace) * time.Second,
	}
}

func asIntOrDef(key string, defaultVal int) int {
	valueStr := os.Getenv(key)
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultVal
}
