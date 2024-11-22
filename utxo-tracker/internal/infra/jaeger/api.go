package jaeger

import (
	"net/http"

	"github.com/hannesdejager/utxo-tracker/pkg/gons/middleware"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
)

// TracingMiddleware is an HTTP middleware that creates and manages spans for incoming requests.
func TracingMiddleware(next http.Handler) http.Handler {
	tracer := otel.Tracer("account-service")
	propagator := otel.GetTextMapPropagator()

	fn := func(w http.ResponseWriter, r *http.Request) {
		// Extract trace context from the incoming request
		ctx := propagator.Extract(
			r.Context(),
			propagation.HeaderCarrier(r.Header),
		)
		peeker := &middleware.ResponsePeeker{ResponseWriter: w}

		// Create a new span for the request
		ctx, span := tracer.Start(ctx, r.Method+" "+r.URL.Path)
		// Defer the logic for panic recovery, capturing the status code and ending the span
		defer func() {
			// Capture HTTP status code in the span
			span.SetAttributes(attribute.Int(
				"http.status_code", peeker.Status))
			// Mark the span as errored if the status code is in the 4xx or 5xx range
			if peeker.Status >= 400 {
				span.SetStatus(codes.Error,
					"HTTP error code returned")
			}

			span.End()
		}()

		// Add common attributes to the span
		span.SetAttributes(
			attribute.String("http.method", r.Method),
			attribute.String("http.url", r.URL.String()),
			attribute.String("http.user_agent", r.UserAgent()),
		)

		// Pass the updated context to the next handler
		r = r.WithContext(ctx)
		next.ServeHTTP(peeker, r)
	}
	return http.HandlerFunc(fn)
}
