package bete

import (
	"context"
	"encoding/json"
	"time"

	"github.com/pkg/errors"
	"github.com/yi-jiayu/datamall/v3"
	"github.com/yi-jiayu/ted"
)

type Clock interface {
	Now() time.Time
}

type RealClock struct{}

func (c RealClock) Now() time.Time {
	return time.Now()
}

type DataMall interface {
	GetBusArrival(busStopCode string, serviceNo string) (datamall.BusArrival, error)
}

type Telegram interface {
	Do(request ted.Request) (ted.Response, error)
}

type Bete struct {
	Clock    Clock
	BusStops BusStopRepository
	DataMall DataMall
	Telegram Telegram
}

func (b Bete) HandleUpdate(ctx context.Context, u ted.Update) {
	switch {
	case u.Message != nil:
		b.HandleMessage(ctx, u.Message)
	case u.CallbackQuery != nil:
		b.HandleCallbackQuery(ctx, u.CallbackQuery)
	}
}

func (b Bete) etaMessageText(ctx context.Context, stopID string, filter []string) (string, error) {
	t := b.Clock.Now()
	arrivals, err := b.DataMall.GetBusArrival(stopID, "")
	if err != nil {
		return "", errors.Wrap(err, "error getting bus arrivals")
	}
	var stop BusStop
	stop, err = b.BusStops.Find(stopID)
	if err != nil {
		captureError(ctx, err)
		stop = BusStop{ID: stopID}
	}
	return FormatArrivalsByService(ArrivalInfo{
		Stop:     stop,
		Time:     t,
		Services: arrivals.Services,
		Filter:   filter,
	})
}

type CallbackData struct {
	Type   string   `json:"t"`
	StopID string   `json:"b,omitempty"`
	Filter []string `json:"s,omitempty"`
	Format string   `json:"f,omitempty"`
}

func (c CallbackData) Encode() string {
	JSON, err := json.Marshal(c)
	if err != nil {
		panic(err)
	}
	return string(JSON)
}

func etaMessageReplyMarkup(stopID string, filter []string) ted.InlineKeyboardMarkup {
	return ted.InlineKeyboardMarkup{
		InlineKeyboard: [][]ted.InlineKeyboardButton{
			{
				{
					Text: "Refresh",
					CallbackData: CallbackData{
						Type:   "refresh",
						StopID: stopID,
						Filter: filter,
					}.Encode(),
				},
			},
		},
	}
}

func etaMessageReplyMarkupP(stopID string, filter []string) *ted.InlineKeyboardMarkup {
	markup := etaMessageReplyMarkup(stopID, filter)
	return &markup
}
