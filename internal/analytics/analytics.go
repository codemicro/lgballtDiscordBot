package analytics

import (
	"errors"
	"github.com/codemicro/lgballtDiscordBot/internal/config"
	"github.com/codemicro/lgballtDiscordBot/internal/state"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
	"net/http"
)

func Start(state *state.State) error {

	http.Handle("/metrics", promhttp.Handler())

	server := &http.Server{
		Addr: config.PrometheusAddress,
	}

	go func() {
		log.Info().Msg("Running Prometheus HTTP server at " + config.PrometheusAddress)
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error().Err(err).Msg("Failed to start Prometheus HTTP server")
		}
	}()

	go func() {
		state.WaitUntilShutdownTrigger()
		_ = server.Close()
		state.FinishGoroutine()
	}()

	return nil
}

var (
	counterCommandsRun = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "lgballtbot_commands_run",
		Help: "Commands run",
	}, []string{"command"})

	counterPluralKitRequests = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "lgballtbot_pluralkit_requests",
		Help: "Requests made to the PluralKit API",
	}, []string{"request"})
)

func ReportCommandUse(commandName string) {
	go counterCommandsRun.WithLabelValues(commandName).Inc()
}

func ReportPluralKitRequest(requestName string) {
	go counterPluralKitRequests.WithLabelValues(requestName).Inc()
}
