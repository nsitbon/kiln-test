package monitor

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func HandlerDuration(inner http.Handler, name string) http.Handler {
	requestDuration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: fmt.Sprintf("http_request_%s_seconds", name),
		Help: fmt.Sprintf("Duration of all %s HTTP requests", name),
	}, []string{"code", "method"})
	prometheus.MustRegister(requestDuration)
	return promhttp.InstrumentHandlerDuration(requestDuration, inner)
}
