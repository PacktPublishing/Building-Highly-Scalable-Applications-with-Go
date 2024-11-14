package prometheus

import (
	"net/http"

	"github.com/hannesdejager/utxo-tracker/internal/domain"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewHandler(info domain.ServiceInstance) http.Handler {
	reg := prometheus.NewRegistry()
	reg.MustRegister(
		newBuildInfoGauge(info),
		apiRequestsTotal,
		apiErrorsTotal,
		apiRequestDuration,
	)
	return promhttp.HandlerFor(
		reg,
		promhttp.HandlerOpts{Registry: reg},
	)
}
