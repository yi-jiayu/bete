package bete

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	telegramUpdatesTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "telegram_updates_total",
			Help: "The total number of Telegram updates received, partitioned by type.",
		},
		[]string{"type"},
	)
	callbackQueriesTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "bete_callback_queries_total",
			Help: "The total number of callback queries received, partitioned by type.",
		},
		[]string{"type"},
	)
	commandsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "bete_commands_total",
			Help: "The total number of commands received, partitioned by command.",
		},
		[]string{"command"},
	)
)
