package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"interview-test/configs"
	"interview-test/internal/client/infura"
	"interview-test/internal/handler"
	"interview-test/internal/util/monitor"
)

func main() {
	conf := configs.NewConfigFromEnv()
	server := getHTTPServer(conf.ListenAddress, getConfiguredRouter(getGasPriceHandler(conf)))
	serverStopped, receivedSignal := make(chan error, 1), make(chan os.Signal, 1)
	signal.Notify(receivedSignal, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("starting server listening on addr '%s'\n", server.Addr)
		serverStopped <- server.ListenAndServe()
	}()

	select {
	case <-receivedSignal:
		log.Println("received termination signal: stopping server")
		_ = server.Close()
	case err := <-serverStopped:
		log.Printf("HTTP server error: %s\n", err.Error())
	}

	log.Println("shutting down")
}

func getGasPriceHandler(conf configs.Config) http.Handler {
	return monitor.HandlerDuration(handler.NewGetGasPriceHandler(getInfuraClient(conf)), "get_gas_price")
}

func getHTTPServer(listenAddress string, getGasPriceHandler http.Handler) *http.Server {
	return &http.Server{Addr: listenAddress, Handler: getConfiguredRouter(getGasPriceHandler)}
}

func getInfuraClient(c configs.Config) infura.Client {
	infuraClient, err := infura.NewClient(c.ApiKey, c.ProjectID, &http.Client{Timeout: c.InfuraClientTimeout})

	if err != nil {
		log.Fatal(err)
	}

	return infura.NewCachingClient(infuraClient, c.InfuraClientCacheRefreshInterval, c.InfuraClientCacheMaxDurationBeforeEviction)
}

func getConfiguredRouter(getGasPriceHandler http.Handler) http.Handler {
	router := http.NewServeMux()
	router.Handle("/metrics", promhttp.Handler())
	router.Handle("/alive", handlerStatusOK())
	router.Handle("/ready", handlerStatusOK())
	router.Handle("/eth/gasPrice", getGasPriceHandler)
	return router
}

func handlerStatusOK() http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, _ *http.Request) { rw.WriteHeader(http.StatusOK) })
}
