package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/getsentry/sentry-go"
	sentryhttp "github.com/getsentry/sentry-go/http"

	_ "github.com/lib/pq"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/yi-jiayu/datamall/v3"
	"github.com/yi-jiayu/ted"

	"github.com/yi-jiayu/bete"
)

var (
	commit      = "master"
	environment = os.Getenv("ENVIRONMENT")
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

func newTelegramWebhookHandler(b bete.Bete) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var update ted.Update
		var buf bytes.Buffer
		tee := io.TeeReader(r.Body, &buf)
		err := json.NewDecoder(tee).Decode(&update)
		if err != nil {
			log.Printf("error decoding update: %v", err)
			return
		}
		log.Println(buf.String())
		b.HandleUpdate(r.Context(), update)
	}
}

func init() {
	if environment == "" {
		environment = "development"
	}
}

func main() {
	if err := sentry.Init(sentry.ClientOptions{
		Release:     commit,
		Environment: environment,
		BeforeSend: func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
			if hint.OriginalException != nil {
				log.Printf("%+v", hint.OriginalException)
			}
			return event
		},
	}); err != nil {
		log.Printf("Sentry initialization failed: %v\n", err)
	}
	accountKey := os.Getenv("DATAMALL_ACCOUNT_KEY")
	if accountKey == "" {
		log.Fatal("DATAMALL_ACCOUNT_KEY environment variable not set")
	}
	dm := datamall.NewClient(accountKey, &http.Client{
		Transport: promhttp.InstrumentRoundTripperDuration(
			httpOutgoingRequestDurationSeconds.MustCurryWith(prometheus.Labels{"service": "datamall"}),
			http.DefaultTransport,
		),
	})
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable not set")
	}
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatalf("error opening postgres database: %v", err)
	}
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN environment variable not set")
	}
	bot := ted.Bot{
		Token: botToken,
		HTTPClient: &http.Client{
			Transport: promhttp.InstrumentRoundTripperDuration(
				httpOutgoingRequestDurationSeconds.MustCurryWith(prometheus.Labels{"service": "telegram"}),
				http.DefaultTransport,
			),
		},
	}
	streetViewStaticApiKey := os.Getenv("STREET_VIEW_STATIC_API_KEY")
	if streetViewStaticApiKey == "" {
		log.Println("warning: STREET_VIEW_STATIC_API_KEY environment variable not set")
	}
	b := bete.Bete{
		Version:                commit,
		StreetViewStaticAPIKey: streetViewStaticApiKey,
		Clock:                  bete.RealClock{},
		BusStops:               bete.SQLBusStopRepository{DB: db},
		Favourites:             bete.SQLFavouriteRepository{DB: db},
		DataMall:               dm,
		Telegram:               bot,
	}
	sentryHandler := sentryhttp.New(sentryhttp.Options{Repanic: true})
	http.Handle(
		"/telegram/updates",
		sentryHandler.Handle(
			promhttp.InstrumentHandlerDuration(
				httpIncomingRequestDurationSeconds.MustCurryWith(prometheus.Labels{"path": "/telegram/updates"}),
				newTelegramWebhookHandler(b),
			),
		),
	)
	http.Handle("/metrics", promhttp.Handler())
	var host string
	if environment == "development" {
		host = "localhost"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := fmt.Sprintf("%s:%s", host, port)
	log.Printf("listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
