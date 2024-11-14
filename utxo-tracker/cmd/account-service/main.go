package main

import (
	"time"

	"github.com/hannesdejager/utxo-tracker/internal/infra/api/restv1"
	"github.com/hannesdejager/utxo-tracker/internal/infra/env"
	"github.com/hannesdejager/utxo-tracker/internal/infra/httpsvr"
	"github.com/hannesdejager/utxo-tracker/internal/infra/logging"
	"github.com/hannesdejager/utxo-tracker/internal/infra/prometheus"
	"github.com/hannesdejager/utxo-tracker/internal/infra/sys"
)

func main() {
	inf := logging.InstanceInfo(time.Now())
	log := logging.NewLogger(inf)
	log.Info("Starting up...",
		"pid", inf.PID,
		"built", inf.Version.BuildDate,
		"commited", inf.Version.CommitDate,
		"committer", inf.Version.Committer,
		"subject", inf.Version.CommitSubject,
	)
	svr := httpsvr.StartAsync(
		env.HTTPConfig(),
		restv1.NewHandler("/rest/v1"),
	)
	_ = httpsvr.StartAsync(env.MonitoringServerConfig(),
		prometheus.NewHandler(inf))

	sys.AwaitTermination()
	log.Info("Shutting down...")
	httpsvr.StopGracefully(svr, 30*time.Second)
	log.Info("Bye!")
}
