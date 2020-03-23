package bete

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	telegramUpdates = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "telegram_updates_total",
			Help: "The total number of Telegram updates received, partitioned by type.",
		},
		[]string{"type"},
	)
)
