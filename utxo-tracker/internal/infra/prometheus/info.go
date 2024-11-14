package prometheus

import (
	"github.com/hannesdejager/utxo-tracker/internal/domain"
	"github.com/prometheus/client_golang/prometheus"
)

// newBuildInfoGauge sets up a prometheus gauge with a static value of 1
// when registered, this ensures the "service_meta_info" metric is always exposed.
func newBuildInfoGauge(info domain.ServiceInstance,
) prometheus.Gauge {
	v := info.Version
	buildInfo := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "service_meta_info",
		Help: "Service build and instance information",
		ConstLabels: map[string]string{
			"name":           info.Name,
			"go_version":     v.GoVersion,
			"build_date":     v.BuildDate,
			"commit_hash":    v.CommitLongHash,
			"committer":      v.Committer,
			"commit_subject": v.CommitSubject,
			"commit_date":    v.CommitDate,
			"start_time":     info.StartupTime.String(),
		},
	})
	buildInfo.Set(1)
	return buildInfo
}
