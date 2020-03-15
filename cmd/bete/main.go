package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/yi-jiayu/datamall/v3"
	"github.com/yi-jiayu/ted"

	"github.com/yi-jiayu/bete"
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
		var update ted.Update
		err := json.NewDecoder(r.Body).Decode(&update)
		if err != nil {
			log.Printf("error decoding update: %v", err)
			return
		}
		if update.Message == nil || update.Message.Text == "" {
			return
		}
		message := update.Message
		var query string
		if command, args := message.CommandAndArgs(); command == "eta" {
			query = args
		} else {
			query = message.Text
		}
		parts := strings.Fields(query)
		if len(parts) == 0 {
			return
		}
		stop := parts[0]
		t := time.Now()
		arrivals, err := dm.GetBusArrival(stop, "")
		if err != nil {
			log.Printf("error getting bus arrivals: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		text, err := bete.FormatArrivalsByService(bete.ArrivalInfo{
			Stop:     bete.BusStop{ID: stop},
			Time:     t,
			Services: arrivals.Services,
			Filter:   parts[1:],
		})
		w.Header().Set("Content-Type", "application/json")
		reply := map[string]interface{}{
			"method":     "sendMessage",
			"chat_id":    message.Chat.ID,
			"text":       text,
			"parse_mode": "HTML",
		}
		err = json.NewEncoder(w).Encode(reply)
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
	log.Printf("listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
