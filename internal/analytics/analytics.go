package analytics

import (
	"github.com/codemicro/lgballtDiscordBot/internal/config"
	"github.com/codemicro/lgballtDiscordBot/internal/logging"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func init() {
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		logging.Info("Running Prometheus HTTP server at " + config.PrometheusAddress)
		err := http.ListenAndServe(config.PrometheusAddress, nil)
		if err != nil {
			logging.Error(err, "Failed to start Prometheus HTTP server")
		}
	}()
}

var (
	counterCommandsRun = promauto.NewCounterVec(prometheus.CounterOpts{
		Name:        "lgballtbot_commands_run",
		Help:        "Commands run",
	}, []string{"command"})

	counterPluralKitRequests = promauto.NewCounterVec(prometheus.CounterOpts{
		Name:        "lgballtbot_pluralkit_requests",
		Help:        "Requests made to the PluralKit API",
	}, []string{"request"})
)

func ReportCommandUse(commandName string) {
	go counterCommandsRun.WithLabelValues(commandName).Inc()
}

func ReportPluralKitRequest(requestName string) {
	go counterPluralKitRequests.WithLabelValues(requestName).Inc()
}