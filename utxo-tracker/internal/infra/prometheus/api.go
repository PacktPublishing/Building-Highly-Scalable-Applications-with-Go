package prometheus

import (
	"fmt"
	"net/http"
	"time"

	"github.com/hannesdejager/utxo-tracker/pkg/gons/middleware"
	"github.com/prometheus/client_golang/prometheus"
)

var apiRequestsTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "api_requests_total",
		Help: "Number of requests to the API",
	},
	[]string{"path", "method", "status_class"},
)

var apiErrorsTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "api_errors_total",
		Help: "Number of API errors",
	},
	[]string{"path", "method", "status_class"},
)

var apiRequestDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "api_request_duration_seconds",
		Help:    "Duration of API requests in seconds",
		Buckets: prometheus.DefBuckets,
	},
	[]string{"path", "method", "status_class"},
)

func APIMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(
		w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &middleware.ResponsePeeker{ResponseWriter: w}

		next.ServeHTTP(rec, r)

		statusClass := fmt.Sprintf(
			"%dxx", rec.Status/100,
		)
		labels := prometheus.Labels{
			"path":         r.URL.Path,
			"method":       r.Method,
			"status_class": statusClass,
		}

		apiRequestsTotal.With(labels).Inc()

		if rec.Status >= 400 {
			apiErrorsTotal.With(labels).Inc()
		}

		duration := time.Since(start).Seconds()
		apiRequestDuration.With(labels).Observe(duration)
	})
}
