package logging

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/hannesdejager/utxo-tracker/internal/domain"
	"github.com/hannesdejager/utxo-tracker/internal/infra/linker"
	"go.opentelemetry.io/otel/trace"
)

func InstanceInfo(startup time.Time) domain.ServiceInstance {
	return domain.ServiceInstance{
		Name:        filepath.Base(os.Args[0]),
		Version:     linker.ServiceVersion(),
		PID:         os.Getpid(),
		StartupTime: startup,
	}
}

func NewLogger(i domain.ServiceInstance) *slog.Logger {
	opts := &slog.HandlerOptions{
		AddSource: false,
		Level:     slog.LevelDebug,
		ReplaceAttr: func(groups []string,
			a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				a.Value = slog.StringValue(
					a.Value.Time().
						In(time.UTC).
						Format(time.RFC3339),
				)
			}
			return a
		},
	}
	log := slog.New(
		CorrelationHandler{
			slog.NewJSONHandler(os.Stdout, opts),
		},
	).With(
		"service", i.Name,
		"commit", i.Version.CommitShortHash,
	)
	return log
}

type CorrelationHandler struct {
	slog.Handler
}

func (h CorrelationHandler) Handle(
	ctx context.Context, r slog.Record) error {
	s := trace.SpanFromContext(ctx).SpanContext()
	if s.HasTraceID() {
		r.Add("trace_id", s.TraceID().String())
		if s.HasSpanID() {
			r.Add("span_id", s.SpanID().String())
		}
	}

	return h.Handler.Handle(ctx, r)
}

func (h CorrelationHandler) WithAttrs(
	attrs []slog.Attr) slog.Handler {
	return CorrelationHandler{h.Handler.WithAttrs(attrs)}
}

func (h CorrelationHandler) WithGroup(
	name string) slog.Handler {
	return CorrelationHandler{h.Handler.WithGroup(name)}
}
