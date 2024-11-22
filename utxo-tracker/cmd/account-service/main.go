package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/hannesdejager/utxo-tracker/internal/infra/api/restv1"
	"github.com/hannesdejager/utxo-tracker/internal/infra/env"
	"github.com/hannesdejager/utxo-tracker/internal/infra/httpsvr"
	"github.com/hannesdejager/utxo-tracker/internal/infra/jaeger"
	"github.com/hannesdejager/utxo-tracker/internal/infra/logging"
	"github.com/hannesdejager/utxo-tracker/internal/infra/prometheus"
	"github.com/hannesdejager/utxo-tracker/internal/infra/sys"
)

func main() {
	inf := logging.InstanceInfo(time.Now())
	log := logging.NewLogger(inf)
	slog.SetDefault(log)
	log.Info("Starting up...",
		"pid", inf.PID,
		"built", inf.Version.BuildDate,
		"commited", inf.Version.CommitDate,
		"committer", inf.Version.Committer,
		"subject", inf.Version.CommitSubject,
	)

	tp, err := jaeger.InitTracing(env.TracingConfig(), inf)
	if err != nil {
		log.Error("Failed to initialize tracing", "error", err)
		os.Exit(1)
	}
	defer func() { _ = tp.Shutdown(context.Background()) }()

	svr := httpsvr.StartAsync(
		env.HTTPConfig(),
		restv1.NewHandler(log, "/rest/v1"),
	)
	_ = httpsvr.StartAsync(env.MonitoringServerConfig(),
		prometheus.NewHandler(inf))

	sys.AwaitTermination()
	log.Info("Shutting down...")
	httpsvr.StopGracefully(svr, 30*time.Second)
	log.Info("Bye!")
}
