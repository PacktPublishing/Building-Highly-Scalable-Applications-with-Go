package config

import "time"

// HTTPServer holds settings for HTTP server configuration.
type HTTPServer struct {
	Port                int
	IdleTimeout         time.Duration
	ReadTimeout         time.Duration
	WriteTimeout        time.Duration
	ReadHeaderTimeout   time.Duration
	ShutdownGracePeriod time.Duration
}
