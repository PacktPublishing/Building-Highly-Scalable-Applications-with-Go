package logging

import (
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/hannesdejager/utxo-tracker/internal/domain"
	"github.com/hannesdejager/utxo-tracker/internal/infra/linker"
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
		AddSource: true,
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
		slog.NewJSONHandler(os.Stdout, opts),
	).With(
		"service", i.Name,
		"commit", i.Version.CommitShortHash,
	)
	return log
}
