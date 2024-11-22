package config

// Tracing holds information regarding
// the tracing configuration
type Tracing struct {
	// ExportEndpoint represents the target endpoint (host and port)
	// the tracing exporter will connect to
	ExportEndpointURL string
}
