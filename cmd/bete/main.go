package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/yi-jiayu/datamall/v3"
)

var (
	httpIncomingRequestDurationSeconds = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_incoming_request_duration_seconds",
			Help: "Duration of incoming HTTP requests by path, status code and method.",
		},
		[]string{"path", "code", "method"},
	)
	httpOutgoingRequestDurationSeconds = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_outgoing_request_duration_seconds",
			Help: "Duration of outgoing HTTP requests by service and status code.",
		},
		[]string{"service", "code"},
	)
)

func newTelegramWebhookHandler(dm datamall.APIClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		stop := r.URL.Query().Get("stop")
		if stop == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		arrivals, err := dm.GetBusArrival(stop, "")
		if err != nil {
			log.Printf("error getting bus arrivals: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(arrivals)
		if err != nil {
			log.Printf("error writing bus arrivals to response: %v", err)
		}
	}
}

func main() {
	client := &http.Client{
		Transport: promhttp.InstrumentRoundTripperDuration(
			httpOutgoingRequestDurationSeconds.MustCurryWith(prometheus.Labels{"service": "datamall"}),
			http.DefaultTransport,
		),
	}
	accountKey := os.Getenv("DATAMALL_ACCOUNT_KEY")
	dm := datamall.NewClient(accountKey, client)
	http.Handle(
		"/telegram/updates",
		promhttp.InstrumentHandlerDuration(
			httpIncomingRequestDurationSeconds.MustCurryWith(prometheus.Labels{"path": "/telegram/updates"}),
			newTelegramWebhookHandler(dm),
		),
	)
	http.Handle("/metrics", promhttp.Handler())
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := fmt.Sprintf(":%s", port)
	log.Fatal(http.ListenAndServe(addr, nil))
}
