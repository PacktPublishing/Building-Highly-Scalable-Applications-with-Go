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

// MonitoringServerConfig loads HTTP server config where metrics
// and probes will be hosted.
func MonitoringServerConfig() config.HTTPServer {
	return config.HTTPServer{
		Port: asIntOrDef("MONITORING_HTTP_PORT", 8081),
	}
}

// TracingConfig loads tracing configuration from
// the environment
func TracingConfig() config.Tracing {
	return config.Tracing{
		ExportEndpointURL: asStringOrDef(
			"TRACING_EXPORTER_ENDPOINT",
			"http://localhost:4318/v1/traces",
		),
	}
}

// AccountServiceDBConfig loads the database
// information for the account service from the
// environment
func AccountServiceDBConfig() config.AccountServiceDB {
	return config.AccountServiceDB{
		DSN: asStringOrDef(
			"ACCOUNT_SERVICE_DSN",
			"postgresql://as:as@localhost:5432/as?sslmode=disable",
		),
	}
}

func asIntOrDef(key string, defaultVal int) int {
	valueStr := os.Getenv(key)
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultVal
}

func asStringOrDef(key string, defaultVal string) string {
	v := os.Getenv(key)
	if v == "" {
		return defaultVal
	}
	return v
}
